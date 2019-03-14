package dcrwalletrpc

import (
	"bytes"
	"context"
	"fmt"
	"io"
	"sort"
	"strings"

	"github.com/decred/dcrd/chaincfg/chainhash"
	"github.com/decred/dcrd/dcrutil"
	"github.com/decred/dcrd/wire"
	"github.com/decred/dcrwallet/rpc/walletrpc"
	"github.com/raedahgroup/dcrlibwallet"
	"github.com/raedahgroup/dcrlibwallet/addresshelper"
	"github.com/raedahgroup/dcrlibwallet/txhelper"
	"github.com/raedahgroup/godcr/app/walletcore"
	"google.golang.org/grpc/codes"
)

func (c *WalletRPCClient) AccountBalance(accountNumber uint32, requiredConfirmations int32) (*walletcore.Balance, error) {
	req := &walletrpc.BalanceRequest{
		AccountNumber:         accountNumber,
		RequiredConfirmations: requiredConfirmations,
	}

	res, err := c.walletService.Balance(context.Background(), req)
	if err != nil {
		return nil, fmt.Errorf("error fetching balance for account: %d \n:%s", accountNumber, err.Error())
	}

	return &walletcore.Balance{
		Total:           dcrutil.Amount(res.Total),
		Spendable:       dcrutil.Amount(res.Spendable),
		LockedByTickets: dcrutil.Amount(res.LockedByTickets),
		VotingAuthority: dcrutil.Amount(res.VotingAuthority),
		Unconfirmed:     dcrutil.Amount(res.Unconfirmed),
	}, nil
}

func (c *WalletRPCClient) AccountsOverview(requiredConfirmations int32) ([]*walletcore.Account, error) {
	accounts, err := c.walletService.Accounts(context.Background(), &walletrpc.AccountsRequest{})
	if err != nil {
		return nil, fmt.Errorf("error fetching accounts: %s", err.Error())
	}

	accountsOverview := make([]*walletcore.Account, 0, len(accounts.Accounts))

	for _, acc := range accounts.Accounts {
		balance, err := c.AccountBalance(acc.AccountNumber, requiredConfirmations)
		if err != nil {
			return nil, err
		}

		// skip zero-balance imported accounts
		if acc.AccountName == "imported" && balance.Total == 0 {
			continue
		}

		account := &walletcore.Account{
			Name:    acc.AccountName,
			Number:  acc.AccountNumber,
			Balance: balance,
		}
		accountsOverview = append(accountsOverview, account)
	}

	return accountsOverview, nil
}

func (c *WalletRPCClient) NextAccount(accountName string, passphrase string) (uint32, error) {
	req := &walletrpc.NextAccountRequest{
		AccountName: accountName,
		Passphrase:  []byte(passphrase),
	}

	nextAccount, err := c.walletService.NextAccount(context.Background(), req)
	if err != nil {
		return 0, err
	}

	return nextAccount.AccountNumber, nil
}

func (c *WalletRPCClient) AccountNumber(accountName string) (uint32, error) {
	req := &walletrpc.AccountNumberRequest{
		AccountName: accountName,
	}

	r, err := c.walletService.AccountNumber(context.Background(), req)
	if err != nil {
		return 0, err
	}

	return r.AccountNumber, nil
}

func (c *WalletRPCClient) AccountName(accountNumber uint32) (string, error) {
	accounts, err := c.walletService.Accounts(context.Background(), &walletrpc.AccountsRequest{})
	if err != nil {
		return "", err
	}

	for _, account := range accounts.Accounts {
		if account.AccountNumber == accountNumber {
			return account.AccountName, nil
		}
	}

	return "", fmt.Errorf("Account not found")
}

func (c *WalletRPCClient) AddressInfo(address string) (*txhelper.AddressInfo, error) {
	req := &walletrpc.ValidateAddressRequest{
		Address: address,
	}

	addressValidationResult, err := c.walletService.ValidateAddress(context.Background(), req)
	if err != nil {
		return nil, err
	}

	addressInfo := &txhelper.AddressInfo{
		IsMine:  addressValidationResult.IsMine,
		Address: address,
	}
	if addressValidationResult.IsMine {
		addressInfo.AccountNumber = addressValidationResult.AccountNumber
		addressInfo.AccountName, _ = c.AccountName(addressValidationResult.AccountNumber)
	}

	return addressInfo, nil
}

