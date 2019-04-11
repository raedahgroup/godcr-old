package widgets

import (
	"github.com/aarzilli/nucular"
	"github.com/raedahgroup/godcr/app/walletcore"
)

type AccountSelector struct {
	accounts     []*walletcore.Account
	accountNames []string
	accountsFetchError error
	selectedAccountIndex int
}

func NewAccountSelectionWidget(prompt string, showBalance bool, wallet walletcore.Wallet) (accountSelector *AccountSelector) {
	accountSelector = &AccountSelector{}

	// print account selection prompt / label
	accountSelector.accounts, accountSelector.accountsFetchError = wallet.AccountsOverview(walletcore.DefaultRequiredConfirmations)
	if accountSelector.accountsFetchError != nil {
		return
	}

	accountSelector.accountNames = make([]string, len(accountSelector.accounts))
	for i, account := range accountSelector.accounts {
		accountSelector.accountNames[i] = account.String()
	}

	return &AccountSelector{
		accounts:             accounts,
		accountNames:         comboItems,
		selectedAccountIndex: 0,
	}
}

func (a *AccountSelector) Render(window *nucular.Window) {
	if len(a.accounts) > 1 {
		a.selectedAccountIndex = window.ComboSimple(a.accountNames, a.selectedAccountIndex, 30)
	} else {
		account := a.accounts[0]
		a.selectedAccountIndex = 0
		window.Label(account.String(), "LC")
	}
}

func (a *AccountSelector) GetSelectedAccount() *walletcore.Account {
	accountName := a.accountNames[a.selectedAccountIndex]
	for _, account := range a.accounts {
		if account.Name == accountName {
			return account
		}
	}

	return nil
}

func (a *AccountSelector) GetSelectedAccountNumber() uint32 {
	selectedAccount := a.GetSelectedAccount()
	if selectedAccount != nil {
		return selectedAccount.Number
	}

	return 0
}

func (a *AccountSelector) Reset() {
	if len(a.accounts) > 1 {
		a.selectedAccountIndex = 0
	}
}