package commands

import (
	"github.com/raedahgroup/godcr/cli"
)

// HistoryCommand enables the user view their transaction history.
type HistoryCommand struct{}

// Execute runs the `history` command.
func (h HistoryCommand) Execute(args []string) error {
	transactions, err := cli.WalletClient.GetTransactions()
	if err != nil {
		return err
	}

	res := &cli.Response{
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

	cli.PrintResult(cli.StdoutWriter, res)
	return nil
}
