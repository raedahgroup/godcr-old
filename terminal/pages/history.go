package pages

import (
	"context"
	"fmt"
	"strconv"

	"github.com/decred/dcrd/dcrutil"
	"github.com/gdamore/tcell"
	"github.com/raedahgroup/godcr/app/walletcore"
	"github.com/raedahgroup/godcr/terminal/primitives"
	"github.com/rivo/tview"
)

func historyPage(wallet walletcore.Wallet, setFocus func(p tview.Primitive) *tview.Application, clearFocus func()) tview.Primitive {
	// parent flexbox layout container to hold other primitives
	body := tview.NewFlex().SetDirection(tview.FlexRow)

	// page title and tip
	titleTextView := primitives.NewCenterAlignedTextView("")
	body.AddItem(titleTextView.SetText("History"), 1, 0, false)

	hintText := primitives.WordWrappedTextView("")
	hintText.SetTextColor(tcell.ColorGray)
	body.AddItem(hintText.SetText("(TIP: Use ARROW UP/DOWN to select txn, ENTER to view details, ESC to return to nav menu)"), 3, 0, false)

	historyTable := tview.NewTable().
		SetBorders(false).
		SetSeparator('\t').
		SetFixed(1, 0).
		SetSelectable(true, false)

	transactionDetailsTable := tview.NewTable().
		SetBorders(false).
		SetSeparator(' ').
		SetSelectable(false, false)

	errorTextView := primitives.WordWrappedTextView("")
	errorTextView.SetTextColor(tcell.ColorOrangeRed)

	displayError := func(errorMessage string) {
		body.RemoveItem(errorTextView)
		errorTextView.SetText(errorMessage)
		body.AddItem(errorTextView, 2, 0, false)
	}

	// clearHistoryPage clear the screen before outputing new data
	clearHistoryPage := func() {
		body.RemoveItem(historyTable)
		body.RemoveItem(transactionDetailsTable)
	}

	// method for getting transaction details
	historyTable.SetSelectedFunc(func(row, column int) {
		clearHistoryPage()
		txHash := historyTable.GetCell(row, 6).Text

		titleTextView.SetText("Transaction Detail")
		hintText.SetText("(TIP: Use ARROW UP/DOWN to scroll, BACKSPACE to view History page, ESC to return to nav menu)")
		
		body.AddItem(transactionDetailsTable, 0, 1, true)

		setFocus(transactionDetailsTable)
		fetchTransactionDetail(txHash, wallet, displayError, transactionDetailsTable)
	})

	var history func()
	// hander for setting focus on transaction Details Table and returning back to history table
	transactionDetailsTable.SetInputCapture(func(event *tcell.EventKey) *tcell.EventKey {
		if event.Key() == tcell.KeyEscape {
			clearFocus()
			return nil
		}
		if event.Key() == tcell.KeyBackspace {
			clearHistoryPage()
			history()
			setFocus(historyTable)
			return nil
		}

		return event
	})

	historyTable.SetDoneFunc(func(key tcell.Key) {
		if key == tcell.KeyEscape {
			clearFocus()
		}
	})

	history = func() {
		clearHistoryPage()

		titleTextView.SetText("History")
		hintText.SetText("(TIP: Use ARROW UP/DOWN to select txn, ENTER to view details, ESC to return to nav menu)")

		tableHeaderCell := func(text string) *tview.TableCell {
			return tview.NewTableCell(text).SetAlign(tview.AlignCenter).SetSelectable(false)
		}

		// Table header
		historyTable.SetCell(0, 0, tableHeaderCell("#"))
		historyTable.SetCell(0, 1, tableHeaderCell("Date"))
		historyTable.SetCell(0, 4, tableHeaderCell("Direction"))
		historyTable.SetCell(0, 2, tableHeaderCell("Amount"))
		historyTable.SetCell(0, 3, tableHeaderCell("Fee"))
		historyTable.SetCell(0, 5, tableHeaderCell("Type"))
		historyTable.SetCell(0, 6, tableHeaderCell("Hash"))

		body.AddItem(historyTable, 0, 1, true)
	}

	history()
	fetchAndDisplayTransactions(-1, wallet, historyTable, displayError)

	body.AddItem(nil, 1, 0, false) // add some "padding" at the bottom
	setFocus(body)

	return body
}

