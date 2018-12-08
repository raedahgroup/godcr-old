package walletrpcclient

import (
	"github.com/decred/dcrd/chaincfg/chainhash"
	pb "github.com/decred/dcrwallet/rpc/walletrpc"
)

func AtomsToCoin(amount int64) float64 {
	return float64(amount) / float64(100000000)
}

func getOutputs(credits []*pb.TransactionDetails_Output) []*TransactionOutput {
	outputs := make([]*TransactionOutput, len(credits))

	for i, v := range credits {
		output := &TransactionOutput{
			Index:        v.Index,
			Account:      v.Account,
			Internal:     v.Internal,
			Amount:       AtomsToCoin(v.Amount),
			Address:      v.Address,
			OutputScript: v.OutputScript,
		}
		outputs[i] = output
	}

	return outputs
}

func getInputs(debits []*pb.TransactionDetails_Input) []*TransactionInput {
	inputs := make([]*TransactionInput, len(debits))

	for i, v := range debits {
		input := &TransactionInput{
			Index:           v.Index,
			PreviousAccount: v.PreviousAccount,
			PreviousAmount:  AtomsToCoin(v.PreviousAmount),
		}
		inputs[i] = input
	}

	return inputs
}

func getTransactionDetails(txDetails *pb.TransactionDetails) *TransactionDetails {
	hash, _ := chainhash.NewHash(txDetails.Hash)

	tx := &TransactionDetails{
		Hash:            hash.String(),
		Transaction:     txDetails.Transaction,
		Debits:          getInputs(txDetails.Debits),
		Credits:         getOutputs(txDetails.Credits),
		Fee:             txDetails.Fee,
		Timestamp:       txDetails.Timestamp,
		TransactionType: int(txDetails.TransactionType),
	}

	return tx
}

func getTransactionsDetails(txDetails []*pb.TransactionDetails) []*TransactionDetails {
	txns := make([]*TransactionDetails, len(txDetails))

	for i, v := range txDetails {
		txns[i] = getTransactionDetails(v)
	}

	return txns
}

func getBlockDetails(blockDetails *pb.BlockDetails) *BlockDetails {
	hash, _ := chainhash.NewHash(blockDetails.Hash)

	b := &BlockDetails{
		Hash:           hash.String(),
		Height:         blockDetails.Height,
		Timestamp:      blockDetails.Timestamp,
		ApprovesParent: blockDetails.ApprovesParent,
		Transactions:   getTransactionsDetails(blockDetails.Transactions),
	}

	return b
}

func getSummary(blockDetails *BlockDetails) []*TransactionSummary {
	summaries := []*TransactionSummary{}

	for _, v := range blockDetails.Transactions {
		for _, j := range v.Credits {
			s := &TransactionSummary{
				Hash:            v.Hash,
				TransactionType: "credit",
				Amount:          j.Amount,
				Index:           j.Index,
				Account:         j.Account,
				Address:         j.Address,
				Internal:        j.Internal,
				OutputScript:    j.OutputScript,
			}

			summaries = append(summaries, s)
		}

		for _, j := range v.Debits {
			s := &TransactionSummary{
				TransactionType: "debit",
				Hash:            v.Hash,
				Index:           j.Index,
				PreviousAccount: j.PreviousAccount,
				PreviousAmount:  j.PreviousAmount,
			}
			summaries = append(summaries, s)
		}
	}

	return summaries
}
