package pages

import (
	"io/ioutil"
	"log"
	"strconv"

	"fyne.io/fyne/theme"

	"fyne.io/fyne/canvas"

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
	icon            *canvas.Image
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

func overviewPage(wallet godcrApp.WalletMiddleware, fyneApp fyne.App) fyne.CanvasObject {
	fyneTheme := fyneApp.Settings().Theme()
	if fyneTheme.BackgroundColor() == theme.LightTheme().BackgroundColor() {
		decredDark, err := ioutil.ReadFile("./fyne/pages/png/decredDark.png")
		if err != nil {
			log.Fatalln("exit png file missing", err)
		}
		overview.icon = canvas.NewImageFromResource(fyne.NewStaticResource("Decred", decredDark)) //NewIcon(fyne.NewStaticResource("deced", decredLogo))
	} else if fyneTheme.BackgroundColor() == theme.DarkTheme().BackgroundColor() {
		decredLight, err := ioutil.ReadFile("./fyne/pages/png/decredLight.png")
		if err != nil {
			log.Fatalln("exit png file missing", err)
		}
		overview.icon = canvas.NewImageFromResource(fyne.NewStaticResource("Decred", decredLight)) //NewIcon(fyne.NewStaticResource("deced", decredLogo))
	}

	iconFix := fyne.NewContainerWithLayout(layout.NewFixedGridLayout(fyne.NewSize(100, 100)), overview.icon)
	name := widget.NewLabelWithStyle(godcrApp.DisplayName, fyne.TextAlignLeading, fyne.TextStyle{Bold: true})
	iconLabel := fyne.NewContainer(iconFix, name)
	name.Move(fyne.NewPos(20, 65))

	balanceLabel := widget.NewLabelWithStyle("Current Total Balance", fyne.TextAlignLeading, fyne.TextStyle{Bold: true})
	activityLabel := widget.NewLabelWithStyle("Recent Activity", fyne.TextAlignLeading, fyne.TextStyle{Bold: true})
	overview.balance = widget.NewLabel(fetchBalance(wallet))

	overview.noActivityLabel = widget.NewLabelWithStyle("No activities yet", fyne.TextAlignCenter, fyne.TextStyle{})

	fetchOverviewTx(&overview.txTable, 0, 5, wallet)

	output := widget.NewVBox(
		iconLabel,
		balanceLabel,
		overview.balance,
		widgets.NewVSpacer(5),
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