func fetchAndDisplayTransactions(startBlockHeight int32, wallet walletcore.Wallet, historyTable *tview.Table, displayError func(errorMessage string)) {
	txns, endBlockHeight, err := wallet.TransactionHistory(context.Background(), startBlockHeight, walletcore.TransactionHistoryCountPerPage)
	if err != nil {
		displayError(err.Error())
		return
	}

	for _, tx := range txns {
		row := historyTable.GetRowCount()
		historyTable.SetCellSimple(row, 0, fmt.Sprintf("%d.", row))
		historyTable.SetCell(row, 1, tview.NewTableCell(tx.FormattedTime).SetAlign(tview.AlignCenter))
		historyTable.SetCell(row, 4, tview.NewTableCell(tx.Direction.String()).SetAlign(tview.AlignCenter))
		historyTable.SetCell(row, 2, tview.NewTableCell(tx.Amount).SetAlign(tview.AlignRight))
		historyTable.SetCell(row, 3, tview.NewTableCell(tx.Fee).SetAlign(tview.AlignRight))
		historyTable.SetCell(row, 5, tview.NewTableCell(tx.Type).SetAlign(tview.AlignCenter))
		historyTable.SetCell(row, 6, tview.NewTableCell(tx.Hash).SetAlign(tview.AlignCenter))
	}

	if endBlockHeight > 0 {
		// set or reset selection changed listener to load more data when the table is almost scrolled to the end
		historyTable.SetSelectionChangedFunc(func(row, column int) {
			if row >= historyTable.GetRowCount()-10 {
				historyTable.SetSelectionChangedFunc(nil) // unset selection change listener until table is populated
				go fetchAndDisplayTransactions(endBlockHeight-1, wallet, historyTable, displayError)
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

	transactionDetailsTable.Clear()
	transactionDetailsTable.SetCellSimple(0, 0, "Hash ")
	transactionDetailsTable.SetCellSimple(1, 0, "Confirmations")
	transactionDetailsTable.SetCellSimple(2, 0, "Included in block  ")
	transactionDetailsTable.SetCellSimple(3, 0, "Type")
	transactionDetailsTable.SetCellSimple(4, 0, "Amount received")
	transactionDetailsTable.SetCellSimple(5, 0, "Date")
	transactionDetailsTable.SetCellSimple(6, 0, "Direction")
	transactionDetailsTable.SetCellSimple(7, 0, "Fee")
	transactionDetailsTable.SetCellSimple(8, 0, "Fee Rate")
	transactionDetailsTable.SetCellSimple(9, 0, "-Inputs-")

	transactionDetailsTable.SetCellSimple(0, 1, tx.Hash)
	transactionDetailsTable.SetCellSimple(1, 1, strconv.Itoa(int(tx.Confirmations)))
	transactionDetailsTable.SetCellSimple(2, 1, strconv.Itoa(int(tx.BlockHeight)))
	transactionDetailsTable.SetCellSimple(3, 1, tx.Type)
	transactionDetailsTable.SetCellSimple(4, 1, tx.Amount)
	transactionDetailsTable.SetCellSimple(5, 1, tx.FormattedTime)
	transactionDetailsTable.SetCellSimple(6, 1, tx.Direction.String())
	transactionDetailsTable.SetCellSimple(7, 1, tx.Fee)
	transactionDetailsTable.SetCellSimple(8, 1, fmt.Sprintf("%s/kB", tx.FeeRate))
	
	for _, txIn := range tx.Inputs {
		row := transactionDetailsTable.GetRowCount()
		transactionDetailsTable.SetCellSimple(row, 0, dcrutil.Amount(txIn.AmountIn).String())
		transactionDetailsTable.SetCellSimple(row, 1, txIn.PreviousOutpoint)
	}

	row := transactionDetailsTable.GetRowCount()
	transactionDetailsTable.SetCellSimple(row, 0, "-Outputs-")
	for _, txOut := range tx.Outputs {
		row := transactionDetailsTable.GetRowCount()
		if len(txOut.Addresses) == 0 {
			transactionDetailsTable.SetCellSimple(row, 0, fmt.Sprintf("  %s (no address)", dcrutil.Amount(txOut.Value).String()))
			continue
		}

		outputAmount := dcrutil.Amount(txOut.Value).String()
		for _, address := range txOut.Addresses {
			accountName := address.AccountName
			if !address.IsMine {
				accountName = "external"
			}

			transactionDetailsTable.SetCellSimple(row, 0, outputAmount)
			transactionDetailsTable.SetCellSimple(row, 1, fmt.Sprintf("%s (%s)", address.Address, accountName))
		}
	}
}
