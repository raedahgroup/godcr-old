package commands

import (
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

func (c CreateAccountCommand) Run(wallet walletcore.Wallet) error {
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
