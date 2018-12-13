package cli

import (
	"fmt"
	"os"

	"github.com/jessevdk/go-flags"
	"github.com/raedahgroup/godcr/cli/commands"
	"github.com/raedahgroup/godcr/config"
	"github.com/raedahgroup/godcr/walletsource"
)

// Root is the entrypoint to the cli application.
// It defines both the commands and the options available.
type Root struct {
	Commands commands.CliCommands
	Config   config.Config
}

// commandHandler provides a type name for the command handler to register on flags.Parser
type commandHandler func(flags.Commander, []string) error

// CommandHandlerWrapper provides a command handler that provides walletrpcclient.Client
// to commands.WalletCommandRunner types. Other command that satisfy flags.Commander and do not
// depend on walletrpcclient.Client will be run as well.
// If the command does not satisfy any of these types, ErrNotSupported will be returned.
func CommandHandlerWrapper(parser *flags.Parser, walletSource ws.WalletSource) commandHandler {
	return func(command flags.Commander, args []string) error {
		if command == nil {
			return brokenCommandError(parser.Command)
		}
		if commandRunner, ok := command.(commands.WalletCommandRunner); ok {
			return commandRunner.Run(walletSource, args)
		}
		return command.Execute(args)
	}
}

func brokenCommandError(command *flags.Command) error {
	return fmt.Errorf("The command %q was not properly setup.\n" +
		"Please report this bug at https://github.com/raedahgroup/godcr/issues",
		commandName(command))
}

func commandName(command *flags.Command) string {
	name := command.Name
	if command.Active != nil {
		return fmt.Sprintf("%s %s", name, commandName(command.Active))
	}
	return name
}
