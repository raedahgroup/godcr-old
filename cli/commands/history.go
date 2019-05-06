package commands

import (
	"context"
	"fmt"
	"math"
	"strconv"
	"strings"

	"github.com/decred/dcrd/dcrutil"
	"github.com/raedahgroup/godcr/app/walletcore"
	"github.com/raedahgroup/godcr/cli/termio"
	"github.com/raedahgroup/godcr/cli/termio/terminalprompt"
)

// HistoryCommand enables the user view their transaction history.
type HistoryCommand struct {
	commanderStub
	txHistoryOffset   int32
	displayedTxHashes []string
}

// Run runs the `history` command.
func (history HistoryCommand) Run(ctx context.Context, wallet walletcore.Wallet) error {
	columns := []string{
		"#",
		"Date",
		"Direction",
		centerAlignAmountHeader("Amount"),
		centerAlignAmountHeader("Fee"),
		"Type",
	}

	txCount, err := wallet.TransactionCount(nil)
	if err != nil {
		return fmt.Errorf("cannot load history, getting tx count failed with error: %s", err.Error())
	}

	var txPerPage int32 = walletcore.TransactionHistoryCountPerPage
	var previous bool
	var previousPageTxCount int

	// show transactions in pages, using infinite loop
	// after displaying transactions for each page,
	// ask user if to show next page, previous page, tx details or exit the loop
	for {
		transactions, err := wallet.TransactionHistory(history.txHistoryOffset, txPerPage, nil)
		if err != nil {
			return err
		}

		if previous {
			history.displayedTxHashes = history.displayedTxHashes[:len(history.displayedTxHashes)-previousPageTxCount]
		}

		previousPageTxCountRowNumber := len(history.displayedTxHashes) + 1

		pageTxRows := make([][]interface{}, len(transactions))

		for i, tx := range transactions {
			pageTxRows[i] = []interface{}{}
			if previous {
				pageTxRows[i] = append(pageTxRows[i], previousPageTxCountRowNumber-(len(transactions)-i))
			} else {
				history.displayedTxHashes = append(history.displayedTxHashes, tx.Hash)
				pageTxRows[i] = append(pageTxRows[i], previousPageTxCountRowNumber+i)
			}
			pageTxRows[i] = append(pageTxRows[i], tx.ShortTime)
			pageTxRows[i] = append(pageTxRows[i], tx.Direction)
			pageTxRows[i] = append(pageTxRows[i], formatAmount(tx.Amount))
			pageTxRows[i] = append(pageTxRows[i], formatAmount(tx.Fee))
			pageTxRows[i] = append(pageTxRows[i], tx.Type)
		}

		previousPageTxCount = len(transactions)
		termio.PrintTabularResult(termio.StdoutWriter, columns, pageTxRows)

		fmt.Println() // print empty line before listing txs for next page

		var pageInfo string
		if previous {
			pageInfo = fmt.Sprintf("Showing transactions %d-%d of %d", previousPageTxCountRowNumber-(len(transactions)),
				previousPageTxCountRowNumber-1, txCount)
		} else {
			pageInfo = fmt.Sprintf("Showing transactions %d-%d of %d", previousPageTxCountRowNumber,
				previousPageTxCountRowNumber+len(transactions)-1, txCount)
		}
		// ask user what to do next
		var prompt string
		if len(history.displayedTxHashes) < txCount && len(history.displayedTxHashes) == int(txPerPage) {
			prompt = fmt.Sprintf("%s, enter # for details, show (n)ext, or (q)uit", pageInfo)
		} else if len(history.displayedTxHashes) >= int(txPerPage) && len(history.displayedTxHashes) < txCount {
			prompt = fmt.Sprintf("%s, enter # for details, show (p)revious, (n)ext, or (q)uit", pageInfo)
		} else if len(history.displayedTxHashes) == txCount && len(history.displayedTxHashes) > int(txPerPage) {
			prompt = fmt.Sprintf("%s, enter # for details, show (p)revious, or (q)uit", pageInfo)
		} else {
			prompt = fmt.Sprintf("%s, enter # for details or (q)uit", pageInfo)
		}

		validateUserInput := func(userInput string) error {
			if strings.EqualFold(userInput, "q") ||
				(strings.EqualFold(userInput, "n") && len(history.displayedTxHashes) < txCount) ||
				(strings.EqualFold(userInput, "p") && (history.txHistoryOffset-txPerPage) >= 0) {
				return nil
			}

			// check if user input is a valid tx #
			txRowNumber, err := strconv.ParseUint(userInput, 10, 32)
			if err != nil || txRowNumber < 1 || int(txRowNumber) > len(history.displayedTxHashes) {
				return fmt.Errorf("invalid response, try again")
			}

			return nil
		}

		userChoice, err := terminalprompt.RequestInput(prompt, validateUserInput)
		if err != nil {
			return fmt.Errorf("error reading response: %s", err.Error())
		}

		if strings.EqualFold(userChoice, "q") {
			break
		} else if strings.EqualFold(userChoice, "n") {
			previous = false
			fmt.Println() // print empty line before listing txs for next page
			history.txHistoryOffset += int32(len(transactions))
			continue
		} else if strings.EqualFold(userChoice, "p") {
			previous = true
			fmt.Println() // print empty line before listing txs for next page
			history.txHistoryOffset -= txPerPage
			continue
		}

		// if the code execution continues to this point, it means user's response was neither "q" nor "n" nor "p"
		// must therefore be a tx # to view tx details
		txRowNumber, _ := strconv.ParseUint(userChoice, 10, 32)
		txHash := history.displayedTxHashes[txRowNumber-1]

		showTransactionCommandArgs := ShowTransactionCommandArgs{txHash}

		showTxDetails := ShowTransactionCommand{
			Args:              showTransactionCommandArgs,
			txHistoryOffset:   history.txHistoryOffset,
			displayedTxHashes: history.displayedTxHashes,
			Detailed:          true,
		}

		fmt.Println()
		err = showTxDetails.Run(ctx, wallet)
		if err == nil {
			fmt.Println()
		}
		return err
	}

	return nil
}

// centerAlignAmountHeader returns the Amount or Fee header as a 17-character string
// padded with equal spaces to the left and right
func centerAlignAmountHeader(header string) string {
	nHeaderCharacters := len(header)
	if nHeaderCharacters < 17 {
		spacesToPad := math.Floor(17.0 - float64(nHeaderCharacters)/2.0)
		nLeftSpaces := int(spacesToPad)
		nRightSpaces := int(17 - nHeaderCharacters - nLeftSpaces)

		header = fmt.Sprintf("%*s", nLeftSpaces, header)   // pad with spaces to the left
		header = fmt.Sprintf("%-*s", nRightSpaces, header) // pad with spaces to the right
	}
	return header
}

// formatAmount returns the amount as a 17-character string padded with spaces to the left
func formatAmount(amount int64) string {
	amountString := dcrutil.Amount(amount).String()
	return fmt.Sprintf("%17s", amountString)
}
