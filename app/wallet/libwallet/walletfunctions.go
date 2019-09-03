package libwallet

import (
	"context"
	"errors"
	"fmt"
	"strings"

	"github.com/decred/dcrd/chaincfg/chainhash"
	"github.com/decred/dcrd/dcrutil"
	"github.com/raedahgroup/dcrlibwallet"
	"github.com/raedahgroup/dcrlibwallet/addresshelper"
	"github.com/raedahgroup/dcrlibwallet/txhelper"
	"github.com/raedahgroup/dcrlibwallet/txindex"
	"github.com/raedahgroup/godcr/app/wallet"
)

func (lib *DcrWalletLib) AccountBalance(accountNumber uint32, requiredConfirmations int32) (*wallet.Balance, error) {
	balance, err := lib.walletLib.GetAccountBalance(accountNumber, requiredConfirmations)
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

func (lib *DcrWalletLib) AccountsOverview(requiredConfirmations int32) ([]*wallet.Account, error) {
	accounts, err := lib.walletLib.GetAccountsRaw(requiredConfirmations)
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

func (lib *DcrWalletLib) NextAccount(accountName string, passphrase string) (uint32, error) {
	return lib.walletLib.NextAccountRaw(accountName, []byte(passphrase))
}

func (lib *DcrWalletLib) AccountNumber(accountName string) (uint32, error) {
	return lib.walletLib.AccountNumber(accountName)
}

func (lib *DcrWalletLib) AccountName(accountNumber uint32) (string, error) {
	return lib.walletLib.AccountName(accountNumber), nil
}

func (lib *DcrWalletLib) AddressInfo(address string) (*dcrlibwallet.AddressInfo, error) {
	return lib.walletLib.AddressInfo(address)
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

func (lib *DcrWalletLib) UnspentOutputs(account uint32, targetAmount int64, requiredConfirmations int32) ([]*wallet.UnspentOutput, error) {
	utxos, err := lib.walletLib.UnspentOutputs(account, requiredConfirmations, targetAmount)
	if err != nil {
		return nil, err
	}

	unspentOutputs := make([]*wallet.UnspentOutput, len(utxos))
	for i, utxo := range utxos {
		hash, err := chainhash.NewHash(utxo.TransactionHash)
		if err != nil {
			return nil, err
		}
		txHash := hash.String()

		addresses, err := addresshelper.PkScriptAddresses(lib.activeNet.Params, utxo.PkScript)
		if err != nil {
			return nil, err
		}
		address := strings.Join(addresses, ", ")

		txn, err := lib.GetTransaction(txHash)
		if err != nil {
			return nil, fmt.Errorf("error reading transaction: %s", err.Error())
		}

		unspentOutputs[i] = &wallet.UnspentOutput{
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

func (lib *DcrWalletLib) SendFromUTXOs(sourceAccount uint32, requiredConfirmations int32, utxoKeys []string,
	txDestinations []txhelper.TransactionDestination, changeDestinations []txhelper.TransactionDestination, passphrase string) (string, error) {

	return lib.walletLib.SendFromCustomInputs(sourceAccount, requiredConfirmations, utxoKeys, txDestinations,
		changeDestinations, []byte(passphrase))
}

func (lib *DcrWalletLib) TransactionCount(filter *txindex.ReadFilter) (int, error) {
	return lib.walletLib.TxCount(filter)
}

func (lib *DcrWalletLib) TransactionHistory(offset, count int32, filter *txindex.ReadFilter) ([]*wallet.Transaction, error) {
	txs, err := lib.walletLib.GetTransactionsRaw(offset, count, filter)
	if err != nil {
		return nil, err
	}

	processedTxs := make([]*wallet.Transaction, len(txs))
	for i, tx := range txs {
		confirmations := txhelper.TxConfirmations(tx.BlockHeight, lib.walletLib.GetBestBlock())
		processedTxs[i] = wallet.TxDetails(tx, confirmations)
	}
	return processedTxs, nil
}

func (lib *DcrWalletLib) GetTransaction(transactionHash string) (*wallet.Transaction, error) {
	hash, err := chainhash.NewHashFromStr(transactionHash)
	if err != nil {
		return nil, fmt.Errorf("invalid hash: %s\n%s", transactionHash, err.Error())
	}

	tx, err := lib.walletLib.GetTransactionRaw(hash[:])
	if err != nil {
		return nil, err
	}

	confirmations := txhelper.TxConfirmations(tx.BlockHeight, lib.walletLib.GetBestBlock())
	return wallet.TxDetails(tx, confirmations), nil
}

func (lib *DcrWalletLib) StakeInfo(ctx context.Context) (*wallet.StakeInfo, error) {
	data, err := lib.walletLib.StakeInfo()
	if err != nil {
		return nil, fmt.Errorf("error getting stake info: %s", err.Error())
	}

	stakeInfo := &wallet.StakeInfo{
		AllMempoolTix: data.AllMempoolTix,
		Expired:       data.Expired,
		Immature:      data.Immature,
		Live:          data.Live,
		Missed:        data.Missed,
		OwnMempoolTix: data.OwnMempoolTix,
		PoolSize:      data.PoolSize,
		Revoked:       data.Revoked,
		TotalSubsidy:  dcrutil.Amount(data.TotalSubsidy).String(),
		Unspent:       data.Unspent,
		Voted:         data.Voted,
	}

	return stakeInfo, nil
}

func (lib *DcrWalletLib) TicketPrice(ctx context.Context) (int64, error) {
	ticketPrice, err := lib.walletLib.TicketPrice(ctx)
	if err != nil {
		return 0, err
	}

	return ticketPrice.TicketPrice, nil
}

func (lib *DcrWalletLib) PurchaseTicket(ctx context.Context, request dcrlibwallet.PurchaseTicketsRequest) ([]string, error) {
	balance, err := lib.AccountBalance(request.Account, int32(request.RequiredConfirmations))
	if err != nil {
		return nil, fmt.Errorf("could not fetch account balance: %s", err.Error())
	}

	ticketPrice, err := lib.TicketPrice(ctx)
	if err != nil {
		return nil, fmt.Errorf("could not determine ticket price: %s", err.Error())
	}

	totalTicketPrice := dcrutil.Amount(ticketPrice * int64(request.NumTickets))
	if balance.Spendable < totalTicketPrice {
		return nil, fmt.Errorf("insufficient funds: spendable account balance (%s) is less than ticket purchase cost %s",
			balance.Spendable, totalTicketPrice)
	}

	tickets, err := lib.walletLib.PurchaseTickets(ctx, &request)
	if err != nil {
		return nil, fmt.Errorf("could not complete ticket(s) purchase, encountered an error:\n%s", err.Error())
	}
	return tickets, nil
}

func (lib *DcrWalletLib) ChangePrivatePassphrase(_ context.Context, oldPass, newPass string) error {
	if oldPass == "" || newPass == "" {
		return errors.New("Passphrase cannot be empty")
	}
	return lib.walletLib.ChangePrivatePassphrase([]byte(oldPass), []byte(newPass))
}

func (lib *DcrWalletLib) NetType() string {
	return lib.activeNet.Params.Name
}
