package walletrpcclient

import (
	"bytes"
	"context"
	"fmt"
	"io"
	"net"
	"sort"

	"github.com/decred/dcrd/chaincfg/chainhash"
	"github.com/decred/dcrd/dcrutil"
	"github.com/decred/dcrd/wire"
	"github.com/decred/dcrdata/txhelpers"
	"github.com/decred/dcrwallet/netparams"
	pb "github.com/decred/dcrwallet/rpc/walletrpc"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"

	"github.com/raedahgroup/godcr/walletrpcclient/walletcore"
)

type Client struct {
	walletServiceClient pb.WalletServiceClient
}

func New(rpcAddress, rpcCert string, noTLS, isTestnet bool) (*Client, error) {
	if rpcAddress == "" {
		rpcAddress = defaultDcrWalletRPCAddress(isTestnet)
	}

	c := &Client{}
	conn, err := c.connect(rpcAddress, rpcCert, noTLS)
	if err != nil {
		return nil, err
	}

	// register clients
	c.walletServiceClient = pb.NewWalletServiceClient(conn)

	return c, nil
}

func defaultDcrWalletRPCAddress(isTestnet bool) string {
	if isTestnet {
		return net.JoinHostPort("localhost", netparams.TestNet3Params.GRPCServerPort)
	} else {
		return net.JoinHostPort("localhost", netparams.MainNetParams.GRPCServerPort)
	}
}

func (c *Client) connect(rpcAddress, rpcCert string, noTLS bool) (*grpc.ClientConn, error) {
	var conn *grpc.ClientConn
	var err error

	if noTLS {
		conn, err = grpc.Dial(rpcAddress, grpc.WithInsecure())
		if err != nil {
			return nil, err
		}
	} else {
		creds, err := credentials.NewClientTLSFromFile(rpcCert, "")
		if err != nil {
			return nil, err
		}

		opts := []grpc.DialOption{
			grpc.WithTransportCredentials(
				creds,
			),
		}

		conn, err = grpc.Dial(rpcAddress, opts...)
		if err != nil {
			return nil, err
		}
	}

	return conn, nil
}

func (c *Client) SendFromAccount(amountInDCR float64, sourceAccount uint32, destinationAddress,
	passphrase string) (*SendResult, error) {

	// convert amount from float64 DCR to int64 Atom as required by dcrwallet ConstructTransaction implementation
	amountInAtom, err := dcrutil.NewAmount(amountInDCR)
	if err != nil {
		return nil, err
	}
	// type of amountInAtom is `dcrutil.Amount` which is an int64 alias
	amount := int64(amountInAtom)

	// construct transaction
	pkScript, err := walletcore.GetPKScript(destinationAddress)
	if err != nil {
		return nil, err
	}
	constructRequest := &pb.ConstructTransactionRequest{
		SourceAccount: sourceAccount,
		NonChangeOutputs: []*pb.ConstructTransactionRequest_Output{{
			Destination: &pb.ConstructTransactionRequest_OutputDestination{
				Script:        pkScript,
				ScriptVersion: 0,
			},
			Amount: amount,
		}},
	}

	constructResponse, err := c.walletServiceClient.ConstructTransaction(context.Background(), constructRequest)
	if err != nil {
		return nil, fmt.Errorf("error constructing transaction: %s", err.Error())
	}

	return c.signAndPublishTransaction(constructResponse.UnsignedTransaction, passphrase)
}

func (c *Client) SendFromUTXOs(utxoKeys []string, amountInDCR float64, sourceAccount uint32,
	destinationAddress, passphrase string) (*SendResult, error) {

	// convert amount to atoms
	amountInAtom, err := dcrutil.NewAmount(amountInDCR)
	amount := int64(amountInAtom)

	// fetch all utxos to extract details for the utxos selected by user
	req := &pb.UnspentOutputsRequest{
		Account:                  sourceAccount,
		TargetAmount:             0,
		RequiredConfirmations:    0,
		IncludeImmatureCoinbases: true,
	}
	stream, err := c.walletServiceClient.UnspentOutputs(context.Background(), req)
	if err != nil {
		return nil, err
	}

	// loop through utxo stream to find user selected utxos
	inputs := make([]*wire.TxIn, 0, len(utxoKeys))
	for {
		item, err := stream.Recv()
		if err == io.EOF {
			break
		}

		if err != nil {
			return nil, err
		}

		transactionHash, err := chainhash.NewHash(item.TransactionHash)
		if err != nil {
			return nil, fmt.Errorf("invalid transaction hash: %s", err.Error())
		}

		outputKey := fmt.Sprintf("%s:%v", transactionHash.String(), item.OutputIndex)
		useUtxo := false
		for _, key := range utxoKeys {
			if outputKey == key {
				useUtxo = true
			}
		}
		if !useUtxo {
			continue
		}

		outpoint := wire.NewOutPoint(transactionHash, item.OutputIndex, int8(item.Tree))
		input := wire.NewTxIn(outpoint, item.Amount, nil)
		inputs = append(inputs, input)

		if len(inputs) == len(utxoKeys) {
			break
		}
	}

	// generate address from sourceAccount to receive change
	receiveResult, err := c.Receive(sourceAccount)
	if err != nil {
		return nil, err
	}
	changeAddress := receiveResult.Address

	unsignedTx, err := walletcore.NewUnsignedTx(inputs, amount, destinationAddress, changeAddress)
	if err != nil {
		return nil, err
	}

	// serialize unsigned tx
	var txBuf bytes.Buffer
	txBuf.Grow(unsignedTx.SerializeSize())
	err = unsignedTx.Serialize(&txBuf)
	if err != nil {
		return nil, fmt.Errorf("error serializing transaction: %s", err.Error())
	}

	return c.signAndPublishTransaction(txBuf.Bytes(), passphrase)
}

