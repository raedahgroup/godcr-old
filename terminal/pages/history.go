package pages

import (
	"fmt"
	"strconv"

	"github.com/decred/dcrd/dcrutil"
	"github.com/gdamore/tcell"
	"github.com/raedahgroup/dcrlibwallet/utils"
	"github.com/raedahgroup/godcr/app/walletcore"
	"github.com/raedahgroup/godcr/terminal/helpers"
	"github.com/raedahgroup/godcr/terminal/primitives"
	"github.com/rivo/tview"
)

var displayedTxHashes []string
var txPerPage int32 = walletcore.TransactionHistoryCountPerPage
var totalTxCount int

func historyPage(wallet walletcore.Wallet, hintTextView *primitives.TextView, tviewApp *tview.Application, clearFocus func()) tview.Primitive {
	// parent flexbox layout container to hold other primitives
	body := tview.NewFlex().SetDirection(tview.FlexRow)

	// handler for returning back to menu column
	body.SetInputCapture(func(event *tcell.EventKey) *tcell.EventKey {
		if event.Key() == tcell.KeyEscape || event.Key() == tcell.KeyBackspace || event.Key() == tcell.KeyBackspace2 {
			clearFocus()
			return nil
		}

		return event
	})

	messageTextView := primitives.WordWrappedTextView("")

	displayMessage := func(message string, error bool) {
		// this function may be called from a goroutine, use tviewApp.QueueUpdateDraw
		tviewApp.QueueUpdateDraw(func() {
			body.RemoveItem(messageTextView)
			messageTextView.SetText(message)
			if error && message != "" {
				messageTextView.SetTextColor(helpers.DecredOrangeColor)
			}

			body.AddItem(messageTextView, 2, 0, false)
		})
	}

	// get total tx count early on, so as to display on history header
	txCount, txCountErr := wallet.TransactionCount(nil)
	totalTxCount = txCount

	// page title
	historyPageTitle := fmt.Sprintf("History (%d transactions)", txCount)
	titleTextView := primitives.NewLeftAlignedTextView(historyPageTitle)
	body.AddItem(titleTextView, 2, 0, false)

	if txCountErr != nil {
		displayMessage(fmt.Sprintf("Cannot load history. Get total tx count error: %s", txCountErr.Error()), true)
		tviewApp.SetFocus(body)
		return body
	}
	if txCount == 0 {
		displayMessage("No transactions yet", false)
		tviewApp.SetFocus(body)
		return body
	}

	historyTable := tview.NewTable().
		SetBorders(false).
		SetFixed(1, 0).
		SetSelectable(true, false)

	transactionDetailsTable := tview.NewTable().SetBorders(false)

	displayHistoryTable := func() {
		body.RemoveItem(transactionDetailsTable)

		titleTextView.SetText(historyPageTitle)
		hintTextView.SetText("TIP: Use ARROW UP/DOWN to select txn,\nENTER to view details, ESC to return to navigation menu")

		body.AddItem(historyTable, 0, 1, true)
		tviewApp.SetFocus(historyTable)
	}

	historyTable.SetDoneFunc(func(key tcell.Key) {
		if key == tcell.KeyEscape {
			clearFocus()
		}
	})

	displayedTxHashes = []string{}

	// method for getting transaction details when a tx is selected from the history table
	historyTable.SetSelectedFunc(func(row, column int) {
		if row >= len(displayedTxHashes) {
			// ignore selected func call for table header
			return
		}

		body.RemoveItem(historyTable)
		txHash := displayedTxHashes[row-1]

		titleTextView.SetText("Transaction Details")
		hintTextView.SetText("TIP: Use ARROW UP/DOWN to scroll, \nBACKSPACE to view History page, ESC to return to navigation menu")

		transactionDetailsTable.Clear()
		body.AddItem(transactionDetailsTable, 0, 1, true)

		tviewApp.SetFocus(transactionDetailsTable)

		displayTxDetails(txHash, wallet, displayMessage, transactionDetailsTable)
	})

	// handler for returning back to history table
	transactionDetailsTable.SetInputCapture(func(event *tcell.EventKey) *tcell.EventKey {
		if event.Key() == tcell.KeyEscape || event.Key() == tcell.KeyBackspace || event.Key() == tcell.KeyBackspace2 {
			displayHistoryTable()
			return nil
		}

		return event
	})

	tableHeaderCell := func(text string) *tview.TableCell {
		return tview.NewTableCell(text).SetAlign(tview.AlignCenter).SetSelectable(false).SetMaxWidth(1).SetExpansion(1)
	}

	// history table header
	historyTable.SetCell(0, 0, tableHeaderCell("Date (UTC)"))
	historyTable.SetCell(0, 1, tableHeaderCell(fmt.Sprintf("%10s", "Direction")))
	historyTable.SetCell(0, 2, tableHeaderCell(fmt.Sprintf("%8s", "Amount")))
	historyTable.SetCell(0, 3, tableHeaderCell(fmt.Sprintf("%5s", "Status")))
	historyTable.SetCell(0, 4, tableHeaderCell(fmt.Sprintf("%-5s", "Type")))

	displayHistoryTable()

	// fetch tx to display in subroutine so the UI isn't blocked
	go fetchAndDisplayTransactions(0, wallet, historyTable, tviewApp, displayMessage)

	hintTextView.SetText("TIP: Use ARROW UP/DOWN to select txn, \nENTER to view details, ESC to return to navigation menu")

	tviewApp.SetFocus(body)

	return body
}

