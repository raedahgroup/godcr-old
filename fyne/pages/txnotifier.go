package pages

import (
	"encoding/json"
	"fmt"
	"log"
	"strconv"

	"fyne.io/fyne/widget"

	"github.com/gen2brain/beeep"
	"github.com/raedahgroup/dcrlibwallet"
)

type multiWalletTxListener struct {
	tabMenu     *widget.TabContainer
	multiWallet *dcrlibwallet.MultiWallet
}

func (app *multiWalletTxListener) OnSyncStarted() {

}

func (app *multiWalletTxListener) OnPeerConnectedOrDisconnected(numberOfConnectedPeers int32) {

}

func (app *multiWalletTxListener) OnHeadersFetchProgress(headersFetchProgress *dcrlibwallet.HeadersFetchProgressReport) {

}

func (app *multiWalletTxListener) OnAddressDiscoveryProgress(addressDiscoveryProgress *dcrlibwallet.AddressDiscoveryProgressReport) {

}

func (app *multiWalletTxListener) OnHeadersRescanProgress(headersRescanProgress *dcrlibwallet.HeadersRescanProgressReport) {

}

func (app *multiWalletTxListener) OnSyncCompleted() {

}

func (app *multiWalletTxListener) OnSyncCanceled(willRestart bool) {

}

func (app *multiWalletTxListener) OnSyncEndedWithError(err error) {

}

func (app *multiWalletTxListener) Debug(debugInfo *dcrlibwallet.DebugInfo) {

}

func (app *multiWalletTxListener) OnTransaction(transaction string) {
	var currentTransaction map[string]interface{}
	err := json.Unmarshal([]byte(transaction), &currentTransaction)
	if err != nil {
		log.Println("could read transaction to json")
		return
	}

	var amount, walletID float64
	var ok bool
	var walletName string

	if walletID, ok = currentTransaction["walletID"].(float64); ok {
		wallet := app.multiWallet.WalletWithID(int(walletID))
		if wallet == nil {
			return
		}
		walletName = wallet.Name

	} else {
		return
	}

	if amount, ok = currentTransaction["amount"].(float64); ok {
		err = beeep.Notify("Decred Fyne Wallet", fmt.Sprintf("You have received %s DCR in wallet %s", strconv.FormatFloat(amount/100000000, 'f', -1, 64), walletName), "assets/information.png")
		if err != nil {
			log.Println("could not start desktop notification")
		}

	} else {
		return
	}

	// place all dynamic widgets here to be updated only when tabmenu is in view.
	if app.tabMenu.CurrentTabIndex() == 0 {

	} else if app.tabMenu.CurrentTabIndex() == 2 {

	} else if app.tabMenu.CurrentTabIndex() == 3 {

	} else if app.tabMenu.CurrentTabIndex() == 2 {

	}
}

func (app *multiWalletTxListener) OnTransactionConfirmed(walletID int, hash string, blockHeight int32) {
	// place all dynamic widgets in a function here, to be updated only when tabmenu is in view.
	if app.tabMenu.CurrentTabIndex() == 0 {

	} else if app.tabMenu.CurrentTabIndex() == 2 {

	} else if app.tabMenu.CurrentTabIndex() == 3 {

	} else if app.tabMenu.CurrentTabIndex() == 2 {

	}
}

func (app *multiWalletTxListener) OnBlockAttached(walletID int, blockHeight int32) {
	// place all dynamic widgets in a function here, to be updated only when tabmenu is in view.
	if app.tabMenu.CurrentTabIndex() == 0 {

	} else if app.tabMenu.CurrentTabIndex() == 2 {

	} else if app.tabMenu.CurrentTabIndex() == 3 {

	} else if app.tabMenu.CurrentTabIndex() == 2 {

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
