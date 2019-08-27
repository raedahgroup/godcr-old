package pages

import (
	"strconv"
	"strings"

	"github.com/raedahgroup/godcr/app/walletcore"

	"fyne.io/fyne"
	"fyne.io/fyne/layout"
	"fyne.io/fyne/widget"
	"github.com/decred/dcrd/dcrutil"
	godcrApp "github.com/raedahgroup/godcr/app"
	"github.com/raedahgroup/godcr/fyne/widgets"
)

type historyPageData struct {
	totalTxOnTable int32
	offset         int32
	currentTxCount int32
	txFilters      *widget.Select
	options        map[string]int
	txTable        widgets.TableStruct
	errorLabel     *widget.Label
	container      *widget.Box
}

var selected bool
var history historyPageData

func historyPageUpdates(wallet godcrApp.WalletMiddleware, window fyne.Window) {
	filters := walletcore.TransactionFilters
	txCountByFilter := make(map[string]int)

	errorHandler := func(err string) {
		history.errorLabel.Show()
		widget.Refresh(history.errorLabel)
		history.errorLabel.SetText(err)
	}

	for _, filter := range filters {
		txCount, err := wallet.TransactionCount(walletcore.BuildTransactionFilter(filter))
		if err != nil {
			errorHandler(err.Error())
			continue
		}
		history.errorLabel.Hide()
		txCountByFilter[filter] = txCount
	}
	var options []string
	count, err := wallet.TransactionCount(nil)
	if err != nil {
		errorHandler(err.Error())
	} else {
		history.errorLabel.Hide()
	}

	txCountByFilter["All"] = count
	options = append(options, "All ("+strconv.Itoa(count)+")")
	options = append(options, "Sent ("+strconv.Itoa(txCountByFilter["Sent"])+")")
	options = append(options, "Received ("+strconv.Itoa(txCountByFilter["Received"])+")")
	options = append(options, "Yourself ("+strconv.Itoa(txCountByFilter["Yourself"])+")")
	options = append(options, "Staking ("+strconv.Itoa(txCountByFilter["Staking"])+")")

	history.txFilters.Options = options
	widget.Refresh(history.txFilters)

	size := history.txTable.Container.Content.Size().Height - history.txTable.Container.Size().Height
	scrollPosition := float64(history.txTable.Container.Offset.Y) / float64(size)

	splittedWord := strings.Split(history.txFilters.Selected, " ")
	var found int32
	if int32(txCountByFilter[splittedWord[0]]) > history.currentTxCount {
		found = int32(txCountByFilter[splittedWord[0]]) - history.currentTxCount
		history.txFilters.Selected = splittedWord[0] + " (" + strconv.Itoa(txCountByFilter[splittedWord[0]]) + ")"
		widget.Refresh(history.txFilters)
	}

	// append to table when scrollbar is at 80% of the scroller.
	if scrollPosition == 1 {
		addToHistoryTable(&history.txTable, history.totalTxOnTable+found, 20, wallet, window, false)
		if txCountByFilter[splittedWord[0]] > int(history.totalTxOnTable+20) {
			history.totalTxOnTable = history.totalTxOnTable + 20
		} else {
			history.totalTxOnTable = int32(txCountByFilter[splittedWord[0]])
		}
		if history.txTable.NumberOfRows() >= 90 {
			history.txTable.Delete(0, 20)
			history.offset = history.offset + 20
		}

	} else if history.txTable.Container.Offset.Y == 0 {
		// if the scroll bar is at the begining, then fetch 1st 50 tx
		if int32(txCountByFilter[splittedWord[0]]) > history.currentTxCount {
			history.txFilters.SetSelected(splittedWord[0] + " (" + strconv.Itoa(txCountByFilter[splittedWord[0]]) + ")")
		}
	} else if scrollPosition < 0.2 {
		if history.offset == 0 {
			return
		}
		addToHistoryTable(&history.txTable, history.offset+found-20, 20+found, wallet, window, true)
		history.offset = history.offset - 20

		rowNo := history.txTable.NumberOfRows()
		if rowNo >= 90 {
			history.txTable.Delete(rowNo-20, rowNo)
			history.totalTxOnTable = int32(rowNo) + history.offset
		}
	}
}

