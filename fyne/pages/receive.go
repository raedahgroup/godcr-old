package pages

import (
	"log"
	"sort"

	"fyne.io/fyne"
	"fyne.io/fyne/widget"

	"github.com/decred/dcrd/dcrutil"
	"github.com/raedahgroup/dcrlibwallet"

	"github.com/raedahgroup/godcr/fyne/pages/multipagecomponents"
	"github.com/raedahgroup/godcr/fyne/pages/receivepagehandler"
	"github.com/raedahgroup/godcr/fyne/widgets"
)

type receivePageDynamicData struct {
	accountBoxes                []*widget.Box
	selectedAccountLabel        *widget.Label
	selectedAccountBalanceLabel *widget.Label
	selectedWalletID            int
	Contents                    *widget.Box
}

var receivePage receivePageDynamicData

func receivePageContent(multiWallet *dcrlibwallet.MultiWallet, window fyne.Window) fyne.CanvasObject {
	openedWalletIDs := multiWallet.OpenedWalletIDsRaw()
	if len(openedWalletIDs) == 0 {
		return widget.NewHBox(widgets.NewHSpacer(10), widget.NewLabelWithStyle("Could not retrieve wallets", fyne.TextAlignCenter, fyne.TextStyle{Bold: true}))
	}
	sort.Ints(openedWalletIDs)

	var selectedWallet = multiWallet.WalletWithID(openedWalletIDs[0])
	if selectedWallet == nil {
		return widget.NewLabelWithStyle("Unable to load MultiWallet", fyne.TextAlignLeading, fyne.TextStyle{})
	}

	selectedWalletAccounts, err := selectedWallet.GetAccountsRaw(dcrlibwallet.DefaultRequiredConfirmations)
	if err != nil {
		log.Println("Error while getting accounts for wallet", err.Error())
		return widget.NewLabel("Error while getting accounts for wallet")
	}

	initReceivePageDynamicContent(openedWalletIDs, selectedWalletAccounts)

	accountSelectorObjects := multipagecomponents.AccountSelectorStruct{
		MultiWallet:             multiWallet,
		WalletIDs:               openedWalletIDs,
		SendingSelectedWalletID: &receivePage.selectedWalletID,

		AccountBoxes:                receivePage.accountBoxes,
		SelectedAccountLabel:        receivePage.selectedAccountLabel,
		SelectedAccountBalanceLabel: receivePage.selectedAccountBalanceLabel,

		PageContents: receivePage.Contents,
		Window:       window,
	}

	initReceivePage := receivepagehandler.ReceivePageObjects{
		Accounts:            accountSelectorObjects,
		MultiWallet:         multiWallet,
		ReceivePageContents: receivePage.Contents,
		Window:              window,
	}

	err = initReceivePage.InitReceivePage()
	if err != nil {
		return widget.NewLabelWithStyle("Unable to load Receive Page components, "+err.Error(), fyne.TextAlignLeading, fyne.TextStyle{})
	}

	return widget.NewHBox(widgets.NewHSpacer(20), receivePage.Contents)
}

func initReceivePageDynamicContent(openedWalletIDs []int, selectedWalletAccounts *dcrlibwallet.Accounts) {
	receivePage.selectedWalletID = openedWalletIDs[0]
	receivePage.accountBoxes = make([]*widget.Box, len(openedWalletIDs))

	receivePage.selectedAccountLabel = widget.NewLabel(selectedWalletAccounts.Acc[0].Name)
	receivePage.selectedAccountBalanceLabel = widget.NewLabel(dcrutil.Amount(selectedWalletAccounts.Acc[0].TotalBalance).String())

	receivePage.Contents = widget.NewVBox()
}
