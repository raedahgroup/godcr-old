package historypagehandler

import (
	"fmt"
	"strings"

	"fyne.io/fyne/widget"

	"github.com/raedahgroup/dcrlibwallet"
	"github.com/raedahgroup/godcr/fyne/assets"
	"github.com/raedahgroup/godcr/fyne/widgets"
)

func (historyPage *HistoryPageData) txFilterDropDown(walletId int) {
	var txTable widgets.Table
	txFilterListWidget := widget.NewVBox()

	var allTxFilterNames = []string{"All", "Sent", "Received", "Transferred", "Coinbase", "Staking"}
	var allTxFilters = map[string]int32{
		"All":         dcrlibwallet.TxFilterAll,
		"Sent":        dcrlibwallet.TxFilterSent,
		"Received":    dcrlibwallet.TxFilterReceived,
		"Transferred": dcrlibwallet.TxFilterTransferred,
		"Coinbase":    dcrlibwallet.TxFilterCoinBase,
		"Staking":     dcrlibwallet.TxFilterStaking,
	}

	if walletId != historyPage.selectedWalletID {
		historyPage.selectedWalletID = walletId
	}

	txCountForFilter, err := historyPage.MultiWallet.WalletWithID(historyPage.selectedWalletID).CountTransactions(allTxFilters["All"])
	if err != nil {
		historyPage.errorMessage = fmt.Sprintf("Cannot load history page page. Error getting transaction count for filter All: %s", err.Error())
		return
	}

	historyPage.allTxCount = txCountForFilter

	historyPage.selectedTxFilterLabel.SetText(fmt.Sprintf("%s (%d)", "All", txCountForFilter))

	for _, filterName := range allTxFilterNames {
		filterId := allTxFilters[filterName]
		txCountForFilter, err := historyPage.MultiWallet.WalletWithID(historyPage.selectedWalletID).CountTransactions(filterId)
		if err != nil {
			historyPage.errorMessage = fmt.Sprintf("Cannot load historyPage page. Error getting transaction count for filter %s: %s", filterName, err.Error())
			return
		}
		if txCountForFilter > 0 {
			filter := fmt.Sprintf("%s (%d)", filterName, txCountForFilter)
			txFilterView := widget.NewHBox(
				widgets.NewHSpacer(5),
				widget.NewLabel(filter),
				widgets.NewHSpacer(5),
			)

			txFilterListWidget.Append(widgets.NewClickableBox(txFilterView, func() {
				selectedFilterName := strings.Split(filter, " ")[0]
				selectedFilterId := allTxFilters[selectedFilterName]
				if allTxCountForSelectedTx, err := historyPage.MultiWallet.WalletWithID(historyPage.selectedWalletID).CountTransactions(selectedFilterId); err == nil {
					historyPage.allTxCount = allTxCountForSelectedTx
				}

				if selectedFilterId != historyPage.selectedFilterId {
					historyPage.selectedTxFilterLabel.SetText(filter)
					historyPage.txTableHeader(&txTable)
					historyPage.fetchTx(&txTable, 0, selectedFilterId, false)
					widget.Refresh(historyPage.txTable.Result)
				}

				historyPage.txFilterSelectionPopup.Hide()
			}))
		}
	}

	// txFilterSelectionPopup create a popup that has tx filter name and tx count
	historyPage.txFilterSelectionPopup = widget.NewPopUp(widget.NewVBox(txFilterListWidget), historyPage.Window.Canvas())
	historyPage.txFilterSelectionPopup.Hide()

	historyPage.txFilterTab = widget.NewHBox(
		historyPage.selectedTxFilterLabel,
		widgets.NewHSpacer(10),
		widget.NewIcon(historyPage.icons[assets.CollapseIcon]),
		widgets.NewHSpacer(10),
	)
	widget.Refresh(historyPage.txFilterTab)
}

func (historyPage *HistoryPageData) txSortDropDown() {
	var txTable widgets.Table
	var allTxSortNames = []string{"Newest", "Oldest"}
	var allTxSortFilters = map[string]bool{
		"Newest": true,
		"Oldest": false,
	}

	historyPage.selectedTxSortFilterLabel.SetText("Newest")

	historyPage.selectedtxSort = allTxSortFilters["Newest"]

	txSortFilterListWidget := widget.NewVBox()
	for _, sortName := range allTxSortNames {
		txSortView := widget.NewHBox(
			widgets.NewHSpacer(5),
			widget.NewLabel(sortName),
			widgets.NewHSpacer(5),
		)
		txSort := allTxSortFilters[sortName]
		newSortName := sortName

		txSortFilterListWidget.Append(widgets.NewClickableBox(txSortView, func() {
			historyPage.selectedTxSortFilterLabel.SetText(newSortName)
			historyPage.selectedtxSort = txSort

			historyPage.txTableHeader(&txTable)
			historyPage.txTable.Result.Children = txTable.Result.Children
			historyPage.fetchTx(&txTable, 0, historyPage.selectedFilterId, false)
			widget.Refresh(historyPage.txTable.Result)
			historyPage.txSortFilterSelectionPopup.Hide()
		}))
	}

	// txSortFilterSelectionPopup create a popup that has tx filter name and tx count
	historyPage.txSortFilterSelectionPopup = widget.NewPopUp(widget.NewVBox(txSortFilterListWidget), historyPage.Window.Canvas())
	historyPage.txSortFilterSelectionPopup.Hide()

	historyPage.txSortFilterTab = widget.NewHBox(
		historyPage.selectedTxSortFilterLabel,
		widgets.NewHSpacer(10),
		widget.NewIcon(historyPage.icons[assets.CollapseIcon]),
	)
	widget.Refresh(historyPage.txSortFilterTab)
}
