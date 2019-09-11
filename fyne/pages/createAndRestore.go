package pages

import (
	"fmt"
	"image/color"
	"io/ioutil"
	"log"
	"math"
	"regexp"
	"strings"

	"fyne.io/fyne"
	"fyne.io/fyne/canvas"
	"fyne.io/fyne/layout"
	"fyne.io/fyne/theme"
	"fyne.io/fyne/widget"
	"github.com/raedahgroup/dcrlibwallet"
	godcrApp "github.com/raedahgroup/godcr/app"
	"github.com/raedahgroup/godcr/fyne/widgets"
)

func (app *AppInterface) ShowCreateAndRestoreWalletPage() {
	app.Window.SetContent(app.createAndRestoreWalletPage())
	app.Window.CenterOnScreen()
	app.Window.Resize(fyne.NewSize(500, 500))
	app.Window.SetFixedSize(true)
	fyne.CurrentApp().Settings().SetTheme(theme.LightTheme())
	app.Window.ShowAndRun()
}

func (app *AppInterface) createAndRestoreWalletPage() fyne.CanvasObject {
	createWallet := widget.NewButtonWithIcon("Create a new Wallet", theme.ContentAddIcon(), func() {
		app.createSpendingPasswordPopup("")
	})

	restoreWallet := widget.NewButtonWithIcon("Restore an existing wallet", theme.ContentRedoIcon(), func() {
		app.Window.SetContent(app.restoreWalletPage())
	})

	image := canvas.NewImageFromFile("fyne/assets/decred.png")
	image.FillMode = canvas.ImageFillOriginal

	createAndRestoreButtons := widget.NewVBox(fyne.NewContainerWithLayout(
		layout.NewFixedGridLayout(restoreWallet.MinSize()), createWallet),
		restoreWallet)

	page := widget.NewVBox(
		image,
		widget.NewLabelWithStyle("Welcome to\nDecred Desktop Wallet", fyne.TextAlignLeading, fyne.TextStyle{Bold: true}),
		layout.NewSpacer(),
		widget.NewHBox(layout.NewSpacer(), createAndRestoreButtons, layout.NewSpacer()))

	return widget.NewHBox(layout.NewSpacer(), page, layout.NewSpacer())
}

