package commands

import (
	"fmt"
	"strings"

	"github.com/raedahgroup/godcr/cli/termio"
	"github.com/raedahgroup/godcr/walletrpcclient"
)

// ShowTransactionCommand requests for transaction details with a transaction hash.
type ShowTransactionCommand struct {
	CommanderStub
	Detailed bool `short:"d" long:"detailed" description:"Display detailed transaction information"`
	Args     struct {
		TxHash string `positional-arg-name:"transaction hash" required:"yes"`
	} `positional-args:"yes"`
}

// Run runs the get-transaction command, displaying the transaction details to the client.
func (showTxCommand ShowTransactionCommand) Run(walletrpcclient *walletrpcclient.Client, args []string) error {
	transaction, err := walletrpcclient.GetTransaction(showTxCommand.Args.TxHash)
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
	if showTxCommand.Detailed {
		detailedOutput.WriteString("\nInputs\t\n")
		for _, input := range transaction.Inputs {
			detailedOutput.WriteString(fmt.Sprintf("%s\t%s\t\n", input.PreviousOutpoint, input.Amount))
		}
		detailedOutput.WriteString("\nOutputs\t\n")
		for _, out := range transaction.Outputs {
			detailedOutput.WriteString(fmt.Sprintf("%s\t%s\t%s", out.Address, out.ScriptClass, out.Value.String()))
			if out.Internal {
				detailedOutput.WriteString(" (internal)")
			}
			detailedOutput.WriteString("\t\n")
		}
		termio.PrintStringResult(basicOutput, strings.TrimRight(detailedOutput.String(), " \n\r"))
	} else {
		termio.PrintStringResult(basicOutput)
	}
	return nil
}
