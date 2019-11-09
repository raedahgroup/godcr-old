package pages

import (
	"fmt"
	"image/color"
	"strconv"
	"strings"
	"time"

	"fyne.io/fyne"
	"fyne.io/fyne/canvas"
	"fyne.io/fyne/layout"
	"fyne.io/fyne/theme"
	"fyne.io/fyne/widget"

	"github.com/decred/dcrd/chaincfg/chainhash"
	"github.com/decred/dcrd/dcrutil"
	"github.com/raedahgroup/dcrlibwallet"
	"github.com/raedahgroup/godcr/fyne/assets"
	"github.com/raedahgroup/godcr/fyne/widgets"
)

const txPerPage int32 = 25

const txPerPage int32 = 25

type txHistoryPageData struct {
	txTable          widgets.Table
	allTxCount       int
	selectedFilterId int32
	errorLabel       *widget.Label
	TotalTxFetched   int32
}

var txHistory txHistoryPageData


func HistoryPageContent(wallet *dcrlibwallet.LibWallet, window fyne.Window, tabmenu *widget.TabContainer) fyne.CanvasObject {
	txHistory.errorLabel = widget.NewLabelWithStyle("", fyne.TextAlignCenter, fyne.TextStyle{Bold: true})
	txHistory.errorLabel.Hide()

	pageTitleLabel := widget.NewLabelWithStyle("Transactions", fyne.TextAlignLeading, fyne.TextStyle{Bold: true, Italic: true})

	txHistoryPageOutput := widget.NewVBox(
		widgets.NewVSpacer(5),
		widget.NewHBox(pageTitleLabel),
		widgets.NewVSpacer(5),
	)

	txFilterDropDown, errorMessage := txFilterDropDown(wallet, window, tabmenu)
	if errorMessage != "" {
		errorHandler(errorMessage, txHistory.errorLabel)
		txHistoryPageOutput.Append(txHistory.errorLabel)
		return widget.NewHBox(widgets.NewHSpacer(18), txHistoryPageOutput)
	}

	txSortFilterDropDown, errorMessage := txSortDropDown(window)
	if errorMessage != "" {
		errorHandler(errorMessage, txHistory.errorLabel)
		txHistoryPageOutput.Append(txHistory.errorLabel)
		return widget.NewHBox(widgets.NewHSpacer(18), txHistoryPageOutput)
	}

	txTableHeader(wallet, &txHistory.txTable, window)
	fetchTx(&txHistory.txTable, 0, dcrlibwallet.TxFilterAll, wallet, window, false, tabmenu)

	txHistoryPageOutput.Append(widget.NewHBox(txSortFilterDropDown, widgets.NewHSpacer(30), txFilterDropDown))
	txHistoryPageOutput.Append(txHistory.errorLabel)
	txHistoryPageOutput.Append(widgets.NewVSpacer(5))
	txHistoryPageOutput.Append(fyne.NewContainerWithLayout(layout.NewFixedGridLayout(fyne.NewSize(txHistory.txTable.Container.MinSize().Width, txHistory.txTable.Container.MinSize().Height+450)), txHistory.txTable.Container))
	txHistoryPageOutput.Append(widgets.NewVSpacer(15))

	return widget.NewHBox(widgets.NewHSpacer(18), txHistoryPageOutput)
}

