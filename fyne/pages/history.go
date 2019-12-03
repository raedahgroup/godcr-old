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

const txPerPage int32 = 8

type historyPageData struct {
	txTable          widgets.Table
	txDetailsTable   widgets.Table
	allTxCount       int
	txns             []*dcrlibwallet.Transaction
	selectedFilterId int32
	errorLabel       *widget.Label
	txOffset         int32
	TotalTxFetched   int32
}

var history historyPageData


func HistoryPageContent(wallet *dcrlibwallet.LibWallet, window fyne.Window, tabmenu *widget.TabContainer) fyne.CanvasObject {
	// error handler
	history.errorLabel = widget.NewLabelWithStyle("", fyne.TextAlignCenter, fyne.TextStyle{Bold: true})
	history.errorLabel.Hide()

	pageTitleLabel := widget.NewLabelWithStyle("Transactions", fyne.TextAlignLeading, fyne.TextStyle{Bold: true, Italic: true})

	txFilterDropDown := prepareTxFilterDropDown(wallet, window, tabmenu)

	txTableHeader(wallet, &history.txTable, window)
	fetchTx(&history.txTable, 0, dcrlibwallet.TxFilterAll, wallet, window, false, tabmenu)
	historyPageOutput := widget.NewVBox(
		widgets.NewVSpacer(5),
		widget.NewHBox(pageTitleLabel),
		widgets.NewVSpacer(5),
		txFilterDropDown,
		widgets.NewVSpacer(5),
		fyne.NewContainerWithLayout(layout.NewFixedGridLayout(fyne.NewSize(history.txTable.Container.MinSize().Width, history.txTable.Container.MinSize().Height+200)), history.txTable.Container),
		widgets.NewVSpacer(15),
		history.errorLabel,
	)

	return widget.NewHBox(widgets.NewHSpacer(18), historyPageOutput)
}

func prepareTxFilterDropDown(wallet *dcrlibwallet.LibWallet, window fyne.Window, tabmenu *widget.TabContainer) *widgets.ClickableBox {
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
		errorMessage := fmt.Sprintf("Cannot load history page. Error getting transaction count for filter All: %s", err.Error())
		errorHandler(errorMessage, history.errorLabel)
		return nil
	}
	history.allTxCount = txCountForFilter

	selectedTxFilterLabel := widget.NewLabel(fmt.Sprintf("%s (%d)", "All", txCountForFilter))

	var txFilterSelectionPopup *widget.PopUp
	txFilterListWidget := widget.NewVBox()
	for _, filterName := range allTxFilterNames {
		filterId := allTxFilters[filterName]
		txCountForFilter, err := wallet.CountTransactions(filterId)
		if err != nil {
			errorMessage := fmt.Sprintf("Cannot load history page. Error getting transaction count for filter %s: %s",
				filterName, err.Error())
			errorHandler(errorMessage, history.errorLabel)
			return nil
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
					history.allTxCount = allTxCountForSelectedTx
				}

				if selectedFilterId != history.selectedFilterId {
					txTableHeader(wallet, &txTable, window)
					history.txTable.Result.Children = txTable.Result.Children
					fetchTx(&txTable, 0, selectedFilterId, wallet, window, false, tabmenu)
					widget.Refresh(history.txTable.Result)
					selectedTxFilterLabel.SetText(filter)
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
		errorHandler(errorMessage, history.errorLabel)
		return nil
	}

	txFilterTab := widget.NewHBox(
		selectedTxFilterLabel,
		widgets.NewHSpacer(50),
		widget.NewIcon(icons[assets.CollapseIcon]),
	)

	var txFilterDropDown *widgets.ClickableBox
	txFilterDropDown = widgets.NewClickableBox(txFilterTab, func() {
		txFilterSelectionPopup.Move(fyne.CurrentApp().Driver().AbsolutePositionForObject(
			txFilterDropDown).Add(fyne.NewPos(0, txFilterDropDown.Size().Height)))
		txFilterSelectionPopup.Show()
	})

	return txFilterDropDown
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
	if filter != history.selectedFilterId {
		txOffset = 0
		history.TotalTxFetched = 0
		history.selectedFilterId = filter
	}

	txns, err := wallet.GetTransactionsRaw(txOffset, txPerPage, filter)
	if err != nil {
		errorHandler(fmt.Sprintf("Error getting transaction for Filter %s: %s", filter, err.Error()), history.errorLabel)
		return
	}
	if len(txns) == 0 {
		errorHandler("No transaction history yet.", history.errorLabel)
		txTable.Container.Hide()
		return
	}

	history.TotalTxFetched += int32(len(txns))

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
				fetchTxDetails(txForTrimmedHash, wallet, window, history.errorLabel, tabmenu)
			}),
		))
	}

	if prepend {
		txTable.Prepend(txBox...)
	} else {
		txTable.Append(txBox...)
	}

	history.txTable.Result.Children = txTable.Result.Children
	widget.Refresh(history.txTable.Result)
	widget.Refresh(history.txTable.Container)

	time.AfterFunc(time.Second*8, func() {
		updateTable(wallet, window, tabmenu)
	})

	history.errorLabel.Hide()
}

