package commands

import (
	"fmt"
	"github.com/decred/dcrd/dcrutil"
	"github.com/raedahgroup/dcrcli/cli"
)

// BalanceCommand displays the user's account balance.
type BalanceCommand struct{}

// Execute runs the `balance` command, displaying the user's account balance.
func (b BalanceCommand) Execute(args []string) error {
	accountBalances, err := cli.WalletClient.Balance()
	if err != nil {
		return err
	}

	summarizeBalance := func(total, spendable dcrutil.Amount) string {
		if total == spendable {
			return total.String()
		} else {
			return fmt.Sprintf("Total %s (Spendable %s)", total.String(), spendable.String())
		}
	}

	if len(accountBalances) == 1 {
		commandOutput := summarizeBalance(accountBalances[0].Total, accountBalances[0].Spendable)
		cli.PrintStringResult(commandOutput)
	} else {
		commandOutput := make([]string, len(accountBalances))
		for i, accountBalance := range accountBalances {
			balanceText := summarizeBalance(accountBalance.Total, accountBalance.Spendable)
			commandOutput[i] = fmt.Sprintf("%s \t %s", accountBalance.AccountName, balanceText)
		}
		cli.PrintStringResult(commandOutput...)
	}

	return nil
}
