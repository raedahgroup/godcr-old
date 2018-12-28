package commands

import (
	"github.com/raedahgroup/dcrcli/cli/termio"
	"github.com/raedahgroup/dcrcli/walletrpcclient"
)

// HistoryCommand enables the user view their transaction history.
type HistoryCommand struct{}

// Execute is a stub method to satisfy the commander interface, so that
// it can be passed to the custom command handler which will inject the
// necessary dependencies to run the command.
func (h HistoryCommand) Execute(args []string) error {
	return nil
}

// Execute runs the `history` command.
func (h HistoryCommand) Run(client *walletrpcclient.Client, args []string) error {
	transactions, err := client.GetTransactions()
	if err != nil {
		return err
	}

	res := &termio.Response{
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

	termio.PrintResult(termio.StdoutWriter, res)
	return nil
}