func txFilterDropDown(wallet *dcrlibwallet.LibWallet, window fyne.Window, tabmenu *widget.TabContainer) (*widgets.ClickableBox, string) {
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

	txCountForFilter, err := wallet.CountTransactions(allTxFilters["All"])
	if err != nil {
		errorMessage := fmt.Sprintf("Cannot load txHistory page. Error getting transaction count for filter All: %s", err.Error())
		return nil, errorMessage
	}
	txHistory.allTxCount = txCountForFilter

	selectedTxFilterLabel := widget.NewLabel(fmt.Sprintf("%s (%d)", "All", txCountForFilter))

	var txFilterSelectionPopup *widget.PopUp
	txFilterListWidget := widget.NewVBox()
	for _, filterName := range allTxFilterNames {
		filterId := allTxFilters[filterName]
		txCountForFilter, err := wallet.CountTransactions(filterId)
		if err != nil {
			errorMessage := fmt.Sprintf("Cannot load txHistory page. Error getting transaction count for filter %s: %s",
				filterName, err.Error())
			return nil, errorMessage
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
				if allTxCountForSelectedTx, err := wallet.CountTransactions(selectedFilterId); err == nil {
					txHistory.allTxCount = allTxCountForSelectedTx
				}

				if selectedFilterId != txHistory.selectedFilterId {
					selectedTxFilterLabel.SetText(filter)
					txTableHeader(wallet, &txTable, window)
					txHistory.txTable.Result.Children = txTable.Result.Children
					fetchTx(&txTable, 0, selectedFilterId, wallet, window, false, tabmenu)
					widget.Refresh(txHistory.txTable.Result)
				}

				txFilterSelectionPopup.Hide()
			}))
		}
	}

	// txFilterSelectionPopup create a popup that has tx filter name and tx count
	txFilterSelectionPopup = widget.NewPopUp(widget.NewVBox(txFilterListWidget), window.Canvas())
	txFilterSelectionPopup.Hide()

	icons, err := assets.GetIcons(assets.CollapseIcon)
	if err != nil {
		errorMessage := fmt.Sprintf("Error: %s", err.Error())
		return nil, errorMessage
	}

	txFilterTab := widget.NewHBox(
		selectedTxFilterLabel,
		widgets.NewHSpacer(60),
		widget.NewIcon(icons[assets.CollapseIcon]),
	)

	var txFilterDropDown *widgets.ClickableBox
	txFilterDropDown = widgets.NewClickableBox(txFilterTab, func() {
		txFilterSelectionPopup.Move(fyne.CurrentApp().Driver().AbsolutePositionForObject(
			txFilterDropDown).Add(fyne.NewPos(0, txFilterDropDown.Size().Height)))
		txFilterSelectionPopup.Show()
	})

	return txFilterDropDown, ""
}

func txSortDropDown(window fyne.Window) (*widgets.ClickableBox, string) {
	var allTxSortNames = []string{"Newest", "Oldest"}
	var allTxSortFilters = map[string]int32{
		"Newest": 0,
		"Oldest": 1,
	}

	selectedTxSortFilterLabel := widget.NewLabel("Newest")

	var txSortFilterSelectionPopup *widget.PopUp
	txSortFilterListWidget := widget.NewVBox()
	for _, sortName := range allTxSortNames {
		sortId := allTxSortFilters[sortName]
		fmt.Println(sortId)

		txSortView := widget.NewHBox(
			widgets.NewHSpacer(5),
			widget.NewLabel(sortName),
			widgets.NewHSpacer(5),
		)

		txSortFilterListWidget.Append(widgets.NewClickableBox(txSortView, func() {
			selectedTxSortFilterLabel.SetText(sortName)
			txSortFilterSelectionPopup.Hide()
		}))
	}

	// txSortFilterSelectionPopup create a popup that has tx filter name and tx count
	txSortFilterSelectionPopup = widget.NewPopUp(widget.NewVBox(txSortFilterListWidget), window.Canvas())
	txSortFilterSelectionPopup.Hide()

	icons, err := assets.GetIcons(assets.CollapseIcon)
	if err != nil {
		errorMessage := fmt.Sprintf("Error: %s", err.Error())
		return nil, errorMessage
	}

	txSortFilterTab := widget.NewHBox(
		selectedTxSortFilterLabel,
		widgets.NewHSpacer(10),
		widget.NewIcon(icons[assets.CollapseIcon]),
	)

	var txSortFilterDropDown *widgets.ClickableBox
	txSortFilterDropDown = widgets.NewClickableBox(txSortFilterTab, func() {
		txSortFilterSelectionPopup.Move(fyne.CurrentApp().Driver().AbsolutePositionForObject(
			txSortFilterDropDown).Add(fyne.NewPos(0, txSortFilterDropDown.Size().Height)))
		txSortFilterSelectionPopup.Show()
	})

	return txSortFilterDropDown, ""
}