func historyPage(wallet godcrApp.WalletMiddleware, window fyne.Window) {
	history.options = make(map[string]int)
	history.options["All"] = 0
	history.options["Sent"] = 1
	history.options["Received"] = 2
	history.options["Yourself"] = 3
	history.options["Staking"] = 4

	history.errorLabel = widget.NewLabel("")
	history.errorLabel.Hide()

	history.txFilters = widget.NewSelect(nil, func(selected string) {
		// if a new type is selected, load the first 50tx
		var txTable widgets.TableStruct
		fetchTxTable(true, &txTable, 0, 50, wallet, window)
		splittedWord := strings.Split(history.txFilters.Selected, " ")
		currentTxCount, err := wallet.TransactionCount(walletcore.BuildTransactionFilter(splittedWord[0]))
		if err != nil {
			history.errorLabel.Show()
			widget.Refresh(history.errorLabel)
			history.errorLabel.SetText(err.Error())
			return
		}
		history.errorLabel.Hide()

		history.currentTxCount = int32(currentTxCount)
		if currentTxCount >= 50 {
			history.totalTxOnTable = 50
		} else {
			history.totalTxOnTable = int32(currentTxCount)
		}

		history.offset = 0
		history.txTable.Result.Children = txTable.Result.Children
		widget.Refresh(history.txTable.Result)
		history.txTable.Container.Offset.Y = 0
		widget.Refresh(history.txTable.Container)
	})

	historyPageUpdates(wallet, window)
	if overview.errorLabel.Hidden {
		history.txFilters.SetSelected(history.txFilters.Options[0])
	}
	history.txTable.Container.Resize(fyne.NewSize(history.txTable.Container.MinSize().Width, 500))

	output := widget.NewVBox(widget.NewLabelWithStyle("History", fyne.TextAlignLeading, fyne.TextStyle{Bold: true}), widgets.NewVSpacer(10),
		widget.NewHBox(layout.NewSpacer(), history.txFilters),
		fyne.NewContainer(history.txTable.Container),
		history.errorLabel)

	history.container.Children = widget.NewHBox(widgets.NewHSpacer(10), output).Children
	widget.Refresh(history.container)
}

func addToHistoryTable(txTable *widgets.TableStruct, offset, count int32, wallet godcrApp.WalletMiddleware, window fyne.Window, prepend bool) {
	splittedWord := strings.Split(history.txFilters.Selected, " ")
	txs, err := wallet.TransactionHistory(offset, count, walletcore.BuildTransactionFilter(splittedWord[0]))
	if err != nil {
		history.errorLabel.Show()
		widget.Refresh(history.errorLabel)
		history.errorLabel.SetText(err.Error())
		return
	}
	history.errorLabel.Hide()

	var hBox []*widget.Box
	for _, tx := range txs {
		trimmedHash := tx.Hash[:15] + "..." + tx.Hash[len(tx.Hash)-15:]
		hBox = append(hBox, widget.NewHBox(
			widget.NewLabelWithStyle(tx.LongTime, fyne.TextAlignCenter, fyne.TextStyle{}),
			widget.NewLabelWithStyle(tx.Type, fyne.TextAlignCenter, fyne.TextStyle{}),
			widget.NewLabelWithStyle(tx.Direction.String(), fyne.TextAlignLeading, fyne.TextStyle{}),
			widget.NewLabelWithStyle(dcrutil.Amount(tx.Amount).String(), fyne.TextAlignTrailing, fyne.TextStyle{}),
			widget.NewLabelWithStyle(dcrutil.Amount(tx.Fee).String(), fyne.TextAlignCenter, fyne.TextStyle{}),
			widget.NewLabelWithStyle(tx.Status, fyne.TextAlignCenter, fyne.TextStyle{}),
			widget.NewButton(trimmedHash, func() {
				getTxDetails(tx.Hash, wallet, window)
			}),
		))
	}
	if prepend {
		history.txTable.Prepend(hBox...)
	} else {
		history.txTable.Append(hBox...)
	}
	widget.Refresh(history.txTable.Container)
}

func getTxDetails(hash string, wallet godcrApp.WalletMiddleware, window fyne.Window) {
	txDetails, err := wallet.GetTransaction(hash)
	if err != nil {
		history.errorLabel.Show()
		widget.Refresh(history.errorLabel)
		history.errorLabel.SetText(err.Error())
		return
	}
	history.errorLabel.Hide()

	label := widget.NewLabelWithStyle("Transaction Details", fyne.TextAlignLeading, fyne.TextStyle{Bold: true, Italic: true})
	form := widget.NewForm()
	form.Append("Date", widget.NewLabel(txDetails.LongTime))
	form.Append("Status", widget.NewLabel(txDetails.Status))
	form.Append("Amount", widget.NewLabel(dcrutil.Amount(txDetails.Amount).String()))
	form.Append("Fee", widget.NewLabel(dcrutil.Amount(txDetails.Fee).String()))
	form.Append("Fee Rate", widget.NewLabel(dcrutil.Amount(txDetails.FeeRate).String()))
	form.Append("Type", widget.NewLabel(txDetails.Type))
	form.Append("Confirmation", widget.NewLabel(strconv.Itoa(int(txDetails.Confirmations))))
	form.Append("Hash", widget.NewLabel(txDetails.Hash))

	var txInput widgets.TableStruct
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

	var txOutput widgets.TableStruct
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

	output := widget.NewHBox(widgets.NewHSpacer(10), widget.NewVBox(
		form,
		widget.NewLabelWithStyle("Inputs", fyne.TextAlignLeading, fyne.TextStyle{Bold: true}),
		txInput.Result,
		widgets.NewVSpacer(10),
		widget.NewLabelWithStyle("Outputs", fyne.TextAlignLeading, fyne.TextStyle{Bold: true}),
		txOutput.Result,
	), widgets.NewHSpacer(10))

	scrollContainer := widget.NewScrollContainer(output)
	scrollContainer.Resize(fyne.NewSize(scrollContainer.MinSize().Width, 500))
	popUp := widget.NewPopUp(widget.NewVBox(label, fyne.NewContainer(scrollContainer)),
		window.Canvas())
	popUp.Show()
}
