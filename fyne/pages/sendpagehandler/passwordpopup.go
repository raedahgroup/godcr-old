package sendpagehandler

import (
	"image/color"
	"log"

	"fyne.io/fyne"
	"fyne.io/fyne/canvas"
	"fyne.io/fyne/layout"
	"fyne.io/fyne/widget"

	"github.com/raedahgroup/dcrlibwallet"

	"github.com/raedahgroup/godcr/fyne/layouts"
	"github.com/raedahgroup/godcr/fyne/widgets"
)

func PasswordPopUp(initOnConfirmation func(string) error, initOnCancel, onError, extraCalls func(), conceal, reveal fyne.Resource, window fyne.Window) {
	errorLabel := canvas.NewText("Wrong spending password. Please try again.", color.RGBA{237, 109, 71, 255})
	errorLabel.Alignment = fyne.TextAlignCenter
	errorLabel.TextSize = 12
	errorLabel.Hide()

	var confirmButton *widgets.Button

	walletPassword := widget.NewPasswordEntry()
	walletPassword.SetPlaceHolder("Spending Password")
	walletPassword.OnChanged = func(value string) {
		if value == "" {
			confirmButton.Disable()
		} else if confirmButton.Disabled() {
			confirmButton.Enable()
		}
	}

	var sendingPasswordPopup *widget.PopUp
	var popupContent *widget.Box

	cancelLabel := canvas.NewText("Cancel", color.RGBA{41, 112, 255, 255})
	cancelLabel.TextStyle.Bold = true

	cancelButton := widgets.NewClickableBox(widget.NewHBox(cancelLabel), func() {
		sendingPasswordPopup.Hide()
		initOnCancel()
	})

	confirmButton = widgets.NewButton(color.RGBA{41, 112, 255, 255}, "Confirm", func() {
		confirmButton.Disable()
		cancelButton.Disable()

		var err error
		if initOnConfirmation != nil {
			err = initOnConfirmation(walletPassword.Text)
		}

		if err != nil {
			// do not exit password popup on invalid passphrase
			if err.Error() == dcrlibwallet.ErrInvalidPassphrase {
				errorLabel.Show()
				// this is an hack as selective refresh to errorLabel doesn't work
				popupContent.Refresh()
				confirmButton.Enable()
				cancelButton.Disable()
			} else {
				log.Println(err)
				sendingPasswordPopup.Hide()
				if onError != nil {
					onError()
				}
			}

			return
		}

		extraCalls()
		sendingPasswordPopup.Hide()
	})
	confirmButton.SetMinSize(fyne.NewSize(91, 40))
	confirmButton.Disable()

	var passwordConceal *widgets.ImageButton
	passwordConceal = widgets.NewImageButton(reveal, nil, func() {
		if walletPassword.Password {
			passwordConceal.SetIcon(conceal)
			walletPassword.Password = false
		} else {
			passwordConceal.SetIcon(reveal)
			walletPassword.Password = true
		}
		// reveal texts
		walletPassword.SetText(walletPassword.Text)
	})

	popupContent = widget.NewHBox(
		widgets.NewHSpacer(24),
		widget.NewVBox(
			widgets.NewVSpacer(24),
			widget.NewLabelWithStyle("Confirm to send", fyne.TextAlignLeading, fyne.TextStyle{Bold: true}),
			widgets.NewVSpacer(40),
			fyne.NewContainerWithLayout(layouts.NewPasswordLayout(fyne.NewSize(312, walletPassword.MinSize().Height)), walletPassword, passwordConceal),
			errorLabel,
			widgets.NewVSpacer(20),
			widget.NewHBox(layout.NewSpacer(), cancelButton, widgets.NewHSpacer(24), confirmButton.Container),
			widgets.NewVSpacer(24),
		),
		widgets.NewHSpacer(24),
	)

	sendingPasswordPopup = widget.NewModalPopUp(popupContent, window.Canvas())
}
