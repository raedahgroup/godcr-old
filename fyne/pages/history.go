package pages

import (
	"fmt"
	"image/color"
	"net/url"
	"sort"
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
	"github.com/raedahgroup/godcr/fyne/helpers"
	"github.com/raedahgroup/godcr/fyne/widgets"
)

const txPerPage int32 = 25
const netType = "testnet"

type txHistoryPageData struct {
	allTxCount                 int
	selectedWalletID           int
	selectedFilterId           int32
	TotalTxFetched             int32
	selectedtxSort             bool
	errorMessage               string
	walletListTab              *widget.Box
	txFilterTab                *widget.Box
	txSortFilterTab            *widget.Box
	txTable                    widgets.Table
	errorLabel                 *widget.Label
	selectedTxFilterLabel      *widget.Label
	selectedTxSortFilterLabel  *widget.Label
	selectedWalletLabel        *widget.Label
	txFilterSelectionPopup     *widget.PopUp
	txSortFilterSelectionPopup *widget.PopUp
	txWalletSelectionPopup     *widget.PopUp
	icons                      map[string]*fyne.StaticResource
}

var txHistory txHistoryPageData

func historyPageContent(app *AppInterface) fyne.CanvasObject {
	txHistory.errorLabel = widget.NewLabelWithStyle("", fyne.TextAlignCenter, fyne.TextStyle{Bold: true})
	txHistory.errorLabel.Hide()

	txHistory.selectedFilterId = dcrlibwallet.TxFilterAll
	txHistory.selectedWalletLabel = widget.NewLabel("")
	txHistory.selectedTxFilterLabel = widget.NewLabel("")
	txHistory.selectedTxSortFilterLabel = widget.NewLabel("")

	// gets all icons used on this page
	icons, err := assets.GetIcons(assets.CollapseIcon, assets.InfoIcon, assets.SendIcon, assets.ReceiveIcon, assets.ReceiveIcon, assets.InfoIcon, assets.RedirectIcon)
	if err != nil {
		errorMessage := fmt.Sprintf("Error: %s", err.Error())
		helpers.ErrorHandler(errorMessage, txHistory.errorLabel)
		return widget.NewHBox(widgets.NewHSpacer(18), txHistory.errorLabel)
	}
	txHistory.icons = icons

	// history page title label
	pageTitleLabel := widget.NewLabelWithStyle("Transactions", fyne.TextAlignLeading, fyne.TextStyle{Bold: true, Italic: true})

	// infoPopUp creates a popup with history page hint-text
	var infoIcon *widgets.ImageButton
	info := "- Tap Hash to view Transaction details.\n\n- Tap Blue Text to Copy."
	infoIcon = widgets.NewImageButton(txHistory.icons[assets.InfoIcon], nil, func() {
		infoLabel := widget.NewLabelWithStyle(info, fyne.TextAlignLeading, fyne.TextStyle{Monospace: true})
		gotItLabel := canvas.NewText("Got it", color.RGBA{41, 112, 255, 255})
		gotItLabel.TextStyle = fyne.TextStyle{Bold: true}
		gotItLabel.TextSize = 14

		var infoPopUp *widget.PopUp
		infoPopUp = widget.NewPopUp(widget.NewVBox(
			widgets.NewVSpacer(5),
			widget.NewHBox(widgets.NewHSpacer(5), infoLabel, widgets.NewHSpacer(5)),
			widgets.NewVSpacer(5),
			widget.NewHBox(layout.NewSpacer(), widgets.NewClickableBox(widget.NewVBox(gotItLabel), func() { infoPopUp.Hide() }), widgets.NewHSpacer(5)),
			widgets.NewVSpacer(5),
		), app.Window.Canvas())

		infoPopUp.Move(fyne.CurrentApp().Driver().AbsolutePositionForObject(infoIcon).Add(fyne.NewPos(0, infoIcon.Size().Height)))
	})

	// history output widget
	txHistoryPageOutput := widget.NewVBox(
		widget.NewHBox(pageTitleLabel, widgets.NewHSpacer(110), infoIcon),
		widgets.NewVSpacer(5),
	)

	// initialize history page data
	txWalletList(app.MultiWallet, app.Window, app.tabMenu)

	// walletDropDown creates a popup like dropdown that holds the list of available wallets.
	var walletDropDown *widgets.ClickableBox
	walletDropDown = widgets.NewClickableBox(txHistory.walletListTab, func() {
		txHistory.txWalletSelectionPopup.Move(fyne.CurrentApp().Driver().AbsolutePositionForObject(
			walletDropDown).Add(fyne.NewPos(0, walletDropDown.Size().Height)))
		txHistory.txWalletSelectionPopup.Show()
	})

	// txFilterDropDown creates a popup like dropdown that holds the list of tx filters.
	var txFilterDropDown *widgets.ClickableBox
	txFilterDropDown = widgets.NewClickableBox(txHistory.txFilterTab, func() {
		if txHistory.allTxCount == 0 {
			txHistory.txFilterSelectionPopup.Hide()
		} else {
			txHistory.txFilterSelectionPopup.Move(fyne.CurrentApp().Driver().AbsolutePositionForObject(
				txFilterDropDown).Add(fyne.NewPos(0, txFilterDropDown.Size().Height)))
			txHistory.txFilterSelectionPopup.Show()
		}
	})

	// txSortFilterDropDown creates a popup like dropdown that holds the list of sort filters.
	var txSortFilterDropDown *widgets.ClickableBox
	txSortFilterDropDown = widgets.NewClickableBox(txHistory.txSortFilterTab, func() {
		if txHistory.allTxCount == 0 {
			txHistory.txSortFilterSelectionPopup.Hide()
		} else {
			txHistory.txSortFilterSelectionPopup.Move(fyne.CurrentApp().Driver().AbsolutePositionForObject(
				txSortFilterDropDown).Add(fyne.NewPos(0, txSortFilterDropDown.Size().Height)))
			txHistory.txSortFilterSelectionPopup.Show()
		}
	})

	// catch all errors when trying to setup and render tx page data.
	if txHistory.errorMessage != "" {
		helpers.ErrorHandler(txHistory.errorMessage, txHistory.errorLabel)
		txHistoryPageOutput.Append(txHistory.errorLabel)
		return widget.NewHBox(widgets.NewHSpacer(18), txHistoryPageOutput)
	}

	txHistoryPageOutput.Append(walletDropDown)
	txHistoryPageOutput.Append(widget.NewHBox(txSortFilterDropDown, widgets.NewHSpacer(30), txFilterDropDown))
	txHistoryPageOutput.Append(widgets.NewVSpacer(5))
	txHistoryPageOutput.Append(txHistory.errorLabel)
	txHistoryPageOutput.Append(fyne.NewContainerWithLayout(layout.NewFixedGridLayout(fyne.NewSize(txHistory.txTable.Container.MinSize().Width, txHistory.txTable.Container.MinSize().Height+450)), txHistory.txTable.Container))
	txHistoryPageOutput.Append(widgets.NewVSpacer(15))

	return widget.NewHBox(widgets.NewHSpacer(18), txHistoryPageOutput)
}

