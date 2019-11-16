package pages

import (
	"fmt"
	"strings"
	"strconv"

	"fyne.io/fyne"
	"fyne.io/fyne/widget"
	"fyne.io/fyne/layout"

	"github.com/raedahgroup/dcrlibwallet"
	"github.com/decred/dcrd/dcrutil"
	"github.com/decred/dcrd/chaincfg/chainhash"
	"github.com/raedahgroup/godcr/fyne/widgets"
	"github.com/raedahgroup/godcr/fyne/assets"
)

const txPerPage int32 = 3

type historyPageData struct {
	txTable         widgets.Table
	txDetailsTable widgets.Table
	currentFilter        int32
	currentPage             int32
	selectedFilterCount int
	txCountErr error
	txCountForFilter int
	txns []*dcrlibwallet.Transaction
	selectedFilterId int32
	txl int
}

var history historyPageData


func HistoryPageContent(wallet *dcrlibwallet.LibWallet, window fyne.Window, tabmenu *widget.TabContainer) fyne.CanvasObject {
	// error handler
	var errorLabel *widget.Label
	errorLabel = widget.NewLabelWithStyle("", fyne.TextAlignCenter, fyne.TextStyle{Bold: true})
	errorLabel.Hide()

	history.currentPage = 1
	history.selectedFilterId = dcrlibwallet.TxFilterAll

	var prevButton *widget.Button
	var nextButton *widget.Button
	prevButton = widget.NewButton("Prev", func() {
		loadPreviousPage(wallet , nextButton, prevButton)
	})

	nextButton = widget.NewButton("Next", func() {
		loadNextPage(wallet , nextButton, prevButton)
	})

	pageTitleLabel := widget.NewLabelWithStyle("Transactions", fyne.TextAlignLeading, fyne.TextStyle{Bold: true, Italic: true})
	
	filterDropdown := txFilterDropDown(wallet, window, errorLabel, nextButton, prevButton)

	fetchAndDisplayTransactions(wallet, &history.txTable, nextButton, prevButton)
	
	output := widget.NewVBox(
		widgets.NewVSpacer(5),
		widget.NewHBox(pageTitleLabel),
		widgets.NewVSpacer(5),
		filterDropdown,
		widgets.NewVSpacer(5),
		fyne.NewContainerWithLayout(layout.NewFixedGridLayout(fyne.NewSize(history.txTable.Container.MinSize().Width, history.txTable.Container.MinSize().Height+200)), history.txTable.Container),
		widgets.NewVSpacer(15),
		widget.NewHBox(prevButton, widgets.NewHSpacer(110), nextButton),
		errorLabel,
	)

	return widget.NewHBox(widgets.NewHSpacer(18), output)
}

