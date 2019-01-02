package dcrlibwallet

import (
	"bytes"
	"fmt"
	"sort"
	"time"

	"github.com/decred/dcrd/chaincfg/chainhash"
	"github.com/decred/dcrd/dcrutil"
	"github.com/decred/dcrd/wire"
	"github.com/raedahgroup/dcrlibwallet/txhelper"
	"github.com/raedahgroup/godcr/app/walletcore"
)

// ideally, we should let user provide this info in settings and use the user provided value
// using a constant now to make it easier to update the code where this value is required/used
const requiredConfirmations = 0

func (lib *DcrWalletLib) AccountBalance(accountNumber uint32) (*walletcore.Balance, error) {
	// pass 0 as requiredConfirmations
	balance, err := lib.walletLib.GetAccountBalance(accountNumber, requiredConfirmations)
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

func (lib *DcrWalletLib) AccountsOverview() ([]*walletcore.Account, error) {
	// pass 0 as requiredConfirmations
	accounts, err := lib.walletLib.GetAccountsRaw(requiredConfirmations)
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

func (lib *DcrWalletLib) NextAccount(accountName string, passphrase string) (uint32, error) {
	return lib.walletLib.NextAccountRaw(accountName, []byte(passphrase))
}

func (lib *DcrWalletLib) AccountNumber(accountName string) (uint32, error) {
	return lib.walletLib.AccountNumber(accountName)
}

func (lib *DcrWalletLib) GenerateReceiveAddress(account uint32) (string, error) {
	return lib.walletLib.CurrentAddress(int32(account))
}

func (lib *DcrWalletLib) ValidateAddress(address string) (bool, error) {
	return lib.walletLib.IsAddressValid(address), nil
}

func (lib *DcrWalletLib) UnspentOutputs(account uint32, targetAmount int64) ([]*walletcore.UnspentOutput, error) {
	utxos, err := lib.walletLib.UnspentOutputs(account, requiredConfirmations, targetAmount)
	if err != nil {
		return nil, err
	}

	unspentOutputs := make([]*walletcore.UnspentOutput, len(utxos))
	for i, utxo := range utxos {
		hash, err := chainhash.NewHash(utxo.TransactionHash)
		if err != nil {
			return nil, err
		}
		txHash := hash.String()

		unspentOutputs[i] = &walletcore.UnspentOutput{
			OutputKey:       fmt.Sprintf("%s:%d", txHash, utxo.OutputIndex),
			TransactionHash: txHash,
			OutputIndex:     utxo.OutputIndex,
			Tree:            utxo.Tree,
			ReceiveTime:     utxo.ReceiveTime,
			Amount:          dcrutil.Amount(utxo.Amount),
		}
	}

	return unspentOutputs, nil
}

func (lib *DcrWalletLib) SendFromAccount(sourceAccount uint32, destinations []txhelper.TransactionDestination, passphrase string) (string, error) {
	txHash, err := lib.walletLib.BulkSendTransaction([]byte(passphrase), destinations, int32(sourceAccount), requiredConfirmations)
	if err != nil {
		return "", err
	}

	transactionHash, err := chainhash.NewHash(txHash)
	if err != nil {
		return "", fmt.Errorf("error parsing successful transaction hash: %s", err.Error())
	}

	return transactionHash.String(), nil
}

func (lib *DcrWalletLib) SendFromUTXOs(sourceAccount uint32, utxoKeys []string, destinations []txhelper.TransactionDestination, passphrase string) (string, error) {
	// fetch all utxos in account to extract details for the utxos selected by user
	// use targetAmount = 0 to fetch ALL utxos in account
	unspentOutputs, err := lib.UnspentOutputs(sourceAccount, 0)
	if err != nil {
		return "", err
	}

	// loop through unspentOutputs to find user selected utxos
	inputs := make([]*wire.TxIn, 0, len(utxoKeys))
	for _, utxo := range unspentOutputs {
		useUtxo := false
		for _, key := range utxoKeys {
			if utxo.OutputKey == key {
				useUtxo = true
			}
		}
		if !useUtxo {
			continue
		}

		// this is a reverse conversion and should not throw an error
		// this string hash was originally chainhash.Hash and was converted to string in `lib.UnspentOutputs`
		txHash, _ := chainhash.NewHashFromStr(utxo.TransactionHash)

		outpoint := wire.NewOutPoint(txHash, utxo.OutputIndex, int8(utxo.Tree))
		input := wire.NewTxIn(outpoint, int64(utxo.Amount), nil)
		inputs = append(inputs, input)

		if len(inputs) == len(utxoKeys) {
			break
		}
	}

	// generate address from sourceAccount to receive change
	changeAddress, err := lib.GenerateReceiveAddress(sourceAccount)
	if err != nil {
		return "", err
	}

	unsignedTx, err := txhelper.NewUnsignedTx(inputs, destinations, changeAddress)
	if err != nil {
		return "", err
	}

	// serialize unsigned tx
	var txBuf bytes.Buffer
	txBuf.Grow(unsignedTx.SerializeSize())
	err = unsignedTx.Serialize(&txBuf)
	if err != nil {
		return "", fmt.Errorf("error serializing transaction: %s", err.Error())
	}

	txHash, err := lib.walletLib.SignAndPublishTransaction(txBuf.Bytes(), []byte(passphrase))
	if err != nil {
		return "", err
	}

	transactionHash, err := chainhash.NewHash(txHash)
	if err != nil {
		return "", fmt.Errorf("error parsing successful transaction hash: %s", err.Error())
	}

	return transactionHash.String(), nil
}

func (lib *DcrWalletLib) TransactionHistory() ([]*walletcore.Transaction, error) {
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
			Amount:        dcrutil.Amount(tx.Amount),
			Fee:           dcrutil.Amount(tx.Fee),
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

func (lib *DcrWalletLib) GetTransaction(transactionHash string) (*walletcore.TransactionDetails, error) {
	return nil, fmt.Errorf("not implemented")
}
