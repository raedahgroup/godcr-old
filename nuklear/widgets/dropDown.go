package widgets

import (
	"fmt"
	"strings"

	"github.com/raedahgroup/dcrlibwallet/txindex"
	"github.com/raedahgroup/godcr/app/walletcore"
	"github.com/aarzilli/nucular/label"
)

type FilterSelector struct {
	wallet walletcore.Wallet
	filters                  []string
	transactionFilter []string
	txCountErr               error
	selectedFilterIndex      int
	totalTxCount             int
	selectedFilter           string
	lastfilter               string
	changed                  bool
	selectedFilterAndCount   string
	txFilter                 *txindex.ReadFilter
	selectionChanged         func()
}

const (
	defaultFilterSelectorWidth = 120
	filterSelectorHeight       = 25
)

func FilterSelectorWidget(wallet walletcore.Wallet, selectionChanged func()) (filterSelector *FilterSelector, error error) {
	filterSelector = &FilterSelector{
		selectionChanged: selectionChanged,
	}

	filterSelector.wallet = wallet
	filterSelector.filters = walletcore.TransactionFilters
	filterSelector.transactionFilter = make([]string, len(filterSelector.filters))

	for index, filter := range filterSelector.filters {
		if filter == "All" {
			txCount, txCountErr := wallet.TransactionCount(nil)
			filterSelector.totalTxCount = txCount
			filterSelector.txCountErr = txCountErr
		}

		txCount, txCountErr := wallet.TransactionCount(walletcore.BuildTransactionFilter(filter))
		filterSelector.txCountErr = txCountErr
		if filterSelector.txCountErr != nil {
			return nil, filterSelector.txCountErr
		}

		filterSelector.transactionFilter[index] = fmt.Sprintf("%s (%d)", filter, txCount)
	}

	return filterSelector, nil
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

	filterSelector.makeDropDown(window)
}

func (filterSelector *FilterSelector) GetSelectedFilter() (int, *txindex.ReadFilter, string, error) {
	if filterSelector.selectedFilterIndex < len(filterSelector.transactionFilter) {
		selectedFilterAndCount := filterSelector.transactionFilter[filterSelector.selectedFilterIndex]

		selectedFilterCount := strings.Split(selectedFilterAndCount, " ")
		filterSelector.selectedFilter = selectedFilterCount[0]

		if filterSelector.selectedFilter == "All" {
			return filterSelector.totalTxCount, nil, "All", nil
		}

		txFilter := txindex.Filter()
		txFilter = walletcore.BuildTransactionFilter(filterSelector.selectedFilter)

		filterSelector.totalTxCount, filterSelector.txCountErr = filterSelector.wallet.TransactionCount(txFilter)
		if filterSelector.txCountErr != nil {
			return filterSelector.totalTxCount, nil, " ", filterSelector.txCountErr
		}

		return filterSelector.totalTxCount, txFilter, filterSelector.selectedFilter, nil
	}

	return filterSelector.totalTxCount, nil, "All", nil
}

// makeDropDown is adapted from nucular's Window.ComboSimple
// to allow triggering a callback when dropdown selection changes.
func (filterSelector *FilterSelector) makeDropDown(window *Window) {
	if len(filterSelector.transactionFilter) == 0 {
		return
	}

	items := filterSelector.transactionFilter
	itemHeight := int(float64(filterSelectorHeight) * window.Master().Style().Scaling)
	itemPadding := window.Master().Style().Combo.ButtonPadding.Y
	maxHeight := (len(items)+1)*itemHeight + itemPadding*3

	if w := window.Combo(label.T(items[filterSelector.selectedFilterIndex]), maxHeight, nil); w != nil {
		w.RowScaled(itemHeight).Dynamic(1)
		for i := range items {
			if w.MenuItem(label.TA(items[i], "LC")) {
				filterSelector.selectedFilterIndex = i
				if filterSelector.selectionChanged != nil {
					filterSelector.selectionChanged()
				}
			}
		}
	}
}

