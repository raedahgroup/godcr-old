package commands

import (
	"context"
	"fmt"
	"strings"

	"github.com/decred/dcrd/dcrutil"
	godcrUtils "github.com/raedahgroup/godcr/app/utils"
	"github.com/raedahgroup/godcr/app/walletcore"
	"github.com/raedahgroup/godcr/cli/termio"
	"github.com/raedahgroup/godcr/cli/termio/terminalprompt"
)

// ShowTransactionCommand requests for transaction details with a transaction hash.
type ShowTransactionCommand struct {
	commanderStub
	Args ShowTransactionCommandArgs `positional-args:"yes"`
	*historyCommandData
}

type ShowTransactionCommandArgs struct {
	TxHash string `positional-arg-name:"transaction hash" required:"yes"`
}

type historyCommandData struct {
	txHistoryOffset                 int32
	historyCommandDisplayedTxHashes []string
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

	// calculate max number of digits after decimal point for inputs and outputs
	inputsAndOutputsAmount := make([]int64, 0, len(transaction.Inputs)+len(transaction.Outputs))
	for _, txIn := range transaction.Inputs {
		inputsAndOutputsAmount = append(inputsAndOutputsAmount, txIn.Amount)
	}
	for _, txOut := range transaction.Outputs {
		inputsAndOutputsAmount = append(inputsAndOutputsAmount, txOut.Amount)
	}
	maxDecimalPlacesForInputsAndOutputsAmounts := godcrUtils.MaxDecimalPlaces(inputsAndOutputsAmount)

	// now format amount having determined the max number of decimal places
	formatAmount := func(amount int64) string {
		return godcrUtils.FormatAmountDisplay(amount, maxDecimalPlacesForInputsAndOutputsAmounts)
	}
	txDetailsOutput := strings.Builder{}
	txDetailsOutput.WriteString("Transaction Details\n")
	txDetailsOutput.WriteString(basicOutput)
	txDetailsOutput.WriteString("-Inputs- \t \n")
	for _, input := range transaction.Inputs {
		inputAmount := formatAmount(input.Amount)
		txDetailsOutput.WriteString(fmt.Sprintf("  %s \t %s\n", inputAmount, input.PreviousOutpoint))
	}
	txDetailsOutput.WriteString("-Outputs- \t \n") // add tabs to maintain tab spacing for previous inputs section and next outputs section
	for _, out := range transaction.Outputs {
		outputAmount := formatAmount(out.Amount)

		if out.Address == "" {
			txDetailsOutput.WriteString(fmt.Sprintf("  %s \t (no address)\n", outputAmount))
			continue
		}
		txDetailsOutput.WriteString(fmt.Sprintf("  %s \t %s (%s)\n", outputAmount, out.Address, out.AccountName))
	}
	termio.PrintStringResult(strings.TrimRight(txDetailsOutput.String(), " \n\r"))

	if showTxCommand.historyCommandData != nil {
		fmt.Println()
		prompt := fmt.Sprintf("Enter (h)istory table, or (q)uit")

		validateUserInput := func(userInput string) error {
			if strings.EqualFold(userInput, "q") || strings.EqualFold(userInput, "h") {
				return nil
			}

			return fmt.Errorf("invalid response, try again")
		}

		userChoice, err := terminalprompt.RequestInput(prompt, validateUserInput)
		if err != nil {
			return fmt.Errorf("error reading response: %s", err.Error())
		}

		if strings.EqualFold(userChoice, "q") {
			return nil
		}

		showTxHistory := HistoryCommand{
			txHistoryOffset:   showTxCommand.historyCommandData.txHistoryOffset,
			displayedTxHashes: showTxCommand.historyCommandData.historyCommandDisplayedTxHashes,
		}

		fmt.Println()
		err = showTxHistory.Run(ctx, wallet)
		if err != nil {
			return err
		}
	}

	return nil
}
