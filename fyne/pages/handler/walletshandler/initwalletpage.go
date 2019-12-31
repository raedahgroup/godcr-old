package walletshandler

import (
	"fyne.io/fyne"
	"fyne.io/fyne/canvas"
	"fyne.io/fyne/widget"

	"github.com/raedahgroup/dcrlibwallet"

	"github.com/raedahgroup/godcr/fyne/assets"
	"github.com/raedahgroup/godcr/fyne/pages/handler/values"
	"github.com/raedahgroup/godcr/fyne/widgets"
)

type WalletPageObject struct {
	icons map[string]*fyne.StaticResource

	walletSelectorBox *widget.Box
	//	WalletsAccountAmountText  [][]*canvas.Text
	//	WalletSpendableAmountText [][]*canvas.Text
	WalletTotalAmountLabel []*canvas.Text

	OpenedWallets   []int
	walletIDToIndex map[int]int // relates wallet ID to iterator where key is iterator and value wallet ID

	WalletPageContents *widget.Box
	MultiWallet        *dcrlibwallet.MultiWallet

	Window fyne.Window
}

func (walletPage *WalletPageObject) InitWalletPage() error {
	var err error
	walletPage.icons, err = assets.GetIcons(assets.AddWallet, assets.Expand, assets.CollapseIcon, assets.WalletIcon,
		assets.WalletAccount, assets.ImportedAccount, assets.MoreIcon, assets.Edit)

	walletPage.WalletPageContents.Append(widgets.NewVSpacer(values.Padding))

	err = walletPage.initBaseWidgets()
	if err != nil {
		return err
	}

	walletPage.WalletPageContents.Append(widgets.NewVSpacer(values.Padding))

	err = walletPage.accountSelector()
	if err != nil {
		return err
	}

	return nil
}
