package nuklear

import (
	"github.com/aarzilli/nucular"
	"github.com/raedahgroup/godcr/app"
	"github.com/raedahgroup/godcr/app/walletcore"
	"github.com/raedahgroup/godcr/nuklear/handlers"
	"github.com/raedahgroup/godcr/nuklear/styles"
	"github.com/raedahgroup/godcr/nuklear/widgets"
)

type navPage struct {
	name    string
	label   string
	handler navPageHandler
}

type navPageHandler interface {
	// BeforeRender initializes all variables that will be needed for displaying the page.
	// It is expected that this method will only be called once i.e. when the page is switched to from a different page.
	// It might be necessary to load some wallet data in background thread
	// after which the app can be notified to repaint the page using `refreshWindowDisplay()`
	// Returns true when done.
	BeforeRender(wallet walletcore.Wallet, refreshWindowDisplay func()) bool

	// Render draws widgets on the provided window.
	// It is usually called several times not only when the page is navigated to.
	// For example, this method will be triggered whenever the mouse is moved, causing the window to repaint.
	Render(window *nucular.Window)
}

type notImplementedNavPageHandler struct {
	pageTitle string
}

func (_ *notImplementedNavPageHandler) BeforeRender(_ walletcore.Wallet, _ func()) bool {
	return true
}
func (p *notImplementedNavPageHandler) Render(window *nucular.Window) {
	widgets.PageContentWindowDefaultPadding(p.pageTitle, window, func(contentWindow *widgets.Window) {
		contentWindow.DisplayMessage("Page not yet implemented", styles.GrayColor)
	})
}

type standalonePage struct {
	name    string
	handler standalonePageHandler
}

type standalonePageHandler interface {
	BeforeRender()
	Render(*nucular.Window, app.WalletMiddleware, func(*nucular.Window, string))
}

func getNavPages() []navPage {
	return []navPage{
		{
			name:    "overview",
			label:   "Overview",
			handler: &handlers.OverviewHandler{},
		},
		{
			name:    "history",
			label:   "History",
			handler: &handlers.HistoryHandler{},
		},
		{
			name:    "send",
			label:   "Send",
			handler: &handlers.SendHandler{},
		},
		{
			name:    "receive",
			label:   "Receive",
			handler: &handlers.ReceiveHandler{},
		},
		{
			name:    "staking",
			label:   "Staking",
			handler: &handlers.StakingHandler{},
		},
		{
			name:    "accounts",
			label:   "Accounts",
			handler: &notImplementedNavPageHandler{"Accounts"},
		},
		{
			name:    "security",
			label:   "Security",
			handler: &notImplementedNavPageHandler{"Security"},
		},
		{
			name:    "settings",
			label:   "Settings",
			handler: &notImplementedNavPageHandler{"Settings"},
		},
	}
}

func getStandalonePages() []standalonePage {
	return []standalonePage{
		{
			name:    "sync",
			handler: &handlers.SyncHandler{},
		},
		{
			name:    "createwallet",
			handler: &handlers.CreateWalletHandler{},
		},
	}
}
