package handlers

import (
	"fmt"
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
	Synced         bool
	Syncing        bool
	SyncProgress   float64
	ConnectedPeers int32

	Balance            []*canvas.Text
	Transactions       []dcrlibwallet.Transaction
	PageBox            fyne.CanvasObject
	SyncStatusText     *canvas.Text
	TimeLeftText       *widget.Label
	ProgressBar        *widget.ProgressBar
	BlockStatus        *fyne.Container
	Table              *widgets.Table
	ConnectedPeersText *widget.Label
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

func (handler *OverviewHandler) UpdateSyncStatus (wallet *dcrlibwallet.MultiWallet, refresh bool) {
	status := handler.SyncStatusText
	progressBar := handler.ProgressBar
	handler.Syncing = wallet.IsSyncing()
	handler.Synced = wallet.IsSynced()
	if wallet.IsSynced() {
		handler.SyncProgress = 1
	}
	status.Text, status.Color, progressBar.Value = handler.blockSyncStatus()
	if refresh {
		widget.Refresh(progressBar)
		canvas.Refresh(status)
	}
}

func (handler *OverviewHandler) UpdateSyncProgress(progressReport *dcrlibwallet.HeadersFetchProgressReport) {
	timeInString := strconv.Itoa(int(progressReport.GeneralSyncProgress.TotalTimeRemainingSeconds))
	handler.TimeLeftText.Text = timeInString + " secs left"
	widget.Refresh(handler.TimeLeftText)
}

func (handler *OverviewHandler) UpdateConnectedPeers(peers int32, refresh bool) {
	handler.ConnectedPeersText.SetText(fmt.Sprintf("Connected peers count  %d", peers))
	if refresh {
		widget.Refresh(handler.ConnectedPeersText)
	}
}

func (handler *OverviewHandler) UpdateTimestamp(wallet *dcrlibwallet.MultiWallet, refresh bool) {
	handler.TimeLeftText.SetText(bestBlockInfo(wallet.GetBestBlock().Height, wallet.GetBestBlock().Timestamp))
	if refresh {
		widget.Refresh(handler.TimeLeftText)
	}
}

func (handler *OverviewHandler) blockSyncStatus() (string, color.Color, float64) {
	if handler.Syncing {
		return "Syncing...", color.Gray{Y: 123}, handler.SyncProgress
	}
	if handler.Synced {
		return "Synced", color.Gray{Y: 123}, handler.SyncProgress
	}

	return "Not Synced", color.Gray{Y: 123}, handler.SyncProgress
}

func bestBlockInfo(blockHeight int32, timestamp int64) string {
	blockTimeStamp := time.Unix(timestamp, 0)
	timeLeft := time.Now().Sub(blockTimeStamp).Round(time.Second).String()
	return timeLeft + " ago"
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
