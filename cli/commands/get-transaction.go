package commands

import (
	"fmt"

	"github.com/raedahgroup/dcrcli/cli"
)

// GetTransactionCommand requests for transaction details with a transaction hash.
type GetTransactionCommand struct {
	Args struct {
		Hash string `positional-arg-name:"transaction hash" required:"yes"`
	} `positional-args:"yes"`
}

// Execute runs the get-transaction command, displaying the transaction details to the client.
func (g GetTransactionCommand) Execute(args []string) error {
	transaction, err := cli.WalletClient.GetTransaction(g.Args.Hash)
	if err != nil {
		return err
	}
	output := "Transaction\t%s\t\n" +
		"Confirmations\t%d\t\n" +
		"Included in block\t%s\t\n" +
		"Type\t%s\t\n" +
		"Total sent\t%s\t\n" +
		"Time\t%s\t\n" +
		"Size\t%s\t\n" +
		"Fee\t%s\t\n" +
		"Rate\t%s/kB\t"

	txSize := fmt.Sprintf("%.1f kB", float64(transaction.Size)/1000)
	output = fmt.Sprintf(output, transaction.Hash, transaction.Confirmations, transaction.BlockHash,
		transaction.Type, transaction.Amount, transaction.FormattedTime, txSize, transaction.Fee, transaction.Rate)
	cli.PrintStringResult(output)
	return nil
}
