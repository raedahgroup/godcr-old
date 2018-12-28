package cli

import (
	"context"
	"fmt"

	"github.com/jessevdk/go-flags"
	"github.com/raedahgroup/dcrcli/app"
	"github.com/raedahgroup/dcrcli/app/config"
	"github.com/raedahgroup/dcrcli/cli/utils"
)

// Run starts the app in cli interface mode
func Run(ctx context.Context, walletMiddleware app.WalletMiddleware, appConfig *config.Config) error {
	if appConfig.CreateWallet {
		return createWallet(walletMiddleware)
	}

	// open wallet, subsequent operations including blockchain sync and command handlers need wallet to be open
	err := openWallet(ctx, walletMiddleware)
	if err != nil {
		return err
	}

	if appConfig.SyncBlockchain {
		err = syncBlockChain(walletMiddleware)
		if err != nil {
			return err
		}
	}

	// Set the core wallet object that will be used by the command handlers
	utils.Wallet = walletMiddleware

	// parser.Parse checks if a command is passed and invokes the Execute method of the command
	// if no command is passed, parser.Parse returns an error of type ErrCommandRequired
	parser := flags.NewParser(appConfig, flags.HelpFlag|flags.PassDoubleDash)
	_, err = parser.Parse()

	// help flag error should have been caught and handled in config.LoadConfig, so only check for ErrCommandRequired
	noCommandPassedError := config.IsFlagErrorType(err, flags.ErrCommandRequired)

	// command mustn't be passed with --sync flag
	if noCommandPassedError && !appConfig.SyncBlockchain {
		displayAvailableCommandsHelpMessage(parser)
	} else if err != nil && !noCommandPassedError {
		fmt.Println(err)
	}

	return err
}
