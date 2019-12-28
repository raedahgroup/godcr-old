package overview

import (
	"fyne.io/fyne"
	"fyne.io/fyne/canvas"
	"fyne.io/fyne/layout"
	"fyne.io/fyne/widget"
	"github.com/raedahgroup/godcr/fyne/values"
	"image/color"
	"sort"
	"time"

	"github.com/raedahgroup/dcrlibwallet"
	"github.com/raedahgroup/godcr/fyne/widgets"
)

const historyPageIndex = 1
const PageTitle = "overview"

type overview struct {
	multiWallet  *dcrlibwallet.MultiWallet
	walletIds    []int
	transactions []dcrlibwallet.Transaction
	tabMenu 	 *widget.TabContainer
}

var overviewHandler = &Handler{}

func PageContent(mw *dcrlibwallet.MultiWallet, tabMenu *widget.TabContainer) (content fyne.CanvasObject, handler *Handler) {
	ov := &overview{}
	//app.Window.Resize(fyne.NewSize(650, 680))
	ov.multiWallet = mw
	ov.tabMenu = tabMenu
	ov.walletIds = ov.multiWallet.OpenedWalletIDsRaw()
	if len(ov.walletIds) == 0 {
		content = widget.NewHBox(
			widgets.NewHSpacer(10),
			widget.NewLabelWithStyle("Could not retrieve wallets", fyne.TextAlignCenter, fyne.TextStyle{Bold: true}))
		return
	}
	sort.Ints(ov.walletIds)

	defer func() {
		go overviewHandler.UpdateBalance(mw)
		go overviewHandler.UpdateBlockStatusBox(mw)
	}()
	go func() {
		time.Sleep(200 * time.Millisecond)
		overviewHandler.UpdateWalletsSyncBox(mw)
	}()

	overviewHandler.Container = ov.container()
	handler = overviewHandler
	content = widget.NewHBox(widgets.NewHSpacer(values.Padding), overviewHandler.Container, widgets.NewHSpacer(values.Padding))
	return
}

