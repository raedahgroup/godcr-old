package commands

import (
	"fmt"
	"strings"

	"github.com/raedahgroup/dcrcli/cli"
)

// GetTransactionCommand requests for transaction details with a transaction hash.
type GetTransactionCommand struct {
	Detailed bool `short:"d" long:"detailed" description:"Display detailed transaction information"`
	Args     struct {
		Hash string `positional-arg-name:"transaction hash" required:"yes"`
	} `positional-args:"yes"`
}

// Execute runs the get-transaction command, displaying the transaction details to the client.
func (g GetTransactionCommand) Execute(args []string) error {
	transaction, err := cli.WalletClient.GetTransaction(g.Args.Hash)
	if err != nil {
		return err
	}
	basicOutput := "Transaction\t%s\t\n" +
		"Confirmations\t%d\t\n" +
		"Included in block\t%s\t\n" +
		"Type\t%s\t\n" +
		"Total sent\t%s\t\n" +
		"Time\t%s\t\n" +
		"Size\t%s\t\n" +
		"Fee\t%s\t\n" +
		"Rate\t%s/kB\t"

	txSize := fmt.Sprintf("%.1f kB", float64(transaction.Size)/1000)
	basicOutput = fmt.Sprintf(basicOutput, transaction.Hash, transaction.Confirmations, transaction.BlockHash,
		transaction.Type, transaction.Amount, transaction.FormattedTime, txSize, transaction.Fee, transaction.Rate)

	detailedOutput := strings.Builder{}
	if g.Detailed {
		detailedOutput.WriteString("\nInputs\t\n")
		for _, input := range transaction.Inputs {
			detailedOutput.WriteString(fmt.Sprintf("%s\t%s\t\n", input.PreviousOutpoint, input.Value))
		}
		detailedOutput.WriteString("\nOutputs\t\n")
		for _, out := range transaction.Outputs {
			detailedOutput.WriteString(fmt.Sprintf("%s\t%s\t%s", out.Address, out.ScriptClass, out.Value.String()))
			if out.Internal {
				detailedOutput.WriteString(" (internal)")
			}
			detailedOutput.WriteString("\t\n")
		}
		cli.PrintStringResult(basicOutput, strings.TrimRight(detailedOutput.String(), " \n\r"))
	} else {
		cli.PrintStringResult(basicOutput)
	}
	return nil
}
