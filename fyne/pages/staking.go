package pages

import (
	"sort"

	"fyne.io/fyne"
	"fyne.io/fyne/canvas"
	"fyne.io/fyne/widget"

	"github.com/decred/dcrd/dcrutil"
	"github.com/raedahgroup/dcrlibwallet"

	"github.com/raedahgroup/godcr/fyne/pages/handler/multipagecomponents"
	"github.com/raedahgroup/godcr/fyne/pages/handler/stakingpagehandler"
	"github.com/raedahgroup/godcr/fyne/pages/handler/values"
	"github.com/raedahgroup/godcr/fyne/widgets"
)

type stakingPageDynamicData struct {
	accountBoxes                []*widget.Box
	selectedAccountLabel        *canvas.Text
	selectedAccountBalanceLabel *canvas.Text
	selectedWalletID            int
	selectedAccountID           int
	Contents                    *widget.Box
}

var stakingPage stakingPageDynamicData

func stakingPageContent(app *AppInterface) fyne.CanvasObject {
	openedWalletIDs := app.MultiWallet.OpenedWalletIDsRaw()
	if len(openedWalletIDs) == 0 {
		return widget.NewHBox(widgets.NewHSpacer(10), widget.NewLabelWithStyle(values.WalletsErr, fyne.TextAlignCenter, fyne.TextStyle{Bold: true}))
	}
	sort.Ints(openedWalletIDs)

	var selectedWallet = app.MultiWallet.WalletWithID(openedWalletIDs[0])
	if selectedWallet == nil {
		return widget.NewLabelWithStyle(values.LoadMultiWalletErr, fyne.TextAlignLeading, fyne.TextStyle{})
	}

	selectedWalletAccounts, err := selectedWallet.GetAccountsRaw(dcrlibwallet.DefaultRequiredConfirmations)
	if err != nil {
		return widget.NewLabel(values.AccountDetailsErr)
	}

	initStakingPageDynamicContent(openedWalletIDs, selectedWalletAccounts)

	accountSelectorObjects := multipagecomponents.AccountSelectorStruct{
		MultiWallet:              app.MultiWallet,
		WalletIDs:                openedWalletIDs,
		SendingSelectedWalletID:  &stakingPage.selectedWalletID,
		SendingSelectedAccountID: &stakingPage.selectedAccountID,

		AccountBoxes:                stakingPage.accountBoxes,
		SelectedAccountLabel:        stakingPage.selectedAccountLabel,
		SelectedAccountBalanceLabel: stakingPage.selectedAccountBalanceLabel,

		PageContents: stakingPage.Contents,
		Window:       app.Window,
	}

	initStakingPage := stakingpagehandler.StakingPageObjects{
		Accounts:            accountSelectorObjects,
		MultiWallet:         app.MultiWallet,
		StakingPageContents: stakingPage.Contents,
		Window:              app.Window,
	}

	err = initStakingPage.InitStakingPage()
	if err != nil {
		return widget.NewLabelWithStyle(values.ReceivePageLoadErr, fyne.TextAlignLeading, fyne.TextStyle{})
	}

	return widget.NewHBox(widgets.NewHSpacer(values.Padding), stakingPage.Contents, widgets.NewHSpacer(values.Padding))
}

func initStakingPageDynamicContent(openedWalletIDs []int, selectedWalletAccounts *dcrlibwallet.Accounts) {
	stakingPage = stakingPageDynamicData{}

	stakingPage.selectedWalletID = openedWalletIDs[0]
	stakingPage.accountBoxes = make([]*widget.Box, len(openedWalletIDs))

	stakingPage.selectedAccountLabel = canvas.NewText(selectedWalletAccounts.Acc[0].Name, values.DefaultTextColor)
	stakingPage.selectedAccountBalanceLabel = canvas.NewText(dcrutil.Amount(selectedWalletAccounts.Acc[0].TotalBalance).String(), values.DefaultTextColor)

	stakingPage.Contents = widget.NewVBox()
}
