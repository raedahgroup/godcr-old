package nuklear

import (
	"reflect"
	"testing"

	"github.com/raedahgroup/godcr/nuklear/handlers"
)

func Test_getNavPages(t *testing.T) {
	tests := []struct {
		name string
		want []navPage
	}{
		{
			name: "nav pages",
			want: []navPage{
				{
					name:    "balance",
					label:   "Balance",
					handler: &handlers.BalanceHandler{},
				},
				{
					name:    "receive",
					label:   "Receive",
					handler: &handlers.ReceiveHandler{},
				},
				{
					name:    "send",
					label:   "Send",
					handler: &handlers.SendHandler{},
				},
				{
					name:    "history",
					label:   "History",
					handler: &handlers.TransactionsHandler{},
				},
				{
					name:    "staking",
					label:   "Staking",
					handler: &handlers.StakingHandler{},
				},
			},
		},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			if got := getNavPages(); !reflect.DeepEqual(got, test.want) {
				t.Errorf("getNavPages() = %v, want %v", got, test.want)
			}
		})
	}
}

func Test_getStandalonePages(t *testing.T) {
	tests := []struct {
		name string
		want []standalonePage
	}{
		{
			name: "get standalone pages",
			want: []standalonePage{
				{
					name:    "sync",
					handler: &handlers.SyncHandler{},
				},
				{
					name:    "createwallet",
					handler: &handlers.CreateWalletHandler{},
				},
			},
		},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			if got := getStandalonePages(); !reflect.DeepEqual(got, test.want) {
				t.Errorf("getStandalonePages() = %v, want %v", got, test.want)
			}
		})
	}
}