func txWalletList(multiWallet *dcrlibwallet.MultiWallet, window fyne.Window, tabMenu *widget.TabContainer) {
	var txTable widgets.Table
	walletListWidget := widget.NewVBox()

	walletsID := multiWallet.OpenedWalletIDsRaw()
	if len(walletsID) == 0 {
		txHistory.errorMessage = "Could not retrieve wallets"
		return
	}
	sort.Ints(walletsID)

	txHistory.selectedWalletLabel.SetText(multiWallet.WalletWithID(walletsID[0]).Name)

	txFilterDropDown(multiWallet, window, tabMenu, walletsID[0])
	txSortDropDown(multiWallet, window, tabMenu)
	txTableHeader(multiWallet, &txHistory.txTable, window)
	fetchTx(&txHistory.txTable, 0, dcrlibwallet.TxFilterAll, multiWallet, window, tabMenu, false)

	for index, walletID := range walletsID {
		wallet := multiWallet.WalletWithID(walletID)
		if wallet == nil {
			continue
		}

		checkmarkIcon := widget.NewIcon(theme.ConfirmIcon())
		if index != 0 || walletID != walletsID[0] {
			checkmarkIcon.Hide()
		}

		walletContainer := widget.NewHBox(
			widget.NewLabel(wallet.Name),
			checkmarkIcon,
			widgets.NewHSpacer(5),
		)

		individualWalletID := walletID

		walletListWidget.Append(widgets.NewClickableBox(walletContainer, func() {
			// hide checkmark icon of other wallets
			for _, children := range walletListWidget.Children {
				if box, ok := children.(*widgets.ClickableBox); !ok {
					continue
				} else {
					if len(box.Children) != 3 {
						continue
					}

					if icon, ok := box.Children[1].(*widget.Icon); !ok {
						continue
					} else {
						icon.Hide()
					}
				}
			}

			checkmarkIcon.Show()
			txHistory.selectedWalletLabel.SetText(wallet.Name)
			txHistory.selectedFilterId = dcrlibwallet.TxFilterAll
			txFilterDropDown(multiWallet, window, tabMenu, individualWalletID)
			txSortDropDown(multiWallet, window, tabMenu)
			txTableHeader(multiWallet, &txTable, window)
			fetchTx(&txTable, 0, txHistory.selectedFilterId, multiWallet, window, tabMenu, false)
			txHistory.txWalletSelectionPopup.Hide()
		}))
	}

	// txWalletSelectionPopup create a popup that has tx wallet
	txHistory.txWalletSelectionPopup = widget.NewPopUp(widget.NewVBox(walletListWidget), window.Canvas())
	txHistory.txWalletSelectionPopup.Hide()

	txHistory.walletListTab = widget.NewHBox(
		txHistory.selectedWalletLabel,
		widgets.NewHSpacer(10),
		widget.NewIcon(txHistory.icons[assets.CollapseIcon]),
	)
}

