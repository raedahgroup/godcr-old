package pages

import (
	"github.com/raedahgroup/godcr/app"
	"github.com/raedahgroup/godcr/app/walletcore"
	"github.com/rivo/tview"
)

func HistoryPage(walletMiddleware app.WalletMiddleware) tview.Primitive {
	body := tview.NewTable().SetBorders(true)
	transactions, err := walletMiddleware.TransactionHistory()
	if err != nil {
		fmt.Println(err)
	}

	body.SetCell(0, 0, tview.NewTableCell("Date").SetAlign(tview.AlignCenter))
	body.SetCell(0, 1, tview.NewTableCell("Amount").SetAlign(tview.AlignCenter))
	body.SetCell(0, 2, tview.NewTableCell("Fee").SetAlign(tview.AlignCenter))
	body.SetCell(0, 3, tview.NewTableCell("Direction").SetAlign(tview.AlignCenter))
	body.SetCell(0, 4, tview.NewTableCell("Type").SetAlign(tview.AlignCenter))
	body.SetCell(0, 5, tview.NewTableCell("Hash").SetAlign(tview.AlignCenter))


	for i, tx := range transactions {
		row := i + 1
		body.SetCell(row, 0, tview.NewTableCell(tx.FormattedTime).SetAlign(tview.AlignCenter))
		body.SetCell(row, 1, tview.NewTableCell(walletcore.AmountToString(tx.Amount.ToCoin())).SetAlign(tview.AlignCenter))
		body.SetCell(row, 2, tview.NewTableCell(walletcore.AmountToString(tx.Fee.ToCoin())).SetAlign(tview.AlignCenter))
		body.SetCell(row, 3, tview.NewTableCell(tx.Direction.String()).SetAlign(tview.AlignCenter))
		body.SetCell(row, 4, tview.NewTableCell(tx.Type).SetAlign(tview.AlignCenter))
		body.SetCell(row, 5, tview.NewTableCell(tx.Hash).SetAlign(tview.AlignCenter))
	}
	return body
}