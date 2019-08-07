package pages

import (
	"fmt"
	"strconv"

	"github.com/gdamore/tcell"
	"github.com/raedahgroup/godcr/app/config"
	"github.com/raedahgroup/godcr/app/walletcore"
	"github.com/raedahgroup/godcr/terminal/helpers"
	"github.com/raedahgroup/godcr/terminal/primitives"
	"github.com/rivo/tview"
)

func accountsPage(wallet walletcore.Wallet, hintTextView *primitives.TextView, settings config.Settings, tviewApp *tview.Application, clearFocus func()) tview.Primitive {
	accountPage := tview.NewFlex().SetDirection(tview.FlexRow)

	messageTextView := primitives.WordWrappedTextView("")
	messageTextView.SetTextColor(helpers.DecredOrangeColor)

	displayMessage := func(message string) {
		accountPage.RemoveItem(messageTextView)
		if message != "" {
			messageTextView.SetText(message)
			accountPage.AddItem(messageTextView, 2, 0, false)
		}
	}

	// page title
	titleTextView := primitives.NewLeftAlignedTextView("Accounts")
	accountPage.AddItem(titleTextView, 2, 0, false)

	accountsTable := tview.NewTable().
		SetBorders(false).
		SetFixed(1, 0).
		SetSelectable(true, false)

	accountPropertiesTable := tview.NewTable().SetBorders(false)

	hideAccount := primitives.NewCheckbox("Hide this account (Account balance will be ignored): ")
	defaultAccount := primitives.NewCheckbox("Default account (Make this account default for all outgoing and incoming transactions): ")

	displayAccountsTable := func() {
		accountPage.RemoveItem(accountPropertiesTable)
		accountPage.RemoveItem(hideAccount)
		accountPage.RemoveItem(defaultAccount)

		titleTextView.SetText("Accounts")
		hintTextView.SetText("TIP: Use ARROW UP/DOWN to select an account,\nENTER to view details, ESC to return to navigation menu")

		accountPage.AddItem(accountsTable, 0, 1, true)
		tviewApp.SetFocus(accountsTable)
	}

	accountsTable.SetDoneFunc(func(key tcell.Key) {
		if key == tcell.KeyEscape {
			clearFocus()
		}
	})

	accounts, err := wallet.AccountsOverview(walletcore.DefaultRequiredConfirmations)
	if err != nil {
		displayMessage(err.Error())
	}

	// method for getting transaction details when a tx is selected from the history table
	var selectedAccount *walletcore.Account
	hiddenAccounts := settings.HiddenAccounts
	accountsTable.SetSelectedFunc(func(row, column int) {
		accountPage.RemoveItem(accountsTable)
		selectedRow := row - 1
		selectedAccount = accounts[selectedRow]

		titleTextView.SetText("Account Details")
		hintTextView.SetText("TIP: Use TAB key to switch between checkbox, \nBACKSPACE to retun to accounts page, ESC to return to navigation menu")

		accountPage.AddItem(accountPropertiesTable, 9, 0, true)
		tviewApp.SetFocus(hideAccount)
		displayAccountsDetails(wallet.NetType(), selectedAccount, accountPropertiesTable, displayMessage)

		if settings.DefaultAccount != selectedAccount.Number {
			defaultAccount.SetChecked(false)
		} else {
			defaultAccount.SetChecked(true)
		}

		var accountIsHidden bool
		for _, hiddenAccount := range hiddenAccounts {
			if selectedAccount.Number == hiddenAccount {
				accountIsHidden = true
			}
		}

		hideAccount.SetChecked(accountIsHidden)

		accountPage.AddItem(hideAccount, 2, 0, false)
		accountPage.AddItem(defaultAccount, 2, 0, false)
	})

	defaultAccount.SetChangedFunc(func(checked bool) {
		if checked {
			err := config.UpdateConfigFile(func(cnfg *config.ConfFileOptions) {
				cnfg.DefaultAccount = selectedAccount.Number
			})
			if err != nil {
				displayMessage(err.Error())
				return
			}

			settings.DefaultAccount = selectedAccount.Number
			return
		}

		err := config.UpdateConfigFile(func(cnfg *config.ConfFileOptions) {
			cnfg.DefaultAccount = 0
		})
		if err != nil {
			displayMessage(err.Error())
			return
		}

		settings.DefaultAccount = 0
		return
	})

	hideAccount.SetChangedFunc(func(checked bool) {
		if checked {
			hiddenAccounts = append(hiddenAccounts, selectedAccount.Number)
			err := config.UpdateConfigFile(func(cnfg *config.ConfFileOptions) {
				cnfg.HiddenAccounts = hiddenAccounts
			})
			if err != nil {
				displayMessage(err.Error())
				return
			}

			settings.HiddenAccounts = hiddenAccounts
			return
		}

		for index := range settings.HiddenAccounts {
			if hiddenAccounts[index] == selectedAccount.Number {
				hiddenAccounts = append(hiddenAccounts[:index], hiddenAccounts[index+1:]...)
				err := config.UpdateConfigFile(func(cnfg *config.ConfFileOptions) {
					cnfg.HiddenAccounts = hiddenAccounts
				})
				if err != nil {
					displayMessage(err.Error())
					return
				}

				settings.HiddenAccounts = hiddenAccounts
				return
			}
		}
	})

	// handler for returning back to accounts page
	accountPropertiesTable.SetInputCapture(func(event *tcell.EventKey) *tcell.EventKey {
		if event.Key() == tcell.KeyEscape || event.Key() == tcell.KeyBackspace || event.Key() == tcell.KeyBackspace2 {
			displayAccountsTable()
			return nil
		}
		if event.Key() == tcell.KeyTab {
			tviewApp.SetFocus(hideAccount)
			return nil
		}

		return event
	})

	// handler for returning back to accounts page
	hideAccount.SetInputCapture(func(event *tcell.EventKey) *tcell.EventKey {
		if event.Key() == tcell.KeyEscape || event.Key() == tcell.KeyBackspace || event.Key() == tcell.KeyBackspace2 {
			displayAccountsTable()
			return nil
		}

		if event.Key() == tcell.KeyTab {
			tviewApp.SetFocus(defaultAccount)
			return nil
		}

		return event
	})

	// handler for returning back to accounts page
	defaultAccount.SetInputCapture(func(event *tcell.EventKey) *tcell.EventKey {
		if event.Key() == tcell.KeyEscape || event.Key() == tcell.KeyBackspace || event.Key() == tcell.KeyBackspace2 {
			displayAccountsTable()
			return nil
		}

		if event.Key() == tcell.KeyTab {
			tviewApp.SetFocus(hideAccount)
			return nil
		}

		return event
	})

	displayAccountsTable()
	displayWalletAcccounts(accounts, accountsTable, displayMessage)

	tviewApp.SetFocus(accountPage)
	return accountPage
}

