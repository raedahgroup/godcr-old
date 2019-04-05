package fyne

import (
	"fmt"

	godcrApp "github.com/raedahgroup/godcr/app"
	"fyne.io/fyne/widget"
)

func (app *fyneApp) showCreateWalletWindow() {
	window := app.NewWindow(fmt.Sprintf("%s Create Wallet", godcrApp.DisplayName))

	createWalletButton := widget.NewButton("Create Wallet", func() {
		// todo this function should not quit the app but actually create a wallet and then open the sync window
		app.Quit()
	})

	// todo complete this create wallet window's content
	window.SetContent(createWalletButton)

	window.CenterOnScreen()
	window.Show()
}
