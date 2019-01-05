package walletcore

import (
	"github.com/decred/dcrd/dcrutil"
	"github.com/raedahgroup/dcrlibwallet/txhelper"
)

type Balance struct {
	Total           dcrutil.Amount `json:"total"`
	Spendable       dcrutil.Amount `json:"spendable"`
	LockedByTickets dcrutil.Amount `json:"locked_by_tickets"`
	VotingAuthority dcrutil.Amount `json:"voting_authority"`
	Unconfirmed     dcrutil.Amount `json:"unconfirmed"`
}

type Account struct {
	Name    string   `json:"name"`
	Number  uint32   `json:"number"`
	Balance *Balance `json:"balance"`
}

type UnspentOutput struct {
	OutputKey       string         `json:"key"`
	TransactionHash string         `json:"transaction_hash"`
	OutputIndex     uint32         `json:"output_index"`
	Tree            int32          `json:"tree"`
	ReceiveTime     int64          `json:"receive_time"`
	Amount          dcrutil.Amount `json:"amount"`
}

type Transaction struct {
	Hash          string                        `json:"hash"`
	Type          string                        `json:"type"`
	Amount        dcrutil.Amount                `json:"amount"`
	AmountCoin    float64                       `json:"coin_amount"`
	Fee           dcrutil.Amount                `json:"fee"`
	FeeCoin       float64                       `json:"coin_fee"`
	FeeRate       dcrutil.Amount                `json:"rate,omitempty"`
	Direction     txhelper.TransactionDirection `json:"direction"`
	Timestamp     int64                         `json:"timestamp"`
	FormattedTime string                        `json:"formatted_time"`
	Size          int                           `json:"size"`
}

type TxInput struct {
	Amount           dcrutil.Amount `json:"value"`
	PreviousOutpoint string         `json:"previousOutpoint"`
}

type TxOutput struct {
	Address  string         `json:"address"`
	Internal bool           `json:"internal"`
	Value    dcrutil.Amount `json:"value"`
}

type TransactionDetails struct {
	BlockHeight   int32                     `json:"blockHeight"`
	Confirmations int32                     `json:"confirmations"`
	Inputs        []*txhelper.DecodedInput  `json:"inputs"`
	Outputs       []*txhelper.DecodedOutput `json:"outputs"`
	*Transaction
}

type Ticket struct {
	Ticket Transaction `json:"ticket"`
	Status string      `json:"status"`
}
