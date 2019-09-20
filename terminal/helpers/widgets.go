package helpers

//import (
//	"github.com/raedahgroup/godcr/app/walletcore"
//	"github.com/raedahgroup/godcr/terminal/primitives"
//)
//
//type AccountSelectionWidgetData struct {
//	Label                 string
//	Accounts              []*walletcore.Account
//	ShowOnlyAccountName   bool
//	SelectedAccountNumber uint32
//}
//
//func AddAccountSelectionWidgetToForm(form *primitives.Form, data *AccountSelectionWidgetData) {
//	accountText := func(account *walletcore.Account) string {
//		if data.ShowOnlyAccountName {
//			return account.Name
//		} else {
//			return account.String()
//		}
//	}
//
//	if len(data.Accounts) == 1 {
//		account := data.Accounts[0]
//		data.SelectedAccountNumber = account.Number
//
//		accountName := accountText(account)
//		accountWidget := primitives.NewLeftAlignedTextView(accountName)
//
//		accountFormItem := primitives.NewTextViewFormItem(accountWidget, 20, 1, true)
//		accountFormItem.SetLabel(data.Label)
//
//		form.AddFormItem(accountFormItem)
//
//		return
//	}
//
//	accountNames := make([]string, len(data.Accounts))
//	accountNumbers := make([]uint32, len(data.Accounts))
//	for index, account := range data.Accounts {
//		accountNames[index] = accountText(account)
//		accountNumbers[index] = account.Number
//	}
//
//	form.AddDropDown(data.Label, accountNames, 0, func(option string, optionIndex int) {
//		data.SelectedAccountNumber = accountNumbers[optionIndex]
//	})
//
//	return
//}
