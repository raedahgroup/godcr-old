package commands

import (
	"context"
	"fmt"
	"strings"

	"github.com/decred/dcrd/dcrutil"
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
func (showTxCommand ShowTransactionCommand) Run(ctx context.Context, wallet walletcore.Wallet) error {
	transaction, err := wallet.GetTransaction(showTxCommand.Args.TxHash)
	if err != nil {
		return err
	}

	basicOutput := "  Hash \t %s\n" +
		"  Confirmations \t %d\n" +
		"  Included in block \t %d\n" +
		"  Type \t %s\n" +
		"  Amount %s \t %s\n" +
		"  Date \t %s\n" +
		"  Size \t %s\n" +
		"  Fee \t %s\n" +
		"  Rate \t %s/kB\n"

	txDirection := strings.ToLower(transaction.Direction.String())
	txSize := fmt.Sprintf("%.1f kB", float64(transaction.Size)/1000)
	basicOutput = fmt.Sprintf(basicOutput,
		transaction.Hash,
		transaction.Confirmations,
		transaction.BlockHeight,
		transaction.Type,
		txDirection, transaction.Amount,
		transaction.FormattedTime,
		txSize,
		transaction.Fee,
		transaction.FeeRate)

	if showTxCommand.Detailed {
		detailedOutput := strings.Builder{}
		detailedOutput.WriteString("Transaction Details\n")
		detailedOutput.WriteString(basicOutput)
		detailedOutput.WriteString("Inputs \t \n")
		for _, input := range transaction.Inputs {
			inputAmount := dcrutil.Amount(input.AmountIn).String()
			detailedOutput.WriteString(fmt.Sprintf("  %s \t %s\n", inputAmount, input.PreviousOutpoint))
		}
		detailedOutput.WriteString("Outputs \t \n") // add tabs to maintain tab spacing for previous inputs section and next outputs section
		for _, out := range transaction.Outputs {
			if len(out.Addresses) == 0 {
				detailedOutput.WriteString(fmt.Sprintf("  %s \t (no address)\n", dcrutil.Amount(out.Value).String()))
				continue
			}

			outputAmount := dcrutil.Amount(out.Value).String()
			for _, address := range out.Addresses {
				accountName := address.AccountName
				if !address.IsMine {
					accountName = "external"
				}
				detailedOutput.WriteString(fmt.Sprintf("  %s \t %s (%s)\n", outputAmount, address.Address, accountName))
			}
		}
		termio.PrintStringResult(strings.TrimRight(detailedOutput.String(), " \n\r"))
	} else {
		termio.PrintStringResult(basicOutput)
	}
	return nil
}