func txTableHeader(wallet *dcrlibwallet.LibWallet, txTable *widgets.Table, window fyne.Window) {
	tableHeading := widget.NewHBox(
		widget.NewLabelWithStyle("Date (UTC)", fyne.TextAlignCenter, fyne.TextStyle{Bold: true}),
		widget.NewLabelWithStyle("Direction", fyne.TextAlignCenter, fyne.TextStyle{Bold: true}),
		widget.NewLabelWithStyle("Status", fyne.TextAlignCenter, fyne.TextStyle{Bold: true}),
		widget.NewLabelWithStyle("Amount", fyne.TextAlignCenter, fyne.TextStyle{Bold: true}),
		widget.NewLabelWithStyle("Fee", fyne.TextAlignCenter, fyne.TextStyle{Bold: true}),
		widget.NewLabelWithStyle("Type", fyne.TextAlignCenter, fyne.TextStyle{Bold: true}),
		widget.NewLabelWithStyle("Hash", fyne.TextAlignCenter, fyne.TextStyle{Bold: true}),
	)

	var hBox []*widget.Box

	txTable.NewTable(tableHeading, hBox...)

	return
}

func fetchTx(txTable *widgets.Table, txOffset, filter int32, wallet *dcrlibwallet.LibWallet, window fyne.Window, prepend bool, tabmenu *widget.TabContainer) {
	if filter != txHistory.selectedFilterId {
		txOffset = 0
		txHistory.TotalTxFetched = 0
		txHistory.selectedFilterId = filter
	}

	txns, err := wallet.GetTransactionsRaw(txOffset, txPerPage, filter)
	if err != nil {
		errorHandler(fmt.Sprintf("Error getting transaction for Filter: %s", err.Error()), txHistory.errorLabel)
		txHistory.txTable.Container.Hide()
		return
	}
	if len(txns) == 0 {
		errorHandler("No transactions yet.", txHistory.errorLabel)
		txHistory.txTable.Container.Hide()
		return
	}

	txHistory.TotalTxFetched += int32(len(txns))

	var txBox []*widget.Box
	for _, tx := range txns {
		status := "Pending"
		confirmations := wallet.GetBestBlock() - tx.BlockHeight + 1
		if tx.BlockHeight != -1 && confirmations > dcrlibwallet.DefaultRequiredConfirmations {
			status = "Confirmed"
		}

		trimmedHash := tx.Hash[:15] + "..." + tx.Hash[len(tx.Hash)-15:]
		txForTrimmedHash := tx.Hash
		txBox = append(txBox, widget.NewHBox(
			widget.NewLabelWithStyle(dcrlibwallet.ExtractDateOrTime(tx.Timestamp), fyne.TextAlignCenter, fyne.TextStyle{}),
			widget.NewLabelWithStyle(dcrlibwallet.TransactionDirectionName(tx.Direction), fyne.TextAlignCenter, fyne.TextStyle{}),
			widget.NewLabelWithStyle(status, fyne.TextAlignLeading, fyne.TextStyle{}),
			widget.NewLabelWithStyle(dcrutil.Amount(tx.Amount).String(), fyne.TextAlignTrailing, fyne.TextStyle{}),
			widget.NewLabelWithStyle(dcrutil.Amount(tx.Fee).String(), fyne.TextAlignCenter, fyne.TextStyle{}),
			widget.NewLabelWithStyle(tx.Type, fyne.TextAlignCenter, fyne.TextStyle{}),
			widgets.NewClickableBox(widget.NewHBox(widget.NewLabelWithStyle(trimmedHash, fyne.TextAlignCenter, fyne.TextStyle{Italic: true})), func() {
				fetchTxDetails(txForTrimmedHash, wallet, window, txHistory.errorLabel, tabmenu)
			}),
		))
	}

	if prepend {
		txTable.Prepend(txBox...)
	} else {
		txTable.Append(txBox...)
	}

	txHistory.txTable.Result.Children = txTable.Result.Children
	widget.Refresh(txHistory.txTable.Result)
	widget.Refresh(txHistory.txTable.Container)

	time.AfterFunc(time.Second*8, func() {
		updateTable(wallet, window, tabmenu)
	})

	txHistory.errorLabel.Hide()
}

