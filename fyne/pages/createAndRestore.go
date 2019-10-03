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
	app.Window.Resize(fyne.NewSize(370, 626))
	app.Window.SetFixedSize(true)
	fyne.CurrentApp().Settings().SetTheme(theme.LightTheme())
	app.Window.ShowAndRun()
}

func (app *AppInterface) createAndRestoreWalletPage() fyne.CanvasObject {
	icons, err := assets.Get(assets.DecredLogo, assets.Restore, assets.Add)
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
			layout.NewFixedGridLayout(fyne.NewSize(308, 56)), createWallet),
		widgets.NewVSpacer(5),
		fyne.NewContainerWithLayout(
			layout.NewFixedGridLayout(fyne.NewSize(308, 56)), restoreWallet))

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

	return widget.NewHBox(widgets.NewHSpacer(24), page)
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

	icons, err := assets.Get(assets.Checkmark, assets.Back)
	if err != nil {
		return app.DisplayLaunchErrorAndExit(err.Error())
	}

	var textbox = make([]*widget.Entry, 33)
	var layouts = make([]*fyne.Container, 33)
	wordlist := dcrlibwallet.PGPWordList()
	horizontalTextBoxes := widget.NewHBox()

	errorLabel := canvas.NewText("Failed to restore. Please verify all words and try again.", color.RGBA{255, 0, 0, 225})
	errorLabel.Alignment = fyne.TextAlignCenter
	errorLabel.Hide()

	wordlistDropdown := func(start, stop, textboxIndex int, val string) {
		if len(val) <= 1 {
			return
		}

		var menuItem []*fyne.MenuItem
		var popup *widget.PopUp

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

		textbox[textboxIndex].OnChanged = func(word string) {
			var allCompleted = true
			wordlistDropdown(0, len(wordlist), textboxIndex, word)

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
		vertical := widget.NewVBox()
		for k := i; k < i+11; k++ {
			number := widget.NewLabel(fmt.Sprintf("%d.", k+1))
			if k+1 > 9 {
				vertical.Append(widget.NewHBox(number, layouts[k]))
			} else {
				vertical.Append(widget.NewHBox(widgets.NewHSpacer(5), number, layouts[k]))
			}
		}
		horizontalTextBoxes.Append(vertical)
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
		widget.NewHBox(backButton, widgets.NewHSpacer(16), widget.NewLabel("Restore from seed phrase")),
		widgets.NewVSpacer(18),
		widget.NewLabelWithStyle("Enter your seed phrase in the correct order.", fyne.TextAlignCenter, fyne.TextStyle{Monospace: true}),
		widgets.NewVSpacer(15),
		errorLabel,
		textBoxContainer,
		widgets.NewVSpacer(10),
		buttonContainer,
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
