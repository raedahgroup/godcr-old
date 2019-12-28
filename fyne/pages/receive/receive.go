package receive

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

type ReceivePageDynamicData struct {
	AccountBoxes                []*widget.Box
	SelectedAccountLabel        *canvas.Text
	SelectedAccountBalanceLabel *canvas.Text
	SelectedWalletID            int
	SelectedAccountID           int
	Contents                    *widget.Box
}

var ReceivePage ReceivePageDynamicData

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

	initReceivePageDynamicContent(openedWalletIDs, selectedWalletAccounts)

	accountSelectorObjects := multipagecomponents.AccountSelectorStruct{
		MultiWallet:              multiWallet,
		WalletIDs:                openedWalletIDs,
		SendingSelectedWalletID:  &ReceivePage.SelectedWalletID,
		SendingSelectedAccountID: &ReceivePage.SelectedAccountID,

		AccountBoxes:                ReceivePage.AccountBoxes,
		SelectedAccountLabel:        ReceivePage.SelectedAccountLabel,
		SelectedAccountBalanceLabel: ReceivePage.SelectedAccountBalanceLabel,

		PageContents: ReceivePage.Contents,
		Window:       window,
	}

	initReceivePage := ReceivePageObjects{
		Accounts:            accountSelectorObjects,
		MultiWallet:         multiWallet,
		ReceivePageContents: ReceivePage.Contents,
		Window:              window,
	}

	err = initReceivePage.InitReceivePage()
	if err != nil {
		return widget.NewLabelWithStyle(values.ReceivePageLoadErr, fyne.TextAlignLeading, fyne.TextStyle{})
	}

	return widget.NewHBox(widgets.NewHSpacer(values.Padding), ReceivePage.Contents, widgets.NewHSpacer(values.Padding))
}

func initReceivePageDynamicContent(openedWalletIDs []int, selectedWalletAccounts *dcrlibwallet.Accounts) {
	ReceivePage = ReceivePageDynamicData{}

	ReceivePage.SelectedWalletID = openedWalletIDs[0]
	ReceivePage.AccountBoxes = make([]*widget.Box, len(openedWalletIDs))

	ReceivePage.SelectedAccountLabel = canvas.NewText(selectedWalletAccounts.Acc[0].Name, values.DefaultTextColor)
	ReceivePage.SelectedAccountBalanceLabel = canvas.NewText(dcrutil.Amount(selectedWalletAccounts.Acc[0].TotalBalance).String(), values.DefaultTextColor)

	ReceivePage.Contents = widget.NewVBox()
}
