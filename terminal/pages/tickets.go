package pages

import (
	"fmt"

	"github.com/rivo/tview"
	"github.com/raedahgroup/godcr/app/walletcore"
)

func PurchaseTicketPage(wallet walletcore.Wallet, setFocus func(p tview.Primitive) *tview.Application, clearFocus func()) tview.Primitive {
	body := tview.NewFlex().SetDirection(tview.FlexRow)
	form := tview.NewForm()

	body.AddItem(tview.NewTextView().SetTextAlign(tview.AlignCenter).SetText("Purchase Tickets"), 4, 1, false)

	accounts, err := wallet.AccountsOverview(walletcore.DefaultRequiredConfirmations)
	if err != nil {
		return body.AddItem(tview.NewTextView().SetTextAlign(tview.AlignCenter).SetText(fmt.Sprintf("Error: %s", err.Error())), 0, 1, false)
	}

	accountNumbers := make([]uint32, len(accounts))
	accountOverviews := make([]string, len(accounts))
	for index, account := range accounts {
		accountOverviews[index] = fmt.Sprintf("%s - Total %s (Spendable %s)", account.Name, account.Balance.Total.String(), account.Balance.Spendable.String())
		accountNumbers[index] = account.Number
	
	
	body.AddItem(form.AddDropDown("Account", []string{accountOverviews[index]}, 0, func(option string, optionIndex int) {
		// accountNum = accountN[optionIndex]
		}).
		AddButton("Generate", func() {
			// address, qr, err := generateAddress(wallet, accountNum)
			// if err != nil {
			// 	body.AddItem(tview.NewTextView().SetTextAlign(tview.AlignCenter).SetText(fmt.Sprintf("Error: %s", err.Error())), 3, 1, false)
			// 	return
			// }
			// body.AddItem(tview.NewTextView().SetTextAlign(tview.AlignLeft).SetText(fmt.Sprintf("Address: %s", address)), 4, 1, false).
			// 	AddItem(tview.NewTextView().SetTextAlign(tview.AlignLeft).SetText(fmt.Sprintf(qr.ToSmallString(false))), 0, 1, false)
		}).SetItemPadding(17).SetHorizontal(true).SetCancelFunc(func() {
		clearFocus()
	}), 4, 1, true)
}
	setFocus(body)
	return body
}
