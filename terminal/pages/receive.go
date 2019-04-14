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
	body.AddItem(receivingDecredHintTextView, 3, 1, false)

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

	generateAndDisplayAddress := func(accountNumber uint32) {
		// clear previously generated address or displayed error before generating new one
		body.RemoveItem(qrCodeTextView)
		body.RemoveItem(addressTextView)
		body.RemoveItem(errorMessageTextView)

		address, qr, err := generateAddressAndQrCode(wallet, accountNumber)
		if err != nil {
			errorText := fmt.Sprintf("Error: %s", err.Error())
			displayErrorMessage(errorText)
			return
		}

		qrCodeTextView.SetText(qr.ToSmallString(false))
		addressTextView.SetText(address)

		body.AddItem(qrCodeTextView, 19, 0, true)
		body.AddItem(addressTextView, 0, 1, true)
	}

	accountNumbers := make([]uint32, len(accounts))
	accountNames := make([]string, len(accounts))
	for index, account := range accounts {
		accountNames[index] = account.Name
		accountNumbers[index] = account.Number
	}

	if len(accounts) == 1 {
		singleAccountTextView := primitives.NewLeftAlignedTextView(fmt.Sprintf("Source Account: %s", accounts[0].Name))
		singleAccountTextView.SetTextColor(helpers.DecredLightBlueColor)
		singleAccountTextView.SetDoneFunc(func(key tcell.Key) {
			if key == tcell.KeyEscape {
				clearFocus()
			}
		})

		body.AddItem(singleAccountTextView, 2, 1, true)
		hintTextView.SetText("TIP: ESC to return to navigation menu")
	} else {
		accountNumbers := make([]uint32, len(accounts))
		accountNames := make([]string, len(accounts))
		for index, account := range accounts {
			accountNames[index] = account.Name
			accountNumbers[index] = account.Number
		}

		form := primitives.NewForm(false)
		form.SetBorderPadding(0, 0, 0, 0)
		form.SetLabelColor(helpers.DecredLightBlueColor)
		form.AddDropDown("Source Account: ", accountNames, 0, func(option string, optionIndex int) {
			accountNumber := accountNumbers[optionIndex]
			generateAndDisplayAddress(accountNumber)
		})

		form.SetCancelFunc(clearFocus)

		body.AddItem(form, 2, 0, true)
		hintTextView.SetText("TIP: Select Prefered Account and hit ENTER to generate Address. ESC to return to navigation menu")
	}

	// always generate and display address for the first account, even if there are multiple accounts
	generateAndDisplayAddress(accounts[0].Number)

	setFocus(body)
	return body
}

func generateAddressAndQrCode(wallet walletcore.Wallet, accountNumber uint32) (string, *qrcode.QRCode, error) {
	generatedAddress, err := wallet.ReceiveAddress(accountNumber)
	if err != nil {
		return "", nil, err
	}

	// generate qrcode
	qr, err := qrcode.New(generatedAddress, qrcode.Medium)
	if err != nil {
		return "", nil, err
	}

	return generatedAddress, qr, nil
}
