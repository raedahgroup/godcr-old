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
	titleTextView := primitives.TitleTextView("")
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
	clearHistoryPage := func(hintText, historyTable, transactionDetailsTable, titleTextView tview.Primitive) {
		body.RemoveItem(historyTable)
		body.RemoveItem(transactionDetailsTable)
		body.RemoveItem(titleTextView)
		body.RemoveItem(hintText)
	}

	// Table header
	detailedTableHeaderCell := func(text string) *tview.TableCell {
		return tview.NewTableCell(text).SetAlign(tview.AlignLeft).SetSelectable(false)
	}

	historyTable.SetSelectedFunc(func(row, column int) {
		clearHistoryPage(hintText, historyTable, transactionDetailsTable, titleTextView)
		setFocus(transactionDetailsTable)
		txHash := historyTable.GetCell(row, 6).Text

		body.AddItem(titleTextView.SetText("Transaction Detail"), 1, 0, false)
		body.AddItem(hintText.SetText("(TIP: Use ARROW UP/DOWN to scroll, BACKSPACE to view History page, ESC to return to nav menu)"), 3, 0, false)

		transactionDetailsTable.SetCell(0, 0, detailedTableHeaderCell("Date"))
		transactionDetailsTable.SetCell(1, 0, detailedTableHeaderCell("Direction"))
		transactionDetailsTable.SetCell(2, 0, detailedTableHeaderCell("Amount"))
		transactionDetailsTable.SetCell(3, 0, detailedTableHeaderCell("Fee"))
		transactionDetailsTable.SetCell(4, 0, detailedTableHeaderCell("Rate"))
		transactionDetailsTable.SetCell(5, 0, detailedTableHeaderCell("Type"))
		transactionDetailsTable.SetCell(6, 0, detailedTableHeaderCell("Confirmation"))
		transactionDetailsTable.SetCell(7, 0, detailedTableHeaderCell("Included in block \n \n \t"))
		transactionDetailsTable.SetCell(8, 0, detailedTableHeaderCell("Hash "))
		transactionDetailsTable.SetCell(9, 0, detailedTableHeaderCell(" "))
		transactionDetailsTable.SetCell(10, 0, detailedTableHeaderCell("Input"))

		body.AddItem(transactionDetailsTable, 0, 1, true)

		fetchTransactionDetail(txHash, wallet, displayError, transactionDetailsTable)
	})

	var history func()
	transactionDetailsTable.SetInputCapture(func(event *tcell.EventKey) *tcell.EventKey {
		if event.Key() == tcell.KeyEscape {
			clearFocus()
			return nil
		}
		if event.Key() == tcell.KeyBackspace {
			clearHistoryPage(hintText, historyTable, transactionDetailsTable, titleTextView)
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
		clearHistoryPage(hintText, historyTable, transactionDetailsTable, titleTextView)

		body.AddItem(titleTextView.SetText("History"), 1, 0, false)
		body.AddItem(hintText.SetText("(TIP: Use ARROW UP/DOWN to select txn, ENTER to view details, ESC to return to nav menu)"), 3, 0, false)

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

		fetchAndDisplayTransactions(-1, wallet, historyTable, displayError)
			body.AddItem(historyTable, 0, 1, true)

	}

	history()

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

	transactionDetailsTable.SetCell(0, 1, tview.NewTableCell(tx.FormattedTime).SetAlign(tview.AlignLeft))
	transactionDetailsTable.SetCell(1, 1, tview.NewTableCell(tx.Direction.String()).SetAlign(tview.AlignLeft))
	transactionDetailsTable.SetCell(2, 1, tview.NewTableCell(tx.Amount).SetAlign(tview.AlignLeft))
	transactionDetailsTable.SetCell(3, 1, tview.NewTableCell(tx.Fee).SetAlign(tview.AlignLeft))
	transactionDetailsTable.SetCell(4, 1, tview.NewTableCell(fmt.Sprintf("%s/kB\n", tx.FeeRate)).SetAlign(tview.AlignLeft))
	transactionDetailsTable.SetCell(5, 1, tview.NewTableCell(tx.Type).SetAlign(tview.AlignLeft))
	transactionDetailsTable.SetCell(6, 1, tview.NewTableCell(strconv.Itoa(int(tx.Confirmations))).SetAlign(tview.AlignLeft))
	transactionDetailsTable.SetCell(7, 1, tview.NewTableCell(strconv.Itoa(int(tx.BlockHeight))).SetAlign(tview.AlignLeft))
	transactionDetailsTable.SetCell(8, 1, tview.NewTableCell(tx.Hash).SetAlign(tview.AlignLeft))
	for _, txIn := range tx.Inputs {
		row := transactionDetailsTable.GetRowCount()
		transactionDetailsTable.SetCell(row, 0, tview.NewTableCell(dcrutil.Amount(txIn.AmountIn).String()).SetAlign(tview.AlignLeft))
		transactionDetailsTable.SetCell(row, 1, tview.NewTableCell(txIn.PreviousOutpoint).SetAlign(tview.AlignLeft))
	}

	row := transactionDetailsTable.GetRowCount()
	transactionDetailsTable.SetCell(row, 0, tview.NewTableCell(" ").SetAlign(tview.AlignLeft))
	transactionDetailsTable.SetCell(row, 0, tview.NewTableCell("Output").SetAlign(tview.AlignLeft))
	for _, txOut := range tx.Outputs {
		row := transactionDetailsTable.GetRowCount()
		if len(txOut.Addresses) == 0 {
			transactionDetailsTable.SetCell(row, 0, tview.NewTableCell(fmt.Sprintf("  %s \t (no address)\n", dcrutil.Amount(txOut.Value).String())).SetAlign(tview.AlignLeft))
			continue
		}

		outputAmount := dcrutil.Amount(txOut.Value).String()
		for _, address := range txOut.Addresses {
			accountName := address.AccountName
			if !address.IsMine {
				accountName = "external"
			}

			transactionDetailsTable.SetCell(row, 0, tview.NewTableCell(outputAmount).SetAlign(tview.AlignLeft))
			transactionDetailsTable.SetCell(row, 1, tview.NewTableCell(fmt.Sprintf("%s (%s)\n", address.Address, accountName)).SetAlign(tview.AlignLeft))
		}
	}
}
