package handlers

import (
	"fmt"
	"fyne.io/fyne/layout"
	"github.com/gen2brain/beeep"
	"image/color"
	"log"
	"sort"
	"strconv"
	"strings"
	"sync"
	"time"

	"fyne.io/fyne"
	"fyne.io/fyne/canvas"
	"fyne.io/fyne/widget"

	"github.com/decred/dcrd/dcrutil"
	"github.com/raedahgroup/dcrlibwallet"
	"github.com/raedahgroup/godcr/fyne/assets"
	"github.com/raedahgroup/godcr/fyne/widgets"
)

var OverviewHandlerLock sync.Mutex

type OverviewHandler struct {
	Synced          bool
	Syncing         bool
	SyncProgress    float64
	ConnectedPeers  int32
	Steps           int32
	hideSyncDetails bool

	Balance              []*canvas.Text
	Transactions         []dcrlibwallet.Transaction
	PageBox              fyne.CanvasObject
	SyncStatusWidget     *canvas.Text
	TimeLeftWidget       *widget.Label
	ProgressBar          *widget.ProgressBar
	BlockStatus          *fyne.Container
	Table                *widgets.Table
	ConnectedPeersWidget *widget.Label
	SyncStepWidget       *widget.Label
	BlockHeadersWidget   *widget.Label
	WalletSyncInfo       *fyne.Container
	WalletSyncInfoToggle *widgets.ClickableBox
	WalletSyncInfoToggleText	 *widget.Label
	Scroll *widget.ScrollContainer
}

type TransactionUpdate struct {
	WalletId    int
	TxnHash     string
	Transaction dcrlibwallet.Transaction
}

func (handler *OverviewHandler) UpdateBalance(multiWallet *dcrlibwallet.MultiWallet) {
	tb, _ := totalBalance(multiWallet)
	mainBalance, subBalance := breakBalance(tb)
	handler.Balance[0].Text = mainBalance
	handler.Balance[1].Text = subBalance
	for _, w := range handler.Balance {
		canvas.Refresh(w)
	}
}

func (handler *OverviewHandler) UpdateBlockStatusBox(wallet *dcrlibwallet.MultiWallet) {
	handler.UpdateTimestamp(wallet, false)
	handler.UpdateSyncStatus(wallet, false)
	handler.UpdateConnectedPeers(wallet.ConnectedPeers(), false)
	widget.Refresh(handler.ProgressBar)
	canvas.Refresh(handler.BlockStatus)
}

func (handler *OverviewHandler) UpdateTransactions(wallet *dcrlibwallet.MultiWallet, update TransactionUpdate) {
	var transactionList []*widget.Box
	var height int32

	height = wallet.GetBestBlock().Height
	if len(handler.Table.Result.Children) > 1 {
		handler.Table.DeleteAll()
	}

	if update.Transaction.Hash != "" {
		var oldTransactions []dcrlibwallet.Transaction
		if len(handler.Transactions) > 0 {
			oldTransactions = handler.Transactions[1:]
		}

		handler.Transactions = []dcrlibwallet.Transaction{update.Transaction}
		handler.Transactions = append(handler.Transactions, oldTransactions...)
	} else if update.WalletId != 0 && update.TxnHash != "" {
		txn, err := wallet.WalletWithID(update.WalletId).GetTransactionRaw([]byte(update.TxnHash))
		if err != nil {
			log.Printf("error fetching transaction %v", err.Error())
			return
		}

		handler.Transactions = append(handler.Transactions, *txn)
	} else {
		var err error
		handler.Transactions, err = recentTransactions(wallet)
		if err != nil {
			log.Printf("error recentTransactions %v", err.Error())
			return
		}
	}

	for _, txn := range handler.Transactions {
		amount := dcrutil.Amount(txn.Amount).String()
		fee := dcrutil.Amount(txn.Fee).String()
		timeDate := dcrlibwallet.ExtractDateOrTime(txn.Timestamp)
		status := transactionStatus(height, txn)
		transactionList = append(transactionList, newTransactionRow(transactionIcon(txn.Direction), amount, fee,
			dcrlibwallet.TransactionDirectionName(txn.Direction), status, timeDate))
	}
	handler.Table.Append(transactionList...)
}

func (handler *OverviewHandler) UpdateSyncStatus(wallet *dcrlibwallet.MultiWallet, refresh bool) {
	status := handler.SyncStatusWidget
	progressBar := handler.ProgressBar
	handler.Syncing = wallet.IsSyncing()
	handler.Synced = wallet.IsSynced()
	if wallet.IsSynced() {
		handler.SyncProgress = 1
	}

	handler.UpdateProgressBar(false)
	status.Text, status.Color = handler.blockSyncStatus()
	if refresh {
		widget.Refresh(progressBar)
		canvas.Refresh(status)
	}
}

func (handler *OverviewHandler) UpdateSyncProgress(progressReport *dcrlibwallet.HeadersFetchProgressReport) {
	timeInString := strconv.Itoa(int(progressReport.GeneralSyncProgress.TotalTimeRemainingSeconds))
	handler.TimeLeftWidget.Text = timeInString + " secs left"
	widget.Refresh(handler.TimeLeftWidget)
}

func (handler *OverviewHandler) UpdateConnectedPeers(peers int32, refresh bool) {
	handler.ConnectedPeersWidget.SetText(fmt.Sprintf("Connected peers count  %d", peers))
	if refresh {
		widget.Refresh(handler.ConnectedPeersWidget)
	}
}

