package commands

import (
	"fmt"
	"strings"

	"github.com/raedahgroup/godcr/app/walletcore"
	"github.com/raedahgroup/godcr/cli/termio"
)

// ShowTransactionCommand requests for transaction details with a transaction hash.
type ShowTransactionCommand struct {
	commanderStub
	Detailed bool                       `short:"d" long:"detailed" description:"Display detailed transaction information"`
	Args     ShowTransactionCommandArgs `positional-args:"yes"`
}
type ShowTransactionCommandArgs struct {
	TxHash string `positional-arg-name:"transaction hash" required:"yes"`
}

// Run runs the get-transaction command, displaying the transaction details to the client.
func (showTxCommand ShowTransactionCommand) Run(wallet walletcore.Wallet) error {
	transaction, err := wallet.GetTransaction(showTxCommand.Args.TxHash)
	if err != nil {
		return err
	}

	basicOutput := "Hash\t%s\n" +
		"Confirmations\t%d\n" +
		"Included in block\t%s\n" +
		"Type\t%s\n" +
		"Amount %s\t%s\n" +
		"Date\t%s\n" +
		"Size\t%s\n" +
		"Fee\t%s\n" +
		"Rate\t%s/kB\n"

	txDirection := strings.ToLower(transaction.Direction.String())
	txSize := fmt.Sprintf("%.1f kB", float64(transaction.Size)/1000)
	basicOutput = fmt.Sprintf(basicOutput,
		transaction.Hash,
		transaction.Confirmations,
		transaction.BlockHash,
		transaction.Type,
		txDirection, transaction.Amount,
		transaction.FormattedTime,
		txSize,
		transaction.Fee,
		transaction.Rate)

	if showTxCommand.Detailed {
		detailedOutput := strings.Builder{}
		detailedOutput.WriteString("General Info\n")
		detailedOutput.WriteString(basicOutput)
		detailedOutput.WriteString("\nInputs\n")
		for _, input := range transaction.Inputs {
			detailedOutput.WriteString(fmt.Sprintf("%s\t%s\n", input.Amount, input.PreviousOutpoint))
		}
		detailedOutput.WriteString("\nOutputs\n")
		for _, out := range transaction.Outputs {
			detailedOutput.WriteString(fmt.Sprintf("%s\t%s", out.Value, out.Address))
			if out.Internal {
				detailedOutput.WriteString(" (internal)")
			}
			detailedOutput.WriteString("\n")
		}
		termio.PrintStringResult(strings.TrimRight(detailedOutput.String(), " \n\r"))
	} else {
		termio.PrintStringResult(basicOutput)
	}
	return nil
}
