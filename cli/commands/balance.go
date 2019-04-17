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
	var Account, Total, Spendable, Locked, Voting, Unconfirmed bool

	checkAndAddColumn := func() {
		if Account {
			columns = append(columns, "Account")
		}
		if Total {
			columns = append(columns, "Total")
		}
		if Spendable {
			columns = append(columns, "Spendable")
		}
		if Locked {
			columns = append(columns, "Locked By Tickets")
		}
		if Voting {
			columns = append(columns, "Voting Authority")
		}
		if Unconfirmed {
			columns = append(columns, "Unconfirmed")
		}
	}

	if len(accounts) == 1 {
		rows := make([][]interface{}, 1)
		rows[0] = []interface{}{}

		if accounts[0].Balance.Total == accounts[0].Balance.Spendable {
			Total = true
			rows[0] = append(rows[0], accounts[0].Balance.Total)
		} else {
			Total, Spendable = true, true
			rows[0] = append(rows[0], accounts[0].Balance.Total)
			rows[0] = append(rows[0], accounts[0].Balance.Spendable)
		}
		if accounts[0].Balance.LockedByTickets != 0 {
			Locked = true
			rows[0] = append(rows[0], accounts[0].Balance.LockedByTickets)
		}
		if accounts[0].Balance.VotingAuthority != 0 {
			Voting = true
			rows[0] = append(rows[0], accounts[0].Balance.VotingAuthority)
		}
		if accounts[0].Balance.Unconfirmed != 0 {
			Unconfirmed = true
			rows[0] = append(rows[0], accounts[0].Balance.Unconfirmed)
		}

		checkAndAddColumn()

		termio.PrintTabularResult(termio.StdoutWriter, columns, rows)
	} else {
		rows := make([][]interface{}, len(accounts))

		for i, account := range accounts {
			rows[i] = []interface{}{}

			Account = true
			rows[i] = append(rows[i], account.Name)
			if account.Balance.Total == account.Balance.Spendable {
				Total = true

				rows[i] = append(rows[i], account.Balance.Total)
			} else {
				Total, Spendable = true, true
				rows[i] = append(rows[i], account.Balance.Total)
				rows[i] = append(rows[i], account.Balance.Spendable)
			}
			if account.Balance.LockedByTickets != 0 {
				Locked = true
				rows[i] = append(rows[i], account.Balance.LockedByTickets)
			}
			if account.Balance.VotingAuthority != 0 {
				Voting = true
				rows[i] = append(rows[i], account.Balance.VotingAuthority)
			}
			if account.Balance.Unconfirmed != 0 {
				Unconfirmed = true
				rows[i] = append(rows[i], account.Balance.Unconfirmed)
			}
		}

		checkAndAddColumn()

		termio.PrintTabularResult(termio.StdoutWriter, columns, rows)
	}

	return nil
}
