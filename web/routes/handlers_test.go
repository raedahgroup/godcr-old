package routes

import (
	"context"
	"net/http"
	"testing"
	"text/template"

	"github.com/raedahgroup/godcr/app"
)

func TestRoutes_createWalletPage(t *testing.T) {
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
		req    *http.Request
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
			routes.createWalletPage(test.res, test.req)
		})
	}
}

func TestRoutes_createWallet(t *testing.T) {
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
		req    *http.Request
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
			routes.createWallet(test.res, test.req)
		})
	}
}

func TestRoutes_balancePage(t *testing.T) {
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
		req    *http.Request
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
			routes.balancePage(test.res, test.req)
		})
	}
}

func TestRoutes_sendPage(t *testing.T) {
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
		req    *http.Request
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
			routes.sendPage(test.res, test.req)
		})
	}
}

func TestRoutes_submitSendTxForm(t *testing.T) {
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
		req    *http.Request
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
			routes.submitSendTxForm(test.res, test.req)
		})
	}
}

func TestRoutes_receivePage(t *testing.T) {
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
		req    *http.Request
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
			routes.receivePage(test.res, test.req)
		})
	}
}

func TestRoutes_generateReceiveAddress(t *testing.T) {
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
		req    *http.Request
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
			routes.generateReceiveAddress(test.res, test.req)
		})
	}
}

func TestRoutes_getUnspentOutputs(t *testing.T) {
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
		req    *http.Request
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
			routes.getUnspentOutputs(test.res, test.req)
		})
	}
}

func TestRoutes_historyPage(t *testing.T) {
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
		req    *http.Request
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
			routes.historyPage(test.res, test.req)
		})
	}
}

func TestRoutes_transactionDetailsPage(t *testing.T) {
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
		req    *http.Request
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
			routes.transactionDetailsPage(test.res, test.req)
		})
	}
}

func TestRoutes_stakeInfoPage(t *testing.T) {
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
		req    *http.Request
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
			routes.stakeInfoPage(test.res, test.req)
		})
	}
}

func TestRoutes_purchaseTicketsPage(t *testing.T) {
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
		req    *http.Request
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
			routes.purchaseTicketsPage(test.res, test.req)
		})
	}
}

func TestRoutes_submitPurchaseTicketsForm(t *testing.T) {
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
		req    *http.Request
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
			routes.submitPurchaseTicketsForm(test.res, test.req)
		})
	}
}