func fetchAndDisplayTransactions(txOffset int, wallet walletcore.Wallet, historyTable *tview.Table, tviewApp *tview.Application, displayMessage func(string, bool)) {

	// show a loading text at the bottom of the table so user knows an op is in progress
	displayMessage("Fetching data...", false)

	txns, err := wallet.TransactionHistory(int32(txOffset), txPerPage, nil)
	if err != nil {
		displayMessage(err.Error(), true)
		return
	}

	// calculate max number of digits after decimal point for all tx amounts
	inputsAndOutputsAmount := make([]int64, len(txns))
	for i, tx := range txns {
		inputsAndOutputsAmount[i] = tx.Amount
	}
	maxDecimalPlacesForTxAmounts := maxDecimalPlaces(inputsAndOutputsAmount)

	// now format amount having determined the max number of decimal places
	formatAmount := func(amount int64) string {
		return formatAmountDisplay(amount, maxDecimalPlacesForTxAmounts)
	}

	// updating the history table from a goroutine, use tviewApp.QueueUpdateDraw
	tviewApp.QueueUpdateDraw(func() {
		for _, tx := range txns {
			nextRowIndex := historyTable.GetRowCount()

			historyTable.SetCell(nextRowIndex, 0, tview.NewTableCell(fmt.Sprintf("%-10s", utils.ExtractDateOrTime(tx.Timestamp))).SetAlign(tview.AlignCenter).SetMaxWidth(1).SetExpansion(1).SetMaxWidth(1).SetExpansion(1))
			historyTable.SetCell(nextRowIndex, 1, tview.NewTableCell(fmt.Sprintf("%-10s", tx.Direction.String())).SetAlign(tview.AlignCenter).SetMaxWidth(2).SetExpansion(1))
			historyTable.SetCell(nextRowIndex, 2, tview.NewTableCell(fmt.Sprintf("%15s", formatAmount(tx.Amount))).SetAlign(tview.AlignCenter).SetMaxWidth(3).SetExpansion(1))
			historyTable.SetCell(nextRowIndex, 3, tview.NewTableCell(fmt.Sprintf("%12s", tx.Status)).SetAlign(tview.AlignCenter).SetMaxWidth(1).SetExpansion(1))
			historyTable.SetCell(nextRowIndex, 4, tview.NewTableCell(fmt.Sprintf("%-8s", tx.Type)).SetAlign(tview.AlignCenter).SetMaxWidth(1).SetExpansion(1))

			displayedTxHashes = append(displayedTxHashes, tx.Hash)
		}

		// clear loading message text
		displayMessage("", false)
	})

	if len(displayedTxHashes) < totalTxCount {
		// set or reset selection changed listener to load more data when the table is almost scrolled to the end
		nextOffset := txOffset + len(txns)
		historyTable.SetSelectionChangedFunc(func(row, column int) {
			if row >= historyTable.GetRowCount()-10 {
				historyTable.SetSelectionChangedFunc(nil) // unset selection change listener until table is populated
				fetchAndDisplayTransactions(nextOffset, wallet, historyTable, tviewApp, displayMessage)
			}
		})
	}

	return
}

