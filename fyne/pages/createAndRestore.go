package pages

import (
	"fmt"
	"image/color"
	"strings"

	"fyne.io/fyne"
	"fyne.io/fyne/canvas"
	"fyne.io/fyne/layout"
	"fyne.io/fyne/theme"
	"fyne.io/fyne/widget"

	"github.com/raedahgroup/dcrlibwallet"

	"github.com/raedahgroup/godcr/fyne/assets"
	"github.com/raedahgroup/godcr/fyne/values"
	"github.com/raedahgroup/godcr/fyne/widgets"
)

func (app *AppInterface) ShowCreateAndRestoreWalletPage() {
	app.Window.SetContent(app.createAndRestoreWalletPage())
	app.Window.CenterOnScreen()
	app.Window.Resize(fyne.NewSize(370, 626))
	app.Window.SetFixedSize(true)
	fyne.CurrentApp().Settings().SetTheme(theme.LightTheme())
	app.Window.ShowAndRun()
	app.tearDown()
}

func (app *AppInterface) createAndRestoreWalletPage() fyne.CanvasObject {
	icons, err := assets.GetIcons(assets.DecredLogo, assets.Restore, assets.Add)
	if err != nil {
		return app.displayErrorPage(err.Error())
	}

	greenBar := canvas.NewRectangle(values.BlueGreen)
	greenBar.SetMinSize(fyne.NewSize(312, 56))

	blueBar := canvas.NewRectangle(values.Blue)
	blueBar.SetMinSize(fyne.NewSize(312, 56))

	restoreWalletLabel := canvas.NewText("Restore an existing wallet", color.White)
	createWalletLabel := canvas.NewText("Create a new wallet", color.White)

	createWalletWidget := widgets.NewClickableBox(
		widget.NewVBox(
			fyne.NewContainerWithLayout(layout.NewBorderLayout(nil, nil, nil, nil), blueBar,
				fyne.NewContainerWithLayout(
					layout.NewHBoxLayout(),
					widgets.NewHSpacer(16), widget.NewIcon(icons[assets.Add]), widgets.NewHSpacer(16), createWalletLabel))),
		func() {
			app.createSpendingPasswordPopup("")
		})

	restoreWalletWidget := widgets.NewClickableBox(widget.NewVBox(
		fyne.NewContainerWithLayout(layout.NewBorderLayout(nil, nil, nil, nil), greenBar,
			fyne.NewContainerWithLayout(layout.NewHBoxLayout(),
				widgets.NewHSpacer(16), widget.NewIcon(icons[assets.Restore]), widgets.NewHSpacer(16), restoreWalletLabel))),
		func() {
			app.Window.SetContent(app.restoreWalletPage())
		})

	decredLogo := canvas.NewImageFromResource(icons[assets.DecredLogo])
	decredLogo.FillMode = canvas.ImageFillOriginal

	createAndRestoreButtons := widget.NewVBox(
		fyne.NewContainerWithLayout(layout.NewFixedGridLayout(fyne.NewSize(308, 56)), createWalletWidget),
		widgets.NewVSpacer(5),
		fyne.NewContainerWithLayout(layout.NewFixedGridLayout(fyne.NewSize(308, 56)), restoreWalletWidget))

	// canvas doesnt support escaping characters therefore the hack
	welcomeLabel := canvas.NewText("Welcome to", color.Black)
	welcomeLabel.Alignment = fyne.TextAlignLeading
	welcomeLabel.TextSize = 24

	godcrLabel := canvas.NewText("GoDCR", color.Black)
	godcrLabel.Alignment = fyne.TextAlignLeading
	godcrLabel.TextSize = 24

	createRestorePage := widget.NewVBox(
		widgets.NewVSpacer(24),
		widget.NewHBox(decredLogo, layout.NewSpacer()),
		widgets.NewVSpacer(24),
		welcomeLabel,
		godcrLabel,
		layout.NewSpacer(),
		createAndRestoreButtons,
		widgets.NewVSpacer(24))

	return widget.NewHBox(widgets.NewHSpacer(24), createRestorePage)
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

	icons, err := assets.GetIcons(assets.Checkmark, assets.Back)
	if err != nil {
		return app.displayErrorPage(err.Error())
	}

	var textbox = make([]*widget.Entry, 33)
	var layouts = make([]*fyne.Container, 33)
	wordlist := dcrlibwallet.PGPWordList()
	horizontalTextBoxes := widget.NewHBox()

	errorLabel := canvas.NewText("Failed to restore. Please verify all words and try again.", color.RGBA{255, 0, 0, 225})
	errorLabel.Alignment = fyne.TextAlignCenter
	errorLabel.Hide()

	wordlistDropdown := func(start, stop, textboxIndex int, wordText string) {
		if len(wordText) <= 1 {
			return
		}

		var menuItem []*fyne.MenuItem
		var wordlistPopup *widget.PopUp

		for i := start; i < stop; i++ {
			index := i
			toLowerWordList := strings.ToLower(wordlist[i])
			toLowerVal := strings.ToLower(wordText)
			if strings.HasPrefix(toLowerWordList, toLowerVal) {
				menuItem = append(menuItem, fyne.NewMenuItem(wordlist[i], func() {
					textbox[textboxIndex].SetText(wordlist[index])
					wordlistPopup.Hide()
				}))
			}
		}

		// do not show wordlistPopup if there's no text to display
		if len(menuItem) == 0 {
			return
		}

		wordlistPopup = widget.NewPopUpMenu(fyne.NewMenu("", menuItem...), app.Window.Canvas())
		wordlistPopup.Move(fyne.CurrentApp().Driver().AbsolutePositionForObject(textbox[textboxIndex]).Add(
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

			app.Window.SetContent(widget.NewVBox(
				layout.NewSpacer(),
				icon,
				widgets.NewVSpacer(24),
				widget.NewLabelWithStyle("Your wallet is successfully restored", fyne.TextAlignCenter, fyne.TextStyle{Bold: true}),
				widgets.NewVSpacer(16),
				widget.NewLabelWithStyle("Now create a spending password to protect your funds.", fyne.TextAlignCenter, fyne.TextStyle{}),
				widgets.NewVSpacer(172),
				widget.NewHBox(layout.NewSpacer(), widget.NewButton("Create a spending password", func() { app.createSpendingPasswordPopup(seed) }),
					layout.NewSpacer()), widgets.NewVSpacer(16)))
		} else {
			errorLabel.Show()
		}
	})
	restoreButton.Disable()

	// initialize all textboxes
	for i := 0; i < 33; i++ {
		textboxIndex := i
		textbox[textboxIndex] = widget.NewEntry()
		maxTextboxSize := fyne.NewSize(110, textbox[textboxIndex].MinSize().Height)
		layouts[textboxIndex] = fyne.NewContainerWithLayout(layout.NewFixedGridLayout(maxTextboxSize), textbox[textboxIndex])

		textbox[textboxIndex].OnChanged = func(wordText string) {
			var allCompleted = true
			wordlistDropdown(0, len(wordlist), textboxIndex, wordText)

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

	for i := 0; i < 33; i += 11 {
		verticalTextBoxes := widget.NewVBox()
		for k := i; k < i+11; k++ {
			number := widget.NewLabel(fmt.Sprintf("%d.", k+1))
			if k+1 > 9 {
				verticalTextBoxes.Append(widget.NewHBox(number, layouts[k]))
			} else {
				verticalTextBoxes.Append(widget.NewHBox(widgets.NewHSpacer(5), number, layouts[k]))
			}
		}

		horizontalTextBoxes.Append(verticalTextBoxes)
	}

	backButton := widgets.NewImageButton(icons[assets.Back], nil, func() {
		app.Window.SetOnClosed(nil)
		app.Window.SetContent(app.createAndRestoreWalletPage())
		app.Window.Resize(fyne.NewSize(370, 626))
	})

	textBoxContainer := widget.NewHBox(
		horizontalTextBoxes.Children[0], layout.NewSpacer(),
		horizontalTextBoxes.Children[1], layout.NewSpacer(),
		horizontalTextBoxes.Children[2])

	buttonContainer := widget.NewHBox(layout.NewSpacer(), restoreButton, layout.NewSpacer())

	return widget.NewHBox(widgets.NewHSpacer(10), widget.NewVBox(
		widgets.NewVSpacer(10),
		widget.NewHBox(backButton, widgets.NewHSpacer(16), widget.NewLabelWithStyle("Restore from seed phrase", fyne.TextAlignTrailing, fyne.TextStyle{Bold: true})),
		widgets.NewVSpacer(18),
		widget.NewLabelWithStyle("Enter your seed phrase in the correct order.", fyne.TextAlignCenter, fyne.TextStyle{}),
		widgets.NewVSpacer(15),
		errorLabel,
		textBoxContainer,
		widgets.NewVSpacer(10),
		buttonContainer,
		widgets.NewVSpacer(10)), widgets.NewHSpacer(10))
}