func txFilterDropDown(multiWallet *dcrlibwallet.MultiWallet, window fyne.Window, tabMenu *widget.TabContainer, walletId int) {
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

	if walletId != txHistory.selectedWalletID {
		txHistory.selectedWalletID = walletId
	}

	txCountForFilter, err := multiWallet.WalletWithID(txHistory.selectedWalletID).CountTransactions(allTxFilters["All"])
	if err != nil {
		txHistory.errorMessage = fmt.Sprintf("Cannot load txHistory page. Error getting transaction count for filter All: %s", err.Error())
		return
	}

	txHistory.allTxCount = txCountForFilter
	txHistory.selectedTxFilterLabel.SetText(fmt.Sprintf("%s (%d)", "All", txCountForFilter))

	for _, filterName := range allTxFilterNames {
		filterId := allTxFilters[filterName]
		txCountForFilter, err := multiWallet.WalletWithID(txHistory.selectedWalletID).CountTransactions(filterId)
		if err != nil {
			txHistory.errorMessage = fmt.Sprintf("Cannot load txHistory page. Error getting transaction count for filter %s: %s", filterName, err.Error())
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
				if allTxCountForSelectedTx, err := multiWallet.WalletWithID(txHistory.selectedWalletID).CountTransactions(selectedFilterId); err == nil {
					txHistory.allTxCount = allTxCountForSelectedTx
				}

				if selectedFilterId != txHistory.selectedFilterId {
					txHistory.selectedTxFilterLabel.SetText(filter)
					txTableHeader(multiWallet, &txTable, window)
					fetchTx(&txTable, 0, selectedFilterId, multiWallet, window, tabMenu, false)
					widget.Refresh(txHistory.txTable.Result)
				}

				txHistory.txFilterSelectionPopup.Hide()
			}))
		}
	}

	// txFilterSelectionPopup create a popup that has tx filter name and tx count
	txHistory.txFilterSelectionPopup = widget.NewPopUp(widget.NewVBox(txFilterListWidget), window.Canvas())
	txHistory.txFilterSelectionPopup.Hide()

	txHistory.txFilterTab = widget.NewHBox(
		txHistory.selectedTxFilterLabel,
		widgets.NewHSpacer(10),
		widget.NewIcon(txHistory.icons[assets.CollapseIcon]),
		widgets.NewHSpacer(10),
	)
	widget.Refresh(txHistory.txFilterTab)
}

