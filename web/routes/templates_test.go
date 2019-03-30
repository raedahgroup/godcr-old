package routes

import (
	"reflect"
	"testing"
)

func Test_templates(t *testing.T) {
	tests := []struct {
		name string
		want []templateData
	}{
		{
			name: "test templates",
			want: []templateData{
				{"error.html", "web/views/error.html"},
				{"createwallet.html", "web/views/createwallet.html"},
				{"balance.html", "web/views/balance.html"},
				{"send.html", "web/views/send.html"},
				{"receive.html", "web/views/receive.html"},
				{"history.html", "web/views/history.html"},
				{"transaction_details.html", "web/views/transaction_details.html"},
				{"staking.html", "web/views/staking.html"},
			},
		},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			if got := templates(); !reflect.DeepEqual(got, test.want) {
				t.Errorf("templates() = %v, want %v", got, test.want)
			}
		})
	}
}
