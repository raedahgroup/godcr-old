package sendpagehandler

import (
	"image/color"

	"fyne.io/fyne"
	"fyne.io/fyne/canvas"
	"fyne.io/fyne/layout"
	"fyne.io/fyne/widget"

	"github.com/raedahgroup/godcr/fyne/widgets"
)

func transactionInfoWithBorder(transactionSizeLabel *widget.Label) *fyne.Container {
	transactionInfoForm := fyne.NewContainerWithLayout(layout.NewVBoxLayout())

	transactionInfoForm.AddObject(fyne.NewContainerWithLayout(
		layout.NewHBoxLayout(), canvas.NewText("Processing time", color.RGBA{89, 109, 129, 255}), widgets.NewHSpacer(46),
		layout.NewSpacer(), widget.NewLabelWithStyle("Approx. 10 mins (2 blocks)", fyne.TextAlignLeading, fyne.TextStyle{})))

	transactionInfoForm.AddObject(fyne.NewContainerWithLayout(
		layout.NewHBoxLayout(), canvas.NewText("Fee rate", color.RGBA{89, 109, 129, 255}), layout.NewSpacer(),
		widget.NewLabelWithStyle("0.0001 DCR/byte", fyne.TextAlignLeading, fyne.TextStyle{})))

	transactionInfoForm.AddObject(fyne.NewContainerWithLayout(layout.NewHBoxLayout(),
		canvas.NewText("Transaction size", color.RGBA{89, 109, 129, 255}), layout.NewSpacer(), transactionSizeLabel))

	return fyne.NewContainerWithLayout(layout.NewCenterLayout(),
		widgets.NewBorder(color.RGBA{158, 158, 158, 180}, fyne.NewSize(20, 30), transactionInfoForm), transactionInfoForm)
}

func TransactionDetails(collapseDropdown, expandDropdown *fyne.StaticResource, transactionFeeLabel,
	transactionSizeLabel, totalCostLabel, balanceAfterSendLabel *widget.Label, contents *widget.Box) (container *fyne.Container) {

	paintedtransactionInfoform := transactionInfoWithBorder(transactionSizeLabel) // widget.NewHBox(layout.NewSpacer(), transactionInfoWithBorder(transactionSizeLabel), layout.NewSpacer())
	paintedtransactionInfoform.Hide()

	var transactionFeeBox *widget.Box

	costAndBalanceAfterSendBox := widget.NewVBox()
	var transactionSizeDropdown *widgets.ClickableBox

	transactionSizeDropdown = widgets.NewClickableBox(widget.NewHBox(widget.NewIcon(expandDropdown)), func() {
		if paintedtransactionInfoform.Hidden {
			transactionSizeDropdown.Box.Children[0] = widget.NewIcon(collapseDropdown)
			paintedtransactionInfoform.Show()
		} else {
			transactionSizeDropdown.Box.Children[0] = widget.NewIcon(expandDropdown)
			paintedtransactionInfoform.Hide()
		}

		costAndBalanceAfterSendBox.Refresh()
		container.Layout = layout.NewFixedGridLayout(
			fyne.NewSize(paintedtransactionInfoform.MinSize().Width, costAndBalanceAfterSendBox.MinSize().Height))

		container.Refresh()
		transactionSizeDropdown.Refresh()
		paintedtransactionInfoform.Refresh()
		transactionFeeBox.Refresh()
		contents.Refresh()
	})

	transactionFeeBox = widget.NewHBox(canvas.NewText("Transaction fee", color.RGBA{89, 109, 129, 255}), layout.NewSpacer(),
		transactionFeeLabel, widgets.NewHSpacer(4), transactionSizeDropdown)

	costAndBalanceAfterSendBox.Append(transactionFeeBox)

	costAndBalanceAfterSendBox.Append(paintedtransactionInfoform)
	costAndBalanceAfterSendBox.Append(widgets.NewVSpacer(10))

	costAndBalanceAfterSendBox.Append(widget.NewHBox(
		canvas.NewText("Total cost", color.RGBA{89, 109, 129, 255}), layout.NewSpacer(), totalCostLabel))

	costAndBalanceAfterSendBox.Append(widget.NewHBox(
		canvas.NewText("Balance after send", color.RGBA{89, 109, 129, 255}), layout.NewSpacer(), balanceAfterSendLabel))

	container = fyne.NewContainerWithLayout(layout.NewFixedGridLayout(paintedtransactionInfoform.MinSize()), costAndBalanceAfterSendBox)
	return
}
