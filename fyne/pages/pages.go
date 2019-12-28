package pages

import (
	"fyne.io/fyne"
	"fyne.io/fyne/widget"

	"github.com/raedahgroup/dcrlibwallet"

	"github.com/raedahgroup/godcr/fyne/assets"
	"github.com/raedahgroup/godcr/fyne/pages/accounts"
	"github.com/raedahgroup/godcr/fyne/pages/history"
	"github.com/raedahgroup/godcr/fyne/pages/overview"
	"github.com/raedahgroup/godcr/fyne/pages/receive"
	"github.com/raedahgroup/godcr/fyne/pages/send"
	"github.com/raedahgroup/godcr/fyne/pages/staking"
)

type PageInitiator interface {
	NewTab() (*widget.TabContainer, error)
	overviewPage() fyne.CanvasObject
	historyPage() fyne.CanvasObject
	sendPage() fyne.CanvasObject
	receivePage() fyne.CanvasObject
	accountsPage() fyne.CanvasObject
	stakingPage() fyne.CanvasObject
	pageHandlers() *pageHandlers
}

type page struct {
	multiWallet *dcrlibwallet.MultiWallet
	tabMenu     *widget.TabContainer
	window 	   	fyne.Window
	handlers  	*pageHandlers
}

type pageHandlers struct {
	overviewHandler *overview.Handler
}

func initiatePages(multiWallet *dcrlibwallet.MultiWallet, tabMenu *widget.TabContainer, window fyne.Window) PageInitiator {
	return &page{
		multiWallet: multiWallet,
		tabMenu:     tabMenu,
		window: 	 window,
		handlers: 	 &pageHandlers{},
	}
}

func (p *page) NewTab() (container *widget.TabContainer, err error) {
	icons, err := assets.GetIcons(assets.OverviewIcon, assets.HistoryIcon, assets.SendIcon,
		assets.ReceiveIcon, assets.AccountsIcon, assets.StakeIcon)
	if err != nil {
		return nil, err
	}

	overviewPage, handler := overview.PageContent(p.multiWallet, p.tabMenu)
	p.handlers.overviewHandler = handler
	container = widget.NewTabContainer(
		widget.NewTabItemWithIcon("Overview", icons[assets.OverviewIcon], overviewPage),
		widget.NewTabItemWithIcon("History", icons[assets.HistoryIcon], widget.NewHBox()),
		widget.NewTabItemWithIcon("Send", icons[assets.SendIcon], widget.NewHBox()),
		widget.NewTabItemWithIcon("Receive", icons[assets.ReceiveIcon], widget.NewHBox()),
		widget.NewTabItemWithIcon("Accounts", icons[assets.AccountsIcon], widget.NewHBox()),
		widget.NewTabItemWithIcon("Staking", icons[assets.StakeIcon], widget.NewHBox()),
	)
	return
}

func (p *page) launchPage() fyne.CanvasObject {
	return nil
}

// get overview page content and initialize its handler
func (p *page) overviewPage() fyne.CanvasObject {
	content, handler := overview.PageContent(p.multiWallet, p.tabMenu)
	p.handlers.overviewHandler = handler
	return content
}

func (p *page) createRestorePage() fyne.CanvasObject {
	return nil
}

func (p *page) historyPage() fyne.CanvasObject {
	return history.PageContent()
}

func (p *page) sendPage() fyne.CanvasObject {
	return send.PageContent(p.multiWallet, p.window)
}

func (p *page) receivePage() fyne.CanvasObject {
	return receive.PageContent(p.multiWallet, p.window)
}

func (p *page) accountsPage() fyne.CanvasObject {
	return accounts.PageContent()
}

func (p *page) stakingPage() fyne.CanvasObject {
	return staking.PageContent()
}

func (p *page) pageHandlers() *pageHandlers {
	return p.handlers
}
