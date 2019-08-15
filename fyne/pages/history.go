package pages

import (
	"fmt"

	"fyne.io/fyne"
	"fyne.io/fyne/layout"
	"fyne.io/fyne/widget"
	"github.com/decred/dcrd/dcrutil"
	godcrApp "github.com/raedahgroup/godcr/app"
	"github.com/raedahgroup/godcr/fyne/widgets"
)

type historyPageData struct {
	txCount int
	txType  string
	txTable widgets.TableStruct
}

var history historyPageData

func historyUpdates(wallet godcrApp.WalletMiddleware) {
	txCount, _ := wallet.TransactionCount(nil)

	if txCount == history.txCount {
		return
	}
	val := txCount - history.txCount
	history.txCount = txCount

	//TODO: need to treat types
	txs, _ := wallet.TransactionHistory(0, int32(val), nil)

	var hBox []*widget.Box
	for _, tx := range txs {
		trimmedHash := tx.Hash[:len(tx.Hash)/2] + "..."
		hBox = append(hBox, widget.NewHBox(
			widget.NewLabelWithStyle(tx.LongTime, fyne.TextAlignCenter, fyne.TextStyle{}),
			widget.NewLabelWithStyle(tx.Type, fyne.TextAlignCenter, fyne.TextStyle{}),
			widget.NewLabelWithStyle(tx.Direction.String(), fyne.TextAlignLeading, fyne.TextStyle{}),
			widget.NewLabelWithStyle(dcrutil.Amount(tx.Amount).String(), fyne.TextAlignTrailing, fyne.TextStyle{}),
			widget.NewLabelWithStyle(dcrutil.Amount(tx.Fee).String(), fyne.TextAlignCenter, fyne.TextStyle{}),
			widget.NewLabelWithStyle(tx.Status, fyne.TextAlignCenter, fyne.TextStyle{}),
			widget.NewButton(trimmedHash, func() { fmt.Println("Hello") }),
		))
	}
	history.txTable.Prepend(hBox...)
}

func historyPage(wallet godcrApp.WalletMiddleware) fyne.CanvasObject {
	fetchHistoryTx(&history.txTable, wallet)
	//TODO: add dropdown to select history types
	return widget.NewHBox(widgets.NewHSpacer(10), fyne.NewContainerWithLayout(layout.NewFixedGridLayout(fyne.NewSize(history.txTable.Container.MinSize().Width, history.txTable.Container.MinSize().Height+600)), history.txTable.Container))
}

func fetchHistoryTx(txTable *widgets.TableStruct, wallet godcrApp.WalletMiddleware) {
	count, _ := wallet.TransactionCount(nil)
	history.txCount = count
	tx, _ := wallet.TransactionHistory(0, int32(count), nil)

	heading := widget.NewHBox(
		widget.NewLabelWithStyle("Date (UTC)", fyne.TextAlignCenter, fyne.TextStyle{Bold: true}),
		widget.NewLabelWithStyle("Type", fyne.TextAlignCenter, fyne.TextStyle{Bold: true}),
		widget.NewLabelWithStyle("Direction", fyne.TextAlignCenter, fyne.TextStyle{Bold: true}),
		widget.NewLabelWithStyle("Amount", fyne.TextAlignCenter, fyne.TextStyle{Bold: true}),
		widget.NewLabelWithStyle("Fee", fyne.TextAlignCenter, fyne.TextStyle{Bold: true}),
		widget.NewLabelWithStyle("Status", fyne.TextAlignCenter, fyne.TextStyle{Bold: true}),
		widget.NewLabelWithStyle("Hash", fyne.TextAlignCenter, fyne.TextStyle{Bold: true}))

	var hBox []*widget.Box
	for i := 0; i < count; i++ {
		trimmedHash := tx[i].Hash[:len(tx[i].Hash)/2] + "..."
		hBox = append(hBox, widget.NewHBox(
			widget.NewLabelWithStyle(tx[i].LongTime, fyne.TextAlignCenter, fyne.TextStyle{}),
			widget.NewLabelWithStyle(tx[i].Type, fyne.TextAlignCenter, fyne.TextStyle{}),
			widget.NewLabelWithStyle(tx[i].Direction.String(), fyne.TextAlignLeading, fyne.TextStyle{}),
			widget.NewLabelWithStyle(dcrutil.Amount(tx[i].Amount).String(), fyne.TextAlignTrailing, fyne.TextStyle{}),
			widget.NewLabelWithStyle(dcrutil.Amount(tx[i].Fee).String(), fyne.TextAlignCenter, fyne.TextStyle{}),
			widget.NewLabelWithStyle(tx[i].Status, fyne.TextAlignCenter, fyne.TextStyle{}),
			widget.NewButton(trimmedHash, func() { fmt.Println("Hello") }),
		))
	}

	fmt.Println(len(hBox))
	txTable.NewTable(heading, hBox...)
	txTable.Refresh()
}
