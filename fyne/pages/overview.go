package pages

import (
	"fmt"
	"fyne.io/fyne"
	"github.com/raedahgroup/godcr/fyne/assets"
	"image/color"

	// "fyne.io/fyne/canvas"
	"fyne.io/fyne/layout"
	"fyne.io/fyne/widget"
	"github.com/raedahgroup/godcr/fyne/widgets"
	// "image/color"
)

const PageTitle = "Overview"

var app *AppInterface

// todo: display overview page (include sync progress UI elements)
// todo: register sync progress listener on overview page to update sync progress views
func overviewPageContent(appInterface *AppInterface) fyne.CanvasObject {
		app = appInterface
		return container()
}
func container () fyne.CanvasObject {
	return widget.NewVBox(
		title(),
		balance(),
		widgets.NewVSpacer(50),
		pageBoxes(),
		)
}

func title () fyne.CanvasObject {
	titleWidget := widget.NewLabelWithStyle(PageTitle, fyne.TextAlignLeading, fyne.TextStyle{Bold: true, Italic: true})
	return widget.NewHBox(widgets.NewHSpacer(18), titleWidget)
}

func balance () fyne.CanvasObject {
	dcrBalance := widgets.NewLargeText("315.08", color.Black)
	dcrDecimals := widgets.NewSmallText("193725 DCR", color.Black)
	decimalsBox := fyne.NewContainerWithLayout(layout.NewVBoxLayout(), widgets.NewVSpacer(6),dcrDecimals)
	return widget.NewHBox(widgets.NewHSpacer(18), widgets.NewVSpacer(10), dcrBalance, decimalsBox)
}

func pageBoxes() (object fyne.CanvasObject) {
	return widget.NewHBox(widgets.NewHSpacer(18), fyne.NewContainerWithLayout(layout.NewGridLayout(2),
		recentTransactionBox(),
		blockStatusBox(),
		))
}

func recentTransactionBox () (object fyne.CanvasObject) {
	icons, err := assets.GetIcons(assets.ReceiveIcon, assets.SendIcon)
	if err != nil {
		app.DisplayLaunchErrorAndExit(fmt.Sprintf("An error occured while loading app icons: %s", err))
		return
	}
	return widgets.NewOverviewBox(
		widgets.NewTransactionList(
			widgets.TransactionLister{"4.08",  "298071 DCR", "Pending", icons[assets.ReceiveIcon]}, func(){}),
		widgets.NewTransactionList(
			widgets.TransactionLister{"34.17", "458878 DCR", "Friday", icons[assets.ReceiveIcon]}, func(){}),
		widgets.NewTransactionList(
			widgets.TransactionLister{"134.17", "018472 DCR", "Jan 20", icons[assets.SendIcon]}, func(){}),
		widgets.NewVSpacer(10),
		widget.NewLabelWithStyle("see all", fyne.TextAlignCenter, fyne.TextStyle{}),
	)
}

func blockStatusBox () (object fyne.CanvasObject){
	return widget.NewLabelWithStyle("block status box", fyne.TextAlignCenter, fyne.TextStyle{})
}

func transactionList () {

}

