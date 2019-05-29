package pages

import (
	"fmt"
	"strconv"

	"github.com/gdamore/tcell"
	"github.com/raedahgroup/godcr/app/walletcore"
	// "github.com/raedahgroup/godcr/terminal/helpers"
	"github.com/raedahgroup/godcr/terminal/primitives"
	"github.com/rivo/tview"
)

func accountsPage(wallet walletcore.Wallet, hintTextView *primitives.TextView, setFocus func(p tview.Primitive) *tview.Application, clearFocus func()) tview.Primitive {
	body := tview.NewFlex().SetDirection(tview.FlexRow)

	// page title
	titleTextView := primitives.NewLeftAlignedTextView("Accounts")
	body.AddItem(titleTextView, 2, 0, false)

	accountsTable := tview.NewTable().
		SetBorders(false).
		SetFixed(1, 0).
		SetSelectable(true, false)

	accountPropertiesTable := tview.NewTable().SetBorders(false)

	displayAccountsTable := func() {
		body.RemoveItem(accountPropertiesTable)

		titleTextView.SetText("Accounts")
		hintTextView.SetText("TIP: Use ARROW UP/DOWN to select an account,\nENTER to view details, ESC to return to navigation menu")

		body.AddItem(accountsTable, 0, 1, true)
		setFocus(accountsTable)
	}

	accountsTable.SetDoneFunc(func(key tcell.Key) {
		if key == tcell.KeyEscape {
			clearFocus()
		}
	})

	// method for getting transaction details when a tx is selected from the history table
	accountsTable.SetSelectedFunc(func(row, column int) {
		body.RemoveItem(accountsTable)
		accountPropertiesTable.Clear()

		titleTextView.SetText("Account Details")
		hintTextView.SetText("TIP: Use ARROW UP/DOWN to scroll, \nBACKSPACE to retun to accounts page, ESC to return to navigation menu")

		body.AddItem(accountPropertiesTable, 0, 1, true)
		setFocus(accountPropertiesTable)
		displayAccountsDetails(wallet, accountPropertiesTable, row)
	})

	// handler for returning back to accounts page
	accountPropertiesTable.SetInputCapture(func(event *tcell.EventKey) *tcell.EventKey {
		if event.Key() == tcell.KeyEscape || event.Key() == tcell.KeyBackspace || event.Key() == tcell.KeyBackspace2 {
			displayAccountsTable()
			return nil
		}

		return event
	})

	displayAccountsTable()
	displayWalletAcccounts(wallet, accountsTable)

	setFocus(body)
	return body
}

func fetchAccounts(wallet walletcore.Wallet) []*walletcore.Account {
	accounts, err := wallet.AccountsOverview(walletcore.DefaultRequiredConfirmations)
	if err != nil {
		// errorTextView := primitives.NewCenterAlignedTextView(err.Error()).SetTextColor(helpers.DecredOrangeColor)
		// body.AddItem(errorTextView, 2, 0, false)
	}

	return accounts
}

func displayWalletAcccounts(wallet walletcore.Wallet, accountsTable *tview.Table) {
	accounts := fetchAccounts(wallet)

	for row, account := range accounts {
		accountsTable.SetCell(row, 0, tview.NewTableCell(fmt.Sprintf("%-10s", account.Name)).SetAlign(tview.AlignCenter).SetMaxWidth(1).SetExpansion(1))
		accountsTable.SetCell(row, 1, tview.NewTableCell(fmt.Sprintf("%-10s", account.Balance.Total)).SetAlign(tview.AlignCenter).SetMaxWidth(2).SetExpansion(1))
		if account.Balance.Total != account.Balance.Spendable {
			accountsTable.SetCell(row, 2, tview.NewTableCell(fmt.Sprintf("%15s", account.Balance.Spendable)).SetAlign(tview.AlignCenter).SetMaxWidth(3).SetExpansion(1))
		}
	}

}

func displayAccountsDetails(wallet walletcore.Wallet, accountPropertiesTable *tview.Table, row int) {
	var networkHDPath string
	if wallet.NetType() == "testnet3" {
		networkHDPath = walletcore.TestnetHDPath
	} else {
		networkHDPath = walletcore.MainnetHDPath
	}

	accountPropertiesTable.SetCellSimple(0, 0, "Account")
	accountPropertiesTable.SetCellSimple(1, 0, "Total Balance")
	accountPropertiesTable.SetCellSimple(2, 0, "Spendable Balance")
	accountPropertiesTable.SetCellSimple(3, 0, "-properties-")
	accountPropertiesTable.SetCellSimple(4, 0, "Account Number")
	accountPropertiesTable.SetCellSimple(5, 0, "HD Path")
	accountPropertiesTable.SetCellSimple(6, 0, "Keys")
	// accountPropertiesTable.SetCellSimple(7, 0, "-Settings-")
	// accountPropertiesTable.SetCellSimple(8, 0, "Fee Rate")

	accounts := fetchAccounts(wallet)

	accountPropertiesTable.SetCellSimple(0, 1, accounts[row].Name)
	accountPropertiesTable.SetCellSimple(1, 1, accounts[row].Balance.Total.String())
	accountPropertiesTable.SetCellSimple(2, 1, accounts[row].Balance.Spendable.String())
	accountPropertiesTable.SetCellSimple(4, 1, strconv.FormatInt(int64(accounts[row].Number), 10))
	accountPropertiesTable.SetCellSimple(5, 1, networkHDPath)
	accountPropertiesTable.SetCellSimple(6, 1, fmt.Sprintf("%d External, %d Internal, %d Imported", accounts[row].ExternalKeyCount,
		accounts[row].InternalKeyCount,
		accounts[row].ImportedKeyCount))

}
