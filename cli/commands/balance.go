package commands

import (
	"context"

	"github.com/raedahgroup/godcr/app/walletcore"
	"github.com/raedahgroup/godcr/cli/termio"
)

// Balance displays the user's account balance.
type BalanceCommand struct {
	commanderStub
}

// Run runs the `balance` command, displaying the user's account balance.
func (balanceCommand BalanceCommand) Run(ctx context.Context, wallet walletcore.Wallet) error {
	accounts, err := wallet.AccountsOverview(walletcore.DefaultRequiredConfirmations)
	if err != nil {
		return err
	}

	termio.PrintStringResult(walletcore.WalletBalance(accounts))

	return nil
}
