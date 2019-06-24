package widgets

import (
	"fmt"
	"strings"

	"github.com/raedahgroup/dcrlibwallet/txindex"
	"github.com/raedahgroup/godcr/app/walletcore"
)

type FilterSelector struct {
	wallet walletcore.Wallet

	filters                  []string
	filter                   string
	txCount                  int
	transactionCountByFilter []string
	txCountErr               error
	selectedFilterIndex      int
	totalTxCount             int
	selectedFilter           string
	lastfilter	string
	changed bool
	selectedFilterAndCount   string
	txFilter                 *txindex.ReadFilter
}

const (
	defaultFilterSelectorWidth = 120
	filterSelectorHeight       = 25
)

func FilterSelectorWidget(wallet walletcore.Wallet) (filterSelector *FilterSelector) {
	filterSelector = &FilterSelector{}

	filterSelector.wallet = wallet

	filterSelector.filters = walletcore.TransactionFilters

	filterSelector.transactionCountByFilter = make([]string, len(filterSelector.filters))

	for index, filter := range filterSelector.filters {
		if filter == "All" {
			filterSelector.filter = filter
			filterSelector.txCount, filterSelector.txCountErr = wallet.TransactionCount(nil)
			filterSelector.totalTxCount = filterSelector.txCount
		}

		filterSelector.txCount, filterSelector.txCountErr = wallet.TransactionCount(walletcore.BuildTransactionFilter(filter))
		if filterSelector.txCountErr != nil {
			return
		}

		if filterSelector.txCount == 0 {
			continue
		}
		filterSelector.filter = filter

		filterSelector.transactionCountByFilter[index] = fmt.Sprintf("%s (%d)", filterSelector.filter, filterSelector.txCount)
	}

	return
}

func (filterSelector *FilterSelector) Render(window *Window, addColumns ...int) {
	filterSelectorWidth := defaultFilterSelectorWidth

	// row with fixed column widths to hold account selection prompt, the account widget, and any other widgets that may be added later
	rowColumns := make([]int, 2)
	rowColumns[0] = window.LabelWidth("filter")
	rowColumns[1] = filterSelectorWidth
	rowColumns = append(rowColumns, addColumns...)
	window.Row(filterSelectorHeight).Static(rowColumns...)

	// print account selection prompt / label
	window.Label("filter", LeftCenterAlign)

	if filterSelector.txCountErr != nil {
		window.DisplayErrorMessage("Fetch tx error", filterSelector.txCountErr)
	}

	filterSelector.selectedFilterIndex = window.ComboSimple(filterSelector.transactionCountByFilter,
		filterSelector.selectedFilterIndex, filterSelectorHeight)
}

func (filterSelector *FilterSelector) GetSelectedFilter() (int, *txindex.ReadFilter, string) {
	if filterSelector.selectedFilterIndex < len(filterSelector.transactionCountByFilter) {
		selectedFilterAndCount := filterSelector.transactionCountByFilter[filterSelector.selectedFilterIndex]

		selectedFilterCount := strings.Split(selectedFilterAndCount, " ")
		filterSelector.selectedFilter = selectedFilterCount[0]

		if filterSelector.selectedFilter == "All"{
			return filterSelector.totalTxCount, nil, "All"
		}

		txFilter := txindex.Filter()
		txFilter = walletcore.BuildTransactionFilter(filterSelector.selectedFilter)

		filterSelector.totalTxCount, _ = filterSelector.wallet.TransactionCount(txFilter)
		// if handler.fetchHistoryError != nil {
		// 	// return
		// }

		return filterSelector.totalTxCount, txFilter, filterSelector.selectedFilter
	}

	return filterSelector.totalTxCount, nil, "All"
}

// func (filterSelector *FilterSelector) Reset() {
// 	if len(filterSelector.accounts) > 1 {
// 		filterSelector.selectedAccountIndex = 0
// 	}
// }
