package commands

import (
	"github.com/raedahgroup/dcrcli/cli"
	"github.com/raedahgroup/dcrcli/walletrpcclient"
)

// HistoryCommand enables the user view their transaction history.
type HistoryCommand struct{}

// Execute runs the `history` command.
func (h HistoryCommand) Execute(_ []string) error {
	res, err := transactionHistory(cli.WalletClient)
	if err != nil {
		return err
	}
	cli.PrintResult(cli.StdoutWriter, res)
	return nil
}

func transactionHistory(walletrpcclient *walletrpcclient.Client) (*cli.Response, error) {
	transactions, err := walletrpcclient.GetTransactions()
	if err != nil {
		return nil, err
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

	return res, nil
}
