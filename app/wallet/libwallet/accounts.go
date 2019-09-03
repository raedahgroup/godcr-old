package libwallet

import (
	"fmt"

	"github.com/raedahgroup/dcrlibwallet"
)

func (lw *LibWallet) Accounts(requiredConfirmations int32) ([]*dcrlibwallet.Account, error) {
	accountsInfo, err := lw.dcrlw.GetAccountsRaw(requiredConfirmations)
	if err != nil {
		return nil, fmt.Errorf("error fetching accounts: %s", err.Error())
	}

	accounts := make([]*dcrlibwallet.Account, 0, len(accountsInfo.Acc))

	for _, acc := range accountsInfo.Acc {
		// skip zero-balance imported accounts
		if acc.Name == "imported" && acc.Balance.Total == 0 {
			continue
		}
		accounts = append(accounts, acc)
	}

	return accounts, nil
}

func (lw *LibWallet) AccountBalance(accountNumber uint32, requiredConfirmations int32) (*dcrlibwallet.Balance, error) {
	return lw.dcrlw.GetAccountBalance(int32(accountNumber), requiredConfirmations)
}

func (lw *LibWallet) CreateAccount(accountName string, passphrase string) (uint32, error) {
	return lw.dcrlw.NextAccountRaw(accountName, []byte(passphrase))
}

func (lw *LibWallet) AccountNumber(accountName string) (uint32, error) {
	return lw.dcrlw.AccountNumber(accountName)
}

func (lw *LibWallet) AccountName(accountNumber uint32) (string, error) {
	return lw.dcrlw.AccountName(int32(accountNumber)), nil
}
