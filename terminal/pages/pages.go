package pages

import (
	"github.com/decred/slog"
	"github.com/raedahgroup/dcrlibwallet"
	"github.com/raedahgroup/godcr/terminal/primitives"
	"github.com/rivo/tview"
)

type page struct {
	Name     string
	Shortcut rune
	Content  func() tview.Primitive
}

var commonPageData struct {
	app                 *tview.Application
	log                 slog.Logger
	wallet              *dcrlibwallet.LibWallet
	hintTextView        *primitives.TextView
	clearAllPageContent func()
}

func Setup(app *tview.Application, log slog.Logger, dcrlw *dcrlibwallet.LibWallet,
	hintTextView *primitives.TextView, clearAllPageContent func()) {

	commonPageData.app = app
	commonPageData.log = log
	commonPageData.wallet = dcrlw
	commonPageData.hintTextView = hintTextView
	commonPageData.clearAllPageContent = clearAllPageContent
}

func All() []*page {
	return []*page{
		{Name: "Accounts", Shortcut: 'a', Content: accountsPage},
	}
}