func (handler *OverviewHandler) UpdateTimestamp(wallet *dcrlibwallet.MultiWallet, refresh bool) {
	handler.TimeLeftWidget.SetText(bestBlockInfo(wallet.GetBestBlock().Height, wallet.GetBestBlock().Timestamp))
	if refresh {
		widget.Refresh(handler.TimeLeftWidget)
	}
}

func (handler *OverviewHandler) UpdateProgressBar(refresh bool) {
	handler.ProgressBar.Value = handler.SyncProgress
	if refresh {
		widget.Refresh(handler.ProgressBar)
	}
}

func (handler *OverviewHandler) UpdateSyncSteps(refresh bool) {
	handler.SyncStepWidget.Text = fmt.Sprintf("Step 1/%d", handler.Steps)
	if refresh {
		widget.Refresh(handler.SyncStepWidget)
	}
}

func (handler *OverviewHandler) UpdateBlockHeadersSync(status float64, refresh bool) {
	handler.BlockHeadersWidget.Text = fmt.Sprintf("Fetching block headers  %f%%", status)
	if refresh {
		widget.Refresh(handler.BlockHeadersWidget)
	}
}

func (handler *OverviewHandler) blockSyncStatus() (string, color.Color) {
	if handler.Syncing {
		return "Syncing...", color.Gray{Y: 123}
	}
	if handler.Synced {
		return "Synced", color.Gray{Y: 123}
	}

	return "Not Synced", color.Gray{Y: 123}
}

func (handler *OverviewHandler) CancelSync(wallet *dcrlibwallet.MultiWallet) {
	if wallet.IsSyncing() {
		wallet.CancelSync()
		err := beeep.Notify("Sync Canceled", "wallet sync stopped until app restart", "assets/information.png")
		if err != nil {
			log.Println("error initiating alert:", err.Error())
		}
	}
}

func (handler *OverviewHandler) UpdateWalletsSyncBox(wallet *dcrlibwallet.MultiWallet) {
	walletIds := wallet.OpenedWalletIDsRaw()
	sort.Ints(walletIds)
	for _, id := range walletIds {
		w := wallet.WalletWithID(id)
		headersFetched := fmt.Sprintf("%d of %d", w.GetBestBlock(), wallet.GetBestBlock().Height)
		progress := fmt.Sprintf("%s behind", dcrlibwallet.CalculateDaysBehind(w.GetBestBlockTimeStamp()))
		handler.WalletSyncInfo.AddObject(
			walletSyncBox(
				w.Name,
				walletSyncStatus(w, wallet.GetBestBlock().Height),
				headersFetched,
				progress,
			),
		)
	}
	handler.Scroll.Resize(fyne.NewSize(510, 70))
	widget.Refresh(handler.Scroll)
	canvas.Refresh(handler.WalletSyncInfo)
	fmt.Printf("SCROLLER %v",handler.Scroll.MinSize())
}

func (handler *OverviewHandler) HideWalletSyncBox() {
	if handler.hideSyncDetails {
		handler.hideSyncDetails = !handler.hideSyncDetails
		handler.WalletSyncInfoToggleText.Text = "show details"
	} else {
		handler.hideSyncDetails  = !handler.hideSyncDetails
		handler.WalletSyncInfoToggleText.Text = "hide details"
	}

	// handler.WalletSyncInfoToggle.Refresh()
}

func walletSyncStatus(wallet *dcrlibwallet.Wallet, bestBlockHeight int32) string {
	if wallet.IsWaiting() {
		return "waiting for other wallets"
	}
	if wallet.GetBestBlock() < bestBlockHeight {
		return "syncing"
	}
	return "synced"
}

func bestBlockInfo(blockHeight int32, timestamp int64) string {
	blockTimeStamp := time.Unix(timestamp, 0)
	timeLeft := time.Now().Sub(blockTimeStamp).Round(time.Second).String()
	return timeLeft + " ago"
}

func totalBalance(multiWallet *dcrlibwallet.MultiWallet) (balance string, err error) {
	var totalWalletBalance int64
	walletIds := multiWallet.OpenedWalletIDsRaw()
	sort.Ints(walletIds)
	for _, id := range walletIds {
		accounts, err := multiWallet.WalletWithID(id).GetAccountsRaw(dcrlibwallet.DefaultRequiredConfirmations)
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
	b2 = balanceParts[1]
	b1 = b1 + "." + b2[:2]
	b2 = b2[2:]
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

func transactionStatus(bestBlockHeight int32, txn dcrlibwallet.Transaction) string {
	confirmations := bestBlockHeight - txn.BlockHeight + 1
	if txn.BlockHeight != -1 && confirmations > dcrlibwallet.DefaultRequiredConfirmations {
		return "confirmed"
	}

	return "pending"
}

func recentTransactions(wallet *dcrlibwallet.MultiWallet) (transactions []dcrlibwallet.Transaction, err error) {
	walletIds := wallet.OpenedWalletIDsRaw()

	// add recent transactions of all wallets to a single slice
	for _, id := range walletIds {
		txns, err := wallet.WalletWithID(id).GetTransactionsRaw(0, 10, 0, true)
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

func newTransactionRow(transactionType, amount, fee, direction, status, date string) *widget.Box {
	icons, _ := assets.GetIcons(assets.ReceiveIcon, assets.SendIcon)
	icon := canvas.NewImageFromResource(icons[transactionType])
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
