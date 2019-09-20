package core

import (
	"fmt"

	"github.com/decred/dcrd/dcrutil/v2"
	"github.com/raedahgroup/dcrlibwallet"
)

func AccountsOverview(walletLib   *dcrlibwallet.LibWallet, requiredConfirmations int32) ([]*Account, error) {
	accounts, err := walletLib.GetAccountsRaw(requiredConfirmations)
	if err != nil {
		return nil, fmt.Errorf("error fetching accounts: %s", err.Error())
	}

	accountsOverview := make([]*Account, 0, len(accounts.Acc))

	for _, acc := range accounts.Acc {
		accountNumber := uint32(acc.Number)

		// skip zero-balance imported accounts
		if acc.Name == "imported" && acc.Balance.Total == 0 {
			continue
		}

		account := &Account{
			Name:   acc.Name,
			Number: accountNumber,
			Balance: &Balance{
				Total:           dcrutil.Amount(acc.Balance.Total),
				Spendable:       dcrutil.Amount(acc.Balance.Spendable),
				LockedByTickets: dcrutil.Amount(acc.Balance.LockedByTickets),
				VotingAuthority: dcrutil.Amount(acc.Balance.VotingAuthority),
				Unconfirmed:     dcrutil.Amount(acc.Balance.UnConfirmed),
			},
			ExternalKeyCount: acc.ExternalKeyCount,
			InternalKeyCount: acc.InternalKeyCount,
			ImportedKeyCount: acc.ImportedKeyCount,
		}
		accountsOverview = append(accountsOverview, account)
	}

	return accountsOverview, nil
}
