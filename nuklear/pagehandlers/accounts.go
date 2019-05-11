package pagehandlers

import (
	"fmt"
	"strconv"

	"github.com/aarzilli/nucular"
	"github.com/raedahgroup/godcr/app/walletcore"
	"github.com/raedahgroup/godcr/nuklear/styles"
	"github.com/raedahgroup/godcr/nuklear/widgets"
)

type account struct {
	isHidden         bool
	isDefaultAccount bool
	account          *walletcore.Account
}

type AccountsHandler struct {
	err                error
	accounts           []*account
	wallet             walletcore.Wallet
	isFetchingAccounts bool
	hdPath             string
}

func (handler *AccountsHandler) BeforeRender(wallet walletcore.Wallet, refreshWindowDisplay func()) bool {
	handler.wallet = wallet
	handler.err = nil
	handler.accounts = nil
	handler.isFetchingAccounts = false

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
			contentWindow.DisplayErrorMessage("Error fetching accounts", handler.err)
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
		acc := &account{
			isHidden:         false,
			isDefaultAccount: false,
			account:          accountItem,
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
				settingsWindow.AddCheckbox("Hide This Account (Account balance will be ignored)", &item.isHidden, handler.hideAccount())
				settingsWindow.AddCheckbox("Default Account (Make this account default for all outgoing and incoming transactions)", &item.isDefaultAccount, handler.setDefaultAccount())
			})

			window.TreePop()
		}
	}
}

func (handler *AccountsHandler) hideAccount() func() {
	return func() {
		fmt.Println("ddd")
	}
}

func (handler *AccountsHandler) setDefaultAccount() func() {
	return func() {
		fmt.Println("www")
	}
}
