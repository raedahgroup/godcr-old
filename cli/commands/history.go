package commands

<<<<<<< HEAD
import (
	"github.com/raedahgroup/godcr/cli/termio"
	ws "github.com/raedahgroup/godcr/walletsource"
)
=======
import "github.com/raedahgroup/dcrcli/cli/utils"
>>>>>>> little refactor

// HistoryCommand enables the user view their transaction history.
type HistoryCommand struct {
	CommanderStub
}

<<<<<<< HEAD
// Run runs the `history` command.
func (h HistoryCommand) Run(walletsource ws.WalletSource, args []string) error {
	transactions, err := walletsource.TransactionHistory()
=======
// Execute runs the `history` command.
func (h HistoryCommand) Execute(args []string) error {
	transactions, err := utils.Wallet.TransactionHistory()
>>>>>>> little refactor
	if err != nil {
		return err
	}

<<<<<<< HEAD
	columns := []string{
		"Date",
		"Amount (DCR)",
		"Direction",
		"Hash",
		"Type",
=======
	res := &utils.Response{
		Columns: []string{
			"Date",
			"Amount (DCR)",
			"Direction",
			"Hash",
			"Type",
		},
		Result: make([][]interface{}, len(transactions)),
>>>>>>> little refactor
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

<<<<<<< HEAD
	termio.PrintTabularResult(termio.StdoutWriter, columns, rows)
=======
	utils.PrintResult(utils.StdoutTabWriter, res)
>>>>>>> little refactor
	return nil
}
