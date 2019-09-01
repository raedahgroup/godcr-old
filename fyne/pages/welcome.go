package pages

import (
	"fyne.io/fyne"
	"fyne.io/fyne/widget"
	"github.com/raedahgroup/godcr/app/wallet"
)

// todo: keep variables here that will be needed by this page at different points in time
// instead of using package-global variables
type welcomePage struct {
	wallet wallet.Wallet
}

func ShowWelcomePage(wallet wallet.Wallet) {
	page := &welcomePage{
		wallet: wallet,
	}
	page.showMainWindow()
}

func (page *welcomePage) showMainWindow() {
	pageContent := widget.NewVBox(
		widget.NewLabel("Welcome to GoDCR Decred Wallet"),
		widget.NewButton("Create New Wallet", func() {
			// todo show create wallet form/flow and use `page.wallet.CreateWallet` to create the wallet
			println("create wallet clicked")
		}),
		widget.NewButton("Restore Existing Wallet", func() {
			// todo show restore wallet form/flow and use `page.wallet.CreateWallet` to restore the wallet
			println("restore wallet clicked")
		}),
	)

	window := fyne.CurrentApp().NewWindow("GoDCR")
	window.SetContent(pageContent)
	window.CenterOnScreen()
	window.ShowAndRun()
}