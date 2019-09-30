package pages

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/decred/dcrd/dcrutil"
	"github.com/gdamore/tcell"
	"github.com/raedahgroup/dcrlibwallet/utils"
	godcrUtils "github.com/raedahgroup/godcr/app/utils"
	"github.com/raedahgroup/godcr/app/walletcore"
	"github.com/raedahgroup/godcr/terminal/helpers"
	"github.com/raedahgroup/godcr/terminal/primitives"
	"github.com/rivo/tview"
)

type messageKind uint32

const (
	MessageKindError messageKind = iota
	MessageKindInfo
)

var historyPageData struct {
	pageContentHolder *tview.Flex
	titleTextView     *primitives.TextView
	hintTextView      *primitives.TextView

	txPerPage    int32
	totalTxCount int

	txFilterDropDown  *primitives.Form
	currentTxFilter   string
	historyTable      *tview.Table
	displayedTxHashes []string

	transactionDetailsTable *tview.Table

	displayMessage func(string, messageKind)
}

func historyPage(wallet walletcore.Wallet, hintTextView *primitives.TextView, tviewApp *tview.Application, clearFocus func()) tview.Primitive {
	// setup initial page data properties
	historyPageData.pageContentHolder = tview.NewFlex().SetDirection(tview.FlexRow)
	historyPageData.titleTextView = primitives.NewLeftAlignedTextView("History")
	historyPageData.hintTextView = hintTextView

	historyPageData.displayMessage = messageDisplayFn(tviewApp)

	// handler for returning back to menu column
	historyPageData.pageContentHolder.SetInputCapture(func(event *tcell.EventKey) *tcell.EventKey {
		if event.Key() == tcell.KeyEscape || event.Key() == tcell.KeyBackspace || event.Key() == tcell.KeyBackspace2 {
			clearFocus()
			return nil
		}
		return event
	})

	historyPageData.pageContentHolder.AddItem(historyPageData.titleTextView, 2, 0, false)

	// get total tx count before setting up further page data props
	totalTxCount, txCountErr := wallet.TransactionCount(nil)
	if txCountErr != nil {
		errorMessage := fmt.Sprintf("Cannot load history page. Error getting total transaction count: %v", txCountErr)
		historyPageData.displayMessage(errorMessage, MessageKindError)
		tviewApp.SetFocus(historyPageData.pageContentHolder)
		return historyPageData.pageContentHolder
	}
	if totalTxCount == 0 {
		historyPageData.displayMessage("No transactions yet", MessageKindInfo)
		hintTextView.SetText("TIP: ESC or BACKSPACE to return to navigation menu")
		tviewApp.SetFocus(historyPageData.pageContentHolder)
		return historyPageData.pageContentHolder
	}

	// setup tx filter dropdown
	historyPageData.txFilterDropDown = prepareTxFilterDropDown(wallet, tviewApp, clearFocus)
	if historyPageData.txFilterDropDown != nil {
		historyPageData.pageContentHolder.AddItem(historyPageData.txFilterDropDown, 2, 0, false)
	} else {
		tviewApp.SetFocus(historyPageData.pageContentHolder)
		return historyPageData.pageContentHolder
	}

	// setup more page data props
	historyPageData.txPerPage = walletcore.TransactionHistoryCountPerPage
	historyPageData.totalTxCount = totalTxCount
	historyPageData.historyTable = prepareHistoryTable(wallet, tviewApp, clearFocus)
	historyPageData.displayedTxHashes = nil
	historyPageData.transactionDetailsTable = prepareTxDetailsTable(tviewApp)

	// fetch tx to display in subroutine so the UI isn't blocked
	go fetchAndDisplayTransactions(0, wallet, "All", tviewApp)
	historyPageData.pageContentHolder.AddItem(historyPageData.historyTable, 0, 1, true)
	hintTextView.SetText("TIP: Use TAB to switch/navigate, ARROW UP/DOWN to select txn, \nENTER to view details," +
		" ESC to return to navigation menu")

	tviewApp.SetFocus(historyPageData.pageContentHolder)
	return historyPageData.pageContentHolder
}

