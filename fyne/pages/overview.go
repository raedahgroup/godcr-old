package pages

import (
	"fmt"
	"image/color"
	"sort"
	"strings"
	"time"

	"fyne.io/fyne"
	"fyne.io/fyne/canvas"
	"fyne.io/fyne/layout"
	"fyne.io/fyne/widget"

	"github.com/decred/dcrd/dcrutil"
	"github.com/raedahgroup/dcrlibwallet"
	"github.com/raedahgroup/godcr/fyne/assets"
	"github.com/raedahgroup/godcr/fyne/pages/handler/values"
	"github.com/raedahgroup/godcr/fyne/widgets"
)

const historyPageIndex = 1
const PageTitle = "Overview"

type overview struct {
	app            *AppInterface
	transactionBox *widget.Box
	multiWallet    *dcrlibwallet.MultiWallet
	walletIds      []int
	transactions   []dcrlibwallet.Transaction
}

// todo: display overview page (include sync progress UI elements)
// todo: register sync progress listener on overview page to update sync progress views
func overviewPageContent(app *AppInterface) fyne.CanvasObject {
	ov := &overview{}
	ov.app = app
	// app.Window.Resize(fyne.NewSize(650, 650))
	ov.multiWallet = app.MultiWallet
	ov.walletIds = ov.multiWallet.OpenedWalletIDsRaw()
	if len(ov.walletIds) == 0 {
		return widget.NewHBox(widgets.NewHSpacer(10), widget.NewLabelWithStyle("Could not retrieve wallets", fyne.TextAlignCenter, fyne.TextStyle{Bold: true}))
	}
	sort.Ints(ov.walletIds)
	return widget.NewHBox(widgets.NewHSpacer(values.Padding), ov.container(), widgets.NewHSpacer(values.Padding))
}

func (ov *overview) container() fyne.CanvasObject {
	return widget.NewVBox(
		widgets.NewVSpacer(values.Padding),
		title(),
		ov.balance(),
		widgets.NewVSpacer(50),
		ov.pageBoxes(),
		widgets.NewVSpacer(values.Padding),
	)
}

func title() fyne.CanvasObject {
	titleWidget := widget.NewLabelWithStyle(PageTitle, fyne.TextAlignLeading, fyne.TextStyle{Bold: true})
	return widget.NewHBox(titleWidget)
}

func (ov *overview) balance() fyne.CanvasObject {
	tb, err := totalBalance(ov)
	if err != nil {
		return widget.NewLabel(fmt.Sprintf("Error: %s", err.Error()))
	}
	mainBalance, subBalance := breakBalance(tb)
	dcrBalance := widgets.NewLargeText(mainBalance, color.Black)
	dcrDecimals := widgets.NewSmallText(subBalance, color.Black)
	decimalsBox := fyne.NewContainerWithLayout(layout.NewVBoxLayout(), widgets.NewVSpacer(6), dcrDecimals)
	return widget.NewHBox(widgets.NewVSpacer(10), dcrBalance, decimalsBox)
}

func (ov *overview) pageBoxes() (object fyne.CanvasObject) {
	return fyne.NewContainerWithLayout(layout.NewVBoxLayout(),
		blockStatusBox(),
		widgets.NewVSpacer(15),
		ov.recentTransactionBox(),
	)
}

func (ov *overview) recentTransactionBox() fyne.CanvasObject {
	var err error
	ov.transactions, err = recentTransactions(ov)
	if err != nil {
		return widget.NewHBox(widgets.NewHSpacer(10), widget.NewLabelWithStyle(err.Error(), fyne.TextAlignCenter, fyne.TextStyle{Bold: true}))
	}

	table := &widgets.Table{}
	table.NewTable(transactionRowHeader())
	for _, txn := range ov.transactions {
		amount := dcrutil.Amount(txn.Amount).String()
		fee := dcrutil.Amount(txn.Fee).String()
		timeDate := dcrlibwallet.ExtractDateOrTime(txn.Timestamp)
		status := transactionStatus(ov, txn)
		table.Append(newTransactionRow(transactionIcon(txn.Direction), amount, fee,
			dcrlibwallet.TransactionDirectionName(txn.Direction), status, timeDate))
	}

	return widget.NewVBox(
		table.Result,
		fyne.NewContainerWithLayout(layout.NewHBoxLayout(),
			layout.NewSpacer(),
			widgets.NewClickableWidget(widget.NewLabelWithStyle("see all", fyne.TextAlignCenter, fyne.TextStyle{Italic: true}),
				func() {
					ov.app.tabMenu.SelectTabIndex(historyPageIndex)
				},
			),
			layout.NewSpacer(),
		),
	)
}

