package pages

import (
	"context"
	"fmt"
	"strconv"

	"github.com/decred/dcrd/dcrutil"
	"github.com/gdamore/tcell"
	"github.com/raedahgroup/godcr/app/walletcore"
	"github.com/raedahgroup/godcr/terminal/helpers"
	"github.com/raedahgroup/godcr/terminal/primitives"
	"github.com/rivo/tview"
)

var displayedTxHashes []string

func historyPage(wallet walletcore.Wallet, hintTextView *primitives.TextView, tviewApp *tview.Application, clearFocus func()) tview.Primitive {
	// parent flexbox layout container to hold other primitives
	body := tview.NewFlex().SetDirection(tview.FlexRow)

	// page title and tip
	titleTextView := primitives.NewLeftAlignedTextView("History")
	body.AddItem(titleTextView, 2, 0, false)

	historyTable := tview.NewTable().
		SetBorders(false).
		SetFixed(1, 0).
		SetSelectable(true, false)

	transactionDetailsTable := tview.NewTable().SetBorders(false)

	displayHistoryTable := func() {
		body.RemoveItem(transactionDetailsTable)

		titleTextView.SetText("History")
		hintTextView.SetText("TIP: Use ARROW UP/DOWN to select txn, ENTER to view details, ESC to return to navigation menu")

		body.AddItem(historyTable, 0, 1, true)
		tviewApp.SetFocus(historyTable)
	}

	errorTextView := primitives.WordWrappedTextView("")
	errorTextView.SetTextColor(helpers.DecredOrangeColor)

	displayError := func(errorMessage string) {
		// this function may be called from a goroutine, use tviewApp.QueueUpdateDraw
		tviewApp.QueueUpdateDraw(func() {
			body.RemoveItem(errorTextView)
			if errorMessage != "" {
				errorTextView.SetText(errorMessage)
				body.AddItem(errorTextView, 2, 0, false)
			}
		})
	}

	historyTable.SetDoneFunc(func(key tcell.Key) {
		if key == tcell.KeyEscape {
			clearFocus()
		}
	})

	displayedTxHashes = []string{}

	// method for getting transaction details when a tx is selected from the history table
	historyTable.SetSelectedFunc(func(row, column int) {
		body.RemoveItem(historyTable)
		txHash := displayedTxHashes[row-1]

		titleTextView.SetText("Transaction Details")
		hintTextView.SetText("TIP: Use ARROW UP/DOWN to scroll, BACKSPACE to view History page, ESC to return to navigation menu")

		transactionDetailsTable.Clear()
		body.AddItem(transactionDetailsTable, 0, 1, true)

		tviewApp.SetFocus(transactionDetailsTable)

		fetchTransactionDetail(txHash, wallet, displayError, transactionDetailsTable)
	})

	// handler for returning back to history table
	transactionDetailsTable.SetInputCapture(func(event *tcell.EventKey) *tcell.EventKey {
		if event.Key() == tcell.KeyEscape {
			clearFocus()
			return nil
		}
		if event.Key() == tcell.KeyBackspace || event.Key() == tcell.KeyBackspace2 {
			displayHistoryTable()
			return nil
		}

		return event
	})

	tableHeaderCell := func(text string) *tview.TableCell {
		return tview.NewTableCell(text).SetAlign(tview.AlignCenter).SetSelectable(false)
	}

	// history table header
	historyTable.SetCell(0, 0, tableHeaderCell("Date (UTC)"))
	historyTable.SetCell(0, 1, tableHeaderCell("Direction"))
	historyTable.SetCell(0, 2, tableHeaderCell("Amount"))
	historyTable.SetCell(0, 3, tableHeaderCell("Status"))
	historyTable.SetCell(0, 4, tableHeaderCell("Type"))

	displayHistoryTable()

	// fetch tx to display in subroutine so the UI isn't blocked
	go fetchAndDisplayTransactions(-1, wallet, historyTable, tviewApp, displayError)

	hintTextView.SetText("TIP: Use ARROW UP/DOWN to select txn, ENTER to view details, ESC to return to navigation menu")

	tviewApp.SetFocus(body)

	return body
}

