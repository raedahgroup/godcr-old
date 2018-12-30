package walletrpcclient

import (
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
			addr, err := dcrutil.DecodeAddress(output.Address)
			if err != nil {
				continue
			}

			isTestnet = !addr.IsForNet(netparams.MainNetParams.Params)
			break
		}

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

		transactions = append(transactions, tx)
	}

	return transactions, nil
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
