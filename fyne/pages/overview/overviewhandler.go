package overview

import (
	"fmt"
	"fyne.io/fyne/layout"
	"github.com/gen2brain/beeep"
	"image/color"
	"log"
	"sort"
	"strconv"
	"strings"
	"time"

	"fyne.io/fyne"
	"fyne.io/fyne/canvas"
	"fyne.io/fyne/widget"

	"github.com/decred/dcrd/dcrutil"
	"github.com/raedahgroup/dcrlibwallet"
	"github.com/raedahgroup/godcr/fyne/assets"
	"github.com/raedahgroup/godcr/fyne/widgets"
)

type Handler struct {
	Synced          bool
	Syncing         bool
	SyncProgress    float64
	ConnectedPeers  int32
	Steps           int32
	hideSyncDetails bool

	Container            fyne.CanvasObject
	Balance              []*canvas.Text
	Transactions         []dcrlibwallet.Transaction
	PageBoxes            fyne.CanvasObject
	SyncStatusWidget     *canvas.Text
	TimeLeftWidget       *widget.Label
	BlockHeightTime      *widget.Label
	ProgressBar          *widget.ProgressBar
	BlockStatus          *fyne.Container
	BlockStatusSynced    fyne.CanvasObject
	BlockStatusSyncing   fyne.CanvasObject
	Table                *widgets.Table
	ConnectedPeersWidget *widget.Label
	SyncStepWidget       *widget.Label
	BlockHeadersWidget   *widget.Label
	WalletSyncInfo       *fyne.Container
	Scroll               *widget.ScrollContainer
	CancelButton         *widget.Button

	StepsChannel chan int32
}

type TransactionUpdate struct {
	WalletId    int
	TxnHash     string
	Transaction dcrlibwallet.Transaction
}

func (h *Handler) UpdateBalance(multiWallet *dcrlibwallet.MultiWallet) {
	tb, _ := totalBalance(multiWallet)
	mainBalance, subBalance := breakBalance(tb)
	h.Balance[0].Text = mainBalance
	h.Balance[1].Text = subBalance
	for _, w := range h.Balance {
		w.Refresh()
	}
}

func (h *Handler) UpdateBlockStatusBox(wallet *dcrlibwallet.MultiWallet) {
	h.UpdateSyncStatus(wallet, false)
	if h.Synced {
		h.BlockStatus.Objects = []fyne.CanvasObject{}
		h.BlockStatus.AddObject(h.BlockStatusSynced)
		h.UpdateConnectedPeers(wallet.ConnectedPeers(), false)
		h.UpdateBlockHeightTime(wallet, false)
	} else {
		h.BlockStatus.Objects = []fyne.CanvasObject{}
		h.BlockStatus.AddObject(h.BlockStatusSyncing)
		h.UpdateTimestamp(wallet, false)
		h.UpdateConnectedPeers(wallet.ConnectedPeers(), false)
		h.UpdateProgressBar(false)
	}

	h.Container.Refresh()
}

func (h *Handler) UpdateTransactions(wallet *dcrlibwallet.MultiWallet, update TransactionUpdate) {
	var transactionList []*widget.Box
	var height int32

	height = wallet.GetBestBlock().Height
	if len(h.Table.Result.Children) > 1 {
		h.Table.DeleteAll()
	}

	if update.Transaction.Hash != "" {
		var oldTransactions []dcrlibwallet.Transaction
		if len(h.Transactions) > 0 {
			oldTransactions = h.Transactions[1:]
		}

		h.Transactions = []dcrlibwallet.Transaction{update.Transaction}
		h.Transactions = append(h.Transactions, oldTransactions...)
	} else if update.WalletId != 0 && update.TxnHash != "" {
		txn, err := wallet.WalletWithID(update.WalletId).GetTransactionRaw([]byte(update.TxnHash))
		if err != nil {
			log.Printf("error fetching transaction %v", err.Error())
			return
		}

		h.Transactions = append(h.Transactions, *txn)
	} else {
		var err error
		h.Transactions, err = recentTransactions(wallet)
		if err != nil {
			log.Printf("error recentTransactions %v", err.Error())
			return
		}
	}

	for _, txn := range h.Transactions {
		amount := dcrutil.Amount(txn.Amount).String()
		fee := dcrutil.Amount(txn.Fee).String()
		timeDate := dcrlibwallet.ExtractDateOrTime(txn.Timestamp)
		status := transactionStatus(height, txn)
		transactionList = append(transactionList, newTransactionRow(transactionIcon(txn.Direction), amount, fee,
			dcrlibwallet.TransactionDirectionName(txn.Direction), status, timeDate))
	}
	h.Table.Append(transactionList...)
}