func displayWalletAcccounts(accounts []*walletcore.Account, accountsTable *tview.Table, displayMessage func(string)) {
	tableHeaderCell := func(text string) *tview.TableCell {
		return tview.NewTableCell(text).SetAlign(tview.AlignLeft).SetSelectable(false).SetMaxWidth(1).SetExpansion(1)
	}

	accountsTable.SetCell(0, 0, tableHeaderCell("Account Name"))
	accountsTable.SetCell(0, 1, tableHeaderCell("Total Balance"))
	accountsTable.SetCell(0, 2, tableHeaderCell("Spendable Balance"))

	for _, account := range accounts {
		row := accountsTable.GetRowCount()

		accountsTable.SetCell(row, 0, tview.NewTableCell(fmt.Sprintf("%-5s", account.Name)).SetAlign(tview.AlignLeft).SetMaxWidth(1).SetExpansion(1))
		accountsTable.SetCell(row, 1, tview.NewTableCell(fmt.Sprintf("%-5s", account.Balance.Total)).SetAlign(tview.AlignLeft).SetMaxWidth(1).SetExpansion(1))
		accountsTable.SetCell(row, 2, tview.NewTableCell(fmt.Sprintf("%5s", account.Balance.Spendable)).SetAlign(tview.AlignLeft).SetMaxWidth(1).SetExpansion(1))
	}
}

func displayAccountsDetails(netType string, account *walletcore.Account, accountPropertiesTable *tview.Table, displayMessage func(string)) {
	var networkHDPath string
	if netType == "testnet3" {
		networkHDPath = walletcore.TestnetHDPath
	} else {
		networkHDPath = walletcore.MainnetHDPath
	}

	accountPropertiesTable.SetCellSimple(0, 0, "Account Name:")
	accountPropertiesTable.SetCellSimple(1, 0, "Total Balance:")
	accountPropertiesTable.SetCellSimple(2, 0, "Spendable Balance: ")
	accountPropertiesTable.SetCellSimple(3, 0, "-properties-")
	accountPropertiesTable.SetCellSimple(4, 0, "Account Number:")
	accountPropertiesTable.SetCellSimple(5, 0, "HD Path:")
	accountPropertiesTable.SetCellSimple(6, 0, "Keys:")

	accountPropertiesTable.SetCellSimple(0, 1, account.Name)
	accountPropertiesTable.SetCellSimple(1, 1, account.Balance.Total.String())
	accountPropertiesTable.SetCellSimple(2, 1, account.Balance.Spendable.String())
	accountPropertiesTable.SetCellSimple(4, 1, strconv.FormatInt(int64(account.Number), 10))
	accountPropertiesTable.SetCellSimple(5, 1, networkHDPath)
	accountPropertiesTable.SetCellSimple(6, 1, fmt.Sprintf("%d External, %d Internal, %d Imported", account.ExternalKeyCount,
		account.InternalKeyCount,
		account.ImportedKeyCount))
}
