package pagehandlers

import (
	"fmt"

	"github.com/aarzilli/nucular"
	"github.com/decred/dcrd/dcrutil"
	"github.com/raedahgroup/dcrlibwallet"
	"github.com/raedahgroup/godcr/nuklear/styles"
	"github.com/raedahgroup/godcr/nuklear/widgets"
)

type OverviewHandler struct {
	err      error
	accounts []*dcrlibwallet.Account
	wallet   *dcrlibwallet.LibWallet
}

func (handler *OverviewHandler) BeforeRender(wallet *dcrlibwallet.LibWallet, _ func()) {
	handler.wallet = wallet
	getAccountResp, err := wallet.GetAccountsRaw(dcrlibwallet.DefaultRequiredConfirmations)
	if err != nil {
		handler.err = err
	} else {
		handler.accounts = getAccountResp.Acc
	}
}

func (handler *OverviewHandler) Render(window *nucular.Window) {
	widgets.PageContentWindowDefaultPadding("Overview", window, func(contentWindow *widgets.Window) {
		contentWindow.AddLabelWithFont("Current Total Balance", widgets.LeftCenterAlign, styles.BoldPageContentFont)

		if handler.err != nil {
			contentWindow.DisplayErrorMessage("Error fetching accounts balance", handler.err)
			return
		}

		var totalBalance, spendableBalance dcrutil.Amount
		for _, account := range handler.accounts {
			totalBalance += dcrutil.Amount(account.Balance.Total)
			spendableBalance += dcrutil.Amount(account.Balance.Total)
		}

		var balance string
		if totalBalance != spendableBalance {
			balance = fmt.Sprintf("Total %s (Spendable %s)", totalBalance.String(), spendableBalance.String())
		} else {
			balance = totalBalance.String()
		}
		contentWindow.AddLabel(balance, widgets.LeftCenterAlign)

		contentWindow.AddHorizontalSpace(20)
		handler.displayRecentActivities(contentWindow)
	})
}

func (handler *OverviewHandler) displayRecentActivities(contentWindow *widgets.Window) {
	contentWindow.AddLabelWithFont("Recent Activity", widgets.LeftCenterAlign, styles.BoldPageContentFont)

	txns, err := handler.wallet.GetTransactionsRaw(0, 5, dcrlibwallet.TxFilterAll)
	if err != nil {
		handler.err = err
	}

	if len(txns) == 0 {
		contentWindow.AddHorizontalSpace(20)
		contentWindow.AddLabel("No Transaction yet", widgets.CenterAlign)
		return
	}

	historyTable := widgets.NewTable()

	// render table header with nav font
	historyTable.AddRowWithFont(styles.NavFont,
		widgets.NewLabelTableCell("Date", "LC"),
		widgets.NewLabelTableCell("Direction", "LC"),
		widgets.NewLabelTableCell("Amount", "LC"),
		widgets.NewLabelTableCell("Fee", "LC"),
		widgets.NewLabelTableCell("Type", "LC"),
		widgets.NewLabelTableCell("Hash", "LC"),
	)

	for _, tx := range txns {
		historyTable.AddRow(
			widgets.NewLabelTableCell(dcrlibwallet.ExtractDateOrTime(tx.Timestamp), "LC"),
			widgets.NewLabelTableCell(dcrlibwallet.TransactionDirectionName(tx.Direction), "LC"),
			widgets.NewLabelTableCell(dcrutil.Amount(tx.Amount).String(), "RC"),
			widgets.NewLabelTableCell(dcrutil.Amount(tx.Fee).String(), "RC"),
			widgets.NewLabelTableCell(tx.Type, "LC"),
			widgets.NewLabelTableCell(tx.Hash, "LC"),
		)
	}

	historyTable.Render(contentWindow)
}
