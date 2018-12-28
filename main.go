package main

import (
	"context"
	"fmt"
	"github.com/raedahgroup/dcrcli/app"
	"github.com/raedahgroup/dcrcli/app/config"
	"github.com/raedahgroup/dcrcli/app/walletmediums/dcrwalletrpc"
	"github.com/raedahgroup/dcrcli/app/walletmediums/mobilewalletlib"
	"github.com/raedahgroup/dcrcli/cli"
	"github.com/raedahgroup/dcrcli/web"
	"os"
	"os/signal"
	"sync"
	"syscall"
)

// triggered after program execution is complete or if interrupt signal is received
var beginShutdown = make(chan bool)

// shutdownOps holds cleanup/shutdown functions that should be executed when shutdown signal is triggered
var shutdownOps []func()

// opError stores any error that occurs while performing an operation
var opError error

func main() {
	appConfig := config.LoadConfig()
	if appConfig == nil {
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

	if appConfig.HTTPMode {
		opError = web.StartHttpServer(ctx, walletMiddleware, appConfig.HTTPServerAddress)
		// only trigger shutdown if some error occurred, ctx.Err cases would already have triggered shutdown, so ignore
		if opError != nil && ctx.Err() == nil {
			beginShutdown <- true
		}
	} else {
		opError = cli.Run(ctx, walletMiddleware, appConfig)
		// cli run done, trigger shutdown
		beginShutdown <- true
	}

	// wait for handleShutdown goroutine, to finish before exiting main
	shutdownWaitGroup.Wait()
}

// connectToWallet opens connection to a wallet via any of the available walletmiddleware
// default walletmiddleware is mobilewallet library, alternative is dcrwallet rpc
func connectToWallet(ctx context.Context, config *config.Config) app.WalletMiddleware {
	var netType string
	if config.UseTestNet {
		netType = "testnet"
	} else {
		netType = "mainnet"
	}

	if !config.UseWalletRPC {
		return mobilewalletlib.New(config.AppDataDir, netType)
	}

	walletMiddleware, err := dcrwalletrpc.New(ctx, netType, config.WalletRPCServer, config.WalletRPCCert, config.NoWalletRPCTLS)
	if err != nil {
		fmt.Println("Connect to dcrwallet rpc failed")
		fmt.Println(err.Error())
		os.Exit(1)
	}

	return walletMiddleware
}

func listenForInterruptRequests() {
	interruptChannel := make(chan os.Signal, 1)
	signal.Notify(interruptChannel, os.Interrupt, syscall.SIGTERM)

	// listen for the initial interrupt request and trigger shutdown signal
	sig := <-interruptChannel
	fmt.Printf(" Received %s signal. Shutting down...\n", sig)
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
	}
}
