package pages

import (
	"sort"
	"strings"

	"fyne.io/fyne"
	"fyne.io/fyne/canvas"
	"fyne.io/fyne/widget"

	"github.com/decred/dcrd/dcrutil"
	"github.com/raedahgroup/dcrlibwallet"

	"github.com/raedahgroup/godcr/fyne/pages/handler/multipagecomponents"
	"github.com/raedahgroup/godcr/fyne/pages/handler/sendpagehandler"
	"github.com/raedahgroup/godcr/fyne/pages/handler/values"
	"github.com/raedahgroup/godcr/fyne/widgets"
)

type sendPageDynamicData struct {
	// houses all clickable box
	sendingAccountBoxes     []*widget.Box
	selfSendingAccountBoxes []*widget.Box

	spendableLabel *canvas.Text

	selfSendingSelectedAccountLabel        *canvas.Text
	selfSendingSelectedAccountBalanceLabel *canvas.Text
	selfSendingSelectedWalletID            int
	selfSendingSelectedAccountID           int

	sendingSelectedAccountLabel        *canvas.Text
	sendingSelectedAccountBalanceLabel *canvas.Text
	sendingSelectedWalletID            int
	sendingSelectedAccountID           int

	Contents *widget.Box
}

var sendPage sendPageDynamicData

func sendPageContent(multiWallet *dcrlibwallet.MultiWallet, window fyne.Window) fyne.CanvasObject {
	openedWalletIDs := multiWallet.OpenedWalletIDsRaw()
	if len(openedWalletIDs) == 0 {
		return widget.NewHBox(widgets.NewHSpacer(10), widget.NewLabelWithStyle(values.WalletsErr, fyne.TextAlignCenter, fyne.TextStyle{Bold: true}))
	}
	sort.Ints(openedWalletIDs)

	var selectedWallet = multiWallet.WalletWithID(openedWalletIDs[0])
	if selectedWallet == nil {
		return widget.NewLabelWithStyle(values.LoadMultiWalletErr, fyne.TextAlignLeading, fyne.TextStyle{})
	}

	selectedWalletAccounts, err := selectedWallet.GetAccountsRaw(dcrlibwallet.DefaultRequiredConfirmations)
	if err != nil {
		return widget.NewLabel(values.AccountDetailsErr)
	}

	initSendPageDynamicContent(openedWalletIDs, selectedWalletAccounts)

	sendingFromAccountSelectorObjects := multipagecomponents.AccountSelectorStruct{
		MultiWallet:              multiWallet,
		WalletIDs:                openedWalletIDs,
		SendingSelectedWalletID:  &sendPage.sendingSelectedWalletID,
		SendingSelectedAccountID: &sendPage.sendingSelectedAccountID,

		AccountBoxes:                sendPage.sendingAccountBoxes,
		SelectedAccountLabel:        sendPage.sendingSelectedAccountLabel,
		SelectedAccountBalanceLabel: sendPage.sendingSelectedAccountBalanceLabel,

		PageContents: sendPage.Contents,
		Window:       window,
	}

	sendingToAccountSelectorObjects := multipagecomponents.AccountSelectorStruct{
		MultiWallet:              multiWallet,
		WalletIDs:                openedWalletIDs,
		SendingSelectedWalletID:  &sendPage.selfSendingSelectedWalletID,
		SendingSelectedAccountID: &sendPage.selfSendingSelectedAccountID,

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
		return widget.NewLabelWithStyle(values.SendPageLoadErr, fyne.TextAlignLeading, fyne.TextStyle{})
	}

	return widget.NewHBox(widgets.NewHSpacer(values.Padding), initSendPage.SendPageContents, widgets.NewHSpacer(values.Padding))
}

func initSendPageDynamicContent(openedWalletIDs []int, selectedWalletAccounts *dcrlibwallet.Accounts) {
	sendPage = sendPageDynamicData{}
	firstWalletWalletID := openedWalletIDs[0]
	defaultAccount := selectedWalletAccounts.Acc[0]

	sendPage.sendingSelectedWalletID = firstWalletWalletID
	sendPage.selfSendingSelectedWalletID = firstWalletWalletID

	sendPage.selfSendingAccountBoxes = make([]*widget.Box, len(openedWalletIDs))
	sendPage.sendingAccountBoxes = make([]*widget.Box, len(openedWalletIDs))

	sendPage.spendableLabel = canvas.NewText(values.SpendableAmountLabel+dcrutil.Amount(defaultAccount.Balance.Spendable).String(), values.DarkerBlueGrayTextColor)
	sendPage.spendableLabel.TextSize = values.TextSize12

	sendPage.sendingSelectedAccountLabel = canvas.NewText(strings.Title(defaultAccount.Name), values.DefaultTextColor)
	sendPage.sendingSelectedAccountBalanceLabel = canvas.NewText(dcrutil.Amount(defaultAccount.TotalBalance).String(), values.DefaultTextColor)

	sendPage.selfSendingSelectedAccountLabel = canvas.NewText(strings.Title(defaultAccount.Name), values.DefaultTextColor)
	sendPage.selfSendingSelectedAccountBalanceLabel = canvas.NewText(dcrutil.Amount(defaultAccount.TotalBalance).String(), values.DefaultTextColor)

	sendPage.Contents = widget.NewVBox()
}
