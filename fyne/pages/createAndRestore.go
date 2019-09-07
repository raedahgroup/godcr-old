package pages

import (
	"context"
	"fmt"
	"image/color"
	"io/ioutil"
	"log"
	"math"
	"regexp"
	"strings"

	"fyne.io/fyne/canvas"

	"fyne.io/fyne/theme"

	"github.com/raedahgroup/dcrlibwallet"

	"fyne.io/fyne"
	"fyne.io/fyne/layout"
	"fyne.io/fyne/widget"
	"github.com/raedahgroup/godcr/app/wallet"
)

func ShowCreateAndRestoreWalletPage(wallet wallet.Wallet, window fyne.Window, ctx context.Context) {
	createWallet := widget.NewButtonWithIcon("Create a new Wallet", theme.ContentAddIcon(), func() {
		var popup *widget.PopUp

		popupContent := widget.NewVBox()
		popup = widget.NewModalPopUp(widget.NewVBox(
			widget.NewLabelWithStyle("Create a Spending Password", fyne.TextAlignLeading, fyne.TextStyle{Bold: true}), popupContent),
			window.Canvas())

		tabMenu := widget.NewTabContainer(
			widget.NewTabItem("     Password    ", passwordTab(wallet, popup, true)), widget.NewTabItem("       Pin        ", passwordTab(wallet, popup, false)))

		popupContent.Children = []fyne.CanvasObject{tabMenu}
		widget.Refresh(popupContent)
	})

	restoreWallet := widget.NewButtonWithIcon("Restore an existing wallet", theme.ContentRedoIcon(), func() {
	})

	image := canvas.NewImageFromFile("fyne/decred.png")
	image.SetMinSize(fyne.NewSize(100, 50))
	page := widget.NewVBox(
		image,
		widget.NewLabelWithStyle("Welcome to\nDecred Desktop Wallet", fyne.TextAlignLeading, fyne.TextStyle{Bold: true}),
		layout.NewSpacer(),
		fyne.NewContainerWithLayout(layout.NewFixedGridLayout(restoreWallet.MinSize()), createWallet),
		restoreWallet)

	window.SetContent(widget.NewHBox(layout.NewSpacer(), page, layout.NewSpacer()))
	window.CenterOnScreen()
	window.Resize(fyne.NewSize(500, 500))
	window.SetFixedSize(true)
	window.ShowAndRun()
}

func passwordTab(wallet wallet.Wallet, popup *widget.PopUp, isPassword bool) fyne.CanvasObject {
	errorLabel := canvas.NewText("Password do not match", color.RGBA{255, 0, 0, 225})
	errorLabel.TextSize = 10
	errorLabel.Hide()

	password := widget.NewPasswordEntry()
	password.SetPlaceHolder("Password")
	confirmPassword := widget.NewPasswordEntry()
	confirmPassword.SetPlaceHolder("Confirm Password")

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
				//todo: using setText, cursor column count doesnt increase or reduce
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
		seed, err := dcrlibwallet.GenerateSeed()
		if err != nil {
			log.Println("could not generate seed", err.Error())
			return
		}
		if wallet.CreateWallet(password.Text, seed) != nil {
			log.Println("could not create wallet", err.Error())
		}
		// move to overview page
	})

	createButton.Disable()
	conceal, _ := ioutil.ReadFile("/Users/macbook/Downloads/ic_conceal_24px.png")
	reveal, _ := ioutil.ReadFile("/Users/macbook/Downloads/ic_reveal_24px.png")

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
		widget.NewHBox(layout.NewSpacer(), fyne.NewContainerWithLayout(layout.NewFixedGridLayout(confirmPassword.MinSize()), password), passwordConceal, layout.NewSpacer()), //fyne.NewContainerWithLayout(layout.NewFixedGridLayout(passwordConceal.MinSize()), passwordConceal), layout.NewSpacer()),
		passwordLength,
		widget.NewHBox(widget.NewLabel("Password Strength"), layout.NewSpacer(), passwordStrength),                                                                                         //passwordWeak, passwordStrong),
		widget.NewHBox(layout.NewSpacer(), fyne.NewContainerWithLayout(layout.NewFixedGridLayout(confirmPassword.MinSize()), confirmPassword), confirmPasswordConceal, layout.NewSpacer()), //, confirmPasswordConceal, layout.NewSpacer()),
		confirmPasswordLength,
		widget.NewHBox(layout.NewSpacer(), widget.NewButton("Cancel", func() { popup.Hide() }), createButton),
		errorLabel)
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

func restoreWalletPage(wallet wallet.Wallet, window fyne.Window) fyne.CanvasObject {
	var textbox = make([]*widget.Entry, 33)
	file, err := ioutil.ReadFile("fyne/wordlist.txt")
	if err != nil {
		log.Fatalln("couldnt read wordlist")
	}
	wordlist := strings.Split(string(file), "\n")
	var horizontalTextBoxes = widget.NewHBox()

	containString := func(start, stop, textboxIndex int, val string) {
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
		if len(menuItem) == 0 {
			return
		}

		popup = widget.NewPopUpMenu(fyne.NewMenu("", menuItem...), window.Canvas())
		popup.Move(fyne.CurrentApp().Driver().AbsolutePositionForObject(textbox[textboxIndex]).Add(fyne.NewPos(0, textbox[textboxIndex].Size().Height)))
	}

	for i := 0; i < 33; i++ {
		textboxIndex := i
		textbox[textboxIndex] = widget.NewEntry()
		textbox[textboxIndex].SetPlaceHolder(fmt.Sprintf("Word %d", i+1))
		textbox[textboxIndex].OnChanged = func(word string) {
			containString(0, len(wordlist), textboxIndex, word)
		}
	}

	for i := 0; i < 3; i++ {
		vertical := widget.NewVBox()
		for k := i; k < 33; k += 3 {
			vertical.Append(textbox[k])
		}
		horizontalTextBoxes.Append(vertical)
	}

	return widget.NewHBox(layout.NewSpacer(), horizontalTextBoxes.Children[0], layout.NewSpacer(), horizontalTextBoxes.Children[1], layout.NewSpacer(), horizontalTextBoxes.Children[2], layout.NewSpacer())
}

func createNewWalletPage(window fyne.Window) fyne.CanvasObject {
	seed, err := dcrlibwallet.GenerateSeed()
	if err != nil {
		return widget.NewLabel(fmt.Sprintf("error generating seed, %s", err.Error()))
	}
	wordings := strings.Split(seed, " ")

	var horizontalTextBoxes = widget.NewHBox()
	for i := 0; i < 3; i++ {
		vertical := widget.NewVBox()
		horizontalTextBoxes.Append(vertical)
		vertical = widget.NewVBox()
		for k := i; k < 33; k += 3 {
			vertical.Append(widget.NewLabel(fmt.Sprintf("%d. %s", k+1, wordings[k])))
		}
		horizontalTextBoxes.Append(vertical)

		if i == 2 {
			break
		}

		for j := 0; j < 11; j++ {
			vertical.Append(layout.NewSpacer())
		}
	}
	horizontalTextBoxes.Append(layout.NewSpacer())

	container := widget.NewVBox()
	container.Append(widget.NewHBox(layout.NewSpacer(), horizontalTextBoxes, layout.NewSpacer()))
	container.Append(layout.NewSpacer())
	container.Append(widget.NewHBox(layout.NewSpacer(), widget.NewButton("Copy Seed Phrase", func() {
		clippy := window.Clipboard()
		clippy.SetContent(seed)
		// todo: After copying seeds, user should verify seed also.
	}), layout.NewSpacer()))
	container.Append(layout.NewSpacer())

	return container
}
