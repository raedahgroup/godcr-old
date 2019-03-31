package walletcore

import (
	"fmt"
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

func (balance *Balance) String() string {
	if balance.Total == balance.Spendable {
		return balance.Total.String()
	} else {
		return fmt.Sprintf("Total %s (Spendable %s)", balance.Total.String(), balance.Spendable.String())
	}
}

type Account struct {
	Name    string   `json:"name"`
	Number  uint32   `json:"number"`
	Balance *Balance `json:"balance"`
}

func (account *Account) String() string {
	return fmt.Sprintf("%s - %s", account.Name, account.Balance.String())
}

type UnspentOutput struct {
	OutputKey       string         `json:"key"`
	TransactionHash string         `json:"transaction_hash"`
	OutputIndex     uint32         `json:"output_index"`
	Tree            int32          `json:"tree"`
	ReceiveTime     int64          `json:"receive_time"`
	Amount          dcrutil.Amount `json:"amount"`
	Address         string         `json:"address"`
	Confirmations   int32          `json:"confirmations"`
}

type Transaction struct {
	Hash          string                        `json:"hash"`
	Type          string                        `json:"type"`
	Amount        string                `json:"amount"`
	Fee           string                `json:"fee"`
	FeeRate       dcrutil.Amount                `json:"rate,omitempty"`
	Direction     txhelper.TransactionDirection `json:"direction"`
	Timestamp     int64                         `json:"timestamp"`
	FormattedTime string                        `json:"formatted_time"`
	Size          int                           `json:"size"`
}

type TransactionDetails struct {
	BlockHeight   int32                     `json:"blockHeight"`
	Confirmations int32                     `json:"confirmations"`
	Inputs        []*txhelper.DecodedInput  `json:"inputs"`
	Outputs       []*txhelper.DecodedOutput `json:"outputs"`
	*Transaction
}

// StakeInfo holds ticket information summary related to the wallet.
type StakeInfo struct {
	// Stake info related to the wallet
	Expired       uint32 `json:"expired"`
	Immature      uint32 `json:"immature"`
	Live          uint32 `json:"live"`
	Missed        uint32 `json:"missed"`
	OwnMempoolTix uint32 `json:"ownMempoolTix"`
	Revoked       uint32 `json:"revoked"`
	Unspent       uint32 `json:"unspent"`
	Voted         uint32 `json:"voted"`

	// General blockchain stake info
	AllMempoolTix uint32 `json:"allMempoolTix"`
	PoolSize      uint32 `json:"poolSize"`
	TotalSubsidy  string `json:"totalSubsidy"`
}
