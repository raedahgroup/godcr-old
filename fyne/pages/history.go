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

const txPerPage int32 = 10

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
	errorLabel     *widget.Label
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

	pageTitleLabel := widget.NewLabelWithStyle("Transactions", fyne.TextAlignLeading, fyne.TextStyle{Bold: true, Italic: true})
	
	filterDropdown := txFilterDropDown(wallet, window, errorLabel)

	fetchAndDisplayTransactions(wallet, &history.txTable, window)
	
	output := widget.NewVBox(
		widgets.NewVSpacer(5),
		widget.NewHBox(pageTitleLabel),
		widgets.NewVSpacer(5),
		filterDropdown,
		widgets.NewVSpacer(5),
		fyne.NewContainerWithLayout(layout.NewFixedGridLayout(fyne.NewSize(history.txTable.Container.MinSize().Width, history.txTable.Container.MinSize().Height+200)), history.txTable.Container),
		widgets.NewVSpacer(15),
		errorLabel,
	)

	return widget.NewHBox(widgets.NewHSpacer(18), output)
}

func txFilterDropDown(wallet *dcrlibwallet.LibWallet, window fyne.Window, errorLabel *widget.Label) *widgets.ClickableBox {
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

				fetchAndDisplayTransactions(wallet, &txTable, window)
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

func fetchAndDisplayTransactions(wallet *dcrlibwallet.LibWallet, txTable *widgets.Table, window fyne.Window) {
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

	var hBox []*widget.Box
	for currentTxIndex, tx := range txns {
		status := "Pending"
		confirmations := wallet.GetBestBlock() - tx.BlockHeight + 1
		if tx.BlockHeight != -1 && confirmations > dcrlibwallet.DefaultRequiredConfirmations {
			status = "Confirmed"
		}
		trimmedHash := tx.Hash[:15] + "..." + tx.Hash[len(tx.Hash)-15:]
		hBox = append(hBox, widget.NewHBox(
			widget.NewLabelWithStyle(fmt.Sprintf("%d", currentTxIndex+1), fyne.TextAlignCenter, fyne.TextStyle{}),
			widget.NewLabelWithStyle(dcrlibwallet.ExtractDateOrTime(tx.Timestamp), fyne.TextAlignCenter, fyne.TextStyle{}),
			widget.NewLabelWithStyle(dcrlibwallet.TransactionDirectionName(tx.Direction), fyne.TextAlignCenter, fyne.TextStyle{}),
			widget.NewLabelWithStyle(status, fyne.TextAlignLeading, fyne.TextStyle{}),
			widget.NewLabelWithStyle(dcrutil.Amount(tx.Amount).String(), fyne.TextAlignTrailing, fyne.TextStyle{}),
			widget.NewLabelWithStyle(dcrutil.Amount(tx.Fee).String(),fyne.TextAlignCenter, fyne.TextStyle{}),
			widget.NewLabelWithStyle(tx.Type, fyne.TextAlignCenter, fyne.TextStyle{}),
			widgets.NewClickableBox(widget.NewHBox(widget.NewLabelWithStyle(trimmedHash, fyne.TextAlignLeading, fyne.TextStyle{Italic: true})), func() {
				fetchTxDetails(tx.Hash, wallet, window)
			}),
		))
	}

	txTable.NewTable(tableHeading, hBox...)

	return
}

func fetchTxDetails(hash string, wallet *dcrlibwallet.LibWallet, window fyne.Window) {
	var confirmations int32 = 0
	// if txDetails.BlockHeight != -1 {
	// 	confirmations = wallet.GetBestBlock() - txDetails.BlockHeight + 1
	// }
	newHash, _ := chainhash.NewHashFromStr(hash)
	txDetails, err := wallet.GetTransactionRaw(newHash[:])
	if err != nil {

	}

	var spendUnconfirmed = wallet.ReadBoolConfigValueForKey(dcrlibwallet.SpendUnconfirmedConfigKey)

	var status string
	// var statusColor color.RGBA
	if spendUnconfirmed || confirmations > dcrlibwallet.DefaultRequiredConfirmations {
		status = "Confirmed"
		// statusColor = styles.DecredGreenColor
	} else {
		status = "Pending"
		// statusColor = styles.DecredOrangeColor
	}

	tableConfirmations := widget.NewHBox(
		widget.NewLabelWithStyle("Confirmations:", fyne.TextAlignTrailing, fyne.TextStyle{Bold: true}),
		widget.NewLabelWithStyle(strconv.Itoa(int(confirmations)), fyne.TextAlignCenter, fyne.TextStyle{}),
	)
	tableHash := widget.NewHBox(
		widget.NewLabelWithStyle("Hash:", fyne.TextAlignTrailing, fyne.TextStyle{Bold: true}),
		widget.NewLabelWithStyle(txDetails.Hash, fyne.TextAlignCenter, fyne.TextStyle{}),
	)
	tableBlockHeight := widget.NewHBox(
		widget.NewLabelWithStyle("Block Height:", fyne.TextAlignTrailing, fyne.TextStyle{Bold: true}),
		widget.NewLabelWithStyle(strconv.Itoa(int(txDetails.BlockHeight)), fyne.TextAlignCenter, fyne.TextStyle{}),
	)
	tableDirection := widget.NewHBox(
		widget.NewLabelWithStyle("Direction:", fyne.TextAlignTrailing, fyne.TextStyle{Bold: true}),
		widget.NewLabelWithStyle(dcrlibwallet.TransactionDirectionName(txDetails.Direction), fyne.TextAlignCenter, fyne.TextStyle{}),
	)
	tableType := widget.NewHBox(
		widget.NewLabelWithStyle("Type:", fyne.TextAlignTrailing, fyne.TextStyle{Bold: true}),
		widget.NewLabelWithStyle(txDetails.Type, fyne.TextAlignCenter, fyne.TextStyle{}),
	)
	tableAmount := widget.NewHBox(
		widget.NewLabelWithStyle("Amount:", fyne.TextAlignTrailing, fyne.TextStyle{Bold: true}),
		widget.NewLabelWithStyle(dcrutil.Amount(txDetails.Amount).String(), fyne.TextAlignCenter, fyne.TextStyle{}),
	)
	tableSize := widget.NewHBox(
		widget.NewLabelWithStyle("Size:", fyne.TextAlignTrailing, fyne.TextStyle{Bold: true}),
		widget.NewLabelWithStyle(strconv.Itoa(txDetails.Size)+" Bytes", fyne.TextAlignCenter, fyne.TextStyle{}),
	)
	tableFee := widget.NewHBox(
		widget.NewLabelWithStyle("Fee:", fyne.TextAlignTrailing, fyne.TextStyle{Bold: true}),
		widget.NewLabelWithStyle(dcrutil.Amount(txDetails.Fee).String(), fyne.TextAlignCenter, fyne.TextStyle{}),
	)
	tableFeeRate := widget.NewHBox(
		widget.NewLabelWithStyle("Fee Rate:", fyne.TextAlignTrailing, fyne.TextStyle{Bold: true}),
		widget.NewLabelWithStyle(dcrutil.Amount(txDetails.FeeRate).String(), fyne.TextAlignCenter, fyne.TextStyle{}),
	)
	tableStatus := widget.NewHBox(
		widget.NewLabelWithStyle("Status:", fyne.TextAlignTrailing, fyne.TextStyle{Bold: true}),
		widget.NewLabelWithStyle(status, fyne.TextAlignCenter, fyne.TextStyle{}),
	)
	tableDate := widget.NewHBox(
		widget.NewLabelWithStyle("Date:", fyne.TextAlignTrailing, fyne.TextStyle{Bold: true}),
		widget.NewLabelWithStyle(fmt.Sprintf("%s UTC", dcrlibwallet.FormatUTCTime(txDetails.Timestamp)), fyne.TextAlignCenter, fyne.TextStyle{}),
	)

	tableData := widget.NewVBox(
		tableConfirmations,
		tableHash,
		tableBlockHeight,
		tableDirection,
		tableType,
		tableAmount,
		tableSize,
		tableFee,
		tableFeeRate,
		tableStatus,
		tableDate,
	)


	var txInput widgets.Table
	heading := widget.NewHBox(
		widget.NewLabelWithStyle("Previous Outpoint", fyne.TextAlignCenter, fyne.TextStyle{Bold: true}),
		widget.NewLabelWithStyle("Account", fyne.TextAlignCenter, fyne.TextStyle{Bold: true}),
		widget.NewLabelWithStyle("Amount", fyne.TextAlignCenter, fyne.TextStyle{Bold: true}))

	var inputBox []*widget.Box
	for i := range txDetails.Inputs {
		inputBox = append(inputBox, widget.NewHBox(
			widget.NewLabelWithStyle(txDetails.Inputs[i].PreviousOutpoint, fyne.TextAlignLeading, fyne.TextStyle{}),
			widget.NewLabelWithStyle(txDetails.Inputs[i].AccountName, fyne.TextAlignCenter, fyne.TextStyle{}),
			widget.NewLabelWithStyle(dcrutil.Amount(txDetails.Inputs[i].Amount).String(), fyne.TextAlignTrailing, fyne.TextStyle{}),
		))
	}
	txInput.NewTable(heading, inputBox...)

	var txOutput widgets.Table
	heading = widget.NewHBox(
		widget.NewLabelWithStyle("Address", fyne.TextAlignCenter, fyne.TextStyle{Bold: true}),
		widget.NewLabelWithStyle("Account", fyne.TextAlignCenter, fyne.TextStyle{Bold: true}),
		widget.NewLabelWithStyle("Value", fyne.TextAlignCenter, fyne.TextStyle{Bold: true}),
		widget.NewLabelWithStyle("Type", fyne.TextAlignCenter, fyne.TextStyle{Bold: true}))

	var outputBox []*widget.Box
	for i := range txDetails.Outputs {
		outputBox = append(outputBox, widget.NewHBox(
			widget.NewLabelWithStyle(txDetails.Outputs[i].Address, fyne.TextAlignLeading, fyne.TextStyle{}),
			widget.NewLabelWithStyle(txDetails.Outputs[i].AccountName, fyne.TextAlignCenter, fyne.TextStyle{}),
			widget.NewLabelWithStyle(dcrutil.Amount(txDetails.Outputs[i].Amount).String(), fyne.TextAlignTrailing, fyne.TextStyle{}),
			widget.NewLabelWithStyle(txDetails.Outputs[i].ScriptType, fyne.TextAlignCenter, fyne.TextStyle{})))
	}
	txOutput.NewTable(heading, outputBox...)

	label := widget.NewLabelWithStyle("Transaction Details", fyne.TextAlignLeading, fyne.TextStyle{Bold: true, Italic: true})

	output := widget.NewHBox(widgets.NewHSpacer(10), widget.NewVBox(
		label,
		widgets.NewVSpacer(10),
		tableData,
		widget.NewLabelWithStyle("Inputs", fyne.TextAlignLeading, fyne.TextStyle{Bold: true}),
		txInput.Result,
		widgets.NewVSpacer(10),
		widget.NewLabelWithStyle("Outputs", fyne.TextAlignLeading, fyne.TextStyle{Bold: true}),
		txOutput.Result,
	), widgets.NewHSpacer(10))

	scrollContainer := widget.NewScrollContainer(output)
	scrollContainer.Resize(fyne.NewSize(scrollContainer.MinSize().Width, 500))
	popUp := widget.NewPopUp(widget.NewVBox(fyne.NewContainer(scrollContainer)),
		window.Canvas())
	popUp.Show()
}

