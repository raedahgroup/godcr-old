package commands

import (
	"context"
	"fmt"
	"github.com/raedahgroup/godcr/app/walletcore"
	"github.com/raedahgroup/godcr/cli/runner"
)

type CreateAccountCommand struct {
	runner.WalletCommand
	Args CreateAccountArgs `positional-args:"yes"`
}
type CreateAccountArgs struct {
	AccountName string `positional-arg-name:"account-name" required:"yes"`
}

func (c CreateAccountCommand) Run(ctx context.Context, wallet walletcore.Wallet, args []string) error {
	passphrase, err := getWalletPassphrase()
	if err != nil {
		return err
	}

	_, err = wallet.NextAccount(c.Args.AccountName, passphrase)
	if err != nil {
		return err
	}

	fmt.Println("Account created successfully")
	return nil
}
