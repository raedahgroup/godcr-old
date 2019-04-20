package main

import (
	"context"
	"fmt"
	"os"
	"os/signal"
	"strings"
	"sync"
	"syscall"

	"github.com/jessevdk/go-flags"
	"github.com/raedahgroup/godcr/app"
	"github.com/raedahgroup/godcr/app/config"
	"github.com/raedahgroup/godcr/app/help"
	"github.com/raedahgroup/godcr/app/walletmediums/dcrwalletrpc"
	"github.com/raedahgroup/godcr/cli"
	"github.com/raedahgroup/godcr/cli/commands"
	"github.com/raedahgroup/godcr/cli/runner"
	"github.com/raedahgroup/godcr/cli/walletloader"
	"github.com/raedahgroup/godcr/fyne"
	"github.com/raedahgroup/godcr/nuklear"
	"github.com/raedahgroup/godcr/terminal"
	"github.com/raedahgroup/godcr/web"
)

// triggered after program execution is complete or if interrupt signal is received
var beginShutdown = make(chan bool)

// shutdownOps holds cleanup/shutdown functions that should be executed when shutdown signal is triggered
var shutdownOps []func()

// opError stores any error that occurs while performing an operation
var opError error

func main() {
	appConfig, args, err := config.LoadConfig()
	if err != nil {
		fmt.Fprintln(os.Stderr, err.Error())
		os.Exit(1)
	}

	// Initialize log rotation.  After log rotation has been initialized, the
	// logger variables may be used.
	initLogRotator(config.LogFile)
	defer func() {
		if logRotator != nil {
			logRotator.Close()
		}
	}()

	// Special show command to list supported subsystems and exit.
	if appConfig.DebugLevel == "show" {
		fmt.Println("Supported subsystems", supportedSubsystems())
		os.Exit(0)
	}

	// Parse, validate, and set debug log level(s).
	if err := parseAndSetDebugLevels(appConfig.DebugLevel); err != nil {
		err := fmt.Errorf("loadConfig: %s", err.Error())
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
		return
	}

	// check if we can execute the needed op without connecting to a wallet
	// if len(args) == 0, then there's nothing to execute as all command-line args were parsed as app options
	if len(args) > 0 {
		if ok, err := attemptExecuteSimpleOp(); ok {
			if err != nil {
				fmt.Fprintln(os.Stderr, err.Error())
				os.Exit(1)
			}
			return
		}
	}

	// check if user passed commands/options/args but is not running in cli mode
	if appConfig.InterfaceMode != "cli" && len(args) > 0 {
		fmt.Fprintf(os.Stderr, "unexpected command or flag in %s mode: %s\n", appConfig.InterfaceMode, strings.Join(args, " "))
		os.Exit(1)
	}

	// use wait group to keep main alive until shutdown completes
	shutdownWaitGroup := &sync.WaitGroup{}

	go listenForInterruptRequests()
	go handleShutdown(shutdownWaitGroup)

	// use ctx to monitor potentially long running operations
	// such operations should listen for ctx.Done and stop further processing
	ctx, cancel := context.WithCancel(context.Background())
	shutdownOps = append(shutdownOps, cancel)

	// open connection to wallet and add wallet close function to shutdownOps
	walletMiddleware, err := connectToWallet(ctx, appConfig)
	if err != nil {
		fmt.Fprintln(os.Stderr, "Failed to connect to wallet.", err.Error())
		fmt.Println("Exiting.")
		os.Exit(1)
	}
	shutdownOps = append(shutdownOps, walletMiddleware.CloseWallet)

	switch appConfig.InterfaceMode {
	case "cli":
		enterCliMode(ctx, walletMiddleware, appConfig)
	case "http":
		enterHttpMode(ctx, walletMiddleware, appConfig)
	case "nuklear":
		enterNuklearMode(ctx, walletMiddleware)
	case "fyne":
		enterFyneMode(ctx, walletMiddleware)
	case "terminal":
		enterTerminalMode(ctx, walletMiddleware)
	}

	// wait for handleShutdown goroutine, to finish before exiting main
	shutdownWaitGroup.Wait()
}

//function for writing to stdOut and file simultanously
func logInfo(message string) {
	log.Info(message)
	fmt.Println(message)
}

func logWarn(message string) {
	log.Warn(message)
	fmt.Println(message)
}

func logError(message error) {
	log.Error(message)
	fmt.Println(message)
}

