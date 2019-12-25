package historypagehandler

import (
	"errors"

	"fyne.io/fyne"
	"fyne.io/fyne/layout"
	"fyne.io/fyne/widget"

	"github.com/raedahgroup/dcrlibwallet"
	"github.com/raedahgroup/godcr/fyne/pages/handler/values"
	"github.com/raedahgroup/godcr/fyne/widgets"
)

type HistoryPageData struct {
	allTxCount                 int
	selectedWalletID           int
	selectedFilterId           int32
	TotalTxFetched             int32
	selectedtxSort             bool
	errorMessage               string
	txFilterTab                *widget.Box
	txSortFilterTab            *widget.Box
	txTable                    widgets.Table
	errorLabel                 *widget.Label
	selectedTxFilterLabel      *widget.Label
	selectedTxSortFilterLabel  *widget.Label
	txFilterSelectionPopup     *widget.PopUp
	txSortFilterSelectionPopup *widget.PopUp
	MultiWallet                *dcrlibwallet.MultiWallet
	icons                      map[string]*fyne.StaticResource
	HistoryPageContents        *widget.Box
	Window                     fyne.Window
	TabMenu                    *widget.TabContainer
}

func (historyPage *HistoryPageData) InitHistoryPage() error {
	historyPage.HistoryPageContents.Append(widgets.NewVSpacer(values.Padding))

	err := historyPage.initBaseObjects()
	if err != nil {
		return err
	}

	historyPage.HistoryPageContents.Append(widgets.NewVSpacer(values.SpacerSize10))

	historyPage.errorLabel = widget.NewLabelWithStyle("", fyne.TextAlignCenter, fyne.TextStyle{Bold: true})
	historyPage.errorLabel.Hide()

	historyPage.selectedFilterId = dcrlibwallet.TxFilterAll
	historyPage.selectedTxFilterLabel = widget.NewLabel("")
	historyPage.selectedTxSortFilterLabel = widget.NewLabel("")

	historyPage.txWalletList()

	// txFilterDropDown creates a popup like dropdown that holds the list of tx filters.
	var txFilterDropDown *widgets.ClickableBox
	txFilterDropDown = widgets.NewClickableBox(historyPage.txFilterTab, func() {
		if historyPage.allTxCount == 0 {
			historyPage.txFilterSelectionPopup.Hide()
		} else {
			historyPage.txFilterSelectionPopup.Move(fyne.CurrentApp().Driver().AbsolutePositionForObject(
				txFilterDropDown).Add(fyne.NewPos(0, txFilterDropDown.Size().Height)))
			historyPage.txFilterSelectionPopup.Show()
		}
	})

	// txSortFilterDropDown creates a popup like dropdown that holds the list of sort filters.
	var txSortFilterDropDown *widgets.ClickableBox
	txSortFilterDropDown = widgets.NewClickableBox(historyPage.txSortFilterTab, func() {
		if historyPage.allTxCount == 0 {
			historyPage.txSortFilterSelectionPopup.Hide()
		} else {
			historyPage.txSortFilterSelectionPopup.Move(fyne.CurrentApp().Driver().AbsolutePositionForObject(
				txSortFilterDropDown).Add(fyne.NewPos(0, txSortFilterDropDown.Size().Height)))
			historyPage.txSortFilterSelectionPopup.Show()
		}
	})

	// catch all errors when trying to setup and render tx page data.
	if historyPage.errorMessage != "" {
		return errors.New(historyPage.errorMessage)
	}

	historyPage.HistoryPageContents.Append(widget.NewHBox(txSortFilterDropDown, widgets.NewHSpacer(30), txFilterDropDown))
	historyPage.HistoryPageContents.Append(widgets.NewVSpacer(5))
	historyPage.HistoryPageContents.Append(historyPage.errorLabel)
	historyPage.HistoryPageContents.Append(fyne.NewContainerWithLayout(layout.NewFixedGridLayout(fyne.NewSize(historyPage.txTable.Result.MinSize().Width, historyPage.txTable.Container.MinSize().Height+450)), historyPage.txTable.Container))
	historyPage.HistoryPageContents.Append(widgets.NewVSpacer(15))
	return nil
}
