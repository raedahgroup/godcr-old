package pages

import (
	"fmt"
	"strconv"

	"github.com/raedahgroup/godcr/app/walletcore"

	"fyne.io/fyne"
	"fyne.io/fyne/layout"
	"fyne.io/fyne/widget"
	"github.com/decred/dcrd/dcrutil"
	godcrApp "github.com/raedahgroup/godcr/app"
	"github.com/raedahgroup/godcr/fyne/widgets"
)

type historyPageData struct {
	loadedCount int
	txFilters   *widget.Select
	//txType  string
	txTable widgets.TableStruct
}

var history historyPageData

func historyUpdates(wallet godcrApp.WalletMiddleware) {
	filters := walletcore.TransactionFilters
	txCountByFilter := make(map[string]int, 0)

	for _, filter := range filters {
		txCount, txCountErr := wallet.TransactionCount(walletcore.BuildTransactionFilter(filter))
		if txCountErr != nil {
			//treat
			return
		}
		if txCount == 0 {
			continue
		}
		txCountByFilter[filter] = txCount
	}
	var options []string
	for value, name := range txCountByFilter {
		options = append(options, value+"("+strconv.Itoa(name)+")")
	}

	history.txFilters.Options = options
}

func historyPage(wallet godcrApp.WalletMiddleware) fyne.CanvasObject {
	history.txFilters = widget.NewSelect(nil, func(selected string) {
		fmt.Println(selected)
	})
	historyUpdates(wallet)

	heading := widget.NewHBox(
		widget.NewLabelWithStyle("Date (UTC)", fyne.TextAlignCenter, fyne.TextStyle{Bold: true}),
		widget.NewLabelWithStyle("Type", fyne.TextAlignCenter, fyne.TextStyle{Bold: true}),
		widget.NewLabelWithStyle("Direction", fyne.TextAlignCenter, fyne.TextStyle{Bold: true}),
		widget.NewLabelWithStyle("Amount", fyne.TextAlignCenter, fyne.TextStyle{Bold: true}),
		widget.NewLabelWithStyle("Fee", fyne.TextAlignCenter, fyne.TextStyle{Bold: true}),
		widget.NewLabelWithStyle("Status", fyne.TextAlignCenter, fyne.TextStyle{Bold: true}),
		widget.NewLabelWithStyle("Hash", fyne.TextAlignCenter, fyne.TextStyle{Bold: true}))
	history.txTable.NewTable(heading)

	output := widget.NewVBox(widget.NewHBox(layout.NewSpacer(), history.txFilters),
		fyne.NewContainerWithLayout(layout.NewFixedGridLayout(fyne.NewSize(history.txTable.Result.MinSize().Width, history.txTable.Result.MinSize().Height)), history.txTable.Container))

	return widget.NewHBox(widgets.NewHSpacer(10), output)
}

func fetchHistoryTx(txTable *widgets.TableStruct, offset, count int32, wallet godcrApp.WalletMiddleware) {
	txs, _ := wallet.TransactionHistory(offset, count, walletcore.BuildTransactionFilter(history.txFilters.Selected))
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

	history.txTable.Append(hBox...)
}
