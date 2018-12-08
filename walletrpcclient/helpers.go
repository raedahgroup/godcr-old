package walletrpcclient

import (
	"time"

	"github.com/decred/dcrd/chaincfg/chainhash"
	pb "github.com/decred/dcrwallet/rpc/walletrpc"
)

func AtomsToCoin(amount int64) float64 {
	return float64(amount) / float64(100000000)
}

func getTransactionDetails(blockDetails *pb.BlockDetails) *TransactionSummary {
	hash, _ := chainhash.NewHash(blockDetails.Hash)

	summary := &TransactionSummary{
		Hash:      hash.String(),
		Timestamp: blockDetails.Timestamp,
		HumanTime: time.Unix(blockDetails.Timestamp, 0).Format("Mon Jan _2 15:04:05 2006"),
	}

	// get total amount
	for _, v := range blockDetails.Transactions {
		for _, k := range v.Credits {
			summary.Total += AtomsToCoin(k.Amount)
		}
	}
	return summary
}
