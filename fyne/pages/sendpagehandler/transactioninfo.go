package sendpagehandler

import (
	"errors"
	"image/color"

	"fyne.io/fyne"
	"fyne.io/fyne/canvas"
	"fyne.io/fyne/layout"
	"fyne.io/fyne/widget"

	"github.com/raedahgroup/godcr/fyne/assets"
	"github.com/raedahgroup/godcr/fyne/widgets"
)

func (sendPage *SendPageObjects) initTransactionDetails() error {

	icons, err := assets.GetIcons(assets.CollapseDropdown, assets.ExpandDropdown)
	if err != nil {
		return errors.New("Could not load transaction details icons")
	}

	sendPage.transactionFeeLabel = widget.NewLabel(NilAmount)
	sendPage.transactionSizeLabel = widget.NewLabel(ZeroByte)
	sendPage.totalCostLabel = widget.NewLabel(NilAmount)
	sendPage.balanceAfterSendLabel = widget.NewLabel(NilAmount)

	paintedtransactionInfoform := sendPage.transactionInfoWithBorder()
	paintedtransactionInfoform.Hide()

	var transactionFeeBox *widget.Box

	costAndBalanceAfterSendBox := widget.NewVBox()
	var transactionSizeDropdown *widgets.ClickableBox
	var container *fyne.Container

	transactionSizeDropdown = widgets.NewClickableBox(widget.NewHBox(widget.NewIcon(icons[assets.ExpandDropdown])), func() {
		if paintedtransactionInfoform.Hidden {
			transactionSizeDropdown.Box.Children[0] = widget.NewIcon(icons[assets.CollapseDropdown])
			paintedtransactionInfoform.Show()
		} else {
			transactionSizeDropdown.Box.Children[0] = widget.NewIcon(icons[assets.ExpandDropdown])
			paintedtransactionInfoform.Hide()
		}

		costAndBalanceAfterSendBox.Refresh()
		container.Layout = layout.NewFixedGridLayout(
			fyne.NewSize(paintedtransactionInfoform.MinSize().Width, costAndBalanceAfterSendBox.MinSize().Height))

		container.Refresh()
		transactionSizeDropdown.Refresh()
		paintedtransactionInfoform.Refresh()
		transactionFeeBox.Refresh()
		sendPage.SendPageContents.Refresh()
	})

	transactionFeeBox = widget.NewHBox(canvas.NewText("Transaction fee", color.RGBA{89, 109, 129, 255}), layout.NewSpacer(),
		sendPage.transactionFeeLabel, widgets.NewHSpacer(4), transactionSizeDropdown)

	costAndBalanceAfterSendBox.Append(transactionFeeBox)

	costAndBalanceAfterSendBox.Append(paintedtransactionInfoform)
	costAndBalanceAfterSendBox.Append(widgets.NewVSpacer(10))

	costAndBalanceAfterSendBox.Append(widget.NewHBox(
		canvas.NewText("Total cost", color.RGBA{89, 109, 129, 255}), layout.NewSpacer(), sendPage.totalCostLabel))

	costAndBalanceAfterSendBox.Append(widget.NewHBox(
		canvas.NewText("Balance after send", color.RGBA{89, 109, 129, 255}), layout.NewSpacer(), sendPage.balanceAfterSendLabel))

	container = fyne.NewContainerWithLayout(layout.NewFixedGridLayout(paintedtransactionInfoform.MinSize()), costAndBalanceAfterSendBox)

	sendPage.SendPageContents.Append(container)

	return nil
}

func (sendPage *SendPageObjects) transactionInfoWithBorder() *fyne.Container {
	transactionInfoForm := fyne.NewContainerWithLayout(layout.NewVBoxLayout())

	transactionInfoForm.AddObject(fyne.NewContainerWithLayout(
		layout.NewHBoxLayout(), canvas.NewText("Processing time", color.RGBA{89, 109, 129, 255}), widgets.NewHSpacer(46),
		layout.NewSpacer(), widget.NewLabelWithStyle("Approx. 10 mins (2 blocks)", fyne.TextAlignLeading, fyne.TextStyle{})))

	transactionInfoForm.AddObject(fyne.NewContainerWithLayout(
		layout.NewHBoxLayout(), canvas.NewText("Fee rate", color.RGBA{89, 109, 129, 255}), layout.NewSpacer(),
		widget.NewLabelWithStyle("0.0001 DCR/byte", fyne.TextAlignLeading, fyne.TextStyle{})))

	transactionInfoForm.AddObject(fyne.NewContainerWithLayout(layout.NewHBoxLayout(),
		canvas.NewText("Transaction size", color.RGBA{89, 109, 129, 255}), layout.NewSpacer(), sendPage.transactionSizeLabel))

	return fyne.NewContainerWithLayout(layout.NewCenterLayout(),
		widgets.NewBorder(color.RGBA{158, 158, 158, 180}, fyne.NewSize(20, 30), transactionInfoForm), transactionInfoForm)
}
