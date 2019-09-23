package widgets

import (
	"github.com/raedahgroup/godcr/app/walletcore"
	"github.com/aarzilli/nucular/label"
)

type AccountSelector struct {
	prompt               string
	accounts             []*walletcore.Account
	accountNames         []string
	accountNumbers       []uint32
	accountsFetchError   error
	selectedAccountIndex int
	selectionChanged     func()
}

const (
	defaultAccountSelectorWidth = 200
	accountSelectorHeight       = 25
)

func AccountSelectorWidget(prompt string, spendUnconfirmed, showBalance bool, wallet walletcore.Wallet,
	selectionChanged func()) (accountSelector *AccountSelector) {

	accountSelector = &AccountSelector{
		prompt: prompt,
		selectionChanged: selectionChanged,
	}

	var confirmations int32 = walletcore.DefaultRequiredConfirmations
	if spendUnconfirmed {
		confirmations = 0
	}

	accountSelector.accounts, accountSelector.accountsFetchError = wallet.AccountsOverview(confirmations)
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
		window.DisplayErrorMessage("Fetch accounts error", accountSelector.accountsFetchError)
	} else if len(accountSelector.accounts) == 1 {
		accountSelector.selectedAccountIndex = 0
		window.Label(accountSelector.accountNames[0], LeftCenterAlign)
	} else {
		accountSelector.makeDropDown(window)
	}
}

// makeDropDown is adapted from nucular's Window.ComboSimple
// to allow triggering a callback when dropdown selection changes.
func (accountSelector *AccountSelector) makeDropDown(window *Window) {
	if len(accountSelector.accountNames) == 0 {
		return
	}

	items := accountSelector.accountNames
	itemHeight := int(float64(accountSelectorHeight) * window.Master().Style().Scaling)
	itemPadding := window.Master().Style().Combo.ButtonPadding.Y
	maxHeight := (len(items)+1)*itemHeight + itemPadding*3

	if w := window.Combo(label.T(items[accountSelector.selectedAccountIndex]), maxHeight, nil); w != nil {
		w.RowScaled(itemHeight).Dynamic(1)
		for i := range items {
			if w.MenuItem(label.TA(items[i], "LC")) {
				accountSelector.selectedAccountIndex = i
				if accountSelector.selectionChanged != nil {
					accountSelector.selectionChanged()
				}
			}
		}
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
