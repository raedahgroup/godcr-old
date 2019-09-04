package main

import (
	"context"
	"fmt"
	"os"
	"path/filepath"
	"sync"

	"github.com/jessevdk/go-flags"
	"github.com/raedahgroup/godcr/app"
	"github.com/raedahgroup/godcr/app/config"
	"github.com/raedahgroup/godcr/app/help"
	"github.com/raedahgroup/godcr/app/walletmediums/dcrlibwallet"
	"github.com/raedahgroup/godcr/app/walletmediums/dcrwalletrpc"
	"github.com/raedahgroup/godcr/cli"
	"github.com/raedahgroup/godcr/cli/commands"
	"github.com/raedahgroup/godcr/cli/runner"
	"github.com/raedahgroup/godcr/cli/walletloader"
)

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

	// use wait group to keep main alive until shutdown completes
	shutdownWaitGroup := &sync.WaitGroup{}

	go listenForShutdownRequests()
	go handleShutdownRequests(shutdownWaitGroup)

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

	if walletMiddleware == nil {
		// there was no error but user did not select a wallet to connect to and did not create a new one
		os.Exit(0)
		return
	}

	shutdownOps = append(shutdownOps, walletMiddleware.CloseWallet)

	err = cli.Run(ctx, walletMiddleware, appConfig)
	if err != nil {
		exitCode = 1
	}
	// cli run done, trigger shutdown
	beginShutdown <- true

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
// default is connecting directly to a wallet database file via dcrlibwallet
// alternative is connecting to wallet database via dcrwallet rpc (if rpc server address is provided)
func connectToWallet(ctx context.Context, cfg *config.Config) (app.WalletMiddleware, error) {
	if cfg.WalletRPCServer == "" {
		walletMiddleware, err := connectViaDcrlibwallet(ctx, cfg)

		// important to return nil, nil explicitly instead of walletMiddleware, err even though they're both nil
		if err == nil && walletMiddleware == nil {
			return nil, nil
		}

		return walletMiddleware, err
	}

	return connectViaDcrWalletRPC(ctx, cfg)
}

// connectViaDcrWalletRPC attempts to load the database at `cfg.DefaultWalletDir`.
// Prompts user to select wallet to connect to if default wallet dir isn't set
// or wallet could not be found at set default dir.
func connectViaDcrlibwallet(ctx context.Context, cfg *config.Config) (*dcrlibwallet.DcrWalletLib, error) {
	// attempt to load default wallet if set and wallet db can be found
	if cfg.DefaultWalletDir != "" {
		netType := filepath.Base(cfg.DefaultWalletDir)
		walletMiddleware, err := dcrlibwallet.Connect(ctx, cfg.DefaultWalletDir, netType)
		if err != nil {
			return nil, err
		}

		defaultWalletExists, walletCheckError := walletMiddleware.WalletExists()
		if walletCheckError != nil {
			return nil, fmt.Errorf("\nError checking default wallet directory for wallet database.\n%s",
				walletCheckError.Error())
		}

		if defaultWalletExists {
			fmt.Println("Using wallet", cfg.DefaultWalletDir)
			return walletMiddleware, nil
		}
	}

	// Scan PC for wallet databases and prompt user to select wallet to connect to or create new one.
	return walletloader.DetectWallets(ctx, cfg)
}

// connectViaDcrWalletRPC attempts an rpc connection to dcrwallet at `cfg.WalletRPCServer`
func connectViaDcrWalletRPC(ctx context.Context, cfg *config.Config) (*dcrwalletrpc.WalletRPCClient, error) {
	rpcWalletMiddleware, rpcConnectionError := dcrwalletrpc.Connect(ctx, cfg)
	if rpcConnectionError != nil {
		return nil, rpcConnectionError
	}

	// confirm that this rpc connection has a wallet created for it
	walletExists, walletCheckError := rpcWalletMiddleware.WalletExists()
	if walletCheckError != nil {
		return nil, fmt.Errorf("\nError checking if wallet has been created with dcrwallet previously.\n%s",
			walletCheckError.Error())
	}
	if !walletExists {
		return nil, fmt.Errorf("\nWallet has not been created with dcrwallet daemon.")
	}

	return rpcWalletMiddleware, nil
}
