package pages

import (
	"fyne.io/fyne"
	"fyne.io/fyne/canvas"
	"fyne.io/fyne/layout"
	"fyne.io/fyne/widget"
	"github.com/raedahgroup/dcrlibwallet"
	"github.com/raedahgroup/godcr/fyne/layouts"
	"github.com/raedahgroup/godcr/fyne/pages/handler"
	"github.com/raedahgroup/godcr/fyne/pages/handler/values"
	"github.com/raedahgroup/godcr/fyne/widgets"
	"image/color"
	"sort"
)

const historyPageIndex = 1
const PageTitle = "Overview"

type overview struct {
	app          *AppInterface
	multiWallet  *dcrlibwallet.MultiWallet
	walletIds    []int
	transactions []dcrlibwallet.Transaction
}

var overviewHandler = &handler.OverviewHandler{}

func overviewPageContent(app *AppInterface) fyne.CanvasObject {
	ov := &overview{}
	ov.app = app
	ov.multiWallet = app.MultiWallet
	ov.walletIds = ov.multiWallet.OpenedWalletIDsRaw()
	if len(ov.walletIds) == 0 {
		return widget.NewHBox(widgets.NewHSpacer(10), widget.NewLabelWithStyle("Could not retrieve wallets", fyne.TextAlignCenter, fyne.TextStyle{Bold: true}))
	}
	sort.Ints(ov.walletIds)

	defer func() {
		go overviewHandler.UpdateBalance(app.MultiWallet)
		go overviewHandler.UpdateBlockStatusBox(app.MultiWallet)
		overviewHandler.UpdateWalletsSyncBox(app.MultiWallet)
	}()

	overviewHandler.Container = ov.container()
	return widget.NewHBox(widgets.NewHSpacer(values.Padding), overviewHandler.Container, widgets.NewHSpacer(values.Padding))
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
	overviewHandler.TableHeader = transactionRowHeader()
	table.NewTable(overviewHandler.TableHeader)
	overviewHandler.Table = table
	overviewHandler.UpdateTransactions(ov.multiWallet, handler.TransactionUpdate{})
	transactionsContainer := fyne.NewContainerWithLayout(layout.NewFixedGridLayout(fyne.NewSize(520, 200)), table.Container)
	overviewHandler.TransactionsContainer = transactionsContainer

	return widget.NewVBox(
		transactionsContainer,
		fyne.NewContainerWithLayout(layout.NewHBoxLayout(),
			layout.NewSpacer(),
			widgets.NewClickableBox(
				widget.NewHBox(widget.NewLabelWithStyle("see all", fyne.TextAlignCenter, fyne.TextStyle{})),
				func() {
					ov.app.tabMenu.SelectTabIndex(historyPageIndex)
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
	syncSteps := widget.NewLabelWithStyle("Step 1/3", fyne.TextAlignTrailing, fyne.TextStyle{})
	blockHeadersStatus := widget.NewLabelWithStyle("Fetching block headers  0%", fyne.TextAlignTrailing, fyne.TextStyle{})
	walletSyncInfo := fyne.NewContainerWithLayout(layout.NewHBoxLayout())
	walletSyncScrollContainer := widget.NewScrollContainer(walletSyncInfo)
	syncButton := widget.NewButton("", func() { overviewHandler.SyncTrigger(ov.multiWallet)})
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
	overviewHandler.SyncButton = syncButton
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
			h.SyncButton,
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
	top := fyne.NewContainerWithLayout(layout.NewFixedGridLayout(fyne.NewSize(515, 24)),
		widget.NewHBox(
			widgets.NewHSpacer(values.NilSpacer),
			h.SyncStatusWidget,
			layout.NewSpacer(),
			h.SyncButton,
		))
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
	dcrBalance := widgets.NewTextWithStyle("0.00", color.Black, fyne.TextStyle{}, fyne.TextAlignLeading, values.TextSize25)
	dcrDecimals := widgets.NewTextWithStyle("00000 DCR", color.Black, fyne.TextStyle{}, fyne.TextAlignLeading, values.TextSize15)
	overviewHandler.Balance = make([]*canvas.Text, 2)
	overviewHandler.Balance[0] = dcrBalance
	overviewHandler.Balance[1] = dcrDecimals
	return fyne.NewContainerWithLayout(layouts.NewHBox(0, true), dcrBalance, dcrDecimals)
}

func transactionRowHeader() *widget.Box {
	hash := widget.NewLabelWithStyle("#", fyne.TextAlignCenter, fyne.TextStyle{Bold: true})
	amount := widget.NewLabelWithStyle("Amount", fyne.TextAlignCenter, fyne.TextStyle{Bold: true})
	fee := widget.NewLabelWithStyle("Fee", fyne.TextAlignCenter, fyne.TextStyle{Bold: true})
	direction := widget.NewLabelWithStyle("Direction", fyne.TextAlignCenter, fyne.TextStyle{Bold: true})
	status := widget.NewLabelWithStyle("Status", fyne.TextAlignCenter, fyne.TextStyle{Bold: true})
	date := widget.NewLabelWithStyle("Date", fyne.TextAlignCenter, fyne.TextStyle{Bold: true})
	return widget.NewHBox(hash, amount, fee, direction, status, date)
}
