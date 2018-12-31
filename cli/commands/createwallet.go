package commands

import (
	"context"
	"github.com/raedahgroup/godcr/app"
	"github.com/raedahgroup/godcr/cli/walletloader"
)

type CreateWalletCommand struct {
	commanderStub
}

func (c CreateWalletCommand) Run(ctx context.Context, walletMiddleware app.WalletMiddleware) error {
	return walletloader.CreateWallet(ctx, walletMiddleware)
}
