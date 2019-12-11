package multipagecomponents

import (
	"errors"
	"image/color"
	"log"

	"fyne.io/fyne"
	"fyne.io/fyne/canvas"
	"fyne.io/fyne/layout"
	"fyne.io/fyne/widget"

	"github.com/raedahgroup/dcrlibwallet"

	"github.com/raedahgroup/godcr/fyne/assets"
	"github.com/raedahgroup/godcr/fyne/layouts"
	"github.com/raedahgroup/godcr/fyne/pages/constantvalues"
	"github.com/raedahgroup/godcr/fyne/widgets"
)

type PasswordPopUpObjects struct {
	InitOnConfirmation        func(string) error
	InitOnCancel, InitOnError func()
	ExtraCalls                func() // ExtraCalls is called when InitOnConfirmation is called and doesnt throw an error

	Window fyne.Window
}

func (objects *PasswordPopUpObjects) PasswordPopUp() error {
	icons, err := assets.GetIcons(assets.Conceal, assets.Reveal)
	if err != nil {
		return errors.New(constantvalues.PasswordPopupIconsErr)
	}

	errorLabel := canvas.NewText(constantvalues.WrongPasswordErr, color.RGBA{237, 109, 71, 255})
	errorLabel.Alignment = fyne.TextAlignLeading
	errorLabel.TextSize = 12
	errorLabel.Hide()

	var confirmButton *widgets.Button

	walletPassword := widget.NewPasswordEntry()
	walletPassword.SetPlaceHolder(constantvalues.SpendingPasswordText)
	walletPassword.OnChanged = func(value string) {
		if value == "" {
			confirmButton.Disable()
		} else if confirmButton.Disabled() {
			confirmButton.Enable()
		}
	}

	var sendingPasswordPopup *widget.PopUp
	var popupContent *widget.Box

	cancelLabel := canvas.NewText(constantvalues.Cancel, color.RGBA{41, 112, 255, 255})
	cancelLabel.TextStyle.Bold = true

	cancelButton := widgets.NewClickableBox(widget.NewHBox(cancelLabel), func() {
		sendingPasswordPopup.Hide()
		objects.InitOnCancel()
	})

	confirmButton = widgets.NewButton(color.RGBA{41, 112, 255, 255}, constantvalues.Confirm, func() {
		confirmButton.Disable()
		cancelButton.Disable()

		var err error
		if objects.InitOnConfirmation != nil {
			err = objects.InitOnConfirmation(walletPassword.Text)
		}

		if err != nil {
			// do not exit password popup on invalid passphrase
			if err.Error() == dcrlibwallet.ErrInvalidPassphrase {
				errorLabel.Show()
				popupContent.Refresh()
				confirmButton.Enable()
				cancelButton.Enable()
			} else {
				log.Println(err)
				sendingPasswordPopup.Hide()
				if objects.InitOnError != nil {
					objects.InitOnError()
				}
			}

			return
		}

		objects.ExtraCalls()
		sendingPasswordPopup.Hide()
	})
	confirmButton.SetMinSize(fyne.NewSize(91, 40))
	confirmButton.Disable()

	var passwordConceal *widgets.ImageButton
	passwordConceal = widgets.NewImageButton(icons[assets.Reveal], nil, func() {
		if walletPassword.Password {
			passwordConceal.SetIcon(icons[assets.Conceal])
			walletPassword.Password = false
		} else {
			passwordConceal.SetIcon(icons[assets.Reveal])
			walletPassword.Password = true
		}
		// reveal texts
		walletPassword.SetText(walletPassword.Text)
	})

	popupContent = widget.NewHBox(
		widgets.NewHSpacer(24),
		widget.NewVBox(
			widgets.NewVSpacer(24),
			widget.NewLabelWithStyle(constantvalues.ConfirmToSend, fyne.TextAlignLeading, fyne.TextStyle{Bold: true}),
			widgets.NewVSpacer(40),
			fyne.NewContainerWithLayout(layouts.NewPasswordLayout(fyne.NewSize(312, walletPassword.MinSize().Height)), walletPassword, passwordConceal),
			errorLabel,
			widgets.NewVSpacer(20),
			widget.NewHBox(layout.NewSpacer(), cancelButton, widgets.NewHSpacer(24), confirmButton.Container),
			widgets.NewVSpacer(24),
		),
		widgets.NewHSpacer(24),
	)

	sendingPasswordPopup = widget.NewModalPopUp(popupContent, objects.Window.Canvas())
	return nil
}