func updateTable(wallet *dcrlibwallet.LibWallet, window fyne.Window, tabmenu *widget.TabContainer) {
	size := txHistory.txTable.Container.Content.Size().Height - txHistory.txTable.Container.Size().Height
	scrollPosition := float64(txHistory.txTable.Container.Offset.Y) / float64(size)
	txTableRowCount := txHistory.txTable.NumberOfColumns()

	if txHistory.allTxCount > int(txHistory.TotalTxFetched) {
		if txHistory.txTable.Container.Offset.Y == 0 {

			time.AfterFunc(time.Second*8, func() {
				updateTable(wallet, window, tabmenu)
			})
		} else if scrollPosition < 0.5 {
			if txHistory.TotalTxFetched <= txPerPage {
				time.AfterFunc(time.Second*8, func() {
					updateTable(wallet, window, tabmenu)
				})
			}
			if txTableRowCount <= int(txPerPage) {
				return
			}

			txHistory.TotalTxFetched -= int32(txPerPage)

			fetchTx(&txHistory.txTable, txHistory.TotalTxFetched, txHistory.selectedFilterId, wallet, window, true, tabmenu)
		} else if scrollPosition >= 0.5 {
			fetchTx(&txHistory.txTable, txHistory.TotalTxFetched, txHistory.selectedFilterId, wallet, window, false, tabmenu)
			if txTableRowCount > 12 {
				txHistory.txTable.Delete(0, txTableRowCount-int(txPerPage))
			}
		}
	}
}

