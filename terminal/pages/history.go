package pages

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/decred/dcrd/dcrutil"
	"github.com/gdamore/tcell"
	"github.com/raedahgroup/dcrlibwallet"
	"github.com/raedahgroup/dcrlibwallet/utils"
	"github.com/raedahgroup/godcr/terminal/helpers"
	"github.com/raedahgroup/godcr/terminal/primitives"
	"github.com/rivo/tview"
)

type messageKind uint32

const (
	MessageKindError messageKind = iota
	MessageKindInfo
)

const txPerPage int32 = 25

var historyPageData struct {
	pageContentHolder *tview.Flex
	titleTextView     *primitives.TextView
	messageTextView   *primitives.TextView

	txFilterDropDown        *primitives.Form
	historyTable            *tview.Table
	transactionDetailsTable *tview.Table

	currentTxFilter              int32
	totalTxCountForCurrentFilter int
	displayedTxs                 []*dcrlibwallet.Transaction
}

func historyPage() tview.Primitive {
	// parent flexbox layout container to hold other primitives
	historyPageData.pageContentHolder = tview.NewFlex().SetDirection(tview.FlexRow)
	// handler for returning back to menu column
	historyPageData.pageContentHolder.SetInputCapture(func(event *tcell.EventKey) *tcell.EventKey {
		if event.Key() == tcell.KeyEscape || event.Key() == tcell.KeyBackspace || event.Key() == tcell.KeyBackspace2 {
			commonPageData.clearAllPageContent()
			return nil
		}
		return event
	})

	historyPageData.titleTextView = primitives.NewLeftAlignedTextView("History")
	historyPageData.pageContentHolder.AddItem(historyPageData.titleTextView, 2, 0, false)

	historyPageData.messageTextView = primitives.WordWrappedTextView("")

	historyPageData.txFilterDropDown = prepareTxFilterDropDown()
	if historyPageData.txFilterDropDown == nil {
		commonPageData.app.SetFocus(historyPageData.pageContentHolder)
		return historyPageData.pageContentHolder
	}

	historyPageData.historyTable = prepareHistoryTable()
	historyPageData.transactionDetailsTable = prepareTxDetailsTable()

	// fetch tx to display in subroutine so the UI isn't blocked
	go fetchAndDisplayTransactions(0)

	commonPageData.hintTextView.SetText("TIP: Use ARROW UP/DOWN to select txn, \n" +
		"ENTER to view details, ESC to return to navigation menu")
	commonPageData.app.SetFocus(historyPageData.pageContentHolder)

	return historyPageData.pageContentHolder
}

func prepareTxFilterDropDown() *primitives.Form {
	var allTxFilters = map[string]int32{
		"All":         dcrlibwallet.TxFilterAll,
		"Sent":        dcrlibwallet.TxFilterSent,
		"Received":    dcrlibwallet.TxFilterReceived,
		"Transferred": dcrlibwallet.TxFilterTransferred,
		"Coinbase":    dcrlibwallet.TxFilterCoinBase,
		"Staking":     dcrlibwallet.TxFilterStaking,
	}

	var txFilterSelectionOptions []string
	var activeFiltersWithTxCounts = map[int32]int{}

	for filterName, filterId := range allTxFilters {
		txCountForFilter, txCountErr := commonPageData.wallet.CountTransactions(filterId)
		if txCountErr != nil {
			errorMessage := fmt.Sprintf("Cannot load history page. Error getting transaction count for filter %s: %s",
				filterName, txCountErr.Error())
			displayMessage(errorMessage, MessageKindError)
			return nil
		}

		if txCountForFilter > 0 {
			activeFiltersWithTxCounts[filterId] = txCountForFilter
			txFilterSelectionOptions = append(txFilterSelectionOptions, fmt.Sprintf("%s (%d)", filterName, txCountForFilter))
		}
	}

	if len(activeFiltersWithTxCounts) == 0 {
		displayMessage("No transactions yet", MessageKindInfo)
		commonPageData.hintTextView.SetText("TIP: ESC or BACKSPACE to return to navigation menu")
		commonPageData.app.SetFocus(historyPageData.pageContentHolder)
		return nil
	}

	txFilterDropDown := primitives.NewForm(false)
	txFilterDropDown.SetBorderPadding(0, 0, 0, 0)

	// dropDown selection change listener
	txFilterDropDown.AddDropDown("", txFilterSelectionOptions, 0, func(selectedOption string, index int) {
		selectedFilterName := strings.Split(selectedOption, " ")[0]
		selectedFilterId := allTxFilters[selectedFilterName]
		if selectedFilterId != historyPageData.currentTxFilter {
			historyPageData.currentTxFilter = selectedFilterId
			historyPageData.totalTxCountForCurrentFilter = activeFiltersWithTxCounts[selectedFilterId]
			go fetchAndDisplayTransactions(0)
		}
	})

	// handler for switching between dropDown and table
	txFilterDropDown.SetInputCapture(func(event *tcell.EventKey) *tcell.EventKey {
		if event.Key() == tcell.KeyEscape || event.Key() == tcell.KeyBackspace || event.Key() == tcell.KeyBackspace2 {
			commonPageData.clearAllPageContent()
			return nil
		}
		if event.Key() == tcell.KeyTab {
			commonPageData.app.SetFocus(historyPageData.historyTable)
			return nil
		}
		return event
	})

	historyPageData.pageContentHolder.AddItem(txFilterDropDown, 2, 0, false)

	historyPageData.displayedTxs = nil
	historyPageData.currentTxFilter = dcrlibwallet.TxFilterAll
	historyPageData.totalTxCountForCurrentFilter = activeFiltersWithTxCounts[dcrlibwallet.TxFilterAll]

	return txFilterDropDown
}

