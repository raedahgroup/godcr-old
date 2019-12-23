package historypagehandler

import (
	"sort"

	"fyne.io/fyne"
	"fyne.io/fyne/layout"
	"fyne.io/fyne/theme"
	"fyne.io/fyne/widget"

	"github.com/raedahgroup/dcrlibwallet"
	"github.com/raedahgroup/godcr/fyne/assets"
	"github.com/raedahgroup/godcr/fyne/pages/handler/values"
	"github.com/raedahgroup/godcr/fyne/widgets"
)

func (historyPage *HistoryPageData) txWalletList() {
	var txTable widgets.Table
	walletListWidget := widget.NewVBox()

	walletsID := historyPage.MultiWallet.OpenedWalletIDsRaw()
	if len(walletsID) == 0 {
		historyPage.errorMessage = values.WalletsErr
		return
	}
	sort.Ints(walletsID)

	// var selectedWalletLabel  *widget.Label
	selectedWalletLabel := widget.NewLabel(historyPage.MultiWallet.WalletWithID(walletsID[0]).Name)
	var txWalletSelectionPopup *widget.PopUp

	historyPage.txFilterDropDown(walletsID[0])
	historyPage.txSortDropDown()
	historyPage.txTableHeader(&historyPage.txTable)
	historyPage.fetchTx(&historyPage.txTable, 0, dcrlibwallet.TxFilterAll, false)

	for index, walletID := range walletsID {
		wallet := historyPage.MultiWallet.WalletWithID(walletID)
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
			selectedWalletLabel.SetText(wallet.Name)
			historyPage.selectedFilterId = dcrlibwallet.TxFilterAll
			historyPage.txFilterDropDown(individualWalletID)
			historyPage.txSortDropDown()
			historyPage.txTableHeader(&txTable)
			historyPage.fetchTx(&txTable, 0, historyPage.selectedFilterId, false)
			txWalletSelectionPopup.Hide()
		}))
	}

	// txWalletSelectionPopup create a popup that has tx wallet
	txWalletSelectionPopup = widget.NewPopUp(fyne.NewContainerWithLayout(
		layout.NewFixedGridLayout(fyne.NewSize(walletListWidget.MinSize().Width, 70)), widget.NewScrollContainer(walletListWidget)), historyPage.Window.Canvas())
	txWalletSelectionPopup.Hide()

	walletListTab := widget.NewHBox(
		selectedWalletLabel,
		widgets.NewHSpacer(10),
		widget.NewIcon(historyPage.icons[assets.CollapseIcon]),
	)

	// walletDropDown creates a popup like dropdown that holds the list of available wallets.
	var walletDropDown *widgets.ClickableBox
	walletDropDown = widgets.NewClickableBox(walletListTab, func() {
		txWalletSelectionPopup.Move(fyne.CurrentApp().Driver().AbsolutePositionForObject(
			walletDropDown).Add(fyne.NewPos(0, walletDropDown.Size().Height)))
		txWalletSelectionPopup.Show()
	})

	historyPage.HistoryPageContents.Append(widget.NewHBox(walletDropDown))
}
