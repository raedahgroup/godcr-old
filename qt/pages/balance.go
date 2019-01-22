package pages

import (
	"context"
	"fmt"
	"strings"

	"github.com/decred/dcrd/dcrutil"
	"github.com/raedahgroup/godcr/app/walletcore"
	"github.com/therecipe/qt/core"
	"github.com/therecipe/qt/widgets"
)

type balancePage struct {
	pageStub
	accountNameLabel     *widgets.QLabel
	balanceLabel         *widgets.QLabel
	spendableLabel       *widgets.QLabel
	lockedByTicketsLabel *widgets.QLabel
	votingAuthorityLabel *widgets.QLabel
	unconfirmedLabel     *widgets.QLabel
	balanceGrid          *widgets.QGridLayout
}

func (b *balancePage) SetupWithWallet(ctx context.Context, wallet walletcore.Wallet) *widgets.QWidget {
	pageContent := widgets.NewQWidget(nil, 0)

	// create Gridlayout to arrange child views in grid and align them to top
	b.balanceGrid = widgets.NewQGridLayout(pageContent)
	b.balanceGrid.SetAlign(core.Qt__AlignTop)

	// create Qlabels widgets to hold values
	b.accountNameLabel = widgets.NewQLabel(nil, 0)
	b.balanceLabel = widgets.NewQLabel(nil, 0)
	b.spendableLabel = widgets.NewQLabel(nil, 0)
	b.lockedByTicketsLabel = widgets.NewQLabel(nil, 0)
	b.votingAuthorityLabel = widgets.NewQLabel(nil, 0)
	b.unconfirmedLabel = widgets.NewQLabel(nil, 0)

	//set widget default values and add to grid at different location
	b.balanceGrid.AddWidget(widgets.NewQLabel2("Account Name", nil, 0), 0, 0, core.Qt__AlignLeft)
	b.balanceGrid.AddWidget(widgets.NewQLabel2("Balance", nil, 0), 0, 1, core.Qt__AlignLeft)
	b.balanceGrid.AddWidget(widgets.NewQLabel2("Spendable", nil, 0), 0, 2, core.Qt__AlignLeft)
	b.balanceGrid.AddWidget(widgets.NewQLabel2("Locked By Tickets", nil, 0), 0, 3, core.Qt__AlignLeft)
	b.balanceGrid.AddWidget(widgets.NewQLabel2("Voting Authority", nil, 0), 0, 4, core.Qt__AlignLeft)
	b.balanceGrid.AddWidget(widgets.NewQLabel2("Unconfirmed", nil, 0), 0, 5, core.Qt__AlignLeft)
	b.balanceGrid.AddWidget(b.accountNameLabel, 1, 0, core.Qt__AlignLeft)
	b.balanceGrid.AddWidget(b.balanceLabel, 1, 1, core.Qt__AlignLeft)
	b.balanceGrid.AddWidget(b.spendableLabel, 1, 2, core.Qt__AlignLeft)
	b.balanceGrid.AddWidget(b.lockedByTicketsLabel, 1, 3, core.Qt__AlignLeft)
	b.balanceGrid.AddWidget(b.votingAuthorityLabel, 1, 4, core.Qt__AlignLeft)
	b.balanceGrid.AddWidget(b.unconfirmedLabel, 1, 5, core.Qt__AlignLeft)
	pageContent.SetLayout(b.balanceGrid)

	// run get balance op in separate goroutine to avoid stalling this method
	go b.getAndDisplayBalance(wallet)

	return pageContent
}

func (b *balancePage) getAndDisplayBalance(wallet walletcore.Wallet) {
	accounts, err := wallet.AccountsOverview(walletcore.DefaultRequiredConfirmations)
	if err != nil {
		b.balanceLabel.SetText(fmt.Sprintf("Error reading account balance: %s", err.Error()))
		return
	}

	summarizeBalance := func(total, spendable dcrutil.Amount) string {

		if total == spendable {
			return total.String()
		} else {
			return fmt.Sprintf("Total %s (Spendable %s)", total.String(), spendable.String())
		}
	}

	//create array to hold values for different accounts linked to wallet
	accountNameSummary := make([]string, len(accounts))
	balanceSummary := make([]string, len(accounts))
	spendableSummary := make([]string, len(accounts))
	lockedByTicketsSummary := make([]string, len(accounts))
	votingAuthoritySummary := make([]string, len(accounts))
	unconfirmedSummary := make([]string, len(accounts))

	for i, account := range accounts {

		balanceText := summarizeBalance(account.Balance.Total, account.Balance.Spendable)

		accountNameSummary[i] = fmt.Sprintf("%s", account.Name)
		balanceSummary[i] = fmt.Sprintf("%s", balanceText)
		spendableSummary[i] = fmt.Sprintf("%s", account.Balance.Spendable.String())
		lockedByTicketsSummary[i] = fmt.Sprintf("%s", account.Balance.LockedByTickets.String())
		votingAuthoritySummary[i] = fmt.Sprintf("%s", account.Balance.VotingAuthority.String())
		unconfirmedSummary[i] = fmt.Sprintf("%s", account.Balance.Unconfirmed.String())

	}

	//set label text with returned value
	b.accountNameLabel.SetText(strings.Join(accountNameSummary, "\n"))
	b.balanceLabel.SetText(strings.Join(balanceSummary, "\n"))
	b.spendableLabel.SetText(strings.Join(spendableSummary, "\n"))
	b.lockedByTicketsLabel.SetText(strings.Join(lockedByTicketsSummary, "\n"))
	b.votingAuthorityLabel.SetText(strings.Join(votingAuthoritySummary, "\n"))
	b.unconfirmedLabel.SetText(strings.Join(unconfirmedSummary, "\n"))

}
