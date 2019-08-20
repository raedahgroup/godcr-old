package pages

import (
	"fmt"
	"strconv"
	"strings"

	"fyne.io/fyne"
	"fyne.io/fyne/layout"
	"fyne.io/fyne/widget"
	"github.com/decred/dcrd/dcrutil"
	godcrApp "github.com/raedahgroup/godcr/app"
	"github.com/raedahgroup/godcr/app/walletcore"
	"github.com/raedahgroup/godcr/fyne/widgets"
)

//this contains widgets that needs to be updated realtime
type overviewPageData struct {
	balance         *widget.Label
	noActivityLabel *widget.Label
	txTable         widgets.TableStruct
}

var overview overviewPageData

func overviewUpdates(wallet godcrApp.WalletMiddleware) {
	overview.balance.SetText(fetchBalance(wallet))
	var txTable widgets.TableStruct
	fetchTxTable(false, &txTable, 0, 5, wallet)
	overview.txTable.Result.Children = txTable.Result.Children
	widget.Refresh(overview.txTable.Result)
}

//this updates peerconn and blkheight
func statusUpdates(wallet godcrApp.WalletMiddleware) {
	info, _ := wallet.WalletConnectionInfo()

	if info.PeersConnected <= 1 {
		menu.peerConn.SetText(strconv.Itoa(int(info.PeersConnected)) + " Peer Connected")
	} else {
		menu.peerConn.SetText(strconv.Itoa(int(info.PeersConnected)) + " Peers Connected")
	}

	menu.blkHeight.SetText(strconv.Itoa(int(info.LatestBlock)) + " Blocks Connected")
}

func overviewPage(wallet godcrApp.WalletMiddleware) fyne.CanvasObject {
	label := widget.NewLabelWithStyle("Overview", fyne.TextAlignLeading, fyne.TextStyle{Italic: true, Bold: true})
	balanceLabel := widget.NewLabelWithStyle("Current Total Balance", fyne.TextAlignLeading, fyne.TextStyle{Bold: true})
	activityLabel := widget.NewLabelWithStyle("Recent Activity", fyne.TextAlignLeading, fyne.TextStyle{Bold: true})
	overview.balance = widget.NewLabel(fetchBalance(wallet))

	overview.noActivityLabel = widget.NewLabelWithStyle("No activities yet", fyne.TextAlignCenter, fyne.TextStyle{})

	fetchTxTable(false, &overview.txTable, 0, 5, wallet)
	output := widget.NewVBox(
		label,
		widgets.NewVSpacer(10),
		balanceLabel,
		overview.balance,
		widgets.NewVSpacer(10),
		activityLabel,
		widgets.NewVSpacer(10),
		overview.noActivityLabel,
		fyne.NewContainerWithLayout(layout.NewFixedGridLayout(fyne.NewSize(overview.txTable.Result.MinSize().Width, overview.txTable.Result.MinSize().Height)), overview.txTable.Container))

	return widget.NewHBox(widgets.NewHSpacer(10), output)
}

func fetchBalance(wallet godcrApp.WalletMiddleware) string {
	accounts, err := wallet.AccountsOverview(walletcore.DefaultRequiredConfirmations)
	if err != nil {
		return err.Error()
	}

	return walletcore.WalletBalance(accounts)
}

func fetchTxTable(isHistoryPage bool, txTable *widgets.TableStruct, offset, counter int32, wallet godcrApp.WalletMiddleware) {
	var txs []*walletcore.Transaction
	if !isHistoryPage {
		txs, _ = wallet.TransactionHistory(offset, counter, nil)
	} else {
		splittedWord := strings.Split(history.txFilters.Selected, " ")
		txs, _ = wallet.TransactionHistory(offset, counter, walletcore.BuildTransactionFilter(splittedWord[0]))
	}
	history.currentTxCount = int32(len(txs))
	fmt.Println("len of txs", len(txs))
	if len(txs) > 0 {
		overview.noActivityLabel.Hide()
	}

	heading := widget.NewHBox(
		widget.NewLabelWithStyle("Date (UTC)", fyne.TextAlignCenter, fyne.TextStyle{Bold: true}),
		widget.NewLabelWithStyle("Type", fyne.TextAlignCenter, fyne.TextStyle{Bold: true}),
		widget.NewLabelWithStyle("Direction", fyne.TextAlignCenter, fyne.TextStyle{Bold: true}),
		widget.NewLabelWithStyle("Amount", fyne.TextAlignCenter, fyne.TextStyle{Bold: true}),
		widget.NewLabelWithStyle("Fee", fyne.TextAlignCenter, fyne.TextStyle{Bold: true}),
		widget.NewLabelWithStyle("Status", fyne.TextAlignCenter, fyne.TextStyle{Bold: true}),
		widget.NewLabelWithStyle("Hash", fyne.TextAlignCenter, fyne.TextStyle{Bold: true}))

	var hBox []*widget.Box
	for _, tx := range txs {
		trimmedHash := tx.Hash[:len(tx.Hash)/2] + "..."
		var hash fyne.CanvasObject
		if isHistoryPage {
			hash = widget.NewButton(trimmedHash, func() {
				fmt.Println("Hello")
			})
		} else {
			hash = widget.NewLabelWithStyle(trimmedHash, fyne.TextAlignCenter, fyne.TextStyle{})
		}

		hBox = append(hBox, widget.NewHBox(
			widget.NewLabelWithStyle(tx.LongTime, fyne.TextAlignCenter, fyne.TextStyle{}),
			widget.NewLabelWithStyle(tx.Type, fyne.TextAlignCenter, fyne.TextStyle{}),
			widget.NewLabelWithStyle(tx.Direction.String(), fyne.TextAlignLeading, fyne.TextStyle{}),
			widget.NewLabelWithStyle(dcrutil.Amount(tx.Amount).String(), fyne.TextAlignTrailing, fyne.TextStyle{}),
			widget.NewLabelWithStyle(dcrutil.Amount(tx.Fee).String(), fyne.TextAlignCenter, fyne.TextStyle{}),
			widget.NewLabelWithStyle(tx.Status, fyne.TextAlignCenter, fyne.TextStyle{}),
			hash,
		))
	}
	txTable.NewTable(heading, hBox...)
}
