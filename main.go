package main

import (
	"fmt"
	"os"

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

	wallet := connectToWallet(appConfig)

	if appConfig.HTTPMode {
		web.StartHttpServer(wallet, appConfig.HTTPServerAddress)
	} else {
		cli.Run(wallet, appConfig)
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
