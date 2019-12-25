package stakingpagehandler

import (
	"fyne.io/fyne"
	"fyne.io/fyne/layout"
	"fyne.io/fyne/theme"
	"fyne.io/fyne/widget"

	"github.com/raedahgroup/godcr/fyne/assets"
	"github.com/raedahgroup/godcr/fyne/widgets"
)

func (stakingPage *StakingPageObjects) summaryWalletList() {
	walletListWidget := widget.NewVBox()

	selectedWalletLabel := widget.NewLabel("wallet-1")
	var txWalletSelectionPopup *widget.PopUp

	index := 0
	checkmarkIcon := widget.NewIcon(theme.ConfirmIcon())
	if index != 0 {
		checkmarkIcon.Hide()
	}

	walletContainer := widget.NewHBox(
		widget.NewLabel("wallet-1"),
		checkmarkIcon,
		widgets.NewHSpacer(5),
	)

	walletListWidget.Append(widgets.NewClickableBox(walletContainer, func() {

	}))

	// txWalletSelectionPopup create a popup that has tx wallet
	txWalletSelectionPopup = widget.NewPopUp(fyne.NewContainerWithLayout(
		layout.NewFixedGridLayout(fyne.NewSize(walletListWidget.MinSize().Width, 50)), widget.NewScrollContainer(walletListWidget)), stakingPage.Window.Canvas())
	txWalletSelectionPopup.Hide()

	walletListTab := widget.NewHBox(
		selectedWalletLabel,
		widgets.NewHSpacer(10),
		widget.NewIcon(stakingPage.icons[assets.CollapseIcon]),
	)

	// walletDropDown creates a popup like dropdown that holds the list of available wallets.
	var walletDropDown *widgets.ClickableBox
	walletDropDown = widgets.NewClickableBox(walletListTab, func() {
		txWalletSelectionPopup.Move(fyne.CurrentApp().Driver().AbsolutePositionForObject(
			walletDropDown).Add(fyne.NewPos(0, walletDropDown.Size().Height)))
		txWalletSelectionPopup.Show()
	})

	stakingPage.StakingPageContents.Append(widget.NewHBox(walletDropDown))
}
