package pages

import (
	"fyne.io/fyne"
	"image/color"

	"fyne.io/fyne/layout"
	"fyne.io/fyne/widget"
	"github.com/raedahgroup/godcr/fyne/widgets"
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
	return fyne.NewContainerWithLayout(layout.NewHBoxLayout(), widgets.NewHSpacer(18), fyne.NewContainerWithLayout(layout.NewGridLayout(2),
		recentTransactionBox(),
		blockStatusBox(),
		))
}

func recentTransactionBox () (object fyne.CanvasObject) {
	//icons, err := assets.GetIcons(assets.ReceiveIcon, assets.SendIcon)
	//if err != nil {
	//	app.DisplayLaunchErrorAndExit(fmt.Sprintf("An error occured while loading app icons: %s", err))
	//	return
	//}
	return widget.NewVBox(

		widget.NewVBox(
			NewTransactionColumn("amount", "fee", "date", "direction"),
			NewTransactionColumn("32.0932334", "0.0000004", "08-11-2019", "yourself"),
			NewTransactionColumn("32.0932334", "0.0000004", "08-11-2019", "yourself"),
			NewTransactionColumn("32.0932334", "0.0000004", "08-11-2019", "yourself"),
			NewTransactionColumn("32.0932334", "0.0000004", "08-11-2019", "yourself"),
		),
	)
}

func blockStatusBox () (object fyne.CanvasObject){
	return widget.NewLabelWithStyle("block status box", fyne.TextAlignCenter, fyne.TextStyle{})
}

func transactionList () {

}

func NewTransactionColumn (amount, fee, date, direction string) fyne.CanvasObject {
	amountLabel := widget.NewLabel(amount)
	feeLabel := widget.NewLabel(fee)
	dateLabel := widget.NewLabel(date)
	directionLabel := widget.NewLabel(direction)
	column := widget.NewHBox(amountLabel, feeLabel, directionLabel, dateLabel)
	column.Resize(fyne.Size{Width:30, Height:10})
	return column
}

func TransactionColumnHeader() fyne.CanvasObject {
	amount := widget.NewLabelWithStyle("amount", fyne.TextAlignLeading, fyne.TextStyle{Bold:true})
	fee := widget.NewLabelWithStyle("fee", fyne.TextAlignLeading, fyne.TextStyle{Bold:true})
	direction := widget.NewLabelWithStyle("direction", fyne.TextAlignLeading, fyne.TextStyle{Bold:true})
	date := widget.NewLabelWithStyle("date", fyne.TextAlignLeading, fyne.TextStyle{Bold:true})
	return widget.NewHBox(amount, fee, direction, date)
}