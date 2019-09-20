package core

import (
	"fmt"
	"github.com/decred/dcrd/dcrutil/v2"
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

