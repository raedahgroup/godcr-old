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

	var columns []string

	if len(accounts) == 1 {
		rows := make([][]interface{}, 1)
		rows[0] = []interface{}{}

		if accounts[0].Balance.Total == accounts[0].Balance.Spendable {
			columns = append(columns, "Total")
			rows[0] = append(rows[0], accounts[0].Balance.Total)
		} else {
			columns = append(columns, "Total", "Spendable")
			rows[0] = append(rows[0], accounts[0].Balance.Total)
			rows[0] = append(rows[0], accounts[0].Balance.Spendable)
		}
		if accounts[0].Balance.LockedByTickets != 0 {
			columns = append(columns, "Locked By Tickets")
			rows[0] = append(rows[0], accounts[0].Balance.LockedByTickets)
		}
		if accounts[0].Balance.VotingAuthority != 0 {
			columns = append(columns, "Voting Authority")
			rows[0] = append(rows[0], accounts[0].Balance.VotingAuthority)
		}
		if accounts[0].Balance.Unconfirmed != 0 {
			columns = append(columns, "Unconfirmed")
			rows[0] = append(rows[0], accounts[0].Balance.Unconfirmed)
		}

		termio.PrintTabularResult(termio.StdoutWriter, columns, rows)
	} else {
		rows := make([][]interface{}, len(accounts))
		for i, account := range accounts {
			rows[i] = []interface{}{}

			columns = append(columns, "Account")
			rows[i] = append(rows[i], account.Name)
			if account.Balance.Total == account.Balance.Spendable {
				columns = append(columns, "Total")
				rows[i] = append(rows[i], account.Balance.Total)
			} else {
				columns = append(columns, " Total", "Spendable")
				rows[i] = append(rows[i], account.Balance.Total)
				rows[i] = append(rows[i], account.Balance.Spendable)
			}
			if account.Balance.LockedByTickets != 0 {
				columns = append(columns, "Locked By Tickets")
				rows[i] = append(rows[i], account.Balance.LockedByTickets)
			}
			if account.Balance.VotingAuthority != 0 {
				columns = append(columns, "Voting Authority")
				rows[i] = append(rows[i], account.Balance.VotingAuthority)
			}
			if account.Balance.Unconfirmed != 0 {
				columns = append(columns, "Unconfirmed")
				rows[i] = append(rows[i], account.Balance.Unconfirmed)
			}
		}
		
		termio.PrintTabularResult(termio.StdoutWriter, columns, rows)
	}

	return nil
}
