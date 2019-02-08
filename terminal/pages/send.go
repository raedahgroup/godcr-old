package pages

import (
	"fmt"

	"github.com/rivo/tview"
	"github.com/raedahgroup/godcr/app/walletcore"
)

func SendPage(wallet walletcore.Wallet, setFocus func(p tview.Primitive) *tview.Application, clearFocus func()) tview.Primitive {
	accounts, err := wallet.AccountsOverview(walletcore.DefaultRequiredConfirmations)
	if err != nil {
		fmt.Sprintf("Error fetching accounts: %s", err)
		// return err
	}
	accountNames := make([]string, len(accounts))
	for index, account := range accounts {
		accountNames[index] = account.Name
	}
	// fmt.Println(accountNames)

	//Form for Sending
	body := tview.NewForm()
	body.AddDropDown("Account", []string{"Dafault", "...."}, 0, nil)
	body.AddInputField("Amount", "", 20, nil, func (value string) {
		
	})
	body.AddInputField("Destination Address", "", 20, nil, nil)
	body.AddCheckbox("Unconfirmed", false, func(checked bool) {
		
	})
	body.AddButton("Send", func() {
		Account := body.GetFormItem(0)
		Amount := body.GetFormItem(1)
		Address := body.GetFormItem(2)
		Unconfirmed := body.GetFormItem(3)
		fmt.Println(Amount, Unconfirmed, Account,Address)
	})
	body.AddButton("Cancel", func() {
		clearFocus()
	})

	setFocus(body)
	
	return body
}
