package pages

import (
	"sort"

	"fyne.io/fyne"
	"fyne.io/fyne/widget"

	"github.com/raedahgroup/godcr/fyne/pages/handler/historypagehandler"
	"github.com/raedahgroup/godcr/fyne/pages/handler/values"
	"github.com/raedahgroup/godcr/fyne/widgets"
)

type historyPageDynamicData struct {
	Contents *widget.Box
}

var txHistory historyPageDynamicData

func historyPageContent(app *AppInterface) fyne.CanvasObject {
	openedWalletIDs := app.MultiWallet.OpenedWalletIDsRaw()
	if len(openedWalletIDs) == 0 {
		return widget.NewHBox(widgets.NewHSpacer(10), widget.NewLabelWithStyle(values.WalletsErr, fyne.TextAlignCenter, fyne.TextStyle{Bold: true}))
	}
	sort.Ints(openedWalletIDs)

	initHistoryPageDynamicContent(openedWalletIDs)

	initHistoryPage := historypagehandler.HistoryPageData{
		MultiWallet:         app.MultiWallet,
		HistoryPageContents: txHistory.Contents,
		Window:              app.Window,
		TabMenu:             app.tabMenu,
	}

	err := initHistoryPage.InitHistoryPage()
	if err != nil {
		return widget.NewLabelWithStyle(values.ReceivePageLoadErr, fyne.TextAlignLeading, fyne.TextStyle{})
	}

	return widget.NewHBox(widgets.NewHSpacer(values.Padding), txHistory.Contents, widgets.NewHSpacer(values.Padding))
}

func initHistoryPageDynamicContent(openedWalletIDs []int) {
	txHistory = historyPageDynamicData{}
	txHistory.Contents = widget.NewVBox()
}
