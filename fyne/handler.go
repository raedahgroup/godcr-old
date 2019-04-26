package fyne

import (
	"github.com/raedahgroup/godcr/app/walletcore"
	"github.com/raedahgroup/godcr/fyne/pages"
)

type Page struct {
	Title   string
	handler pageHandler
}

type pageHandler interface {
	// BeforeRender initializes all variables that will be needed for displaying the page.
	// It is expected that this method will only be called once i.e. when the page is switched to from a different page.
	// It might be necessary to load some wallet data in background thread
	// after which the app can be notified to repaint the page using `refreshWindowDisplay()`
	// Returns true when done.
	BeforeRender(wallet walletcore.Wallet, refreshWindowDisplay func()) bool

	// Render draws widgets on the provided window.
	// It is usually called several times not only when the page is navigated to.
	// For example, this method will be triggered whenever the mouse is moved, causing the window to repaint.
	Render()
}

func getPages() []*Page {
	return []*Page{
		{
			"Overview",
			&pages.OverviewHandler{},
		},
		{
			"History",
			defaultPageNotImplemented,
		},
		{
			"Send",
			defaultPageNotImplemented,
		},
		{
			"Receive",
			defaultPageNotImplemented,
		},
		{
			"Staking",
			defaultPageNotImplemented,
		},
		{
			"Accounts",
			defaultPageNotImplemented,
		},
		{
			"Security",
			defaultPageNotImplemented,
		},
		{
			"Settings",
			defaultPageNotImplemented,
		},
	}
}

type pageNotImplemented struct{}

func (_ *pageNotImplemented) BeforeRender(wallet walletcore.Wallet, refreshWindowDisplay func()) bool {
	return true
}

func (_ *pageNotImplemented) Render() {

}

var defaultPageNotImplemented = &pageNotImplemented{}
