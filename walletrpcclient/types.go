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
	OutputKey       string `json:"key"`
	TransactionHash string `json:"transaction_hash"`
	OutputIndex     uint32 `json:"output_index"`
	ReceiveTime     int64  `json:"receive_time"`
	FromCoinbase    bool   `json:"from_coinbase"`
	Tree            int32  `json:"tree"`
	Amount			int64  `json:"amount"`
	AmountString    string  `json:"amount_string"`
	PkScript        []byte `json:"-"`
	AmountSum       string `json:"amount_sum"`
}
type FetchHeadersResult struct {
	HeadersCount            uint32 `json:"headers_count"`
	FirstNewBlockHash       []byte `json:"first_new_block_hash"`
	FirstNewBlockHeight     int32  `json:"first_new_block_height"`
	MainChainTipBlockHash   []byte `json:"main_chain_tip_block_hash"`
	MainChainTipBlockHeight int32  `json:"main_chain_tip_block_height"`
}
type TransactionInput struct {
	Index           uint32  `json:"index"`
	PreviousAccount uint32  `json:"previous_account"`
	PreviousAmount  float64 `json:"previous_amount"`
}
type TransactionOutput struct {
	Index        uint32  `json:"index"`
	Account      uint32  `json:"account"`
	Internal     bool    `json:"internal"`
	Amount       float64 `json:"amount"`
	Address      string  `json:"address"`
	OutputScript []byte  `json:"output_script"`
}
type TransactionDetails struct {
	Hash            string               `json:"hash"`
	Transaction     []byte               `json:"-"`
	Debits          []*TransactionInput  `json:"debits"`
	Credits         []*TransactionOutput `json:"credits"`
	Fee             int64                `json:"fee"`
	Timestamp       int64                `json:"timestamp"`
	TransactionType int                  `json:"transaction_type"`
}

type BlockDetails struct {
	Hash           string                `json:"hash"`
	Height         int32                 `json:"height"`
	Timestamp      int64                 `json:"timestamp"`
	ApprovesParent bool                  `json:"approves_parent"`
	Transactions   []*TransactionDetails `json:"transactions"`
}

type TransactionSummary struct {
	Hash            string  `json:"hash"`
	TransactionType string  `json:"transaction_type"`
	Amount          float64 `json:"amount"`
	Index           uint32  `json:"index"`
	PreviousAccount uint32  `json:"previous_account"`
	PreviousAmount  float64 `json:"previous_amount"`
	Account         uint32  `json:"account"`
	Internal        bool    `json:"internal"`
	Address         string  `json:"address"`
	OutputScript    []byte  `json:"output_script"`
}

type GetTransactionsResult struct {
	MinedTransactions   *BlockDetails         `json:"mined_transactions"`
	UnminedTransactions []*TransactionDetails `json:"unmined_transactions"`
	Summary             []*TransactionSummary `json:"summary"`
}
