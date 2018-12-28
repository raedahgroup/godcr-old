package cli

import (
	"github.com/raedahgroup/dcrcli/cli/commands"
	"github.com/raedahgroup/dcrcli/config"
)

// AppRoot is the entrypoint to the cli application.
// It defines both the commands and the options available.
type AppRoot struct {
	Commands commands.CliCommands
	Config   config.Config
}
