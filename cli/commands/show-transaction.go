package commands

import (
	"context"
	"fmt"
	"strings"

	"github.com/decred/dcrd/dcrutil"
	"github.com/raedahgroup/godcr/app/walletcore"
	"github.com/raedahgroup/godcr/cli/termio"
	"github.com/raedahgroup/godcr/cli/termio/terminalprompt"
)

// ShowTransactionCommand requests for transaction details with a transaction hash.
type ShowTransactionCommand struct {
	commanderStub
	Detailed          bool                       `short:"d" long:"detailed" description:"Display detailed transaction information"`
	Args              ShowTransactionCommandArgs `positional-args:"yes"`
	txHistoryOffset   int32
	displayedTxHashes []string
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
		"  Fee Rate \t %s/kB\n"

	txDirection := strings.ToLower(transaction.Direction.String())
	txSize := fmt.Sprintf("%.1f kB", float64(transaction.Size)/1000)
	basicOutput = fmt.Sprintf(basicOutput,
		transaction.Hash,
		transaction.Confirmations,
		transaction.BlockHeight,
		transaction.Type,
		txDirection, dcrutil.Amount(transaction.Amount).String(),
		fmt.Sprintf("%s UTC", transaction.LongTime),
		txSize,
		dcrutil.Amount(transaction.Fee).String(),
		dcrutil.Amount(transaction.FeeRate).String(),
	)

	if showTxCommand.Detailed {
		detailedOutput := strings.Builder{}
		detailedOutput.WriteString("Transaction Details\n")
		detailedOutput.WriteString(basicOutput)
		detailedOutput.WriteString("Inputs \t \n")
		for _, input := range transaction.Inputs {
			inputAmount := dcrutil.Amount(input.Amount).String()
			detailedOutput.WriteString(fmt.Sprintf("  %s \t %s  (%s)\n", inputAmount, input.PreviousOutpoint, input.AccountName))
		}
		detailedOutput.WriteString("Outputs \t \n") // add tabs to maintain tab spacing for previous inputs section and next outputs section
		for _, out := range transaction.Outputs {
			outputAmount := dcrutil.Amount(out.Amount).String()

			if out.Address == "" {
				detailedOutput.WriteString(fmt.Sprintf("  %s \t (no address)\n", outputAmount))
				continue
			}
			detailedOutput.WriteString(fmt.Sprintf("  %s \t %s (%s)\n", outputAmount, out.Address, out.AccountName))
		}
		termio.PrintStringResult(strings.TrimRight(detailedOutput.String(), " \n\r"))

		// var prompt string
		if len(showTxCommand.displayedTxHashes) > 0 {
			fmt.Println()
			prompt := fmt.Sprintf("Enter (h)istory table, or (q)uit")

			validateUserInput := func(userInput string) error {
				if strings.EqualFold(userInput, "q") || strings.EqualFold(userInput, "h") {
					return nil
				}
				return nil
			}

			userChoice, err := terminalprompt.RequestInput(prompt, validateUserInput)
			if err != nil {
				return fmt.Errorf("error reading response: %s", err.Error())
			}

			if strings.EqualFold(userChoice, "q") {
				return nil
			}

			var displayedTxHashes []string
			displayedTxHashes = showTxCommand.displayedTxHashes
			displayedTxHashes = displayedTxHashes[:len(displayedTxHashes)-(len(displayedTxHashes)-int(showTxCommand.txHistoryOffset))]

			showTxHistory := HistoryCommand{
				txHistoryOffset:   showTxCommand.txHistoryOffset,
				displayedTxHashes: displayedTxHashes,
			}

			err = showTxHistory.Run(ctx, wallet)
			if err == nil {
				fmt.Println()
			}
			return err
		}
	} else {
		termio.PrintStringResult(basicOutput)
	}
	return nil
}
