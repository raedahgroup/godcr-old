package routes

import (
	"context"
	"html/template"
	"net/http"
	"testing"

	"github.com/raedahgroup/godcr/app"
)

func TestRoutes_render(t *testing.T) {
	type fields struct {
		walletMiddleware app.WalletMiddleware
		templates        map[string]*template.Template
		blockchain       *Blockchain
		ctx              context.Context
	}
	tests := []struct {
		name    string
		fields  fields
		tplName string
		data    interface{}
		res     http.ResponseWriter
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
			routes.render(test.tplName, test.data, test.res)
		})
	}
}

func TestRoutes_renderError(t *testing.T) {
	type fields struct {
		walletMiddleware app.WalletMiddleware
		templates        map[string]*template.Template
		blockchain       *Blockchain
		ctx              context.Context
	}
	tests := []struct {
		name         string
		fields       fields
		errorMessage string
		res          http.ResponseWriter
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
			routes.renderError(test.errorMessage, test.res)
		})
	}
}

func TestRoutes_renderNoWalletError(t *testing.T) {
	type fields struct {
		walletMiddleware app.WalletMiddleware
		templates        map[string]*template.Template
		blockchain       *Blockchain
		ctx              context.Context
	}
	tests := []struct {
		name   string
		fields fields
		res    http.ResponseWriter
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
			routes.renderNoWalletError(test.res)
		})
	}
}

func Test_renderJSON(t *testing.T) {
	tests := []struct {
		name string
		data interface{}
		res  http.ResponseWriter
	}{
		// TODO: Add test cases.
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			renderJSON(test.data, test.res)
		})
	}
}
