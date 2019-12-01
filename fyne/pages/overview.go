package pages

import (
	"fyne.io/fyne"
	"fyne.io/fyne/canvas"
	"fyne.io/fyne/layout"
	"fyne.io/fyne/widget"
	"github.com/raedahgroup/godcr/fyne/assets"
	"github.com/raedahgroup/godcr/fyne/widgets"
	"image/color"
)

const PageTitle = "Overview"

type Overview struct {
	transactionBox *widget.Box
}

// todo: display overview page (include sync progress UI elements)
// todo: register sync progress listener on overview page to update sync progress views
func overviewPageContent() fyne.CanvasObject {
		return widget.NewHBox(widgets.NewHSpacer(18), container())
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
	return widget.NewHBox(titleWidget)
}

func balance () fyne.CanvasObject {
	dcrBalance := widgets.NewLargeText("315.08", color.Black)
	dcrDecimals := widgets.NewSmallText("193725 DCR", color.Black)
	decimalsBox := fyne.NewContainerWithLayout(layout.NewVBoxLayout(), widgets.NewVSpacer(6),dcrDecimals)
	return widget.NewHBox(widgets.NewVSpacer(10), dcrBalance, decimalsBox)
}

func pageBoxes() (object fyne.CanvasObject) {
	return fyne.NewContainerWithLayout(layout.NewVBoxLayout(),
		blockStatusBox(),
		recentTransactionBox(),
		)
}

func recentTransactionBox () fyne.CanvasObject {
	table := &widgets.Table{}
	table.NewTable(transactionColumnHeader(),
		newTransactionColumn(assets.ReceiveIcon,"0.0000004 DCR", "0.0000004 DCR", "yourself", "confirmed", "08-11-2019" ),
		newTransactionColumn(assets.SendIcon,"0.0000004 DCR", "0.0000004 DCR", "yourself", "confirmed", "08-11-2019" ),
		newTransactionColumn(assets.ReceiveIcon,"0.0000004 DCR", "0.0000004 DCR", "yourself", "confirmed", "08-11-2019" ),
		newTransactionColumn(assets.SendIcon,"0.0000004 DCR", "0.0000004 DCR", "yourself", "confirmed", "08-11-2019" ),
	)
	return table.Result
}

func blockStatusBox() fyne.CanvasObject {
	top := fyne.NewContainerWithLayout(layout.NewFixedGridLayout(fyne.NewSize(515, 24)),
				widget.NewHBox(
				widgets.NewSmallText("Syncing...", color.Black),
				layout.NewSpacer(),
				widget.NewButton("Cancel", func(){}),
				))
	progressBar := fyne.NewContainerWithLayout(layout.NewFixedGridLayout(fyne.NewSize(515, 20)),
			widget.NewProgressBar(),
			)
	timeLeft := widget.NewLabelWithStyle("6 min left", fyne.TextAlignLeading, fyne.TextStyle{Italic:true})
	connectedPeers := widget.NewLabelWithStyle("Connected peers count  6", fyne.TextAlignTrailing, fyne.TextStyle{Italic:true})
	syncSteps := widget.NewLabelWithStyle("Step 1/3", fyne.TextAlignTrailing, fyne.TextStyle{Italic:true})
	blockHeadersStatus := widget.NewLabelWithStyle("Fetching block headers  89%", fyne.TextAlignTrailing, fyne.TextStyle{Italic:true})
	syncDuration := fyne.NewContainerWithLayout(layout.NewBorderLayout(nil, nil, timeLeft, connectedPeers),
		timeLeft, connectedPeers)
	syncStatus := fyne.NewContainerWithLayout(layout.NewBorderLayout(nil, nil, syncSteps, blockHeadersStatus),
		syncSteps, blockHeadersStatus)

	bottom := fyne.NewContainerWithLayout(layout.NewGridLayout(2),
		walletSyncBox("Default", "waiting for other wallets", "6000 of 164864", "220 days behind"),
		walletSyncBox("Default", "waiting for other wallets", "6000 of 164864", "220 days behind"),
		)
	return fyne.NewContainerWithLayout(layout.NewVBoxLayout(),
		widgets.NewVSpacer(5),
		top,
		progressBar,
		syncDuration,
		syncStatus,
		widgets.NewVSpacer(15),
		bottom,
		)
}

