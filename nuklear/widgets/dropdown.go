package widgets

import (
	"fmt"
	"strings"

	"github.com/aarzilli/nucular/label"
	"github.com/raedahgroup/dcrlibwallet"
)

type FilterSelector struct {
	wallet                 *dcrlibwallet.LibWallet
	filterSelectionOptions []string
	txCountForAllFilters   map[string]int
	selectedFilterIndex    int
	selectionChanged       func()
}

const (
	defaultFilterSelectorWidth = 120
	filterSelectorHeight       = 25
)

var allTxFilterNames = []string{"All", "Sent", "Received", "Transferred", "Coinbase", "Staking"}
var allTxFilters = map[string]int32{
	"All":         dcrlibwallet.TxFilterAll,
	"Sent":        dcrlibwallet.TxFilterSent,
	"Received":    dcrlibwallet.TxFilterReceived,
	"Transferred": dcrlibwallet.TxFilterTransferred,
	"Coinbase":    dcrlibwallet.TxFilterCoinBase,
	"Staking":     dcrlibwallet.TxFilterStaking,
}

func FilterSelectorWidget(wallet *dcrlibwallet.LibWallet, selectionChanged func()) (*FilterSelector, error) {
	filterSelector := &FilterSelector{
		selectionChanged:     selectionChanged,
		wallet:               wallet,
		txCountForAllFilters: make(map[string]int),
	}

	for _, filterName := range allTxFilterNames {
		filterId := allTxFilters[filterName]
		txCountForFilter, txCountErr := wallet.CountTransactions(filterId)
		if txCountErr != nil {
			return nil, fmt.Errorf("error counting tx for filter %s: %s", filterName, txCountErr.Error())
		}
		if txCountForFilter > 0 {
			filterWithCount := fmt.Sprintf("%s (%d)", filterName, txCountForFilter)
			filterSelector.filterSelectionOptions = append(filterSelector.filterSelectionOptions, filterWithCount)
			filterSelector.txCountForAllFilters[filterName] = txCountForFilter
		}
	}

	return filterSelector, nil
}

func (filterSelector *FilterSelector) Render(window *Window) {
	prompt := "Filter"

	// row with fixed column widths to hold filter selection prompt and the actual dropdown widget
	window.Row(filterSelectorHeight).Static(window.LabelWidth(prompt), defaultFilterSelectorWidth)

	// print account selection prompt / label to first column
	window.Label("filter", LeftCenterAlign)

	// render actual dropdown to second column
	filterSelector.makeDropDown(window)
}

// makeDropDown is adapted from nucular's Window.ComboSimple
// to allow triggering a callback when dropdown selection changes.
func (filterSelector *FilterSelector) makeDropDown(window *Window) {
	if len(filterSelector.filterSelectionOptions) == 0 {
		return
	}

	items := filterSelector.filterSelectionOptions
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

// GetSelectedFilter returns the id of the selected filter and the tx count for the filter.
func (filterSelector *FilterSelector) GetSelectedFilter() (int32, int) {
	selectedOption := filterSelector.filterSelectionOptions[filterSelector.selectedFilterIndex]
	selectedFilterName := strings.Split(selectedOption, " ")[0]
	return allTxFilters[selectedFilterName], filterSelector.txCountForAllFilters[selectedFilterName]
}