func blockStatusBox() fyne.CanvasObject {
	top := fyne.NewContainerWithLayout(layout.NewFixedGridLayout(fyne.NewSize(515, 24)),
		widget.NewHBox(
			widgets.NewSmallText("Syncing...", color.Black),
			layout.NewSpacer(),
			widget.NewButton("Cancel", func() {}),
		))
	progressBar := fyne.NewContainerWithLayout(layout.NewFixedGridLayout(fyne.NewSize(515, 20)),
		widget.NewProgressBar(),
	)
	timeLeft := widget.NewLabelWithStyle("6 min left", fyne.TextAlignLeading, fyne.TextStyle{Italic: true})
	connectedPeers := widget.NewLabelWithStyle("Connected peers count  6", fyne.TextAlignTrailing, fyne.TextStyle{Italic: true})
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

	return fyne.NewContainerWithLayout(layout.NewVBoxLayout(),
		widgets.NewVSpacer(5),
		top,
		progressBar,
		syncDuration,
		syncStatus,
		widgets.NewVSpacer(15),
		bottom,
	)
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

func newTransactionRow(transactionType, amount, fee, direction, status, date string) *widget.Box {
	icons, _ := assets.GetIcons(assets.ReceiveIcon, assets.SendIcon)
	icon := canvas.NewImageFromResource(icons[transactionType])
	// spacer := widgets.NewHSpacer(10)
	icon.SetMinSize(fyne.NewSize(5, 20))
	iconBox := widget.NewVBox(widgets.NewVSpacer(4), icon)
	amountLabel := widget.NewLabel(amount)
	feeLabel := widget.NewLabel(fee)
	dateLabel := widget.NewLabel(date)
	statusLabel := widget.NewLabel(status)
	directionLabel := widget.NewLabel(direction)
	column := widget.NewHBox(iconBox, amountLabel, feeLabel, directionLabel, statusLabel, dateLabel)
	return column
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

func totalBalance(overview *overview) (balance string, err error) {
	var totalWalletBalance int64
	mw := overview.multiWallet
	for _, id := range overview.walletIds {
		accounts, err := mw.WalletWithID(id).GetAccountsRaw(dcrlibwallet.DefaultRequiredConfirmations)
		if err != nil {
			return "", err
		}
		for _, acc := range accounts.Acc {
			totalWalletBalance += acc.TotalBalance
		}
	}
	balance = dcrutil.Amount(totalWalletBalance).String()
	return
}

func breakBalance(balance string) (b1, b2 string) {
	balanceParts := strings.Split(balance, ".")
	b1 = balanceParts[0]
	if len(balanceParts) > 1 {
		b2 = balanceParts[1]
		b1 = b1 + "." + b2[:2]
		b2 = b2[2:]
	}
	return
}

func recentTransactions(overview *overview) (transactions []dcrlibwallet.Transaction, err error) {
	mw := overview.multiWallet

	// add recent transactions of all wallets to a single slice
	for _, id := range overview.walletIds {
		txns, err := mw.WalletWithID(id).GetTransactionsRaw(0, 10, 0, true)
		transactions = append(transactions, txns...)
		if err != nil {
			return nil, err
		}
	}
	sort.SliceStable(transactions, func(i, j int) bool {
		backTime := time.Unix(transactions[j].Timestamp, 0)
		frontTime := time.Unix(transactions[i].Timestamp, 0)
		return backTime.Before(frontTime)
	})
	if len(transactions) > 5 {
		transactions = transactions[:5]
	}
	return
}

func transactionIcon(direction int32) string {
	switch direction {
	case 0:
		return assets.SendIcon
	case 1:
		return assets.ReceiveIcon
	case 2:
		return assets.ReceiveIcon
	default:
		return assets.InfoIcon
	}
}

func transactionStatus(overview *overview, txn dcrlibwallet.Transaction) string {
	confirmations := overview.multiWallet.GetBestBlock().Height - txn.BlockHeight + 1
	if txn.BlockHeight != -1 && confirmations > dcrlibwallet.DefaultRequiredConfirmations {
		return "confirmed"
	}
	return "pending"
}
