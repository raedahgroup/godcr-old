package walletrpcclient

import "github.com/decred/dcrd/dcrutil"

type ReceiveResult struct {
	Address string
}

type AccountBalanceResult struct {
	AccountName     string         `json:"account_name"`
	AccountNumber   uint32         `json:"account_number"`
	Total           dcrutil.Amount `json:"total"`
	Spendable       dcrutil.Amount `json:"spendable"`
	LockedByTickets dcrutil.Amount `json:"locked_by_tickets"`
	VotingAuthority dcrutil.Amount `json:"voting_authority"`
	Unconfirmed     dcrutil.Amount `json:"unconfirmed"`
}

type SendResult struct {
	TransactionHash string `json:"transaction_hash"`
}

type UnspentOutputsResult struct {
	TransactionHash string `json:"transaction_hash"`
	OutputIndex     uint32 `json:"output_index"`
	ReceiveTime     int64  `json:"receive_time"`
	FromCoinbase    bool   `json:"from_coinbase"`
	Tree            int32  `json:"tree"`
	Amount          int64  `json:"amount"`
	PkScript        []byte `json:"-"`
	AmountSum       int64  `json:"amount_sum"`
}

type transactionInput struct {
	Index           uint32 `json:"index"`
	PreviousAccount uint32 `json:"previous_account"`
	PreviousAmount  int64  `json:"previous_amount"`
}

type transactionOutput struct {
	Index        uint32 `json:"index"`
	Account      uint32 `json:"account"`
	Internal     bool   `json:"internal"`
	Amount       int64  `json:"amount"`
	Address      string `json:"address"`
	OutputScript []byte `json:"-"`
}

type transaction struct {
	Input           []*transactionInput  `json:"input"`
	Output          []*transactionOutput `json:"output"`
	Hash            []byte               `json:"-"`
	Transaction     []byte               `json:"-"`
	Timestamp       int64                `json:"timestamp"`
	Fee             int64                `json:"fee"`
	TransactionType int                  `json:"transaction_type"`
}

type TransactionResult struct {
	TransactionDetails *transaction `json:"transaction"`
	Confirmations      int32        `json:"confirmations"`
	BlockHash          []byte       `json:"block_hash"`
}
