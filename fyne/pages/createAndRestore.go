package pages

import (
	"fmt"
	"image/color"
	"log"
	"strings"

	"fyne.io/fyne"
	"fyne.io/fyne/canvas"
	"fyne.io/fyne/layout"
	"fyne.io/fyne/theme"
	"fyne.io/fyne/widget"

	"github.com/decred/slog"
	"github.com/raedahgroup/dcrlibwallet"
	"github.com/raedahgroup/godcr/fyne/assets"
	"github.com/raedahgroup/godcr/fyne/widgets"
)

type AppInterface struct {
	Log            slog.Logger
	Dcrlw          *dcrlibwallet.LibWallet
	Window         fyne.Window
	AppDisplayName string

	tabMenu *widget.TabContainer
}

func ShowCreateAndRestoreWalletPage(dcrlw *dcrlibwallet.LibWallet, window fyne.Window, tabmenu *widget.TabContainer, log slog.Logger) {
	var app = AppInterface{
		Log:            log,
		Dcrlw:          dcrlw,
		Window:         window,
		AppDisplayName: "GoDCR",
		tabMenu:        tabmenu,
	}

	app.Window.SetContent(app.createAndRestoreWalletPage())
	app.Window.CenterOnScreen()
	app.Window.Resize(fyne.NewSize(370, 616))
	app.Window.SetFixedSize(true)
	fyne.CurrentApp().Settings().SetTheme(theme.LightTheme())
	app.Window.ShowAndRun()
}

func (app *AppInterface) createAndRestoreWalletPage() fyne.CanvasObject {
	icons, err := assets.GetIcons(assets.DecredLogo, assets.Restore, assets.Add)
	if err != nil {
		return app.DisplayLaunchErrorAndExit(err.Error())
	}

	greenBar := canvas.NewRectangle(color.RGBA{45, 216, 163, 255})
	blueBar := canvas.NewRectangle(color.RGBA{41, 112, 255, 255})

	greenBar.SetMinSize(fyne.NewSize(312, 56))
	blueBar.SetMinSize(fyne.NewSize(312, 56))

	restoreLabel := canvas.NewText("Restore an existing wallet", color.White)
	createLabel := canvas.NewText("Create a new wallet", color.White)

	createWallet := widgets.NewClickableBox(widget.NewVBox(
		fyne.NewContainerWithLayout(layout.NewBorderLayout(nil, nil, nil, nil), blueBar,
			fyne.NewContainerWithLayout(layout.NewHBoxLayout(),
				widgets.NewHSpacer(16), widget.NewIcon(icons[assets.Add]), widgets.NewHSpacer(16), createLabel))),
		func() {
			app.createSpendingPasswordPopup("")
		})

	restoreWallet := widgets.NewClickableBox(widget.NewVBox(
		fyne.NewContainerWithLayout(layout.NewBorderLayout(nil, nil, nil, nil), greenBar,
			fyne.NewContainerWithLayout(layout.NewHBoxLayout(),
				widgets.NewHSpacer(16), widget.NewIcon(icons[assets.Restore]), widgets.NewHSpacer(16), restoreLabel))),
		func() {
			app.Window.SetContent(app.restoreWalletPage())
		})

	image := canvas.NewImageFromResource(icons[assets.DecredLogo])
	image.FillMode = canvas.ImageFillOriginal

	createAndRestoreButtons := widget.NewVBox(
		fyne.NewContainerWithLayout(
			layout.NewFixedGridLayout(fyne.NewSize(285, 56)), createWallet),
		widgets.NewVSpacer(5),
		fyne.NewContainerWithLayout(
			layout.NewFixedGridLayout(fyne.NewSize(285, 56)), restoreWallet))

	// canvas doesnt support escaping characters therefore the hack
	godcrLabel := canvas.NewText("Welcome to", color.Black)
	godcrText := canvas.NewText("GoDCR", color.Black)

	godcrText.Alignment = fyne.TextAlignLeading
	godcrText.TextSize = 24
	godcrLabel.Alignment = fyne.TextAlignLeading
	godcrLabel.TextSize = 24

	page := widget.NewVBox(
		widgets.NewVSpacer(24),
		widget.NewHBox(image, layout.NewSpacer()),
		widgets.NewVSpacer(24),
		godcrLabel,
		godcrText,
		layout.NewSpacer(),
		createAndRestoreButtons,
		widgets.NewVSpacer(24))

	return widget.NewHBox(layout.NewSpacer(), page, layout.NewSpacer())
}

