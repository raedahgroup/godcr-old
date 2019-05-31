package pages

import (
	"fmt"
	"strconv"

	"github.com/gdamore/tcell"
	"github.com/raedahgroup/godcr/app/walletcore"
	"github.com/raedahgroup/godcr/terminal/helpers"
	"github.com/raedahgroup/godcr/terminal/primitives"
	"github.com/rivo/tview"
)

func accountsPage(wallet walletcore.Wallet, hintTextView *primitives.TextView, tviewApp *tview.Application, clearFocus func()) tview.Primitive {
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

	accountSettingFlex := tview.NewFlex().SetDirection(tview.FlexRow)

	displayAccountsTable := func() {
		accountPage.RemoveItem(accountPropertiesTable)
		accountPage.RemoveItem(accountSettingFlex)

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

	h := tview.NewCheckbox()
	d := tview.NewCheckbox()

	// method for getting transaction details when a tx is selected from the history table
	accountsTable.SetSelectedFunc(func(row, column int) {
		accountPage.RemoveItem(accountsTable)
		accountPage.RemoveItem(accountSettingFlex)

		titleTextView.SetText("Account Details")
		hintTextView.SetText("TIP: Use ARROW UP/DOWN to scroll, \nBACKSPACE to retun to accounts page, ESC to return to navigation menu")

		accountPage.AddItem(accountPropertiesTable, 9, 0, true)
		tviewApp.SetFocus(accountPropertiesTable)
		displayAccountsDetails(wallet, accountPropertiesTable, row, displayMessage)

		accountSettingFlex.AddItem(tview.NewFlex().SetDirection(tview.FlexColumn).
			AddItem(h, 3, 0, true).
			AddItem(primitives.NewLeftAlignedTextView("Hide this account (Account balance will be ignored)"), 0, 1, false), 2, 0, false).
			AddItem(tview.NewFlex().SetDirection(tview.FlexColumn).
				AddItem(d, 3, 0, true).
				AddItem(primitives.NewLeftAlignedTextView("Default account (Make this account default for all outgoing and incoming transactions)"), 0, 1, false), 2, 0, false)
		accountPage.AddItem(accountSettingFlex, 2, 0, false)

	})

	// handler for returning back to accounts page
	accountPropertiesTable.SetInputCapture(func(event *tcell.EventKey) *tcell.EventKey {
		if event.Key() == tcell.KeyEscape || event.Key() == tcell.KeyBackspace || event.Key() == tcell.KeyBackspace2 {
			displayAccountsTable()
			return nil
		}
		if event.Key() == tcell.KeyTab {
			tviewApp.SetFocus(h)
			return nil
		}

		return event
	})

	// handler for returning back to accounts page
	h.SetInputCapture(func(event *tcell.EventKey) *tcell.EventKey {
		if event.Key() == tcell.KeyTab {
			tviewApp.SetFocus(d)
			return nil
		}

		return event
	})

	// handler for returning back to accounts page
	d.SetInputCapture(func(event *tcell.EventKey) *tcell.EventKey {
		if event.Key() == tcell.KeyTab {
			tviewApp.SetFocus(h)
			return nil
		}

		return event
	})

	displayAccountsTable()
	displayWalletAcccounts(wallet, accountsTable, displayMessage)

	tviewApp.SetFocus(accountPage)
	return accountPage
}

func fetchAccounts(wallet walletcore.Wallet, displayMessage func(string)) []*walletcore.Account {
	accounts, err := wallet.AccountsOverview(walletcore.DefaultRequiredConfirmations)
	if err != nil {
		displayMessage(err.Error())
	}

	return accounts
}

func displayWalletAcccounts(wallet walletcore.Wallet, accountsTable *tview.Table, displayMessage func(string)) {
	accounts := fetchAccounts(wallet, displayMessage)

	for row, account := range accounts {
		accountsTable.SetCell(row, 0, tview.NewTableCell(fmt.Sprintf("%-5s", account.Name)).SetAlign(tview.AlignLeft).SetMaxWidth(1).SetExpansion(1))
		accountsTable.SetCell(row, 1, tview.NewTableCell(fmt.Sprintf("%-5s", account.Balance.Total)).SetAlign(tview.AlignLeft).SetMaxWidth(1).SetExpansion(1))
		if account.Balance.Total != account.Balance.Spendable {
			accountsTable.SetCell(row, 2, tview.NewTableCell(fmt.Sprintf("%5s", account.Balance.Spendable)).SetAlign(tview.AlignLeft).SetMaxWidth(1).SetExpansion(1))
		}
	}

}

func displayAccountsDetails(wallet walletcore.Wallet, accountPropertiesTable *tview.Table, row int, displayMessage func(string)) {
	var networkHDPath string
	if wallet.NetType() == "testnet3" {
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

	accounts := fetchAccounts(wallet, displayMessage)

	accountPropertiesTable.SetCellSimple(0, 1, accounts[row].Name)
	accountPropertiesTable.SetCellSimple(1, 1, accounts[row].Balance.Total.String())
	accountPropertiesTable.SetCellSimple(2, 1, accounts[row].Balance.Spendable.String())
	accountPropertiesTable.SetCellSimple(4, 1, strconv.FormatInt(int64(accounts[row].Number), 10))
	accountPropertiesTable.SetCellSimple(5, 1, networkHDPath)
	accountPropertiesTable.SetCellSimple(6, 1, fmt.Sprintf("%d External, %d Internal, %d Imported", accounts[row].ExternalKeyCount,
		accounts[row].InternalKeyCount,
		accounts[row].ImportedKeyCount))
}
