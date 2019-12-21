package pages

import (
	"fyne.io/fyne"
	"fyne.io/fyne/canvas"
	"fyne.io/fyne/layout"
	"fyne.io/fyne/widget"
	"image/color"
	"sort"
	"time"

	"github.com/raedahgroup/dcrlibwallet"
	"github.com/raedahgroup/godcr/fyne/handlers"
	"github.com/raedahgroup/godcr/fyne/widgets"
)

const historyPageIndex = 1
const PageTitle = "overview"

type overview struct {
	app            *AppInterface
	multiWallet    *dcrlibwallet.MultiWallet
	walletIds      []int
	transactions   []dcrlibwallet.Transaction
}

var overviewHandler = &handlers.OverviewHandler{}

// todo: display overview page (include sync progress UI elements)
// todo: register sync progress listener on overview page to update sync progress views
func overviewPageContent(app *AppInterface) fyne.CanvasObject {
	ov := &overview{}
	ov.app = app
	app.Window.Resize(fyne.NewSize(650, 650))
	ov.multiWallet = app.MultiWallet
	ov.walletIds = ov.multiWallet.OpenedWalletIDsRaw()
	if len(ov.walletIds) == 0 {
		return widget.NewHBox(widgets.NewHSpacer(10), widget.NewLabelWithStyle("Could not retrieve wallets", fyne.TextAlignCenter, fyne.TextStyle{Bold: true}))
	}
	sort.Ints(ov.walletIds)

	defer func (){
		overviewHandler.UpdateBalance(app.MultiWallet)
		handlers.OverviewHandlerLock.Lock()
		overviewHandler.UpdateBlockStatusBox(app.MultiWallet)
		handlers.OverviewHandlerLock.Unlock()
	}()
	go func() {
		time.Sleep(200 * time.Millisecond)
		overviewHandler.UpdateWalletsSyncBox(app.MultiWallet)
	}()

	return widget.NewHBox(widgets.NewHSpacer(18), ov.container())
}

func (ov *overview) container() fyne.CanvasObject {
	overviewContainer := widget.NewVBox(
		title(),
		balance(),
		widgets.NewVSpacer(50),
		ov.pageBoxes(),
	)
	return overviewContainer
}

func (ov *overview) pageBoxes() (object fyne.CanvasObject) {
	return fyne.NewContainerWithLayout(layout.NewVBoxLayout(),
		ov.blockStatusBox(),
		widgets.NewVSpacer(15),
		ov.recentTransactionBox(),
	)
}

func (ov *overview) recentTransactionBox() fyne.CanvasObject {
	//var err error
	//overviewHandler.Transactions, err = recentTransactions(ov)
	//if err != nil {
	//	return widget.NewHBox(widgets.NewHSpacer(10), widget.NewLabelWithStyle(err.Error(), fyne.TextAlignCenter, fyne.TextStyle{Bold: true}))
	//}

	table := &widgets.Table{}
	table.NewTable(transactionRowHeader())
	overviewHandler.Table = table
	overviewHandler.UpdateTransactions(ov.multiWallet, handlers.TransactionUpdate{})
	return widget.NewVBox(
		table.Result,
		fyne.NewContainerWithLayout(layout.NewHBoxLayout(),
			layout.NewSpacer(),
			widgets.NewClickableBox(
				widget.NewHBox(widget.NewLabelWithStyle("see all", fyne.TextAlignCenter, fyne.TextStyle{Italic: true})),
				func() {
					ov.app.tabMenu.SelectTabIndex(historyPageIndex)
				},
			),
			layout.NewSpacer(),
		),
	)
}

