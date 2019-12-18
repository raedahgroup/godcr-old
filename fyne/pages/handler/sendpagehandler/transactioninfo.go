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

	sendPage.transactionFeeLabel = canvas.NewText(values.NilAmount, values.NilAmountColor)
	sendPage.totalCostLabel = canvas.NewText(values.NilAmount, values.NilAmountColor)
	sendPage.balanceAfterSendLabel = canvas.NewText(values.NilAmount, values.NilAmountColor)
	sendPage.transactionSizeLabel = canvas.NewText(values.ZeroByte, values.DefaultTextColor)

	borderedtransactionInfoform := sendPage.transactionInfoWithBorder()
	borderedtransactionInfoform.Hide()

	var transactionFeeBox *widget.Box

	costAndBalanceAfterSendBox := widget.NewVBox()
	var transactionSizeDropdown *widgets.ImageButton
	var transactionInfoContainer *fyne.Container

	transactionSizeDropdown = widgets.NewImageButton(icons[assets.ExpandDropdown], nil, func() {
		if borderedtransactionInfoform.Hidden {
			transactionSizeDropdown.SetIcon(icons[assets.CollapseDropdown])
			borderedtransactionInfoform.Show()
		} else {
			transactionSizeDropdown.SetIcon(icons[assets.ExpandDropdown])
			borderedtransactionInfoform.Hide()
		}

		transactionInfoContainer.Layout = layout.NewFixedGridLayout(
			fyne.NewSize(borderedtransactionInfoform.MinSize().Width, costAndBalanceAfterSendBox.MinSize().Height))
		sendPage.Window.Resize(sendPage.SendPageContents.MinSize().Add(fyne.NewSize(0, values.SpacerSize10)))
	})

	transactionFeeBox = widget.NewHBox(canvas.NewText(values.TransactionFee, values.TransactionInfoColor), layout.NewSpacer(),
		sendPage.transactionFeeLabel, canvas.NewText(values.DCR, values.DefaultTextColor), widgets.NewHSpacer(values.SpacerSize4), transactionSizeDropdown)

	costAndBalanceAfterSendBox.Append(transactionFeeBox)

	costAndBalanceAfterSendBox.Append(borderedtransactionInfoform)
	costAndBalanceAfterSendBox.Append(widgets.NewVSpacer(values.SpacerSize4))

	costAndBalanceAfterSendBox.Append(widget.NewHBox(
		canvas.NewText(values.TotalCost, values.TransactionInfoColor), layout.NewSpacer(), sendPage.totalCostLabel, canvas.NewText(values.DCR, values.DefaultTextColor)))

	costAndBalanceAfterSendBox.Append(widgets.NewVSpacer(values.SpacerSize6))

	costAndBalanceAfterSendBox.Append(widget.NewHBox(
		canvas.NewText(values.BalanceAfterSend, values.TransactionInfoColor), layout.NewSpacer(), sendPage.balanceAfterSendLabel, canvas.NewText(values.DCR, values.DefaultTextColor)))

	transactionInfoContainer = fyne.NewContainerWithLayout(layout.NewFixedGridLayout(borderedtransactionInfoform.MinSize()), costAndBalanceAfterSendBox)

	sendPage.SendPageContents.Append(transactionInfoContainer)

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
