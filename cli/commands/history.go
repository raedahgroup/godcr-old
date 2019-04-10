package commands

import (
	"context"
	"fmt"
	// "math"
	"strconv"
	"strings"
	"time"

	"github.com/raedahgroup/godcr/app/walletcore"
	"github.com/raedahgroup/godcr/cli/termio"
	"github.com/raedahgroup/godcr/cli/termio/terminalprompt"
)

// HistoryCommand enables the user view their transaction history.
type HistoryCommand struct {
	commanderStub
}

// Run runs the `history` command.
func (h HistoryCommand) Run(ctx context.Context, wallet walletcore.Wallet) error {
	var startBlockHeight int32 = -1
	var displayedTxHashes []string

	columns := []string{
		"Date (UTC)",
		"Direction",
		"Amount",
		"Status",
		"Type",
	}

	// show transactions in pages, using infinite loop
	// after displaying transactions for each page,
	// ask user if to show next page, previous page, tx details or exit the loop
	for {
		transactions, endBlockHeight, err := wallet.TransactionHistory(ctx, startBlockHeight, walletcore.TransactionHistoryCountPerPage)
		if err != nil {
			return err
		}

		// next start block should be the block immediately preceding the current end block
		startBlockHeight = endBlockHeight - 1

		lastTxRowNumber := len(displayedTxHashes) + 1

		loc, _ := time.LoadLocation("UTC")
		currentDate := time.Now().In(loc).Add(1 * time.Hour)
		timeDifference, _ := time.ParseDuration("24h")

		var confirmations int32 
		confirmations = walletcore.DefaultRequiredConfirmations

		pageTxRows := make([][]interface{}, len(transactions))
		for i, tx := range transactions {
			displayedTxHashes = append(displayedTxHashes, tx.Hash)

			transactionDate := time.Unix(tx.Timestamp, 0).In(loc).Add(1 * time.Hour)
			transactionDuration := currentDate.Sub(transactionDate)

			dateOutput := strings.Split(tx.FormattedTime, " ")

			var txDate string
			if transactionDuration > timeDifference {
		  		txDate = dateOutput[0]
			}else{
		    	txDate = dateOutput[1]
			}

			txns, err := wallet.GetTransaction(tx.Hash)
			if err != nil {
				return fmt.Errorf("Error: %s",err.Error())
			}

			var txStatus string
			if txns.Confirmations > confirmations{
				txStatus = "Confirmed"
			}else{
				txStatus = "Unconfirmed"
			}

			pageTxRows[i] = []interface{}{
				txDate,
				tx.Direction,
				tx.Amount,
				txStatus,
				tx.Type,
			}
		}
		termio.PrintTabularResult(termio.StdoutWriter, columns, pageTxRows)

		// ask user what to do next
		var prompt string
		pageInfo := fmt.Sprintf("Showing transactions %d-%d", lastTxRowNumber, lastTxRowNumber+len(transactions)-1)
		if startBlockHeight >= 0 {
			prompt = fmt.Sprintf("%s, enter # for details, show (m)ore, or (q)uit", pageInfo)
		} else {
			prompt = fmt.Sprintf("%s, enter # for details or (q)uit", pageInfo)
		}

		validateUserInput := func(userInput string) error {
			if strings.EqualFold(userInput, "q") ||
				(strings.EqualFold(userInput, "m") && startBlockHeight >= 0) {
				return nil
			}

			// check if user input is a valid tx #
			txRowNumber, err := strconv.ParseUint(userInput, 10, 32)
			if err != nil || txRowNumber < 1 || int(txRowNumber) > len(displayedTxHashes) {
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
		} else if strings.EqualFold(userChoice, "m") {
			fmt.Println() // print empty line before listing txs for next page
			continue
		}

		// if the code execution continues to this point, it means user's response was neither "q" nor "m"
		// must therefore be a tx # to view tx details
		txRowNumber, _ := strconv.ParseUint(userChoice, 10, 32)
		txHash := displayedTxHashes[txRowNumber-1]

		showTransactionCommandArgs := ShowTransactionCommandArgs{txHash}
		showTxDetails := ShowTransactionCommand{
			Args:     showTransactionCommandArgs,
			Detailed: true,
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
