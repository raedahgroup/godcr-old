package utils

import (
	"fmt"
	"math"
	"strings"
)

func DecimalPortion(n float64) string {
	decimalPlaces := fmt.Sprintf("%f", n-math.Floor(n))          // produces 0.xxxx0000
	decimalPlaces = strings.Replace(decimalPlaces, "0.", "", -1) // remove 0.
	decimalPlaces = strings.TrimRight(decimalPlaces, "0")        // remove trailing 0s
	return decimalPlaces
}

//func accountNameFromTransaction()