// attemptExecuteSimpleOp checks if the operation requested by the user does not require a connection to a decred wallet
// such operations may include cli commands like `help`, ergo a flags parser object is created with cli commands and flags
// help flag errors (-h, --help) are also handled here, since they do not require access to wallet
func attemptExecuteSimpleOp() (isSimpleOp bool, err error) {
	configWithCommands := &cli.AppConfigWithCliCommands{}
	parser := flags.NewParser(configWithCommands, flags.HelpFlag|flags.PassDoubleDash)

	// use command handler wrapper function to check if any command passed by user can be executed simply
	parser.CommandHandler = func(command flags.Commander, args []string) error {
		if runner.CommandRequiresWallet(command) {
			return nil
		}

		isSimpleOp = true
		commandRunner := runner.New(parser, nil, nil)
		return commandRunner.RunNoneWalletCommands(command, args)
	}

	// re-parse command-line args to catch help flag or execute any commands passed
	_, err = parser.Parse()
	if config.IsFlagErrorType(err, flags.ErrHelp) {
		err = nil
		isSimpleOp = true

		if parser.Active != nil {
			help.PrintCommandHelp(os.Stdout, parser.Name, parser.Active)
		} else {
			help.PrintGeneralHelp(os.Stdout, commands.HelpParser(), commands.Categories())
		}
	}

	return
}

// connectToWallet opens connection to a wallet via any of the available walletmiddleware
// default is connecting directly to a wallet database file via dcrlibwallet
// alternative is connecting to wallet database via dcrwallet rpc (if rpc server address is provided)
func connectToWallet(ctx context.Context, cfg *config.Config) (walletMiddleware app.WalletMiddleware, err error) {
	if cfg.WalletRPCServer == "" {
		// use drlibwallet
		// scan PC for wallet databases and prompt user to select wallet to connect to or create new one
		return walletloader.DetectWallets(ctx, cfg)
	}

	// attempt rpc connection at `cfg.WalletRPCServer`
	rpcWalletMiddleware, rpcConnectionError := dcrwalletrpc.Connect(ctx, cfg.WalletRPCServer, cfg.WalletRPCCert,
		cfg.NoWalletRPCTLS)
	if rpcConnectionError != nil {
		return nil, rpcConnectionError
	}

	// confirm that this rpc connection has a wallet created for it
	walletExists, walletCheckError := rpcWalletMiddleware.WalletExists()
	if walletCheckError != nil {
		return nil, fmt.Errorf("\nError checking if wallet has been created with dcrwallet previously.")
	}
	if !walletExists {
		return nil, fmt.Errorf("\nWallet has not been created with dcrwallet daemon.")
	}

	return rpcWalletMiddleware, nil
}

func enterCliMode(ctx context.Context, walletMiddleware app.WalletMiddleware, appConfig *config.Config) {
	opError = cli.Run(ctx, walletMiddleware, appConfig)
	// cli run done, trigger shutdown
	beginShutdown <- true
}

func enterHttpMode(ctx context.Context, walletMiddleware app.WalletMiddleware, appConfig *config.Config) {
	opError = web.StartServer(ctx, walletMiddleware, appConfig.HTTPHost, appConfig.HTTPPort, &appConfig.Settings)
	// only trigger shutdown if some error occurred, ctx.Err cases would already have triggered shutdown, so ignore
	if opError != nil && ctx.Err() == nil {
		beginShutdown <- true
	}
}

func enterNuklearMode(ctx context.Context, walletMiddleware app.WalletMiddleware) {
	logInfo("Launching desktop app with nuklear")
	nuklear.LaunchApp(ctx, walletMiddleware)
	// todo need to properly listen for shutdown and trigger shutdown
	beginShutdown <- true
}

func enterFyneMode(ctx context.Context, walletMiddleware app.WalletMiddleware) {
	logInfo("Launching desktop app with fyne")
	fyne.LaunchApp(ctx, walletMiddleware)
	beginShutdown <- true
}

func enterTerminalMode(ctx context.Context, walletMiddleware app.WalletMiddleware) {
	fmt.Println("Launching Terminal...")
	opError = terminal.StartTerminalApp(ctx, walletMiddleware)
	// Terminal app closed, trigger shutdown
	beginShutdown <- true
}

func listenForInterruptRequests() {
	interruptChannel := make(chan os.Signal, 1)
	signal.Notify(interruptChannel, os.Interrupt, syscall.SIGTERM)

	// listen for the initial interrupt request and trigger shutdown signal
	sig := <-interruptChannel
	logWarn(fmt.Sprintf("\nReceived %s signal. Shutting down...\n", sig))
	beginShutdown <- true

	// continue to listen for interrupt requests and log that shutdown has already been signaled
	for {
		<-interruptChannel
		logInfo(" Already shutting down... Please wait")
		fmt.Println(" Already shutting down... Please wait")
	}
}

func handleShutdown(wg *sync.WaitGroup) {
	// make wait group wait till shutdownSignal is received and shutdownOps performed
	wg.Add(1)

	<-beginShutdown
	for _, shutdownOp := range shutdownOps {
		shutdownOp()
	}

	// shutdown complete
	wg.Done()

	// check if error occurred while program was running
	if opError != nil {
		os.Exit(1)
	} else {
		os.Exit(0)
	}
}
