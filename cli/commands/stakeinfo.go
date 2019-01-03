package commands

import (
	"context"
	"errors"
	"fmt"
	"strings"

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
	output := strings.Builder{}
	output.WriteString(fmt.Sprintf("\nStake Info\n"+"Total tickets:\t%d\n", stakeInfo.Total))
	if stakeInfo.Immature > 0 {
		output.WriteString(fmt.Sprintf("Immature:\t%d\n", stakeInfo.Immature))
	}
	if stakeInfo.Unspent > 0 {
		output.WriteString(fmt.Sprintf("Live:\t%d\n", stakeInfo.Unspent))
	}
	if stakeInfo.OwnMempoolTix > 0 {
		output.WriteString(fmt.Sprintf("Unmined:\t%d\n", stakeInfo.OwnMempoolTix))
	}
	output.WriteString("\nTickets\n")
	for _, ticket := range stakeInfo.Tickets {
		output.WriteString(fmt.Sprintf("%s\t%s\n", ticket.Hash, ticket.Status))
	}
	termio.PrintStringResult(strings.TrimSpace(output.String()))
	return nil
}
