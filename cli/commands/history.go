package commands

import (
	"github.com/raedahgroup/godcr/app/walletcore"
	"github.com/raedahgroup/godcr/cli/termio"
)

// HistoryCommand enables the user view their transaction history.
type HistoryCommand struct {
	commanderStub
}

// Run runs the `history` command.
func (h HistoryCommand) Run(wallet walletcore.Wallet) error {
	transactions, err := wallet.TransactionHistory()
	if err != nil {
		return err
	}

	columns := []string{
		"Date",
		"Amount (DCR)",
		"Direction",
		"Hash",
		"Type",
	}
	rows := make([][]interface{}, len(transactions))

	for i, tx := range transactions {
		rows[i] = []interface{}{
			tx.FormattedTime,
			tx.Amount,
			tx.Direction,
			tx.Hash,
			tx.Type,
		}
	}

	termio.PrintTabularResult(termio.StdoutWriter, columns, rows)
	return nil
}
