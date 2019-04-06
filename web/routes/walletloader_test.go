package routes

import (
	"reflect"
	"sync"
	"testing"

	"github.com/raedahgroup/godcr/app/walletcore"
)

func TestBlockchain_status(t *testing.T) {
	type fields struct {
		RWMutex sync.RWMutex
		_status walletcore.SyncStatus
		_report string
	}
	tests := []struct {
		name   string
		fields fields
		want   walletcore.SyncStatus
	}{
		{
			name: "blockchain status started",
			fields: fields{
				RWMutex: sync.RWMutex{},
				_status: walletcore.SyncStatusInProgress,
				_report: "Blockchain sync started...",
			},
			want: walletcore.SyncStatusInProgress,
		},
		{
			name: "blockchain status success",
			fields: fields{
				RWMutex: sync.RWMutex{},
				_status: walletcore.SyncStatusSuccess,
				_report: "Blockchain sync completed successfully",
			},
			want: walletcore.SyncStatusSuccess,
		},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			b := &Blockchain{
				RWMutex: test.fields.RWMutex,
				_status: test.fields._status,
				_report: test.fields._report,
			}
			if got := b.status(); !reflect.DeepEqual(got, test.want) {
				t.Errorf("Blockchain.status() = %v, want %v", got, test.want)
			}
		})
	}
}

func TestBlockchain_report(t *testing.T) {
	type fields struct {
		RWMutex sync.RWMutex
		_status walletcore.SyncStatus
		_report string
	}
	tests := []struct {
		name   string
		fields fields
		want   string
	}{
		{
			name: "blockchain report progress",
			fields: fields{
				RWMutex: sync.RWMutex{},
				_status: walletcore.SyncStatusInProgress,
				_report: "Blockchain sync started...",
			},
			want: "Blockchain sync started...",
		},
		{
			name: "blockchain report success",
			fields: fields{
				RWMutex: sync.RWMutex{},
				_status: walletcore.SyncStatusSuccess,
				_report: "Blockchain sync completed successfully",
			},
			want: "Blockchain sync completed successfully",
		},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			b := &Blockchain{
				RWMutex: test.fields.RWMutex,
				_status: test.fields._status,
				_report: test.fields._report,
			}
			if got := b.report(); got != test.want {
				t.Errorf("Blockchain.report() = %v, want %v", got, test.want)
			}
		})
	}
}
