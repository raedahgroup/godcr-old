package send

import (
	"sort"

	"fyne.io/fyne"
	"fyne.io/fyne/canvas"
	"fyne.io/fyne/widget"

	"github.com/decred/dcrd/dcrutil"
	"github.com/raedahgroup/dcrlibwallet"

	"github.com/raedahgroup/godcr/fyne/pages/multipagecomponents"
	"github.com/raedahgroup/godcr/fyne/values"
	"github.com/raedahgroup/godcr/fyne/widgets"
)

type sendPageDynamicData struct {
	// houses all clickable box
	SendingAccountBoxes     []*widget.Box
	SelfSendingAccountBoxes []*widget.Box

	SpendableLabel *canvas.Text

	SelfSendingSelectedAccountLabel        *canvas.Text
	SelfSendingSelectedAccountBalanceLabel *canvas.Text
	SelfSendingSelectedWalletID            int
	SelfSendingSelectedAccountID           int

	SendingSelectedAccountLabel        *canvas.Text
	SendingSelectedAccountBalanceLabel *canvas.Text
	SendingSelectedWalletID            int
	SendingSelectedAccountID           int

	Contents *widget.Box
}

var SendPage sendPageDynamicData

func PageContent(multiWallet *dcrlibwallet.MultiWallet, window fyne.Window) fyne.CanvasObject {
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
		SendingSelectedWalletID:  &SendPage.SendingSelectedWalletID,
		SendingSelectedAccountID: &SendPage.SendingSelectedAccountID,

		AccountBoxes:                SendPage.SendingAccountBoxes,
		SelectedAccountLabel:        SendPage.SendingSelectedAccountLabel,
		SelectedAccountBalanceLabel: SendPage.SendingSelectedAccountBalanceLabel,

		PageContents: SendPage.Contents,
		Window:       window,
	}

	sendingToAccountSelectorObjects := multipagecomponents.AccountSelectorStruct{
		MultiWallet:              multiWallet,
		WalletIDs:                openedWalletIDs,
		SendingSelectedWalletID:  &SendPage.SelfSendingSelectedWalletID,
		SendingSelectedAccountID: &SendPage.SelfSendingSelectedAccountID,

		AccountBoxes:                SendPage.SelfSendingAccountBoxes,
		SelectedAccountLabel:        SendPage.SelfSendingSelectedAccountLabel,
		SelectedAccountBalanceLabel: SendPage.SelfSendingSelectedAccountBalanceLabel,

		PageContents: SendPage.Contents,
		Window:       window,
	}

	initSendPage := SendPageObjects{
		MultiWallet:      multiWallet,
		SpendableLabel:   SendPage.SpendableLabel,
		SendPageContents: SendPage.Contents,
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
	SendPage = sendPageDynamicData{}
	firstWalletWalletID := openedWalletIDs[0]
	defaultAccount := selectedWalletAccounts.Acc[0]

	SendPage.SendingSelectedWalletID = firstWalletWalletID
	SendPage.SelfSendingSelectedWalletID = firstWalletWalletID

	SendPage.SelfSendingAccountBoxes = make([]*widget.Box, len(openedWalletIDs))
	SendPage.SendingAccountBoxes = make([]*widget.Box, len(openedWalletIDs))

	SendPage.SpendableLabel = canvas.NewText(values.SpendableAmountLabel+dcrutil.Amount(defaultAccount.Balance.Spendable).String(), values.DarkerBlueGrayTextColor)
	SendPage.SpendableLabel.TextSize = values.TextSize12

	SendPage.SendingSelectedAccountLabel = canvas.NewText(defaultAccount.Name, values.DefaultTextColor)
	SendPage.SendingSelectedAccountBalanceLabel = canvas.NewText(dcrutil.Amount(defaultAccount.TotalBalance).String(), values.DefaultTextColor)

	SendPage.SelfSendingSelectedAccountLabel = canvas.NewText(defaultAccount.Name, values.DefaultTextColor)
	SendPage.SelfSendingSelectedAccountBalanceLabel = canvas.NewText(dcrutil.Amount(defaultAccount.TotalBalance).String(), values.DefaultTextColor)

	SendPage.Contents = widget.NewVBox()
}
