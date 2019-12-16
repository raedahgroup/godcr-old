package sendpagehandler

import (
	"errors"

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

	borderedtransactionInfoform := sendPage.transactionInfoWithBorder()
	borderedtransactionInfoform.Hide()

	var transactionFeeBox *widget.Box

	sendPage.costAndBalanceAfterSendBox = widget.NewVBox()
	var transactionSizeDropdown *widgets.ImageButton

	transactionSizeDropdown = widgets.NewImageButton(icons[assets.ExpandDropdown], nil, func() {
		if borderedtransactionInfoform.Hidden {
			transactionSizeDropdown.SetIcon(icons[assets.CollapseDropdown])
			borderedtransactionInfoform.Show()
		} else {
			transactionSizeDropdown.SetIcon(icons[assets.ExpandDropdown])
			borderedtransactionInfoform.Hide()
		}

		sendPage.transactionInfoContainer.Layout = layout.NewFixedGridLayout(
			fyne.NewSize(borderedtransactionInfoform.MinSize().Width, sendPage.costAndBalanceAfterSendBox.MinSize().Height))
		sendPage.Window.Resize(sendPage.SendPageContents.MinSize().Add(fyne.NewSize(0, values.SpacerSize10)))
	})

	transactionFeeBox = widget.NewHBox(canvas.NewText(values.TransactionFee, values.TransactionInfoColor), layout.NewSpacer(),
		sendPage.transactionFeeLabel, widgets.NewHSpacer(values.SpacerSize4), transactionSizeDropdown)

	sendPage.costAndBalanceAfterSendBox.Append(transactionFeeBox)

	sendPage.costAndBalanceAfterSendBox.Append(borderedtransactionInfoform)
	sendPage.costAndBalanceAfterSendBox.Append(widgets.NewVSpacer(values.SpacerSize4))

	sendPage.costAndBalanceAfterSendBox.Append(widget.NewHBox(
		canvas.NewText(values.TotalCost, values.TransactionInfoColor), layout.NewSpacer(), sendPage.totalCostLabel))

	sendPage.costAndBalanceAfterSendBox.Append(widgets.NewVSpacer(values.SpacerSize6))

	sendPage.costAndBalanceAfterSendBox.Append(widget.NewHBox(
		canvas.NewText(values.BalanceAfterSend, values.TransactionInfoColor), layout.NewSpacer(), sendPage.balanceAfterSendLabel))

	sendPage.transactionInfoContainer = fyne.NewContainerWithLayout(layout.NewFixedGridLayout(borderedtransactionInfoform.MinSize()), sendPage.costAndBalanceAfterSendBox)

	sendPage.SendPageContents.Append(sendPage.transactionInfoContainer)

	return nil
}

func (sendPage *SendPageObjects) transactionInfoWithBorder() *fyne.Container {
	transactionInfoForm := fyne.NewContainerWithLayout(layout.NewVBoxLayout())

	transactionInfoForm.AddObject(fyne.NewContainerWithLayout(
		layout.NewHBoxLayout(), canvas.NewText(values.ProcessingTime, values.DefaultTextColor), widgets.NewHSpacer(values.SpacerSize46),
		layout.NewSpacer(), widget.NewLabelWithStyle(values.ProcessingTimeInfo, fyne.TextAlignLeading, fyne.TextStyle{})))

	transactionInfoForm.AddObject(fyne.NewContainerWithLayout(
		layout.NewHBoxLayout(), canvas.NewText(values.FeeRate, values.DefaultTextColor), layout.NewSpacer(),
		widget.NewLabelWithStyle(values.FeeRateInfo, fyne.TextAlignLeading, fyne.TextStyle{})))

	transactionInfoForm.AddObject(fyne.NewContainerWithLayout(layout.NewHBoxLayout(),
		canvas.NewText(values.TransactionSize, values.DefaultTextColor), layout.NewSpacer(), sendPage.transactionSizeLabel))

	return fyne.NewContainerWithLayout(layout.NewCenterLayout(),
		widgets.NewBorder(values.TransactionInfoBorderColor, fyne.NewSize(20, 30), transactionInfoForm), transactionInfoForm)
}
