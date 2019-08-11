package pages

import (
	"strconv"

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

	fetchOverviewTx(&txTable, 0, 5, wallet)
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

	fetchOverviewTx(&overview.txTable, 0, 5, wallet)

	output := widget.NewVBox(
		label,
		widgets.NewVSpacer(10),
		balanceLabel,
		overview.balance,
		widgets.NewVSpacer(10),
		activityLabel,
		overview.noActivityLabel,
		fyne.NewContainerWithLayout(layout.NewFixedGridLayout(fyne.NewSize(overview.txTable.Container.MinSize().Width, overview.txTable.Container.MinSize().Height+200)), overview.txTable.Container))

	return widget.NewHBox(widgets.NewHSpacer(20), output)
}

func fetchBalance(wallet godcrApp.WalletMiddleware) string {
	accounts, err := wallet.AccountsOverview(walletcore.DefaultRequiredConfirmations)
	if err != nil {
		return err.Error()
	}

	return walletcore.WalletBalance(accounts)
}

func fetchOverviewTx(txTable *widgets.TableStruct, offset, counter int32, wallet godcrApp.WalletMiddleware) {
	txs, _ := wallet.TransactionHistory(offset, counter, nil)
	if len(txs) > 0 && overview.noActivityLabel.Hidden == false {
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
	for i := 0; int32(i) < counter; i++ {
		trimmedHash := txs[i].Hash[:len(txs[i].Hash)/2] + "..."
		hBox = append(hBox, widget.NewHBox(
			widget.NewLabelWithStyle(txs[i].LongTime, fyne.TextAlignCenter, fyne.TextStyle{}),
			widget.NewLabelWithStyle(txs[i].Type, fyne.TextAlignCenter, fyne.TextStyle{}),
			widget.NewLabelWithStyle(txs[i].Direction.String(), fyne.TextAlignLeading, fyne.TextStyle{}),
			widget.NewLabelWithStyle(dcrutil.Amount(txs[i].Amount).String(), fyne.TextAlignTrailing, fyne.TextStyle{}),
			widget.NewLabelWithStyle(dcrutil.Amount(txs[i].Fee).String(), fyne.TextAlignCenter, fyne.TextStyle{}),
			widget.NewLabelWithStyle(txs[i].Status, fyne.TextAlignCenter, fyne.TextStyle{}),
			widget.NewLabelWithStyle(trimmedHash, fyne.TextAlignLeading, fyne.TextStyle{}),
		))
	}
	txTable.NewTable(heading, hBox...)
	txTable.Refresh()
}
