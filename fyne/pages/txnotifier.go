package pages

import (
	"encoding/json"
	"fmt"
	"github.com/raedahgroup/godcr/fyne/handlers"
	"log"
	"math"
	"strconv"

	"fyne.io/fyne/widget"

	"github.com/gen2brain/beeep"
	"github.com/raedahgroup/dcrlibwallet"

	"github.com/raedahgroup/godcr/fyne/pages/handler/multipagecomponents"
)

type multiWalletTxListener struct {
	tabMenu     *widget.TabContainer
	multiWallet *dcrlibwallet.MultiWallet
}

func (app *multiWalletTxListener) OnSyncStarted() {
	mw := app.multiWallet
	handlers.OverviewHandlerLock.Lock()
	defer handlers.OverviewHandlerLock.Unlock()
	overviewHandler.Synced = mw.IsSynced()
	overviewHandler.Syncing = mw.IsSyncing()
	overviewHandler.UpdateBlockStatusBox(mw)
}

func (app *multiWalletTxListener) OnPeerConnectedOrDisconnected(numberOfConnectedPeers int32) {
	overviewHandler.UpdateConnectedPeers(numberOfConnectedPeers, true)
}

func (app *multiWalletTxListener) OnHeadersFetchProgress(headersFetchProgress *dcrlibwallet.HeadersFetchProgressReport) {
	handlers.OverviewHandlerLock.Lock()
	defer handlers.OverviewHandlerLock.Unlock()
	overviewHandler.SyncProgress = float64(headersFetchProgress.FetchedHeadersCount) / float64(headersFetchProgress.TotalHeadersToFetch)
	overviewHandler.SyncProgress = math.Round(overviewHandler.SyncProgress * 100) / 100
	overviewHandler.UpdateBlockStatusBox(app.multiWallet)
}

func (app *multiWalletTxListener) OnAddressDiscoveryProgress(addressDiscoveryProgress *dcrlibwallet.AddressDiscoveryProgressReport) {

}

func (app *multiWalletTxListener) OnHeadersRescanProgress(headersRescanProgress *dcrlibwallet.HeadersRescanProgressReport) {

}

func (app *multiWalletTxListener) OnSyncCompleted() {
	handlers.OverviewHandlerLock.Lock()
	defer handlers.OverviewHandlerLock.Unlock()
	mw := app.multiWallet
	overviewHandler.Synced = mw.IsSynced()
	overviewHandler.Syncing = mw.IsSyncing()
	overviewHandler.SyncProgress = 1
	overviewHandler.UpdateBlockStatusBox(mw)
}

func (app *multiWalletTxListener) OnSyncCanceled(willRestart bool) {
	handlers.OverviewHandlerLock.Lock()
	defer handlers.OverviewHandlerLock.Unlock()
	mw := app.multiWallet
	overviewHandler.Synced = false
	overviewHandler.Syncing = false
	overviewHandler.UpdateBlockStatusBox(mw)
}

func (app *multiWalletTxListener) OnSyncEndedWithError(err error) {

}

func (app *multiWalletTxListener) Debug(debugInfo *dcrlibwallet.DebugInfo) {

}

func (app *multiWalletTxListener) OnTransaction(transaction string) {
	var currentTransaction dcrlibwallet.Transaction
	err := json.Unmarshal([]byte(transaction), &currentTransaction)
	if err != nil {
		log.Println("could read transaction to json")
		return
	}
	app.desktopNotifier(currentTransaction)

	// place all dynamic widgets here to be updated only when tabmenu is in view.
	if app.tabMenu.CurrentTabIndex() == 0 {
		transactionUpdate := handlers.TransactionUpdate{
			Transaction: currentTransaction,
		}
		overviewHandler.UpdateTransactions(app.multiWallet, transactionUpdate)
		overviewHandler.UpdateBalance(app.multiWallet)
	} else if app.tabMenu.CurrentTabIndex() == 2 {
		multipagecomponents.UpdateAccountSelectorOnNotification(sendPage.sendingAccountBoxes, sendPage.sendingSelectedAccountBalanceLabel,
			sendPage.spendableLabel, app.multiWallet, sendPage.sendingSelectedWalletID, sendPage.sendingSelectedAccountID, sendPage.Contents)

		multipagecomponents.UpdateAccountSelectorOnNotification(sendPage.selfSendingAccountBoxes, sendPage.selfSendingSelectedAccountBalanceLabel,
			nil, app.multiWallet, sendPage.selfSendingSelectedWalletID, sendPage.selfSendingSelectedAccountID, sendPage.Contents)

	} else if app.tabMenu.CurrentTabIndex() == 3 {
		multipagecomponents.UpdateAccountSelectorOnNotification(receivePage.accountBoxes, receivePage.selectedAccountBalanceLabel,
			nil, app.multiWallet, receivePage.selectedWalletID, receivePage.selectedAccountID, receivePage.Contents)
	} else if app.tabMenu.CurrentTabIndex() == 4 {

	}
}

