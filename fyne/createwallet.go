package fyne

import (
	"fmt"

	"fyne.io/fyne/widget"
	godcrApp "github.com/raedahgroup/godcr/app"
)

func (app *fyneApp) showCreateWalletWindow() {
	app.mainWindow.SetTitle(fmt.Sprintf("%s Create Wallet", godcrApp.DisplayName))

	createWalletButton := widget.NewButton("Create Wallet", func() {
		// todo this function should not quit the app but actually create a wallet
		// and then open the sync window using app.showSyncWindow()
		app.Quit()
	})

	// todo complete this create wallet window's content
	app.mainWindow.SetContent(createWalletButton)

	app.mainWindow.CenterOnScreen()
	app.mainWindow.Show()
}
