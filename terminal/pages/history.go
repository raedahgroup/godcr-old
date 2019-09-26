package pages

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/decred/dcrd/dcrutil"
	"github.com/gdamore/tcell"
	"github.com/raedahgroup/dcrlibwallet/txindex"
	"github.com/raedahgroup/dcrlibwallet/utils"
	godcrUtils "github.com/raedahgroup/godcr/app/utils"
	"github.com/raedahgroup/godcr/app/walletcore"
	"github.com/raedahgroup/godcr/terminal/helpers"
	"github.com/raedahgroup/godcr/terminal/primitives"
	"github.com/rivo/tview"
)

var displayedTxHashes []string
var txPerPage int32 = walletcore.TransactionHistoryCountPerPage
var totalTxCount int
var selectedFilter, lastSelectedFilter string
type messageKind uint32
const (
	MessageKindError messageKind = iota
	InfoMessage messageKind = iota
)

func historyPage(wallet walletcore.Wallet, hintTextView *primitives.TextView, tviewApp *tview.Application, clearFocus func()) tview.Primitive {
	// parent flexbox layout container to hold other primitives
	historyPage := tview.NewFlex().SetDirection(tview.FlexRow)

	// handler for returning back to menu column
	historyPage.SetInputCapture(func(event *tcell.EventKey) *tcell.EventKey {
		if event.Key() == tcell.KeyEscape || event.Key() == tcell.KeyBackspace || event.Key() == tcell.KeyBackspace2 {
			clearFocus()
			return nil
		}

		return event
	})

	messageTextView := primitives.WordWrappedTextView("")
	displayMessage := func(message string, kind messageKind) {
		// this function may be called from a goroutine, use tviewApp.QueueUpdateDraw
		tviewApp.QueueUpdateDraw(func() {
			historyPage.RemoveItem(messageTextView)
			if message != "" {
				if kind == MessageKindError {
					messageTextView.SetTextColor(helpers.DecredOrangeColor)
				} else {
					messageTextView.SetTextColor(tcell.ColorWhite)
				}

				messageTextView.SetText(message)
				historyPage.AddItem(messageTextView, 2, 0, false)
			}
		})
	}

	// page title
	titleTextView := primitives.NewLeftAlignedTextView("History")
	historyPage.AddItem(titleTextView, 2, 0, false)

	// get total tx count
	txCount, txCountErr := wallet.TransactionCount(nil)
	if txCountErr != nil {
		displayMessage(fmt.Sprintf("Cannot load history page. Error getting total transaction count: %s", txCountErr.Error()), MessageKindError)
		tviewApp.SetFocus(historyPage)
		return historyPage
	}
	totalTxCount = txCount

	// tx filter
	var transactionFilters []string
	var selectedTxFilter *txindex.ReadFilter
	for _, filter := range walletcore.TransactionFilters {
		if filter == "All" {
			selectedTxFilter = nil
		} else {
			selectedTxFilter = walletcore.BuildTransactionFilter(filter)
		}

		filteredTxCount, txCountErr := wallet.TransactionCount(selectedTxFilter)
		if txCountErr != nil {
			displayMessage(fmt.Sprintf("Cannot load history page. Error getting transaction count for filter %s: %s", selectedTxFilter, txCountErr.Error()), MessageKindError)
			tviewApp.SetFocus(historyPage)
			return historyPage
		}
		if filteredTxCount == 0 {
			continue
		}

		transactionFilters = append(transactionFilters, fmt.Sprintf("%s (%d)", filter, filteredTxCount))
	}

	if totalTxCount == 0 {
		displayMessage("No transactions yet", InfoMessage)
		hintTextView.SetText("TIP: ESC or BACKSPACE to return to navigation menu")
		tviewApp.SetFocus(historyPage)
		return historyPage
	}

	txFilterDropdown := primitives.NewForm(false)
	txFilterDropdown.SetBorderPadding(0, 0, 0, 0)
	historyPage.AddItem(txFilterDropdown, 2, 0, false)

	historyTable := tview.NewTable().
	SetBorders(false).
	SetFixed(1, 0).
	SetSelectable(true, false)

	transactionDetailsTable := tview.NewTable().SetBorders(false)

	// filter dropDown sleection options
	txFilterDropdown.AddDropDown("", transactionFilters, 0, func(filterAndCount string, index int) {
		lastSelectedFilter = selectedFilter
		selectedFilterAndCount := strings.Split(filterAndCount, " ")
		selectedFilter = selectedFilterAndCount[0]

		filter := walletcore.BuildTransactionFilter(selectedFilter)
		if selectedFilter != lastSelectedFilter {
			if selectedFilter == "All"{
				go fetchAndDisplayTransactions(0, wallet, nil, historyTable, tviewApp, displayMessage)
			}else{
				go fetchAndDisplayTransactions(0, wallet, filter, historyTable, tviewApp, displayMessage)
			}
		}
	})

	displayHistoryTable := func() {
		historyPage.RemoveItem(transactionDetailsTable)

		titleTextView.SetText("History")
		hintTextView.SetText("TIP: Use ARROW UP/DOWN to select txn,\nENTER to view details, ESC to return to navigation menu")

		historyPage.AddItem(historyTable, 0, 1, true)

		tviewApp.SetFocus(historyTable)

	}

	historyTable.SetDoneFunc(func(key tcell.Key) {
		if key == tcell.KeyEscape {
			clearFocus()
		}
	})

	// method for getting transaction details when a tx is selected from the history table
	historyTable.SetSelectedFunc(func(row, column int) {
		if row >= len(displayedTxHashes) {
			// ignore selected func call for table header
			return
		}

		historyPage.RemoveItem(historyTable)
		historyPage.RemoveItem(txFilterDropdown)

		txHash := displayedTxHashes[row-1]

		titleTextView.SetText("Transaction Details")
		hintTextView.SetText("TIP: Use ARROW UP/DOWN to scroll, \nBACKSPACE to view History page, ESC to return to navigation menu")

		transactionDetailsTable.Clear()
		historyPage.AddItem(transactionDetailsTable, 0, 1, true)

		tviewApp.SetFocus(transactionDetailsTable)

		displayTxDetails(txHash, wallet, displayMessage, transactionDetailsTable)
	})

	// handler for returning back to history table
	transactionDetailsTable.SetInputCapture(func(event *tcell.EventKey) *tcell.EventKey {
		if event.Key() == tcell.KeyEscape || event.Key() == tcell.KeyBackspace || event.Key() == tcell.KeyBackspace2 {
			historyPage.AddItem(txFilterDropdown, 2, 0, false)
			displayHistoryTable()
			return nil
		}

		return event
	})

	// handler for switching between dropDown and table
	txFilterDropdown.SetInputCapture(func(event *tcell.EventKey) *tcell.EventKey {
		if event.Key() == tcell.KeyEscape || event.Key() == tcell.KeyBackspace || event.Key() == tcell.KeyBackspace2 {
			clearFocus()
			return nil
		}
		if event.Key() == tcell.KeyTab {
			tviewApp.SetFocus(historyTable)
			return nil
		}

		return event
	})

	// handler for switching between dropDown and table
	historyTable.SetInputCapture(func(event *tcell.EventKey) *tcell.EventKey {
		if event.Key() == tcell.KeyEscape || event.Key() == tcell.KeyBackspace || event.Key() == tcell.KeyBackspace2 {
			clearFocus()
			return nil
		}
		if event.Key() == tcell.KeyTab {
			tviewApp.SetFocus(txFilterDropdown)
			return nil
		}

		return event
	})

	displayHistoryTable()

	// fetch tx to display in subroutine so the UI isn't blocked
	selectedFilter, lastSelectedFilter = "All", "All"
	go fetchAndDisplayTransactions(0, wallet, nil, historyTable, tviewApp, displayMessage)

	hintTextView.SetText("TIP: Use TAB to switch/navigate, ARROW UP/DOWN to select txn, \nENTER to view details, ESC to return to navigation menu")

	tviewApp.SetFocus(historyPage)

	return historyPage
}