func (app *AppInterface) restoreWalletPage() fyne.CanvasObject {
	app.Window.SetOnClosed(func() {
		app.Window = fyne.CurrentApp().NewWindow(app.AppDisplayName)
		app.Window.SetContent(app.createAndRestoreWalletPage())
		app.Window.CenterOnScreen()
		app.Window.Resize(fyne.NewSize(360, 616))
		app.Window.SetFixedSize(true)
		app.Window.Show()
	})

	icons, err := assets.GetIcons(assets.Checkmark)
	if err != nil {
		return app.DisplayLaunchErrorAndExit(err.Error())
	}

	var textbox = make([]*widget.Entry, 33)
	wordlist := dcrlibwallet.PGPWordList()
	horizontalTextBoxes := widget.NewHBox()

	errorLabel := canvas.NewText("Failed to restore. Please verify all words and try again.", color.RGBA{255, 0, 0, 225})
	errorLabel.Alignment = fyne.TextAlignCenter
	errorLabel.Hide()

	wordlistDropdown := func(start, stop, textboxIndex int, val string) {
		if len(val) <= 1 {
			return
		}
		var popup *widget.PopUp
		var menuItem []*fyne.MenuItem

		for i := start; i < stop; i++ {
			index := i
			toLowerWordList := strings.ToLower(wordlist[i])
			toLowerVal := strings.ToLower(val)
			if strings.HasPrefix(toLowerWordList, toLowerVal) {
				menuItem = append(menuItem, fyne.NewMenuItem(wordlist[i], func() {
					textbox[textboxIndex].SetText(wordlist[index])
					popup.Hide()
				}))
			}
		}

		// do not show popup if there's no text to display
		if len(menuItem) == 0 {
			return
		}

		popup = widget.NewPopUpMenu(fyne.NewMenu("", menuItem...), app.Window.Canvas())
		popup.Move(fyne.CurrentApp().Driver().AbsolutePositionForObject(textbox[textboxIndex]).Add(
			fyne.NewPos(0, textbox[textboxIndex].Size().Height)))
	}

	var restoreButton = widget.NewButton("Restore", func() {
		var seed string
		for i := 0; i < 32; i++ {
			seed += textbox[i].Text + " "
		}
		seed += textbox[32].Text

		if dcrlibwallet.VerifySeed(seed) {
			icon := canvas.NewImageFromResource(icons[assets.Checkmark])
			icon.FillMode = canvas.ImageFillOriginal

			windowContent := app.Window.Content()

			if box, ok := windowContent.(*widget.Box); ok {
				box.Children = []fyne.CanvasObject{
					layout.NewSpacer(),
					icon,
					widgets.NewVSpacer(24),
					widget.NewLabelWithStyle("Your wallet is successfully restored", fyne.TextAlignCenter, fyne.TextStyle{Bold: true}),
					widgets.NewVSpacer(16),
					widget.NewLabelWithStyle("Now create a spending password to protect your funds.", fyne.TextAlignCenter, fyne.TextStyle{}),
					widgets.NewVSpacer(172),
					widget.NewHBox(layout.NewSpacer(), widget.NewButton("Create a spending password", func() {
						app.createSpendingPasswordPopup(seed)
					}),
						layout.NewSpacer()), widgets.NewVSpacer(16)}

				widget.Refresh(box)
			}
		} else {
			errorLabel.Show()
		}
	})
	restoreButton.Disable()

	// initialize all textboxes
	for i := 0; i < 33; i++ {
		textboxIndex := i
		textbox[textboxIndex] = widget.NewEntry()
		textbox[textboxIndex].SetPlaceHolder(fmt.Sprintf("Word %d", i+1))

		textbox[textboxIndex].OnChanged = func(word string) {
			wordlistDropdown(0, len(wordlist), textboxIndex, word)
			var allCompleted = true
			for j := 0; j < 33; j++ {
				if textbox[j].Text == "" {
					allCompleted = false
				}
			}

			if allCompleted == true {
				restoreButton.Enable()
			} else {
				restoreButton.Disable()
			}
		}
	}

	for i := 0; i < 3; i++ {
		vertical := widget.NewVBox()
		for k := i; k < 33; k += 3 {
			vertical.Append(textbox[k])
		}
		horizontalTextBoxes.Append(vertical)
	}

	textBoxContainer := widget.NewHBox(layout.NewSpacer(), horizontalTextBoxes.Children[0], layout.NewSpacer(),
		horizontalTextBoxes.Children[1], layout.NewSpacer(), horizontalTextBoxes.Children[2], layout.NewSpacer())

	buttonContainer := widget.NewHBox(layout.NewSpacer(), restoreButton, layout.NewSpacer())

	return widget.NewVBox(widgets.NewVSpacer(10), errorLabel, textBoxContainer, widgets.NewVSpacer(10), buttonContainer)
}

func (app *AppInterface) createSpendingPasswordPopup(seed string) {
	var popup *widget.PopUp
	popupContent := widget.NewVBox()

	popup = widget.NewModalPopUp(widget.NewVBox(
		widget.NewLabelWithStyle("Create a spending password", fyne.TextAlignLeading, fyne.TextStyle{Bold: true}), popupContent),
		app.Window.Canvas())

	popupContent.Children = []fyne.CanvasObject{app.passwordPopup(popup, seed)}
	widget.Refresh(popupContent)
}

