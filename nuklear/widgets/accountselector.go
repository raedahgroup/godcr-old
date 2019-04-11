package widgets

import (
	"github.com/raedahgroup/godcr/app/walletcore"
)

type AccountSelector struct {
	prompt               string
	accounts             []*walletcore.Account
	accountNames         []string
	accountNumbers       []uint32
	accountsFetchError   error
	selectedAccountIndex int
}

const (
	defaultAccountSelectorWidth = 200
	accountSelectorHeight       = 25
)

func AccountSelectorWidget(prompt string, showBalance bool, wallet walletcore.Wallet) (accountSelector *AccountSelector) {
	accountSelector = &AccountSelector{
		prompt: prompt,
	}

	accountSelector.accounts, accountSelector.accountsFetchError = wallet.AccountsOverview(walletcore.DefaultRequiredConfirmations)
	if accountSelector.accountsFetchError != nil {
		return
	}

	accountSelector.accountNames = make([]string, len(accountSelector.accounts))
	accountSelector.accountNumbers = make([]uint32, len(accountSelector.accounts))

	for i, account := range accountSelector.accounts {
		if showBalance {
			accountSelector.accountNames[i] = account.String()
		} else {
			accountSelector.accountNames[i] = account.Name
		}
		accountSelector.accountNumbers[i] = account.Number
	}

	return
}

func (accountSelector *AccountSelector) Render(window *Window, addColumns ...int) {
	accountSelectorWidth := defaultAccountSelectorWidth
	if len(accountSelector.accounts) == 1 {
		accountSelectorWidth = window.LabelWidth(accountSelector.accountNames[0])
	}

	// row with fixed column widths to hold account selection prompt, the account widget, and any other widgets that may be added later
	rowColumns := make([]int, 2)
	rowColumns[0] = window.LabelWidth(accountSelector.prompt)
	rowColumns[1] = accountSelectorWidth
	rowColumns = append(rowColumns, addColumns...)
	window.Row(accountSelectorHeight).Static(rowColumns...)

	// print account selection prompt / label
	window.Label(accountSelector.prompt, LeftCenterAlign)

	if accountSelector.accountsFetchError != nil {
		window.DisplayErrorMessage(accountSelector.accountsFetchError.Error())
	} else if len(accountSelector.accounts) == 1 {
		accountSelector.selectedAccountIndex = 0
		window.Label(accountSelector.accountNames[0], LeftCenterAlign)
	} else {
		accountSelector.selectedAccountIndex = window.ComboSimple(accountSelector.accountNames,
			accountSelector.selectedAccountIndex, accountSelectorHeight)
	}
}

func (accountSelector *AccountSelector) GetSelectedAccountNumber() uint32 {
	if accountSelector.selectedAccountIndex < len(accountSelector.accountNumbers) {
		return accountSelector.accountNumbers[accountSelector.selectedAccountIndex]
	}
	return 0
}

func (accountSelector *AccountSelector) Reset() {
	if len(accountSelector.accounts) > 1 {
		accountSelector.selectedAccountIndex = 0
	}
}
