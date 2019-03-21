package helpers

import (
	"fmt"
	"math"
)

func AmountToString(amount float64) string {
	amount = math.Round(amount)
	return fmt.Sprintf("%d DCR", int(amount))
}
