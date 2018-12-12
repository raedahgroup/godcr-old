package walletrpcclient

import (
	"math"
	"strings"
	"time"

	"github.com/decred/dcrd/chaincfg/chainhash"
	"github.com/decred/dcrd/dcrutil"
	pb "github.com/decred/dcrwallet/rpc/walletrpc"
)

const (
	// TransactionDirectionSent for transactions sent to external address(es) from wallet
	TransactionDirectionSent = "Sent"

	// TransactionDirectionReceived for transactions received from external address(es) into wallet
	TransactionDirectionReceived = "Received"

	// TransactionDirectionTransferred for transactions sent from wallet to internal address(es)
	TransactionDirectionTransferred = "Transferred"
)

func (c *Client) processTransactions(transactionDetails []*pb.TransactionDetails) ([]*Transaction, error) {
	transactions := make([]*Transaction, 0, len(transactionDetails))
	
	for _, txDetail := range transactionDetails {
		// use any of the addresses in inputs/outputs to determine if this is a testnet tx
		var isTestnet bool
		for _, output := range txDetail.Credits {
			isTestnet = strings.HasPrefix(output.Address, "T")
			break
		}
	
		hash, err := chainhash.NewHash(txDetail.Hash)
		if err != nil {
			return nil, err
		}

		amount, direction := transactionAmountAndDirection(txDetail)
	
		tx := &Transaction{
			Hash:      hash.String(),
			Amount:    dcrutil.Amount(amount).ToCoin(),
			Fee: 	dcrutil.Amount(txDetail.Fee).ToCoin(),
			Type: txDetail.TransactionType.String(),
			Direction: direction,
			IsTestnet: isTestnet,
			Timestamp: txDetail.Timestamp,
			FormattedTime: time.Unix(txDetail.Timestamp, 0).Format("Mon Jan 2, 2006 3:04PM"),
		}

		transactions = append(transactions, tx)
	}
	
	return transactions, nil
}

func transactionAmountAndDirection(txDetail *pb.TransactionDetails) (int64, string) {
	var outputAmounts int64
	for _, credit := range txDetail.Credits {
		outputAmounts += int64(credit.Amount)
	}

	var inputAmounts int64
	for _, debit := range txDetail.Debits {
		inputAmounts += int64(debit.PreviousAmount)
	}

	var amount int64
	var direction string

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
