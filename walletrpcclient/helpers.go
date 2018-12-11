package walletrpcclient

import (
	"strings"
	"time"

	"github.com/decred/dcrd/chaincfg/chainhash"
	"github.com/decred/dcrd/dcrutil"
	pb "github.com/decred/dcrwallet/rpc/walletrpc"
)

func (c *Client) processTransactions(transactionDetails []*pb.TransactionDetails) ([]*Transaction, error) {
	transactions := make([]*Transaction, 0, len(transactionDetails))
	
	for _, txDetail := range transactionDetails {
		// sum credits
		var transactionAmount int64
		var isTestnet bool
		for _, output := range txDetail.Credits {
			isTestnet = strings.HasPrefix(output.Address, "T")
			if !output.Internal {
				transactionAmount += output.Amount
			}
		}
	
		hash, err := chainhash.NewHash(txDetail.Hash)
		if err != nil {
			return nil, err
		}

	
		tx := &Transaction{
			Hash:      hash.String(),
			Amount:    dcrutil.Amount(transactionAmount).ToCoin(),
			Fee: 	dcrutil.Amount(txDetail.Fee).ToCoin(),
			Type: txDetail.TransactionType.String(),
			IsTestnet: isTestnet,
			Timestamp: txDetail.Timestamp,
			FormattedTime: time.Unix(txDetail.Timestamp, 0).Format("Mon Jan 2, 2006 3:04PM"),
		}

		transactions = append(transactions, tx)
	}
	
	return transactions, nil
}

// func transactionAmountAndType() (amount uint64, type)

// received or sent?
// amountDifference := outputAmounts - inputAmounts
// if amountDifference < 0 && (float64(transaction.Fee) == math.Abs(float64(amountDifference))) {
// 	//Transfered
// 	direction = 2
// 	amount = int64(transaction.Fee)
// } else if amountDifference > 0 {
// 	//Received
// 	direction = 1
// 	amount = outputAmounts
// } else {
// 	//Sent
// 	direction = 0
// 	amount = inputAmounts
// 	amount -= outputAmounts
// 	amount -= int64(transaction.Fee)
// }