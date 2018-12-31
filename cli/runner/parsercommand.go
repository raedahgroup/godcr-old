package runner

import (
	"context"
	"github.com/jessevdk/go-flags"
)

// ParserCommandRunner defines the Run method that cli commands that depends on
// flags.Parser can implement to have it injected at run time
type ParserCommandRunner interface {
	Run(ctx context.Context, parser *flags.Parser, args []string) error
	flags.Commander
}

// ParserCommand implements `flags.Commander`, using a noop Execute method to satisfy `flags.Commander` interface
// Commands embedding this struct should ideally implement `ParserCommandRunner` so that their `Run` method can
// be invoked by `CommandRunner.Run` which will inject the necessary dependencies to run the command
type ParserCommand struct{}

// Noop Execute method added to satisfy `flags.Commander` interface
func (w ParserCommand) Execute(args []string) error {
	return nil
}
