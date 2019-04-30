package walletcore

import (
	"fmt"
	"strings"
	
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
	Name             string   `json:"name"`
	Number           uint32   `json:"number"`
	Balance          *Balance `json:"balance"`
	ExternalKeyCount int32    `json:"external_key_count"`
	InternalKeyCount int32    `json:"internal_key_count"`
	ImportedKeyCount int32    `json:"imported_key_count"`
}

func (account *Account) String() string {
	return fmt.Sprintf("%s [%s]", account.Name, account.Balance.Spendable.String())
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

// ConnectionInfo holds connection information for the wallet
type ConnectionInfo struct {
	NetworkType    string `json:"networkType"`
	PeersConnected int32  `json:"peersConnected"`
	TotalBalance   string `json:"totalBalance"`
	LatestBlock    uint32 `json:"latestBlock"`
}

type Transaction struct {
	*txhelper.Transaction
	// Following additional properties are not constant but change with time.
	// Need to update these fields before returning the tx to the caller.
	Status        string `json:"status"`
	Confirmations int32  `json:"confirmations"`
	ShortTime     string `json:"short_time"`
	LongTime      string `json:"long_time"`
}

func (tx *Transaction) WalletAccountForTx() string {
	var accountNames []string
	isInArray := func(accountName string) bool {
		for _, name := range accountNames {
			if name == accountName {
				return true
			}
		}
		return false
	}
	if tx.Direction == txhelper.TransactionDirectionReceived {
		for _, output := range tx.Outputs {
			if output.AccountNumber == -1 || isInArray(output.AccountName) {
				continue
			}
			accountNames = append(accountNames, output.AccountName)
		}
	} else {
		for _, input := range tx.Inputs {
			if input.AccountNumber == -1 || isInArray(input.AccountName) {
				continue
			}
			accountNames = append(accountNames, input.AccountName)
		}
	}
	return strings.Join(accountNames, ", ")
}