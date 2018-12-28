package walletcore

import (
	"github.com/decred/dcrd/dcrutil"
)

var (
	transactionDirectionNames = []string{"Sent", "Received", "Transferred", "Unclear"}
)

const (
	// TransactionDirectionSent for transactions sent to external address(es) from wallet
	TransactionDirectionSent TransactionDirection = iota

	// TransactionDirectionReceived for transactions received from external address(es) into wallet
	TransactionDirectionReceived

	// TransactionDirectionTransferred for transactions sent from wallet to internal address(es)
	TransactionDirectionTransferred

	// TransactionDirectionUnclear for unrecognized transaction directions
	TransactionDirectionUnclear
)

type TransactionDirection int8

func (direction TransactionDirection) String() string {
	if direction <= TransactionDirectionUnclear {
		return transactionDirectionNames[direction]
	} else {
		return transactionDirectionNames[TransactionDirectionUnclear]
	}
}

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
	Hash          string               `json:"hash"`
	Type          string               `json:"type"`
	Amount        string               `json:"amount"`
	Fee           string               `json:"fee"`
	Direction     TransactionDirection `json:"direction"`
	Timestamp     int64                `json:"timestamp"`
	FormattedTime string               `json:"formatted_time"`
}
