package walletshandler

import (
	"time"

	"fyne.io/fyne"
	"fyne.io/fyne/canvas"
	"fyne.io/fyne/layout"
	"fyne.io/fyne/widget"

	"github.com/raedahgroup/dcrlibwallet"

	"github.com/raedahgroup/godcr/fyne/assets"
	"github.com/raedahgroup/godcr/fyne/pages/handler/values"
	"github.com/raedahgroup/godcr/fyne/widgets"
)

type WalletPageObject struct {
	icons        map[string]*fyne.StaticResource
	successLabel *widgets.BorderedText

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
		assets.WalletAccount, assets.ImportedAccount, assets.MoreIcon, assets.Edit, assets.Alert, assets.InfoIcon,
		assets.Checkmark, assets.Crossmark)

	walletPage.walletAccountsBox = make([]*widgets.Box, len(walletPage.OpenedWallets))
	walletPage.walletExpandCollapseIcon = make([]*widget.Icon, len(walletPage.OpenedWallets))

	walletPage.successLabel = widgets.NewBorderedText("Wallet renamed", fyne.NewSize(16, 20), values.Green)
	walletPage.successLabel.Container.Hide()

	walletPage.WalletPageContents.Append(widgets.NewVSpacer(values.Padding))

	err = walletPage.initBaseWidgets()
	if err != nil {
		return err
	}

	walletPage.WalletPageContents.Append(widgets.NewVSpacer(values.SpacerSize8))
	walletPage.WalletPageContents.Append(widget.NewHBox(layout.NewSpacer(), walletPage.successLabel.Container, layout.NewSpacer()))
	walletPage.WalletPageContents.Append(widgets.NewVSpacer(values.SpacerSize8))

	err = walletPage.accountSelector()
	if err != nil {
		return err
	}

	return nil
}

func (walletPage *WalletPageObject) showLabel(Text string, object *widgets.BorderedText) {
	object.SetText(Text)
	object.Container.Show()
	walletPage.WalletPageContents.Refresh()
	time.AfterFunc(time.Second*5, func() {
		object.Container.Hide()
		walletPage.WalletPageContents.Refresh()
	})
}
