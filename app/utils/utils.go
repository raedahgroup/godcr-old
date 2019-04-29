package utils

import (
	"fmt"
	"math"
	"strconv"
	"strings"
)

func DecimalPortion(n float64) string {
	decimalPlaces := fmt.Sprintf("%f", n-math.Floor(n))          // produces 0.xxxx0000
	decimalPlaces = strings.Replace(decimalPlaces, "0.", "", -1) // remove 0.
	decimalPlaces = strings.TrimRight(decimalPlaces, "0")        // remove trailing 0s
	return decimalPlaces
}

func SplitAmountIntoParts(amount float64) []string {
	balanceParts := make([]string, 3)

	wholeNumber := int(math.Floor(amount))
	balanceParts[0] = strconv.Itoa(wholeNumber)

	decimalPortion := DecimalPortion(amount)
	if len(decimalPortion) == 0 {
		balanceParts[2] = " DCR"
	} else if len(decimalPortion) <= 2 {
		balanceParts[1] = fmt.Sprintf(".%s DCR", decimalPortion)
	} else {
		balanceParts[1] = fmt.Sprintf(".%s", decimalPortion[0:2])
		balanceParts[2] = fmt.Sprintf("%s DCR", decimalPortion[2:])
	}

	return balanceParts
}

//func accountNameFromTransaction()
