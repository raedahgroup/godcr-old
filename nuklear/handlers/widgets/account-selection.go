package widgets

import (
	"github.com/aarzilli/nucular"
	"github.com/raedahgroup/godcr/app/walletcore"
)

type AccountSelection struct {
	accounts   []*walletcore.Account
	comboItems []string

	selectedAccountIndex int
}

func NewAccountSelectionWidget(accounts []*walletcore.Account) *AccountSelection {
	comboItems := make([]string, len(accounts))
	for index, account := range accounts {
		comboItems[index] = account.String()
	}

	return &AccountSelection{
		accounts:             accounts,
		comboItems:           comboItems,
		selectedAccountIndex: 0,
	}
}

func (a *AccountSelection) Render(window *nucular.Window) {
	if len(a.accounts) > 1 {
		a.selectedAccountIndex = window.ComboSimple(a.comboItems, a.selectedAccountIndex, 30)
	} else {
		account := a.accounts[0]
		a.selectedAccountIndex = 0
		window.Label(account.String(), "LC")
	}
}

func (a *AccountSelection) GetSelectedAccount() *walletcore.Account {
	accountName := a.comboItems[a.selectedAccountIndex]
	for _, account := range a.accounts {
		if account.Name == accountName {
			return account
		}
	}

	return nil
}

func (a *AccountSelection) GetSelectedAccountNumber() uint32 {
	selectedAccount := a.GetSelectedAccount()
	if selectedAccount != nil {
		return selectedAccount.Number
	}

	return 0
}

func (a *AccountSelection) Reset() {
	if len(a.accounts) > 1 {
		a.selectedAccountIndex = 0
	}
}
