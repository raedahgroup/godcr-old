package handler

import (
	"fmt"
	"fyne.io/fyne/layout"
	"github.com/gen2brain/beeep"
	"github.com/raedahgroup/godcr/fyne/pages/handler/values"
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

type OverviewHandler struct {
	Synced          bool
	Syncing         bool
	SyncProgress    float64
	ConnectedPeers  int32
	Steps           int32
	hideSyncDetails bool

	Container             fyne.CanvasObject
	TransactionsContainer *fyne.Container
	Balance               []*canvas.Text
	Transactions          []dcrlibwallet.Transaction
	PageBoxes             fyne.CanvasObject
	SyncStatusWidget      *canvas.Text
	TimeLeftWidget        *widget.Label
	BlockHeightTime       *widget.Label
	ProgressBar           *widget.ProgressBar
	BlockStatus           *fyne.Container
	BlockStatusSynced     fyne.CanvasObject
	BlockStatusSyncing    fyne.CanvasObject
	Table                 *widgets.Table
	TableHeader           *widget.Box
	ConnectedPeersWidget  *widget.Label
	SyncStepWidget        *widget.Label
	BlockHeadersWidget    *widget.Label
	WalletSyncInfo        *fyne.Container
	Scroll                *widget.ScrollContainer
	SyncButton            *widget.Button

	StepsChannel chan int32
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
		w.Refresh()
	}
}

func (handler *OverviewHandler) UpdateBlockStatusBox(wallet *dcrlibwallet.MultiWallet) {
	handler.UpdateSyncStatus(wallet, false)
	if !handler.Syncing && !handler.Synced {
		handler.SyncButton.SetText("Reconnect")
	} else {
		handler.SyncButton.SetText("Cancel")
	}

	if handler.Synced {
		handler.BlockStatus.Objects = []fyne.CanvasObject{}
		handler.SyncButton.SetText("Disconnect")
		handler.BlockStatus.AddObject(handler.BlockStatusSynced)
		handler.UpdateConnectedPeers(wallet.ConnectedPeers(), false)
		handler.UpdateBlockHeightTime(wallet, false)
	} else {
		handler.BlockStatus.Objects = []fyne.CanvasObject{}
		handler.BlockStatus.AddObject(handler.BlockStatusSyncing)
		handler.UpdateTimestamp(wallet, false)
		handler.UpdateConnectedPeers(wallet.ConnectedPeers(), false)
		handler.UpdateProgressBar(false)
	}

	handler.Container.Refresh()
}

func (handler *OverviewHandler) UpdateTransactions(wallet *dcrlibwallet.MultiWallet, update TransactionUpdate) {
	var transactionList []*widget.Box
	var height int32

	// handler.TransactionsContainer.Hide()
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
	if handler.Synced {
		handler.SyncProgress = 1
		handler.hideSyncDetails = true
	}

	handler.UpdateProgressBar(false)
	status.Text, status.Color = handler.blockSyncStatus()
	if refresh {
		progressBar.Refresh()
		status.Refresh()
	}
}

func (handler *OverviewHandler) UpdateSyncProgress(progressReport *dcrlibwallet.HeadersFetchProgressReport) {
	timeInString := strconv.Itoa(int(progressReport.GeneralSyncProgress.TotalTimeRemainingSeconds))
	handler.TimeLeftWidget.Text = timeInString + " secs left"
	handler.TimeLeftWidget.Refresh()
}

func (handler *OverviewHandler) UpdateConnectedPeers(peers int32, refresh bool) {
	handler.ConnectedPeersWidget.SetText(fmt.Sprintf("Connected peers count  %d", peers))
	if refresh {
		handler.ConnectedPeersWidget.Refresh()
	}
}

func (handler *OverviewHandler) UpdateTimestamp(wallet *dcrlibwallet.MultiWallet, refresh bool) {
	handler.TimeLeftWidget.SetText(bestBlockInfo(wallet.GetBestBlock().Height, wallet.GetBestBlock().Timestamp))
	if refresh {
		handler.TimeLeftWidget.Refresh()
	}
}

func (handler *OverviewHandler) UpdateProgressBar(refresh bool) {
	handler.ProgressBar.Value = handler.SyncProgress
	if refresh {
		handler.ProgressBar.Refresh()
	}
}

func (handler *OverviewHandler) UpdateSyncSteps(refresh bool) {
	fmt.Printf("handler steps %v\n \n", handler.Steps)
	handler.SyncStepWidget.SetText(fmt.Sprintf("Step %d/3", handler.Steps))
	if refresh {
		handler.SyncStepWidget.Refresh()
	}
}

