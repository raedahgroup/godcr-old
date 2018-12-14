package walletsource

import "github.com/decred/dcrd/dcrutil"

type BlockChainSyncListener struct {
	SyncStarted         func()
	SyncEnded           func(err error)
	OnHeadersFetched    func(percentageProgress int64)
	OnDiscoveredAddress func(state string)
	OnRescanningBlocks  func(percentageProgress int64)
}

type Balance struct {
	Total         dcrutil.Amount `json:"total"`
	Spendable     dcrutil.Amount `json:"spendable"`
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
	ReceiveTime     int64          `json:"receive_time"`
	Amount          dcrutil.Amount `json:"amount"`
}

type Transaction struct {
	Hash          string               `json:"hash"`
	Type          string               `json:"type"`
	Amount        float64              `json:"amount"`
	Fee           float64              `json:"fee"`
	Direction     TransactionDirection `json:"direction"`
	Testnet       bool                 `json:"testnet"`
	Timestamp     int64                `json:"timestamp"`
	FormattedTime string               `json:"formatted_time"`
}