func (ov *overview) blockStatusBox() fyne.CanvasObject {
	syncStatusText := widgets.NewSmallText("", color.Black)
	timeLeft := widget.NewLabelWithStyle("", fyne.TextAlignLeading, fyne.TextStyle{Italic: true})
	connectedPeers := widget.NewLabelWithStyle("", fyne.TextAlignTrailing, fyne.TextStyle{Italic: true})
	progressBar := widget.NewProgressBar()
	syncSteps := widget.NewLabelWithStyle("Step 1/3", fyne.TextAlignTrailing, fyne.TextStyle{Italic: true})
	blockHeadersStatus := widget.NewLabelWithStyle("Fetching block headers  89%", fyne.TextAlignTrailing, fyne.TextStyle{Italic: true})
	walletSyncInfo := fyne.NewContainerWithLayout(layout.NewGridLayout(2))
	walletSyncScrollContainer := widget.NewScrollContainer(walletSyncInfo)
	walletSyncInfoToggleText := widget.NewLabelWithStyle("hide details", fyne.TextAlignCenter, fyne.TextStyle{Italic: true})
	walletSyncInfoToggle := widgets.NewClickableBox(widget.NewHBox(walletSyncInfoToggleText),
		func() {
			overviewHandler.HideWalletSyncBox()
		},
	)

	overviewHandler.SyncStatusWidget = syncStatusText
	overviewHandler.TimeLeftWidget = timeLeft
	overviewHandler.ProgressBar = progressBar
	overviewHandler.ConnectedPeersWidget = connectedPeers
	overviewHandler.SyncStepWidget = syncSteps
	overviewHandler.BlockHeadersWidget = blockHeadersStatus
	overviewHandler.WalletSyncInfo = walletSyncInfo
	overviewHandler.WalletSyncInfoToggle = walletSyncInfoToggle
	overviewHandler.WalletSyncInfoToggleText = walletSyncInfoToggleText
	overviewHandler.Scroll = walletSyncScrollContainer
	top := fyne.NewContainerWithLayout(layout.NewFixedGridLayout(fyne.NewSize(515, 24)),
		widget.NewHBox(
			syncStatusText,
			layout.NewSpacer(),
			widget.NewButton("Cancel", func() {
				overviewHandler.CancelSync(ov.multiWallet)
			}),
		))
	progressBarContainer := fyne.NewContainerWithLayout(layout.NewFixedGridLayout(fyne.NewSize(515, 20)),
		progressBar,
	)
	syncDuration := fyne.NewContainerWithLayout(layout.NewBorderLayout(nil, nil, timeLeft, connectedPeers),
		timeLeft, connectedPeers)
	syncStatus := fyne.NewContainerWithLayout(layout.NewBorderLayout(nil, nil, syncSteps, blockHeadersStatus),
		syncSteps, blockHeadersStatus)

	blockStatus := fyne.NewContainerWithLayout(layout.NewVBoxLayout(),
		widgets.NewVSpacer(5),
		top,
		progressBarContainer,
		syncDuration,
		syncStatus,
		widgets.NewVSpacer(15),
		walletSyncScrollContainer,
		fyne.NewContainerWithLayout(layout.NewHBoxLayout(),
			layout.NewSpacer(),
			widgets.NewClickableBox(
				widget.NewHBox(walletSyncInfoToggleText),
				func() {
					overviewHandler.HideWalletSyncBox()
				},
			),
			layout.NewSpacer(),
		),
	)
	overviewHandler.BlockStatus = blockStatus
	return blockStatus
}

func title() fyne.CanvasObject {
	titleWidget := widget.NewLabelWithStyle(PageTitle, fyne.TextAlignLeading, fyne.TextStyle{Bold: true, Italic: true})
	return widget.NewHBox(titleWidget)
}

func balance() fyne.CanvasObject {
	dcrBalance := widgets.NewLargeText("0.00", color.Black)
	dcrDecimals := widgets.NewSmallText("00000 DCR", color.Black)
	overviewHandler.Balance = make([]*canvas.Text, 2)
	overviewHandler.Balance[0] = dcrBalance
	overviewHandler.Balance[1] = dcrDecimals
	decimalsBox := fyne.NewContainerWithLayout(layout.NewVBoxLayout(), widgets.NewVSpacer(6), dcrDecimals)
	return widget.NewHBox(widgets.NewVSpacer(10), dcrBalance, decimalsBox)
}

func transactionRowHeader() *widget.Box {
	hash := widget.NewLabelWithStyle("#", fyne.TextAlignCenter, fyne.TextStyle{Bold: true})
	amount := widget.NewLabelWithStyle("amount", fyne.TextAlignCenter, fyne.TextStyle{Bold: true})
	fee := widget.NewLabelWithStyle("fee", fyne.TextAlignCenter, fyne.TextStyle{Bold: true})
	direction := widget.NewLabelWithStyle("direction", fyne.TextAlignCenter, fyne.TextStyle{Bold: true})
	status := widget.NewLabelWithStyle("status", fyne.TextAlignCenter, fyne.TextStyle{Bold: true})
	date := widget.NewLabelWithStyle("date", fyne.TextAlignCenter, fyne.TextStyle{Bold: true})
	return widget.NewHBox(hash, amount, fee, direction, status, date)
}
