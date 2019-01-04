package commands

import (
	"context"
	"errors"
	"fmt"
	"github.com/raedahgroup/godcr/app"
	"github.com/raedahgroup/godcr/cli/termio"
)

// StakeInfoCommand requests statistics about the wallet stakes.
type StakeInfoCommand struct {
	commanderStub
}

// Run displays information about wallet stakes, tickets and their statuses.
func (g StakeInfoCommand) Run(ctx context.Context, wallet app.WalletMiddleware) error {
	stakeInfo, err := wallet.StakeInfo(ctx)
	if err != nil {
		return err
	}
	if stakeInfo == nil {
		return errors.New("no tickets in wallet")
	}
	output := fmt.Sprintf("stake info for wallet:\n" + "total %d  ", stakeInfo.Total)
	if stakeInfo.Immature > 0 {
		output += fmt.Sprintf("immature %d  ", stakeInfo.Immature)
	}
	if stakeInfo.Unspent > 0 {
		output += fmt.Sprintf("live %d  ", stakeInfo.Unspent)
	}
	if stakeInfo.OwnMempoolTix > 0 {
		output += fmt.Sprintf("unmined %d", stakeInfo.OwnMempoolTix)
	}
	termio.PrintStringResult(output)
	return nil
}
