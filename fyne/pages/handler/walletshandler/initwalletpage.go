package walletshandler

import (
	"fyne.io/fyne/canvas"
	"fyne.io/fyne/widget"
	"github.com/raedahgroup/dcrlibwallet"
)

type WalletPageObject struct {
	WalletsAccountAmountText  [][]*canvas.Text
	WalletSpendableAmountText [][]*canvas.Text
	WalletTotalAmountText     []*canvas.Text

	OpenedWallets []int

	WalletPageContents *widget.Box
	MultiWallet        *dcrlibwallet.MultiWallet
}
