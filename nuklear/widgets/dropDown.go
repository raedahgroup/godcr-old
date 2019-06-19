package widgets

import (
	"fmt"
	"github.com/raedahgroup/godcr/app/walletcore"
	// "github.com/raedahgroup/dcrlibwallet/txindex"
)

type FilterSelector struct {
	filters                  []string
	filter                   string
	txCount                  int
	transactionCountByFilter []string
	txCountErr               error
	selectedFilterIndex      int
	totalTxCount             int
}

const (
	defaultFilterSelectorWidth = 120
	filterSelectorHeight       = 25
)

func FilterSelectorWidget(wallet walletcore.Wallet) (filterSelector *FilterSelector) {
	filterSelector = &FilterSelector{}

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
