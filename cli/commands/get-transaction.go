package commands

import (
	"fmt"
	"strings"

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
		"Rate\t%s/kB\t\n"

	txSize := fmt.Sprintf("%.1f kB", float64(transaction.Size)/1000)
	output = fmt.Sprintf(output, transaction.Hash, transaction.Confirmations, transaction.BlockHash,
		transaction.Type, transaction.Amount, transaction.FormattedTime, txSize, transaction.Fee, transaction.Rate)

	inputsBuilder := strings.Builder{}
	inputsBuilder.Grow(len(transaction.Inputs))
	txInputs := "Inputs\t\n"
	inputsBuilder.WriteString(txInputs)
	for _, input := range transaction.Inputs {
		inputsBuilder.WriteString(fmt.Sprintf("%s\t%s\t\n", input.PreviousOutpoint, input.Value))
	}

	outputsBuilder := strings.Builder{}
	outputsBuilder.Grow(len(transaction.Outputs))
	txOutputs := "Outputs\t\n"
	outputsBuilder.WriteString(txOutputs)
	for _, out := range transaction.Outputs {
		outputsBuilder.WriteString(fmt.Sprintf("%s\t%s\t%s", out.Address, out.ScriptClass, out.Value.String()))
		if out.Internal {
			outputsBuilder.WriteString(" (internal)")
		}
		outputsBuilder.WriteString("\t\n")
	}
	cli.PrintStringResult(output, inputsBuilder.String(), strings.TrimSpace(outputsBuilder.String()))
	return nil
}