func prepareHistoryTable() *tview.Table {
	historyTable := tview.NewTable().
		SetBorders(false).
		SetFixed(1, 0). // keep first row (column headers) fixed during scroll
		SetSelectable(true, false)

	historyTable.SetDoneFunc(func(key tcell.Key) {
		if key == tcell.KeyEscape {
			commonPageData.clearAllPageContent()
		}
	})

	// handler for switching between dropDown and table
	historyTable.SetInputCapture(func(event *tcell.EventKey) *tcell.EventKey {
		if event.Key() == tcell.KeyEscape || event.Key() == tcell.KeyBackspace || event.Key() == tcell.KeyBackspace2 {
			commonPageData.clearAllPageContent()
			return nil
		}
		if event.Key() == tcell.KeyTab {
			commonPageData.app.SetFocus(historyPageData.txFilterDropDown)
			return nil
		}
		return event
	})

	// method for getting transaction details when a tx is selected from the history table
	historyTable.SetSelectedFunc(func(row, column int) {
		if row >= len(historyPageData.displayedTxs) {
			// ignore selected func call for table header
			return
		}

		historyPageData.pageContentHolder.RemoveItem(historyTable)
		historyPageData.pageContentHolder.RemoveItem(historyPageData.txFilterDropDown)

		historyPageData.titleTextView.SetText("Transaction Details")
		commonPageData.hintTextView.SetText("TIP: Use ARROW UP/DOWN to scroll, \nBACKSPACE to view History page, ESC to return to navigation menu")

		historyPageData.transactionDetailsTable.Clear()
		historyPageData.pageContentHolder.AddItem(historyPageData.transactionDetailsTable, 0, 1, true)
		commonPageData.app.SetFocus(historyPageData.transactionDetailsTable)

		selectedTx := historyPageData.displayedTxs[row-1]
		displayTxDetails(selectedTx, historyPageData.transactionDetailsTable)
	})

	return historyTable
}

func prepareTxDetailsTable() *tview.Table {
	transactionDetailsTable := tview.NewTable().SetBorders(false)

	// handler for returning back to history table from tx details table
	transactionDetailsTable.SetInputCapture(func(event *tcell.EventKey) *tcell.EventKey {
		if event.Key() == tcell.KeyEscape || event.Key() == tcell.KeyBackspace || event.Key() == tcell.KeyBackspace2 {
			historyPageData.pageContentHolder.AddItem(historyPageData.txFilterDropDown, 2, 0, false)

			historyPageData.pageContentHolder.RemoveItem(historyPageData.transactionDetailsTable)

			historyPageData.titleTextView.SetText("History")
			commonPageData.hintTextView.SetText("TIP: Use ARROW UP/DOWN to select txn,\nENTER to view details, " +
				"ESC to return to navigation menu")

			historyPageData.pageContentHolder.AddItem(historyPageData.historyTable, 0, 1, true)
			commonPageData.app.SetFocus(historyPageData.historyTable)

			return nil
		}
		return event
	})

	return transactionDetailsTable
}

