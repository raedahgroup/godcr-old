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
	stakes := "\nStake Info\n" +
		"Total tickets:\t%d\n" +
		"Immature:\t%d\n" +
		"Live:\t%d\n" +
		"Unspent:\t%d\n"
	output.WriteString(fmt.Sprintf(stakes,
		stakeInfo.Total,
		stakeInfo.Immature,
		stakeInfo.Live,
		stakeInfo.Unspent))
	output.WriteString("\nTickets\n")
	for _, ticket := range stakeInfo.Tickets {
		output.WriteString(fmt.Sprintf("%s\t%s\n", ticket.Hash, ticket.Status))
	}
	termio.PrintStringResult(strings.TrimSpace(output.String()))
	return nil
}
