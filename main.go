package main

import (
	"context"
	"fmt"
	"os"
	"os/signal"
	"path/filepath"
	"strings"
	"sync"
	"syscall"

	"github.com/raedahgroup/godcr/app"
	"github.com/raedahgroup/godcr/app/config"
	"github.com/raedahgroup/godcr/app/walletmediums/dcrlibwallet"
	"github.com/raedahgroup/godcr/fyne"
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

	if walletMiddleware == nil {
		// there was no error but user did not select a wallet to connect to and did not create a new one
		os.Exit(0)
		return
	}

	shutdownOps = append(shutdownOps, walletMiddleware.CloseWallet)

	enterFyneMode(ctx, walletMiddleware)

	// wait for handleShutdown goroutine, to finish before exiting main
	shutdownWaitGroup.Wait()
}

//function for writing to stdOut and file simultaneously
func logInfo(message string) {
	log.Info(message)
	fmt.Println(message)
}

func logWarn(message string) {
	log.Warn(message)
	fmt.Println(message)
}

// connectToWallet opens connection to a wallet via any of the available walletmiddleware
// default is connecting directly to a wallet database file via dcrlibwallet
// alternative is connecting to wallet database via dcrwallet rpc (if rpc server address is provided)
func connectToWallet(ctx context.Context, cfg *config.Config) (app.WalletMiddleware, error) {
	return connectViaDcrlibwallet(ctx, cfg)
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

	// todo: no longer detecting wallets, so should not require default dir
	return nil, fmt.Errorf("default wallet not configured")
}

func enterFyneMode(ctx context.Context, walletMiddleware app.WalletMiddleware) {
	logInfo("Launching desktop app with fyne")
	fyne.LaunchFyne(ctx, walletMiddleware)
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