func (app *AppInterface) passwordTab(popup *widget.PopUp, isPassword bool, seed string) fyne.CanvasObject {
	displayError := func(err error) {
		log.Println("could not generate seed", err.Error())
		newWindow := fyne.CurrentApp().NewWindow(godcrApp.DisplayName)
		newWindow.SetContent(widget.NewVBox(
			widget.NewLabelWithStyle(fmt.Sprintf("Could not generate seed, %s", err.Error()), fyne.TextAlignCenter, fyne.TextStyle{}),
			widget.NewButton("Close", func() { newWindow.Close() })))
	}

	errorLabel := canvas.NewText("Password do not match", color.RGBA{255, 0, 0, 225})
	errorLabel.TextSize = 10
	errorLabel.Hide()

	placeholder := "Spending Password"
	if !isPassword {
		placeholder = "Spending PIN"
	}

	password := widget.NewPasswordEntry()
	password.SetPlaceHolder(placeholder)
	confirmPassword := widget.NewPasswordEntry()
	confirmPassword.SetPlaceHolder(fmt.Sprintf("Confirm %s", placeholder))

	passwordLength := canvas.NewText("0", color.Black)
	passwordLength.TextSize = 10
	passwordLength.Alignment = fyne.TextAlignTrailing
	confirmPasswordLength := canvas.NewText("0", color.Black)
	confirmPasswordLength.TextSize = 10
	confirmPasswordLength.Alignment = fyne.TextAlignTrailing

	passwordStrength := widget.NewProgressBar()
	var createButton *widget.Button

	pinExpression, err := regexp.Compile("\\D")
	if err != nil {
		log.Println(err)
	}

	password.OnChanged = func(val string) {
		if !isPassword && len(val) > 0 && pinExpression.MatchString(val) {
			if len(val) == 1 {
				password.SetText("")
			} else {
				val = val[:password.CursorColumn-1] + val[password.CursorColumn:]
				//todo: using setText, cursor column count doesnt increase or reduce. Create an issue on this
				password.CursorColumn--
				password.SetText(val)
			}
			return
		}

		// check if password and confirm password matches only when the user fills confirmPassword textbox
		if confirmPassword.Text != "" {
			if confirmPassword.Text != password.Text {
				errorLabel.Show()
			} else {
				errorLabel.Hide()
			}
		}

		passwordLength.Text = fmt.Sprintf("%d", len(val))
		canvas.Refresh(passwordLength)
		strength := (shannonEntropy(val) / 4.0)
		passwordStrength.SetValue(strength)
	}

	confirmPassword.OnChanged = func(val string) {
		if !isPassword && len(val) > 0 && !pinExpression.MatchString(val) {
			if len(val) == 1 {
				confirmPassword.SetText("")
			} else {
				val = val[:confirmPassword.CursorColumn-1] + val[confirmPassword.CursorColumn:]
				confirmPassword.CursorColumn--
				confirmPassword.SetText(val)
			}
			return
		}

		confirmPasswordLength.Text = fmt.Sprintf("%d", len(val))
		canvas.Refresh(confirmPasswordLength)
		if password.Text != val {
			errorLabel.Show()
		} else if password.Text != "" && password.Text == confirmPassword.Text {
			errorLabel.Hide()
			createButton.Enable()
		}
	}

	createButton = widget.NewButton("Create", func() {
		if seed == "" {
			seed, err = dcrlibwallet.GenerateSeed()
			if err != nil {
				displayError(err)
				return
			}
		}

		if app.Wallet.CreateWallet(password.Text, seed) != nil {
			displayError(err)
			log.Println("could not create wallet", err.Error())
			return
		}
		popup.Hide()
		app.Window.SetFixedSize(false)
		app.MenuPage()
	})

	createButton.Disable()
	conceal, err := ioutil.ReadFile("fyne/assets/ic_conceal_24px.png")
	if err != nil {
		log.Fatalln(err)
	}
	reveal, err := ioutil.ReadFile("fyne/assets/ic_reveal_24px.png")
	if err != nil {
		log.Fatalln(err)
	}

	var passwordConceal *widget.Button
	passwordConceal = widget.NewButtonWithIcon("", fyne.NewStaticResource("", reveal), func() {
		if password.Password {
			passwordConceal.SetIcon(fyne.NewStaticResource("", conceal))
			password.Password = false
		} else {
			passwordConceal.SetIcon(fyne.NewStaticResource("", reveal))
			password.Password = true
		}
		// reveal texts
		password.SetText(password.Text)
	})

	var confirmPasswordConceal *widget.Button
	confirmPasswordConceal = widget.NewButtonWithIcon("", fyne.NewStaticResource("", reveal), func() {
		if confirmPassword.Password {
			confirmPasswordConceal.SetIcon(fyne.NewStaticResource("", conceal))
			confirmPassword.Password = false
		} else {
			confirmPasswordConceal.SetIcon(fyne.NewStaticResource("", reveal))
			confirmPassword.Password = true
		}
		// reveal texts
		confirmPassword.SetText(confirmPassword.Text)
	})

	return widget.NewVBox(
		widget.NewHBox(layout.NewSpacer(), fyne.NewContainerWithLayout(layout.NewFixedGridLayout(confirmPassword.MinSize()), password), passwordConceal, layout.NewSpacer()),
		passwordLength,
		widget.NewHBox(widget.NewLabel("Password Strength"), layout.NewSpacer(), passwordStrength),
		widget.NewHBox(layout.NewSpacer(), fyne.NewContainerWithLayout(layout.NewFixedGridLayout(confirmPassword.MinSize()), confirmPassword), confirmPasswordConceal, layout.NewSpacer()),
		confirmPasswordLength,
		widget.NewHBox(layout.NewSpacer(), widget.NewButton("Cancel", func() { popup.Hide() }), createButton),
		errorLabel)
}

