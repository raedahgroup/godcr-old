package commands

import (
	"context"

	"github.com/raedahgroup/godcr/app/walletcore"
	"github.com/raedahgroup/godcr/cli/clilog"
)

type CreateAccountCommand struct {
	commanderStub
	Args CreateAccountArgs `positional-args:"yes"`
}
type CreateAccountArgs struct {
	AccountName string `positional-arg-name:"account-name" required:"yes"`
}

func (c CreateAccountCommand) Run(ctx context.Context, wallet walletcore.Wallet) error {
	passphrase, err := getWalletPassphrase()
	if err != nil {
		return err
	}

	_, err = wallet.NextAccount(c.Args.AccountName, passphrase)
	if err != nil {
		return err
	}
	clilog.LogInfo("Account created successfully")
	return nil
}
