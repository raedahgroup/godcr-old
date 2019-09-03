package pages

import (
	"context"

	"fyne.io/fyne"
	"fyne.io/fyne/widget"

	"github.com/raedahgroup/godcr/app/wallet"
)

func ShowCreateAndRestoreWalletPage(wallet wallet.Wallet, window fyne.Window, ctx context.Context) {
	tabs := widget.NewTabContainer(
		widget.NewTabItem("Create a new wallet", widget.NewLabel("Not Implemented yet")),
		widget.NewTabItem("Restore an existing wallet", widget.NewLabel("Not Implemented yet")))
	window.SetContent(tabs)

	window.CenterOnScreen()
	window.ShowAndRun()
}

// todo: create an user interface for create and restore wallet and then pass to tab item.
