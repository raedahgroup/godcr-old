package help

import (
	"testing"
)

var (
	commandCategories = []*CommandCategory{
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
	}
)

func Test_commandCategoryName(t *testing.T) {
	tests := []struct {
		name              string
		commandName       string
		commandCategories []*CommandCategory
		want              string
	}{
		{
			name:              "help command category",
			commandName:       "help",
			commandCategories: commandCategories,
			want:              "Available Commands",
		},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			if got := commandCategoryName(test.commandName, test.commandCategories); got != test.want {
				t.Errorf("commandCategoryName() = %v, want %v", got, test.want)
			}
		})
	}
}
