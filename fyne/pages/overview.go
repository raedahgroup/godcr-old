package pages

import (
	"fyne.io/fyne"
	"fyne.io/fyne/canvas"
	"fyne.io/fyne/layout"
	"fyne.io/fyne/widget"
	"github.com/raedahgroup/godcr/fyne/pages/handler"
	"image/color"
	"sort"
	"time"

	"github.com/raedahgroup/dcrlibwallet"
	"github.com/raedahgroup/godcr/fyne/widgets"
)

const historyPageIndex = 1
const PageTitle = "overview"

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
	//app.Window.Resize(fyne.NewSize(650, 680))
	ov.multiWallet = app.MultiWallet
	ov.walletIds = ov.multiWallet.OpenedWalletIDsRaw()
	if len(ov.walletIds) == 0 {
		return widget.NewHBox(widgets.NewHSpacer(10), widget.NewLabelWithStyle("Could not retrieve wallets", fyne.TextAlignCenter, fyne.TextStyle{Bold: true}))
	}
	sort.Ints(ov.walletIds)

	defer func() {
		go overviewHandler.UpdateBalance(app.MultiWallet)
		go overviewHandler.UpdateBlockStatusBox(app.MultiWallet)
	}()
	go func() {
		time.Sleep(200 * time.Millisecond)
		overviewHandler.UpdateWalletsSyncBox(app.MultiWallet)
	}()

	overviewHandler.Container = ov.container()
	return widget.NewHBox(widgets.NewHSpacer(18), overviewHandler.Container)
}

func (ov *overview) container() fyne.CanvasObject {
	overviewHandler.PageBoxes = ov.pageBoxes()
	overviewContainer := widget.NewVBox(
		title(),
		balance(),
		widgets.NewVSpacer(25),
		overviewHandler.PageBoxes,
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
	overviewHandler.UpdateTransactions(ov.multiWallet, handler.TransactionUpdate{})
	return widget.NewVBox(
		table.Result,
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

func (ov *overview) blockStatusBox() fyne.CanvasObject {
	if overviewHandler.Synced {
		return ov.blockStatusBoxInSync()
	}
	return ov.blockStatusBoxSyncing()
}

func (ov *overview) blockStatusBoxSyncing() fyne.CanvasObject {
	syncStatusText := widgets.NewSmallText("", color.Black)
	timeLeft := widget.NewLabelWithStyle("", fyne.TextAlignLeading, fyne.TextStyle{})
	connectedPeers := widget.NewLabelWithStyle("", fyne.TextAlignTrailing, fyne.TextStyle{})
	progressBar := widget.NewProgressBar()
	syncSteps := widget.NewLabelWithStyle("Step 0/3", fyne.TextAlignTrailing, fyne.TextStyle{})
	blockHeadersStatus := widget.NewLabelWithStyle("Fetching block headers  0%", fyne.TextAlignTrailing, fyne.TextStyle{})
	walletSyncInfo := fyne.NewContainerWithLayout(layout.NewHBoxLayout())
	walletSyncScrollContainer := widget.NewScrollContainer(walletSyncInfo)
	cancelButton := widget.NewButton("Cancel", func() {overviewHandler.CancelSync(ov.multiWallet)})

	overviewHandler.SyncStatusWidget = syncStatusText
	overviewHandler.TimeLeftWidget = timeLeft
	overviewHandler.ProgressBar = progressBar
	overviewHandler.ConnectedPeersWidget = connectedPeers
	overviewHandler.SyncStepWidget = syncSteps
	overviewHandler.BlockHeadersWidget = blockHeadersStatus
	overviewHandler.WalletSyncInfo = walletSyncInfo
	overviewHandler.Scroll = walletSyncScrollContainer
	overviewHandler.CancelButton = cancelButton
	overviewHandler.BlockHeightTime = widget.NewLabelWithStyle("", fyne.TextAlignLeading, fyne.TextStyle{})
	top := fyne.NewContainerWithLayout(layout.NewFixedGridLayout(fyne.NewSize(515, 24)),
		widget.NewHBox(
			syncStatusText,
			layout.NewSpacer(),
			cancelButton,
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
		fyne.NewContainerWithLayout(layout.NewFixedGridLayout(fyne.NewSize(510, 80)), walletSyncScrollContainer),
	)
	overviewHandler.BlockStatus = blockStatus
	return blockStatus
}

func (ov *overview) blockStatusBoxInSync() fyne.CanvasObject {
	syncStatusText := widgets.NewSmallText("", color.Black)
	blockHeightTime := widget.NewLabelWithStyle("", fyne.TextAlignLeading, fyne.TextStyle{})
	connectedPeers := widget.NewLabelWithStyle("", fyne.TextAlignTrailing, fyne.TextStyle{})

	overviewHandler.SyncStatusWidget = syncStatusText
	overviewHandler.BlockHeightTime = blockHeightTime
	overviewHandler.ConnectedPeersWidget = connectedPeers
	top := fyne.NewContainerWithLayout(layout.NewFixedGridLayout(fyne.NewSize(515, 24)),
		widget.NewHBox(
			syncStatusText,
			layout.NewSpacer(),
		))
	syncedStatus := fyne.NewContainerWithLayout(layout.NewBorderLayout(nil, nil, blockHeightTime, connectedPeers),
		blockHeightTime, connectedPeers)

	blockStatus := fyne.NewContainerWithLayout(layout.NewVBoxLayout(),
		widgets.NewVSpacer(5),
		top,
		syncedStatus,
		widgets.NewVSpacer(15),
	)
	overviewHandler.BlockStatus = blockStatus
	return blockStatus
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