func fetchAndDisplayTransactions(txOffset int) {
	// show a loading text at the bottom of the table so user knows an op is in progress
	displayMessage("Fetching data...", MessageKindInfo)

	txns, err := commonPageData.wallet.GetTransactionsRaw(int32(txOffset), txPerPage, historyPageData.currentTxFilter)
	if err != nil {
		displayMessage(err.Error(), MessageKindError)
		return
	}

	// calculate max number of digits after decimal point for all tx amounts
	inputsAndOutputsAmount := make([]int64, len(txns))
	for i, tx := range txns {
		inputsAndOutputsAmount[i] = tx.Amount
	}
	maxDecimalPlacesForTxAmounts := utils.MaxDecimalPlaces(inputsAndOutputsAmount)

	// updating the history table from a goroutine, use tviewApp.QueueUpdateDraw
	commonPageData.app.QueueUpdateDraw(func() {
		for _, tx := range txns {
			nextRowIndex := historyPageData.historyTable.GetRowCount()

			dateCell := tview.NewTableCell(fmt.Sprintf("%-10s", dcrlibwallet.ExtractDateOrTime(tx.Timestamp))).
				SetAlign(tview.AlignCenter).
				SetMaxWidth(1).
				SetExpansion(1)
			historyPageData.historyTable.SetCell(nextRowIndex, 0, dateCell)

			directionCell := tview.NewTableCell(fmt.Sprintf("%-10s", dcrlibwallet.TransactionDirectionName(tx.Direction))).
				SetAlign(tview.AlignCenter).
				SetMaxWidth(2).
				SetExpansion(1)
			historyPageData.historyTable.SetCell(nextRowIndex, 1, directionCell)

			formattedAmount := utils.FormatAmountDisplay(tx.Amount, maxDecimalPlacesForTxAmounts)
			amountCell := tview.NewTableCell(fmt.Sprintf("%15s", formattedAmount)).
				SetAlign(tview.AlignCenter).
				SetMaxWidth(3).
				SetExpansion(1)
			historyPageData.historyTable.SetCell(nextRowIndex, 2, amountCell)

			status := "Pending"
			confirmations := commonPageData.wallet.GetBestBlock() - tx.BlockHeight + 1
			if tx.BlockHeight != -1 && confirmations > dcrlibwallet.DefaultRequiredConfirmations {
				status = "Confirmed"
			}
			statusCell := tview.NewTableCell(fmt.Sprintf("%12s", status)).
				SetAlign(tview.AlignCenter).
				SetMaxWidth(1).
				SetExpansion(1)
			historyPageData.historyTable.SetCell(nextRowIndex, 3, statusCell)

			typeCell := tview.NewTableCell(fmt.Sprintf("%-8s", tx.Type)).
				SetAlign(tview.AlignCenter).
				SetMaxWidth(1).
				SetExpansion(1)
			historyPageData.historyTable.SetCell(nextRowIndex, 4, typeCell)

			historyPageData.displayedTxs = append(historyPageData.displayedTxs, tx)
		}

		// clear loading message text
		displayMessage("", MessageKindInfo)
	})

	if len(historyPageData.displayedTxs) < historyPageData.totalTxCountForCurrentFilter {
		// set or reset selection changed listener to load more data when the table is almost scrolled to the end
		nextOffset := txOffset + len(txns)
		historyPageData.historyTable.SetSelectionChangedFunc(func(row, column int) {
			if row >= historyPageData.historyTable.GetRowCount()-10 {
				historyPageData.historyTable.SetSelectionChangedFunc(nil) // unset selection change listener until table is populated
				fetchAndDisplayTransactions(nextOffset)
			}
		})
	}

	return
}