func (c *Client) signAndPublishTransaction(serializedTx []byte, passphrase string) (*SendResult, error) {
	ctx := context.Background()

	// sign transaction
	signRequest := &pb.SignTransactionRequest{
		Passphrase:            []byte(passphrase),
		SerializedTransaction: serializedTx,
	}

	signResponse, err := c.walletServiceClient.SignTransaction(ctx, signRequest)
	if err != nil {
		return nil, fmt.Errorf("error signing transaction: %s", err.Error())
	}

	// publish transaction
	publishRequest := &pb.PublishTransactionRequest{
		SignedTransaction: signResponse.Transaction,
	}

	publishResponse, err := c.walletServiceClient.PublishTransaction(ctx, publishRequest)
	if err != nil {
		return nil, fmt.Errorf("error publishing transaction: %s", err.Error())
	}

	transactionHash, _ := chainhash.NewHash(publishResponse.TransactionHash)

	response := &SendResult{
		TransactionHash: transactionHash.String(),
	}

	return response, nil
}

func (c *Client) Balance() ([]*AccountBalanceResult, error) {
	ctx := context.Background()
	accounts, err := c.walletServiceClient.Accounts(ctx, &pb.AccountsRequest{})
	if err != nil {
		return nil, fmt.Errorf("error fetching accounts: %s", err.Error())
	}

	balanceResult := make([]*AccountBalanceResult, 0, len(accounts.Accounts))

	for _, v := range accounts.Accounts {
		accountBalance, err := c.SingleAccountBalance(v.AccountNumber, ctx)
		if err != nil {
			return nil, err
		}

		if v.AccountName == "imported" && accountBalance.Total == 0 {
			continue
		}

		accountBalance.AccountName = v.AccountName
		balanceResult = append(balanceResult, accountBalance)
	}

	return balanceResult, nil
}

func (c *Client) SingleAccountBalance(accountNumber uint32, ctx context.Context) (*AccountBalanceResult, error) {
	if ctx == nil {
		ctx = context.Background()
	}

	req := &pb.BalanceRequest{
		AccountNumber:         accountNumber,
		RequiredConfirmations: 0,
	}

	res, err := c.walletServiceClient.Balance(ctx, req)
	if err != nil {
		return nil, fmt.Errorf("error fetching balance for account: %d :%s", accountNumber, err.Error())
	}

	return &AccountBalanceResult{
		AccountNumber:   accountNumber,
		Total:           dcrutil.Amount(res.Total),
		Spendable:       dcrutil.Amount(res.Spendable),
		LockedByTickets: dcrutil.Amount(res.LockedByTickets),
		VotingAuthority: dcrutil.Amount(res.VotingAuthority),
		Unconfirmed:     dcrutil.Amount(res.Unconfirmed),
	}, nil
}

func (c *Client) Receive(accountNumber uint32) (*ReceiveResult, error) {
	ctx := context.Background()

	req := &pb.NextAddressRequest{
		Account:   accountNumber,
		GapPolicy: pb.NextAddressRequest_GAP_POLICY_WRAP,
		Kind:      pb.NextAddressRequest_BIP0044_EXTERNAL,
	}

	r, err := c.walletServiceClient.NextAddress(ctx, req)
	if err != nil {
		return nil, fmt.Errorf("error generating receive address: %s", err.Error())
	}

	res := &ReceiveResult{
		Address: r.Address,
	}

	return res, nil
}

