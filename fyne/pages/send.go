package pages

import (
	"fmt"
	"image/color"
	"io/ioutil"
	"log"

	"fyne.io/fyne"
	"fyne.io/fyne/canvas"
	"fyne.io/fyne/layout"
	"fyne.io/fyne/theme"
	"fyne.io/fyne/widget"
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
	send.fromAccountSelect.Options = name
	send.toAccountSelect.Options = fullStatus
}

func sendPage(wallet godcrApp.WalletMiddleware, window fyne.Window) fyne.CanvasObject {
	label := widget.NewLabelWithStyle("Sending Decred", fyne.TextAlignLeading, fyne.TextStyle{Bold: true, Italic: true})
	accountLabel := widget.NewLabelWithStyle("From:", fyne.TextAlignLeading, fyne.TextStyle{Bold: true})

	send.errorLabel = canvas.NewText("", color.RGBA{255, 0, 0, 0})
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

	button := widget.NewButton("Submit", func() {
		fmt.Println("Hello")
	})
	button.Hide()

	var name []string
	var fullStatus []string
	for _, account := range accounts {
		name = append(name, account.Name)
		fullStatus = append(fullStatus, account.String())
	}

	send.fromAccountSelect = widget.NewSelect(fullStatus, func(s string) {
		if button.Disabled() == true {
			button.Enable()
		}
	})
	send.toAccountSelect = widget.NewSelect(name, func(b string) {
		//	fmt.Println(b)
	})
	send.toAccountSelect.Hide()
	send.toAccountSelect.Hide()

	address := widget.NewEntry()
	address.SetPlaceHolder("Destination Address")
	amount := widget.NewEntry()
	amount.SetPlaceHolder("Amount")

	//place a menu toolbar to view options
	file, err := ioutil.ReadFile("fyne/pages/png/menu.png")
	if err != nil {
		log.Fatalln("could not read file menu.png", err.Error())
	}

	var popup *widget.PopUp
	menuIcon := widget.NewToolbar(widget.NewToolbarAction(fyne.NewStaticResource("menu", file), func() { popup.Show() }))

	transactionInfo := widget.NewHBox(address, menuIcon)

	popup = widget.NewModalPopUp(widget.NewVBox(
		widget.NewButton("Send between accounts", func() {
			address.Hide()
			send.toAccountSelect.Show()
			widget.Refresh(address)
			widget.Refresh(send.toAccountSelect)
			transactionInfo.Children = []fyne.CanvasObject{send.toAccountSelect, menuIcon}
			widget.Refresh(transactionInfo)
			popup.Hide()
		}),
		widget.NewButton("Send to others", func() {
			send.toAccountSelect.Hide()
			address.Show()
			widget.Refresh(address)
			widget.Refresh(send.toAccountSelect)
			transactionInfo.Children = []fyne.CanvasObject{address, menuIcon}
			widget.Refresh(transactionInfo)
			popup.Hide()
		})), window.Canvas())

	var spendUnconfirmed bool
	spendUnconfirmedCheck := widget.NewCheck("Spend Unconfirmed", func(check bool) {
		spendUnconfirmed = check
		fmt.Println(spendUnconfirmed)
	})

	//this is the button that shows information of on sendpage as a popup
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
		amount,
		fyne.NewContainerWithLayout(layout.NewFixedGridLayout(button.MinSize()), button),
		send.errorLabel,
	)
	return widget.NewHBox(widgets.NewHSpacer(10), output)
}
