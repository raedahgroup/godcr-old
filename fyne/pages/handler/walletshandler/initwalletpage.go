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
	walletAccountsBox []*widgets.Box
	// changing icons when other wallet propertie box's are collapsed,
	// we need to define all in a struct for easy accessibility
	walletExpandCollapseIcon []*widget.Icon

	WalletsAccountAmountText  [][]*fyne.Container
	WalletSpendableAmountText [][]*canvas.Text
	WalletTotalAmountLabel    []*canvas.Text

	OpenedWallets   []int
	WalletIDToIndex map[int]int // relates wallet ID to iterator where key is iterator and value wallet ID

	WalletPageContents *widget.Box
	MultiWallet        *dcrlibwallet.MultiWallet

	TabMenu *widget.TabContainer
	Window  fyne.Window
}

func (walletPage *WalletPageObject) InitWalletPage() error {
	var err error
	walletPage.icons, err = assets.GetIcons(assets.AddWallet, assets.Expand, assets.CollapseIcon, assets.WalletIcon,
		assets.WalletAccount, assets.ImportedAccount, assets.MoreIcon, assets.Edit)

	walletPage.walletAccountsBox = make([]*widgets.Box, len(walletPage.OpenedWallets))
	walletPage.walletExpandCollapseIcon = make([]*widget.Icon, len(walletPage.OpenedWallets))

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
