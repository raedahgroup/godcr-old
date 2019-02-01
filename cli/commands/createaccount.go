package commands

import (
	"context"
	"fmt"

	"github.com/raedahgroup/godcr/app/walletcore"
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
	log.Infof("Account created successfully")
	fmt.Println("Account created successfully")
	return nil
}
