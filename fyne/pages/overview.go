package pages

import (
	"image/color"
	"sort"
	"time"

	"fyne.io/fyne"
	"fyne.io/fyne/canvas"
	"fyne.io/fyne/layout"
	"fyne.io/fyne/widget"

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
		blockStatusBox(),
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

func blockStatusBox() fyne.CanvasObject {
	syncStatusText := widgets.NewSmallText("", color.Black)
	timeLeft := widget.NewLabelWithStyle("", fyne.TextAlignLeading, fyne.TextStyle{Italic: true})
	connectedPeers := widget.NewLabelWithStyle("", fyne.TextAlignTrailing, fyne.TextStyle{Italic: true})
	progressBar := widget.NewProgressBar()
	overviewHandler.SyncStatusText = syncStatusText
	overviewHandler.TimeLeftText = timeLeft
	overviewHandler.ProgressBar = progressBar
	overviewHandler.ConnectedPeersText = connectedPeers

	top := fyne.NewContainerWithLayout(layout.NewFixedGridLayout(fyne.NewSize(515, 24)),
		widget.NewHBox(
			syncStatusText,
			layout.NewSpacer(),
			widget.NewButton("Cancel", func() {}),
		))
	progressBarContainer := fyne.NewContainerWithLayout(layout.NewFixedGridLayout(fyne.NewSize(515, 20)),
		progressBar,
	)
	syncSteps := widget.NewLabelWithStyle("Step 1/3", fyne.TextAlignTrailing, fyne.TextStyle{Italic: true})
	blockHeadersStatus := widget.NewLabelWithStyle("Fetching block headers  89%", fyne.TextAlignTrailing, fyne.TextStyle{Italic: true})
	syncDuration := fyne.NewContainerWithLayout(layout.NewBorderLayout(nil, nil, timeLeft, connectedPeers),
		timeLeft, connectedPeers)
	syncStatus := fyne.NewContainerWithLayout(layout.NewBorderLayout(nil, nil, syncSteps, blockHeadersStatus),
		syncSteps, blockHeadersStatus)

	bottom := fyne.NewContainerWithLayout(layout.NewGridLayout(2),
		walletSyncBox("Default", "waiting for other wallets", "6000 of 164864", "220 days behind"),
		walletSyncBox("Wallet 2", "Syncing", "100 of 164864", "320 days behind"),
	)
	blockStatus := fyne.NewContainerWithLayout(layout.NewVBoxLayout(),
		widgets.NewVSpacer(5),
		top,
		progressBarContainer,
		syncDuration,
		syncStatus,
		widgets.NewVSpacer(15),
	)
	overviewHandler.BlockStatus = blockStatus
	go func() {
		time.Sleep(time.Millisecond * 200)
		blockStatus.AddObject(bottom)
		canvas.Refresh(blockStatus)
	}()
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

func walletSyncBox(name, status, headerFetched, progress string) fyne.CanvasObject {
	blackColor := color.Black
	nameText := widgets.NewTextWithSize(name, blackColor, 12)
	statusText := widgets.NewTextWithSize(status, blackColor, 10)
	headerFetchedTitleText := widgets.NewTextWithSize("Block header fetched", blackColor, 12)
	headerFetchedText := widgets.NewTextWithSize(headerFetched, blackColor, 10)
	progressTitleText := widgets.NewTextWithSize("Syncing progress", blackColor, 12)
	progressText := widgets.NewTextWithSize(progress, blackColor, 10)
	top := fyne.NewContainerWithLayout(layout.NewHBoxLayout(),
		widgets.NewHSpacer(2),
		nameText, layout.NewSpacer(),
		statusText,
		widgets.NewHSpacer(2))
	middle := fyne.NewContainerWithLayout(layout.NewHBoxLayout(),
		widgets.NewHSpacer(2),
		headerFetchedTitleText,
		layout.NewSpacer(),
		headerFetchedText,
		widgets.NewHSpacer(2),
	)
	bottom := fyne.NewContainerWithLayout(layout.NewHBoxLayout(),
		widgets.NewHSpacer(2),
		progressTitleText,
		layout.NewSpacer(),
		progressText,
		widgets.NewHSpacer(2),
	)
	background := canvas.NewRectangle(color.RGBA{0, 0, 0, 7})
	background.SetMinSize(fyne.NewSize(250, 70))
	walletSyncContent := fyne.NewContainerWithLayout(layout.NewFixedGridLayout(fyne.NewSize(250, 70)),
		fyne.NewContainerWithLayout(layout.NewVBoxLayout(), top, layout.NewSpacer(), middle, layout.NewSpacer(), bottom),
	)

	return fyne.NewContainerWithLayout(layout.NewFixedGridLayout(fyne.NewSize(250, 70)),
		fyne.NewContainerWithLayout(layout.NewBorderLayout(nil, nil, nil, nil),
			background,
			walletSyncContent,
		),
	)
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