func walletSyncBox (name, status, headerFetched, progress string) fyne.CanvasObject {
	blackColor := color.Black
	nameText := widgets.NewTextWithSize(name, blackColor, 12)
	statusText := widgets.NewTextWithSize(status, blackColor, 10)
	headerFetchedTitleText := widgets.NewTextWithSize("Block header fetched", blackColor, 12)
	headerFetchedText := widgets.NewTextWithSize(headerFetched, blackColor, 10)
	progressTitleText := widgets.NewTextWithSize("Syncing progress", blackColor, 12)
	progressText := widgets.NewTextWithSize(progress, blackColor, 10)
	top := fyne.NewContainerWithLayout(layout.NewHBoxLayout(), nameText, layout.NewSpacer(), statusText)
	middle := fyne.NewContainerWithLayout(layout.NewHBoxLayout(), headerFetchedTitleText, layout.NewSpacer(), headerFetchedText)
	bottom := fyne.NewContainerWithLayout(layout.NewHBoxLayout(), progressTitleText, layout.NewSpacer(), progressText)
	//middle := fyne.NewContainerWithLayout(layout.NewBorderLayout(nil, nil, headerFetchedTitleLabel, headerFetchedLabel),
	//	nameLabel, statusLabel)
	//bottom := fyne.NewContainerWithLayout(layout.NewBorderLayout(nil, nil, progressTitleLabel, progressLabel),
	//	nameLabel, statusLabel)
	return fyne.NewContainerWithLayout(layout.NewFixedGridLayout(fyne.NewSize(250, 70)),
			top,
			layout.NewSpacer(),
			middle,
			layout.NewSpacer(),
			bottom,
	)
}

func newTransactionColumn (transactionType, amount, fee, direction, status, date string) *widget.Box {
	icons, _ := assets.GetIcons(assets.ReceiveIcon, assets.SendIcon)
	icon := canvas.NewImageFromResource(icons[transactionType])
	// spacer := widgets.NewHSpacer(10)
	icon.SetMinSize(fyne.NewSize(5, 20))
	iconBox := widget.NewVBox(widgets.NewVSpacer(4), icon)
	amountLabel := widget.NewLabel(amount)
	feeLabel := widget.NewLabel(fee)
	dateLabel := widget.NewLabel(date)
	statusLabel := widget.NewLabel(status)
	directionLabel := widget.NewLabel(direction)
	column := widget.NewHBox(iconBox, amountLabel, feeLabel, directionLabel, statusLabel,  dateLabel)
	return column
}

//func newTransactionColumn (transactionType, amount, fee, direction, status, date string) fyne.CanvasObject {
//	icons, _ := assets.GetIcons(assets.ReceiveIcon, assets.SendIcon)
//	icon := canvas.NewImageFromResource(icons[transactionType])
//	// spacer := widgets.NewHSpacer(5)
//	amountLabel := widget.NewLabel(amount)
//	feeLabel := widget.NewLabel(fee)
//	dateLabel := widget.NewLabel(date)
//	statusLabel := widget.NewLabel(status)
//	directionLabel := widget.NewLabel(direction)
//	column := fyne.NewContainerWithLayout(layout.NewGridLayout(6), icon, amountLabel, feeLabel, directionLabel, statusLabel,  dateLabel)
//	return column
//}

func transactionColumnHeader() *widget.Box {
	// spacer := widgets.NewHSpacer(20)
	hash := widget.NewLabelWithStyle("#", fyne.TextAlignCenter, fyne.TextStyle{Bold: true})
	amount := widget.NewLabelWithStyle("amount", fyne.TextAlignCenter, fyne.TextStyle{Bold: true})
	fee := widget.NewLabelWithStyle("fee", fyne.TextAlignCenter, fyne.TextStyle{Bold: true})
	direction := widget.NewLabelWithStyle("direction", fyne.TextAlignCenter, fyne.TextStyle{Bold: true})
	status := widget.NewLabelWithStyle("status", fyne.TextAlignCenter, fyne.TextStyle{Bold: true})
	date := widget.NewLabelWithStyle("date", fyne.TextAlignCenter, fyne.TextStyle{Bold: true})
	return widget.NewHBox(hash, amount, fee, direction, status, date)
}

func (tb *Overview) addHeader () {
	spacer := widgets.NewHSpacer(10)
	amountLabel := widget.NewLabelWithStyle("amount", fyne.TextAlignCenter, fyne.TextStyle{Bold: true})
	feeLabel := widget.NewLabelWithStyle("fee", fyne.TextAlignCenter, fyne.TextStyle{Bold: true})
	directionLabel := widget.NewLabelWithStyle("direction", fyne.TextAlignCenter, fyne.TextStyle{Bold: true})
	statusLabel := widget.NewLabelWithStyle("status", fyne.TextAlignCenter, fyne.TextStyle{Bold: true})
	dateLabel := widget.NewLabelWithStyle("date", fyne.TextAlignCenter, fyne.TextStyle{Bold: true})
	header := fyne.NewContainerWithLayout(layout.NewGridLayout(6), spacer, amountLabel, feeLabel, directionLabel, statusLabel, dateLabel)
	//for _, child := range header.Children {
	//	child.Resize(fyne.NewSize(100, 100))
	//}
	tb.transactionBox.Children = append(tb.transactionBox.Children, header)
}

func (tb *Overview) addTransactionColumn (transactionType, amount, fee, direction, status, date string) {
	transactionBox := tb.transactionBox
	if len(transactionBox.Children) > 0 {
		transactionBox.Children = append(transactionBox.Children, newTransactionColumn(transactionType, amount, fee, direction, status, date))
	} else {
		tb.addHeader()
		transactionBox.Children = append(transactionBox.Children, newTransactionColumn(transactionType, amount, fee, direction, status, date))
	}
}