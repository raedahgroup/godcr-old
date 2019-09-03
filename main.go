package main

import (
	"context"
	"fmt"
	"sync"

	"github.com/raedahgroup/godcr/app/config"
	w "github.com/raedahgroup/godcr/app/wallet"
	"github.com/raedahgroup/godcr/app/wallet/libwallet"
	"github.com/raedahgroup/godcr/fyne"
)

func main() {
	// Initialize log rotation.  After log rotation has been initialized, the
	// logger variables may be used.
	initLogRotator(config.LogFile)
	defer func() {
		if logRotator != nil {
			logRotator.Close()
		}
	}()

	fyneUI := fyne.InitializeUserInterface()

	// nb: cli support will require loading from a config file
	cfg, err := config.LoadConfigFromDb()
	if err != nil {
		errorMessage := fmt.Sprintf("Error loading config from db: %v", err)
		log.Errorf(errorMessage)
		fyneUI.DisplayPreLaunchError(errorMessage)
		return
	}

	// Parse, validate, and set debug log level(s).
	if err := parseAndSetDebugLevels(cfg.DebugLevel); err != nil {
		errorMessage := fmt.Sprintf("Error setting log levels: %v", err)
		log.Errorf(errorMessage)
		fyneUI.DisplayPreLaunchError(errorMessage)
		return
	}

	// use wait group to keep main alive until shutdown completes
	shutdownWaitGroup := &sync.WaitGroup{}
	go handleShutdownRequests(shutdownWaitGroup)
	go listenForShutdownRequests()

	// open connection to wallet and add wallet shutdown function to shutdownOps
	wallet, err := connectToWallet(cfg)
	if err != nil {
		errorMessage := fmt.Sprintf("Failed to connect to wallet: %v", err)
		log.Errorf(errorMessage)
		fyneUI.DisplayPreLaunchError(errorMessage)
		return
	}
	shutdownOps = append(shutdownOps, wallet.Shutdown)

	// use ctx to monitor potentially long running operations
	// such operations should listen for ctx.Done and stop further processing
	ctx, cancel := context.WithCancel(context.Background())
	shutdownOps = append(shutdownOps, cancel)
	fyneUI.LaunchApp(ctx, cfg, wallet)

	// fyneUI.LaunchApp calls fyne's showandrun function which blocks until the fyne app is exited
	// beginshutdown calls for an exit to app when fyneUI quits
	beginShutdown <- true

	// wait for handleShutdownRequests goroutine, to finish before exiting main
	shutdownWaitGroup.Wait()
}

// connectToWallet opens connection to a wallet via dcrlibwallet (LibWallet)
// or dcrwalletrpc (RpcWallet, currently unimplemented.unsuported)
func connectToWallet(cfg *config.Config) (w.Wallet, error) {
	netType := "mainnet"
	if cfg.UseTestnet {
		netType = "testnet3"
	}
	return libwallet.Init(cfg.AppDataDir, netType)
}
