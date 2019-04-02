package pages

import (
	"context"
	"fmt"

	"github.com/gdamore/tcell"
	"github.com/raedahgroup/godcr/app/walletcore"
	"github.com/raedahgroup/godcr/terminal/primitives"
	"github.com/rivo/tview"
)

func historyPage(wallet walletcore.Wallet, setFocus func(p tview.Primitive) *tview.Application, clearFocus func()) tview.Primitive {
	// parent flexbox layout container to hold other primitives
	body := tview.NewFlex().SetDirection(tview.FlexRow)

	// page title and tip
	body.AddItem(primitives.NewCenterAlignedTextView("History"), 1, 0, false)
	hintText := primitives.WordWrappedTextView("(TIP: Use ARROW UP/DOWN to select txn, ENTER to view details, ESC to return to nav menu)")
	hintText.SetTextColor(tcell.ColorGray)
	body.AddItem(hintText, 2, 0, false)

	historyTable := tview.NewTable().
		SetBorders(false).
		SetSeparator('\t').
		SetFixed(1, 0).
		SetSelectable(true, false)

	historyTable.SetSelectedFunc(func(row, column int) {
		// todo show tx details for selected row
	})

	historyTable.SetDoneFunc(func(key tcell.Key) {
		if key == tcell.KeyEscape {
			clearFocus()
		}
	})

	// historyTable header
	tableHeaderCell := func(text string) *tview.TableCell {
		return tview.NewTableCell(text).SetAlign(tview.AlignCenter).SetSelectable(false)
	}
	historyTable.SetCell(0, 0, tableHeaderCell("#"))
	historyTable.SetCell(0, 1, tableHeaderCell("Date"))
	historyTable.SetCell(0, 4, tableHeaderCell("Direction"))
	historyTable.SetCell(0, 2, tableHeaderCell("Amount"))
	historyTable.SetCell(0, 3, tableHeaderCell("Fee"))
	historyTable.SetCell(0, 5, tableHeaderCell("Type"))
	historyTable.SetCell(0, 6, tableHeaderCell("Hash"))

	body.AddItem(historyTable, 0, 1, true)

	errorTextView := primitives.WordWrappedTextView("")
	errorTextView.SetTextColor(tcell.ColorOrangeRed)

	displayError := func(errorMessage string) {
		body.RemoveItem(errorTextView)
		errorTextView.SetText(errorMessage)
		body.AddItem(errorTextView, 2, 0, false)
	}

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
