package pages

import (
	"fmt"
    // "strconv"

	"github.com/rivo/tview"
	"github.com/raedahgroup/godcr/app/walletcore"
)

func SendPage(wallet walletcore.Wallet, setFocus func(p tview.Primitive) *tview.Application, clearFocus func()) tview.Primitive {
	// errMsg := tview.NewTextView().SetTextAlign(tview.AlignCenter)
	accounts, err := wallet.AccountsOverview(walletcore.DefaultRequiredConfirmations)
	if err != nil {
		fmt.Sprintf("Error fetching accounts: %s", err)
		// return err
	}

	//Form for Sending
	var  amount, destination string 
	var accountNum uint32
	var checked bool
	body := tview.NewForm()

	accountNames := make([]string, len(accounts))
	accountN := make([]uint32, len(accounts))
	for index, account := range accounts {
		accountNames[index] = account.Name
		body.AddDropDown("Account", []string{accountNames[index]}, 0, func(option string, optionIndex int) {
			accountNum = accountN[optionIndex]
			accountBalance, err := wallet.AccountBalance(accountNum, walletcore.DefaultRequiredConfirmations)
			if err != nil {
				fmt.Println(err)
			}
			fmt.Println(accountBalance)
		})
	}
	body.AddInputField("Amount", "", 20, nil, func (text string) {
		// tx, err := strconv.Atoi(text)
		// if err != nil {
		// 	// return errMsg.SetText(err.Error())
		// }
		// if tx == 0 {
		// 	// return errMsg.SetText("field must be greater than 0")
		// }
		// if tx == ""{
		// 	// return errMsg.SetText("field cannot be empty")
		// }

		amount = text
	})
	body.AddInputField("Destination Address", "", 20, nil, func (text string) {
		destination = text
	})
	body.AddCheckbox("Unconfirmed", false, func(checked bool) {
		if checked != false {
			checked = true
		}
	})
	body.AddButton("Send", func() {
		fmt.Println(destination, amount, checked)
	})
	body.AddButton("Cancel", func() {
		clearFocus()
	})

	setFocus(body)
	
	return body
}
