package historypagehandler

import (
	"fmt"
	"time"

	"fyne.io/fyne"
	"fyne.io/fyne/widget"
	// "fyne.io/fyne/layout"

	"github.com/decred/dcrd/dcrutil"
	"github.com/raedahgroup/dcrlibwallet"
	"github.com/raedahgroup/godcr/fyne/helpers"
	"github.com/raedahgroup/godcr/fyne/widgets"
	// "github.com/raedahgroup/godcr/fyne/assets"
	// "github.com/raedahgroup/godcr/fyne/pages/handler/values"
)

const txPerPage int32 = 10
const maxTxPerPage int32 = 28

func (historyPage *HistoryPageData) txTableHeader(txTable *widgets.Table) {
	tableHeading := widget.NewHBox(
		widget.NewLabelWithStyle("#", fyne.TextAlignLeading, fyne.TextStyle{Bold: true}),
		widget.NewLabelWithStyle("Date (UTC)", fyne.TextAlignCenter, fyne.TextStyle{Bold: true}),
		widget.NewLabelWithStyle("Status", fyne.TextAlignCenter, fyne.TextStyle{Bold: true}),
		widget.NewLabelWithStyle("Amount", fyne.TextAlignCenter, fyne.TextStyle{Bold: true}),
		widget.NewLabelWithStyle("Fee", fyne.TextAlignCenter, fyne.TextStyle{Bold: true}),
		widget.NewLabelWithStyle("Hash", fyne.TextAlignCenter, fyne.TextStyle{Bold: true}),
	)
	fmt.Println(tableHeading)
	return
}

func (historyPage *HistoryPageData) fetchTx(txTable *widgets.Table, txOffset, filter int32, prepend bool) {
	if filter != historyPage.selectedFilterId {
		txOffset = 0
		historyPage.selectedFilterId = filter
	}

	pageTxOffset := (historyPage.pageCount - 1) * txPerPage
	txns, err := historyPage.MultiWallet.WalletWithID(historyPage.selectedWalletID).GetTransactionsRaw(pageTxOffset, txPerPage, filter, historyPage.selectedtxSort)
	if err != nil {
		helpers.ErrorHandler(fmt.Sprintf("Error getting transaction for Filter: %s", err.Error()), historyPage.errorLabel)
		historyPage.txTable.Container.Hide()
		return
	}
	if len(txns) == 0 {
		helpers.ErrorHandler(fmt.Sprintf("No transactions for %s yet.", historyPage.MultiWallet.WalletWithID(historyPage.selectedWalletID).Name), historyPage.errorLabel)
		historyPage.txTable.Container.Hide()
		return
	}

	tableHeading := widget.NewHBox(
		widget.NewLabelWithStyle("#", fyne.TextAlignLeading, fyne.TextStyle{Bold: true}),
		widget.NewLabelWithStyle("Date (UTC)", fyne.TextAlignCenter, fyne.TextStyle{Bold: true}),
		widget.NewLabelWithStyle("Status", fyne.TextAlignCenter, fyne.TextStyle{Bold: true}),
		widget.NewLabelWithStyle("Amount", fyne.TextAlignCenter, fyne.TextStyle{Bold: true}),
		widget.NewLabelWithStyle("Fee", fyne.TextAlignCenter, fyne.TextStyle{Bold: true}),
		widget.NewLabelWithStyle("Hash", fyne.TextAlignCenter, fyne.TextStyle{Bold: true}),
	)

	historyPage.lastTx = len(txns)
	if prepend {
		historyPage.TotalTxFetched -= int32(historyPage.lastTx)
	} else {
		historyPage.TotalTxFetched += int32(historyPage.lastTx)
	}

	fmt.Println(historyPage.TotalTxFetched, historyPage.lastTx, historyPage.pageCount)

	var txBox []*widget.Box
	for _, tx := range txns {
		status := "Pending"
		confirmations := historyPage.MultiWallet.WalletWithID(historyPage.selectedWalletID).GetBestBlock() - tx.BlockHeight + 1
		if tx.BlockHeight != -1 && confirmations > dcrlibwallet.DefaultRequiredConfirmations {
			status = "Confirmed"
		}

		trimmedHash := tx.Hash[:7] + "..." + tx.Hash[len(tx.Hash)-7:]
		txForTrimmedHash := tx.Hash
		txDirectionIcon := widget.NewIcon(historyPage.icons[helpers.TxDirectionIcon(tx.Direction)])
		txBox = append(txBox, widget.NewHBox(
			widget.NewHBox(txDirectionIcon, widget.NewLabel("")),
			widget.NewLabelWithStyle(dcrlibwallet.ExtractDateOrTime(tx.Timestamp), fyne.TextAlignCenter, fyne.TextStyle{}),
			widget.NewLabelWithStyle(status, fyne.TextAlignCenter, fyne.TextStyle{}),
			widget.NewLabelWithStyle(dcrutil.Amount(tx.Amount).String(), fyne.TextAlignTrailing, fyne.TextStyle{}),
			widget.NewLabelWithStyle(dcrutil.Amount(tx.Fee).String(), fyne.TextAlignCenter, fyne.TextStyle{}),
			widgets.NewClickableBox(widget.NewHBox(widget.NewLabelWithStyle(trimmedHash, fyne.TextAlignCenter, fyne.TextStyle{Italic: true})), func() {
				historyPage.fetchTxDetails(txForTrimmedHash)
			}),
		))
	}

	txTable.NewTable(tableHeading, txBox...)
	historyPage.txTable.Result.Children = txTable.Result.Children
	historyPage.txTable.Container.Offset.Y = 2
	widget.Refresh(historyPage.txTable.Result)
	widget.Refresh(historyPage.txTable.Container)
	historyPage.txTable.Container.Show()

	historyPage.updateTable()

	historyPage.errorLabel.Hide()
}

func (historyPage *HistoryPageData) updateTable() {
	var txTable widgets.Table
	size := historyPage.txTable.Container.Content.Size().Height - historyPage.txTable.Container.Size().Height

	if historyPage.allTxCount > int(historyPage.TotalTxFetched) {
		if historyPage.txTable.Container.Offset.Y == size {
			historyPage.pageCount += 1
			historyPage.fetchTx(&txTable, historyPage.TotalTxFetched, historyPage.selectedFilterId, false)
		} else if historyPage.txTable.Container.Offset.Y == 0 && historyPage.pageCount > 1 {
			historyPage.pageCount -= 1
			historyPage.fetchTx(&txTable, historyPage.TotalTxFetched, historyPage.selectedFilterId, true)
		} else {
			time.AfterFunc(time.Millisecond*500, func() {
				if historyPage.TabMenu.CurrentTabIndex() != 1 {
					return
				}
				historyPage.updateTable()
			})
		}
	} else if historyPage.allTxCount == int(historyPage.TotalTxFetched) && historyPage.pageCount > 1 {
		if historyPage.lastTx < 10 {
			historyPage.backIcon.Show()
		}
	}
}
