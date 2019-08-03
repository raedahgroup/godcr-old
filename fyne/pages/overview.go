package pages

import (
	"strconv"

	"fyne.io/fyne"
	"fyne.io/fyne/widget"
	godcrApp "github.com/raedahgroup/godcr/app"
	"github.com/raedahgroup/godcr/app/walletcore"
	"github.com/raedahgroup/godcr/fyne/widgets"
)

//this contains widgets that needs to be updated realtime
type overviewPageData struct {
	balance         *widget.Label
	noActivityLabel *widget.Label
}

var overview overviewPageData

func overviewUpdates(wallet godcrApp.WalletMiddleware) {
	tx, _ := wallet.TransactionHistory(0, 5, nil)

	if len(tx) > 0 {
		if overview.noActivityLabel.Text != "" {
			overview.noActivityLabel.SetText("")
			overview.noActivityLabel.Hide()
		}
	}

	overview.balance.SetText(fetchBalance(wallet))
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

	tx, _ := wallet.TransactionHistory(0, 5, nil)
	if len(tx) == 0 {
		overview.noActivityLabel.Hide()
	}

	output := widget.NewVBox(
		label,
		widgets.NewVSpacer(10),
		balanceLabel,
		overview.balance,
		widgets.NewVSpacer(10),
		activityLabel,
		widgets.NewVSpacer(10),
		overview.noActivityLabel)

	return widget.NewHBox(widgets.NewHSpacer(10), output)
}

func fetchBalance(wallet godcrApp.WalletMiddleware) string {
	accounts, err := wallet.AccountsOverview(walletcore.DefaultRequiredConfirmations)
	if err != nil {
		return err.Error()
	}

	return walletcore.WalletBalance(accounts)
}
