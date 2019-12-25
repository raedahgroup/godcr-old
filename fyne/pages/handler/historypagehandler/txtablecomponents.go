package historypagehandler

import (
	"fmt"
	"time"

	"fyne.io/fyne"
	"fyne.io/fyne/widget"

	"github.com/decred/dcrd/dcrutil"
	"github.com/raedahgroup/dcrlibwallet"
	"github.com/raedahgroup/godcr/fyne/helpers"
	"github.com/raedahgroup/godcr/fyne/widgets"
)

const txPerPage int32 = 15

func (historyPage *HistoryPageData) txTableHeader(txTable *widgets.Table) {
	tableHeading := widget.NewHBox(
		widget.NewLabelWithStyle("#", fyne.TextAlignLeading, fyne.TextStyle{Bold: true}),
		widget.NewLabelWithStyle("Date (UTC)", fyne.TextAlignCenter, fyne.TextStyle{Bold: true}),
		widget.NewLabelWithStyle("Status", fyne.TextAlignCenter, fyne.TextStyle{Bold: true}),
		widget.NewLabelWithStyle("Amount", fyne.TextAlignCenter, fyne.TextStyle{Bold: true}),
		widget.NewLabelWithStyle("Fee", fyne.TextAlignCenter, fyne.TextStyle{Bold: true}),
		widget.NewLabelWithStyle("Hash", fyne.TextAlignCenter, fyne.TextStyle{Bold: true}),
	)

	var hBox []*widget.Box

	txTable.NewTable(tableHeading, hBox...)
	historyPage.txTable.Result.Children = txTable.Result.Children
	return
}

func (historyPage *HistoryPageData) fetchTx(txTable *widgets.Table, txOffset, filter int32, prepend bool) {
	if filter != historyPage.selectedFilterId {
		txOffset = 0
		historyPage.TotalTxFetched = 0
		historyPage.selectedFilterId = filter
	}

	txns, err := historyPage.MultiWallet.WalletWithID(historyPage.selectedWalletID).GetTransactionsRaw(txOffset, txPerPage, filter, historyPage.selectedtxSort)
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

	historyPage.TotalTxFetched += int32(len(txns))

	var txBox []*widget.Box
	for _, tx := range txns {
		status := "Pending"
		confirmations := historyPage.MultiWallet.WalletWithID(historyPage.selectedWalletID).GetBestBlock() - tx.BlockHeight + 1
		if tx.BlockHeight != -1 && confirmations > dcrlibwallet.DefaultRequiredConfirmations {
			status = "Confirmed"
		}

		trimmedHash := tx.Hash[:10] + "..." + tx.Hash[len(tx.Hash)-5:]
		txForTrimmedHash := tx.Hash
		txDirectionIcon := widget.NewIcon(historyPage.icons[helpers.TxDirectionIcon(tx.Direction)])
		txBox = append(txBox, widget.NewHBox(
			widget.NewHBox(txDirectionIcon, widget.NewLabel("")),
			widget.NewLabelWithStyle(dcrlibwallet.ExtractDateOrTime(tx.Timestamp), fyne.TextAlignCenter, fyne.TextStyle{}),
			widget.NewLabelWithStyle(status, fyne.TextAlignLeading, fyne.TextStyle{}),
			widget.NewLabelWithStyle(dcrutil.Amount(tx.Amount).String(), fyne.TextAlignTrailing, fyne.TextStyle{}),
			widget.NewLabelWithStyle(dcrutil.Amount(tx.Fee).String(), fyne.TextAlignCenter, fyne.TextStyle{}),
			widgets.NewClickableBox(widget.NewHBox(widget.NewLabelWithStyle(trimmedHash, fyne.TextAlignCenter, fyne.TextStyle{Italic: true})), func() {
				historyPage.fetchTxDetails(txForTrimmedHash)
			}),
		))
	}

	if prepend {
		txTable.Prepend(txBox...)
	} else {
		txTable.Append(txBox...)
	}

	historyPage.txTable.Result.Children = txTable.Result.Children
	widget.Refresh(historyPage.txTable.Result)
	widget.Refresh(historyPage.txTable.Container)
	historyPage.txTable.Container.Show()

	// wait four sec then update tx table
	time.AfterFunc(time.Second*2, func() {
		if historyPage.TabMenu.CurrentTabIndex() != 1 {
			return
		}
		historyPage.updateTable()
	})

	historyPage.errorLabel.Hide()
}

func (historyPage *HistoryPageData) updateTable() {
	size := historyPage.txTable.Container.Content.Size().Height - historyPage.txTable.Container.Size().Height
	scrollPosition := float64(historyPage.txTable.Container.Offset.Y) / float64(size)
	txTableRowCount := historyPage.txTable.NumberOfColumns()

	if historyPage.allTxCount > int(historyPage.TotalTxFetched) {
		if historyPage.txTable.Container.Offset.Y == 0 {
			// table not yet scrolled wait 4 secs and update
			time.AfterFunc(time.Second*2, func() {
				if historyPage.TabMenu.CurrentTabIndex() != 1 {
					return
				}
				historyPage.updateTable()
			})
		} else if scrollPosition < 0.5 {
			if historyPage.TotalTxFetched == txPerPage {
				time.AfterFunc(time.Second*2, func() {
					if historyPage.TabMenu.CurrentTabIndex() != 1 {
						return
					}
					historyPage.updateTable()
				})
			}
			if historyPage.TotalTxFetched >= 50 {
				historyPage.TotalTxFetched -= txPerPage * 2
				if txTableRowCount >= 50 {
					historyPage.txTable.Delete(txTableRowCount-int(txPerPage), txTableRowCount)
				}
				historyPage.fetchTx(&historyPage.txTable, historyPage.TotalTxFetched, historyPage.selectedFilterId, true)
			}
		} else if scrollPosition >= 0.5 {
			if txTableRowCount >= 50 {
				historyPage.txTable.Delete(0, txTableRowCount-int(txPerPage))
			}
			historyPage.fetchTx(&historyPage.txTable, historyPage.TotalTxFetched, historyPage.selectedFilterId, false)
		}
	}
}
