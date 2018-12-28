package commands

import (
	"github.com/raedahgroup/dcrcli/cli/io"
	"github.com/raedahgroup/dcrcli/cli/walletclient"
)

// HistoryCommand enables the user view their transaction history.
type HistoryCommand struct{}

// Execute runs the `history` command.
func (h HistoryCommand) Execute(args []string) error {
	transactions, err := walletclient.WalletClient.GetTransactions()
	if err != nil {
		return err
	}

	res := &io.Response{
		Columns: []string{
			"Date",
			"Amount (DCR)",
			"Direction",
			"Hash",
			"Type",
		},
		Result: make([][]interface{}, len(transactions)),
	}

	for i, tx := range transactions {
		res.Result[i] = []interface{}{
			tx.FormattedTime,
			tx.Amount,
			tx.Direction,
			tx.Hash,
			tx.Type,
		}
	}

	io.PrintResult(io.StdoutWriter, res)
	return nil
}
