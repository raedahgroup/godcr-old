package pages

import (
	"time"

	"fyne.io/fyne"
	"fyne.io/fyne/widget"
	godcrApp "github.com/raedahgroup/godcr/app"
	"github.com/raedahgroup/godcr/app/walletcore"
	"github.com/raedahgroup/godcr/fyne/widgets"
)

type overviewPageData struct {
	balance  *widget.Label
	showText *widget.Label
}

var overview overviewPageData

func init() {
	overview.balance = widget.NewLabel("")
	overview.showText = widget.NewLabel("")
}

func overviewUpdates(wallet godcrApp.WalletMiddleware) {
	tx, _ := wallet.TransactionHistory(0, 5, nil)
	if len(tx) > 0 {
		if overview.showText != nil {
			widget.Refresh(overview.showText)
			overview.showText.Hide()
			widget.Refresh(overview.showText)
		}
	}
	widget.Refresh(overview.balance)
	overview.balance.SetText(fetchBalance(wallet))
	widget.Refresh(overview.balance)
}

func overviewPage(wallet godcrApp.WalletMiddleware) fyne.CanvasObject {
	label := widget.NewLabelWithStyle("Overview", fyne.TextAlignLeading, fyne.TextStyle{Italic: true, Bold: true})
	balanceLabel := widget.NewLabel("Current Total Balance")
	activityLabel := widget.NewLabel("Recent Activity")
	overview.balance = widget.NewLabel(fetchBalance(wallet))

	balanceLabel.TextStyle = fyne.TextStyle{Bold: true}
	activityLabel.TextStyle = fyne.TextStyle{Bold: true}
	var output fyne.CanvasObject
	overview.showText = widget.NewLabel("No activities yet")
	overview.showText.Alignment = fyne.TextAlignCenter

	tx, _ := wallet.TransactionHistory(0, 5, nil)
	if len(tx) == 0 {
		overview.showText.Hide()
		widget.Refresh(overview.showText)
	}
	output = widget.NewVBox(
		label,
		widgets.NewVSpacer(10),
		balanceLabel,
		overview.balance,
		widgets.NewVSpacer(10),
		activityLabel,
		widgets.NewVSpacer(10),
		overview.showText)

	//this would update all labels for all pages every seconds
	go func() {
		for {
			overviewUpdates(wallet)
			time.Sleep(time.Second * 1)
		}
	}()

	return widget.NewHBox(widgets.NewHSpacer(10), output)
}

func fetchBalance(wallet godcrApp.WalletMiddleware) string {
	accounts, err := wallet.AccountsOverview(walletcore.DefaultRequiredConfirmations)
	if err != nil {
		return err.Error()
	}

	return walletcore.WalletBalance(accounts)
}