func fetchAndDisplayTransactions(startBlockHeight int32, wallet walletcore.Wallet, historyTable *tview.Table, tviewApp *tview.Application,
	displayError func(errorMessage string)) {

	// show a loading text at the bottom of the table so user knows an op is in progress
	displayError("Fetching data...")

	txns, endBlockHeight, err := wallet.TransactionHistory(context.Background(), startBlockHeight, walletcore.TransactionHistoryCountPerPage)
	if err != nil {
		displayError(err.Error())
		return
	}

	// calculate max number of digits after decimal point for all tx amounts
	inputsAndOutputsAmount := make([]int64, len(txns))
	for i, tx := range txns {
		inputsAndOutputsAmount[i] = tx.RawAmount
	}
	maxDecimalPlacesForTxAmounts := maxDecimalPlaces(inputsAndOutputsAmount)

	// now format amount having determined the max number of decimal places
	formatAmount := func(amount int64) string {
		return formatAmountDisplay(amount, maxDecimalPlacesForTxAmounts)
	}

	// updating the history table from a goroutine, use tviewApp.QueueUpdateDraw
	tviewApp.QueueUpdateDraw(func() {
		for _, tx := range txns {
			row := historyTable.GetRowCount()

			historyTable.SetCell(row, 0, tview.NewTableCell(tx.FormattedTime).SetAlign(tview.AlignCenter))
			historyTable.SetCell(row, 1, tview.NewTableCell(tx.Direction.String()).SetAlign(tview.AlignCenter))
			historyTable.SetCell(row, 2, tview.NewTableCell(formatAmount(tx.RawAmount)).SetAlign(tview.AlignRight))
			historyTable.SetCell(row, 3, tview.NewTableCell(tx.Status).SetAlign(tview.AlignCenter))
			historyTable.SetCell(row, 4, tview.NewTableCell(tx.Type).SetAlign(tview.AlignCenter))

			displayedTxHashes = append(displayedTxHashes, tx.Hash)
		}

		// clear loading message text
		displayError("")
	})

	if endBlockHeight > 0 {
		// set or reset selection changed listener to load more data when the table is almost scrolled to the end
		historyTable.SetSelectionChangedFunc(func(row, column int) {
			if row >= historyTable.GetRowCount()-10 {
				historyTable.SetSelectionChangedFunc(nil) // unset selection change listener until table is populated
				fetchAndDisplayTransactions(endBlockHeight-1, wallet, historyTable, tviewApp, displayError)
			}
		})
	}

	return
}

func fetchTransactionDetail(txHash string, wallet walletcore.Wallet, displayError func(errorMessage string), transactionDetailsTable *tview.Table) {
	tx, err := wallet.GetTransaction(txHash)
	if err != nil {
		displayError(err.Error())
		return
	}

	transactionDetailsTable.SetCellSimple(0, 0, "Hash ")
	transactionDetailsTable.SetCellSimple(1, 0, "Confirmations")
	transactionDetailsTable.SetCellSimple(2, 0, "Included in block")
	transactionDetailsTable.SetCellSimple(3, 0, "Type")
	transactionDetailsTable.SetCellSimple(4, 0, "Amount received")
	transactionDetailsTable.SetCellSimple(5, 0, "Date")
	transactionDetailsTable.SetCellSimple(6, 0, "Direction")
	transactionDetailsTable.SetCellSimple(7, 0, "Fee")
	transactionDetailsTable.SetCellSimple(8, 0, "Fee Rate")

	transactionDetailsTable.SetCellSimple(0, 1, tx.Hash)
	transactionDetailsTable.SetCellSimple(1, 1, strconv.Itoa(int(tx.Confirmations)))
	transactionDetailsTable.SetCellSimple(2, 1, strconv.Itoa(int(tx.BlockHeight)))
	transactionDetailsTable.SetCellSimple(3, 1, tx.Type)
	transactionDetailsTable.SetCellSimple(4, 1, tx.Amount)
	transactionDetailsTable.SetCellSimple(5, 1, tx.FormattedTime)
	transactionDetailsTable.SetCellSimple(6, 1, tx.Direction.String())
	transactionDetailsTable.SetCellSimple(7, 1, tx.Fee)
	transactionDetailsTable.SetCellSimple(8, 1, fmt.Sprintf("%s/kB", tx.FeeRate))

	// calculate max number of digits after decimal point for inputs and outputs
	inputsAndOutputsAmount := make([]int64, 0, len(tx.Inputs) + len(tx.Outputs))
	for _, txIn := range tx.Inputs {
		inputsAndOutputsAmount = append(inputsAndOutputsAmount, txIn.AmountIn)
	}
	for _, txOut := range tx.Outputs {
		inputsAndOutputsAmount = append(inputsAndOutputsAmount, txOut.Value)
	}
	maxDecimalPlacesForInputsAndOutputsAmounts := maxDecimalPlaces(inputsAndOutputsAmount)

	// now format amount having determined the max number of decimal places
	formatAmount := func(amount int64) string {
		return formatAmountDisplay(amount, maxDecimalPlacesForInputsAndOutputsAmounts)
	}

	transactionDetailsTable.SetCellSimple(9, 0, "-Inputs-")
	for _, txIn := range tx.Inputs {
		row := transactionDetailsTable.GetRowCount()
		transactionDetailsTable.SetCell(row, 0, tview.NewTableCell(formatAmount(txIn.AmountIn)).SetAlign(tview.AlignRight))
		transactionDetailsTable.SetCellSimple(row, 1, txIn.PreviousOutpoint)
	}

	row := transactionDetailsTable.GetRowCount()
	transactionDetailsTable.SetCellSimple(row, 0, "-Outputs-")
	for _, txOut := range tx.Outputs {
		row++
		if len(txOut.Addresses) == 0 {
			transactionDetailsTable.SetCellSimple(row, 0, fmt.Sprintf("  %s (no address)", dcrutil.Amount(txOut.Value).String()))
			continue
		}

		outputAmount := formatAmount(txOut.Value)
		for _, address := range txOut.Addresses {
			accountName := address.AccountName
			if !address.IsMine {
				accountName = "external"
			}

			transactionDetailsTable.SetCell(row, 0, tview.NewTableCell(outputAmount).SetAlign(tview.AlignRight))
			transactionDetailsTable.SetCellSimple(row, 1, fmt.Sprintf("%s (%s)", address.Address, accountName))
		}
	}
}