func fetchAndDisplayTransactions(txOffset int, 
	wallet walletcore.Wallet, 
	filter *txindex.ReadFilter, 
	historyTable *tview.Table, 
	tviewApp *tview.Application, 
	displayMessage func(string, messageKind)) {
		// show a loading text at the bottom of the table so user knows an op is in progress
		displayMessage("Fetching data...", InfoMessage)

		txns, err := wallet.TransactionHistory(int32(txOffset), txPerPage, filter)
		if err != nil {
			displayMessage(err.Error(), MessageKindError)
			return
		}

		// calculate max number of digits after decimal point for all tx amounts
		inputsAndOutputsAmount := make([]int64, len(txns))
		for i, tx := range txns {
			inputsAndOutputsAmount[i] = tx.Amount
		}
		maxDecimalPlacesForTxAmounts := godcrUtils.MaxDecimalPlaces(inputsAndOutputsAmount)

		// now format amount having determined the max number of decimal places
		formatAmount := func(amount int64) string {
			return godcrUtils.FormatAmountDisplay(amount, maxDecimalPlacesForTxAmounts)
		}

		tableHeaderCell := func(text string) *tview.TableCell {
			return tview.NewTableCell(text).SetAlign(tview.AlignCenter).SetSelectable(false).SetMaxWidth(1).SetExpansion(1)
		}

		if selectedFilter != lastSelectedFilter {
			historyTable.Clear()
			displayedTxHashes = nil
		}

		// history table header
		historyTable.SetCell(0, 0, tableHeaderCell("Date (UTC)"))
		historyTable.SetCell(0, 1, tableHeaderCell(fmt.Sprintf("%10s", "Direction")))
		historyTable.SetCell(0, 2, tableHeaderCell(fmt.Sprintf("%8s", "Amount")))
		historyTable.SetCell(0, 3, tableHeaderCell(fmt.Sprintf("%5s", "Status")))
		historyTable.SetCell(0, 4, tableHeaderCell(fmt.Sprintf("%-5s", "Type")))

		// updating the history table from a goroutine, use tviewApp.QueueUpdateDraw
		tviewApp.QueueUpdateDraw(func() {
			for _, tx := range txns {
				nextRowIndex := historyTable.GetRowCount()

				historyTable.SetCell(nextRowIndex, 0, tview.NewTableCell(fmt.Sprintf("%-10s", utils.ExtractDateOrTime(tx.Timestamp))).SetAlign(tview.AlignCenter).SetMaxWidth(1).SetExpansion(1).SetMaxWidth(1).SetExpansion(1))
				historyTable.SetCell(nextRowIndex, 1, tview.NewTableCell(fmt.Sprintf("%-10s", tx.Direction.String())).SetAlign(tview.AlignCenter).SetMaxWidth(2).SetExpansion(1))
				historyTable.SetCell(nextRowIndex, 2, tview.NewTableCell(fmt.Sprintf("%15s", formatAmount(tx.Amount))).SetAlign(tview.AlignCenter).SetMaxWidth(3).SetExpansion(1))
				historyTable.SetCell(nextRowIndex, 3, tview.NewTableCell(fmt.Sprintf("%12s", tx.Status)).SetAlign(tview.AlignCenter).SetMaxWidth(1).SetExpansion(1))
				historyTable.SetCell(nextRowIndex, 4, tview.NewTableCell(fmt.Sprintf("%-8s", tx.Type)).SetAlign(tview.AlignCenter).SetMaxWidth(1).SetExpansion(1))

				if len(displayedTxHashes) < totalTxCount {
					displayedTxHashes = append(displayedTxHashes, tx.Hash)
				}
			}

			// clear loading message text
			displayMessage("", InfoMessage)
			
			if len(displayedTxHashes) < totalTxCount {
				nextOffset := txOffset + len(txns)
				historyTable.SetSelectionChangedFunc(func(row, column int) {
					if row >= historyTable.GetRowCount()-10 {
						historyTable.SetSelectionChangedFunc(nil) // unset selection change listener until table is populated
						fetchAndDisplayTransactions(nextOffset, wallet, filter, historyTable, tviewApp, displayMessage)
					}
				})
			} 
		})

		return
	}

	func displayTxDetails(txHash string, wallet walletcore.Wallet, displayError func(string, messageKind), transactionDetailsTable *tview.Table) {
		tx, err := wallet.GetTransaction(txHash)
		if err != nil {
			displayError(err.Error(), MessageKindError)
		}

		transactionDetailsTable.SetCellSimple(0, 0, "Hash")
		transactionDetailsTable.SetCellSimple(1, 0, "Confirmations")
		transactionDetailsTable.SetCellSimple(2, 0, "Included in block")
		transactionDetailsTable.SetCellSimple(3, 0, "Type")
		transactionDetailsTable.SetCellSimple(4, 0, "Amount")
		transactionDetailsTable.SetCellSimple(5, 0, "Date")
		transactionDetailsTable.SetCellSimple(6, 0, "Direction")
		transactionDetailsTable.SetCellSimple(7, 0, "Fee")
		transactionDetailsTable.SetCellSimple(8, 0, "Fee Rate")

		transactionDetailsTable.SetCellSimple(0, 1, tx.Hash)
		transactionDetailsTable.SetCellSimple(1, 1, strconv.Itoa(int(tx.Confirmations)))
		transactionDetailsTable.SetCellSimple(2, 1, strconv.Itoa(int(tx.BlockHeight)))
		transactionDetailsTable.SetCellSimple(3, 1, tx.Type)
		transactionDetailsTable.SetCellSimple(4, 1, dcrutil.Amount(tx.Amount).String())
		transactionDetailsTable.SetCellSimple(5, 1, fmt.Sprintf("%s UTC", tx.LongTime))
		transactionDetailsTable.SetCellSimple(6, 1, tx.Direction.String())
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
		maxDecimalPlacesForInputsAndOutputsAmounts := godcrUtils.MaxDecimalPlaces(inputsAndOutputsAmount)

		// now txFilterDropdownat amount having determined the max number of decimal places
		txFilterDropdownatAmount := func(amount int64) string {
			return godcrUtils.FormatAmountDisplay(amount, maxDecimalPlacesForInputsAndOutputsAmounts)
		}

		transactionDetailsTable.SetCellSimple(9, 0, "-Inputs-")
		for _, txIn := range tx.Inputs {
			row := transactionDetailsTable.GetRowCount()
			transactionDetailsTable.SetCell(row, 0, tview.NewTableCell(txFilterDropdownatAmount(txIn.Amount)).SetAlign(tview.AlignRight))
			transactionDetailsTable.SetCellSimple(row, 1, txIn.PreviousOutpoint)
		}

		row := transactionDetailsTable.GetRowCount()
		transactionDetailsTable.SetCellSimple(row, 0, "-Outputs-")
		for _, txOut := range tx.Outputs {
			row++
			outputAmount := txFilterDropdownatAmount(txOut.Amount)

			if txOut.Address == "" {
				transactionDetailsTable.SetCellSimple(row, 0, fmt.Sprintf("  %s (no address)", outputAmount))
				continue
			}

			transactionDetailsTable.SetCell(row, 0, tview.NewTableCell(outputAmount).SetAlign(tview.AlignRight))
			transactionDetailsTable.SetCellSimple(row, 1, fmt.Sprintf("%s (%s)", txOut.Address, txOut.AccountName))
		}
	}
