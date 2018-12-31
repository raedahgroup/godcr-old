package cli

import (
	"context"
	"fmt"

	"github.com/jessevdk/go-flags"
	"github.com/raedahgroup/godcr/app"
	"github.com/raedahgroup/godcr/app/config"
	"github.com/raedahgroup/godcr/cli/commands"
)

// Root is the entrypoint to the cli application.
// It defines both the commands and the options available.
type Root struct {
	Commands commands.Commands
	Config   config.Config
}

// Run starts the app in cli interface mode
func Run(ctx context.Context, walletMiddleware app.WalletMiddleware, appConfig config.Config) error {
	parser, err := runCliParse(appConfig)
	if err != nil {
		return err
	}

	if appConfig.CreateWallet {
		return createWallet(ctx, walletMiddleware)
	}

	// open wallet, subsequent operations including blockchain sync and command handlers need wallet to be open
	walletExists, err := openWallet(ctx, walletMiddleware)
	if err != nil || !walletExists {
		return err
	}

	if appConfig.SyncBlockchain {
		err = syncBlockChain(ctx, walletMiddleware)
		if err != nil {
			return err
		}
	}

	// parser.Parse checks if a command is passed and invokes the Execute method of the command
	// if no command is passed, parser.Parse returns an error of type ErrCommandRequired
	parser.CommandHandler = commands.CommandHandlerWrapper(parser, walletMiddleware)
	_, err = parser.Parse()
	return err
}

func runCliParse(appConfig config.Config) (*flags.Parser, error) {
	appRoot := Root{Config: appConfig}
	parser := flags.NewParser(&appRoot, flags.HelpFlag|flags.PassDoubleDash)

	// stub out the command handler so that the commands are not executed while loading configuration
	parser.CommandHandler = func(command flags.Commander, args []string) error {
		return nil
	}

	_, err := parser.Parse()
	noCommandPassedError := config.IsFlagErrorType(err, flags.ErrCommandRequired)

	// command mustn't be passed with --sync flag
	if noCommandPassedError && !appConfig.SyncBlockchain {
		displayAvailableCommandsHelpMessage(parser)
	} else if err != nil && !noCommandPassedError {
		fmt.Println(err)
	}

	return parser, err
}
