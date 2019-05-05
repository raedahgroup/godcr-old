package pages

import (
	"github.com/decred/dcrd/dcrutil"

	"fyne.io/fyne"
	"fyne.io/fyne/widget"
	godcrApp "github.com/raedahgroup/godcr/app"
	"github.com/raedahgroup/godcr/app/walletcore"
	"github.com/raedahgroup/godcr/fyne/widgets"
)

func OverviewPage(windows fyne.Window, App fyne.App) fyne.CanvasObject {
	label := widget.NewLabelWithStyle("Overview", fyne.TextAlignLeading, fyne.TextStyle{Italic: true, Bold: true})
	balanceLabel := widget.NewLabel("Current Total Balance")
	activityLabel := widget.NewLabel("Recent Activity")
	balance := widget.NewLabel(FetchBalance(Wallet))
	balanceLabel.TextStyle = fyne.TextStyle{Bold: true}
	activityLabel.TextStyle = fyne.TextStyle{Bold: true}
	table := widgets.NewTable()
	table, _ = FetchRecentActivity(Wallet, table, 5, false)
	return widget.NewVBox(
		label,
		balanceLabel,
		widgets.NewVSpacer(10),
		balance,
		activityLabel,
		table.CondensedTable())
}

func FetchBalance(wallet godcrApp.WalletMiddleware) string {
	accounts, err := wallet.AccountsOverview(walletcore.DefaultRequiredConfirmations)
	if err != nil {
		return err.Error()
	}

	return walletcore.WalletBalance(accounts)
}

func FetchRecentActivity(wallet godcrApp.WalletMiddleware, table *widgets.Table, noOfTransaction int, button bool) (*widgets.Table, error) {
	//table.Clear()
	txns, err := wallet.TransactionHistory(0, int32(noOfTransaction), nil)
	if err != nil {
		return nil, err
	}
	table.AddRowWithTextCells(
		widgets.NewTableTextCell("Account", fyne.TextAlignCenter, fyne.TextStyle{}, nil),
		widgets.NewTableTextCell("Date (UTC)", fyne.TextAlignCenter, fyne.TextStyle{}, nil),
		widgets.NewTableTextCell("Type", fyne.TextAlignCenter, fyne.TextStyle{}, nil),
		widgets.NewTableTextCell("Direction", fyne.TextAlignCenter, fyne.TextStyle{}, nil),
		widgets.NewTableTextCell("Amount", fyne.TextAlignCenter, fyne.TextStyle{}, nil),
		widgets.NewTableTextCell("Fee", fyne.TextAlignCenter, fyne.TextStyle{}, nil),
		widgets.NewTableTextCell("Status", fyne.TextAlignCenter, fyne.TextStyle{}, nil),
		widgets.NewTableTextCell("Hash", fyne.TextAlignCenter, fyne.TextStyle{}, nil))
	for _, tx := range txns {
		trimmedHash := tx.Hash[:len(tx.Hash)/2] + "..."
		hashButton := widgets.NewTableTextCell(trimmedHash, fyne.TextAlignLeading, fyne.TextStyle{}, nil)
		if button {
			hashButton = widgets.NewTableTextCell(trimmedHash, fyne.TextAlignLeading, fyne.TextStyle{}, func() {
				//todo
			})
		}
		table.AddRowWithTextCells(
			widgets.NewTableTextCell(tx.AccountName(), fyne.TextAlignLeading, fyne.TextStyle{}, nil),
			widgets.NewTableTextCell(tx.LongTime, fyne.TextAlignLeading, fyne.TextStyle{}, nil),
			widgets.NewTableTextCell(tx.Type, fyne.TextAlignLeading, fyne.TextStyle{}, nil),
			widgets.NewTableTextCell(tx.Direction.String(), fyne.TextAlignLeading, fyne.TextStyle{}, nil),
			widgets.NewTableTextCell(dcrutil.Amount(tx.Amount).String(), fyne.TextAlignTrailing, fyne.TextStyle{}, nil),
			widgets.NewTableTextCell(dcrutil.Amount(tx.Fee).String(), fyne.TextAlignTrailing, fyne.TextStyle{}, nil),
			widgets.NewTableTextCell(tx.Status, fyne.TextAlignLeading, fyne.TextStyle{}, nil),
			hashButton,
		)
	}
	return table, nil
}
