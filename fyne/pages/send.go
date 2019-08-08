package pages

import (
	"fmt"
	"image/color"
	"io/ioutil"
	"log"
	"strconv"
	"strings"

	"fyne.io/fyne"
	"fyne.io/fyne/canvas"
	"fyne.io/fyne/layout"
	"fyne.io/fyne/theme"
	"fyne.io/fyne/widget"
	"github.com/raedahgroup/dcrlibwallet/txhelper"
	godcrApp "github.com/raedahgroup/godcr/app"
	"github.com/raedahgroup/godcr/app/walletcore"
	"github.com/raedahgroup/godcr/fyne/widgets"
)

type sendPageData struct {
	fromAccountSelect *widget.Select
	toAccountSelect   *widget.Select
	errorLabel        *canvas.Text
}

//both send page and send page update would be in a function
var send sendPageData

func sendPageUpdates(wallet godcrApp.WalletMiddleware) {
	accounts, _ := wallet.AccountsOverview(walletcore.DefaultRequiredConfirmations)

	var name []string
	var fullStatus []string
	for _, account := range accounts {
		name = append(name, account.Name)
		fullStatus = append(fullStatus, account.String())
	}
	send.fromAccountSelect.Options = fullStatus
	send.toAccountSelect.Options = name
}

func sendPage(wallet godcrApp.WalletMiddleware, window fyne.Window) fyne.CanvasObject {
	label := widget.NewLabelWithStyle("Sending Decred", fyne.TextAlignLeading, fyne.TextStyle{Bold: true, Italic: true})
	accountLabel := widget.NewLabelWithStyle("From:", fyne.TextAlignLeading, fyne.TextStyle{Bold: true})

	send.errorLabel = canvas.NewText("", color.RGBA{255, 0, 0, menu.alphaTheme})
	send.errorLabel.Alignment = fyne.TextAlignCenter
	send.errorLabel.TextStyle = fyne.TextStyle{Bold: true}
	send.errorLabel.Hide()

	accounts, err := wallet.AccountsOverview(walletcore.DefaultRequiredConfirmations)
	if err != nil {
		send.errorLabel.Text = "Could not retrieve account information"
		//todo: log to file
		fmt.Println(err.Error())
		send.errorLabel.Show()
	}

	var button *widget.Button
	var name []string
	var fullStatus []string
	var requiredConfirmations int32 = walletcore.DefaultRequiredConfirmations

	for _, account := range accounts {
		name = append(name, account.Name)
		fullStatus = append(fullStatus, account.String())
	}

	send.fromAccountSelect = widget.NewSelect(fullStatus, func(s string) {
	})

	send.toAccountSelect = widget.NewSelect(name, func(b string) {
	})

	address := widget.NewEntry()
	address.SetPlaceHolder("Destination Address")
	amount := widget.NewEntry()
	amount.SetPlaceHolder("Amount")

	file, err := ioutil.ReadFile("fyne/pages/png/max.png")
	if err != nil {
		log.Fatalln("could not read file max.png", err.Error())
	}

	//on clicking, it set amount text to the total amount in wallet
	sendMaxButton := widget.NewToolbar(widget.NewToolbarAction(fyne.NewStaticResource("max", file), func() {
		splittedVals := strings.Split(send.fromAccountSelect.Selected, " ")
		accoutNo, _ := wallet.AccountNumber(splittedVals[0])
		bal, _ := wallet.AccountBalance(accoutNo, requiredConfirmations)
		amount.SetText(strconv.FormatFloat(bal.Spendable.ToCoin(), 'f', -1, 64))
	}))

	//place a menu toolbar to view options
	file, err = ioutil.ReadFile("fyne/pages/png/menu.png")
	if err != nil {
		log.Fatalln("could not read file menu.png", err.Error())
	}

	var transactionTypePopup *widget.PopUp
	menuIcon := widget.NewToolbar(widget.NewToolbarAction(fyne.NewStaticResource("menu", file), func() { transactionTypePopup.Show() }))

	transactionInfo := widget.NewHBox(address, menuIcon)

	var sendToOthers = true

	//this popup enables users choose between sending to other addresses or sending between accounts
	transactionTypePopup = widget.NewModalPopUp(widget.NewVBox(
		widget.NewButton("Send between accounts", func() {
			sendToOthers = false

			address.Hide()
			send.toAccountSelect.Show()
			widget.Refresh(address)
			widget.Refresh(send.toAccountSelect)

			transactionInfo.Children = []fyne.CanvasObject{send.toAccountSelect, menuIcon}
			widget.Refresh(transactionInfo)
			transactionTypePopup.Hide()
		}),
		widget.NewButton("Send to others", func() {
			sendToOthers = true

			send.toAccountSelect.Hide()
			address.Show()
			widget.Refresh(address)
			widget.Refresh(send.toAccountSelect)

			transactionInfo.Children = []fyne.CanvasObject{address, menuIcon}
			widget.Refresh(transactionInfo)
			transactionTypePopup.Hide()
		})), window.Canvas())

	spendUnconfirmedCheck := widget.NewCheck("Spend Unconfirmed", func(check bool) {
		if check {
			requiredConfirmations = 0
		} else {
			requiredConfirmations = walletcore.DefaultRequiredConfirmations
		}
	})

	password := widget.NewPasswordEntry()
	password.SetPlaceHolder("Password")
	var getPasswordPopup *widget.PopUp

	getPasswordPopup = widget.NewModalPopUp(
		widget.NewVBox(
			password,
			widget.NewHBox(
				layout.NewSpacer(),
				widget.NewButton("Cancel", func() {
					getPasswordPopup.Hide()
					password.SetText("")
				}),

				//when user click on the accept button, create an inifinite blocking popup
				widget.NewButton("Accept", func() {
					hash, err := sendTransaction(window, wallet, sendToOthers, amount.Text, address.Text, password.Text, requiredConfirmations)

					//clear password
					password.SetText("")
					if err != nil {
						send.errorLabel.Text = "error performing transaction " + err.Error()
						canvas.Refresh(send.errorLabel)
						//todo: log to file
						fmt.Println(err.Error())
						send.errorLabel.Show()
						return
					}

					getPasswordPopup.Hide()

					confPopup := confirmationPopUp(window, err, hash)
					confPopup.Show()
				}),
			)), window.Canvas())

	button = widget.NewButton("submit", func() {
		getPasswordPopup.Show()
	})

	//this is the button that shows information of sendpage as a popup
	infoIcon := widget.NewToolbar(widget.NewToolbarAction(theme.InfoIcon(), func() {
		var button *widget.Button
		var popUp *widget.PopUp

		button = widget.NewButton("Got it", func() {
			popUp.Hide()
		})

		header := widget.NewLabelWithStyle("Send DCR", fyne.TextAlignLeading, fyne.TextStyle{Bold: true})
		info := widget.NewLabelWithStyle("Input the destination wallet address and the amount to send funds", fyne.TextAlignCenter, fyne.TextStyle{Monospace: true})

		data := widget.NewVBox(header, widgets.NewVSpacer(10), info,
			widget.NewHBox(layout.NewSpacer(), button))

		popUp = widget.NewModalPopUp(data, window.Canvas())
	}))

	output := widget.NewVBox(
		widget.NewHBox(label, infoIcon),
		widget.NewHBox(accountLabel, send.fromAccountSelect),
		spendUnconfirmedCheck,
		transactionInfo,
		widget.NewHBox(amount, sendMaxButton),
		fyne.NewContainerWithLayout(layout.NewFixedGridLayout(button.MinSize()), button),
		send.errorLabel,
	)
	return widget.NewHBox(widgets.NewHSpacer(10), output)
}

