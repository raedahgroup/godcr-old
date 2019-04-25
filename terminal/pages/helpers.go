package pages

import (
	"fmt"
	"math"

	"github.com/decred/dcrd/dcrutil"
	"github.com/raedahgroup/godcr/app/walletcore"
)

func maxDecimalPlaces(amounts []int64) (maxDecimalPlaces int) {
	for _, amount := range amounts {
		decimalPortion := walletcore.DecimalPortion(dcrutil.Amount(amount).ToCoin())
		nDecimalPlaces := len(decimalPortion)
		if nDecimalPlaces > maxDecimalPlaces {
			maxDecimalPlaces = nDecimalPlaces
		}
	}
	return
}

func formatAmountDisplay(amount int64, maxDecimalPlaces int) string {
	dcrAmount := dcrutil.Amount(amount).ToCoin()
	wholeNumber := int(math.Floor(dcrAmount))
	decimalPortion := walletcore.DecimalPortion(dcrAmount)

	if len(decimalPortion) == 0 {
		return fmt.Sprintf("%d%-*s DCR", wholeNumber, maxDecimalPlaces+1, decimalPortion)
	} else {
		return fmt.Sprintf("%d.%-*s DCR", wholeNumber, maxDecimalPlaces, decimalPortion)
	}
}
