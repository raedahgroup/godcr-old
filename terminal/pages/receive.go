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

	body.AddItem(primitives.NewLeftAlignedTextView("Receiving Decred"), 2, 1, false)
	body.AddItem(primitives.NewLeftAlignedTextView("Each time you request a payment, a new address is created to protect your privacy.").SetTextColor(helpers.HintTextColor), 3, 1, false)

	accounts, err := wallet.AccountsOverview(walletcore.DefaultRequiredConfirmations)
	if err != nil {
		return body.AddItem(primitives.NewLeftAlignedTextView(fmt.Sprintf("Error: %s", err.Error())), 0, 1, false)
	}

	outputMessageTextView := primitives.WordWrappedTextView("")
	outputMessageTextView.SetTextColor(helpers.DecredOrangeColor)

	displayErrorMessage := func(message string) {
		body.RemoveItem(outputMessageTextView)
		outputMessageTextView.SetText(message)
		body.AddItem(outputMessageTextView, 2, 0, false)
	}

	qrCodeTextView := primitives.NewCenterAlignedTextView("")
	addressTextView := primitives.NewCenterAlignedTextView("")
	addressTextView.SetTextColor(helpers.DecredLightBlueColor)

	displayOutput := func (qrCode, address string) {
		body.RemoveItem(qrCodeTextView)
		body.RemoveItem(addressTextView)
		body.AddItem(qrCodeTextView.SetText(qrCode), 19, 0, true)
		body.AddItem(addressTextView.SetText(address), 0, 1, true)
	}

	accountNumbers := make([]uint32, len(accounts))
	accountNames := make([]string, len(accounts))
	for index, account := range accounts {
		accountNames[index] = account.Name
		accountNumbers[index] = account.Number
	}

	if len(accounts) != 1 {
		address, qr, err := generateAddress(wallet, accounts[0].Number)
		if err != nil {
			errorText := fmt.Sprintf("Error: %s", err.Error())
			displayErrorMessage(errorText)
			return nil
		}

		body.AddItem(primitives.NewLeftAlignedTextView(fmt.Sprintf("Source Account: %s", accounts[0].Name)).SetTextColor(helpers.DecredLightBlueColor).SetDoneFunc(func(key tcell.Key) {
			if key == tcell.KeyEscape {
				clearFocus()
			}
		}), 2, 1, true)
		qrCode := fmt.Sprintf(qr.ToSmallString(false))
		displayOutput(qrCode, address)
	} else {
		// receive form 
		form := primitives.NewForm()
		form.SetBorderPadding(0, 0, 0, 0)
		form.SetHorizontal(true).
		SetLabelColor(helpers.DecredLightBlueColor)
		body.AddItem(form, 2, 0, true)

		accountNumbers := make([]uint32, len(accounts))
		accountNames := make([]string, len(accounts))
		for index, account := range accounts {
			accountNames[index] = account.Name
			accountNumbers[index] = account.Number
		}

		form.AddDropDown("Source Account: ", accountNames, 0, func(option string, optionIndex int) {
			accountNumber := accountNumbers[optionIndex]

			address, qr, err := generateAddress(wallet, accountNumber)
			if err != nil {
				errorText := fmt.Sprintf("Error: %s", err.Error())
				displayErrorMessage(errorText)
				return
			}

			qrCode := fmt.Sprintf(qr.ToSmallString(false))
			displayOutput(qrCode, address)	
		})

		form.SetCancelFunc(clearFocus)
	}

	hintTextView.SetText("TIP: Navigate with TAB and SHIFT+TAB, hit ENTER to generate Address. ESC to return to navigation menu")

	setFocus(body)
	return body
}

func generateAddress(wallet walletcore.Wallet, accountNumber uint32) (string, *qrcode.QRCode, error) {
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
