package widgets

import (
	"fmt"
	"strings"

	"github.com/aarzilli/nucular/label"
	"github.com/raedahgroup/godcr/app/walletcore"
)

type FilterSelector struct {
	wallet                 walletcore.Wallet
	filterSelectionOptions []string
	txCountForAllFilters   map[string]int
	selectedFilterIndex    int
	selectionChanged       func()
}

const (
	defaultFilterSelectorWidth = 120
	filterSelectorHeight       = 25
)

func FilterSelectorWidget(wallet walletcore.Wallet, selectionChanged func()) (*FilterSelector, error) {
	filterSelector := &FilterSelector{
		selectionChanged:     selectionChanged,
		wallet:               wallet,
		txCountForAllFilters: make(map[string]int),
	}

	for _, filter := range walletcore.TransactionFilters {
		txCount, txCountErr := wallet.TransactionCount(walletcore.BuildTransactionFilter(filter))
		if txCountErr != nil {
			return nil, fmt.Errorf("error counting tx for filter %s: %s", filter, txCountErr.Error())
		}
		if txCount > 0 {
			filterWithCount := fmt.Sprintf("%s (%d)", filter, txCount)
			filterSelector.filterSelectionOptions = append(filterSelector.filterSelectionOptions, filterWithCount)
			filterSelector.txCountForAllFilters[filter] = txCount
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

// GetSelectedFilter returns the name of the selected filter and the tx count for the filter.
func (filterSelector *FilterSelector) GetSelectedFilter() (string, int) {
	selectedOption := filterSelector.filterSelectionOptions[filterSelector.selectedFilterIndex]
	selectedFilterName := strings.Split(selectedOption, " ")[0]
	return selectedFilterName, filterSelector.txCountForAllFilters[selectedFilterName]
}