func (h *Handler) UpdateSyncStatus(wallet *dcrlibwallet.MultiWallet, refresh bool) {
	status := h.SyncStatusWidget
	progressBar := h.ProgressBar
	h.Syncing = wallet.IsSyncing()
	h.Synced = wallet.IsSynced()
	if h.Synced {
		h.SyncProgress = 1
		h.hideSyncDetails = true
	}

	h.UpdateProgressBar(false)
	status.Text, status.Color = h.blockSyncStatus()
	if refresh {
		progressBar.Refresh()
		status.Refresh()
	}
}

func (h *Handler) UpdateSyncProgress(progressReport *dcrlibwallet.HeadersFetchProgressReport) {
	timeInString := strconv.Itoa(int(progressReport.GeneralSyncProgress.TotalTimeRemainingSeconds))
	h.TimeLeftWidget.Text = timeInString + " secs left"
	h.TimeLeftWidget.Refresh()
}

func (h *Handler) UpdateConnectedPeers(peers int32, refresh bool) {
	h.ConnectedPeersWidget.SetText(fmt.Sprintf("Connected peers count  %d", peers))
	if refresh {
		h.ConnectedPeersWidget.Refresh()
	}
}

func (h *Handler) UpdateTimestamp(wallet *dcrlibwallet.MultiWallet, refresh bool) {
	h.TimeLeftWidget.SetText(bestBlockInfo(wallet.GetBestBlock().Height, wallet.GetBestBlock().Timestamp))
	if refresh {
		h.TimeLeftWidget.Refresh()
	}
}

func (h *Handler) UpdateProgressBar(refresh bool) {
	h.ProgressBar.Value = h.SyncProgress
	if refresh {
		h.ProgressBar.Refresh()
	}
}

func (h *Handler) UpdateSyncSteps(refresh bool) {
	h.SyncStepWidget.SetText(fmt.Sprintf("Step %d/3", h.Steps))
	if refresh {
		h.SyncStepWidget.Refresh()
	}
}

func (h *Handler) UpdateBlockHeadersSync(status int32, refresh bool) {
	h.BlockHeadersWidget.SetText(fmt.Sprintf("Fetching block headers  %v%%", status))
	if refresh {
		h.BlockHeadersWidget.Refresh()
	}
}

func (h *Handler) blockSyncStatus() (string, color.Color) {
	if h.Syncing {
		return "Syncing...", color.Gray{Y: 123}
	}
	if h.Synced {
		return "Synced", color.Gray{Y: 123}
	}

	return "Not Synced", color.Gray{Y: 123}
}

func (h *Handler) CancelSync(wallet *dcrlibwallet.MultiWallet) {
	if wallet.IsSyncing() {
		wallet.CancelSync()
		h.ConnectedPeers = 0
		err := beeep.Notify("Sync Canceled", "wallet sync stopped until app restart", "assets/information.png")
		if err != nil {
			log.Println("error initiating alert:", err.Error())
		}
	}
}

func (h *Handler) UpdateWalletsSyncBox(wallet *dcrlibwallet.MultiWallet) {
	walletIds := wallet.OpenedWalletIDsRaw()
	sort.Ints(walletIds)
	for _, id := range walletIds {
		w := wallet.WalletWithID(id)
		headersFetched := fmt.Sprintf("%d of %d", w.GetBestBlock(), wallet.GetBestBlock().Height)
		progress := fmt.Sprintf("%s behind", dcrlibwallet.CalculateDaysBehind(w.GetBestBlockTimeStamp()))
		h.WalletSyncInfo.Objects = []fyne.CanvasObject{}
		h.WalletSyncInfo.AddObject(
			walletSyncBox(
				w.Name,
				walletSyncStatus(w, wallet.GetBestBlock().Height),
				headersFetched,
				progress,
			),
		)
	}

	h.WalletSyncInfo.Refresh()
	return
}

func (h *Handler) UpdateBlockHeightTime(wallet *dcrlibwallet.MultiWallet, refresh bool) {
	blockTimeStamp := time.Unix(wallet.GetBestBlock().Timestamp, 0)
	timeLeft := time.Now().Sub(blockTimeStamp).Round(time.Second).String()
	text := fmt.Sprintf("Latest block %v . %v ago", wallet.GetBestBlock().Height, timeLeft)
	h.BlockHeightTime.SetText(text)
	if refresh {
		h.BlockHeightTime.Refresh()
	}
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
	if len(balanceParts) == 1 {
		return balanceParts[0], ""
	}
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

func (h *Handler) PreserveSyncSteps() {
	h.StepsChannel = make(chan int32)
	defer func() {
		fmt.Printf("\n \n Preserve sync step go routine has been killed!\n \n")
	}()
	h.Steps = 0
	for {
		select {
		case progress := <-h.StepsChannel:
			if progress == 100 && h.Steps < 3 {
				h.Steps += 1
			}
		}
	}
}
