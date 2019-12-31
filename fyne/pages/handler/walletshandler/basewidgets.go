package walletshandler

import (
	"fmt"

	"fyne.io/fyne"
	"fyne.io/fyne/layout"
	"fyne.io/fyne/widget"

	"github.com/raedahgroup/godcr/fyne/assets"
	"github.com/raedahgroup/godcr/fyne/pages/handler/values"
	"github.com/raedahgroup/godcr/fyne/widgets"
)

func (walletPage *WalletPageObject) initBaseWidgets() error {
	walletLabel := widget.NewLabelWithStyle("Wallets", fyne.TextAlignLeading, fyne.TextStyle{Bold: true})

	addWallet := widgets.NewImageButton(walletPage.icons[assets.AddWallet], nil, func() {
		fmt.Println("Helllll0")
	})

	walletPage.WalletPageContents.Append(widget.NewHBox(walletLabel, layout.NewSpacer(), fyne.NewContainerWithLayout(
		layout.NewFixedGridLayout(fyne.NewSize(30, 30)), addWallet), widgets.NewHSpacer(values.Padding)))
		
	return nil
}
