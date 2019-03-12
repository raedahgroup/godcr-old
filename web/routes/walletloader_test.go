package routes

import (
	"context"
	"net/http"
	"reflect"
	"sync"
	"testing"
	"text/template"

	"github.com/raedahgroup/godcr/app"
	"github.com/raedahgroup/godcr/app/walletcore"
)

func TestRoutes_walletLoaderMiddleware(t *testing.T) {
	type fields struct {
		walletMiddleware app.WalletMiddleware
		templates        map[string]*template.Template
		blockchain       *Blockchain
		ctx              context.Context
	}
	tests := []struct {
		name   string
		fields fields
		want   func(http.Handler) http.Handler
	}{
		// TODO: Add test cases.
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			routes := &Routes{
				walletMiddleware: test.fields.walletMiddleware,
				templates:        test.fields.templates,
				blockchain:       test.fields.blockchain,
				ctx:              test.fields.ctx,
			}
			if got := routes.walletLoaderMiddleware(); !reflect.DeepEqual(got, test.want) {
				t.Errorf("Routes.walletLoaderMiddleware() = %v, want %v", got, test.want)
			}
		})
	}
}

func TestRoutes_walletLoaderFn(t *testing.T) {
	type fields struct {
		walletMiddleware app.WalletMiddleware
		templates        map[string]*template.Template
		blockchain       *Blockchain
		ctx              context.Context
	}
	tests := []struct {
		name   string
		fields fields
		next   http.Handler
		want   http.Handler
	}{
		// TODO: Add test cases.
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			routes := &Routes{
				walletMiddleware: test.fields.walletMiddleware,
				templates:        test.fields.templates,
				blockchain:       test.fields.blockchain,
				ctx:              test.fields.ctx,
			}
			if got := routes.walletLoaderFn(test.next); !reflect.DeepEqual(got, test.want) {
				t.Errorf("Routes.walletLoaderFn() = %v, want %v", got, test.want)
			}
		})
	}
}

func TestRoutes_syncBlockchain(t *testing.T) {
	type fields struct {
		walletMiddleware app.WalletMiddleware
		templates        map[string]*template.Template
		blockchain       *Blockchain
		ctx              context.Context
	}
	tests := []struct {
		name   string
		fields fields
	}{
		// TODO: Add test cases.
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			routes := &Routes{
				walletMiddleware: test.fields.walletMiddleware,
				templates:        test.fields.templates,
				blockchain:       test.fields.blockchain,
				ctx:              test.fields.ctx,
			}
			routes.syncBlockchain()
		})
	}
}

func TestBlockchain_updateStatus(t *testing.T) {
	type fields struct {
		RWMutex sync.RWMutex
		_status walletcore.SyncStatus
		_report string
	}
	tests := []struct {
		name   string
		fields fields
		report string
		status walletcore.SyncStatus
	}{
		// TODO: Add test cases.
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			b := &Blockchain{
				RWMutex: test.fields.RWMutex,
				_status: test.fields._status,
				_report: test.fields._report,
			}
			b.updateStatus(test.report, test.status)
		})
	}
}

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
		// TODO: Add test cases.
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
		// TODO: Add test cases.
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
