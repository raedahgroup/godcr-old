package commands

import (
	"context"
	"github.com/raedahgroup/godcr/cli/walletloader"
)

type CreateWalletCommand struct {
	commanderStub
}

func (c CreateWalletCommand) Run(ctx context.Context) error {
	// any errors encountered are printed to terminal directly, no need to return the error to parser
	walletloader.CreateWallet(ctx)
	return nil
}
