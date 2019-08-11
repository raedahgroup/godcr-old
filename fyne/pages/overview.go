package pages

import (
	"fmt"

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
	txCount, _ := wallet.TransactionCount(nil)

	if txCount > 0 {
		if overview.noActivityLabel.Text != "" {
			overview.noActivityLabel.SetText("")
			overview.noActivityLabel.Hide()
		}
	}

	overview.balance.SetText(fetchBalance(wallet))
	var txTable widgets.TableStruct

	fetchOverviewTx(&txTable, 0, 5, wallet)
	overview.txTable.Result.Children = txTable.Result.Children
	widget.Refresh(overview.txTable.Result)
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
		widgets.NewVSpacer(10),
		overview.noActivityLabel,
		fyne.NewContainerWithLayout(layout.NewFixedGridLayout(fyne.NewSize(overview.txTable.Container.MinSize().Width, overview.txTable.Container.MinSize().Height+200)), overview.txTable.Container))

	return widget.NewHBox(widgets.NewHSpacer(10), output)
}

func fetchBalance(wallet godcrApp.WalletMiddleware) string {
	accounts, err := wallet.AccountsOverview(walletcore.DefaultRequiredConfirmations)
	if err != nil {
		return err.Error()
	}

	return walletcore.WalletBalance(accounts)
}

func fetchOverviewTx(txTable *widgets.TableStruct, offset, counter int32, wallet godcrApp.WalletMiddleware) {
	tx, _ := wallet.TransactionHistory(offset, counter, nil)
	if len(tx) == 0 {
		overview.noActivityLabel.Hide()
		return
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
		trimmedHash := tx[i].Hash[:len(tx[i].Hash)/2] + "..."
		hBox = append(hBox, widget.NewHBox(
			widget.NewLabelWithStyle(tx[i].LongTime, fyne.TextAlignCenter, fyne.TextStyle{}),
			widget.NewLabelWithStyle(tx[i].Type, fyne.TextAlignCenter, fyne.TextStyle{}),
			widget.NewLabelWithStyle(tx[i].Direction.String(), fyne.TextAlignLeading, fyne.TextStyle{}),
			widget.NewLabelWithStyle(dcrutil.Amount(tx[i].Amount).String(), fyne.TextAlignTrailing, fyne.TextStyle{}),
			widget.NewLabelWithStyle(dcrutil.Amount(tx[i].Fee).String(), fyne.TextAlignCenter, fyne.TextStyle{}),
			widget.NewLabelWithStyle(tx[i].Status, fyne.TextAlignCenter, fyne.TextStyle{}),
			widget.NewLabelWithStyle(trimmedHash, fyne.TextAlignCenter, fyne.TextStyle{}),
		))
	}
	txTable.NewTable(heading, hBox...)
	txTable.Refresh()

	fmt.Println(txTable.Container.MinSize())
}
