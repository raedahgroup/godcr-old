package main

import (
	"fmt"
	"github.com/raedahgroup/dcrcli/cli"
	"os"

	"github.com/jessevdk/go-flags"
	"github.com/raedahgroup/dcrcli/cli/commands"
	"github.com/raedahgroup/dcrcli/config"
	"github.com/raedahgroup/dcrcli/core"
	"github.com/raedahgroup/dcrcli/core/middlewares/dcrwalletrpc"
	"github.com/raedahgroup/dcrcli/core/middlewares/mobilewalletlib"
	"github.com/raedahgroup/dcrcli/web"
)

func main() {
	appConfig := config.Default()

	// create parser to parse flags/options from config and commands
	parser := flags.NewParser(&commands.CliCommands{Config: appConfig}, flags.HelpFlag)

	// continueExecution will be false if an error is encountered while parsing or if `-h` or `-v` is encountered
	continueExecution := config.ParseConfig(appConfig, parser)
	if !continueExecution {
		os.Exit(1)
	}

	wallet := connectToWallet(appConfig)

	if appConfig.HTTPMode {
		web.StartHttpServer(wallet, appConfig.HTTPServerAddress)
	} else {
		cli.Run(wallet, appConfig)
	}
}

// makeWalletSource opens connection to a wallet via the selected source/medium
// default is mobile wallet library, alternative is dcrwallet rpc
func connectToWallet(config *config.Config) core.Wallet {
	var netType string
	if config.TestNet {
		netType = "testnet"
	} else {
		netType = "mainnet"
	}

	if !config.UseWalletRPC {
		return mobilewalletlib.New(config.AppDataDir, netType)
	}

	wallet, err := dcrwalletrpc.New(netType, config.WalletRPCServer, config.RPCCert, config.NoDaemonTLS)
	if err != nil {
		fmt.Println("Connect to dcrwallet rpc failed")
		fmt.Println(err.Error())
		os.Exit(1)
	}

	return wallet
}