func (c *Client) IsAddressValid(address string) (bool, error) {
	r, err := c.ValidateAddress(address)
	if err != nil {
		return false, err
	}

	return r.IsValid, nil
}

func (c *Client) ValidateAddress(address string) (*pb.ValidateAddressResponse, error) {
	req := &pb.ValidateAddressRequest{
		Address: address,
	}

	return c.walletServiceClient.ValidateAddress(context.Background(), req)
}

func (c *Client) NextAccount(accountName string, passphrase string) (uint32, error) {
	req := &pb.NextAccountRequest{
		AccountName: accountName,
		Passphrase:  []byte(passphrase),
	}

	r, err := c.walletServiceClient.NextAccount(context.Background(), req)
	if err != nil {
		return 0, err
	}

	return r.AccountNumber, nil
}

func (c *Client) AccountNumber(accountName string) (uint32, error) {
	req := &pb.AccountNumberRequest{
		AccountName: accountName,
	}

	r, err := c.walletServiceClient.AccountNumber(context.Background(), req)
	if err != nil {
		return 0, err
	}

	return r.AccountNumber, nil
}

func (c *Client) UnspentOutputs(account uint32, targetAmount int64) ([]*UnspentOutputsResult, error) {
	req := &pb.UnspentOutputsRequest{
		Account:                  account,
		TargetAmount:             targetAmount,
		RequiredConfirmations:    0,
		IncludeImmatureCoinbases: true,
	}

	stream, err := c.walletServiceClient.UnspentOutputs(context.Background(), req)
	if err != nil {
		return nil, err
	}

	outputs := []*UnspentOutputsResult{}

	for {
		item, err := stream.Recv()
		if err == io.EOF {
			break
		}

		if err != nil {
			return nil, err
		}

		transactionHash, _ := chainhash.NewHash(item.TransactionHash)

		outputItem := &UnspentOutputsResult{
			OutputKey:       fmt.Sprintf("%s:%v", transactionHash.String(), item.OutputIndex),
			TransactionHash: transactionHash.String(),
			OutputIndex:     item.OutputIndex,
			Amount:          item.Amount,
			AmountString:    dcrutil.Amount(item.Amount).String(),
			PkScript:        item.PkScript,
			AmountSum:       dcrutil.Amount(item.AmountSum).String(),
			ReceiveTime:     item.ReceiveTime,
			Tree:            item.Tree,
			FromCoinbase:    item.FromCoinbase,
		}
		outputs = append(outputs, outputItem)
	}

	return outputs, nil
}

func (c *Client) GetTransactions() ([]*Transaction, error) {
	req := &pb.GetTransactionsRequest{}

	stream, err := c.walletServiceClient.GetTransactions(context.Background(), req)
	if err != nil {
		return nil, err
	}

	var transactions []*Transaction

	for {
		in, err := stream.Recv()
		if err == io.EOF {
			break
		}
		if err != nil {
			return nil, err
		}

		var transactionDetails []*pb.TransactionDetails
		if in.MinedTransactions != nil {
			transactionDetails = append(transactionDetails, in.MinedTransactions.Transactions...)
		}
		if in.UnminedTransactions != nil {
			transactionDetails = append(transactionDetails, in.UnminedTransactions...)
		}

		txs, err := c.processTransactions(transactionDetails)
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

func (c *Client) GetTransaction(transactionHash string) (*GetTransactionResponse, error) {
	ctx := context.Background()

	hash, err := chainhash.NewHashFromStr(transactionHash)
	if err != nil {
		return nil, err
	}
	getTransactionRequest := &pb.GetTransactionRequest{TransactionHash: hash[:]}
	transactionResponse, err := c.walletServiceClient.GetTransaction(ctx, getTransactionRequest)
	if err != nil {
		return nil, err
	}

	transactionHex := fmt.Sprintf("%x", transactionResponse.GetTransaction().GetTransaction())
	msgTx, err := txhelpers.MsgTxFromHex(transactionHex)
	if err != nil {
		return nil, err
	}
	txFee, txFeeRate := txhelpers.TxFeeRate(msgTx)

	txInfos, err := c.processTransactions([]*pb.TransactionDetails{transactionResponse.GetTransaction()})
	if err != nil {
		return nil, err
	}
	transaction := txInfos[0]

	transaction.Fee, transaction.Rate, transaction.Size = txFee, txFeeRate, msgTx.SerializeSize()

	return &GetTransactionResponse{
		BlockHash:     fmt.Sprintf("%x", transactionResponse.GetBlockHash()),
		Confirmations: transactionResponse.GetConfirmations(),
		Transaction:   transaction,
	}, nil
}
