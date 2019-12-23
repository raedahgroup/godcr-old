package walletshandler

import (
	"fmt"

	"fyne.io/fyne"
	"fyne.io/fyne/layout"
	"fyne.io/fyne/widget"

	"github.com/raedahgroup/godcr/fyne/assets"
	"github.com/raedahgroup/godcr/fyne/widgets"
)

func (walletPage *WalletPageObject) initBaseWidgets() error {
	walletLabel := widget.NewLabelWithStyle("Wallets", fyne.TextAlignLeading, fyne.TextStyle{Bold: true})

	icon, err := assets.GetIcons(assets.AddWallet)
	if err != nil {
		return err
	}

	addWallet := widgets.NewImageButton(icon[assets.AddWallet], nil, func() {
		fmt.Println("Helllll0")
	})

	walletPage.WalletPageContents.Append(widget.NewHBox(walletLabel, layout.NewSpacer(), fyne.NewContainerWithLayout(
		layout.NewFixedGridLayout(fyne.NewSize(30, 30)), addWallet)))
	return nil
}
