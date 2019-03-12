package routes

import (
	"context"
	"html/template"
	"reflect"
	"testing"

	"github.com/go-chi/chi"
	"github.com/raedahgroup/godcr/app"
)

func TestSetup(t *testing.T) {
	tests := []struct {
		name             string
		ctx              context.Context
		walletMiddleware app.WalletMiddleware
		router           chi.Router
		want             func()
	}{
		// TODO: Add test cases.
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			if got := Setup(test.ctx, test.walletMiddleware, test.router); !reflect.DeepEqual(got, test.want) {
				t.Errorf("Setup() = %v, want %v", got, test.want)
			}
		})
	}
}

func TestRoutes_loadTemplates(t *testing.T) {
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
			routes.loadTemplates()
		})
	}
}

func TestRoutes_loadRoutes(t *testing.T) {
	type fields struct {
		walletMiddleware app.WalletMiddleware
		templates        map[string]*template.Template
		blockchain       *Blockchain
		ctx              context.Context
	}
	tests := []struct {
		name   string
		fields fields
		router chi.Router
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
			routes.loadRoutes(test.router)
		})
	}
}

func TestRoutes_registerRoutesRequiringWallet(t *testing.T) {
	type fields struct {
		walletMiddleware app.WalletMiddleware
		templates        map[string]*template.Template
		blockchain       *Blockchain
		ctx              context.Context
	}
	tests := []struct {
		name   string
		fields fields
		router chi.Router
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
			routes.registerRoutesRequiringWallet(test.router)
		})
	}
}
