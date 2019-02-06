package pages

import (
	"fmt"
	"math"

	"github.com/rivo/tview"
	"github.com/raedahgroup/godcr/app/walletcore"
)

func HistoryPage(wallet walletcore.Wallet) tview.Primitive {
	body := tview.NewTable().SetBorders(true)
	transactions, err := wallet.TransactionHistory()
	if err != nil {
		panic(err) 
	}

	body.SetCell(0, 0, tview.NewTableCell("Date").SetAlign(tview.AlignCenter).SetBackgroundColor(color))
	body.SetCell(0, 1, tview.NewTableCell("Amount").SetAlign(tview.AlignCenter))
	body.SetCell(0, 2, tview.NewTableCell("Fee").SetAlign(tview.AlignCenter))
	body.SetCell(0, 3, tview.NewTableCell("Direction").SetAlign(tview.AlignCenter))
	body.SetCell(0, 4, tview.NewTableCell("Type").SetAlign(tview.AlignCenter))
	body.SetCell(0, 5, tview.NewTableCell("Hash").SetAlign(tview.AlignCenter))

	for i , tx := range transactions {
		row = i + 1
		body.SetCell(row, 0, tview.NewTableCell(tx.FormattedTime).SetAlign(tview.AlignCenter))
		body.SetCell(row, 1, tview.NewTableCell(AmountToString(tx.Amount.ToCoin())).SetAlign(tview.AlignCenter))
		body.SetCell(row, 2, tview.NewTableCell(AmountToString(tx.Fee.ToCoin())).SetAlign(tview.AlignCenter))
		body.SetCell(row, 3, tview.NewTableCell(tx.Direction.String()).SetAlign(tview.AlignCenter))
		body.SetCell(row, 4, tview.NewTableCell(tx.Type).SetAlign(tview.AlignCenter))
		body.SetCell(row, 5, tview.NewTableCell(tx.Hash).SetAlign(tview.AlignCenter))
	}
	return body
}

func AmountToString(amount float64) string {
	amount = math.Round(amount)
	return fmt.Sprintf("%d DCR", int(amount))
}