package dcrlibwallet

import (
	"bytes"
	"context"
	"fmt"
	"sort"
	"time"

	"github.com/decred/dcrd/chaincfg/chainhash"
	"github.com/decred/dcrd/dcrutil"
	"github.com/decred/dcrd/wire"
	"github.com/raedahgroup/dcrlibwallet"
	"github.com/raedahgroup/dcrlibwallet/txhelper"
	"github.com/raedahgroup/godcr/app/walletcore"
)

func (lib *DcrWalletLib) AccountBalance(accountNumber uint32, requiredConfirmations int32) (*walletcore.Balance, error) {
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

func (lib *DcrWalletLib) AccountsOverview(requiredConfirmations int32) ([]*walletcore.Account, error) {
	accounts, err := lib.walletLib.GetAccountsRaw(requiredConfirmations)
	if err != nil {
		return nil, fmt.Errorf("error fetching accounts: %s", err.Error())
	}

	accountsOverview := make([]*walletcore.Account, 0, len(accounts.Acc))

	for _, acc := range accounts.Acc {
		accountNumber := uint32(acc.Number)

		// skip zero-balance imported accounts
		if acc.Name == "imported" && acc.Balance.Total == 0 {
			continue
		}

		account := &walletcore.Account{
			Name:   acc.Name,
			Number: accountNumber,
			Balance: &walletcore.Balance{
				Total:           dcrutil.Amount(acc.Balance.Total),
				Spendable:       dcrutil.Amount(acc.Balance.Spendable),
				LockedByTickets: dcrutil.Amount(acc.Balance.LockedByTickets),
				VotingAuthority: dcrutil.Amount(acc.Balance.VotingAuthority),
				Unconfirmed:     dcrutil.Amount(acc.Balance.UnConfirmed),
			},
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

func (lib *DcrWalletLib) AccountName(accountNumber uint32) (string, error) {
	return lib.walletLib.AccountName(accountNumber), nil
}

func (lib *DcrWalletLib) AddressInfo(address string) (*txhelper.AddressInfo, error) {
	return lib.AddressInfo(address)
}

func (lib *DcrWalletLib) ValidateAddress(address string) (bool, error) {
	return lib.walletLib.IsAddressValid(address), nil
}

func (lib *DcrWalletLib) ReceiveAddress(account uint32) (string, error) {
	return lib.walletLib.CurrentAddress(int32(account))
}

func (lib *DcrWalletLib) GenerateNewAddress(account uint32) (string, error) {
	return lib.walletLib.NextAddress(int32(account))
}

func (lib *DcrWalletLib) UnspentOutputs(account uint32, targetAmount int64, requiredConfirmations int32) ([]*walletcore.UnspentOutput, error) {
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

		address, err := walletcore.GetAddressFromPkScript(lib.activeNet.Params, utxo.PkScript)
		if err != nil {
			return nil, err
		}

		txn, err := lib.GetTransaction(txHash)
		if err != nil {
			return nil, fmt.Errorf("error reading transaction: %s", err.Error())
		}

		unspentOutputs[i] = &walletcore.UnspentOutput{
			OutputKey:       fmt.Sprintf("%s:%d", txHash, utxo.OutputIndex),
			TransactionHash: txHash,
			OutputIndex:     utxo.OutputIndex,
			Tree:            utxo.Tree,
			ReceiveTime:     utxo.ReceiveTime,
			Amount:          dcrutil.Amount(utxo.Amount),
			Address:         address,
			Confirmations:   txn.Confirmations,
		}
	}

	return unspentOutputs, nil
}

func (lib *DcrWalletLib) SendFromAccount(sourceAccount uint32, requiredConfirmations int32, destinations []txhelper.TransactionDestination, passphrase string) (string, error) {
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

func (lib *DcrWalletLib) SendFromUTXOs(sourceAccount uint32, requiredConfirmations int32, utxoKeys []string, txDestinations []txhelper.TransactionDestination, changeDestinations []txhelper.TransactionDestination, passphrase string) (string, error) {
	// fetch all utxos in account to extract details for the utxos selected by user
	// use targetAmount = 0 to fetch ALL utxos in account
	unspentOutputs, err := lib.UnspentOutputs(sourceAccount, 0, requiredConfirmations)
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

	unsignedTx, err := txhelper.NewUnsignedTx(inputs, txDestinations, changeDestinations)
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

	transactions := make([]*walletcore.Transaction, len(txs))
	for i, tx := range txs {
		_, txFee, txSize, txFeeRate, err := txhelper.MsgTxFeeSizeRate(tx.Transaction)
		if err != nil {
			return nil, err
		}

		transactions[i] = &walletcore.Transaction{
			Hash:          tx.Hash,
			Amount:        dcrutil.Amount(tx.Amount),
			Fee:           txFee,
			FeeRate:       txFeeRate,
			Size:          txSize,
			Type:          tx.Type,
			Direction:     tx.Direction,
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
	hash, err := chainhash.NewHashFromStr(transactionHash)
	if err != nil {
		return nil, fmt.Errorf("invalid hash: %s\n%s", transactionHash, err.Error())
	}

	txInfo, err := lib.walletLib.GetTransactionRaw(hash[:])
	if err != nil {
		return nil, err
	}

	decodedTx, err := txhelper.DecodeTransaction(hash, txInfo.Transaction, lib.activeNet.Params, lib.walletLib.AddressInfo)
	if err != nil {
		return nil, err
	}

	tx := &walletcore.Transaction{
		Hash:          txInfo.Hash,
		Amount:        dcrutil.Amount(txInfo.Amount),
		FormattedTime: time.Unix(txInfo.Timestamp, 0).Format("Mon Jan 2, 2006 3:04PM"),
		Timestamp:     txInfo.Timestamp,
		Fee:           dcrutil.Amount(decodedTx.Fee),
		Direction:     txInfo.Direction,
		Type:          txInfo.Type,
		FeeRate:       dcrutil.Amount(decodedTx.FeeRate),
		Size:          decodedTx.Size,
	}

	return &walletcore.TransactionDetails{
		BlockHeight:   txInfo.BlockHeight,
		Confirmations: txInfo.Confirmations,
		Transaction:   tx,
		Inputs:        decodedTx.Inputs,
		Outputs:       decodedTx.Outputs,
	}, nil
}

func (lib *DcrWalletLib) StakeInfo(ctx context.Context) (*walletcore.StakeInfo, error) {
	data, err := lib.walletLib.StakeInfo()
	if err != nil {
		return nil, fmt.Errorf("error getting stake info: %s", err.Error())
	}

	stakeInfo := &walletcore.StakeInfo{
		AllMempoolTix: data.AllMempoolTix,
		Expired:       data.Expired,
		Immature:      data.Immature,
		Live:          data.Live,
		Missed:        data.Missed,
		OwnMempoolTix: data.OwnMempoolTix,
		PoolSize:      data.PoolSize,
		Revoked:       data.Revoked,
		TotalSubsidy:  int64(data.TotalSubsidy),
		Unspent:       data.Unspent,
		Voted:         data.Voted,
	}

	return stakeInfo, nil
}

func (lib *DcrWalletLib) PurchaseTickets(ctx context.Context, request dcrlibwallet.PurchaseTicketsRequest) ([]string, error) {
	balance, err := lib.AccountBalance(request.Account, int32(request.RequiredConfirmations))
	if err != nil {
		return nil, fmt.Errorf("could not fetch account balance: %s", err.Error())
	}

	ticketPrice, err := lib.walletLib.TicketPrice(ctx)
	if err != nil {
		return nil, fmt.Errorf("could not determine ticket price: %s", err.Error())
	}

	if balance.Spendable < dcrutil.Amount(ticketPrice.TicketPrice*int64(request.NumTickets)) {
		return nil, fmt.Errorf("insufficient funds: spendable account balance (%s) is less than ticket price %s",
			balance.Spendable, dcrutil.Amount(ticketPrice.TicketPrice))
	}

	tickets, err := lib.walletLib.PurchaseTickets(ctx, &request)
	if err != nil {
		return nil, fmt.Errorf("could not complete ticket(s) purchase, encountered an error:\n%s", err.Error())
	}
	return tickets, nil
}

func (lib *DcrWalletLib) ConnectionInfo() *walletcore.CInfo {
	return &walletcore.CInfo{
		NumOfPeers: NumOfPeers,
		NBestBlock: int64(lib.walletLib.GetBestBlock()),
	}
}
