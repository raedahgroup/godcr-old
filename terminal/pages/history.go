package pages

import (
	"github.com/gdamore/tcell"
	"github.com/raedahgroup/godcr/app/walletcore"
	"github.com/rivo/tview"
)

func HistoryPage(wallet walletcore.Wallet, setFocus func(p tview.Primitive) *tview.Application, clearFocus func()) tview.Primitive {
	errmsg := tview.NewTextView().SetTextAlign(tview.AlignCenter)
	transactions, err := wallet.TransactionHistory()
	if err != nil {
		return errmsg.SetText(err.Error())
	}

	body := tview.NewTable().SetBorders(true)
	body.SetDoneFunc(func(key tcell.Key) {
		if key == tcell.KeyEscape {
			body.ScrollToBeginning()
			clearFocus()
		}
	})
	body.SetCell(0, 0, tview.NewTableCell("Date").SetAlign(tview.AlignCenter))
	body.SetCell(0, 1, tview.NewTableCell("Amount").SetAlign(tview.AlignCenter))
	body.SetCell(0, 2, tview.NewTableCell("Fee").SetAlign(tview.AlignCenter))
	body.SetCell(0, 3, tview.NewTableCell("Direction").SetAlign(tview.AlignCenter))
	body.SetCell(0, 4, tview.NewTableCell("Type").SetAlign(tview.AlignCenter))
	body.SetCell(0, 5, tview.NewTableCell("Hash").SetAlign(tview.AlignCenter))

	for i, tx := range transactions {
		row := i + 1
		body.SetCell(row, 0, tview.NewTableCell(tx.FormattedTime).SetAlign(tview.AlignCenter))
		body.SetCell(row, 1, tview.NewTableCell(tx.Amount.String()).SetAlign(tview.AlignCenter))
		body.SetCell(row, 2, tview.NewTableCell(tx.Fee.String()).SetAlign(tview.AlignCenter))
		body.SetCell(row, 3, tview.NewTableCell(tx.Direction.String()).SetAlign(tview.AlignCenter))
		body.SetCell(row, 4, tview.NewTableCell(tx.Type).SetAlign(tview.AlignCenter))
		body.SetCell(row, 5, tview.NewTableCell(tx.Hash).SetAlign(tview.AlignCenter))
	}

	setFocus(body)


	return body
}
