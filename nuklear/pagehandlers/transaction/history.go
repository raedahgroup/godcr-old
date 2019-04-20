package transaction

import (
	"context"
	"fmt"

	"github.com/aarzilli/nucular"
	"github.com/raedahgroup/godcr/app/walletcore"
	"github.com/raedahgroup/godcr/nuklear/styles"
	"github.com/raedahgroup/godcr/nuklear/widgets"
)

type paginationData struct {
	itemsPerPage    int
	nextBlockHeight int32
	endBlockHeight  int32
}

type HistoryHandler struct {
	fetchHistoryError      error
	ctx                    context.Context
	transactions           []*walletcore.Transaction
	isFetchingTransactions bool
	isRendering            bool
	paginationData         paginationData

	selectedTxHash      string
	selectedTxDetails   *walletcore.TransactionDetails
	isFetchingTxDetails bool
	fetchTxDetailsError error

	wallet walletcore.Wallet
}

func (handler *HistoryHandler) BeforeRender(wallet walletcore.Wallet, refreshWindowDisplay func()) bool {
	// todo: caller should ideally pass a context parameter, propagated from main.go
	handler.ctx = context.Background()

	handler.clearTxDetails()

	handler.fetchHistoryError = nil
	handler.transactions = nil
	handler.isRendering = false

	handler.wallet = wallet

	handler.paginationData = paginationData{
		nextBlockHeight: -1,
		endBlockHeight:  0,
		itemsPerPage:    walletcore.TransactionHistoryCountPerPage,
	}

	return true
}

func (handler *HistoryHandler) Render(window *nucular.Window) {
	if handler.selectedTxHash == "" {
		if !handler.isRendering {
			handler.isRendering = true
			go handler.fetchTransactions(window)
		}
		handler.renderHistoryPage(window)
		return
	}
	handler.renderTransactionDetailsPage(window)
}

func (handler *HistoryHandler) renderHistoryPage(window *nucular.Window) {
	widgets.PageContentWindowDefaultPadding("History", window, func(contentWindow *widgets.Window) {
		if handler.fetchHistoryError != nil {
			contentWindow.DisplayErrorMessage("Error fetching txs", handler.fetchHistoryError)
		} else if len(handler.transactions) > 0 {
			handler.displayTransactions(contentWindow)
		}

		// show loading indicator if tx is being fetched
		if handler.isFetchingTransactions {
			contentWindow.DisplayIsLoadingMessage()
		}
	})
}

func (handler *HistoryHandler) fetchTransactions(window *nucular.Window) {
	handler.isFetchingTransactions = true
	window.Master().Changed()

	transactions, endBlockHeight, err := handler.wallet.TransactionHistory(handler.ctx, handler.paginationData.nextBlockHeight,
		handler.paginationData.itemsPerPage)

	handler.fetchHistoryError = err
	handler.transactions = append(handler.transactions, transactions...)

	handler.paginationData.endBlockHeight = endBlockHeight
	handler.paginationData.nextBlockHeight = endBlockHeight - 1

	window.Master().Changed()

	handler.isFetchingTransactions = false
}

func (handler *HistoryHandler) displayTransactions(contentWindow *widgets.Window) {
	historyTable := widgets.NewTable()

	// render table header with nav font
	historyTable.AddRowWithFont(styles.NavFont,
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
			widgets.NewLinkTableCell(tx.Hash, "Click to see transaction details", handler.gotoTransactionDetails),
		)
	}

	historyTable.Render(contentWindow)

	if !handler.isFetchingTransactions {
		contentWindow.Row(40).Static(130)
		contentWindow.AddButtonToCurrentRow("Load more", func() {
			handler.gotoNextPage(contentWindow)
		})
	}
}

func (handler *HistoryHandler) gotoNextPage(window *widgets.Window) {
	go handler.fetchTransactions(window.Window)
}
