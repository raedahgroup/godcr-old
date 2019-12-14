package pages

import (
	"fmt"
	"image/color"
	"log"

	"fyne.io/fyne"
	"fyne.io/fyne/canvas"
	"fyne.io/fyne/layout"
	"fyne.io/fyne/widget"

	"github.com/raedahgroup/dcrlibwallet"

	"github.com/raedahgroup/godcr/fyne/assets"
	"github.com/raedahgroup/godcr/fyne/widgets"
)

func (app *AppInterface) createSpendingPasswordPopup(seed string) {
	var passwordPopup *widget.PopUp
	popupContent := widget.NewVBox()

	passwordPopup = widget.NewModalPopUp(widget.NewHBox(widgets.NewHSpacer(10),
		widget.NewVBox(widgets.NewVSpacer(10),
			widget.NewLabelWithStyle("Create a spending password", fyne.TextAlignLeading, fyne.TextStyle{Bold: true}),
			popupContent, widgets.NewVSpacer(10)),
		widgets.NewHSpacer(10)),
		app.Window.Canvas())

	popupContent.Children = []fyne.CanvasObject{app.passwordPopup(passwordPopup, seed)}
	widget.Refresh(popupContent)
}

func (app *AppInterface) passwordPopup(passwordPopup *widget.PopUp, seed string) fyne.CanvasObject {
	displayError := func(err error) {
		log.Println("Could not generate seed", err.Error())
		newWindow := fyne.CurrentApp().NewWindow(app.AppDisplayName)
		newWindow.SetContent(widget.NewVBox(
			widget.NewLabelWithStyle(fmt.Sprintf("Could not generate seed, %s", err.Error()), fyne.TextAlignCenter, fyne.TextStyle{}),
			widget.NewHBox(layout.NewSpacer(), widget.NewButton("Close", func() { newWindow.Close() }), layout.NewSpacer())))

		newWindow.CenterOnScreen()
		newWindow.Show()
		newWindow.SetFixedSize(true)
	}

	icons, err := assets.GetIcons(assets.Loader)
	if err != nil {
		return app.displayErrorPage(err.Error())
	}

	errorLabel := canvas.NewText("Password do not match", color.RGBA{255, 0, 0, 255})
	errorLabel.TextSize = 10
	errorLabel.Hide()

	password := widget.NewPasswordEntry()
	password.SetPlaceHolder("Spending Password")
	confirmPassword := widget.NewPasswordEntry()
	confirmPassword.SetPlaceHolder("Confirm Spending Password")

	passwordLength := canvas.NewText("0", color.Black)
	passwordLength.TextSize = 10
	passwordLength.Alignment = fyne.TextAlignTrailing
	confirmPasswordLength := canvas.NewText("0", color.Black)
	confirmPasswordLength.TextSize = 10
	confirmPasswordLength.Alignment = fyne.TextAlignTrailing

	passwordStrength := widget.NewProgressBar()
	var createButton *widget.Button

	password.OnChanged = func(val string) {
		// check if password and confirm password matches only when the user fills confirmPassword textbox
		if confirmPassword.Text != "" {
			if confirmPassword.Text != password.Text {
				errorLabel.Show()
				if !createButton.Disabled() {
					createButton.Disable()
				}
			} else {
				errorLabel.Hide()
				createButton.Enable()
			}
		}

		passwordLength.Text = fmt.Sprintf("%d", len(val))
		canvas.Refresh(passwordLength)

		strength := (dcrlibwallet.ShannonEntropy(val) / 4.0)
		passwordStrength.SetValue(strength)
	}

	confirmPassword.OnChanged = func(val string) {
		confirmPasswordLength.Text = fmt.Sprintf("%d", len(val))
		canvas.Refresh(confirmPasswordLength)
		if password.Text != val {
			errorLabel.Show()
			if !createButton.Disabled() {
				createButton.Disable()
			}
		} else if password.Text != "" && password.Text == confirmPassword.Text {
			errorLabel.Hide()
			createButton.Enable()
		}
	}

	cancelLabel := canvas.NewText("Cancel", color.RGBA{41, 112, 255, 255})
	cancelLabel.TextStyle.Bold = true
	cancelButton := widgets.NewClickableBox(widget.NewHBox(cancelLabel), func() { passwordPopup.Hide() })

	createButton = widget.NewButton("Create", func() {
		createButton.SetText("")
		createButton.SetIcon(icons[assets.Loader])
		createButton.Disable()

		// disable cancel OnTapped function
		cancelButton.OnTapped = nil
		cancelLabel.Color = color.RGBA{196, 203, 210, 255}
		canvas.Refresh(cancelLabel)

		enableCancelButton := func() {
			cancelButton.OnTapped = nil
			cancelLabel.Color = color.RGBA{41, 112, 255, 255}
		}

		var err error
		var wallet *dcrlibwallet.Wallet
		if seed == "" {
			wallet, err = app.MultiWallet.CreateNewWallet("", password.Text, 0)
			if err != nil {
				enableCancelButton()
				displayError(err)
				log.Println("Could not create wallet", err.Error())
				return
			}
		} else {
			wallet, err = app.MultiWallet.RestoreWallet(seed, "", password.Text, 0)
			if err != nil {
				enableCancelButton()
				displayError(err)
				log.Println("Could not create wallet", err.Error())
				return
			}

			err = wallet.UnlockWallet([]byte(password.Text))
			if err != nil {
				log.Println("could not unlock wallet to discover account")
			}
		}

		passwordPopup.Hide()
		app.Window.SetFixedSize(false)
		app.Window.SetOnClosed(nil)
		app.setupNavigationMenu()
		app.Window.SetContent(app.tabMenu)
	})

	createButton.Disable()

	return widget.NewVBox(widgets.NewVSpacer(10),
		password,
		passwordLength,
		widget.NewHBox(layout.NewSpacer(),
			fyne.NewContainerWithLayout(
				layout.NewFixedGridLayout(fyne.NewSize(150, widget.NewLabel("0%").MinSize().Height)), passwordStrength)),
		confirmPassword,
		confirmPasswordLength,
		widget.NewHBox(layout.NewSpacer(), widgets.NewHSpacer(170), cancelButton, widgets.NewHSpacer(24), createButton),
		errorLabel,
		widgets.NewVSpacer(10))
}