func (handler *OverviewHandler) UpdateBlockHeadersSync(status int32, refresh bool) {
	handler.BlockHeadersWidget.SetText(fmt.Sprintf("Fetching block headers  %v%%", status))
	if refresh {
		handler.BlockHeadersWidget.Refresh()
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

func (handler *OverviewHandler) SyncTrigger(wallet *dcrlibwallet.MultiWallet) {
	var notifyTitle, notifyMessage string
	synced, syncing := wallet.IsSynced(), wallet.IsSyncing()
	if wallet.ConnectedPeers() > 0 && (synced || syncing) {
		if synced {
			notifyTitle = "Wallet Disconnected"
			notifyMessage = "wallet has been disconnected from peers"
		} else if syncing {
			notifyTitle = "Sync Canceled"
			notifyMessage = "wallet sync stopped until wallet reconnect"
		}
		wallet.CancelSync()
		handler.ConnectedPeers = 0
		err := beeep.Notify(notifyTitle, notifyMessage, "assets/information.png")
		if err != nil {
			log.Println("error initiating alert:", err.Error())
		}
	} else if !wallet.IsSynced() && !wallet.IsSyncing() {
		err := wallet.RestartSpvSync()
		if err != nil {
			notifyError := beeep.Notify("Sync Error", "error restarting spv sync", "assets/information.png")
			if notifyError != nil {
				log.Println("error initiating alert:", err.Error())
			}
			log.Printf("error restarting SPV sync %v", err.Error())
		}
	}
	handler.Container.Refresh()
}

func (handler *OverviewHandler) UpdateWalletsSyncBox(wallet *dcrlibwallet.MultiWallet) {
	walletIds := wallet.OpenedWalletIDsRaw()
	sort.Ints(walletIds)

	handler.WalletSyncInfo.Objects = []fyne.CanvasObject{}
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

	handler.WalletSyncInfo.Refresh()
	return
}

func (handler *OverviewHandler) UpdateBlockHeightTime(wallet *dcrlibwallet.MultiWallet, refresh bool) {
	blockTimeStamp := time.Unix(wallet.GetBestBlock().Timestamp, 0)
	timeLeft := time.Now().Sub(blockTimeStamp).Round(time.Second).String()
	text := fmt.Sprintf("Latest block %v . %v ago", wallet.GetBestBlock().Height, timeLeft)
	handler.BlockHeightTime.SetText(text)
	if refresh {
		handler.BlockHeightTime.Refresh()
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
	transactions = []dcrlibwallet.Transaction{}
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
	amountLabel := widget.NewLabelWithStyle(amount, fyne.TextAlignTrailing, fyne.TextStyle{})
	feeLabel := widget.NewLabelWithStyle(fee, fyne.TextAlignCenter, fyne.TextStyle{})
	dateLabel := widget.NewLabelWithStyle(date, fyne.TextAlignCenter, fyne.TextStyle{})
	statusLabel := widget.NewLabelWithStyle(status, fyne.TextAlignCenter, fyne.TextStyle{})
	directionLabel := widget.NewLabelWithStyle(direction, fyne.TextAlignCenter, fyne.TextStyle{})
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
	background.SetMinSize(fyne.NewSize(250, 80))
	walletSyncContent := fyne.NewContainerWithLayout(layout.NewFixedGridLayout(fyne.NewSize(250, 80)),
		fyne.NewContainerWithLayout(layout.NewVBoxLayout(), widgets.NewVSpacer(values.SpacerSize2), top, layout.NewSpacer(),
			middle, layout.NewSpacer(), bottom, widgets.NewVSpacer(values.SpacerSize2)),
	)

	return fyne.NewContainerWithLayout(layout.NewFixedGridLayout(fyne.NewSize(250, 80)),
		fyne.NewContainerWithLayout(layout.NewBorderLayout(nil, nil, nil, nil),
			background,
			walletSyncContent,
		),
	)
}

func (handler *OverviewHandler) PreserveSyncSteps() {
	handler.StepsChannel = make(chan int32)
	defer func() {
		fmt.Printf("\n \n Preserve sync step go routine has been killed!\n \n")
	}()
	handler.Steps = 1
	for {
		select {
		case progress := <-handler.StepsChannel:
			if progress == 100 && handler.Steps < 3 {
				handler.Steps += 1
			}
		}
	}
}
