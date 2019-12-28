package pages

import (
	"encoding/json"
	"fmt"
	"fyne.io/fyne/widget"
	"github.com/raedahgroup/godcr/fyne/pages/receive"
	"github.com/raedahgroup/godcr/fyne/pages/send"
	"log"
	"math"
	"strconv"

	"github.com/gen2brain/beeep"
	"github.com/raedahgroup/dcrlibwallet"
	"github.com/raedahgroup/godcr/fyne/pages/multipagecomponents"
	"github.com/raedahgroup/godcr/fyne/pages/overview"
)

type multiWalletTxListener struct {
	tabMenu     *widget.TabContainer
	multiWallet *dcrlibwallet.MultiWallet
	handlers 	*pageHandlers
}

func (app *multiWalletTxListener) OnSyncStarted() {
	mw := app.multiWallet
	go app.handlers.overviewHandler.UpdateBlockStatusBox(mw)
}

func (app *multiWalletTxListener) OnPeerConnectedOrDisconnected(numberOfConnectedPeers int32) {
	app.handlers.overviewHandler.UpdateConnectedPeers(numberOfConnectedPeers, true)
}

func (app *multiWalletTxListener) OnHeadersFetchProgress(headersFetchProgress *dcrlibwallet.HeadersFetchProgressReport) {
	app.handlers.overviewHandler.SyncProgress = float64(headersFetchProgress.FetchedHeadersCount) / float64(headersFetchProgress.TotalHeadersToFetch)
	app.handlers.overviewHandler.SyncProgress = math.Round(app.handlers.overviewHandler.SyncProgress*100) / 100
	if headersFetchProgress.HeadersFetchProgress == 100 && app.handlers.overviewHandler.Steps == 0 {
		app.handlers.overviewHandler.StepsChannel <- headersFetchProgress.HeadersFetchProgress
		app.handlers.overviewHandler.UpdateSyncSteps(true)
	}
	app.handlers.overviewHandler.UpdateBlockHeadersSync(headersFetchProgress.HeadersFetchProgress, true)
	go app.handlers.overviewHandler.UpdateWalletsSyncBox(app.multiWallet)
}

func (app *multiWalletTxListener) OnAddressDiscoveryProgress(addressDiscoveryProgress *dcrlibwallet.AddressDiscoveryProgressReport) {
	if addressDiscoveryProgress.AddressDiscoveryProgress == 100 {
		app.handlers.overviewHandler.StepsChannel <- addressDiscoveryProgress.AddressDiscoveryProgress
		app.handlers.overviewHandler.UpdateSyncSteps(true)
	}
}

func (app *multiWalletTxListener) OnHeadersRescanProgress(headersRescanProgress *dcrlibwallet.HeadersRescanProgressReport) {
	if headersRescanProgress.RescanProgress == 100 {
		app.handlers.overviewHandler.StepsChannel <- headersRescanProgress.RescanProgress
		app.handlers.overviewHandler.UpdateSyncSteps(true)
	}
}

func (app *multiWalletTxListener) OnSyncCompleted() {
	app.handlers.overviewHandler.SyncProgress = 1
	go app.handlers.overviewHandler.UpdateBlockStatusBox(app.multiWallet)
	go app.handlers.overviewHandler.UpdateBalance(app.multiWallet)
	go app.handlers.overviewHandler.UpdateTransactions(app.multiWallet, overview.TransactionUpdate{})
}

func (app *multiWalletTxListener) OnSyncCanceled(willRestart bool) {
	go app.handlers.overviewHandler.UpdateBlockStatusBox(app.multiWallet)
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
		transactionUpdate := overview.TransactionUpdate{
			Transaction: currentTransaction,
		}
		app.handlers.overviewHandler.UpdateTransactions(app.multiWallet, transactionUpdate)
		app.handlers.overviewHandler.UpdateBalance(app.multiWallet)
	} else if app.tabMenu.CurrentTabIndex() == 2 {
		multipagecomponents.UpdateAccountSelectorOnNotification(send.SendPage.SendingAccountBoxes, send.SendPage.SendingSelectedAccountBalanceLabel,
			send.SendPage.SpendableLabel, app.multiWallet, send.SendPage.SendingSelectedWalletID, send.SendPage.SendingSelectedAccountID, send.SendPage.Contents)

		multipagecomponents.UpdateAccountSelectorOnNotification(send.SendPage.SelfSendingAccountBoxes, send.SendPage.SelfSendingSelectedAccountBalanceLabel,
			nil, app.multiWallet, send.SendPage.SelfSendingSelectedWalletID, send.SendPage.SelfSendingSelectedAccountID, send.SendPage.Contents)

	} else if app.tabMenu.CurrentTabIndex() == 3 {
		multipagecomponents.UpdateAccountSelectorOnNotification(receive.ReceivePage.AccountBoxes, receive.ReceivePage.SelectedAccountBalanceLabel,
			nil, app.multiWallet, receive.ReceivePage.SelectedWalletID, receive.ReceivePage.SelectedAccountID, receive.ReceivePage.Contents)
	} else if app.tabMenu.CurrentTabIndex() == 4 {

	}
}

func (app *multiWalletTxListener) OnTransactionConfirmed(walletID int, hash string, blockHeight int32) {
	// place all dynamic widgets in a function here, to be updated only when tabmenu is in view.
	if app.tabMenu.CurrentTabIndex() == 0 {
		transactionUpdate := overview.TransactionUpdate{
			WalletId: walletID,
			TxnHash:  hash,
		}
		app.handlers.overviewHandler.UpdateTransactions(app.multiWallet, transactionUpdate)
		app.handlers.overviewHandler.UpdateBalance(app.multiWallet)
	} else if app.tabMenu.CurrentTabIndex() == 2 {
		multipagecomponents.UpdateAccountSelectorOnNotification(send.SendPage.SendingAccountBoxes, send.SendPage.SendingSelectedAccountBalanceLabel,
			send.SendPage.SpendableLabel, app.multiWallet, send.SendPage.SendingSelectedWalletID, send.SendPage.SendingSelectedAccountID, send.SendPage.Contents)

		multipagecomponents.UpdateAccountSelectorOnNotification(send.SendPage.SelfSendingAccountBoxes, send.SendPage.SelfSendingSelectedAccountBalanceLabel,
			nil, app.multiWallet, send.SendPage.SelfSendingSelectedWalletID, send.SendPage.SelfSendingSelectedAccountID, send.SendPage.Contents)

	} else if app.tabMenu.CurrentTabIndex() == 3 {
		multipagecomponents.UpdateAccountSelectorOnNotification(receive.ReceivePage.AccountBoxes, receive.ReceivePage.SelectedAccountBalanceLabel,
			nil, app.multiWallet, receive.ReceivePage.SelectedWalletID, receive.ReceivePage.SelectedAccountID, receive.ReceivePage.Contents)

	} else if app.tabMenu.CurrentTabIndex() == 4 {

	}
}

func (app *multiWalletTxListener) OnBlockAttached(walletID int, blockHeight int32) {
	// place all dynamic widgets in a function here, to be updated only when tabmenu is in view.
	if app.tabMenu.CurrentTabIndex() == 0 {
		app.handlers.overviewHandler.UpdateBlockStatusBox(app.multiWallet)
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
	dcrListener.handlers = app.handlers

	err := app.MultiWallet.AddSyncProgressListener(&dcrListener, "")
	if err != nil {
		log.Fatalln("could not start progress multiWalletTxListener")
	}
}
