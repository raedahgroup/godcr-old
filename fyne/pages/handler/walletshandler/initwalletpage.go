package walletshandler

import (
	"fyne.io/fyne/canvas"
	"fyne.io/fyne/widget"

	"github.com/raedahgroup/dcrlibwallet"

	"github.com/raedahgroup/godcr/fyne/pages/handler/values"
	"github.com/raedahgroup/godcr/fyne/widgets"
)

type WalletPageObject struct {
	walletSelectorBox *widget.Box
	//	WalletsAccountAmountText  [][]*canvas.Text
	//	WalletSpendableAmountText [][]*canvas.Text
	WalletTotalAmountLabel []*canvas.Text

	OpenedWallets   []int
	walletIDToIndex map[int]int // relates wallet ID to iterator where key is iterator and value wallet ID

	WalletPageContents *widget.Box
	MultiWallet        *dcrlibwallet.MultiWallet
}

func (walletPage *WalletPageObject) InitWalletPage() error {
	walletPage.WalletPageContents.Append(widgets.NewVSpacer(values.Padding))

	err := walletPage.initBaseWidgets()
	if err != nil {
		return err
	}

	walletPage.WalletPageContents.Append(widgets.NewVSpacer(14))

	err = walletPage.accountSelector()
	if err != nil {
		return err
	}

	return nil
}
