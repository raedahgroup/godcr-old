package pages

import (
	"sort"

	"fyne.io/fyne"
	"fyne.io/fyne/canvas"
	"fyne.io/fyne/widget"

	"github.com/raedahgroup/dcrlibwallet"

	"github.com/raedahgroup/godcr/fyne/pages/handler/values"
	"github.com/raedahgroup/godcr/fyne/pages/handler/walletshandler"
	"github.com/raedahgroup/godcr/fyne/widgets"
)

type walletPageDynamicData struct {
	walletTotalAmountLabel []*canvas.Text
	walletPageContents     *widget.Box
}

var walletPage walletPageDynamicData

func walletPageContent(tabMenu *widget.TabContainer, multiWallet *dcrlibwallet.MultiWallet, window fyne.Window) fyne.CanvasObject {
	openedWalletIDs := multiWallet.OpenedWalletIDsRaw()
	if len(openedWalletIDs) == 0 {
		return widget.NewHBox(widgets.NewHSpacer(10), widget.NewLabelWithStyle(values.WalletsErr, fyne.TextAlignCenter, fyne.TextStyle{Bold: true}))
	}
	sort.Ints(openedWalletIDs)

	initWalletDynamicContent(openedWalletIDs)
	initWalletPage := walletshandler.WalletPageObject{
		WalletTotalAmountLabel: walletPage.walletTotalAmountLabel,
		WalletPageContents:     walletPage.walletPageContents,

		OpenedWallets: openedWalletIDs,
		MultiWallet:   multiWallet,
		Window:        window,
		TabMenu:       tabMenu,
	}

	err := initWalletPage.InitWalletPage()
	if err != nil {
		return widget.NewHBox(widgets.NewHSpacer(10), widget.NewLabelWithStyle(values.WalletPageLoadErr, fyne.TextAlignCenter, fyne.TextStyle{Bold: true}))
	}

	return widget.NewHBox(widgets.NewHSpacer(values.Padding), initWalletPage.WalletPageContents)
}

func initWalletDynamicContent(openedWalletIDs []int) {
	walletPage = walletPageDynamicData{}
	walletPage.walletPageContents = widget.NewVBox()

	walletPage.walletTotalAmountLabel = make([]*canvas.Text, len(openedWalletIDs))
	for index := range openedWalletIDs {
		walletPage.walletTotalAmountLabel[index] = canvas.NewText("", values.TransactionInfoColor)
	}
}

func isAllWalletVerified(multiWallet *dcrlibwallet.MultiWallet) bool {
	openedWalletIDs := multiWallet.OpenedWalletIDsRaw()
	if len(openedWalletIDs) == 0 {
		return true
	}
	sort.Ints(openedWalletIDs)

	for _, walletID := range openedWalletIDs {
		if multiWallet.WalletWithID(walletID).Seed != "" {
			return false
		}
	}

	return true
}
