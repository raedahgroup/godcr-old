package pages

import (
	"fmt"

	"github.com/rivo/tview"
	"github.com/raedahgroup/godcr/app/walletcore"
)

func SendPage(wallet walletcore.Wallet, setFocus func(p tview.Primitive) *tview.Application, clearFocus func()) tview.Primitive {
	body := tview.NewFlex().SetDirection(tview.FlexRow)
	// textView := tview.NewTextView().SetTextAlign(tview.AlignCenter)
	form := tview.NewForm()

	accounts, err := wallet.AccountsOverview(walletcore.DefaultRequiredConfirmations)
	if err != nil {
		return body.AddItem(tview.NewTextView().SetTextAlign(tview.AlignCenter).SetText(fmt.Sprintf("Error: %s", err.Error())), 0, 1, false)
	}

	//Form for Sending
	var  amount, destination string
	var accountNum uint32
	var checked bool

	accountNames := make([]string, len(accounts))
	accountN := make([]uint32, len(accounts))
	for index, account := range accounts {
		accountNames[index] = account.Name
		body.AddItem(form.AddDropDown("Account", []string{accountNames[index]}, 0, func(option string, optionIndex int) {
			accountNum = accountN[optionIndex]
		}).
		AddInputField("Amount", "", 20, nil, func (text string) {
			// if text == "" {
			// 	// errMsg.SetText("field cannot be  0")
			// 	// return
			// }
			amount = text
		}).
		AddInputField("Destination Address", "", 20, nil, func (text string) {
			destination = text
		}).
		AddCheckbox("Unconfirmed", false, func(checked bool) {
			if checked != false {
				checked = true
			}
		}).
		AddPasswordField("Password", "", 20, '*', nil).
		AddButton("Send", func() {
			err := confBalance(accountNum, wallet)
			if err != nil {
				body.AddItem(tview.NewTextView().SetTextAlign(tview.AlignCenter).SetText(fmt.Sprintf("Error: %s", err.Error())), 0, 1, false)
				return 
			}
			fmt.Println(destination, amount, checked)
		}).
		AddButton("Cancel", func() {
			clearFocus()
		}), 0, 1, true)
	}

	setFocus(body)
	
	return body
}

// func SendPage(wallet walletcore.Wallet, setFocus func(p tview.Primitive) *tview.Application, clearFocus func()) tview.Primitive {
// 	textView := tview.NewTextView().SetTextAlign(tview.AlignCenter)
// 	flex := tview.NewFlex()
// }

func confBalance(accountNum uint32, wallet walletcore.Wallet) error{
	accountBalance, err := wallet.AccountBalance(accountNum, walletcore.DefaultRequiredConfirmations)
	if err != nil {
		return err
	}
	if accountBalance.Total != 0 {
		return fmt.Errorf("Selected account has 0 balance. Cannot proceed")
	}

	fmt.Println(accountBalance)
	return nil
}