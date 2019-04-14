package pages

import (
	"fmt"
	"math"
	"strings"

	"github.com/decred/dcrd/dcrutil"
)

func decimalPlaces(n float64) string {
	decimalPlaces := fmt.Sprintf("%f", n-math.Floor(n))
	decimalPlaces = strings.Replace(decimalPlaces, "0", "", -1)
	decimalPlaces = strings.Replace(decimalPlaces, ".", "", -1)
	return decimalPlaces
}

func maxDecimalPlaces(amounts []int64) (maxDecimalPlaces int) {
	for _, amount := range amounts {
		decimalPlaces := decimalPlaces(dcrutil.Amount(amount).ToCoin())
		nDecimalPlaces := len(decimalPlaces)
		if nDecimalPlaces > maxDecimalPlaces {
			maxDecimalPlaces = nDecimalPlaces
		}
	}
	return
}

func formatAmountDisplay(amount int64, maxDecimalPlaces int) string {
	dcrAmount := dcrutil.Amount(amount).ToCoin()
	wholeNumber := int(math.Floor(dcrAmount))
	decimalPlaces := decimalPlaces(dcrAmount)

	if len(decimalPlaces) == 0 {
		return fmt.Sprintf("%d%-*s DCR", wholeNumber, maxDecimalPlaces+1, decimalPlaces)
	} else {
		return fmt.Sprintf("%d.%-*s DCR", wholeNumber, maxDecimalPlaces, decimalPlaces)
	}
}
