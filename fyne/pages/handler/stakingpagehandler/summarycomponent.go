package stakingpagehandler

import (
	"fyne.io/fyne"
	"fyne.io/fyne/canvas"
	"fyne.io/fyne/layout"

	"github.com/raedahgroup/godcr/fyne/pages/handler/values"
	"github.com/raedahgroup/godcr/fyne/widgets"
)

func (stakingPage *StakingPageObjects) getStakingSummary() {
	colOne := fyne.NewContainerWithLayout(layout.NewVBoxLayout())
	colOne.AddObject(fyne.NewContainerWithLayout(layout.NewHBoxLayout(), canvas.NewText("Unmined:", values.DarkerBlueGrayTextColor), widgets.NewHSpacer(values.SpacerSize46),
		layout.NewSpacer(), canvas.NewText("8", values.DefaultTextColor)))
	colOne.AddObject(widgets.NewVSpacer(values.SpacerSize12))
	colOne.AddObject(fyne.NewContainerWithLayout(
		layout.NewHBoxLayout(), canvas.NewText("Missed:", values.DarkerBlueGrayTextColor), layout.NewSpacer(),
		canvas.NewText("4", values.DefaultTextColor)))
	colOne.AddObject(widgets.NewVSpacer(values.SpacerSize12))
	colOne.AddObject(fyne.NewContainerWithLayout(layout.NewHBoxLayout(),
		canvas.NewText("Allmempooltix:", values.DarkerBlueGrayTextColor), layout.NewSpacer(), canvas.NewText("0", values.DefaultTextColor)))

	colTwo := fyne.NewContainerWithLayout(layout.NewVBoxLayout())
	colTwo.AddObject(fyne.NewContainerWithLayout(layout.NewHBoxLayout(), canvas.NewText("Immature:", values.DarkerBlueGrayTextColor), widgets.NewHSpacer(values.SpacerSize46),
		layout.NewSpacer(), canvas.NewText("0", values.DefaultTextColor)))
	colTwo.AddObject(widgets.NewVSpacer(values.SpacerSize12))
	colTwo.AddObject(fyne.NewContainerWithLayout(
		layout.NewHBoxLayout(), canvas.NewText("Expired:", values.DarkerBlueGrayTextColor), layout.NewSpacer(),
		canvas.NewText("0", values.DefaultTextColor)))
	colTwo.AddObject(widgets.NewVSpacer(values.SpacerSize12))
	colTwo.AddObject(fyne.NewContainerWithLayout(layout.NewHBoxLayout(),
		canvas.NewText("Poolsize:", values.DarkerBlueGrayTextColor), layout.NewSpacer(), canvas.NewText("0", values.DefaultTextColor)))

	colThree := fyne.NewContainerWithLayout(layout.NewVBoxLayout())
	colThree.AddObject(fyne.NewContainerWithLayout(layout.NewHBoxLayout(), canvas.NewText("Live:", values.DarkerBlueGrayTextColor), widgets.NewHSpacer(values.SpacerSize46),
		layout.NewSpacer(), canvas.NewText("0", values.DefaultTextColor)))
	colThree.AddObject(widgets.NewVSpacer(values.SpacerSize12))
	colThree.AddObject(fyne.NewContainerWithLayout(
		layout.NewHBoxLayout(), canvas.NewText("Revoked:", values.DarkerBlueGrayTextColor), layout.NewSpacer(),
		canvas.NewText("0", values.DefaultTextColor)))
	colThree.AddObject(widgets.NewVSpacer(values.SpacerSize12))
	colThree.AddObject(fyne.NewContainerWithLayout(layout.NewHBoxLayout(),
		canvas.NewText("Total Subsidy:", values.DarkerBlueGrayTextColor), layout.NewSpacer(), canvas.NewText("100 DCR", values.DefaultTextColor)))

	colFour := fyne.NewContainerWithLayout(layout.NewVBoxLayout())
	colFour.AddObject(fyne.NewContainerWithLayout(layout.NewHBoxLayout(), canvas.NewText("Voted:", values.DarkerBlueGrayTextColor), widgets.NewHSpacer(values.SpacerSize46),
		layout.NewSpacer(), canvas.NewText("0", values.DefaultTextColor)))
	colFour.AddObject(widgets.NewVSpacer(values.SpacerSize12))
	colFour.AddObject(fyne.NewContainerWithLayout(
		layout.NewHBoxLayout(), canvas.NewText("Unspent:", values.DarkerBlueGrayTextColor), layout.NewSpacer(),
		canvas.NewText("0", values.DefaultTextColor)))
	colFour.AddObject(widgets.NewVSpacer(values.SpacerSize12))

	summaryData := fyne.NewContainerWithLayout(layout.NewHBoxLayout())
	summaryData.AddObject(colOne)
	summaryData.AddObject(widgets.NewHSpacer(values.SpacerSize14))
	summaryData.AddObject(colTwo)
	summaryData.AddObject(widgets.NewHSpacer(values.SpacerSize14))
	summaryData.AddObject(colThree)
	summaryData.AddObject(widgets.NewHSpacer(values.SpacerSize14))
	summaryData.AddObject(colFour)

	summaryDataLayout := fyne.NewContainerWithLayout(layout.NewCenterLayout(),
		widgets.NewBorder(values.TransactionInfoBorderColor, fyne.NewSize(30, 20), summaryData), summaryData)

	stakingPage.StakingPageContents.Append(summaryDataLayout)
}