func (app *AppInterface) passwordPopup(popup *widget.PopUp, seed string) fyne.CanvasObject {
	displayError := func(err error) {
		log.Println("could not generate seed", err.Error())
		newWindow := fyne.CurrentApp().NewWindow(app.AppDisplayName)
		newWindow.SetContent(widget.NewVBox(
			widget.NewLabelWithStyle(fmt.Sprintf("Could not generate seed, %s", err.Error()), fyne.TextAlignCenter, fyne.TextStyle{}),
			widget.NewHBox(layout.NewSpacer(), widget.NewButton("Close", func() { newWindow.Close() }), layout.NewSpacer())))

		newWindow.CenterOnScreen()
		newWindow.Show()
		newWindow.SetFixedSize(true)
	}

	icons, err := assets.GetIcons(assets.Reveal, assets.Conceal, assets.Loader)
	if err != nil {
		return app.DisplayLaunchErrorAndExit(err.Error())
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
		} else if password.Text != "" && password.Text == confirmPassword.Text {
			errorLabel.Hide()
			createButton.Enable()
		}
	}

	cancelLabel := canvas.NewText("Cancel", color.RGBA{41, 112, 255, 255})
	cancelLabel.TextStyle.Bold = true
	cancelButton := widgets.NewClickableBox(widget.NewHBox(cancelLabel), func() { popup.Hide() })

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
		if seed == "" {
			seed, err = dcrlibwallet.GenerateSeed()
			if err != nil {
				enableCancelButton()
				displayError(err)
				return
			}
		}

		err = app.Dcrlw.CreateWallet(password.Text, seed)

		if err != nil {
			enableCancelButton()
			displayError(err)
			log.Println("could not create wallet", err.Error())
			return
		}

		popup.Hide()
		app.Window.SetFixedSize(false)
		app.Window.SetOnClosed(nil)
		app.Window.SetContent(app.tabMenu)
		app.tabMenu.CreateRenderer().ApplyTheme()
		// apparently tabMenu was initialized in fyne.go line 57,
		// this sets theme of tabmenu to default black
		// resetting theme to light theme when tabmenu is in view fixes the font.
		fyne.CurrentApp().Settings().SetTheme(theme.LightTheme())
	})

	createButton.Disable()

	var passwordConceal *widgets.ImageButton
	passwordConceal = widgets.NewImageButton(icons[assets.Reveal], nil, func() {
		if password.Password {
			passwordConceal.SetIcon(icons[assets.Conceal])
			password.Password = false
		} else {
			passwordConceal.SetIcon(icons[assets.Reveal])
			password.Password = true
		}
		// reveal texts
		password.SetText(password.Text)
	})

	var confirmPasswordConceal *widgets.ImageButton
	confirmPasswordConceal = widgets.NewImageButton(icons[assets.Reveal], nil, func() {
		if confirmPassword.Password {
			confirmPasswordConceal.SetIcon(icons[assets.Conceal])
			confirmPassword.Password = false
		} else {
			confirmPasswordConceal.SetIcon(icons[assets.Reveal])
			confirmPassword.Password = true
		}
		// reveal texts
		confirmPassword.SetText(confirmPassword.Text)
	})

	return widget.NewHBox(widgets.NewHSpacer(10), widget.NewVBox(
		widgets.NewVSpacer(10), widget.NewHBox(layout.NewSpacer(),
			fyne.NewContainerWithLayout(layout.NewFixedGridLayout(fyne.NewSize(264, password.MinSize().Height)), password),
			passwordConceal, layout.NewSpacer()),
		passwordLength, widget.NewHBox(layout.NewSpacer(),
			fyne.NewContainerWithLayout(layout.NewFixedGridLayout(fyne.NewSize(160, widget.NewLabel("0%").MinSize().Height)), passwordStrength)),
		widget.NewHBox(layout.NewSpacer(),
			fyne.NewContainerWithLayout(layout.NewFixedGridLayout(fyne.NewSize(264, confirmPassword.MinSize().Height)), confirmPassword),
			confirmPasswordConceal, layout.NewSpacer()),
		confirmPasswordLength,
		widget.NewHBox(layout.NewSpacer(), cancelButton, widgets.NewHSpacer(24), createButton),
		errorLabel,
		widgets.NewVSpacer(10)), widgets.NewHSpacer(10))
}

// DisplayLaunchErrorAndExit displays the error message to users.
func (app *AppInterface) DisplayLaunchErrorAndExit(errorMessage string) fyne.CanvasObject {
	return widget.NewVBox(
		widget.NewLabelWithStyle(errorMessage, fyne.TextAlignCenter, fyne.TextStyle{}),

		widget.NewHBox(
			layout.NewSpacer(),
			widget.NewButton("Exit", app.Window.Close), // closing the window will trigger app.tearDown()
			layout.NewSpacer(),
		))
}
