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

	var showAccount, showTotal, showSpendable, showUnconfirmed bool

	rows := make([][]interface{}, len(accounts))

	if len(accounts) == 1 {
		rows[0] = []interface{}{}

		if accounts[0].Balance.Total == accounts[0].Balance.Spendable {
			showTotal = true
			rows[0] = append(rows[0], accounts[0].Balance.Total)
		} else {
			showTotal, showSpendable = true, true
			rows[0] = append(rows[0], accounts[0].Balance.Total)
			rows[0] = append(rows[0], accounts[0].Balance.Spendable)
		}
		if accounts[0].Balance.LockedByTickets != 0 {
			showSpendable = true
			rows[0] = append(rows[0], accounts[0].Balance.LockedByTickets)
		}
		if accounts[0].Balance.Unconfirmed != 0 {
			showUnconfirmed = true
			rows[0] = append(rows[0], accounts[0].Balance.Unconfirmed)
		}

	} else {
		for i, account := range accounts {
			rows[i] = []interface{}{}

			showAccount = true
			rows[i] = append(rows[i], account.Name)
			if account.Balance.Total == account.Balance.Spendable {
				showTotal = true

				rows[i] = append(rows[i], account.Balance.Total)
			} else {
				showTotal, showSpendable = true, true
				rows[i] = append(rows[i], account.Balance.Total)
				rows[i] = append(rows[i], account.Balance.Spendable)
			}
			if account.Balance.LockedByTickets != 0 {
				showSpendable = true
				rows[i] = append(rows[i], account.Balance.LockedByTickets)
			}
			if account.Balance.Unconfirmed != 0 {
				showUnconfirmed = true
				rows[i] = append(rows[i], account.Balance.Unconfirmed)
			}
		}
	}

	var columns []string
	if showAccount {
		columns = append(columns, "Account")
	}
	if showTotal {
		columns = append(columns, "Total")
	}
	if showSpendable {
		columns = append(columns, "Spendable")
	}
	if showSpendable {
		columns = append(columns, "Locked")
	}
	if showUnconfirmed {
		columns = append(columns, "Unconfirmed")
	}

	termio.PrintTabularResult(termio.StdoutWriter, columns, rows)

	return nil
}
