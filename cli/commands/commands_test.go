package commands

import (
	"reflect"
	"testing"

	"github.com/raedahgroup/godcr/app/help"
)

func TestCategories(t *testing.T) {
	tests := []struct {
		name string
		want []*help.CommandCategory
	}{
		{
			name: "categories",
			want: []*help.CommandCategory{
				{
					Name:         "Available Commands",
					ShortName:    "Commands",
					CommandNames: []string{"create", "detect", "balance", "send", "receive", "history", "showtransaction", "help", "stakeinfo", "purchaseticket"},
				},
				{
					Name:         "Experimental Commands",
					ShortName:    "Experimental",
					CommandNames: []string{"sendcustom"},
				},
			},
		},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			if got := Categories(); !reflect.DeepEqual(got, test.want) {
				t.Errorf("Categories() = %v, want %v", got, test.want)
			}
		})
	}
}
