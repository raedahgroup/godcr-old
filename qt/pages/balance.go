package pages

import (
	"context"
	"fmt"
	"github.com/decred/dcrd/dcrutil"
	"github.com/raedahgroup/godcr/app/walletcore"
	"github.com/therecipe/qt/core"
	"github.com/therecipe/qt/widgets"
	"strings"
)

type balancePage struct {
	pageStub
	balanceLabel *widgets.QLabel
}

func (b *balancePage) SetupWithWallet(ctx context.Context, wallet walletcore.Wallet) *widgets.QWidget {
	pageContent := widgets.NewQWidget(nil, 0)

	// create layout to arrange child views vertically and center them
	pageLayout := widgets.NewQVBoxLayout()
	pageLayout.SetAlign(core.Qt__AlignCenter)
	pageContent.SetLayout(pageLayout)

	// add views to page layout
	b.balanceLabel = widgets.NewQLabel(nil, 0)
	pageContent.Layout().AddWidget(b.balanceLabel)

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

	balanceSummary := make([]string, len(accounts))
	for i, account := range accounts {
		balanceText := summarizeBalance(account.Balance.Total, account.Balance.Spendable)
		balanceSummary[i] = fmt.Sprintf("%s \t %s", account.Name, balanceText)
	}

	b.balanceLabel.SetText(strings.Join(balanceSummary, "\n"))
}
