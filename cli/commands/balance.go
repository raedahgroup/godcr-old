package commands

import (
	"fmt"

	"github.com/decred/dcrd/dcrutil"
	"github.com/raedahgroup/dcrcli/walletsource"
	"github.com/raedahgroup/godcr/cli/termio"
	ws "github.com/raedahgroup/godcr/walletsource"

	"github.com/raedahgroup/dcrcli/cli/utils"
	"github.com/raedahgroup/dcrcli/core"
)

// BalanceCommand displays the user's account balance.
type BalanceCommand struct {
	CommanderStub
	Detailed bool `short:"d" long:"detailed" description:"Display detailed account balance report"`
}

// Run runs the `balance` command, displaying the user's account balance.
func (balanceCommand BalanceCommand) Run(walletsource ws.WalletSource, args []string) error {
	accounts, err := walletsource.AccountsOverview()
	if err != nil {
		return err
	}

	if balanceCommand.Detailed {
		showDetailedBalance(accounts)
	} else {
		showBalanceSummary(accounts)
	}

	return nil
}

func showDetailedBalance(accountBalances []*walletsource.Account) {
	columns := []string{
		"Account",
		"Total",
		"Spendable",
		"Locked By Tickets",
		"Voting Authority",
		"Unconfirmed",
	}
	rows := make([][]interface{}, len(accountBalances))
	for i, account := range accountBalances {
		rows[i] = []interface{}{
			account.Name,
			account.Balance.Total,
			account.Balance.Spendable,
			account.Balance.LockedByTickets,
			account.Balance.VotingAuthority,
			account.Balance.Unconfirmed,
		}
	}

	termio.PrintTabularResult(termio.StdoutWriter, columns, rows)
}

func showBalanceSummary(accounts []*core.Account) {
	summarizeBalance := func(total, spendable dcrutil.Amount) string {
		if total == spendable {
			return total.String()
		} else {
			return fmt.Sprintf("Total %s (Spendable %s)", total.String(), spendable.String())
		}
	}

	if len(accounts) == 1 {
		commandOutput := summarizeBalance(accounts[0].Balance.Total, accounts[0].Balance.Spendable)
		termio.PrintStringResult(commandOutput)
	} else {
		commandOutput := make([]string, len(accounts))
		for i, account := range accounts {
			balanceText := summarizeBalance(account.Balance.Total, account.Balance.Spendable)
			commandOutput[i] = fmt.Sprintf("%s \t %s", account.Name, balanceText)
		}
		termio.PrintStringResult(commandOutput...)
	}
}
