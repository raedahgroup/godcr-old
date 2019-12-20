package walletshandler

import (
	"errors"
	"image/color"
	"sort"

	"fyne.io/fyne"
	"fyne.io/fyne/canvas"
	"fyne.io/fyne/widget"

	"github.com/raedahgroup/dcrlibwallet"

	"github.com/raedahgroup/godcr/fyne/widgets"
)

func (walletPage *WalletPageObject) accountSelector() error {
	openedWalletIDs := walletPage.MultiWallet.OpenedWalletIDsRaw()
	if len(openedWalletIDs) == 0 {
		return errors.New("Not wallet found")
	}
	sort.Ints(openedWalletIDs)

	return nil
}

func (walletPage *WalletPageObject) getAccountsInWallet(expandIcon, accountIcon, importedAccountIcon,
	addAccountIcon fyne.Resource, selectedWalletID int) {

	selectedWallet := walletPage.MultiWallet.WalletWithID(selectedWalletID)
	accts, err := selectedWallet.GetAccountsRaw(dcrlibwallet.DefaultRequiredConfirmations)
	if err != nil {
		return
	}

	var totalBalance int64
	for _, acc := range accts.Acc {
		totalBalance += acc.TotalBalance
	}

	// walletPage.WalletTotalAmountText[selectedWalletID].Text=selected
	accountBox := widget.NewHBox(widget.NewIcon(expandIcon), widgets.NewHSpacer(4),
		widget.NewIcon(accountIcon), widgets.NewHSpacer(12),
		canvas.NewText(selectedWallet.Name, color.RGBA{9, 20, 64, 255}), widgets.NewHSpacer(50))

}
