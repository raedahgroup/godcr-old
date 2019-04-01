package pages

import (
	"fmt"

	"github.com/gdamore/tcell"
	"github.com/raedahgroup/godcr/app/walletcore"
	"github.com/raedahgroup/godcr/terminal/primitives"
	"github.com/rivo/tview"
	qrcode "github.com/skip2/go-qrcode"
)

func receivePage(wallet walletcore.Wallet, setFocus func(p tview.Primitive) *tview.Application, clearFocus func()) tview.Primitive {
	body := tview.NewFlex().SetDirection(tview.FlexRow)
	form := tview.NewForm()

	body.AddItem(primitives.NewCenterAlignedTextView("Generate Receive Address"), 4, 1, false)

	accounts, err := wallet.AccountsOverview(walletcore.DefaultRequiredConfirmations)
	if err != nil {
		return body.AddItem(primitives.NewCenterAlignedTextView(fmt.Sprintf("Error: %s", err.Error())), 0, 1, false)
	}
	if len(accounts) == 1 {
		address, qr, err := generateAddress(wallet, accounts[0].Number)
		if err != nil {
			return body.AddItem(primitives.NewCenterAlignedTextView(fmt.Sprintf("Error: %s", err.Error())), 0, 1, false)
		}
		body.AddItem(primitives.NewCenterAlignedTextView(fmt.Sprintf("Address: %s", address)).SetDoneFunc(func(key tcell.Key) {
			if key == tcell.KeyEscape {
				clearFocus()
			}
		}), 2, 1, true).
			AddItem(primitives.NewCenterAlignedTextView(fmt.Sprintf(qr.ToSmallString(false))).SetDoneFunc(func(key tcell.Key) {
				if key == tcell.KeyEscape {
					clearFocus()
				}
			}), 0, 1, true)
	} else {
		var accountNum uint32
		accountN := make([]uint32, len(accounts))
		accountNames := make([]string, len(accounts))
		for index, account := range accounts {
			accountNames[index] = account.Name
			body.AddItem(form.AddDropDown("Account", []string{accountNames[index]}, 0, func(option string, optionIndex int) {
				accountNum = accountN[optionIndex]
			}).
				AddButton("Generate", func() {
					address, qr, err := generateAddress(wallet, accountNum)
					if err != nil {
						body.AddItem(primitives.NewCenterAlignedTextView(fmt.Sprintf("Error: %s", err.Error())), 3, 1, false)
						return
					}
					body.AddItem(primitives.NewCenterAlignedTextView(fmt.Sprintf("Address: %s", address)), 4, 1, false).
						AddItem(primitives.NewCenterAlignedTextView(fmt.Sprintf(qr.ToSmallString(false))), 0, 1, false)
				}).SetItemPadding(17).SetHorizontal(true).SetCancelFunc(func() {
				clearFocus()
			}), 4, 1, true)
		}
	}

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
