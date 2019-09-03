package libwallet

import (
	"context"
	"errors"
	"fmt"

	"github.com/decred/dcrd/dcrutil"
	"github.com/raedahgroup/dcrlibwallet"
	"github.com/raedahgroup/godcr/app/wallet"
)

func (lw *LibWallet) AccountBalance(accountNumber uint32, requiredConfirmations int32) (*wallet.Balance, error) {
	balance, err := lw.dcrlw.GetAccountBalance(int32(accountNumber), requiredConfirmations)
	if err != nil {
		return nil, err
	}

	return &wallet.Balance{
		Total:           dcrutil.Amount(balance.Total),
		Spendable:       dcrutil.Amount(balance.Spendable),
		LockedByTickets: dcrutil.Amount(balance.LockedByTickets),
		VotingAuthority: dcrutil.Amount(balance.VotingAuthority),
		Unconfirmed:     dcrutil.Amount(balance.UnConfirmed),
	}, nil
}

func (lw *LibWallet) AccountsOverview(requiredConfirmations int32) ([]*wallet.Account, error) {
	accounts, err := lw.dcrlw.GetAccountsRaw(requiredConfirmations)
	if err != nil {
		return nil, fmt.Errorf("error fetching accounts: %s", err.Error())
	}

	accountsOverview := make([]*wallet.Account, 0, len(accounts.Acc))

	for _, acc := range accounts.Acc {
		accountNumber := uint32(acc.Number)

		// skip zero-balance imported accounts
		if acc.Name == "imported" && acc.Balance.Total == 0 {
			continue
		}

		account := &wallet.Account{
			Name:   acc.Name,
			Number: accountNumber,
			Balance: &wallet.Balance{
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

func (lw *LibWallet) NextAccount(accountName string, passphrase string) (uint32, error) {
	return lw.dcrlw.NextAccountRaw(accountName, []byte(passphrase))
}

func (lw *LibWallet) AccountNumber(accountName string) (uint32, error) {
	return lw.dcrlw.AccountNumber(accountName)
}

func (lw *LibWallet) AccountName(accountNumber uint32) (string, error) {
	return lw.dcrlw.AccountName(int32(accountNumber)), nil
}

func (lw *LibWallet) AddressInfo(address string) (*dcrlibwallet.AddressInfo, error) {
	return lw.dcrlw.AddressInfo(address)
}

func (lw *LibWallet) ValidateAddress(address string) (bool, error) {
	return lw.dcrlw.IsAddressValid(address), nil
}

func (lw *LibWallet) ReceiveAddress(account uint32) (string, error) {
	return lw.dcrlw.CurrentAddress(int32(account))
}

func (lw *LibWallet) GenerateNewAddress(account uint32) (string, error) {
	return lw.dcrlw.NextAddress(int32(account))
}

func (lw *LibWallet) ChangePrivatePassphrase(_ context.Context, oldPass, newPass string) error {
	if oldPass == "" || newPass == "" {
		return errors.New("Passphrase cannot be empty")
	}
	return lw.dcrlw.ChangePrivatePassphrase([]byte(oldPass), []byte(newPass))
}

func (lw *LibWallet) NetType() string {
	return lw.activeNet.Params.Name
}
