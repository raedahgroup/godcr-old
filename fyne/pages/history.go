package pages

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/raedahgroup/godcr/app/walletcore"

	"fyne.io/fyne"
	"fyne.io/fyne/layout"
	"fyne.io/fyne/widget"
	"github.com/decred/dcrd/dcrutil"
	godcrApp "github.com/raedahgroup/godcr/app"
	"github.com/raedahgroup/godcr/fyne/widgets"
)

type historyPageData struct {
	totalTx        int
	currentTxCount int
	txFilters      *widget.Select
	options        map[string]int
	//txType  string
	txTable widgets.TableStruct
}

var offsetYMin, offsetYMax = 180, 10
var sizeMin, sizeMax = 433, 610
var selected bool
var history historyPageData

func historyUpdates(wallet godcrApp.WalletMiddleware) {
	filters := walletcore.TransactionFilters
	txCountByFilter := make(map[string]int)

	for _, filter := range filters {
		txCount, txCountErr := wallet.TransactionCount(walletcore.BuildTransactionFilter(filter))
		if txCountErr != nil {
			//treat
			return
		}
		txCountByFilter[filter] = txCount
	}
	var options []string
	count, _ := wallet.TransactionCount(nil)
	options = append(options, "All ("+strconv.Itoa(count)+")")
	options = append(options, "Sent ("+strconv.Itoa(txCountByFilter["Sent"])+")")
	options = append(options, "Received ("+strconv.Itoa(txCountByFilter["Received"])+")")
	options = append(options, "Yourself ("+strconv.Itoa(txCountByFilter["Yourself"])+")")
	options = append(options, "Staking ("+strconv.Itoa(txCountByFilter["Staking"])+")")

	history.txFilters.Options = options
	widget.Refresh(history.txFilters)
	if !selected {
		selected = true
		history.txFilters.SetSelected(options[0])
	}

	if history.txTable.Container.Offset.Y >= offsetYMin && history.txTable.Container.Size().Height >= sizeMin || history.txTable.Container.Offset.Y >= offsetYMax && history.txTable.Container.Size().Height >= sizeMax {

	} else if history.txTable.Container.Offset.Y == 0 {
		// if the scroll bar is at the begining, then fetch 1st 50 tx
		if count > history.currentTxCount {
			splittedWord := strings.Split(history.txFilters.Selected, " ")
			history.txFilters.SetSelected(history.txFilters.Options[history.options[splittedWord[0]]])
			history.currentTxCount = count
		}
	}
}

func historyPage(wallet godcrApp.WalletMiddleware) fyne.CanvasObject {
	history.options = make(map[string]int)
	history.options["All"] = 0
	history.options["Sent"] = 1
	history.options["Received"] = 2
	history.options["Yourself"] = 3
	history.options["Staking"] = 4

	history.txFilters = widget.NewSelect(nil, func(selected string) {
		// if a new type is selected, load the first 50tx so as to allow the scroller move to the starting point
		var txTable widgets.TableStruct
		fetchTxTable(true, &txTable, 0, 50, wallet)
		history.txTable.Result.Children = txTable.Result.Children
		widget.Refresh(history.txTable.Result)
	})

	heading := widget.NewHBox(
		widget.NewLabelWithStyle("Date (UTC)", fyne.TextAlignCenter, fyne.TextStyle{Bold: true}),
		widget.NewLabelWithStyle("Type", fyne.TextAlignCenter, fyne.TextStyle{Bold: true}),
		widget.NewLabelWithStyle("Direction", fyne.TextAlignCenter, fyne.TextStyle{Bold: true}),
		widget.NewLabelWithStyle("Amount", fyne.TextAlignCenter, fyne.TextStyle{Bold: true}),
		widget.NewLabelWithStyle("Fee", fyne.TextAlignCenter, fyne.TextStyle{Bold: true}),
		widget.NewLabelWithStyle("Status", fyne.TextAlignCenter, fyne.TextStyle{Bold: true}),
		widget.NewLabelWithStyle("Hash", fyne.TextAlignCenter, fyne.TextStyle{Bold: true}))
	history.txTable.NewTable(heading)
	history.currentTxCount, _ = wallet.TransactionCount(nil)
	historyUpdates(wallet)

	output := widget.NewVBox(widget.NewHBox(layout.NewSpacer(), history.txFilters),
		fyne.NewContainerWithLayout(layout.NewFixedGridLayout(fyne.NewSize(history.txTable.Result.MinSize().Width, (history.txTable.Result.MinSize().Height/3)+10)), history.txTable.Container))

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
