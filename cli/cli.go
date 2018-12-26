package cli

import (
	"fmt"
	"os"

	"github.com/jessevdk/go-flags"
	"github.com/raedahgroup/dcrcli/app"
	"github.com/raedahgroup/dcrcli/app/config"
	"github.com/raedahgroup/dcrcli/cli/utils"
)

// Run starts the app in cli interface mode
func Run(walletMiddleware app.WalletMiddleware, appConfig *config.Config) {
	if appConfig.CreateWallet {
		createWallet(walletMiddleware)
		os.Exit(0)
	}

	// open wallet, subsequent operations including blockchain sync and command handlers need wallet to be open
	openWallet(walletMiddleware)

	if appConfig.SyncBlockchain {
		syncBlockChain(walletMiddleware)
	}

	// Set the core wallet object that will be used by the command handlers
	utils.Wallet = walletMiddleware

	// parser.Parse checks if a command is passed and invokes the Execute method of the command
	// if no command is passed, parser.Parse returns an error of type ErrCommandRequired
	parser := flags.NewParser(appConfig, flags.HelpFlag|flags.PassDoubleDash)
	_, err := parser.Parse()
	if err == nil {
		os.Exit(0)
	}

	// help flag error should have been caught and handled in config.LoadConfig, so only check for ErrCommandRequired
	noCommandPassed := config.IsFlagErrorType(err, flags.ErrCommandRequired)

	if noCommandPassed && appConfig.SyncBlockchain {
		// command mustn't be passed with --sync flag
		os.Exit(0)
	} else if noCommandPassed {
		displayAvailableCommandsHelpMessage(parser)
	} else {
		fmt.Println(err)
	}
	os.Exit(1)
}
