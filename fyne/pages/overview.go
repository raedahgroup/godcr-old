package pages

import (
	"context"
	"fmt"

	"fyne.io/fyne/widget"
	"github.com/decred/dcrd/dcrutil"
	"github.com/raedahgroup/godcr/app/walletcore"
	"github.com/raedahgroup/godcr/fyne/widgets"
)

type overview struct {
	err           error
	balanceWidget *widget.Label
	view          *widgets.Box
}

type recentActivity struct {
	err          error
	isFetching   bool
	view         *widgets.Box
	transactions []*walletcore.Transaction
}

type OverviewHandler struct {
	err error

	ctx       context.Context
	wallet    walletcore.Wallet
	container *widgets.Box

	overview       *overview
	recentActivity *recentActivity
}

func (handler *OverviewHandler) Render(ctx context.Context, wallet walletcore.Wallet, container *widgets.Box) {
	handler.ctx = ctx
	handler.wallet = wallet
	handler.container = container
	handler.overview = &overview{
		view: widgets.NewVBox(),
	}
	handler.recentActivity = &recentActivity{
		view: widgets.NewVBox(),
	}

	accounts, err := wallet.AccountsOverview(walletcore.DefaultRequiredConfirmations)
	if err != nil {
		handler.overview.err = err
	} else {
		var totalBalance dcrutil.Amount
		for _, account := range accounts {
			totalBalance += account.Balance.Total
		}
		handler.overview.balanceWidget = widget.NewLabel(totalBalance.String())
	}
	handler.renderOverview()

	// start fetching activity in background
	go handler.fetchRecentActivity()
}

func (handler *OverviewHandler) renderOverview() {
	if handler.overview.err != nil {
		handler.overview.view.AddLabel(handler.overview.err.Error())
	} else {
		handler.overview.view.AddBoldLabel("Current Total Balance")
		handler.overview.view.Add(handler.overview.balanceWidget)
	}

	handler.container.Add(handler.overview.view)
	handler.container.Add(widgets.NewVSpacer(20))
}

func (handler *OverviewHandler) fetchRecentActivity() {
	handler.recentActivity.isFetching = true
	handler.renderRecentActivity()

	handler.recentActivity.transactions, handler.recentActivity.err = handler.wallet.TransactionHistory(0, 5, nil)
	handler.recentActivity.isFetching = false
	handler.renderRecentActivity()
}

func (handler *OverviewHandler) renderRecentActivity() {
	handler.recentActivity.view.Empty()
	handler.recentActivity.view.AddBoldLabel("Recent Activity")

	if handler.recentActivity.isFetching {
		handler.recentActivity.view.AddItalicLabel("Fetching recent activity...")
	} else if handler.recentActivity.err != nil {
		handler.recentActivity.view.AddItalicLabel(handler.recentActivity.err.Error())
	} else {
		table := widgets.NewTable()
		table.AddRowSimple("#", "Date", "Direction", "Amount", "Fee", "Type", "Hash")
		for index, txn := range handler.recentActivity.transactions {
			table.AddRow(
				widget.NewLabel(fmt.Sprintf("%d", index+1)),
				widget.NewLabel(txn.ShortTime),
				widget.NewLabel(txn.Direction.String()),
				widget.NewLabel(dcrutil.Amount(txn.Amount).String()),
				widget.NewLabel(dcrutil.Amount(txn.Fee).String()),
				widget.NewLabel(txn.Type),
				widgets.NewLink(txn.Hash, func() {
					NewTransactionDetailsHandler(txn.Hash, handler.ctx, handler.wallet).Render("Overview", handler.container)
				}),
			)
		}
		handler.recentActivity.view.Add(table.CondensedTable())
	}

	if len(handler.container.Children) > 3 {
		handler.container.Children[3] = handler.recentActivity.view.Box
		handler.container.Update()
	} else {
		handler.container.Add(handler.recentActivity.view)
	}
}