func fetchTxDetails(hash string, wallet *dcrlibwallet.LibWallet, window fyne.Window, errorLabel *widget.Label, tabmenu *widget.TabContainer) {
	messageLabel := widget.NewLabelWithStyle("Fetching data..", fyne.TextAlignCenter, fyne.TextStyle{})
	if messageLabel.Hidden == false {
		time.AfterFunc(time.Millisecond*300, func() {
			if tabmenu.CurrentTabIndex() == 1 {
				messageLabel.Hide()
			}
		})
	}

	var txDetailsPopUp *widget.PopUp

	txDetailslabel := widget.NewLabelWithStyle("Transaction Details", fyne.TextAlignLeading, fyne.TextStyle{Bold: true, Italic: true})

	minimizeIcon := widgets.NewClickableIcon(theme.CancelIcon(), nil, func() { txDetailsPopUp.Hide() })
	errorMessageLabel := widget.NewLabelWithStyle("", fyne.TextAlignCenter, fyne.TextStyle{Bold: true})

	txDetailsErrorMethod := func() {
		txErrorDetailsOutput := widget.NewVBox(
			widgets.NewHSpacer(10),
			widget.NewHBox(
				txDetailslabel,
				widgets.NewHSpacer(150),
				messageLabel,
				layout.NewSpacer(),
				minimizeIcon,
				widgets.NewHSpacer(30),
			),
			errorMessageLabel,
		)
		txDetailsPopUp = widget.NewModalPopUp(widget.NewVBox(fyne.NewContainer(txErrorDetailsOutput)),
			window.Canvas())
		txDetailsPopUp.Show()
	}

	chainHash, err := chainhash.NewHashFromStr(hash)
	if err != nil {
		errorHandler(fmt.Sprintf("fetching generating chainhash from for \n %s \n %s ", hash, err.Error()), errorMessageLabel)
		txDetailsErrorMethod()
		return
	}

	txDetails, err := wallet.GetTransactionRaw(chainHash[:])
	if err != nil {
		errorHandler(fmt.Sprintf("Error fetching transaction details for \n %s \n %s ", hash, err.Error()), errorMessageLabel)
		txDetailsErrorMethod()
		return
	}

	var confirmations int32 = 0
	if txDetails.BlockHeight != -1 {
		confirmations = wallet.GetBestBlock() - txDetails.BlockHeight + 1
	}

	var status string
	var spendUnconfirmed = wallet.ReadBoolConfigValueForKey(dcrlibwallet.SpendUnconfirmedConfigKey)
	if spendUnconfirmed || confirmations > dcrlibwallet.DefaultRequiredConfirmations {
		status = "Confirmed"
	} else {
		status = "Pending"
	}

	tableConfirmations := widget.NewHBox(
		widget.NewLabelWithStyle("Confirmations:", fyne.TextAlignTrailing, fyne.TextStyle{Bold: true}),
		widget.NewLabelWithStyle(strconv.Itoa(int(confirmations)), fyne.TextAlignCenter, fyne.TextStyle{}),
	)

	copyAbleText := func(text string, copyAble bool) *widgets.ClickableBox {
		var textToCopy *canvas.Text
		if copyAble {
			textToCopy = canvas.NewText(text, color.RGBA{0, 255, 255, 1})
		} else {
			textToCopy = canvas.NewText(text, color.RGBA{255, 255, 255, 1})
		}
		textToCopy.TextSize = 14
		textToCopy.Alignment = fyne.TextAlignTrailing

		return widgets.NewClickableBox(widget.NewHBox(textToCopy),
			func() {
				messageLabel.SetText("Data Copied")
				clipboard := window.Clipboard()
				clipboard.SetContent(text)
				messageLabel.Show()

				if messageLabel.Hidden == false {
					time.AfterFunc(time.Second*2, func() {
						if tabmenu.CurrentTabIndex() == 1 {
							messageLabel.Hide()
						}
					})
				}
			},
		)
	}

	tableHash := widget.NewHBox(
		widget.NewLabelWithStyle("Transaction ID:", fyne.TextAlignTrailing, fyne.TextStyle{Bold: true}),
		copyAbleText(txDetails.Hash, true),
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
	inputTableColumnLabels := widget.NewHBox(
		widget.NewLabelWithStyle("Previous Outpoint", fyne.TextAlignCenter, fyne.TextStyle{Bold: true}),
		widget.NewLabelWithStyle("Account", fyne.TextAlignCenter, fyne.TextStyle{Bold: true}),
		widget.NewLabelWithStyle("Amount", fyne.TextAlignCenter, fyne.TextStyle{Bold: true}))

	var inputBox []*widget.Box
	for i := range txDetails.Inputs {
		inputBox = append(inputBox, widget.NewHBox(
			copyAbleText(txDetails.Inputs[i].PreviousOutpoint, true),
			copyAbleText(txDetails.Inputs[i].AccountName, false),
			copyAbleText(dcrutil.Amount(txDetails.Inputs[i].Amount).String(), false),
		))
	}
	txInput.NewTable(inputTableColumnLabels, inputBox...)

	var txOutput widgets.Table
	outputTableColumnLabels := widget.NewHBox(
		widget.NewLabelWithStyle("Address", fyne.TextAlignCenter, fyne.TextStyle{Bold: true}),
		widget.NewLabelWithStyle("Account", fyne.TextAlignCenter, fyne.TextStyle{Bold: true}),
		widget.NewLabelWithStyle("Value", fyne.TextAlignCenter, fyne.TextStyle{Bold: true}),
		widget.NewLabelWithStyle("Type", fyne.TextAlignCenter, fyne.TextStyle{Bold: true}))

	var outputBox []*widget.Box
	for i := range txDetails.Outputs {
		outputBox = append(outputBox, widget.NewHBox(
			copyAbleText(txDetails.Outputs[i].AccountName, false),
			copyAbleText(txDetails.Outputs[i].Address, true),
			copyAbleText(dcrutil.Amount(txDetails.Outputs[i].Amount).String(), false),
			copyAbleText(txDetails.Outputs[i].ScriptType, false),
		))
	}
	txOutput.NewTable(outputTableColumnLabels, outputBox...)

	txDetailsData := widget.NewVBox(
		widgets.NewHSpacer(10),
		tableData,
		widget.NewLabelWithStyle("Inputs", fyne.TextAlignLeading, fyne.TextStyle{Bold: true}),
		txInput.Result,
		widgets.NewVSpacer(10),
		widget.NewLabelWithStyle("Outputs", fyne.TextAlignLeading, fyne.TextStyle{Bold: true}),
		txOutput.Result,
		widgets.NewHSpacer(10),
	)

	txDetailsScrollContainer := widget.NewScrollContainer(txDetailsData)

	txDetailsOutput := widget.NewVBox(
		widgets.NewHSpacer(10),
		widget.NewHBox(
			txDetailslabel,
			widgets.NewHSpacer(150),
			messageLabel,
			layout.NewSpacer(),
			minimizeIcon,
			widgets.NewHSpacer(30),
		),
		widgets.NewHSpacer(10),
		fyne.NewContainerWithLayout(layout.NewFixedGridLayout(fyne.NewSize(txDetailsScrollContainer.MinSize().Width+10, txDetailsScrollContainer.MinSize().Height+400)), txDetailsScrollContainer),
		widgets.NewHSpacer(10),
	)

	txDetailsPopUp = widget.NewModalPopUp(widget.NewVBox(fyne.NewContainer(txDetailsOutput)),
		window.Canvas())
	txDetailsPopUp.Show()
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
				// selectedFilterName := strings.Split(filter, " ")[0]
				// selectedFilterId := allTxFilters[selectedFilterName]
				// if selectedFilterId != historyPageData.currentTxFilter {
				// 	go fetchAndDisplayTransactions(0, selectedFilterId)
				// }

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

func fetchAndDisplayTransactions(wallet *dcrlibwallet.LibWallet, txOffset int, filter int32) *widget.Box {
	tableHeading := widget.NewHBox(
		widget.NewLabelWithStyle("Date (UTC)", fyne.TextAlignCenter, fyne.TextStyle{Bold: true}),
		widget.NewLabelWithStyle("Type", fyne.TextAlignCenter, fyne.TextStyle{Bold: true}),
		widget.NewLabelWithStyle("Direction", fyne.TextAlignCenter, fyne.TextStyle{Bold: true}),
		widget.NewLabelWithStyle("Amount", fyne.TextAlignCenter, fyne.TextStyle{Bold: true}),
		widget.NewLabelWithStyle("Fee", fyne.TextAlignCenter, fyne.TextStyle{Bold: true}),
		widget.NewLabelWithStyle("Status", fyne.TextAlignCenter, fyne.TextStyle{Bold: true}),
		widget.NewLabelWithStyle("Hash", fyne.TextAlignCenter, fyne.TextStyle{Bold: true}))


	// txns, err := wallet.GetTransactionsRaw(int32(txOffset), txPerPage, dcrlibwallet.TxFilterAll)
	// if err != nil {
	// 	// displayMessage(err.Error(), MessageKindError)
	// 	// return
	// }

	// // calculate max number of digits after decimal point for all tx amounts
	// inputsAndOutputsAmount := make([]int64, len(txns))
	// for i, tx := range txns {
	// 	inputsAndOutputsAmount[i] = tx.Amount
	// }
	// maxDecimalPlacesForTxAmounts := helpers.MaxDecimalPlaces(inputsAndOutputsAmount)

	// var hBox []*widget.Box
	// for i, tx := range txns {
	// 	status := "Pending"
	// 	confirmations := wallet.GetBestBlock() - tx.BlockHeight + 1
	// 	if tx.BlockHeight != -1 && confirmations > dcrlibwallet.DefaultRequiredConfirmations {
	// 		status = "Confirmed"
	// 	}

	// 	formattedAmount := helpers.FormatAmountDisplay(tx.Amount, maxDecimalPlacesForTxAmounts)
	// 	trimmedHash := txns[i].Hash[:25] + "..."
	// 	hBox = append(hBox, widget.NewHBox(
	// 		widget.NewLabelWithStyle(fmt.Sprintf("%-10s", dcrlibwallet.ExtractDateOrTime(tx.Timestamp)), fyne.TextAlignCenter, fyne.TextStyle{}),
	// 		widget.NewLabelWithStyle(fmt.Sprintf("%-10s", dcrlibwallet.TransactionDirectionName(tx.Direction)), fyne.TextAlignCenter, fyne.TextStyle{}),
	// 		widget.NewLabelWithStyle(fmt.Sprintf("%12s", status), fyne.TextAlignLeading, fyne.TextStyle{}),
	// 		widget.NewLabelWithStyle(fmt.Sprintf("%15s", formattedAmount), fyne.TextAlignTrailing, fyne.TextStyle{}),
	// 		widget.NewLabelWithStyle(fmt.Sprintf("%-8s", tx.Type), fyne.TextAlignCenter, fyne.TextStyle{}),
	// 		widget.NewLabelWithStyle(trimmedHash, fyne.TextAlignLeading, fyne.TextStyle{}),
	// 	))
	// }

	h := widgets.NewTable(tableHeading, nil)
	// txTable.Refresh()

	return h
}

