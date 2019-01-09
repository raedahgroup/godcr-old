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
	output := fmt.Sprintf("stake info for wallet:\n"+
		"expired %d  immature %d  live %d  revoked %d  unmined %d  unspent %d  "+
		"allmempooltix %d  poolsize %d  missed %d  voted %d  total subsidy %d",
		stakeInfo.Expired, stakeInfo.Immature, stakeInfo.Live, stakeInfo.Revoked,
		stakeInfo.OwnMempoolTix, stakeInfo.Unspent, stakeInfo.AllMempoolTix,
		stakeInfo.PoolSize, stakeInfo.Missed, stakeInfo.Voted, stakeInfo.TotalSubsidy)
	termio.PrintStringResult(output)
	return nil
}
