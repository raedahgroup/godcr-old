package multipagecomponents

import (
	"fmt"
	"strings"

	"fyne.io/fyne"

	"github.com/raedahgroup/godcr/fyne/layouts"
	"github.com/raedahgroup/godcr/fyne/pages/handler/values"
	"github.com/raedahgroup/godcr/fyne/widgets"
)

func AmountFormatBox(amountInString string, textSize1, textSize2 int) *fyne.Container {
	//	amountInString := strconv.FormatFloat(dcrlibwallet.AmountCoin(amount), 'f', -1, 64)
	trailingDotForAmount := strings.Split(amountInString, ".")
	// if amount is a float
	amountLabelBox := fyne.NewContainerWithLayout(layouts.NewHBox(0, true))
	if len(trailingDotForAmount) > 1 && len(trailingDotForAmount[1]) > 2 {
		trailingAmountLabel := widgets.NewTextWithStyle(fmt.Sprintf("%s %s", trailingDotForAmount[1][2:], values.DCR),
			values.DefaultTextColor, fyne.TextStyle{Bold: true, Monospace: true}, fyne.TextAlignLeading, textSize2)
		leadingAmountLabel := widgets.NewTextWithStyle(trailingDotForAmount[0]+"."+trailingDotForAmount[1][:2],
			values.DefaultTextColor, fyne.TextStyle{Bold: true, Monospace: true}, fyne.TextAlignLeading, textSize1)

		amountLabelBox.AddObject(leadingAmountLabel)
		amountLabelBox.AddObject(trailingAmountLabel)

	} else {
		amountLabel := widgets.NewTextWithStyle(amountInString, values.DefaultTextColor,
			fyne.TextStyle{Bold: true, Monospace: true}, fyne.TextAlignLeading, values.TextSize24)

		DCRLabel := widgets.NewTextWithStyle(values.DCR, values.DefaultTextColor,
			fyne.TextStyle{Bold: true, Monospace: true}, fyne.TextAlignLeading, values.TextSize14)

		amountLabelBox.Layout = layouts.NewHBox(values.SpacerSize4, true)
		amountLabelBox.AddObject(amountLabel)
		amountLabelBox.AddObject(DCRLabel)
	}

	return amountLabelBox
}
