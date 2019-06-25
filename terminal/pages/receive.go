package pages

import (
	"fmt"

	"github.com/gdamore/tcell"
	"github.com/raedahgroup/godcr/app/walletcore"
	"github.com/raedahgroup/godcr/terminal/helpers"
	"github.com/raedahgroup/godcr/terminal/primitives"
	"github.com/rivo/tview"
	qrcode "github.com/skip2/go-qrcode"
)

func receivePage(wallet walletcore.Wallet, hintTextView *primitives.TextView, setFocus func(p tview.Primitive) *tview.Application, clearFocus func()) tview.Primitive {
	body := tview.NewFlex().SetDirection(tview.FlexRow)

	body.AddItem(primitives.NewLeftAlignedTextView("Receiving Decred"), 1, 1, false)
	receivingDecredHintTextView := primitives.NewLeftAlignedTextView(walletcore.ReceivingDecredHint).
		SetTextColor(helpers.HintTextColor)
	body.AddItem(receivingDecredHintTextView, 2, 1, false)

	generateAdressFunc := func(formButton tview.Primitive) {
		newAddressFlex := tview.NewFlex().
			AddItem(primitives.NewLeftAlignedTextView("You can also manually generate a").SetTextColor(helpers.HintTextColor), 33, 0, false).
			AddItem(formButton, 0, 1, false)
		body.AddItem(newAddressFlex, 3, 1, false)
	}

	accounts, err := wallet.AccountsOverview(walletcore.DefaultRequiredConfirmations)
	if err != nil {
		return body.AddItem(primitives.NewLeftAlignedTextView(fmt.Sprintf("Error: %s", err.Error())), 0, 1, false)
	}

	errorMessageTextView := primitives.WordWrappedTextView("").
		SetTextColor(helpers.DecredOrangeColor)

	displayErrorMessage := func(message string) {
		body.RemoveItem(errorMessageTextView)
		errorMessageTextView.SetText(message)
		body.AddItem(errorMessageTextView, 2, 0, false)
	}

	qrCodeTextView := primitives.NewCenterAlignedTextView("")
	addressTextView := primitives.NewCenterAlignedTextView("").
		SetTextColor(helpers.DecredLightBlueColor)

	generateAndDisplayAddress := func(accountNumber uint32, newAddress bool) {
		// clear previously generated address or displayed error before generating new one
		body.RemoveItem(qrCodeTextView)
		body.RemoveItem(addressTextView)
		body.RemoveItem(errorMessageTextView)

		var address string
		var qr *qrcode.QRCode
		var err error

		if newAddress {
			address, qr, err = generateNewAddressAndQrCode(wallet, accountNumber)
		} else {
			address, qr, err = generateAddressAndQrCode(wallet, accountNumber)
		}

		if err != nil {
			errorText := fmt.Sprintf("Error: %s", err.Error())
			displayErrorMessage(errorText)
			return
		}

		qrCodeTextView.SetText(qr.ToSmallString(false))
		addressTextView.SetText(address)

		body.AddItem(addressTextView, 2, 0, true)
		body.AddItem(qrCodeTextView, 0, 1, true)
	}

	accountNumbers := make([]uint32, len(accounts))
	accountNames := make([]string, len(accounts))
	for index, account := range accounts {
		accountNames[index] = account.Name
		accountNumbers[index] = account.Number
	}

	formButton := primitives.NewForm(false)
	formButton.SetBorderPadding(0, 0, 0, 0)
	formButton.SetCancelFunc(clearFocus)

	if len(accounts) == 1 {
		singleAccountTextView := primitives.NewLeftAlignedTextView(fmt.Sprintf("Source Account: %s", accounts[0].Name))
		singleAccountTextView.SetTextColor(helpers.DecredLightBlueColor)
		singleAccountTextView.SetDoneFunc(func(key tcell.Key) {
			if key == tcell.KeyEscape {
				clearFocus()
			}
		})
		formButton.AddButton("new address.", func() {
			generateAndDisplayAddress(accounts[0].Number, true)
		})

		generateAdressFunc(formButton)
		body.AddItem(singleAccountTextView, 2, 1, true)

		singleAccountTextView.SetInputCapture(func(event *tcell.EventKey) *tcell.EventKey {
			if event.Key() == tcell.KeyTab {
				setFocus(formButton)
				return nil
			}

			return event
		})

		formButton.GetButton(0).SetInputCapture(func(event *tcell.EventKey) *tcell.EventKey {
			if event.Key() == tcell.KeyTab {
				setFocus(singleAccountTextView)
				return nil
			}

			return event
		})

		hintTextView.SetText("TIP: Move around with TAB. ESC to return to navigation menu")
	} else {
		accountNumbers := make([]uint32, len(accounts))
		accountNames := make([]string, len(accounts))
		for index, account := range accounts {
			accountNames[index] = account.Name
			accountNumbers[index] = account.Number
		}

		formDropdown := primitives.NewForm(false)
		formDropdown.SetBorderPadding(0, 0, 0, 0)
		formDropdown.SetLabelColor(helpers.DecredLightBlueColor)
		formDropdown.SetCancelFunc(clearFocus)

		var accountNumber uint32
		formDropdown.AddDropDown("Source Account: ", accountNames, 0, func(option string, optionIndex int) {
			accountNumber = accountNumbers[optionIndex]
			generateAndDisplayAddress(accountNumber, false)
		})

		formButton.AddButton("NEW ADDRESS.", func() {
			generateAndDisplayAddress(accountNumber, true)
		})

		generateAdressFunc(formButton)

		body.AddItem(formDropdown, 2, 0, true)
		hintTextView.SetText("TIP: Select Prefered Account and hit ENTER to generate Address, \nMove around with TAB, ESC to return to navigation menu")

		formDropdown.GetFormItemBox(0).SetInputCapture(func(event *tcell.EventKey) *tcell.EventKey {
			if event.Key() == tcell.KeyTab {
				setFocus(formButton)
				return nil
			}

			return event
		})

		formButton.GetButton(0).SetInputCapture(func(event *tcell.EventKey) *tcell.EventKey {
			if event.Key() == tcell.KeyTab {
				setFocus(formDropdown)
				return nil
			}

			return event
		})
	}

	// always generate and display address for the first account, even if there are multiple accounts
	generateAndDisplayAddress(accounts[0].Number, false)

	setFocus(body)
	return body
}

func generateAddressAndQrCode(wallet walletcore.Wallet, accountNumber uint32) (string, *qrcode.QRCode, error) {
	generatedAddress, err := wallet.ReceiveAddress(accountNumber)
	if err != nil {
		return "", nil, err
	}

	qrCode, err := generateQrcode(generatedAddress)
	if err != nil {
		return "", nil, err
	}

	return generatedAddress, qrCode, nil
}

func generateNewAddressAndQrCode(wallet walletcore.Wallet, accountNumber uint32) (string, *qrcode.QRCode, error) {
	generatedAddress, err := wallet.GenerateNewAddress(accountNumber)
	if err != nil {
		return "", nil, err
	}

	qrCode, err := generateQrcode(generatedAddress)
	if err != nil {
		return "", nil, err
	}

	return generatedAddress, qrCode, nil
}

func generateQrcode(generatedAddress string) (*qrcode.QRCode, error) {
	// generate qrcode
	qr, err := qrcode.New(generatedAddress, qrcode.Medium)
	if err != nil {
		return nil, err
	}

	return qr, nil
}
