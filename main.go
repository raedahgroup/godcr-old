package main

import (
	"context"
	"fmt"
	"os"
	"os/signal"
	"syscall"

	"github.com/raedahgroup/dcrcli/app"
	"github.com/raedahgroup/dcrcli/app/config"
	"github.com/raedahgroup/dcrcli/app/walletmediums/dcrwalletrpc"
	"github.com/raedahgroup/dcrcli/app/walletmediums/mobilewalletlib"
	"github.com/raedahgroup/dcrcli/cli"
	"github.com/raedahgroup/dcrcli/web"
)

func main() {
	appConfig := config.LoadConfig()
	if appConfig == nil {
		os.Exit(1)
	}

	walletMiddleware := connectToWallet(appConfig)

	// listen for shutdown signals and trigger walletMiddleware.CloseWallet
	ctx, cancel := context.WithCancel(context.Background())
	shutdown := func() {
		cancel()
		walletMiddleware.CloseWallet()
		os.Exit(1)
	}
	go listenForShutdown(shutdown)

	if appConfig.HTTPMode {
		web.StartHttpServer(walletMiddleware, appConfig.HTTPServerAddress, ctx)
	} else {
		cli.Run(walletMiddleware, appConfig)
	}
}

// connectToWallet opens connection to a wallet via any of the available walletmiddleware
// default walletmiddleware is mobilewallet library, alternative is dcrwallet rpc
func connectToWallet(config *config.Config) app.WalletMiddleware {
	var netType string
	if config.UseTestNet {
		netType = "testnet"
	} else {
		netType = "mainnet"
	}

	if !config.UseWalletRPC {
		return mobilewalletlib.New(config.AppDataDir, netType)
	}

	walletMiddleware, err := dcrwalletrpc.New(netType, config.WalletRPCServer, config.WalletRPCCert, config.NoWalletRPCTLS)
	if err != nil {
		fmt.Println("Connect to dcrwallet rpc failed")
		fmt.Println(err.Error())
		os.Exit(1)
	}

	return walletMiddleware
}

func listenForShutdown(shutdown func()) {
	interruptChannel := make(chan os.Signal, 1)
	signal.Notify(interruptChannel, os.Interrupt, syscall.SIGTERM)

	// listen for the initial shutdown signal and begin shutdown/cleanup process in separate goroutine
	// that way, we are able to continue listening for repeated shutdown signals and remind user that a shutdown operation is ongoing
	sig := <-interruptChannel
	fmt.Printf("\nReceived %s signal. Shutting down...\n", sig)
	go shutdown()

	// continue to listen for any more shutdown signals and log that shutdown has already been signaled
	for {
		<-interruptChannel
		fmt.Println("Already shutting down... Please wait")
	}
}
