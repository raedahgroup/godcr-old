package sendpagehandler

import (
	"errors"
	"image/color"

	"fyne.io/fyne"
	"fyne.io/fyne/canvas"
	"fyne.io/fyne/layout"
	"fyne.io/fyne/widget"

	"github.com/raedahgroup/godcr/fyne/assets"
	"github.com/raedahgroup/godcr/fyne/pages/handler/constantvalues"
	"github.com/raedahgroup/godcr/fyne/widgets"
)

func (sendPage *SendPageObjects) initTransactionDetails() error {

	icons, err := assets.GetIcons(assets.CollapseDropdown, assets.ExpandDropdown)
	if err != nil {
		return errors.New(constantvalues.TransactionDetailsIconErr)
	}

	sendPage.transactionFeeLabel = widget.NewLabel(constantvalues.NilAmount)
	sendPage.transactionSizeLabel = widget.NewLabel(constantvalues.ZeroByte)
	sendPage.totalCostLabel = widget.NewLabel(constantvalues.NilAmount)
	sendPage.balanceAfterSendLabel = widget.NewLabel(constantvalues.NilAmount)

	paintedtransactionInfoform := sendPage.transactionInfoWithBorder()
	paintedtransactionInfoform.Hide()

	var transactionFeeBox *widget.Box

	costAndBalanceAfterSendBox := widget.NewVBox()
	var transactionSizeDropdown *widgets.ImageButton
	var container *fyne.Container

	transactionSizeDropdown = widgets.NewImageButton(icons[assets.ExpandDropdown], nil, func() {
		if paintedtransactionInfoform.Hidden {
			transactionSizeDropdown.SetIcon(icons[assets.CollapseDropdown])
			paintedtransactionInfoform.Show()
		} else {
			transactionSizeDropdown.SetIcon(icons[assets.ExpandDropdown])
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

	transactionFeeBox = widget.NewHBox(canvas.NewText(constantvalues.TransactionFee, color.RGBA{89, 109, 129, 255}), layout.NewSpacer(),
		sendPage.transactionFeeLabel, widgets.NewHSpacer(4), transactionSizeDropdown)

	costAndBalanceAfterSendBox.Append(transactionFeeBox)

	costAndBalanceAfterSendBox.Append(paintedtransactionInfoform)
	costAndBalanceAfterSendBox.Append(widgets.NewVSpacer(10))

	costAndBalanceAfterSendBox.Append(widget.NewHBox(
		canvas.NewText(constantvalues.TotalCost, color.RGBA{89, 109, 129, 255}), layout.NewSpacer(), sendPage.totalCostLabel))

	costAndBalanceAfterSendBox.Append(widget.NewHBox(
		canvas.NewText(constantvalues.BalanceAfterSend, color.RGBA{89, 109, 129, 255}), layout.NewSpacer(), sendPage.balanceAfterSendLabel))

	container = fyne.NewContainerWithLayout(layout.NewFixedGridLayout(paintedtransactionInfoform.MinSize()), costAndBalanceAfterSendBox)

	sendPage.SendPageContents.Append(container)

	return nil
}

func (sendPage *SendPageObjects) transactionInfoWithBorder() *fyne.Container {
	transactionInfoForm := fyne.NewContainerWithLayout(layout.NewVBoxLayout())

	transactionInfoForm.AddObject(fyne.NewContainerWithLayout(
		layout.NewHBoxLayout(), canvas.NewText(constantvalues.ProcessingTime, color.RGBA{89, 109, 129, 255}), widgets.NewHSpacer(46),
		layout.NewSpacer(), widget.NewLabelWithStyle(constantvalues.ProcessingTimeInfo, fyne.TextAlignLeading, fyne.TextStyle{})))

	transactionInfoForm.AddObject(fyne.NewContainerWithLayout(
		layout.NewHBoxLayout(), canvas.NewText(constantvalues.FeeRate, color.RGBA{89, 109, 129, 255}), layout.NewSpacer(),
		widget.NewLabelWithStyle(constantvalues.FeeRateInfo, fyne.TextAlignLeading, fyne.TextStyle{})))

	transactionInfoForm.AddObject(fyne.NewContainerWithLayout(layout.NewHBoxLayout(),
		canvas.NewText(constantvalues.TransactionSize, color.RGBA{89, 109, 129, 255}), layout.NewSpacer(), sendPage.transactionSizeLabel))

	return fyne.NewContainerWithLayout(layout.NewCenterLayout(),
		widgets.NewBorder(color.RGBA{158, 158, 158, 180}, fyne.NewSize(20, 30), transactionInfoForm), transactionInfoForm)
}
