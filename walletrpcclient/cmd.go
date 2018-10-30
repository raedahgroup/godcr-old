package walletrpcclient

import (
	"context"
	"fmt"

	"github.com/Baozisoftware/qrcode-terminal-go"
	"github.com/decred/dcrd/dcrutil"
	"github.com/decred/dcrd/txscript"
	pb "github.com/decred/dcrwallet/rpc/walletrpc"
)

func (c *Client) cmdSendTransaction(ctx context.Context, sourceAccount uint32, destAddr string, amount int64, passphrase string) (*Response, error) {

	// decode destination address
	addr, err := dcrutil.DecodeAddress(destAddr)
	if err != nil {
		return nil, fmt.Errorf("error decoding destination address: %s", err.Error())
	}

	// get script
	pkScript, err := txscript.PayToAddrScript(addr)
	if err != nil {
		return nil, err
	}

	// construct transaction
	cReq := &pb.ConstructTransactionRequest{
		SourceAccount: sourceAccount,
		NonChangeOutputs: []*pb.ConstructTransactionRequest_Output{{
			Destination: &pb.ConstructTransactionRequest_OutputDestination{
				Script:        pkScript,
				ScriptVersion: 0,
			},
			Amount: amount,
		}},
	}

	cRes, err := c.wc.ConstructTransaction(ctx, cReq)
	if err != nil {
		return nil, fmt.Errorf("Error constructing transaction: %s", err.Error())
	}

	// sign transaction
	sReq := &pb.SignTransactionRequest{
		Passphrase:            []byte(passphrase),
		SerializedTransaction: cRes.UnsignedTransaction,
	}

	sRes, err := c.wc.SignTransaction(ctx, sReq)
	if err != nil {
		return nil, fmt.Errorf("error signing transaction: %s", err.Error())
	}

	// publish transaction
	pReq := &pb.PublishTransactionRequest{
		SignedTransaction: sRes.Transaction,
	}

	pRes, err := c.wc.PublishTransaction(ctx, pReq)
	if err != nil {
		return nil, fmt.Errorf("error publishing transaction: %s", err.Error())
	}

	res := &Response{
		Columns: []string{"Result", "Hash"},
	}

	resultRow := []interface{}{
		"Successfull",
		string(pRes.TransactionHash),
	}

	res.Result = [][]interface{}{resultRow}
	return res, nil
}

func (c *Client) cmdBalance(ctx context.Context) (*Response, error) {
	accounts, err := c.wc.Accounts(ctx, &pb.AccountsRequest{})
	if err != nil {
		return nil, fmt.Errorf("error fetching accounts. err: %s", err.Error())
	}

	balances := make([][]interface{}, len(accounts.Accounts))
	for i, v := range accounts.Accounts {
		req := &pb.BalanceRequest{
			AccountNumber:         v.AccountNumber,
			RequiredConfirmations: 0,
		}

		res, err := c.wc.Balance(ctx, req)
		if err != nil {
			return nil, fmt.Errorf("error fetching balance for account: %d :%s", v.AccountNumber, err.Error())
		}

		balances[i] = []interface{}{
			v.AccountName,
			dcrutil.Amount(res.Total),
			dcrutil.Amount(res.Spendable),
			dcrutil.Amount(res.LockedByTickets),
			dcrutil.Amount(res.VotingAuthority),
			dcrutil.Amount(res.Unconfirmed),
		}
	}

	balanceColumns := []string{
		"Account",
		"Total",
		"Spendable",
		"Locked By Tickets",
		"Voting Authority",
		"Unconfirmed",
	}

	res := &Response{
		Columns: balanceColumns,
		Result:  balances,
	}

	return res, nil
}

func (c *Client) cmdReceive(ctx context.Context, accountNumber uint32) (*Response, error) {
	req := &pb.NextAddressRequest{
		Account:   accountNumber,
		GapPolicy: pb.NextAddressRequest_GAP_POLICY_WRAP,
		Kind:      pb.NextAddressRequest_BIP0044_EXTERNAL,
	}

	r, err := c.wc.NextAddress(ctx, req)
	if err != nil {
		return nil, fmt.Errorf("error fetching receive address: %s", err.Error())
	}

	res := &Response{
		Columns: []string{
			"Address",
			"QR Code",
		},
		Result: [][]interface{}{
			[]interface{}{
				r.Address,
				"",
			},
		},
	}
	obj := qrcodeTerminal.New()
	obj.Get(r.Address).Print()

	return res, nil
}

// walletVersion fetches and returns version of wallet we are connected to
func (c *Client) cmdWalletVersion(ctx context.Context) (*Response, error) {
	r, err := c.vc.Version(ctx, &pb.VersionRequest{})
	if err != nil {
		return nil, err
	}

	res := &Response{
		Columns: []string{
			"Version",
		},
		Result: [][]interface{}{
			[]interface{}{r.VersionString},
		},
	}
	return res, nil
}

// listCommands lists all supported commands
// requires no parameter
func (c *Client) cmdListCommands(ctx context.Context) (*Response, error) {
	res := &Response{
		Columns: []string{"Command", "Description"},
	}
	for i, v := range c.commands {
		item := []interface{}{
			v,
			c.descriptions[i],
		}
		res.Result = append(res.Result, item)
	}
	return res, nil
}
