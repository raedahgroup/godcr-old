package commands

import (
	"fmt"
	"github.com/raedahgroup/dcrcli/cli/utils"
)

type CreateAccountCommand struct {
	Args CreateAccountArgs `positional-args:"yes"`
}
type CreateAccountArgs struct {
	AccountName string `positional-arg-name:"account-name" required:"yes"`
}

func (c CreateAccountCommand) Execute(args []string) error {
	passphrase, err := utils.GetWalletPassphrase()
	if err != nil {
		return err
	}

	_, err = utils.Wallet.NextAccount(c.Args.AccountName, passphrase)
	if err != nil {
		return err
	}

	fmt.Println("Account created successfully")
	return nil
}
