package commands

import (
	"context"
	"errors"
	"fmt"
	"github.com/raedahgroup/godcr/app/walletcore"
	"github.com/raedahgroup/godcr/cli/termio"
)

// StakeInfoCommand requests statistics about the wallet stakes.
type StakeInfoCommand struct {
	commanderStub
}

// Run displays information about wallet stakes, tickets and their statuses.
func (g StakeInfoCommand) Run(ctx context.Context, wallet walletcore.Wallet) error {
	stakeInfo, err := wallet.StakeInfo(ctx)
	if err != nil {
		return err
	}
	if stakeInfo == nil {
		return errors.New("no tickets in wallet")
	}
	output := fmt.Sprintf("stake info for wallet:\n"+"total %d  ", stakeInfo.Total)
	if stakeInfo.Expired > 0 {
		output += fmt.Sprintf("expired %d  ", stakeInfo.Expired)
	}
	if stakeInfo.Immature > 0 {
		output += fmt.Sprintf("immature %d  ", stakeInfo.Immature)
	}
	if stakeInfo.Live > 0 {
		output += fmt.Sprintf("live %d  ", stakeInfo.Live)
	}
	if stakeInfo.Revoked > 0 {
		output += fmt.Sprintf("revoked %d  ", stakeInfo.Revoked)
	}
	if stakeInfo.OwnMempoolTix > 0 {
		output += fmt.Sprintf("unmined %d", stakeInfo.OwnMempoolTix)
	}
	if stakeInfo.Unspent > 0 {
		output += fmt.Sprintf("unspent %d  ", stakeInfo.Unspent)
	}
	termio.PrintStringResult(output)
	return nil
}