// ValidateAddress tries to decode an address for the given network params, if error is encountered, address is not valid
func (c *WalletRPCClient) ValidateAddress(address string) (bool, error) {
	_, err := addresshelper.DecodeForNetwork(address, c.activeNet.Params)
	return err == nil, nil
}

// ReceiveAddress uses GAP_POLICY_WRAP which returns previously generated unused addresses ONLY if the gap limit is exceeded
// Ideally, ReceiveAddress should always return the last generated address that has not been used
func (c *WalletRPCClient) ReceiveAddress(account uint32) (string, error) {
	req := &walletrpc.NextAddressRequest{
		Account:   account,
		GapPolicy: walletrpc.NextAddressRequest_GAP_POLICY_WRAP,
		Kind:      walletrpc.NextAddressRequest_BIP0044_EXTERNAL,
	}

	nextAddress, err := c.walletService.NextAddress(context.Background(), req)
	if err != nil {
		return "", err
	}

	return nextAddress.Address, nil
}

// GenerateNewAddress uses GAP_POLICY_WRAP which returns previously generated unused addresses ONLY if the gap limit is exceeded
// Ideally, GenerateNewAddress should always generate new addresses but we have to be wary of issues that could arise if the gap limit is exceeded
// GenerateNewAddress will continue to generate new addresses until/unless the gap limit is met, then it'll revert to previously generated addresses
func (c *WalletRPCClient) GenerateNewAddress(account uint32) (string, error) {
	req := &walletrpc.NextAddressRequest{
		Account:   account,
		GapPolicy: walletrpc.NextAddressRequest_GAP_POLICY_WRAP,
		Kind:      walletrpc.NextAddressRequest_BIP0044_EXTERNAL,
	}

	nextAddress, err := c.walletService.NextAddress(context.Background(), req)
	if err != nil {
		return "", err
	}

	return nextAddress.Address, nil
}

func (c *WalletRPCClient) UnspentOutputs(account uint32, targetAmount int64, requiredConfirmations int32) ([]*walletcore.UnspentOutput, error) {
	utxoStream, err := c.unspentOutputStream(account, targetAmount, requiredConfirmations)
	if err != nil {
		return nil, err
	}

	var unspentOutputs []*walletcore.UnspentOutput

	for {
		utxo, err := utxoStream.Recv()
		if err == io.EOF {
			break
		}
		if err != nil {
			return nil, err
		}

		hash, err := chainhash.NewHash(utxo.TransactionHash)
		if err != nil {
			return nil, err
		}
		txHash := hash.String()

		addresses, err := addresshelper.PkScriptAddresses(c.activeNet.Params, utxo.PkScript)
		if err != nil {
			return nil, err
		}
		address := strings.Join(addresses, ", ")

		txn, err := c.GetTransaction(txHash)
		if err != nil {
			return nil, fmt.Errorf("error reading transaction: %s", err.Error())
		}

		unspentOutput := &walletcore.UnspentOutput{
			OutputKey:       fmt.Sprintf("%s:%d", txHash, utxo.OutputIndex),
			TransactionHash: txHash,
			OutputIndex:     utxo.OutputIndex,
			Tree:            utxo.Tree,
			ReceiveTime:     utxo.ReceiveTime,
			Amount:          dcrutil.Amount(utxo.Amount),
			Address:         address,
			Confirmations:   txn.Confirmations,
		}
		unspentOutputs = append(unspentOutputs, unspentOutput)
	}

	return unspentOutputs, nil
}