func displayTxDetails(tx *dcrlibwallet.Transaction, transactionDetailsTable *tview.Table) {
	transactionDetailsTable.SetCellSimple(0, 0, "Hash")
	transactionDetailsTable.SetCellSimple(1, 0, "Confirmations")
	transactionDetailsTable.SetCellSimple(2, 0, "Included in block")
	transactionDetailsTable.SetCellSimple(3, 0, "Type")
	transactionDetailsTable.SetCellSimple(4, 0, "Amount")
	transactionDetailsTable.SetCellSimple(5, 0, "Date")
	transactionDetailsTable.SetCellSimple(6, 0, "Direction")
	transactionDetailsTable.SetCellSimple(7, 0, "Fee")
	transactionDetailsTable.SetCellSimple(8, 0, "Fee Rate")

	var confirmations int32 = 0
	if tx.BlockHeight != -1 {
		confirmations = commonPageData.wallet.GetBestBlock() - tx.BlockHeight + 1
	}

	transactionDetailsTable.SetCellSimple(0, 1, tx.Hash)
	transactionDetailsTable.SetCellSimple(1, 1, strconv.Itoa(int(confirmations)))
	transactionDetailsTable.SetCellSimple(2, 1, strconv.Itoa(int(tx.BlockHeight)))
	transactionDetailsTable.SetCellSimple(3, 1, tx.Type)
	transactionDetailsTable.SetCellSimple(4, 1, dcrutil.Amount(tx.Amount).String())
	transactionDetailsTable.SetCellSimple(5, 1, fmt.Sprintf("%s UTC", dcrlibwallet.FormatUTCTime(tx.Timestamp)))
	transactionDetailsTable.SetCellSimple(6, 1, dcrlibwallet.TransactionDirectionName(tx.Direction))
	transactionDetailsTable.SetCellSimple(7, 1, dcrutil.Amount(tx.Fee).String())
	transactionDetailsTable.SetCellSimple(8, 1, fmt.Sprintf("%s/kB", dcrutil.Amount(tx.FeeRate)))

	// calculate max number of digits after decimal point for inputs and outputs
	inputsAndOutputsAmount := make([]int64, 0, len(tx.Inputs)+len(tx.Outputs))
	for _, txIn := range tx.Inputs {
		inputsAndOutputsAmount = append(inputsAndOutputsAmount, txIn.Amount)
	}
	for _, txOut := range tx.Outputs {
		inputsAndOutputsAmount = append(inputsAndOutputsAmount, txOut.Amount)
	}
	maxDecimalPlacesForInputsAndOutputsAmounts := utils.MaxDecimalPlaces(inputsAndOutputsAmount)

	// now format amount having determined the max number of decimal places
	formatAmount := func(amount int64) string {
		return utils.FormatAmountDisplay(amount, maxDecimalPlacesForInputsAndOutputsAmounts)
	}

	transactionDetailsTable.SetCellSimple(9, 0, "-Inputs-")
	for _, txIn := range tx.Inputs {
		row := transactionDetailsTable.GetRowCount()
		transactionDetailsTable.SetCell(row, 0, tview.NewTableCell(formatAmount(txIn.Amount)).SetAlign(tview.AlignRight))
		transactionDetailsTable.SetCellSimple(row, 1, txIn.PreviousOutpoint)
	}

	row := transactionDetailsTable.GetRowCount()
	transactionDetailsTable.SetCellSimple(row, 0, "-Outputs-")
	for _, txOut := range tx.Outputs {
		row++
		outputAmount := formatAmount(txOut.Amount)

		if txOut.Address == "" {
			transactionDetailsTable.SetCellSimple(row, 0, fmt.Sprintf("  %s (no address)", outputAmount))
			continue
		}

		transactionDetailsTable.SetCell(row, 0, tview.NewTableCell(outputAmount).SetAlign(tview.AlignRight))
		transactionDetailsTable.SetCellSimple(row, 1, fmt.Sprintf("%s (%s)", txOut.Address, txOut.AccountName))
	}
}

func displayMessage(message string, kind messageKind) {
	// this function may be called from a goroutine, use tviewApp.QueueUpdateDraw
	commonPageData.app.QueueUpdateDraw(func() {
		historyPageData.pageContentHolder.RemoveItem(historyPageData.messageTextView)
		if message != "" {
			if kind == MessageKindError {
				historyPageData.messageTextView.SetTextColor(helpers.DecredOrangeColor)
			} else {
				historyPageData.messageTextView.SetTextColor(tcell.ColorWhite)
			}

			historyPageData.messageTextView.SetText(message)
			historyPageData.pageContentHolder.AddItem(historyPageData.messageTextView, 2, 0, false)
		}
	})
}
