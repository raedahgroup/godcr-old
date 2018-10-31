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