func txSortDropDown(multiWallet *dcrlibwallet.MultiWallet, window fyne.Window, tabMenu *widget.TabContainer) {
	var txTable widgets.Table
	var allTxSortNames = []string{"Newest", "Oldest"}
	var allTxSortFilters = map[string]bool{
		"Newest": true,
		"Oldest": false,
	}

	txHistory.selectedTxSortFilterLabel.SetText("Newest")
	txHistory.selectedtxSort = allTxSortFilters["Newest"]

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
			txHistory.selectedTxSortFilterLabel.SetText(newSortName)
			txHistory.selectedtxSort = txSort

			txTableHeader(multiWallet, &txTable, window)
			txHistory.txTable.Result.Children = txTable.Result.Children
			fetchTx(&txTable, 0, txHistory.selectedFilterId, multiWallet, window, tabMenu, false)
			widget.Refresh(txHistory.txTable.Result)
			txHistory.txSortFilterSelectionPopup.Hide()
		}))
	}

	// txSortFilterSelectionPopup create a popup that has tx filter name and tx count
	txHistory.txSortFilterSelectionPopup = widget.NewPopUp(widget.NewVBox(txSortFilterListWidget), window.Canvas())
	txHistory.txSortFilterSelectionPopup.Hide()

	txHistory.txSortFilterTab = widget.NewHBox(
		txHistory.selectedTxSortFilterLabel,
		widgets.NewHSpacer(10),
		widget.NewIcon(txHistory.icons[assets.CollapseIcon]),
	)
	widget.Refresh(txHistory.txSortFilterTab)
}

func txTableHeader(wallet *dcrlibwallet.MultiWallet, txTable *widgets.Table, window fyne.Window) {
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
	txHistory.txTable.Result.Children = txTable.Result.Children
	return
}

