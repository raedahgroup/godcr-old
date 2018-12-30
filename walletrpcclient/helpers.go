package walletrpcclient

import (
	"bytes"
	"github.com/decred/dcrd/chaincfg"
	"github.com/decred/dcrd/txscript"
	"github.com/decred/dcrd/wire"
	"math"
	"time"

	"github.com/decred/dcrd/chaincfg/chainhash"
	"github.com/decred/dcrd/dcrutil"
	"github.com/decred/dcrwallet/netparams"
	pb "github.com/decred/dcrwallet/rpc/walletrpc"
)

type TransactionDirection int8

const (
	// TransactionDirectionSent for transactions sent to external address(es) from wallet
	TransactionDirectionSent TransactionDirection = iota

	// TransactionDirectionReceived for transactions received from external address(es) into wallet
	TransactionDirectionReceived

	// TransactionDirectionTransferred for transactions sent from wallet to internal address(es)
	TransactionDirectionTransferred
)

func (d TransactionDirection) String() string {
	switch d {
	case 0:
		return "Sent"

	case 1:
		return "Received"

	case 2:
		return "Transferred"
	}

	return ""
}

func (c *Client) processTransactions(transactionDetails []*pb.TransactionDetails) ([]*Transaction, error) {
	transactions := make([]*Transaction, 0, len(transactionDetails))

	for _, txDetail := range transactionDetails {
		// use any of the addresses in inputs/outputs to determine if this is a testnet tx
		var isTestnet bool
		for _, output := range txDetail.Credits {
			isMainNet, err := addressIsForNet(output.GetAddress(), netparams.MainNetParams.Params)
			if err != nil {
				continue
			}
			isTestnet = !isMainNet
			break
		}

		tx, err := processTransaction(txDetail, isTestnet)
		if err != nil {
			return nil, err
		}

		transactions = append(transactions, tx)
	}

	return transactions, nil
}

func processTransaction(txDetail *pb.TransactionDetails, isTestnet bool) (*Transaction, error) {
	hash, err := chainhash.NewHash(txDetail.Hash)
		if err != nil {
			return nil, err
		}

		amount, direction := transactionAmountAndDirection(txDetail)

		tx := &Transaction{
			Hash:          hash.String(),
			Amount:        dcrutil.Amount(amount),
			Fee:           dcrutil.Amount(txDetail.Fee),
			Type:          txDetail.TransactionType.String(),
			Direction:     direction,
			Testnet:       isTestnet,
			Timestamp:     txDetail.Timestamp,
			FormattedTime: time.Unix(txDetail.Timestamp, 0).Format("Mon Jan 2, 2006 3:04PM"),
		}
		return tx, nil
}

func addressIsForNet(address string, net *chaincfg.Params) (bool, error) {
	addr, err := dcrutil.DecodeAddress(address)
	if err != nil {
		return false, err
	}
	return addr.IsForNet(net), nil
}

func transactionAmountAndDirection(txDetail *pb.TransactionDetails) (int64, TransactionDirection) {
	var outputAmounts int64
	for _, credit := range txDetail.Credits {
		outputAmounts += int64(credit.Amount)
	}

	var inputAmounts int64
	for _, debit := range txDetail.Debits {
		inputAmounts += int64(debit.PreviousAmount)
	}

	var amount int64
	var direction TransactionDirection

	if txDetail.TransactionType == pb.TransactionDetails_REGULAR {
		amountDifference := outputAmounts - inputAmounts
		if amountDifference < 0 && (float64(txDetail.Fee) == math.Abs(float64(amountDifference))) {
			// transfered internally, the only real amount spent was transaction fee
			direction = TransactionDirectionTransferred
			amount = int64(txDetail.Fee)
		} else if amountDifference > 0 {
			// received
			direction = TransactionDirectionReceived

			for _, credit := range txDetail.Credits {
				amount += int64(credit.Amount)
			}
		} else {
			// sent
			direction = TransactionDirectionSent

			for _, debit := range txDetail.Debits {
				amount += int64(debit.PreviousAmount)
			}
			for _, credit := range txDetail.Credits {
				amount -= int64(credit.Amount)
			}
			amount -= int64(txDetail.Fee)
		}
	}

	return amount, direction
}

func inputsFromMsgTxIn(txIn []*wire.TxIn) []TxInput {
	txInputs := make([]TxInput, len(txIn))
	for i, input := range txIn {
		txInputs[i] = TxInput{
			Amount: dcrutil.Amount(input.ValueIn),
			PreviousOutpoint: input.PreviousOutPoint.String(),
		}
	}
	return txInputs
}

func outputsFromMsgTxOut(txOut []*wire.TxOut, walletCredits []*pb.TransactionDetails_Output, chainParams *chaincfg.Params) ([]TxOutput, error) {
	txOutputs := make([]TxOutput, len(txOut))
	for i, output := range txOut {
		_, addrs, _, err := txscript.ExtractPkScriptAddrs(output.Version, output.PkScript, chainParams)
		if err != nil {
			return nil, err
		}
		txOutputs[i] = TxOutput{Value: dcrutil.Amount(output.Value), Address: addrs[0].String()}
		for _, credit := range walletCredits {
			if bytes.Equal(output.PkScript, credit.GetOutputScript()) {
				txOutputs[i].Internal = credit.GetInternal()
			}
		}
	}
	return txOutputs, nil
}
