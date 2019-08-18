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
	txCount     int
	loadedCount int
	txFilters   *widget.Select
	//txType  string
	txTable widgets.TableStruct
}

var history historyPageData

func historyUpdates(wallet godcrApp.WalletMiddleware) {
	filters := walletcore.TransactionFilters
	transactionCountByFilter := make(map[string]int, 0)

	for _, filter := range filters {
		txCount, txCountErr := wallet.TransactionCount(walletcore.BuildTransactionFilter(filter))
		if txCountErr != nil {
			//treat
			return
		}
		if txCount == 0 {
			continue
		}
		transactionCountByFilter[filter] = txCount
	}
	var options []string
	for value, name := range transactionCountByFilter {
		options = append(options, value+"("+strconv.Itoa(name)+")")
	}

	history.txFilters.Options = options
}

func historyPage(wallet godcrApp.WalletMiddleware) fyne.CanvasObject {
	history.txFilters = widget.NewSelect([]string{}, func(selected string) {
		fmt.Println("Hello")
	})
	historyUpdates(wallet)

	return widget.NewHBox(layout.NewSpacer(), history.txFilters)
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
}