func fetchTxDetails(hash string, wallet *dcrlibwallet.LibWallet, window fyne.Window, errorLabel *widget.Label, tabmenu *widget.TabContainer) {
	messageLabel := widget.NewLabelWithStyle("Fetching data..", fyne.TextAlignCenter, fyne.TextStyle{Bold: true})
	if messageLabel.Hidden == false {
		time.AfterFunc(time.Millisecond*200, func() {
			if tabmenu.CurrentTabIndex() == 1 {
				messageLabel.Hide()
			}
		})
	}

	chainHash, err := chainhash.NewHashFromStr(hash)
	if err != nil {
		errorHandler(fmt.Sprintf("Error: %s", err.Error()), history.errorLabel)
		return
	}

	txDetails, err := wallet.GetTransactionRaw(chainHash[:])
	if err != nil {
		errorHandler(fmt.Sprintf("Error fetching transaction details for %s: %s", hash, err.Error()), history.errorLabel)
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

	txDetailslabel := widget.NewLabelWithStyle("Transaction Details", fyne.TextAlignLeading, fyne.TextStyle{Bold: true, Italic: true})
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

	var txDetailsPopUp *widget.PopUp
	txDetailsScrollContainer := widget.NewScrollContainer(txDetailsData)
	minimizeIcon := widgets.NewClickableIcon(theme.CancelIcon(), nil, func() { txDetailsPopUp.Hide() })

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

func updateTable(wallet *dcrlibwallet.LibWallet, window fyne.Window, tabmenu *widget.TabContainer) {
	size := history.txTable.Container.Content.Size().Height - history.txTable.Container.Size().Height
	scrollPosition := float64(history.txTable.Container.Offset.Y) / float64(size)
	fmt.Println(scrollPosition, size)

	if history.allTxCount > int(history.TotalTxFetched) {
		if history.txTable.Container.Offset.Y == 0 {
			fetchTx(&history.txTable, history.TotalTxFetched, history.selectedFilterId, wallet, window, false, tabmenu)
		}else if scrollPosition < 0.8 {
			time.AfterFunc(time.Second*8, func() {
				updateTable(wallet, window, tabmenu)
			})
		} else if scrollPosition >= 0.8 {
			fetchTx(&history.txTable, history.TotalTxFetched, history.selectedFilterId, wallet, window, false, tabmenu)
		} else if scrollPosition < 0.2 {
			if history.TotalTxFetched <= txPerPage {

			}
			fetchTx(&history.txTable, history.TotalTxFetched, history.selectedFilterId, wallet, window, true, tabmenu)
		}
	}
}
