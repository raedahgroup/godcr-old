package commands

import "github.com/raedahgroup/dcrcli/cli/commands/utils"

// HistoryCommand enables the user view their transaction history.
type HistoryCommand struct{}

// Execute runs the `history` command.
func (h HistoryCommand) Execute(args []string) error {
	transactions, err := utils.Wallet.TransactionHistory()
	if err != nil {
		return err
	}

	res := &utils.Response{
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

	utils.PrintResult(utils.StdoutTabWriter, res)
	return nil
}