func messageDisplayFn(tviewApp *tview.Application) func(string, messageKind) {
	messageTextView := primitives.WordWrappedTextView("")
	return func(message string, kind messageKind) {
		// this function may be called from a goroutine, use tviewApp.QueueUpdateDraw
		tviewApp.QueueUpdateDraw(func() {
			historyPageData.pageContentHolder.RemoveItem(messageTextView)

			if message != "" {
				messageTextView.SetText(message)
				if kind == MessageKindError {
					messageTextView.SetTextColor(helpers.DecredOrangeColor)
				} else {
					messageTextView.SetTextColor(tcell.ColorWhite)
				}
				historyPageData.pageContentHolder.AddItem(messageTextView, 2, 0, false)
			}
		})
	}
}

func prepareTxFilterDropDown(wallet walletcore.Wallet, tviewApp *tview.Application, clearFocus func()) *primitives.Form {
	txFilterDropDown := primitives.NewForm(false)
	txFilterDropDown.SetBorderPadding(0, 0, 0, 0)

	var txFilterSelectionOptions []string
	for _, filter := range walletcore.TransactionFilters {
		txCountForFilter, txCountErr := wallet.TransactionCount(walletcore.BuildTransactionFilter(filter))
		if txCountErr != nil {
			errorMessage := fmt.Sprintf("Cannot load history page. Error getting transaction count for filter %s: %s",
				filter, txCountErr.Error())
			historyPageData.displayMessage(errorMessage, MessageKindError)
			return nil
		}

		if txCountForFilter > 0 {
			txFilterSelectionOptions = append(txFilterSelectionOptions, fmt.Sprintf("%s (%d)", filter, txCountForFilter))
		}
	}

	// dropDown selection change listener
	txFilterDropDown.AddDropDown("", txFilterSelectionOptions, 0, func(selectedOption string, index int) {
		selectedFilter := strings.Split(selectedOption, " ")[0]
		if selectedFilter != historyPageData.currentTxFilter {
			go fetchAndDisplayTransactions(0, wallet, selectedFilter, tviewApp)
		}
	})

	// handler for switching between dropDown and table
	txFilterDropDown.SetInputCapture(func(event *tcell.EventKey) *tcell.EventKey {
		if event.Key() == tcell.KeyEscape || event.Key() == tcell.KeyBackspace || event.Key() == tcell.KeyBackspace2 {
			clearFocus()
			return nil
		}
		if event.Key() == tcell.KeyTab {
			tviewApp.SetFocus(historyPageData.historyTable)
			return nil
		}
		return event
	})

	return txFilterDropDown
}

func prepareHistoryTable(wallet walletcore.Wallet, tviewApp *tview.Application, clearFocus func()) *tview.Table {
	historyTable := tview.NewTable().
		SetBorders(false).
		SetFixed(1, 0). // keep first row (column headers) fixed during scroll
		SetSelectable(true, false)

	historyTable.SetDoneFunc(func(key tcell.Key) {
		if key == tcell.KeyEscape {
			clearFocus()
		}
	})

	// handler for switching between dropDown and table
	historyTable.SetInputCapture(func(event *tcell.EventKey) *tcell.EventKey {
		if event.Key() == tcell.KeyEscape || event.Key() == tcell.KeyBackspace || event.Key() == tcell.KeyBackspace2 {
			clearFocus()
			return nil
		}
		if event.Key() == tcell.KeyTab {
			tviewApp.SetFocus(historyPageData.txFilterDropDown)
			return nil
		}
		return event
	})

	// method for getting transaction details when a tx is selected from the history table
	historyTable.SetSelectedFunc(func(row, column int) {
		if row >= len(historyPageData.displayedTxHashes) {
			// ignore selected func call for table header
			return
		}

		historyPageData.pageContentHolder.RemoveItem(historyTable)
		historyPageData.pageContentHolder.RemoveItem(historyPageData.txFilterDropDown)

		historyPageData.titleTextView.SetText("Transaction Details")
		historyPageData.hintTextView.SetText("TIP: Use ARROW UP/DOWN to scroll, \nBACKSPACE to view History page, ESC to return to navigation menu")

		historyPageData.transactionDetailsTable.Clear()
		historyPageData.pageContentHolder.AddItem(historyPageData.transactionDetailsTable, 0, 1, true)

		tviewApp.SetFocus(historyPageData.transactionDetailsTable)

		txHash := historyPageData.displayedTxHashes[row-1]
		displayTxDetails(txHash, wallet)
	})

	return historyTable
}

