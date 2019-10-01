package pages

import (
	"fmt"

	"github.com/gdamore/tcell"
	"github.com/raedahgroup/dcrlibwallet"
	"github.com/raedahgroup/godcr/terminal/helpers"
	"github.com/raedahgroup/godcr/terminal/primitives"
	"github.com/rivo/tview"
	"github.com/skip2/go-qrcode"
)

const receivingDecredHint = "Each time you request payment, a new address is generated to protect your privacy."

func receivePage() tview.Primitive {
	body := tview.NewFlex().SetDirection(tview.FlexRow)

	body.AddItem(primitives.NewLeftAlignedTextView("Receiving Decred"), 1, 1, false)
	receivingDecredHintTextView := primitives.NewLeftAlignedTextView(receivingDecredHint).
		SetTextColor(helpers.HintTextColor)
	body.AddItem(receivingDecredHintTextView, 2, 1, false)

	generateAdressFunc := func(formButton tview.Primitive) {
		newAddressFlex := tview.NewFlex().
			AddItem(primitives.NewLeftAlignedTextView("You can also manually generate a").SetTextColor(helpers.HintTextColor), 33, 0, false).
			AddItem(formButton, 0, 1, false)
		body.AddItem(newAddressFlex, 3, 1, false)
	}

	accounts, err := commonPageData.wallet.GetAccountsRaw(dcrlibwallet.DefaultRequiredConfirmations)
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

	generateAndDisplayAddress := func(accountNumber int32, newAddress bool) {
		// clear previously generated address or displayed error before generating new one
		body.RemoveItem(qrCodeTextView)
		body.RemoveItem(addressTextView)
		body.RemoveItem(errorMessageTextView)

		var address string
		var qr *qrcode.QRCode
		var err error

		if newAddress {
			address, qr, err = generateNewAddressAndQrCode(accountNumber)
		} else {
			address, qr, err = generateAddressAndQrCode(accountNumber)
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

	accountNumbers := make([]int32, len(accounts.Acc))
	accountNames := make([]string, len(accounts.Acc))
	for index, account := range accounts.Acc {
		accountNames[index] = account.Name
		accountNumbers[index] = account.Number
	}

	formButton := primitives.NewForm(false)
	formButton.SetBorderPadding(0, 0, 0, 0)
	formButton.SetCancelFunc(commonPageData.clearAllPageContent)

	if len(accounts.Acc) == 1 {
		singleAccountTextView := primitives.NewLeftAlignedTextView(fmt.Sprintf("Source Account: %s", accounts.Acc[0].Name))
		singleAccountTextView.SetTextColor(helpers.DecredLightBlueColor)
		singleAccountTextView.SetDoneFunc(func(key tcell.Key) {
			if key == tcell.KeyEscape {
				commonPageData.clearAllPageContent()
			}
		})
		formButton.AddButton("new address.", func() {
			generateAndDisplayAddress(accounts.Acc[0].Number, true)
		})

		generateAdressFunc(formButton)
		body.AddItem(singleAccountTextView, 2, 1, true)

		singleAccountTextView.SetInputCapture(func(event *tcell.EventKey) *tcell.EventKey {
			if event.Key() == tcell.KeyTab {
				commonPageData.app.SetFocus(formButton)
				return nil
			}
			return event
		})

		formButton.GetButton(0).SetInputCapture(func(event *tcell.EventKey) *tcell.EventKey {
			if event.Key() == tcell.KeyTab {
				commonPageData.app.SetFocus(singleAccountTextView)
				return nil
			}

			return event
		})

		commonPageData.hintTextView.SetText("TIP: Move around with TAB. ESC to return to navigation menu")
	} else {
		accountNumbers := make([]int32, len(accounts.Acc))
		accountNames := make([]string, len(accounts.Acc))
		for index, account := range accounts.Acc {
			accountNames[index] = account.Name
			accountNumbers[index] = account.Number
		}

		formDropdown := primitives.NewForm(false)
		formDropdown.SetBorderPadding(0, 0, 0, 0)
		formDropdown.SetLabelColor(helpers.DecredLightBlueColor)
		formDropdown.SetCancelFunc(commonPageData.clearAllPageContent)

		var accountNumber int32
		formDropdown.AddDropDown("Source Account: ", accountNames, 0, func(option string, optionIndex int) {
			accountNumber = accountNumbers[optionIndex]
			generateAndDisplayAddress(accountNumber, false)
		})

		formButton.AddButton("new address.", func() {
			generateAndDisplayAddress(accountNumber, true)
		})

		generateAdressFunc(formButton)

		body.AddItem(formDropdown, 2, 0, true)
		commonPageData.hintTextView.SetText("TIP: Select Preferred Account and hit ENTER to generate Address, \n" +
			"Move around with TAB, ESC to return to navigation menu")

		formDropdown.GetFormItemBox(0).SetInputCapture(func(event *tcell.EventKey) *tcell.EventKey {
			if event.Key() == tcell.KeyTab {
				commonPageData.app.SetFocus(formButton)
				return nil
			}

			return event
		})

		formButton.GetButton(0).SetInputCapture(func(event *tcell.EventKey) *tcell.EventKey {
			if event.Key() == tcell.KeyTab {
				commonPageData.app.SetFocus(formDropdown)
				return nil
			}

			return event
		})
	}

	// always generate and display address for the first account, even if there are multiple accounts
	generateAndDisplayAddress(accounts.Acc[0].Number, false)

	commonPageData.app.SetFocus(body)
	return body
}

func generateAddressAndQrCode(accountNumber int32) (string, *qrcode.QRCode, error) {
	generatedAddress, err := commonPageData.wallet.CurrentAddress(accountNumber)
	if err != nil {
		return "", nil, err
	}

	qrCode, err := generateQrcode(generatedAddress)
	if err != nil {
		return "", nil, err
	}

	return generatedAddress, qrCode, nil
}

func generateNewAddressAndQrCode(accountNumber int32) (string, *qrcode.QRCode, error) {
	generatedAddress, err := commonPageData.wallet.NextAddress(accountNumber)
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
