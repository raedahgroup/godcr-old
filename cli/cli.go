package cli

import (
	"fmt"

	"github.com/jessevdk/go-flags"
	"github.com/raedahgroup/dcrcli/cli/commands"
	"github.com/raedahgroup/dcrcli/config"
	"github.com/raedahgroup/dcrcli/walletrpcclient"
)

// CliRoot is the entrypoint to the cli application.
// It defines both the commands and the options available.
type CliRoot struct {
	Commands commands.CliCommands
	Config   config.Config
}

// CommandHandler provides a type name for the command handler to register on flags.Parser
type CommandHandler func(flags.Commander, []string) error

// CommandHandlerWrapper provides a command handler that provides walletrpcclient.Client
// to commands.WalletCommandRunner types. Other command that satisfy flags.Commander and do not
// depend on walletrpcclient.Client will be run as well.
// If the command does not satisfy any of these types, ErrNotSupported will be returned.
func CommandHandlerWrapper(client *walletrpcclient.Client) CommandHandler {
	return func(command flags.Commander, args []string) error {
		if command == nil {
			return fmt.Errorf("unsupported command")
		}
		if commandRunner, ok := command.(commands.WalletCommandRunner); ok {
			return commandRunner.Run(client, args)
		}
		return command.Execute(args)
	}
}
