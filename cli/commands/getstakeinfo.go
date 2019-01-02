package commands

import (
	"context"
	"fmt"
	"github.com/raedahgroup/godcr/app"
	"github.com/raedahgroup/godcr/cli/termio"
	"strings"
)

// GetStakeInfoCommand requests statistics about the wallet stakes.
type GetStakeInfoCommand struct {
	commanderStub
}

// Run displays information about wallet sakes, tickets and their statuses.
func (g GetStakeInfoCommand) Run(ctx context.Context, wallet app.WalletMiddleware) error {
	stakeInfo, err := wallet.StakeInfo(ctx)
	if err != nil {
		return err
	}
	output := strings.Builder{}
	output.WriteString(fmt.Sprintf("Tickets: %d\n", stakeInfo.Total))
	for _, ticket := range stakeInfo.Tickets {
		output.WriteString(fmt.Sprintf("%s\t%s\n", ticket.Hash, ticket.Status))
	}
	termio.PrintStringResult(strings.TrimSpace(output.String()))
	return nil
}
