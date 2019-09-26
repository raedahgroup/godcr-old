package widgets

import (
	"fmt"
	"strings"

	"github.com/raedahgroup/dcrlibwallet/txindex"
	"github.com/raedahgroup/godcr/app/walletcore"
	"github.com/aarzilli/nucular/label"
)

type FilterSelector struct {
	wallet 					walletcore.Wallet
	transactionFilters 		[]string
	txCountErr              error
	selectedFilterIndex     int
	totalTxCount            int
	selectedFilter          string
	lastfilter              string
	changed                 bool
	selectedFilterAndCount  string
	selectedTxFilter        *txindex.ReadFilter
	selectionChanged        func()
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

	for _, filter := range walletcore.TransactionFilters {
		if filter == "All" {
			filterSelector.selectedTxFilter = nil
		} else {
			filterSelector.selectedTxFilter = walletcore.BuildTransactionFilter(filter)
		}

		txCount, txCountErr := wallet.TransactionCount(filterSelector.selectedTxFilter)
		if txCountErr != nil {
			return nil, txCountErr
		}
		if txCount == 0 {
			continue
		}

		filterSelector.transactionFilters = append(filterSelector.transactionFilters, fmt.Sprintf("%s (%d)", filter, txCount))
	}

	return filterSelector, nil
}

func (filterSelector *FilterSelector) Render(window *Window, addColumns ...int) {
	filterSelectorWidth := defaultFilterSelectorWidth

	// row with fixed column widths to hold account selection prompt, the account widget, and any other widgets that may be added later
	rowColumns := []int {
		window.LabelWidth("filter"),
		filterSelectorWidth,
	}
	rowColumns = append(rowColumns, addColumns...)
	window.Row(filterSelectorHeight).Static(rowColumns...)

	// print account selection prompt / label
	window.Label("filter", LeftCenterAlign)

	if filterSelector.txCountErr != nil {
		window.DisplayErrorMessage("Fetch tx error", filterSelector.txCountErr)
	}

	filterSelector.makeDropDown(window)
}

func (filterSelector *FilterSelector) GetSelectedFilter() (string, error) {
	if filterSelector.selectedFilterIndex < len(filterSelector.transactionFilters) {
		selectedFilterAndCount := filterSelector.transactionFilters[filterSelector.selectedFilterIndex]

		selectedFilterCount := strings.Split(selectedFilterAndCount, " ")
		filterSelector.selectedFilter = selectedFilterCount[0]

		if filterSelector.selectedFilter == "All" {
			return "All", nil
		}

		txFilter := walletcore.BuildTransactionFilter(filterSelector.selectedFilter)

		filterSelector.totalTxCount, filterSelector.txCountErr = filterSelector.wallet.TransactionCount(txFilter)
		if filterSelector.txCountErr != nil {
			return " ", filterSelector.txCountErr
		}

		return filterSelector.selectedFilter, nil
	}

	return "All", nil
}

// makeDropDown is adapted from nucular's Window.ComboSimple
// to allow triggering a callback when dropdown selection changes.
func (filterSelector *FilterSelector) makeDropDown(window *Window) {
	if len(filterSelector.transactionFilters) == 0 {
		return
	}

	items := filterSelector.transactionFilters
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

