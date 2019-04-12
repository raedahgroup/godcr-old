package pagehandlers

import (
	"context"
	"fmt"

	"github.com/aarzilli/nucular"
	"github.com/raedahgroup/godcr/app/walletcore"
	"github.com/raedahgroup/godcr/nuklear/helpers"
	"github.com/raedahgroup/godcr/nuklear/widgets"
)

type HistoryHandler struct {
	err                    error
	ctx                    context.Context
	transactions           []*walletcore.Transaction
	isFetchingTransactions bool
	nextBlockHeight        int32
}

func (handler *HistoryHandler) BeforeRender(wallet walletcore.Wallet, refreshWindowDisplay func()) bool {
	// todo: caller should ideally pass a context parameter, propagated from main.go
	handler.ctx = context.Background()

	handler.isFetchingTransactions = true
	handler.err = nil
	handler.transactions = nil

	go handler.fetchTransactions(wallet, refreshWindowDisplay)

	return true
}

func (handler *HistoryHandler) Render(window *nucular.Window) {
	widgets.PageContentWindowDefaultPadding("History", window, func(contentWindow *widgets.Window) {
		if handler.err != nil {
			contentWindow.DisplayErrorMessage("Error fetching txs", handler.err)
		} else if len(handler.transactions) > 0 {
			handler.displayTransactions(contentWindow)
		}

		// show loading indicator if tx is being fetched
		if handler.isFetchingTransactions {
			contentWindow.DisplayIsLoadingMessage()
		}
	})
}

func (handler *HistoryHandler) fetchTransactions(wallet walletcore.Wallet, refreshWindowDisplay func()) {
	if len(handler.transactions) == 0 {
		// first page
		handler.nextBlockHeight = -1
	}

	transactions, endBlockHeight, err := wallet.TransactionHistory(handler.ctx, handler.nextBlockHeight,
		walletcore.TransactionHistoryCountPerPage)

	// next start block should be the block immediately preceding the current end block
	handler.err = err
	handler.transactions = append(handler.transactions, transactions...)
	handler.nextBlockHeight = endBlockHeight - 1

	refreshWindowDisplay()

	// load more if possible
	if handler.nextBlockHeight >= 0 {
		handler.fetchTransactions(wallet, refreshWindowDisplay)
	} else {
		handler.isFetchingTransactions = false
	}
}

func (handler *HistoryHandler) displayTransactions(contentWindow *widgets.Window) {
	historyTable := widgets.NewTable()

	// render table header with nav font
	historyTable.AddRowWithFont(helpers.NavFont,
		widgets.NewLabelTableCell("#", "LC"),
		widgets.NewLabelTableCell("Date", "LC"),
		widgets.NewLabelTableCell("Direction", "LC"),
		widgets.NewLabelTableCell("Amount", "LC"),
		widgets.NewLabelTableCell("Fee", "LC"),
		widgets.NewLabelTableCell("Type", "LC"),
		widgets.NewLabelTableCell("Hash", "LC"),
	)

	for i, tx := range handler.transactions {
		historyTable.AddRow(
			widgets.NewLabelTableCell(fmt.Sprintf("%d", i+1), "LC"),
			widgets.NewLabelTableCell(tx.FormattedTime, "LC"),
			widgets.NewLabelTableCell(tx.Direction.String(), "LC"),
			widgets.NewLabelTableCell(tx.Amount, "RC"),
			widgets.NewLabelTableCell(tx.Fee, "RC"),
			widgets.NewLabelTableCell(tx.Type, "LC"),
			widgets.NewLabelTableCell(tx.Hash, "LC"),
		)
	}

	historyTable.Render(contentWindow)
}