func sendTransaction(window fyne.Window, wallet godcrApp.WalletMiddleware, sendToOthers bool, amount, address, passphrase string, requiredConfirmations int32) (string, error) {
	progress := widget.NewProgressBarInfinite()
	blockingPopup := widget.NewModalPopUp(progress, window.Canvas())
	blockingPopup.Show()

	var hash string
	var txerr error
	var generatedAddress = address

	if !sendToOthers {
		walletNo, err := wallet.AccountNumber(send.toAccountSelect.Selected)
		if err != nil {
			send.errorLabel.Text = "error getting receiving wallet account " + err.Error()
			fmt.Println(err.Error())
			canvas.Refresh(send.errorLabel)
			send.errorLabel.Show()

			blockingPopup.Hide()
			return "", err
		}

		generatedAddress, err = wallet.GenerateNewAddress(walletNo)
		if err != nil {
			send.errorLabel.Text = "error generating address " + err.Error()
			fmt.Println(err.Error())
			canvas.Refresh(send.errorLabel)
			send.errorLabel.Show()

			blockingPopup.Hide()
			return "", err
		}
	}

	formatAmount, err := strconv.ParseFloat(amount, 64)
	if err != nil {
		send.errorLabel.Text = "error formatting amount " + err.Error()
		fmt.Println(err.Error())
		canvas.Refresh(send.errorLabel)
		send.errorLabel.Show()

		blockingPopup.Hide()
		return "", err
	}

	sendDestinations := []txhelper.TransactionDestination{
		{Address: generatedAddress,
			Amount: formatAmount}}

	//split to get name before space
	splittedVals := strings.Split(send.fromAccountSelect.Selected, " ")
	fmt.Println(splittedVals)
	walletNo, err := wallet.AccountNumber(splittedVals[0])
	if err != nil {
		send.errorLabel.Text = "error getting sending wallet account " + err.Error()
		fmt.Println(err.Error())
		canvas.Refresh(send.errorLabel)
		send.errorLabel.Show()

		blockingPopup.Hide()
		return "", err
	}

	hash, txerr = wallet.SendFromAccount(walletNo, requiredConfirmations, sendDestinations, passphrase)

	blockingPopup.Hide()

	return hash, txerr
}

func confirmationPopUp(window fyne.Window, err error, text string) fyne.CanvasObject {
	errorLabel := canvas.NewText("Transaction was successful", color.RGBA{11, 156, 49, menu.alphaTheme})
	errorLabel.TextStyle = fyne.TextStyle{Bold: true}
	errorLabel.Alignment = fyne.TextAlignCenter

	var popUp *widget.PopUp
	if err != nil {
		errorLabel.Text = "Transaction was not successful"
		canvas.Refresh(errorLabel)

		popUp = widget.NewModalPopUp(widget.NewVBox(
			errorLabel,
			widget.NewLabel(err.Error()),
			widget.NewHBox(layout.NewSpacer(), widget.NewButton("Ok", func() { popUp.Hide() }), layout.NewSpacer())),

			window.Canvas())
		return popUp
	}

	popUp = widget.NewModalPopUp(widget.NewVBox(
		errorLabel,
		widget.NewLabel("Hash: "+text),
		widget.NewHBox(layout.NewSpacer(), widget.NewButton("Ok", func() { popUp.Hide() }), layout.NewSpacer())),
		window.Canvas())

	return popUp
}
