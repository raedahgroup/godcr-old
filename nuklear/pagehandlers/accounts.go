package pagehandlers

import (
	"fmt"
	"strconv"

	"github.com/aarzilli/nucular"
	"github.com/raedahgroup/godcr/app/config"
	"github.com/raedahgroup/godcr/app/walletcore"
	"github.com/raedahgroup/godcr/nuklear/styles"
	"github.com/raedahgroup/godcr/nuklear/widgets"
)

type account struct {
	setAsHidden         bool
	setAsDefaultAccount bool
	account             *walletcore.Account
}

type AccountsHandler struct {
	err                error
	accounts           []*account
	wallet             walletcore.Wallet
	settings           *config.Settings
	isFetchingAccounts bool
	hdPath             string
}

func (handler *AccountsHandler) BeforeRender(wallet walletcore.Wallet, settings *config.Settings, refreshWindowDisplay func()) bool {
	handler.wallet = wallet
	handler.err = nil
	handler.accounts = nil
	handler.isFetchingAccounts = false
	handler.settings = settings

	return true
}

func (handler *AccountsHandler) Render(window *nucular.Window) {
	if handler.accounts == nil && handler.err == nil {
		go handler.fetchAccounts(window.Master().Changed)
	}

	widgets.PageContentWindowDefaultPadding("Accounts", window, func(contentWindow *widgets.Window) {
		// show accounts if any
		if len(handler.accounts) > 0 {
			handler.renderAccounts(contentWindow)
		}

		// show error if any
		if handler.err != nil {
			contentWindow.DisplayErrorMessage("", handler.err)
		}

		// show loading indicator if fetching
		if handler.isFetchingAccounts {
			contentWindow.DisplayIsLoadingMessage()
		}
	})
}

func (handler *AccountsHandler) fetchAccounts(refreshWindowDisplay func()) {
	handler.isFetchingAccounts = true
	defer func() {
		handler.isFetchingAccounts = false
		refreshWindowDisplay()
	}()

	accounts, err := handler.wallet.AccountsOverview(walletcore.DefaultRequiredConfirmations)
	if err != nil {
		handler.err = err
		return
	}

	for _, accountItem := range accounts {
		var setAsHidden, setAsDefaultAccount bool
		for _, hiddenAccount := range handler.settings.HiddenAccounts {
			if uint32(hiddenAccount) == accountItem.Number {
				setAsHidden = true
				break
			}
		}

		acc := &account{
			setAsHidden:         setAsHidden,
			setAsDefaultAccount: setAsDefaultAccount,
			account:             accountItem,
		}
		handler.accounts = append(handler.accounts, acc)
	}
}

func (handler *AccountsHandler) renderAccounts(window *widgets.Window) {
	for _, item := range handler.accounts {
		headerLabel := item.account.Name + " - " + item.account.Balance.Total.String()
		if item.account.Balance.Total != item.account.Balance.Spendable {
			headerLabel += fmt.Sprintf(" (Spendable: %s )", item.account.Balance.Spendable.String())
		}

		window.Row(30).Dynamic(1)
		if window.TreePush(nucular.TreeNode, headerLabel, false) {
			window.AddLabelWithFont("Properties", "LC", styles.BoldPageContentFont)

			window.Row(90).Dynamic(1)
			widgets.GroupWindow("", window.Window, 0, func(propertyWindow *widgets.Window) {
				table := widgets.NewTable()
				table.AddRow(
					widgets.NewLabelTableCell("Account Number", "LC"),
					widgets.NewLabelTableCell(strconv.Itoa(int(item.account.Number)), "LC"),
				)
				table.AddRow(
					widgets.NewLabelTableCell("HD Path", "LC"),
					widgets.NewLabelTableCell("ggg", "LC"),
				)
				table.AddRow(
					widgets.NewLabelTableCell("Keys", "LC"),
					widgets.NewLabelTableCell(fmt.Sprintf("%d External, %d Internal, %d Imported", item.account.ExternalKeyCount, item.account.InternalKeyCount, item.account.ImportedKeyCount), "LC"),
				)
				table.Render(propertyWindow)
			})

			window.AddLabelWithFont("Settings", "LC", styles.BoldPageContentFont)

			window.Row(60).Dynamic(1)
			widgets.GroupWindow("", window.Window, 0, func(settingsWindow *widgets.Window) {
				settingsWindow.AddCheckbox("Hide This Account (Account balance will be ignored)", &item.setAsHidden, handler.toggleAccountVisibilty(item, settingsWindow))
				settingsWindow.AddCheckbox("Default Account (Make this account default for all outgoing and incoming transactions)", &item.setAsDefaultAccount, handler.toggleDefaultAccount(item, settingsWindow))
			})

			window.TreePop()
		}
	}
}

