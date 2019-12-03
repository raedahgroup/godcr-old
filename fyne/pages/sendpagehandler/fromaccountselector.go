package sendpagehandler

import (
	"image/color"

	"fyne.io/fyne"
	"fyne.io/fyne/canvas"
	"fyne.io/fyne/widget"

	"github.com/raedahgroup/dcrlibwallet"
)

func FromAccountSelector(onAccountChange func(), accountLabel string, receiveAccountIcon, collapseIcon fyne.Resource,
	multiWallet *dcrlibwallet.MultiWallet, walletIDs []int, sendingSelectedWalletID *int,
	accountBoxes []*widget.Box, selectedAccountLabel *widget.Label, selectedAccountBalanceLabel *widget.Label,
	selectedWalletLabel *canvas.Text, contents *widget.Box) (box *widget.Box) {

	fromLabel := canvas.NewText("From", color.RGBA{61, 88, 115, 255})
	fromLabel.TextStyle.Bold = true

	accountBox := CreateAccountSelector(onAccountChange, accountLabel, receiveAccountIcon, collapseIcon, multiWallet, walletIDs, sendingSelectedWalletID,
		accountBoxes, selectedAccountLabel, selectedAccountBalanceLabel, selectedWalletLabel, contents)

	box = widget.NewVBox(fromLabel, accountBox)

	return
}
