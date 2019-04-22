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
	currentPage     int
	itemsPerPage    int
	nextBlockHeight int32
	endBlockHeight  int32
}

type transactionsData struct {
	numberingStartsAt int
	transactions      []*walletcore.Transaction
}

type HistoryHandler struct {
	fetchHistoryError      error
	ctx                    context.Context
	transactions           []*transactionsData // each []*walletcore.Transaction slice is a page
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

	handler.wallet = wallet

	handler.paginationData = paginationData{
		currentPage:     0,
		nextBlockHeight: -1,
		endBlockHeight:  0,
		itemsPerPage:    walletcore.TransactionHistoryCountPerPage,
	}

	handler.clearTxDetails()

	handler.fetchHistoryError = nil
	handler.transactions = nil
	handler.isRendering = false

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

func (handler *HistoryHandler) fetchTransactions(window *nucular.Window) {
	handler.isFetchingTransactions = true
	window.Master().Changed()

	transactions, endBlockHeight, err := handler.wallet.TransactionHistory(handler.ctx, handler.paginationData.nextBlockHeight,
		handler.paginationData.itemsPerPage)

	numberingStartsAt := 1
	if len(handler.transactions) > 0 {
		lastSavedTransaction := handler.transactions[len(handler.transactions)-1]
		numberingStartsAt = lastSavedTransaction.numberingStartsAt + len(lastSavedTransaction.transactions)
	}

	transactionsData := &transactionsData{
		numberingStartsAt: numberingStartsAt,
		transactions:      transactions,
	}

	handler.fetchHistoryError = err
	handler.transactions = append(handler.transactions, transactionsData)

	handler.paginationData.endBlockHeight = endBlockHeight
	handler.paginationData.nextBlockHeight = endBlockHeight - 1
	handler.paginationData.currentPage += 1

	window.Master().Changed()

	handler.isFetchingTransactions = false
}

func (handler *HistoryHandler) renderHistoryPage(window *nucular.Window) {
	widgets.PageContentWindowDefaultPadding("History", window, func(contentWindow *widgets.Window) {
		if handler.isFetchingTransactions {
			contentWindow.DisplayIsLoadingMessage()
			return
		}

		if handler.fetchHistoryError != nil {
			contentWindow.DisplayErrorMessage("Error fetching txs", handler.fetchHistoryError)
		} else if len(handler.transactions) > 0 {
			handler.displayTransactions(contentWindow)
		}
	})
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

	// get current page transactions
	transactions := handler.transactions[handler.paginationData.currentPage-1]
	currentNumber := transactions.numberingStartsAt

	for _, tx := range transactions.transactions {
		historyTable.AddRow(
			widgets.NewLabelTableCell(fmt.Sprintf("%d", currentNumber), "LC"),
			widgets.NewLabelTableCell(tx.FormattedTime, "LC"),
			widgets.NewLabelTableCell(tx.Direction.String(), "LC"),
			widgets.NewLabelTableCell(tx.Amount, "RC"),
			widgets.NewLabelTableCell(tx.Fee, "RC"),
			widgets.NewLabelTableCell(tx.Type, "LC"),
			widgets.NewLinkTableCell(tx.Hash, "Click to see transaction details", handler.gotoTransactionDetails),
		)
		currentNumber++
	}

	historyTable.Render(contentWindow)

	if !handler.isFetchingTransactions {
		contentWindow.Row(40).Static(130, 130)
		// show previous button only if current page is greater than 1
		if handler.paginationData.currentPage > 1 {
			contentWindow.AddButtonToCurrentRow("Previous", func() {
				handler.loadPreviousPage(contentWindow)
			})
		}

		contentWindow.AddButtonToCurrentRow("Next Page", func() {
			handler.loadNextPage(window)
		})
	}
}

func (handler *HistoryHandler) loadNextPage(window *widgets.Window) {
	// check if transactions for the page we are navigating to is already loaded
	if len(handler.transactions) >= handler.paginationData.currentPage+1 {
		handler.paginationData.currentPage++
		window.Master().Changed()
		return
	}
	go handler.fetchTransactions(window.Window)
}

func (handler *HistoryHandler) loadPreviousPage(window *widgets.Window) {
	defer window.Master().Changed()

	// perform this check even if it's unlikely to occur
	if handler.paginationData.currentPage == 1 {
		return
	}
	handler.paginationData.currentPage -= 1
}
