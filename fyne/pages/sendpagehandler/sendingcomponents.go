package sendpagehandler

import (
	"image/color"

	"github.com/decred/dcrd/dcrutil"

	"fyne.io/fyne"
	"fyne.io/fyne/canvas"
	"fyne.io/fyne/layout"
	"fyne.io/fyne/widget"

	"github.com/raedahgroup/dcrlibwallet"
	"github.com/raedahgroup/godcr/fyne/widgets"
)

func SendingDestinationComponents(onAccountChange func(), accountLabel string, receiveAccountIcon, collapseIcon fyne.Resource,
	multiWallet *dcrlibwallet.MultiWallet, walletIDs []int, sendingSelectedWalletID *int, accountBoxes []*widget.Box,
	selectedAccountLabel *widget.Label, selectedAccountBalanceLabel *widget.Label, selectedWalletLabel *canvas.Text,
	transactionFee, transactionCost, balance, size *widget.Label, amountEntry *widget.Entry, amountErrorLabel *canvas.Text,
	contents *widget.Box, nextButton *widgets.Button) (container *fyne.Container, destinationAddressEntry *widget.Entry, destinationAddressErrorLabel *canvas.Text) {

	fromLabel := canvas.NewText("To", color.RGBA{61, 88, 115, 255})
	fromLabel.TextStyle.Bold = true

	accountBox := CreateAccountSelector(onAccountChange, accountLabel, receiveAccountIcon, collapseIcon, multiWallet, walletIDs, sendingSelectedWalletID,
		accountBoxes, selectedAccountLabel, selectedAccountBalanceLabel, selectedWalletLabel, contents)

	accountBox.Hide()

	destinationAddressEntry, destinationAddressErrorLabel = destinationAddressEntryComponent(transactionFee, transactionCost, balance, size,
		amountEntry, amountErrorLabel, contents, nextButton)

	destinationAddressErrorLabel.Hide()

	sendToAccountLabel := canvas.NewText("Send to account", color.RGBA{R: 41, G: 112, B: 255, A: 255})
	sendToAccountLabel.TextSize = 12

	destinationAddressContainer := fyne.NewContainerWithLayout(layout.NewFixedGridLayout(
		fyne.NewSize(widget.NewLabel("TsfDLrRkk9ciUuwfp2b8PawwnukYD7yAjGd").MinSize().Width, destinationAddressEntry.MinSize().Height)), destinationAddressEntry)

	spacer := widgets.NewVSpacer(10)

	switchingComponentButton := widgets.NewClickableBox(widget.NewVBox(sendToAccountLabel), func() {
		if accountBox.Hidden {
			sendToAccountLabel.Text = switchToSendToAddress
			accountBox.Show()
			destinationAddressContainer.Hide()
			destinationAddressErrorLabel.Hide()
			spacer.Hide()

			if amountEntry.Text != "" {
				nextButton.Enable()
			} else {
				nextButton.Disable()
			}

		} else {
			sendToAccountLabel.Text = switchToSendToAccount
			destinationAddressContainer.Show()
			accountBox.Hide()
			spacer.Show()

			if amountEntry.Text != "" && destinationAddressEntry.Text != "" && destinationAddressErrorLabel.Hidden {
				nextButton.Enable()
			} else {
				nextButton.Disable()
			}
		}

		container.Refresh()
		contents.Refresh()
		amountEntry.OnChanged(amountEntry.Text)
	})

	box := widget.NewVBox(
		widget.NewHBox(fromLabel, layout.NewSpacer(), switchingComponentButton),
		accountBox,
		destinationAddressContainer,
		destinationAddressErrorLabel,
		spacer)

	container = fyne.NewContainerWithLayout(layout.NewFixedGridLayout(box.MinSize()), box)
	return
}

func destinationAddressEntryComponent(transactionFee, transactionCost, balance, size *widget.Label, amountEntry *widget.Entry, amountErrorLabel *canvas.Text,
	contents *widget.Box, nextButton *widgets.Button) (destinationAddressEntry *widget.Entry, destinationAddressErrorLabel *canvas.Text) {

	destinationAddressErrorLabel = canvas.NewText("", color.RGBA{237, 109, 71, 255})
	destinationAddressErrorLabel.TextSize = 12
	destinationAddressErrorLabel.Hide()

	destinationAddressEntry = widget.NewEntry()
	destinationAddressEntry.SetPlaceHolder(destinationAddressPlaceHolder)

	destinationAddressEntry.OnChanged = func(address string) {
		if destinationAddressEntry.Text == "" {
			destinationAddressErrorLabel.Hide()
			contents.Refresh()
			return
		}

		_, err := dcrutil.DecodeAddress(address)
		if err != nil {
			destinationAddressErrorLabel.Text = invalidAddress
			destinationAddressErrorLabel.Show()
			setLabelText(nilAmount, transactionFee, transactionCost, balance)
			setLabelText(nilByte, size)

		} else {
			destinationAddressErrorLabel.Hide()
		}

		if amountEntry.Text != "" && amountErrorLabel.Hidden && destinationAddressErrorLabel.Hidden {
			nextButton.Enable()
		} else {
			nextButton.Disable()
		}

		contents.Refresh()
	}

	return
}

func setLabelText(Text string, objects ...*widget.Label) {
	for _, object := range objects {
		object.SetText(Text)
		object.Refresh()
	}
}