func txFilterDropDown(wallet *dcrlibwallet.LibWallet, window fyne.Window, errorLabel *widget.Label, nextButton, prevButton *widget.Button) *widgets.ClickableBox {
	var txTable widgets.Table

	var allTxFilterNames = []string{"All", "Sent", "Received", "Transferred", "Coinbase", "Staking"}
	var allTxFilters = map[string]int32{
		"All":         dcrlibwallet.TxFilterAll,
		"Sent":        dcrlibwallet.TxFilterSent,
		"Received":    dcrlibwallet.TxFilterReceived,
		"Transferred": dcrlibwallet.TxFilterTransferred,
		"Coinbase":    dcrlibwallet.TxFilterCoinBase,
		"Staking":     dcrlibwallet.TxFilterStaking,
	}

	history.txCountForFilter, _ = wallet.CountTransactions(allTxFilters["All"])
	selectedAccountLabel := widget.NewLabel(fmt.Sprintf("%s (%d)", "All", history.txCountForFilter))
	history.selectedFilterCount = history.txCountForFilter

	var accountSelectionPopup *widget.PopUp
	accountListWidget := widget.NewVBox()
	for _, filterName := range allTxFilterNames {
		filterId := allTxFilters[filterName]
		history.txCountForFilter, history.txCountErr = wallet.CountTransactions(filterId)
		if history.txCountErr != nil {
			errorMessage := fmt.Sprintf("Cannot load history page. Error getting transaction count for filter %s: %s",
				filterName, history.txCountErr.Error())
			errorHandler(errorMessage, errorLabel)
			return nil
		}

		if history.txCountForFilter > 0 {
			filter := fmt.Sprintf("%s (%d)", filterName, history.txCountForFilter)
			accountsView := widget.NewHBox(
				widgets.NewHSpacer(5),
				widget.NewLabel(filter),
				widgets.NewHSpacer(5),
			)

			accountListWidget.Append(widgets.NewClickableBox(accountsView, func() {
				selectedFilterName := strings.Split(filter, " ")[0]
				history.selectedFilterId = allTxFilters[selectedFilterName]
				history.selectedFilterCount, _ = strconv.Atoi(strings.Split(filter, " ")[1])

				fetchAndDisplayTransactions(wallet, &txTable, nextButton, prevButton)
				history.txTable.Result.Children = txTable.Result.Children
				widget.Refresh(history.txTable.Result)
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

	// accountTab shows the selected account
	icons, _ := assets.GetIcons(assets.CollapseIcon)
	accountTab := widget.NewHBox(
		selectedAccountLabel,
		widgets.NewHSpacer(50),
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

func fetchAndDisplayTransactions(wallet *dcrlibwallet.LibWallet, txTable *widgets.Table, nextButton, prevButton *widget.Button) {
	tableHeading := widget.NewHBox(
		widget.NewLabelWithStyle("#", fyne.TextAlignCenter, fyne.TextStyle{Bold: true}),
		widget.NewLabelWithStyle("Date (UTC)", fyne.TextAlignCenter, fyne.TextStyle{Bold: true}),
		widget.NewLabelWithStyle("Direction", fyne.TextAlignCenter, fyne.TextStyle{Bold: true}),
		widget.NewLabelWithStyle("Status", fyne.TextAlignCenter, fyne.TextStyle{Bold: true}),		
		widget.NewLabelWithStyle("Amount", fyne.TextAlignCenter, fyne.TextStyle{Bold: true}),
		widget.NewLabelWithStyle("Fee", fyne.TextAlignCenter, fyne.TextStyle{Bold: true}),
		widget.NewLabelWithStyle("Type", fyne.TextAlignCenter, fyne.TextStyle{Bold: true}),
		widget.NewLabelWithStyle("Hash", fyne.TextAlignCenter, fyne.TextStyle{Bold: true}),
	)

	txOffset := 0
	if history.txns != nil {
		txOffset = len(history.txns)
	}

	txns, err := wallet.GetTransactionsRaw(int32(txOffset), txPerPage, history.selectedFilterId)
	if err != nil {
		// displayMessage(err.Error(), MessageKindError)
		// return
	}

	history.txl = len(txns)
	history.txns = append(history.txns, txns...)
	pageTxOffset := (history.currentPage - 1) * txPerPage
	maxTxIndexForCurrentPage := pageTxOffset + txPerPage

	var hBox []*widget.Box
	for currentTxIndex, tx := range history.txns {
		if currentTxIndex < int(pageTxOffset) {
			continue // skip txs not belonging to this page
		}
		if currentTxIndex >= int(maxTxIndexForCurrentPage) {
			break // max number of tx displayed for this page
		}

		status := "Pending"
		confirmations := wallet.GetBestBlock() - tx.BlockHeight + 1
		if tx.BlockHeight != -1 && confirmations > dcrlibwallet.DefaultRequiredConfirmations {
			status = "Confirmed"
		}

		hBox = append(hBox, widget.NewHBox(
			widget.NewLabelWithStyle(fmt.Sprintf("%d", currentTxIndex+1), fyne.TextAlignCenter, fyne.TextStyle{}),
			widget.NewLabelWithStyle(dcrlibwallet.ExtractDateOrTime(tx.Timestamp), fyne.TextAlignCenter, fyne.TextStyle{}),
			widget.NewLabelWithStyle(dcrlibwallet.TransactionDirectionName(tx.Direction), fyne.TextAlignCenter, fyne.TextStyle{}),
			widget.NewLabelWithStyle(status, fyne.TextAlignLeading, fyne.TextStyle{}),
			widget.NewLabelWithStyle(dcrutil.Amount(tx.Amount).String(), fyne.TextAlignTrailing, fyne.TextStyle{}),
			widget.NewLabelWithStyle(dcrutil.Amount(tx.Fee).String(),fyne.TextAlignCenter, fyne.TextStyle{}),
			widget.NewLabelWithStyle(tx.Type, fyne.TextAlignCenter, fyne.TextStyle{}),
			widgets.NewClickableBox(widget.NewHBox(widget.NewLabelWithStyle(tx.Hash, fyne.TextAlignLeading, fyne.TextStyle{Italic: true})), func() {
				history.txTable.Container.Hide()
				history.txDetailsTable.Container.Show()
				fetchTxDetail(&history.txDetailsTable, wallet, tx.Hash)
			}),
		))
	}

	txTable.NewTable(tableHeading, hBox...)
	txTable.Refresh()

	if history.currentPage > 1 {
		prevButton.Enable()
	}else{
		prevButton.Disable()
	}

	if history.selectedFilterCount > int(maxTxIndexForCurrentPage) {
		nextButton.Enable()
	}else{
		nextButton.Disable()
	}

	return
}

func loadPreviousPage(wallet *dcrlibwallet.LibWallet, nextButton, prevButton *widget.Button) {
	var txTable widgets.Table

	history.currentPage--
	history.txns = history.txns[:len(history.txns)-(int(txPerPage) + history.txl)]
	
	fetchAndDisplayTransactions(wallet, &txTable, nextButton, prevButton)
	history.txTable.Result.Children = txTable.Result.Children
	widget.Refresh(history.txTable.Result)
	return
}

func loadNextPage(wallet *dcrlibwallet.LibWallet, nextButton, prevButton *widget.Button) {
	var txTable widgets.Table

	nextPage := history.currentPage + 1
	history.currentPage = nextPage
	nextPageTxOffset := (nextPage - 1) * txPerPage
	
	if int(nextPageTxOffset) >= len(history.txns) {
		// we've not loaded txs for this page
		fetchAndDisplayTransactions(wallet, &txTable, nextButton, prevButton)
		history.txTable.Result.Children = txTable.Result.Children
		widget.Refresh(history.txTable.Result)
	}

	return
}

