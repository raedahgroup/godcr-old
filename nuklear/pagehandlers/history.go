package pagehandlers

import (
	"fmt"
	// "strings"

	"github.com/aarzilli/nucular"
	"github.com/decred/dcrd/dcrutil"
	"github.com/raedahgroup/godcr/app/config"
	"github.com/raedahgroup/godcr/app/walletcore"
	"github.com/raedahgroup/godcr/nuklear/styles"
	"github.com/raedahgroup/godcr/nuklear/widgets"
)

type HistoryHandler struct {
	wallet               walletcore.Wallet
	refreshWindowDisplay func()

	filterSelectorWidget *widgets.FilterSelector
	// selectedFilter string
	filter *txindex.ReadFilter

	totalTxCount int

	currentPage            int
	txPerPage              int
	transactions           []*walletcore.Transaction
	fetchHistoryError      error
	isFetchingTransactions bool

	selectedTxHash      string
	selectedTxDetails   *walletcore.Transaction
	isFetchingTxDetails bool
	fetchTxDetailsError error
}

func (handler *HistoryHandler) BeforeRender(wallet walletcore.Wallet, settings *config.Settings, refreshWindowDisplay func()) bool {
	handler.wallet = wallet
	handler.refreshWindowDisplay = refreshWindowDisplay

	handler.filterSelectorWidget = widgets.FilterSelectorWidget(wallet)

	handler.currentPage = 1
	handler.txPerPage = walletcore.TransactionHistoryCountPerPage
	handler.transactions = nil
	handler.fetchHistoryError = nil
	handler.isFetchingTransactions = false

	handler.clearTxDetails()

	handler.totalTxCount, handler.fetchHistoryError = wallet.TransactionCount(nil)
	if handler.fetchHistoryError == nil {
		go handler.fetchTransactions(handler.filter)
	}

	return true
}

func (handler *HistoryHandler) fetchTransactions(filter *txindex.ReadFilter) {
	handler.isFetchingTransactions = true
	handler.refreshWindowDisplay() // refresh display to show loading indicator
	txHistoryOffset := 0

	if handler.transactions != nil {
		txHistoryOffset = len(handler.transactions)
		fmt.Println(txHistoryOffset)

	}
	transactions, err := handler.wallet.TransactionHistory(int32(txHistoryOffset), int32(handler.txPerPage), filter)
	fmt.Println(len(transactions))

	handler.fetchHistoryError = err
	handler.transactions = append(handler.transactions, transactions...)

	handler.isFetchingTransactions = false
	handler.refreshWindowDisplay()
}

func (handler *HistoryHandler) Render(window *nucular.Window) {
	if handler.selectedTxHash == "" {
		handler.renderHistoryPage(window)
		return
	}
	handler.renderTransactionDetailsPage(window)
}

func (handler *HistoryHandler) renderHistoryPage(window *nucular.Window) {
	widgets.PageContentWindowDefaultPadding("History", window, func(contentWindow *widgets.Window) {
		handler.filterSelectorWidget.Render(contentWindow)

		handler.totalTxCount, handler.filter = handler.filterSelectorWidget.GetSelectedFilter()

		// show transactions first, if any
		if len(handler.transactions) > 0 {
			handler.displayTransactions(contentWindow)
		}
		// show fetch error if any
		if handler.fetchHistoryError != nil {
			contentWindow.DisplayErrorMessage("Error loading history", handler.fetchHistoryError)
		}
		// show loading indicator if tx fetching is progress
		if handler.isFetchingTransactions {
			contentWindow.DisplayIsLoadingMessage()
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

	pageTxOffset := (handler.currentPage - 1) * handler.txPerPage
	maxTxIndexForCurrentPage := pageTxOffset + handler.txPerPage
	for currentTxIndex, tx := range handler.transactions {
		if currentTxIndex < pageTxOffset {
			continue // skip txs not belonging to this page
		}
		if currentTxIndex >= maxTxIndexForCurrentPage {
			break // max number of tx displayed for this page
		}

		historyTable.AddRow(
			widgets.NewLabelTableCell(fmt.Sprintf("%d", currentTxIndex+1), "LC"),
			widgets.NewLabelTableCell(tx.ShortTime, "LC"),
			widgets.NewLabelTableCell(tx.Direction.String(), "LC"),
			widgets.NewLabelTableCell(dcrutil.Amount(tx.Amount).String(), "RC"),
			widgets.NewLabelTableCell(dcrutil.Amount(tx.Fee).String(), "RC"),
			widgets.NewLabelTableCell(tx.Type, "LC"),
			widgets.NewLinkTableCell(tx.Hash, "Click to see transaction details", handler.gotoTransactionDetails),
		)
	}
	historyTable.Render(contentWindow)

	if !handler.isFetchingTransactions {
		contentWindow.Row(40).Static(110, 110)

		// show previous button only if current page is greater than 1
		if handler.currentPage > 1 {
			contentWindow.AddButtonToCurrentRow("Previous", func() {
				handler.loadPreviousPage(contentWindow)
			})
		}

		// show next button only if there are more txs to be loaded
		if handler.totalTxCount > maxTxIndexForCurrentPage {
			contentWindow.AddButtonToCurrentRow("Next", func() {
				handler.loadNextPage(contentWindow)
			})
		}
	}
}

func (handler *HistoryHandler) loadPreviousPage(window *widgets.Window) {
	handler.currentPage--
	window.Master().Changed()
}

func (handler *HistoryHandler) loadNextPage(window *widgets.Window) {
	nextPage := handler.currentPage + 1
	handler.currentPage = nextPage

	nextPageTxOffset := (nextPage - 1) * handler.txPerPage
	if nextPageTxOffset >= len(handler.transactions) {
		// we've not loaded txs for this page
		go handler.fetchTransactions(nil)
	}

	handler.refreshWindowDisplay()
}

// func (handler *HistoryHandler) {
// 	selectedFilterAndCount := strings.Split(handler.filterSelectorWidget.GetSelectedFilter(), " ")
// 	handler.selectedFilter = selectedFilterAndCount[0]
// 	fmt.Println(handler.selectedFilter)

// 	if handler.selectedFilter == "All" {
// 		go handler.fetchTransactions(nil)
// 		return false
// 	}

// 	handler.filter = txindex.Filter()
// 	handler.filter = walletcore.BuildTransactionFilter(handler.selectedFilter)

// 	handler.totalTxCount, handler.fetchHistoryError = handler.wallet.TransactionCount(handler.filter)
// 	if handler.fetchHistoryError != nil {
// 		return true
// 	}

// }
