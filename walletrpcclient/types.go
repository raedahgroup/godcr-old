package walletrpcclient

import (
	"github.com/decred/dcrd/dcrutil"
	"github.com/decred/dcrd/wire"
)

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
	Amount          int64  `json:"amount"`
	AmountString    string `json:"amount_string"`
	PkScript        []byte `json:"-"`
	AmountSum       string `json:"amount_sum"`
}

type Transaction struct {
	Hash          string               `json:"hash"`
	Type          string               `json:"type"`
	Amount        dcrutil.Amount       `json:"amount"`
	Fee           dcrutil.Amount       `json:"fee"`
	Rate          dcrutil.Amount       `json:"rate,omitempty"`
	Direction     TransactionDirection `json:"direction"`
	Testnet       bool                 `json:"testnet"`
	Timestamp     int64                `json:"timestamp"`
	FormattedTime string               `json:"formatted_time"`
	Size          int                  `json:"size"`
}

type TxInput struct {
	Amount           dcrutil.Amount `json:"value"`
	PreviousOutpoint wire.OutPoint  `json:"previousOutpoint"`
}

type TxOutput struct {
	Address     string         `json:"address"`
	Internal    bool           `json:"internal"`
	Value       dcrutil.Amount `json:"value"`
	ScriptClass string         `json:"scriptClass"`
}

type TransactionDetails struct {
	BlockHash     string     `json:"blockHash"`
	Confirmations int32      `json:"confirmations"`
	Inputs        []TxInput  `json:"inputs"`
	Outputs       []TxOutput `json:"outputs"`
	*Transaction
}
