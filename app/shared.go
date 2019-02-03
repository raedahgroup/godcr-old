package app

import (
	"fmt"

	"github.com/raedahgroup/godcr/app/walletcore"
)

func SimpleBalance(balance *walletcore.Balance, detailed bool) string {
	if detailed || balance.Total == balance.Spendable {
		return balance.Total.String()
	} else {
		return fmt.Sprintf("Total %s (Spendable %s)", balance.Total.String(), balance.Spendable.String())
	}
}
