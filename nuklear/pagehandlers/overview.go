package pagehandlers

import (
	"github.com/aarzilli/nucular"
	"github.com/decred/dcrd/dcrutil"
	"github.com/raedahgroup/godcr/app/config"
	"github.com/raedahgroup/godcr/app/walletcore"
	"github.com/raedahgroup/godcr/nuklear/styles"
	"github.com/raedahgroup/godcr/nuklear/widgets"
)

type OverviewHandler struct {
	err      error
	accounts []*walletcore.Account
	wallet   walletcore.Wallet
}

func (handler *OverviewHandler) BeforeRender(wallet walletcore.Wallet, settings *config.Settings, _ func()) bool {
	handler.wallet = wallet
	handler.accounts, handler.err = wallet.AccountsOverview(walletcore.DefaultRequiredConfirmations)
	return true
}

func (handler *OverviewHandler) Render(window *nucular.Window) {
	widgets.PageContentWindowDefaultPadding("Overview", window, func(contentWindow *widgets.Window) {
		contentWindow.AddLabelWithFont("Current Total Balance", widgets.LeftCenterAlign, styles.BoldPageContentFont)

		if handler.err != nil {
			contentWindow.DisplayErrorMessage("Error fetching accounts balance", handler.err)
		} else {
			contentWindow.AddLabel(walletcore.WalletBalance(handler.accounts), widgets.LeftCenterAlign)
			contentWindow.AddHorizontalSpace(20)
			handler.displayRecentActivities(contentWindow)
		}
	})
}

func (handler *OverviewHandler) displayRecentActivities(contentWindow *widgets.Window) {
	contentWindow.AddLabelWithFont("Recent Activity", widgets.LeftCenterAlign, styles.BoldPageContentFont)

	txns, err := handler.wallet.TransactionHistory(0, 5, nil)
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
			widgets.NewLabelTableCell(tx.ShortTime, "LC"),
			widgets.NewLabelTableCell(tx.Direction.String(), "LC"),
			widgets.NewLabelTableCell(dcrutil.Amount(tx.Amount).String(), "RC"),
			widgets.NewLabelTableCell(dcrutil.Amount(tx.Fee).String(), "RC"),
			widgets.NewLabelTableCell(tx.Type, "LC"),
			widgets.NewLabelTableCell(tx.Hash, "LC"),
		)
	}

	historyTable.Render(contentWindow)
}