func (app *multiWalletTxListener) OnTransactionConfirmed(walletID int, hash string, blockHeight int32) {
	// place all dynamic widgets in a function here, to be updated only when tabmenu is in view.
	if app.tabMenu.CurrentTabIndex() == 0 {
		transactionUpdate := handlers.TransactionUpdate{
			WalletId:    walletID,
			TxnHash:     hash,
		}
		overviewHandler.UpdateTransactions(app.multiWallet, transactionUpdate)
		overviewHandler.UpdateBalance(app.multiWallet)
	} else if app.tabMenu.CurrentTabIndex() == 2 {
		multipagecomponents.UpdateAccountSelectorOnNotification(sendPage.sendingAccountBoxes, sendPage.sendingSelectedAccountBalanceLabel,
			sendPage.spendableLabel, app.multiWallet, sendPage.sendingSelectedWalletID, sendPage.sendingSelectedAccountID, sendPage.Contents)

		multipagecomponents.UpdateAccountSelectorOnNotification(sendPage.selfSendingAccountBoxes, sendPage.selfSendingSelectedAccountBalanceLabel,
			nil, app.multiWallet, sendPage.selfSendingSelectedWalletID, sendPage.selfSendingSelectedAccountID, sendPage.Contents)

	} else if app.tabMenu.CurrentTabIndex() == 3 {
		multipagecomponents.UpdateAccountSelectorOnNotification(receivePage.accountBoxes, receivePage.selectedAccountBalanceLabel,
			nil, app.multiWallet, receivePage.selectedWalletID, receivePage.selectedAccountID, receivePage.Contents)

	} else if app.tabMenu.CurrentTabIndex() == 4 {

	}
}

func (app *multiWalletTxListener) OnBlockAttached(walletID int, blockHeight int32) {
	// place all dynamic widgets in a function here, to be updated only when tabmenu is in view.
	if app.tabMenu.CurrentTabIndex() == 0 {

	} else if app.tabMenu.CurrentTabIndex() == 2 {

	} else if app.tabMenu.CurrentTabIndex() == 3 {

	} else if app.tabMenu.CurrentTabIndex() == 4 {

	}
}

func (app *multiWalletTxListener) desktopNotifier(currentTransaction dcrlibwallet.Transaction) {
	amount := dcrlibwallet.AmountCoin(currentTransaction.Amount)
	// remove trailing zeros from amount
	if currentTransaction.Direction == 1 {
		var notification string

		if app.multiWallet.OpenedWalletsCount() > 1 {
			wallet := app.multiWallet.WalletWithID(currentTransaction.WalletID)
			if wallet == nil {
				return
			}

			notification = fmt.Sprintf("[%s] You have received %s DCR", wallet.Name, strconv.FormatFloat(amount, 'f', -1, 64))
		} else {

			notification = fmt.Sprintf("You have received %s DCR", strconv.FormatFloat(amount, 'f', -1, 64))
		}

		err := beeep.Notify("Decred Fyne Wallet", notification, "assets/information.png")
		if err != nil {
			log.Println("could not initiate desktop notification, reason:", err.Error())
		}
	}
}

func (app *AppInterface) walletNotificationListener() {
	var dcrListener multiWalletTxListener
	dcrListener.tabMenu = app.tabMenu
	dcrListener.multiWallet = app.MultiWallet

	err := app.MultiWallet.AddSyncProgressListener(&dcrListener, "")
	if err != nil {
		log.Fatalln("could not start progress multiWalletTxListener")
	}
}