func prepareTxDetailsTable(tviewApp *tview.Application) *tview.Table {
	transactionDetailsTable := tview.NewTable().SetBorders(false)

	// handler for returning back to history table from tx details table
	transactionDetailsTable.SetInputCapture(func(event *tcell.EventKey) *tcell.EventKey {
		if event.Key() == tcell.KeyEscape || event.Key() == tcell.KeyBackspace || event.Key() == tcell.KeyBackspace2 {
			historyPageData.pageContentHolder.AddItem(historyPageData.txFilterDropDown, 2, 0, false)

			historyPageData.pageContentHolder.RemoveItem(historyPageData.transactionDetailsTable)

			historyPageData.titleTextView.SetText("History")
			historyPageData.hintTextView.SetText("TIP: Use ARROW UP/DOWN to select txn,\nENTER to view details, " +
				"ESC to return to navigation menu")

			historyPageData.pageContentHolder.AddItem(historyPageData.historyTable, 0, 1, true)
			tviewApp.SetFocus(historyPageData.historyTable)

			return nil
		}
		return event
	})

	return transactionDetailsTable
}

func fetchAndDisplayTransactions(txOffset int, wallet walletcore.Wallet, selectedFilter string, tviewApp *tview.Application) {
	// show a loading text at the bottom of the table so user knows an op is in progress
	historyPageData.displayMessage("Fetching data...", MessageKindInfo)

	if selectedFilter != historyPageData.currentTxFilter {
		historyPageData.historyTable.Clear()
		historyPageData.displayedTxHashes = nil
		txOffset = 0
	}

	historyPageData.currentTxFilter = selectedFilter
	filter := walletcore.BuildTransactionFilter(selectedFilter)

	txns, err := wallet.TransactionHistory(int32(txOffset), historyPageData.txPerPage, filter)
	if err != nil {
		historyPageData.displayMessage(err.Error(), MessageKindError)
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

	// history table header
	historyPageData.historyTable.SetCell(0, 0, tableHeaderCell("Date (UTC)"))
	historyPageData.historyTable.SetCell(0, 1, tableHeaderCell(fmt.Sprintf("%10s", "Direction")))
	historyPageData.historyTable.SetCell(0, 2, tableHeaderCell(fmt.Sprintf("%8s", "Amount")))
	historyPageData.historyTable.SetCell(0, 3, tableHeaderCell(fmt.Sprintf("%5s", "Status")))
	historyPageData.historyTable.SetCell(0, 4, tableHeaderCell(fmt.Sprintf("%-5s", "Type")))

	// updating the history table from a goroutine, use tviewApp.QueueUpdateDraw
	tviewApp.QueueUpdateDraw(func() {
		for _, tx := range txns {
			nextRowIndex := historyPageData.historyTable.GetRowCount()

			historyPageData.historyTable.SetCell(nextRowIndex, 0, tview.NewTableCell(fmt.Sprintf("%-10s", utils.ExtractDateOrTime(tx.Timestamp))).SetAlign(tview.AlignCenter).SetMaxWidth(1).SetExpansion(1).SetMaxWidth(1).SetExpansion(1))
			historyPageData.historyTable.SetCell(nextRowIndex, 1, tview.NewTableCell(fmt.Sprintf("%-10s", tx.Direction.String())).SetAlign(tview.AlignCenter).SetMaxWidth(2).SetExpansion(1))
			historyPageData.historyTable.SetCell(nextRowIndex, 2, tview.NewTableCell(fmt.Sprintf("%15s", formatAmount(tx.Amount))).SetAlign(tview.AlignCenter).SetMaxWidth(3).SetExpansion(1))
			historyPageData.historyTable.SetCell(nextRowIndex, 3, tview.NewTableCell(fmt.Sprintf("%12s", tx.Status)).SetAlign(tview.AlignCenter).SetMaxWidth(1).SetExpansion(1))
			historyPageData.historyTable.SetCell(nextRowIndex, 4, tview.NewTableCell(fmt.Sprintf("%-8s", tx.Type)).SetAlign(tview.AlignCenter).SetMaxWidth(1).SetExpansion(1))

			historyPageData.displayedTxHashes = append(historyPageData.displayedTxHashes, tx.Hash)
		}

		// clear loading message text
		historyPageData.displayMessage("", MessageKindInfo)

		if len(historyPageData.displayedTxHashes) < historyPageData.totalTxCount {
			nextOffset := txOffset + len(txns)
			historyPageData.historyTable.SetSelectionChangedFunc(func(row, column int) {
				if row >= historyPageData.historyTable.GetRowCount()-10 {
					historyPageData.historyTable.SetSelectionChangedFunc(nil) // unset selection change listener until table is populated
					fetchAndDisplayTransactions(nextOffset, wallet, historyPageData.currentTxFilter, tviewApp)
				}
			})
		}
	})

	return
}

func displayTxDetails(txHash string, wallet walletcore.Wallet) {
	tx, err := wallet.GetTransaction(txHash)
	if err != nil {
		historyPageData.displayMessage(err.Error(), MessageKindError)
	}

	historyPageData.transactionDetailsTable.SetCellSimple(0, 0, "Hash")
	historyPageData.transactionDetailsTable.SetCellSimple(1, 0, "Confirmations")
	historyPageData.transactionDetailsTable.SetCellSimple(2, 0, "Included in block")
	historyPageData.transactionDetailsTable.SetCellSimple(3, 0, "Type")
	historyPageData.transactionDetailsTable.SetCellSimple(4, 0, "Amount")
	historyPageData.transactionDetailsTable.SetCellSimple(5, 0, "Date")
	historyPageData.transactionDetailsTable.SetCellSimple(6, 0, "Direction")
	historyPageData.transactionDetailsTable.SetCellSimple(7, 0, "Fee")
	historyPageData.transactionDetailsTable.SetCellSimple(8, 0, "Fee Rate")

	historyPageData.transactionDetailsTable.SetCellSimple(0, 1, tx.Hash)
	historyPageData.transactionDetailsTable.SetCellSimple(1, 1, strconv.Itoa(int(tx.Confirmations)))
	historyPageData.transactionDetailsTable.SetCellSimple(2, 1, strconv.Itoa(int(tx.BlockHeight)))
	historyPageData.transactionDetailsTable.SetCellSimple(3, 1, tx.Type)
	historyPageData.transactionDetailsTable.SetCellSimple(4, 1, dcrutil.Amount(tx.Amount).String())
	historyPageData.transactionDetailsTable.SetCellSimple(5, 1, fmt.Sprintf("%s UTC", tx.LongTime))
	historyPageData.transactionDetailsTable.SetCellSimple(6, 1, tx.Direction.String())
	historyPageData.transactionDetailsTable.SetCellSimple(7, 1, dcrutil.Amount(tx.Fee).String())
	historyPageData.transactionDetailsTable.SetCellSimple(8, 1, fmt.Sprintf("%s/kB", dcrutil.Amount(tx.FeeRate)))

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

	historyPageData.transactionDetailsTable.SetCellSimple(9, 0, "-Inputs-")
	for _, txIn := range tx.Inputs {
		row := historyPageData.transactionDetailsTable.GetRowCount()
		historyPageData.transactionDetailsTable.SetCell(row, 0, tview.NewTableCell(txFilterDropdownatAmount(txIn.Amount)).SetAlign(tview.AlignRight))
		historyPageData.transactionDetailsTable.SetCellSimple(row, 1, txIn.PreviousOutpoint)
	}

	row := historyPageData.transactionDetailsTable.GetRowCount()
	historyPageData.transactionDetailsTable.SetCellSimple(row, 0, "-Outputs-")
	for _, txOut := range tx.Outputs {
		row++
		outputAmount := txFilterDropdownatAmount(txOut.Amount)

		if txOut.Address == "" {
			historyPageData.transactionDetailsTable.SetCellSimple(row, 0, fmt.Sprintf("  %s (no address)", outputAmount))
			continue
		}

		historyPageData.transactionDetailsTable.SetCell(row, 0, tview.NewTableCell(outputAmount).SetAlign(tview.AlignRight))
		historyPageData.transactionDetailsTable.SetCellSimple(row, 1, fmt.Sprintf("%s (%s)", txOut.Address, txOut.AccountName))
	}
}
