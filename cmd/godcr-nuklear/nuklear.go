package main

import (
	"context"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"sync"

	"github.com/raedahgroup/godcr/app"
	"github.com/raedahgroup/godcr/app/config"
	"github.com/raedahgroup/godcr/app/walletmediums/dcrlibwallet"
	"github.com/raedahgroup/godcr/app/walletmediums/dcrwalletrpc"
	"github.com/raedahgroup/godcr/cli/walletloader"
	"github.com/raedahgroup/godcr/nuklear"
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

	// Parse, validate, and set debug log level(s).
	if err := parseAndSetDebugLevels(appConfig.DebugLevel); err != nil {
		err := fmt.Errorf("loadConfig: %s", err.Error())
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
		return
	}

	// check if user passed commands/options/args in non-cli (nuklear) mode
	if len(args) > 0 {
		fmt.Fprintf(os.Stderr, "Unexpected command or flag in %s mode: %s.\n",
			"nuklear",
			strings.Join(args, " "))
		os.Exit(1)
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

	log.Info("Launching desktop app with nuklear")
	nuklear.LaunchApp(ctx, walletMiddleware, &appConfig.Settings)
	// todo need to properly listen for shutdown and trigger shutdown
	beginShutdown <- true

	// wait for handleShutdown goroutine, to finish before exiting main
	shutdownWaitGroup.Wait()
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