func displayTxDetails(txHash string, wallet walletcore.Wallet, displayError func(string, bool), transactionDetailsTable *tview.Table) {
	tx, err := wallet.GetTransaction(txHash)
	if err != nil {
		displayError(err.Error(), true)
		return
	}

	transactionDetailsTable.SetCellSimple(0, 0, "Hash")
	transactionDetailsTable.SetCellSimple(1, 0, "Confirmations")
	transactionDetailsTable.SetCellSimple(2, 0, "Included in block")
	transactionDetailsTable.SetCellSimple(3, 0, "Type")
	transactionDetailsTable.SetCellSimple(4, 0, "Amount")
	transactionDetailsTable.SetCellSimple(5, 0, "Date")
	transactionDetailsTable.SetCellSimple(6, 0, "Direction")
	transactionDetailsTable.SetCellSimple(7, 0, "Fee")
	transactionDetailsTable.SetCellSimple(8, 0, "Fee Rate")

	transactionDetailsTable.SetCellSimple(0, 1, tx.Hash)
	transactionDetailsTable.SetCellSimple(1, 1, strconv.Itoa(int(tx.Confirmations)))
	transactionDetailsTable.SetCellSimple(2, 1, strconv.Itoa(int(tx.BlockHeight)))
	transactionDetailsTable.SetCellSimple(3, 1, tx.Type)
	transactionDetailsTable.SetCellSimple(4, 1, dcrutil.Amount(tx.Amount).String())
	transactionDetailsTable.SetCellSimple(5, 1, fmt.Sprintf("%s UTC", tx.LongTime))
	transactionDetailsTable.SetCellSimple(6, 1, tx.Direction.String())
	transactionDetailsTable.SetCellSimple(7, 1, dcrutil.Amount(tx.Fee).String())
	transactionDetailsTable.SetCellSimple(8, 1, fmt.Sprintf("%s/kB", dcrutil.Amount(tx.FeeRate)))

	// calculate max number of digits after decimal point for inputs and outputs
	inputsAndOutputsAmount := make([]int64, 0, len(tx.Inputs)+len(tx.Outputs))
	for _, txIn := range tx.Inputs {
		inputsAndOutputsAmount = append(inputsAndOutputsAmount, txIn.Amount)
	}
	for _, txOut := range tx.Outputs {
		inputsAndOutputsAmount = append(inputsAndOutputsAmount, txOut.Amount)
	}
	maxDecimalPlacesForInputsAndOutputsAmounts := maxDecimalPlaces(inputsAndOutputsAmount)

	// now format amount having determined the max number of decimal places
	formatAmount := func(amount int64) string {
		return formatAmountDisplay(amount, maxDecimalPlacesForInputsAndOutputsAmounts)
	}

	transactionDetailsTable.SetCellSimple(9, 0, "-Inputs-")
	for _, txIn := range tx.Inputs {
		row := transactionDetailsTable.GetRowCount()
		transactionDetailsTable.SetCell(row, 0, tview.NewTableCell(formatAmount(txIn.Amount)).SetAlign(tview.AlignRight))
		transactionDetailsTable.SetCellSimple(row, 1, txIn.PreviousOutpoint)
	}

	row := transactionDetailsTable.GetRowCount()
	transactionDetailsTable.SetCellSimple(row, 0, "-Outputs-")
	for _, txOut := range tx.Outputs {
		row++
		outputAmount := formatAmount(txOut.Amount)

		if txOut.Address == "" {
			transactionDetailsTable.SetCellSimple(row, 0, fmt.Sprintf("  %s (no address)", outputAmount))
			continue
		}

		transactionDetailsTable.SetCell(row, 0, tview.NewTableCell(outputAmount).SetAlign(tview.AlignRight))
		transactionDetailsTable.SetCellSimple(row, 1, fmt.Sprintf("%s (%s)", txOut.Address, txOut.AccountName))
	}
}
