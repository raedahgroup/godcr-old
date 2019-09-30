package widgets

import (
	"fmt"

	"github.com/aarzilli/nucular/label"
	"github.com/decred/dcrd/dcrutil"
	"github.com/raedahgroup/dcrlibwallet"
)

type AccountSelector struct {
	prompt               string
	accountNames         []string
	accountNumbers       []int32
	accountsFetchError   error
	selectedAccountIndex int
	selectionChanged     func()
}

const (
	defaultAccountSelectorWidth = 200
	accountSelectorHeight       = 25
)

func AccountSelectorWidget(prompt string, wallet *dcrlibwallet.LibWallet, spendUnconfirmed, showBalance bool,
	selectionChanged func()) *AccountSelector {

	accountSelector := &AccountSelector{
		prompt: prompt,
		selectionChanged: selectionChanged,
	}

	var confirmations int32 = dcrlibwallet.DefaultRequiredConfirmations
	if spendUnconfirmed {
		confirmations = 0
	}

	getAccountsResp, err := wallet.GetAccountsRaw(confirmations)
	if err != nil {
		accountSelector.accountsFetchError = err
		return accountSelector
	}

	accountSelector.accountNames = make([]string, len(getAccountsResp.Acc))
	accountSelector.accountNumbers = make([]int32, len(getAccountsResp.Acc))

	for i, account := range getAccountsResp.Acc {
		if showBalance {
			accountNameWithBalance := fmt.Sprintf("%s [%s]", account.Name, dcrutil.Amount(account.Balance.Spendable).String())
			accountSelector.accountNames[i] = accountNameWithBalance
		} else {
			accountSelector.accountNames[i] = account.Name
		}
		accountSelector.accountNumbers[i] = account.Number
	}

	return accountSelector
}

func (accountSelector *AccountSelector) Render(window *Window, addColumns ...int) {
	accountSelectorWidth := defaultAccountSelectorWidth
	if len(accountSelector.accountNames) == 1 {
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
	} else if len(accountSelector.accountNames) == 1 {
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

func (accountSelector *AccountSelector) GetSelectedAccountNumber() int32 {
	if accountSelector.selectedAccountIndex < len(accountSelector.accountNumbers) {
		return accountSelector.accountNumbers[accountSelector.selectedAccountIndex]
	}
	return 0
}

func (accountSelector *AccountSelector) Reset() {
	accountSelector.selectedAccountIndex = 0
}
