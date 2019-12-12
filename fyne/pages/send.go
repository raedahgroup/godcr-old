package pages

import (
	"image/color"
	"log"
	"sort"

	"fyne.io/fyne"
	"fyne.io/fyne/canvas"
	"fyne.io/fyne/widget"

	"github.com/decred/dcrd/dcrutil"
	"github.com/raedahgroup/dcrlibwallet"

	"github.com/raedahgroup/godcr/fyne/pages/constantvalues"
	"github.com/raedahgroup/godcr/fyne/pages/multipagecomponents"
	"github.com/raedahgroup/godcr/fyne/pages/sendpagehandler"
	"github.com/raedahgroup/godcr/fyne/widgets"
)

type sendPageDynamicData struct {
	// houses all clickable box
	sendingAccountBoxes     []*widget.Box
	selfSendingAccountBoxes []*widget.Box

	spendableLabel *canvas.Text

	selfSendingSelectedAccountLabel        *widget.Label
	selfSendingSelectedAccountBalanceLabel *widget.Label
	selfSendingSelectedWalletID            int

	sendingSelectedAccountLabel        *widget.Label
	sendingSelectedAccountBalanceLabel *widget.Label
	sendingSelectedWalletID            int

	Contents *widget.Box
}

var sendPage sendPageDynamicData

func sendPageContent(multiWallet *dcrlibwallet.MultiWallet, window fyne.Window) fyne.CanvasObject {
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

	initSendPageDynamicContent(openedWalletIDs, selectedWalletAccounts)

	sendingFromAccountSelectorObjects := multipagecomponents.AccountSelectorStruct{
		MultiWallet:             multiWallet,
		WalletIDs:               openedWalletIDs,
		SendingSelectedWalletID: &sendPage.sendingSelectedWalletID,

		AccountBoxes:                sendPage.sendingAccountBoxes,
		SelectedAccountLabel:        sendPage.sendingSelectedAccountLabel,
		SelectedAccountBalanceLabel: sendPage.sendingSelectedAccountBalanceLabel,

		PageContents: sendPage.Contents,
		Window:       window,
	}

	sendingToAccountSelectorObjects := multipagecomponents.AccountSelectorStruct{
		MultiWallet:             multiWallet,
		WalletIDs:               openedWalletIDs,
		SendingSelectedWalletID: &sendPage.selfSendingSelectedWalletID,

		AccountBoxes:                sendPage.selfSendingAccountBoxes,
		SelectedAccountLabel:        sendPage.selfSendingSelectedAccountLabel,
		SelectedAccountBalanceLabel: sendPage.selfSendingSelectedAccountBalanceLabel,

		PageContents: sendPage.Contents,
		Window:       window,
	}

	initSendPage := sendpagehandler.SendPageObjects{
		MultiWallet:      multiWallet,
		SpendableLabel:   sendPage.spendableLabel,
		SendPageContents: sendPage.Contents,
		Sending:          sendingFromAccountSelectorObjects,
		SelfSending:      sendingToAccountSelectorObjects,
		Window:           window,
	}
	err = initSendPage.InitAllSendPageComponents()
	if err != nil {
		return widget.NewLabelWithStyle("Unable to load Send Page components, "+err.Error(), fyne.TextAlignLeading, fyne.TextStyle{})
	}

	return widget.NewHBox(widgets.NewHSpacer(30), initSendPage.SendPageContents)
}

func initSendPageDynamicContent(openedWalletIDs []int, selectedWalletAccounts *dcrlibwallet.Accounts) {
	sendPage.sendingSelectedWalletID = openedWalletIDs[0]
	sendPage.selfSendingSelectedWalletID = openedWalletIDs[0]

	sendPage.selfSendingAccountBoxes = make([]*widget.Box, len(openedWalletIDs))
	sendPage.sendingAccountBoxes = make([]*widget.Box, len(openedWalletIDs))

	sendPage.spendableLabel = canvas.NewText(constantvalues.SpendableAmountLabel+dcrutil.Amount(selectedWalletAccounts.Acc[0].Balance.Spendable).String(), color.Black)
	sendPage.spendableLabel.TextSize = 12

	sendPage.sendingSelectedAccountLabel = widget.NewLabel(selectedWalletAccounts.Acc[0].Name)
	sendPage.sendingSelectedAccountBalanceLabel = widget.NewLabel(dcrutil.Amount(selectedWalletAccounts.Acc[0].TotalBalance).String())

	sendPage.selfSendingSelectedAccountLabel = widget.NewLabel(selectedWalletAccounts.Acc[0].Name)
	sendPage.selfSendingSelectedAccountBalanceLabel = widget.NewLabel(dcrutil.Amount(selectedWalletAccounts.Acc[0].TotalBalance).String())

	sendPage.Contents = widget.NewVBox()
}
