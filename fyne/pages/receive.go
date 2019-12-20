package pages

import (
	"sort"

	"fyne.io/fyne"
	"fyne.io/fyne/canvas"
	"fyne.io/fyne/widget"

	"github.com/decred/dcrd/dcrutil"
	"github.com/raedahgroup/dcrlibwallet"

	"github.com/raedahgroup/godcr/fyne/pages/handler/multipagecomponents"
	"github.com/raedahgroup/godcr/fyne/pages/handler/receivepagehandler"
	"github.com/raedahgroup/godcr/fyne/pages/handler/values"
	"github.com/raedahgroup/godcr/fyne/widgets"
)

type receivePageDynamicData struct {
	accountBoxes                []*widget.Box
	selectedAccountLabel        *canvas.Text
	selectedAccountBalanceLabel *canvas.Text
	selectedWalletID            int
	selectedAccountID           int
	Contents                    *widget.Box
}

var receivePage receivePageDynamicData

func receivePageContent(multiWallet *dcrlibwallet.MultiWallet, window fyne.Window) fyne.CanvasObject {
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

	initReceivePageDynamicContent(openedWalletIDs, selectedWalletAccounts)

	accountSelectorObjects := multipagecomponents.AccountSelectorStruct{
		MultiWallet:              multiWallet,
		WalletIDs:                openedWalletIDs,
		SendingSelectedWalletID:  &receivePage.selectedWalletID,
		SendingSelectedAccountID: &receivePage.selectedAccountID,

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
		return widget.NewLabelWithStyle(values.ReceivePageLoadErr, fyne.TextAlignLeading, fyne.TextStyle{})
	}

	return widget.NewHBox(widgets.NewHSpacer(values.Padding), receivePage.Contents, widgets.NewHSpacer(values.Padding))
}

func initReceivePageDynamicContent(openedWalletIDs []int, selectedWalletAccounts *dcrlibwallet.Accounts) {
	receivePage = receivePageDynamicData{}

	receivePage.selectedWalletID = openedWalletIDs[0]
	receivePage.accountBoxes = make([]*widget.Box, len(openedWalletIDs))

	receivePage.selectedAccountLabel = canvas.NewText(selectedWalletAccounts.Acc[0].Name, values.DefaultTextColor)
	receivePage.selectedAccountBalanceLabel = canvas.NewText(dcrutil.Amount(selectedWalletAccounts.Acc[0].TotalBalance).String(), values.DefaultTextColor)

	receivePage.Contents = widget.NewVBox()
}
