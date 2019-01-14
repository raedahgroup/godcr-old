package main

import (
	"context"
	"fmt"
	"github.com/raedahgroup/godcr/terminal"
	"os"
	"os/signal"
	"strings"
	"sync"
	"syscall"

	flags "github.com/jessevdk/go-flags"
	"github.com/raedahgroup/godcr/app"
	"github.com/raedahgroup/godcr/app/config"
	"github.com/raedahgroup/godcr/app/help"
	"github.com/raedahgroup/godcr/app/walletmediums/dcrlibwallet"
	"github.com/raedahgroup/godcr/app/walletmediums/dcrwalletrpc"
	"github.com/raedahgroup/godcr/cli"
	"github.com/raedahgroup/godcr/cli/commands"
	"github.com/raedahgroup/godcr/cli/runner"
	"github.com/raedahgroup/godcr/desktop"
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

	// check if we can execute the needed op without connecting to a wallet
	// if len(args) == 0, then there's nothing to execute as all command-line args were parsed as app options
	if len(args) > 0 {
		if  ok, err := attemptExecuteSimpleOp(); ok {
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
	walletMiddleware := connectToWallet(ctx, appConfig)
	shutdownOps = append(shutdownOps, walletMiddleware.CloseWallet)

	if appConfig.InterfaceMode == "http" {
		enterHttpMode(ctx, walletMiddleware, appConfig)
	} else if appConfig.InterfaceMode == "nuklear" {
		enterDesktopMode(ctx, walletMiddleware)
	} else if appConfig.InterfaceMode == "terminal" {
		//enterDesktopMode(ctx, walletMiddleware)
		enterTerminalMode()
	} else {
		enterCliMode(ctx, walletMiddleware, appConfig)
	}

	// wait for handleShutdown goroutine, to finish before exiting main
	shutdownWaitGroup.Wait()
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
// default walletmiddleware is dcrlibwallet, alternative is dcrwalletrpc
func connectToWallet(ctx context.Context, config config.Config) app.WalletMiddleware {
	if !config.UseWalletRPC {
		var netType string
		if config.UseTestNet {
			netType = "testnet"
		} else {
			netType = "mainnet"
		}
		return dcrlibwallet.New(config.AppDataDir, netType)
	}

	walletMiddleware, err := dcrwalletrpc.New(ctx, config.WalletRPCServer, config.WalletRPCCert, config.NoWalletRPCTLS)
	if err != nil {
		fmt.Println("Connect to dcrwallet rpc failed")
		fmt.Println(err.Error())
		os.Exit(1)
	}

	return walletMiddleware
}

func enterHttpMode(ctx context.Context, walletMiddleware app.WalletMiddleware, appConfig config.Config) {
	opError = web.StartHttpServer(ctx, walletMiddleware, appConfig.HTTPHost, appConfig.HTTPPort)
	// only trigger shutdown if some error occurred, ctx.Err cases would already have triggered shutdown, so ignore
	if opError != nil && ctx.Err() == nil {
		beginShutdown <- true
	}
}

// todo need to add shutdown functionality to this mode
func enterDesktopMode(ctx context.Context, walletMiddleware app.WalletMiddleware) {
	fmt.Println("Launching desktop app")
	desktop.StartDesktopApp(ctx, walletMiddleware)
	// desktop app closed, trigger shutdown
	beginShutdown <- true
}

func enterTerminalMode() {
	fmt.Println("Launching Terminal...")
	terminal.OpenTerminal()

	beginShutdown <- true
}

func enterCliMode(ctx context.Context, walletMiddleware app.WalletMiddleware, appConfig config.Config) {
	opError = cli.Run(ctx, walletMiddleware, appConfig)
	// cli run done, trigger shutdown
	beginShutdown <- true
}

func listenForInterruptRequests() {
	interruptChannel := make(chan os.Signal, 1)
	signal.Notify(interruptChannel, os.Interrupt, syscall.SIGTERM)

	// listen for the initial interrupt request and trigger shutdown signal
	sig := <-interruptChannel
	fmt.Printf("\nReceived %s signal. Shutting down...\n", sig)
	beginShutdown <- true

	// continue to listen for interrupt requests and log that shutdown has already been signaled
	for {
		<-interruptChannel
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
