package pages

import (
	"fmt"
	// "image/color"

	"fyne.io/fyne"
	"fyne.io/fyne/widget"
	// "fyne.io/fyne/canvas"

	"github.com/raedahgroup/dcrlibwallet"
	"github.com/raedahgroup/godcr/fyne/widgets"
	"github.com/raedahgroup/godcr/fyne/assets"
)


func HistoryPageContent(wallet *dcrlibwallet.LibWallet, window fyne.Window, tabmenu *widget.TabContainer) fyne.CanvasObject {
		// error handler
	var errorLabel *widget.Label
	errorLabel = widget.NewLabelWithStyle("", fyne.TextAlignCenter, fyne.TextStyle{Bold: true})
	errorLabel.Hide()

	pageTitleLabel := widget.NewLabelWithStyle("Transactions", fyne.TextAlignLeading, fyne.TextStyle{Bold: true, Italic: true})
	
	t := prepareTxFilterDropDown(wallet, window, errorLabel)
	output := widget.NewVBox(
		widgets.NewVSpacer(5),
		widget.NewHBox(pageTitleLabel),
		widgets.NewVSpacer(5),
		t,
		errorLabel,
	)	

	return widget.NewHBox(widgets.NewHSpacer(18), output)
}

func prepareTxFilterDropDown(wallet *dcrlibwallet.LibWallet, window fyne.Window, errorLabel *widget.Label) *widgets.ClickableBox {
	var allTxFilterNames = []string{"All", "Sent", "Received", "Transferred", "Coinbase", "Staking"}
	var allTxFilters = map[string]int32{
		"All":         dcrlibwallet.TxFilterAll,
		"Sent":        dcrlibwallet.TxFilterSent,
		"Received":    dcrlibwallet.TxFilterReceived,
		"Transferred": dcrlibwallet.TxFilterTransferred,
		"Coinbase":    dcrlibwallet.TxFilterCoinBase,
		"Staking":     dcrlibwallet.TxFilterStaking,
	}

	txCountForFilter, _ := wallet.CountTransactions(allTxFilters["All"])
	selectedAccountLabel := widget.NewLabel(fmt.Sprintf("%s (%d)", "All", txCountForFilter))

	activeFiltersWithTxCounts := make(map[int32]int)

	var accountSelectionPopup *widget.PopUp
	var accountsView *widget.Box
	accountListWidget := widget.NewVBox()
	for _, filterName := range allTxFilterNames {
		filterId := allTxFilters[filterName]
		txCountForFilter, txCountErr := wallet.CountTransactions(filterId)
		if txCountErr != nil {
			errorMessage := fmt.Sprintf("Cannot load history page. Error getting transaction count for filter %s: %s",
				filterName, txCountErr.Error())
			errorHandler(errorMessage, errorLabel)
			return nil
		}

		if txCountForFilter > 0 {
			activeFiltersWithTxCounts[filterId] = txCountForFilter
			filter := fmt.Sprintf("%s (%d)", filterName, txCountForFilter)
			accountsView = widget.NewHBox(
				widgets.NewHSpacer(5),
				widget.NewLabel(filter),
				widgets.NewHSpacer(5),
			)

				accountListWidget.Append(widgets.NewClickableBox(accountsView, func() {
		selectedAccountLabel.SetText(filter)
		accountSelectionPopup.Hide()
	}))
		}
	}



	// accountSelectionPopup create a popup that has account names with spendable amount
	accountSelectionPopup = widget.NewPopUp(
		widget.NewVBox(
			accountListWidget,
		), window.Canvas(),
	)
	accountSelectionPopup.Hide()



	// // dropDown selection change listener
	// txFilterDropDown.AddDropDown("", txFilterSelectionOptions, 0, func(selectedOption string, index int) {
	// 	selectedFilterName := strings.Split(selectedOption, " ")[0]
	// 	selectedFilterId := allTxFilters[selectedFilterName]
	// 	if selectedFilterId != historyPageData.currentTxFilter {
	// 		go fetchAndDisplayTransactions(0, selectedFilterId)
	// 	}
	// })

	// accountTab shows the selected account
	icons, _ := assets.GetIcons(assets.CollapseIcon)
	accountTab := widget.NewHBox(
		selectedAccountLabel,
		widgets.NewHSpacer(8),
		widget.NewIcon(icons[assets.CollapseIcon]),
	)

	var accountDropdown *widgets.ClickableBox
	accountDropdown = widgets.NewClickableBox(accountTab, func() {
		accountSelectionPopup.Move(fyne.CurrentApp().Driver().AbsolutePositionForObject(
			accountDropdown).Add(fyne.NewPos(0, accountDropdown.Size().Height)))
		accountSelectionPopup.Show()
	})

	return accountDropdown
}

// func errorHandler(err string, errorLabel *widget.Label) {
// 	errorLabel.SetText(err)
// 	errorLabel.Show()
// }