func (handler *AccountsHandler) toggleAccountVisibilty(account *account, window *widgets.Window) func() {
	return func() {
		if account.setAsHidden {
			handler.hideAccount(account, window)
		} else {
			handler.revealAccount(account, window)
		}
	}
}

func (handler *AccountsHandler) toggleDefaultAccount(account *account, window *widgets.Window) func() {
	return func() {
		if account.setAsDefaultAccount {
			handler.setDefaultAccount(account, window)
		} else {
			handler.unsetDefaultAccount(account, window)
		}
	}
}

func (handler *AccountsHandler) hideAccount(account *account, window *widgets.Window) {
	accountToBeHidden := account.account.Number
	hiddenAccounts := handler.settings.HiddenAccounts

	// make sure the account is not already set to be hidden
	for _, hiddenAccount := range hiddenAccounts {
		if hiddenAccount == accountToBeHidden {
			widgets.NewAlertWidget("Account is already hidden", true, window)
			return
		}
	}

	hiddenAccounts = append(hiddenAccounts, accountToBeHidden)
	err := config.UpdateConfigFile(func(cnfg *config.ConfFileOptions) {
		cnfg.HiddenAccounts = hiddenAccounts
	})
	if err != nil {
		widgets.NewAlertWidget(fmt.Sprintf("Error hidding account: %s", err.Error()), true, window)
		return
	}
	handler.settings.HiddenAccounts = hiddenAccounts
	widgets.NewAlertWidget("Successfully set account as hidden", false, window)
}

func (handler *AccountsHandler) revealAccount(account *account, window *widgets.Window) {
	defer window.Master().Changed()

	hiddenAccounts := handler.settings.HiddenAccounts
	for index := range handler.settings.HiddenAccounts {
		if hiddenAccounts[index] == account.account.Number {
			hiddenAccounts = append(hiddenAccounts[:index], hiddenAccounts[index+1:]...)
			err := config.UpdateConfigFile(func(cnfg *config.ConfFileOptions) {
				cnfg.HiddenAccounts = hiddenAccounts
			})
			if err != nil {
				widgets.NewAlertWidget(fmt.Sprintf("Error revealing account: %s", err.Error()), true, window)
				return
			}

			handler.settings.HiddenAccounts = hiddenAccounts
			widgets.NewAlertWidget("Successfully revealed account", false, window)
			return
		}
	}
	widgets.NewAlertWidget("Error revealing account", true, window)
}

func (handler *AccountsHandler) setDefaultAccount(account *account, window *widgets.Window) {
	// first check if account is alreadr default
	if handler.settings.DefaultAccount == account.account.Number {
		widgets.NewAlertWidget("Account already default the default account account", true, window)
		return
	}

	// uncheck the other default account checkbox
	for _, acc := range handler.accounts {
		if acc.setAsDefaultAccount && acc.account.Number != account.account.Number {
			acc.setAsDefaultAccount = false
		}
	}

	err := config.UpdateConfigFile(func(cnfg *config.ConfFileOptions) {
		cnfg.DefaultAccount = account.account.Number
	})
	if err != nil {
		widgets.NewAlertWidget(fmt.Sprintf("Error setting default account: %s", err.Error()), true, window)
		return
	}
	handler.settings.DefaultAccount = account.account.Number
	widgets.NewAlertWidget("Successfully set default account", false, window)
}

func (handler *AccountsHandler) unsetDefaultAccount(account *account, window *widgets.Window) {
	// first check if account is default
	if handler.settings.DefaultAccount != account.account.Number {
		widgets.NewAlertWidget("Account is not set as default", true, window)
		return
	}

	err := config.UpdateConfigFile(func(cnfg *config.ConfFileOptions) {
		cnfg.DefaultAccount = 0
	})
	if err != nil {
		widgets.NewAlertWidget(fmt.Sprintf("Error updating default account: %s", err.Error()), true, window)
		return
	}
	widgets.NewAlertWidget("Successfully updated default account", false, window)
	handler.settings.DefaultAccount = 0
}