func (c *WalletRPCClient) SendFromAccount(sourceAccount uint32, requiredConfirmations int32, destinations []txhelper.TransactionDestination, passphrase string) (string, error) {
	// construct non-change outputs for all recipients
	outputs := make([]*walletrpc.ConstructTransactionRequest_Output, len(destinations))
	for i, destination := range destinations {
		amountInAtom, err := txhelper.AmountToAtom(destination.Amount)
		if err != nil {
			return "", err
		}

		outputs[i] = &walletrpc.ConstructTransactionRequest_Output{
			Destination: &walletrpc.ConstructTransactionRequest_OutputDestination{
				Address: destination.Address,
			},
			Amount: amountInAtom,
		}
	}

	// construct transaction
	constructRequest := &walletrpc.ConstructTransactionRequest{
		SourceAccount:         sourceAccount,
		NonChangeOutputs:      outputs,
		RequiredConfirmations: requiredConfirmations,
	}

	constructResponse, err := c.walletService.ConstructTransaction(context.Background(), constructRequest)
	if err != nil {
		return "", fmt.Errorf("error constructing transaction: %s", err.Error())
	}

	return c.signAndPublishTransaction(constructResponse.UnsignedTransaction, passphrase)
}

func (c *WalletRPCClient) SendFromUTXOs(sourceAccount uint32, requiredConfirmations int32, utxoKeys []string, txDestinations []txhelper.TransactionDestination, changeDestinations []txhelper.TransactionDestination, passphrase string) (string, error) {
	// fetch all utxos in account to extract details for the utxos selected by user
	// passing 0 as targetAmount to c.unspentOutputStream fetches ALL utxos in account
	utxoStream, err := c.unspentOutputStream(sourceAccount, 0, requiredConfirmations)
	if err != nil {
		return "", err
	}

	// loop through utxo stream to find user selected utxos
	inputs := make([]*wire.TxIn, 0, len(utxoKeys))
	for {
		utxo, err := utxoStream.Recv()
		if err == io.EOF {
			break
		}
		if err != nil {
			return "", err
		}

		transactionHash, err := chainhash.NewHash(utxo.TransactionHash)
		if err != nil {
			return "", fmt.Errorf("invalid utxo transaction hash: %s", err.Error())
		}
		utxoKey := fmt.Sprintf("%s:%d", transactionHash.String(), utxo.OutputIndex)

		useUtxo := false
		for _, key := range utxoKeys {
			if utxoKey == key {
				useUtxo = true
			}
		}
		if !useUtxo {
			continue
		}

		outpoint := wire.NewOutPoint(transactionHash, utxo.OutputIndex, int8(utxo.Tree))
		input := wire.NewTxIn(outpoint, utxo.Amount, nil)
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

	return c.signAndPublishTransaction(txBuf.Bytes(), passphrase)
}

func (c *WalletRPCClient) TransactionHistory() ([]*walletcore.Transaction, error) {
	req := &walletrpc.GetTransactionsRequest{}

	stream, err := c.walletService.GetTransactions(context.Background(), req)
	if err != nil {
		return nil, err
	}

	var transactions []*walletcore.Transaction

	for {
		in, err := stream.Recv()
		if err == io.EOF {
			break
		}
		if err != nil {
			return nil, err
		}

		var transactionDetails []*walletrpc.TransactionDetails
		if in.MinedTransactions != nil {
			transactionDetails = append(transactionDetails, in.MinedTransactions.Transactions...)
		}
		if in.UnminedTransactions != nil {
			transactionDetails = append(transactionDetails, in.UnminedTransactions...)
		}

		txs, err := processTransactions(transactionDetails)
		if err != nil {
			return nil, err
		}

		transactions = append(transactions, txs...)
	}

	// sort transactions by date (list newer first)
	sort.SliceStable(transactions, func(i1, i2 int) bool {
		return transactions[i1].Timestamp > transactions[i2].Timestamp
	})

	return transactions, nil
}

func (c *WalletRPCClient) GetTransaction(transactionHash string) (*walletcore.TransactionDetails, error) {
	ctx := context.Background()
	hash, err := chainhash.NewHashFromStr(transactionHash)
	if err != nil {
		return nil, fmt.Errorf("invalid hash: %s\n%s", transactionHash, err.Error())
	}
	getTxRequest := &walletrpc.GetTransactionRequest{TransactionHash: hash[:]}
	getTxResponse, err := c.walletService.GetTransaction(ctx, getTxRequest)
	if isRpcErrorCode(err, codes.NotFound) {
		return nil, fmt.Errorf("transaction not found")
	} else if err != nil {
		return nil, err
	}

	decodedTx, err := txhelper.DecodeTransaction(hash, getTxResponse.GetTransaction().GetTransaction(), c.activeNet.Params, c.AddressInfo)
	if err != nil {
		return nil, err
	}

	transaction, err := processTransaction(getTxResponse.GetTransaction())
	if err != nil {
		return nil, err
	}
	transaction.Fee, transaction.FeeRate, transaction.Size = dcrutil.Amount(decodedTx.Fee), dcrutil.Amount(decodedTx.FeeRate), decodedTx.Size

	var blockHeight int32 = -1
	if getTxResponse.BlockHash != nil {
		blockInfo, err := c.walletService.BlockInfo(ctx, &walletrpc.BlockInfoRequest{BlockHash: getTxResponse.BlockHash})
		if err == nil {
			blockHeight = blockInfo.BlockHeight
		}
	}

	return &walletcore.TransactionDetails{
		BlockHeight:   blockHeight,
		Confirmations: getTxResponse.GetConfirmations(),
		Transaction:   transaction,
		Inputs:        decodedTx.Inputs,
		Outputs:       decodedTx.Outputs,
	}, nil
}

func (c *WalletRPCClient) StakeInfo(ctx context.Context) (*walletcore.StakeInfo, error) {
	stakeInfo, err := c.walletService.StakeInfo(ctx, &walletrpc.StakeInfoRequest{})
	if err != nil {
		return nil, fmt.Errorf("error getting stake info: %s", err.Error())
	}

	return &walletcore.StakeInfo{
		AllMempoolTix: stakeInfo.AllMempoolTix,
		Expired:       stakeInfo.Expired,
		Immature:      stakeInfo.Immature,
		Live:          stakeInfo.Live,
		Missed:        stakeInfo.Missed,
		OwnMempoolTix: stakeInfo.OwnMempoolTix,
		PoolSize:      stakeInfo.PoolSize,
		Revoked:       stakeInfo.Revoked,
		TotalSubsidy:  stakeInfo.TotalSubsidy,
		Unspent:       stakeInfo.Unspent,
		Voted:         stakeInfo.Voted,
	}, nil
}

func (c *WalletRPCClient) TicketPrice(ctx context.Context) (int64, error) {
	ticketPrice, err := c.walletService.TicketPrice(ctx, &walletrpc.TicketPriceRequest{})
	if err != nil {
		return 0, err
	}

	return ticketPrice.TicketPrice, nil
}

func (c *WalletRPCClient) PurchaseTicket(ctx context.Context, request dcrlibwallet.PurchaseTicketsRequest) ([]string, error) {
	ticketPrice, err := c.TicketPrice(ctx)
	if err != nil {
		return nil, fmt.Errorf("could not determine ticket price: %s", err.Error())
	}

	balance, err := c.AccountBalance(request.Account, int32(request.RequiredConfirmations))
	if err != nil {
		return nil, fmt.Errorf("could not fetch account balance: %s", err.Error())
	}

	totalTicketPrice := dcrutil.Amount(ticketPrice * int64(request.NumTickets))
	if balance.Spendable < totalTicketPrice {
		return nil, fmt.Errorf("insufficient funds: spendable account balance (%s) is less than ticket purchase cost %s",
			balance.Spendable, totalTicketPrice)
	}

	response, err := c.walletService.PurchaseTickets(ctx, &walletrpc.PurchaseTicketsRequest{
		Account:               request.Account,
		Expiry:                request.Expiry,
		NumTickets:            request.NumTickets,
		Passphrase:            request.Passphrase,
		PoolAddress:           request.PoolAddress,
		PoolFees:              request.PoolFees,
		RequiredConfirmations: request.RequiredConfirmations,
		SpendLimit:            ticketPrice,
		TicketAddress:         request.TicketAddress,
		TicketFee:             request.TicketFee,
		TxFee:                 request.TxFee,
	})
	if err != nil {
		return nil, fmt.Errorf("could not complete ticket(s) purchase, encountered an error:\n%s", err.Error())
	}
	ticketHashes := make([]string, len(response.GetTicketHashes()))
	for i, ticketHash := range response.GetTicketHashes() {
		hash, err := chainhash.NewHash(ticketHash)
		if err != nil {
			return ticketHashes, fmt.Errorf("encountered an error while processing purchased ticket(s):\n%s", err.Error())
		}
		ticketHashes[i] = hash.String()
	}
	return ticketHashes, nil
}
