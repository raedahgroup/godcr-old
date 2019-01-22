package runner

import (
	"context"
	"github.com/jessevdk/go-flags"
	"github.com/raedahgroup/godcr/app/walletcore"
)

// WalletCommandRunner defines the Run method that cli commands that interact with the decred wallet must implement
// to have access to walletcore.Wallet at execution time
type WalletCommandRunner interface {
	Run(ctx context.Context, wallet walletcore.Wallet, args []string) error
	flags.Commander
}

// WalletCommand implements `flags.Commander`, using a noop Execute method to satisfy `flags.Commander` interface
// Commands embedding this struct should ideally implement `WalletCommandRunner` so that their `Run` method can
// be invoked by `CommandRunner.Run` which will inject the necessary dependencies to run the command
type WalletCommand struct{}

// Noop Execute method added to satisfy `flags.Commander` interface
func (w WalletCommand) Execute(args []string) error {
	return nil
}
