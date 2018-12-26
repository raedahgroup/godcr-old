package mobilewalletlib

import (
	"fmt"
	"sort"
	"time"

	"github.com/decred/dcrd/chaincfg/chainhash"
	"github.com/decred/dcrd/dcrutil"
	"github.com/raedahgroup/dcrcli/app/walletcore"
)

func (lib *MobileWalletLib) AccountBalance(accountNumber uint32) (*walletcore.Balance, error) {
	// pass 0 as requiredConfirmations
	balance, err := lib.walletLib.GetAccountBalance(accountNumber, 0)
	if err != nil {
		return nil, err
	}

	return &walletcore.Balance{
		Total:           dcrutil.Amount(balance.Total),
		Spendable:       dcrutil.Amount(balance.Spendable),
		LockedByTickets: dcrutil.Amount(balance.LockedByTickets),
		VotingAuthority: dcrutil.Amount(balance.VotingAuthority),
		Unconfirmed:     dcrutil.Amount(balance.UnConfirmed),
	}, nil
}

func (lib *MobileWalletLib) AccountsOverview() ([]*walletcore.Account, error) {
	// pass 0 as requiredConfirmations
	accounts, err := lib.walletLib.GetAccountsRaw(0)
	if err != nil {
		return nil, fmt.Errorf("error fetching accounts: %s", err.Error())
	}

	accountsOverview := make([]*walletcore.Account, 0, len(accounts.Acc))

	for _, acc := range accounts.Acc {
		accountNumber := uint32(acc.Number)

		balance, err := lib.AccountBalance(accountNumber)
		if err != nil {
			return nil, err
		}

		// skip zero-balance imported accounts
		if acc.Name == "imported" && balance.Total == 0 {
			continue
		}

		account := &walletcore.Account{
			Name:    acc.Name,
			Number:  accountNumber,
			Balance: balance,
		}
		accountsOverview = append(accountsOverview, account)
	}

	return accountsOverview, nil
}

func (lib *MobileWalletLib) NextAccount(accountName string, passphrase string) (uint32, error) {
	return lib.walletLib.NextAccountRaw(accountName, []byte(passphrase))
}

func (lib *MobileWalletLib) AccountNumber(accountName string) (uint32, error) {
	return lib.walletLib.AccountNumber(accountName)
}

func (lib *MobileWalletLib) GenerateReceiveAddress(account uint32) (string, error) {
	return lib.walletLib.CurrentAddress(int32(account))
}

func (lib *MobileWalletLib) ValidateAddress(address string) (bool, error) {
	return lib.walletLib.IsAddressValid(address), nil
}

func (lib *MobileWalletLib) UnspentOutputs(account uint32, targetAmount int64) ([]*walletcore.UnspentOutput, error) {
	return nil, fmt.Errorf("not yet implemented")
}

func (lib *MobileWalletLib) SendFromAccount(amountInDCR float64, sourceAccount uint32, destinationAddress, passphrase string) (string, error) {
	// convert amount from float64 DCR to int64 Atom
	amountInAtom, err := dcrutil.NewAmount(amountInDCR)
	if err != nil {
		return "", err
	}
	amount := int64(amountInAtom)

	txHash, err := lib.walletLib.SendTransaction([]byte(passphrase), destinationAddress, amount,
		int32(sourceAccount), 0, false)

	if err != nil {
		return "", err
	}

	transactionHash, err := chainhash.NewHash(txHash)
	if err != nil {
		return "", fmt.Errorf("error parsing successful transaction hash: %s", err.Error())
	}

	return transactionHash.String(), nil
}

func (lib *MobileWalletLib) SendFromUTXOs(utxoKeys []string, dcrAmount float64, account uint32, destAddress, passphrase string) (string, error) {
	return "", fmt.Errorf("not yet implemented")
}

func (lib *MobileWalletLib) TransactionHistory() ([]*walletcore.Transaction, error) {
	txs, err := lib.walletLib.GetTransactionsRaw()
	if err != nil {
		return nil, err
	}

	txDirection := func(direction int32) walletcore.TransactionDirection {
		if direction < int32(walletcore.TransactionDirectionUnclear) {
			return walletcore.TransactionDirection(direction)
		} else {
			return walletcore.TransactionDirectionUnclear
		}
	}

	transactions := make([]*walletcore.Transaction, len(txs))
	for i, tx := range txs {
		transactions[i] = &walletcore.Transaction{
			Hash:          tx.Hash,
			Amount:        dcrutil.Amount(tx.Amount).ToCoin(),
			Fee:           dcrutil.Amount(tx.Fee).ToCoin(),
			Type:          tx.Type,
			Direction:     txDirection(tx.Direction),
			Timestamp:     tx.Timestamp,
			FormattedTime: time.Unix(tx.Timestamp, 0).Format("Mon Jan 2, 2006 3:04PM"),
		}
	}

	// sort transactions by date (list newer first)
	sort.SliceStable(transactions, func(i1, i2 int) bool {
		return transactions[i1].Timestamp > transactions[i2].Timestamp
	})

	return transactions, nil
}