func (app *AppInterface) restoreWalletPage() fyne.CanvasObject {
	app.Window.SetOnClosed(func() {
		app.Window = fyne.CurrentApp().NewWindow(godcrApp.DisplayName)
		app.Window.SetContent(app.createAndRestoreWalletPage())
		app.Window.CenterOnScreen()
		app.Window.Resize(fyne.NewSize(500, 500))
		app.Window.SetFixedSize(true)
		app.Window.Show()
	})

	var textbox = make([]*widget.Entry, 33)
	file, err := ioutil.ReadFile("fyne/assets/wordlist.txt")
	if err != nil {
		log.Fatalln("couldnt read wordlist", err)
	}

	wordlist := strings.Split(string(file), "\n")
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
		popup.Move(fyne.CurrentApp().Driver().AbsolutePositionForObject(textbox[textboxIndex]).Add(fyne.NewPos(0, textbox[textboxIndex].Size().Height)))
	}

	var restoreButton = widget.NewButton("Restore", func() {
		var seed string
		for i := 0; i < 32; i++ {
			seed += textbox[i].Text + " "
		}
		seed += textbox[32].Text

		if dcrlibwallet.VerifySeed(seed) {
			check, _ := ioutil.ReadFile("fyne/assets/ic_checkmark_64px.png")
			icon := canvas.NewImageFromResource(fyne.NewStaticResource("", check))
			icon.FillMode = canvas.ImageFillOriginal
			windowContent := app.Window.Content()
			if box, ok := windowContent.(*widget.Box); ok {
				box.Children = []fyne.CanvasObject{
					layout.NewSpacer(),
					icon,
					layout.NewSpacer(),
					widget.NewLabelWithStyle("Your wallet is successfully restored", fyne.TextAlignCenter, fyne.TextStyle{Bold: true}),
					widget.NewLabelWithStyle("Now create a spending password to protect your funds.", fyne.TextAlignCenter, fyne.TextStyle{}),
					widget.NewHBox(layout.NewSpacer(), widget.NewButton("Create a spending password", func() { app.createSpendingPasswordPopup(seed) }), layout.NewSpacer()),
					widgets.NewVSpacer(20)}
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

	textBoxContainer := widget.NewHBox(layout.NewSpacer(), horizontalTextBoxes.Children[0], layout.NewSpacer(), horizontalTextBoxes.Children[1], layout.NewSpacer(), horizontalTextBoxes.Children[2], layout.NewSpacer())
	buttonContainer := widget.NewHBox(layout.NewSpacer(), restoreButton, layout.NewSpacer())
	return widget.NewVBox(errorLabel, textBoxContainer, buttonContainer)
}

func (app *AppInterface) createSpendingPasswordPopup(seed string) {
	var popup *widget.PopUp

	popupContent := widget.NewVBox()
	popup = widget.NewModalPopUp(widget.NewVBox(
		widget.NewLabelWithStyle("Create a Spending Password", fyne.TextAlignLeading, fyne.TextStyle{Bold: true}), popupContent),
		app.Window.Canvas())

	tabMenu := widget.NewTabContainer(
		widget.NewTabItem("     Password    ", app.passwordTab(popup, true, seed)), widget.NewTabItem("       Pin        ", app.passwordTab(popup, false, seed)))

	popupContent.Children = []fyne.CanvasObject{tabMenu}
	widget.Refresh(popupContent)
}

func shannonEntropy(data string) (entropy float64) {
	if data == "" {
		return 0
	}
	for i := 0; i < 256; i++ {
		px := float64(strings.Count(data, string(byte(i)))) / float64(len(data))
		if px > 0 {
			entropy += -px * math.Log2(px)
		}
	}
	return entropy
}