func (ov *overview) container() fyne.CanvasObject {
	overviewHandler.PageBoxes = ov.pageBoxes()
	overviewContainer := widget.NewVBox(
		widgets.NewVSpacer(values.Padding),
		title(),
		balance(),
		widgets.NewVSpacer(25),
		overviewHandler.PageBoxes,
		widgets.NewVSpacer(values.Padding),
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
	table := &widgets.Table{}
	table.NewTable(transactionRowHeader())
	overviewHandler.Table = table
	overviewHandler.UpdateTransactions(ov.multiWallet, TransactionUpdate{})
	return widget.NewVBox(
		table.Result,
		fyne.NewContainerWithLayout(layout.NewHBoxLayout(),
			layout.NewSpacer(),
			widgets.NewClickableBox(
				widget.NewHBox(widget.NewLabelWithStyle("see all", fyne.TextAlignCenter, fyne.TextStyle{})),
				func() {
					ov.tabMenu.SelectTabIndex(historyPageIndex)
				},
			),
			widgets.NewVSpacer(40),
			layout.NewSpacer(),
		),
	)
}

func (ov *overview) initializeBlockStatusWidgets() {
	syncStatusText := widgets.NewSmallText("Syncing...", values.DefaultTextColor)
	timeLeft := widget.NewLabelWithStyle("", fyne.TextAlignLeading, fyne.TextStyle{})
	connectedPeers := widget.NewLabelWithStyle("", fyne.TextAlignTrailing, fyne.TextStyle{})
	progressBar := widget.NewProgressBar()
	syncSteps := widget.NewLabelWithStyle("Step 0/3", fyne.TextAlignTrailing, fyne.TextStyle{})
	blockHeadersStatus := widget.NewLabelWithStyle("Fetching block headers  0%", fyne.TextAlignTrailing, fyne.TextStyle{})
	walletSyncInfo := fyne.NewContainerWithLayout(layout.NewHBoxLayout())
	walletSyncScrollContainer := widget.NewScrollContainer(walletSyncInfo)
	cancelButton := widget.NewButton("Cancel", func() { overviewHandler.CancelSync(ov.multiWallet)})
	blockHeightTime := widget.NewLabelWithStyle("Latest block 0 . 0s ago", fyne.TextAlignLeading, fyne.TextStyle{})
	overviewHandler.BlockHeightTime = widget.NewLabelWithStyle("", fyne.TextAlignLeading, fyne.TextStyle{})

	overviewHandler.SyncStatusWidget = syncStatusText
	overviewHandler.TimeLeftWidget = timeLeft
	overviewHandler.ProgressBar = progressBar
	overviewHandler.ConnectedPeersWidget = connectedPeers
	overviewHandler.SyncStepWidget = syncSteps
	overviewHandler.BlockHeadersWidget = blockHeadersStatus
	overviewHandler.WalletSyncInfo = walletSyncInfo
	overviewHandler.Scroll = walletSyncScrollContainer
	overviewHandler.CancelButton = cancelButton
	overviewHandler.BlockHeightTime = blockHeightTime
}

func (ov *overview) blockStatusBox() fyne.CanvasObject {
	ov.initializeBlockStatusWidgets()
	overviewHandler.BlockStatusSyncing = ov.blockStatusBoxSyncing()
	overviewHandler.BlockStatusSynced = ov.blockStatusBoxSynced()
	overviewHandler.BlockStatus = fyne.NewContainerWithLayout(layout.NewVBoxLayout())
	return overviewHandler.BlockStatus
}

func (ov *overview) blockStatusBoxSyncing() fyne.CanvasObject {
	h := overviewHandler
	top := fyne.NewContainerWithLayout(layout.NewFixedGridLayout(fyne.NewSize(515, 24)),
		widget.NewHBox(
			widgets.NewHSpacer(values.NilSpacer),
			h.SyncStatusWidget,
			layout.NewSpacer(),
			h.CancelButton,
		))
	progressBarContainer := fyne.NewContainerWithLayout(layout.NewFixedGridLayout(fyne.NewSize(515, 20)),
		h.ProgressBar)
	syncDuration := fyne.NewContainerWithLayout(layout.NewBorderLayout(nil, nil, h.TimeLeftWidget, h.ConnectedPeersWidget),
		h.TimeLeftWidget, h.ConnectedPeersWidget)
	syncStatus := fyne.NewContainerWithLayout(layout.NewBorderLayout(nil, nil, h.SyncStepWidget, h.BlockHeadersWidget),
		h.SyncStepWidget, h.BlockHeadersWidget)

	return fyne.NewContainerWithLayout(layout.NewVBoxLayout(),
		widgets.NewVSpacer(5),
		top,
		progressBarContainer,
		syncDuration,
		syncStatus,
		widgets.NewVSpacer(15),
		fyne.NewContainerWithLayout(layout.NewFixedGridLayout(fyne.NewSize(510, 80)), h.Scroll),
	)
}

func (ov *overview) blockStatusBoxSynced() fyne.CanvasObject {
	h := overviewHandler
	top := fyne.NewContainerWithLayout(layout.NewHBoxLayout(),
		widgets.NewHSpacer(values.NilSpacer),
		h.SyncStatusWidget,
		layout.NewSpacer())
	syncedStatus := fyne.NewContainerWithLayout(layout.NewBorderLayout(nil, nil, h.BlockHeightTime, h.ConnectedPeersWidget),
		h.BlockHeightTime, h.ConnectedPeersWidget)

	return fyne.NewContainerWithLayout(layout.NewVBoxLayout(),
		widgets.NewVSpacer(5),
		top,
		syncedStatus,
		widgets.NewVSpacer(15),
	)
}

func title() fyne.CanvasObject {
	titleWidget := widget.NewLabelWithStyle(PageTitle, fyne.TextAlignLeading, fyne.TextStyle{Bold: true})
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
