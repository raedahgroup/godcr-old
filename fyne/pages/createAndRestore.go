package pages

import (
	"fyne.io/fyne"
	"fyne.io/fyne/widget"

	"github.com/raedahgroup/dcrlibwallet"
)

func ShowCreateAndRestoreWalletPage(dcrlw *dcrlibwallet.LibWallet, window fyne.Window) {
	tabs := widget.NewTabContainer(
		widget.NewTabItem("Create a new wallet", widget.NewLabel("Not Implemented yet")),
		widget.NewTabItem("Restore an existing wallet", widget.NewLabel("Not Implemented yet")))
	window.SetContent(tabs)

	window.CenterOnScreen()
	window.ShowAndRun()
}

// todo: create an user interface for create and restore wallet and then pass to tab item.