func fetchTx(txTable *widgets.Table, txOffset, filter int32, multiWallet *dcrlibwallet.MultiWallet, window fyne.Window, tabMenu *widget.TabContainer, prepend bool) {
	if filter != txHistory.selectedFilterId {
		txOffset = 0
		txHistory.TotalTxFetched = 0
		txHistory.selectedFilterId = filter
	}

	txns, err := multiWallet.WalletWithID(txHistory.selectedWalletID).GetTransactionsRaw(txOffset, txPerPage, filter, txHistory.selectedtxSort)
	if err != nil {
		helpers.ErrorHandler(fmt.Sprintf("Error getting transaction for Filter: %s", err.Error()), txHistory.errorLabel)
		txHistory.txTable.Container.Hide()
		return
	}
	if len(txns) == 0 {
		helpers.ErrorHandler(fmt.Sprintf("No transactions for %s yet.", multiWallet.WalletWithID(txHistory.selectedWalletID).Name), txHistory.errorLabel)
		txHistory.txTable.Container.Hide()
		return
	}

	txHistory.TotalTxFetched += int32(len(txns))

	var txBox []*widget.Box
	for _, tx := range txns {
		status := "Pending"
		confirmations := multiWallet.WalletWithID(txHistory.selectedWalletID).GetBestBlock() - tx.BlockHeight + 1
		if tx.BlockHeight != -1 && confirmations > dcrlibwallet.DefaultRequiredConfirmations {
			status = "Confirmed"
		}

		trimmedHash := tx.Hash[:15] + "..." + tx.Hash[len(tx.Hash)-15:]
		txForTrimmedHash := tx.Hash
		txDirectionLabel := widget.NewLabelWithStyle(dcrlibwallet.TransactionDirectionName(tx.Direction), fyne.TextAlignCenter, fyne.TextStyle{})
		txDirectionIcon := widget.NewIcon(txHistory.icons[helpers.TxDirectionIcon(tx.Direction)])
		txBox = append(txBox, widget.NewHBox(
			widget.NewLabelWithStyle(dcrlibwallet.ExtractDateOrTime(tx.Timestamp), fyne.TextAlignCenter, fyne.TextStyle{}),
			widget.NewHBox(txDirectionIcon, txDirectionLabel),
			widget.NewLabelWithStyle(status, fyne.TextAlignLeading, fyne.TextStyle{}),
			widget.NewLabelWithStyle(dcrutil.Amount(tx.Amount).String(), fyne.TextAlignTrailing, fyne.TextStyle{}),
			widget.NewLabelWithStyle(dcrutil.Amount(tx.Fee).String(), fyne.TextAlignCenter, fyne.TextStyle{}),
			widget.NewLabelWithStyle(tx.Type, fyne.TextAlignCenter, fyne.TextStyle{}),
			widgets.NewClickableBox(widget.NewHBox(widget.NewLabelWithStyle(trimmedHash, fyne.TextAlignCenter, fyne.TextStyle{Italic: true})), func() {
				fetchTxDetails(txForTrimmedHash, multiWallet, window, txHistory.errorLabel, tabMenu)
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
	txHistory.txTable.Container.Show()

	// wait four sec then update tx table
	time.AfterFunc(time.Second*4, func() {
		updateTable(multiWallet, window, tabMenu)
	})

	txHistory.errorLabel.Hide()
}

func updateTable(multiWallet *dcrlibwallet.MultiWallet, window fyne.Window, tabMenu *widget.TabContainer) {
	size := txHistory.txTable.Container.Content.Size().Height - txHistory.txTable.Container.Size().Height
	scrollPosition := float64(txHistory.txTable.Container.Offset.Y) / float64(size)
	txTableRowCount := txHistory.txTable.NumberOfColumns()

	if txHistory.allTxCount > int(txHistory.TotalTxFetched) {
		if txHistory.txTable.Container.Offset.Y == 0 {
			if txTableRowCount < int(txPerPage) {
				return
			}
			// table not yet scrolled wait 4 secs and update
			time.AfterFunc(time.Second*4, func() {
				updateTable(multiWallet, window, tabMenu)
			})
		} else if scrollPosition < 0.5 {
			if txHistory.TotalTxFetched == txPerPage {
				time.AfterFunc(time.Second*4, func() {
					updateTable(multiWallet, window, tabMenu)
				})
			}
			if txHistory.TotalTxFetched >= 50 {
				txHistory.TotalTxFetched -= txPerPage * 2
				if txTableRowCount >= 50 {
					txHistory.txTable.Delete(txTableRowCount-int(txPerPage), txTableRowCount)
				}
				fetchTx(&txHistory.txTable, txHistory.TotalTxFetched, txHistory.selectedFilterId, multiWallet, window, tabMenu, true)
			}
		} else if scrollPosition >= 0.5 {
			if txTableRowCount >= 50 {
				txHistory.txTable.Delete(0, txTableRowCount-int(txPerPage))
			}
			fetchTx(&txHistory.txTable, txHistory.TotalTxFetched, txHistory.selectedFilterId, multiWallet, window, tabMenu, false)
		}
	}
}

func fetchTxDetails(hash string, multiWallet *dcrlibwallet.MultiWallet, window fyne.Window, errorLabel *widget.Label, tabMenu *widget.TabContainer) {
	messageLabel := widget.NewLabelWithStyle("Fetching data..", fyne.TextAlignCenter, fyne.TextStyle{})
	time.AfterFunc(time.Millisecond*300, func() {
		if tabMenu.CurrentTabIndex() == 1 {
			messageLabel.SetText("")
		}
	})

	txDetailslabel := widget.NewLabelWithStyle("Transaction Details", fyne.TextAlignLeading, fyne.TextStyle{Bold: true, Italic: true})

	var txDetailsPopUp *widget.PopUp
	minimizeIcon := widgets.NewImageButton(theme.CancelIcon(), nil, func() { txDetailsPopUp.Hide() })
	errorMessageLabel := widget.NewLabelWithStyle("", fyne.TextAlignCenter, fyne.TextStyle{Bold: true})

	txDetailsErrorMethod := func() {
		txErrorDetailsOutput := widget.NewVBox(
			widgets.NewHSpacer(10),
			widget.NewHBox(
				txDetailslabel,
				widgets.NewHSpacer(txDetailslabel.MinSize().Width+180),
				minimizeIcon,
			),
			errorMessageLabel,
		)
		txDetailsPopUp = widget.NewModalPopUp(widget.NewVBox(fyne.NewContainer(txErrorDetailsOutput)),
			window.Canvas())
		txDetailsPopUp.Show()
	}

	chainHash, err := chainhash.NewHashFromStr(hash)
	if err != nil {
		helpers.ErrorHandler(fmt.Sprintf("fetching generating chainhash from for \n %s \n %s ", hash, err.Error()), errorMessageLabel)
		txDetailsErrorMethod()
		return
	}

	txDetails, err := multiWallet.WalletWithID(txHistory.selectedWalletID).GetTransactionRaw(chainHash[:])
	if err != nil {
		helpers.ErrorHandler(fmt.Sprintf("Error fetching transaction details for \n %s \n %s ", hash, err.Error()), errorMessageLabel)
		txDetailsErrorMethod()
		return
	}

	var confirmations int32 = 0
	if txDetails.BlockHeight != -1 {
		confirmations = multiWallet.WalletWithID(txHistory.selectedWalletID).GetBestBlock() - txDetails.BlockHeight + 1
	}

	var status string
	var spendUnconfirmed = multiWallet.ReadBoolConfigValueForKey(dcrlibwallet.SpendUnconfirmedConfigKey, true)
	if spendUnconfirmed || confirmations > dcrlibwallet.DefaultRequiredConfirmations {
		status = "Confirmed"
	} else {
		status = "Pending"
	}

	copyAbleText := func(text string, copyAble bool) *widgets.ClickableBox {
		var textToCopy *canvas.Text
		if copyAble {
			textToCopy = canvas.NewText(text, color.RGBA{0x44, 0x8a, 0xff, 0xff})
		} else {
			textToCopy = canvas.NewText(text, color.RGBA{0x00, 0x00, 0x00, 0xff})
		}
		textToCopy.TextSize = 14
		textToCopy.Alignment = fyne.TextAlignTrailing

		return widgets.NewClickableBox(widget.NewHBox(textToCopy),
			func() {
				messageLabel.SetText("Data Copied")
				clipboard := window.Clipboard()
				clipboard.SetContent(text)

				time.AfterFunc(time.Second*2, func() {
					if tabMenu.CurrentTabIndex() == 1 {
						messageLabel.SetText("")
					}
				})
			},
		)
	}

	tableConfirmations := widget.NewHBox(
		widget.NewLabelWithStyle("Confirmations:", fyne.TextAlignTrailing, fyne.TextStyle{Bold: true}),
		widget.NewLabelWithStyle(strconv.Itoa(int(confirmations)), fyne.TextAlignCenter, fyne.TextStyle{}),
	)

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

	link, err := url.Parse(fmt.Sprintf("https://%s.dcrdata.org/tx/%s", netType, txDetails.Hash))
	if err != nil {
		helpers.ErrorHandler(fmt.Sprintf("Error: ", err.Error()), errorMessageLabel)
		txDetailsErrorMethod()
		return
	}
	redirectWidget := widget.NewHBox(
		widget.NewHyperlinkWithStyle("View on dcrdata", link, fyne.TextAlignLeading, fyne.TextStyle{Bold: true}),
		widgets.NewHSpacer(20),
		widget.NewIcon(txHistory.icons[assets.RedirectIcon]),
	)

	txDetailsData := widget.NewVBox(
		widgets.NewHSpacer(10),
		tableData,
		widgets.NewHSpacer(15),
		redirectWidget,
		widgets.NewHSpacer(15),
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
			widgets.NewHSpacer(txDetailsScrollContainer.MinSize().Width-txDetailslabel.MinSize().Width-30),
			minimizeIcon,
		),
		widget.NewHBox(widgets.NewHSpacer(txDetailsScrollContainer.MinSize().Width/2-30), messageLabel),
		fyne.NewContainerWithLayout(layout.NewFixedGridLayout(fyne.NewSize(txDetailsScrollContainer.MinSize().Width+10, txDetailsScrollContainer.MinSize().Height+400)), txDetailsScrollContainer),
		widgets.NewVSpacer(10),
	)

	txDetailsPopUp = widget.NewModalPopUp(widget.NewVBox(fyne.NewContainer(txDetailsOutput)), window.Canvas())
}
