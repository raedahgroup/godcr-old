package nuklear

import (
	"context"
	"testing"

	"github.com/aarzilli/nucular"
	"github.com/raedahgroup/godcr/app"
)

func TestLaunchApp(t *testing.T) {
	tests := []struct {
		name             string
		ctx              context.Context
		walletMiddleware app.WalletMiddleware
		wantErr          bool
	}{
		// TODO: Add test cases.
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			if err := LaunchApp(test.ctx, test.walletMiddleware); (err != nil) != test.wantErr {
				t.Errorf("LaunchApp() error = %v, wantErr %v", err, test.wantErr)
			}
		})
	}
}

func TestDesktop_render(t *testing.T) {
	type fields struct {
		masterWindow     nucular.MasterWindow
		walletMiddleware app.WalletMiddleware
		currentPage      string
		pageChanged      bool
		navPages         map[string]navPageHandler
		standalonePages  map[string]standalonePageHandler
	}
	tests := []struct {
		name   string
		fields fields
		window *nucular.Window
	}{
		// TODO: Add test cases.
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			desktop := &Desktop{
				masterWindow:     test.fields.masterWindow,
				walletMiddleware: test.fields.walletMiddleware,
				currentPage:      test.fields.currentPage,
				pageChanged:      test.fields.pageChanged,
				navPages:         test.fields.navPages,
				standalonePages:  test.fields.standalonePages,
			}
			desktop.render(test.window)
		})
	}
}

func TestDesktop_renderStandalonePage(t *testing.T) {
	type fields struct {
		masterWindow     nucular.MasterWindow
		walletMiddleware app.WalletMiddleware
		currentPage      string
		pageChanged      bool
		navPages         map[string]navPageHandler
		standalonePages  map[string]standalonePageHandler
	}
	tests := []struct {
		name    string
		fields  fields
		window  *nucular.Window
		handler standalonePageHandler
	}{
		// TODO: Add test cases.
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			desktop := &Desktop{
				masterWindow:     test.fields.masterWindow,
				walletMiddleware: test.fields.walletMiddleware,
				currentPage:      test.fields.currentPage,
				pageChanged:      test.fields.pageChanged,
				navPages:         test.fields.navPages,
				standalonePages:  test.fields.standalonePages,
			}
			desktop.renderStandalonePage(test.window, test.handler)
		})
	}
}

func TestDesktop_renderNavPage(t *testing.T) {
	type fields struct {
		masterWindow     nucular.MasterWindow
		walletMiddleware app.WalletMiddleware
		currentPage      string
		pageChanged      bool
		navPages         map[string]navPageHandler
		standalonePages  map[string]standalonePageHandler
	}
	tests := []struct {
		name    string
		fields  fields
		window  *nucular.Window
		handler navPageHandler
	}{
		// TODO: Add test cases.
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			desktop := &Desktop{
				masterWindow:     test.fields.masterWindow,
				walletMiddleware: test.fields.walletMiddleware,
				currentPage:      test.fields.currentPage,
				pageChanged:      test.fields.pageChanged,
				navPages:         test.fields.navPages,
				standalonePages:  test.fields.standalonePages,
			}
			desktop.renderNavPage(test.window, test.handler)
		})
	}
}

func TestDesktop_changePage(t *testing.T) {
	type fields struct {
		masterWindow     nucular.MasterWindow
		walletMiddleware app.WalletMiddleware
		currentPage      string
		pageChanged      bool
		navPages         map[string]navPageHandler
		standalonePages  map[string]standalonePageHandler
	}
	tests := []struct {
		name   string
		fields fields
		page   string
	}{
		// TODO: Add test cases.
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			desktop := &Desktop{
				masterWindow:     test.fields.masterWindow,
				walletMiddleware: test.fields.walletMiddleware,
				currentPage:      test.fields.currentPage,
				pageChanged:      test.fields.pageChanged,
				navPages:         test.fields.navPages,
				standalonePages:  test.fields.standalonePages,
			}
			desktop.changePage(test.page)
		})
	}
}
