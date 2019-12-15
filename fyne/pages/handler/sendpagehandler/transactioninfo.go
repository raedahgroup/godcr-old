package sendpagehandler

import (
	"errors"
	"image/color"

	"fyne.io/fyne"
	"fyne.io/fyne/canvas"
	"fyne.io/fyne/layout"
	"fyne.io/fyne/widget"

	"github.com/raedahgroup/godcr/fyne/assets"
	"github.com/raedahgroup/godcr/fyne/pages/handler/values"
	"github.com/raedahgroup/godcr/fyne/widgets"
)

func (sendPage *SendPageObjects) initTransactionDetails() error {

	icons, err := assets.GetIcons(assets.CollapseDropdown, assets.ExpandDropdown)
	if err != nil {
		return errors.New(values.TransactionDetailsIconErr)
	}

	sendPage.transactionFeeLabel = widget.NewLabel(values.NilAmount)
	sendPage.transactionSizeLabel = widget.NewLabel(values.ZeroByte)
	sendPage.totalCostLabel = widget.NewLabel(values.NilAmount)
	sendPage.balanceAfterSendLabel = widget.NewLabel(values.NilAmount)

	paintedtransactionInfoform := sendPage.transactionInfoWithBorder()
	paintedtransactionInfoform.Hide()

	var transactionFeeBox *widget.Box

	costAndBalanceAfterSendBox := widget.NewVBox()
	var transactionSizeDropdown *widgets.ImageButton

	transactionSizeDropdown = widgets.NewImageButton(icons[assets.ExpandDropdown], nil, func() {
		if paintedtransactionInfoform.Hidden {
			transactionSizeDropdown.SetIcon(icons[assets.CollapseDropdown])
			paintedtransactionInfoform.Show()
		} else {
			transactionSizeDropdown.SetIcon(icons[assets.ExpandDropdown])
			paintedtransactionInfoform.Hide()
		}

		sendPage.transactionInfoContainer.Layout = layout.NewFixedGridLayout(
			fyne.NewSize(paintedtransactionInfoform.MinSize().Width, costAndBalanceAfterSendBox.MinSize().Height))

		sendPage.Window.Resize(sendPage.SendPageContents.MinSize())
	})

	transactionFeeBox = widget.NewHBox(canvas.NewText(values.TransactionFee, color.RGBA{89, 109, 129, 255}), layout.NewSpacer(),
		sendPage.transactionFeeLabel, widgets.NewHSpacer(values.SpacerSize4), transactionSizeDropdown)

	costAndBalanceAfterSendBox.Append(transactionFeeBox)

	costAndBalanceAfterSendBox.Append(paintedtransactionInfoform)

	costAndBalanceAfterSendBox.Append(widget.NewHBox(
		canvas.NewText(values.TotalCost, color.RGBA{89, 109, 129, 255}), layout.NewSpacer(), sendPage.totalCostLabel))

	costAndBalanceAfterSendBox.Append(widget.NewHBox(
		canvas.NewText(values.BalanceAfterSend, color.RGBA{89, 109, 129, 255}), layout.NewSpacer(), sendPage.balanceAfterSendLabel))

	sendPage.transactionInfoContainer = fyne.NewContainerWithLayout(layout.NewFixedGridLayout(paintedtransactionInfoform.MinSize()), costAndBalanceAfterSendBox)

	sendPage.SendPageContents.Append(sendPage.transactionInfoContainer)

	return nil
}

func (sendPage *SendPageObjects) transactionInfoWithBorder() *fyne.Container {
	transactionInfoForm := fyne.NewContainerWithLayout(layout.NewVBoxLayout())

	transactionInfoForm.AddObject(fyne.NewContainerWithLayout(
		layout.NewHBoxLayout(), canvas.NewText(values.ProcessingTime, color.RGBA{89, 109, 129, 255}), widgets.NewHSpacer(values.SpacerSize46),
		layout.NewSpacer(), widget.NewLabelWithStyle(values.ProcessingTimeInfo, fyne.TextAlignLeading, fyne.TextStyle{})))

	transactionInfoForm.AddObject(fyne.NewContainerWithLayout(
		layout.NewHBoxLayout(), canvas.NewText(values.FeeRate, color.RGBA{89, 109, 129, 255}), layout.NewSpacer(),
		widget.NewLabelWithStyle(values.FeeRateInfo, fyne.TextAlignLeading, fyne.TextStyle{})))

	transactionInfoForm.AddObject(fyne.NewContainerWithLayout(layout.NewHBoxLayout(),
		canvas.NewText(values.TransactionSize, color.RGBA{89, 109, 129, 255}), layout.NewSpacer(), sendPage.transactionSizeLabel))

	return fyne.NewContainerWithLayout(layout.NewCenterLayout(),
		widgets.NewBorder(color.RGBA{158, 158, 158, 180}, fyne.NewSize(20, 30), transactionInfoForm), transactionInfoForm)
}
