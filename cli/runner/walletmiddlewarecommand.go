package runner

import (
	"context"
	"github.com/jessevdk/go-flags"
	"github.com/raedahgroup/godcr/app"
)

// WalletMiddlewareCommandRunner defines the Run method that cli commands must implement to have access to app.WalletMiddleware
// in order to perform wallet creation/opening/closing and blockchain syncing operations
type WalletMiddlewareCommandRunner interface {
	Run(ctx context.Context, walletMiddleware app.WalletMiddleware, args []string) error
	flags.Commander
}

// WalletMiddlewareCommand implements `flags.Commander`, using a noop Execute method to satisfy `flags.Commander` interface
// Commands embedding this struct should ideally implement `WalletMiddlewareCommandRunner` so that their `Run` method can
// be invoked by `CommandRunner.Run` which will inject the necessary dependencies to run the command
type WalletMiddlewareCommand struct{}

// Noop Execute method added to satisfy `flags.Commander` interface
func (w WalletMiddlewareCommand) Execute(args []string) error {
	return nil
}
