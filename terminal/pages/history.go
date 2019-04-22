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

	displayMessage := func(message string) {
		// this function may be called from a goroutine, use tviewApp.QueueUpdateDraw
		tviewApp.QueueUpdateDraw(func() {
			body.RemoveItem(errorTextView)
			if message != "" {
				errorTextView.SetText(message)
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

		fetchTransactionDetail(txHash, wallet, displayMessage, transactionDetailsTable)
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
	historyTable.SetCell(0, 3, tableHeaderCell("Status"))
	historyTable.SetCell(0, 4, tableHeaderCell("Type"))

	displayHistoryTable()

	// fetch tx to display in subroutine so the UI isn't blocked
	go fetchAndDisplayTransactions(-1, wallet, historyTable, tviewApp, displayMessage)

	hintTextView.SetText("TIP: Use ARROW UP/DOWN to select txn, ENTER to view details, ESC to return to navigation menu")

	tviewApp.SetFocus(body)

	return body
}

func fetchAndDisplayTransactions(startBlockHeight int32, wallet walletcore.Wallet, historyTable *tview.Table, tviewApp *tview.Application,
	displayMessage func(string)) {

	// show a loading text at the bottom of the table so user knows an op is in progress
	displayMessage("Fetching data...")

	txns, endBlockHeight, err := wallet.TransactionHistory(context.Background(), startBlockHeight, walletcore.TransactionHistoryCountPerPage)
	if err != nil {
		displayMessage(err.Error())
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
			nextRowIndex := historyTable.GetRowCount()

			historyTable.SetCell(nextRowIndex, 0, tview.NewTableCell(tx.FormattedTime).SetAlign(tview.AlignCenter).SetMaxWidth(1).SetExpansion(1).SetMaxWidth(1).SetExpansion(1))
			historyTable.SetCell(nextRowIndex, 1, tview.NewTableCell(fmt.Sprintf("%10s", tx.Direction.String())).SetAlign(tview.AlignCenter).SetMaxWidth(2).SetExpansion(1))
			historyTable.SetCell(nextRowIndex, 2, tview.NewTableCell(fmt.Sprintf("%15s", formatAmount(tx.RawAmount))).SetAlign(tview.AlignCenter).SetMaxWidth(3).SetExpansion(1))
			historyTable.SetCell(nextRowIndex, 3, tview.NewTableCell(fmt.Sprintf("%12s", tx.Status)).SetAlign(tview.AlignCenter).SetMaxWidth(1).SetExpansion(1))
			historyTable.SetCell(nextRowIndex, 4, tview.NewTableCell(fmt.Sprintf("%7s", tx.Type)).SetAlign(tview.AlignCenter).SetMaxWidth(1).SetExpansion(1))

			displayedTxHashes = append(displayedTxHashes, tx.Hash)
		}

		// clear loading message text
		displayMessage("")
	})

	if endBlockHeight > 0 {
		// set or reset selection changed listener to load more data when the table is almost scrolled to the end
		historyTable.SetSelectionChangedFunc(func(row, column int) {
			if row >= historyTable.GetRowCount()-10 {
				historyTable.SetSelectionChangedFunc(nil) // unset selection change listener until table is populated
				fetchAndDisplayTransactions(endBlockHeight-1, wallet, historyTable, tviewApp, displayMessage)
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
	transactionDetailsTable.SetCellSimple(5, 0,  "Date")
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
	inputsAndOutputsAmount := make([]int64, 0, len(tx.Inputs)+len(tx.Outputs))
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
