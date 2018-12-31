package commands

import (
	"context"
	"github.com/raedahgroup/godcr/app"
	"github.com/raedahgroup/godcr/cli/runner"
	"github.com/raedahgroup/godcr/cli/walletloader"
)

type CreateWalletCommand struct {
	runner.WalletMiddlewareCommand
}

func (c CreateWalletCommand) Run(ctx context.Context, walletMiddleware app.WalletMiddleware, args []string) error {
	return walletloader.CreateWallet(ctx, walletMiddleware)
}
