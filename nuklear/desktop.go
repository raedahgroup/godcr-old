package nuklear

import (
	"context"
	"errors"

	"fmt"
	"github.com/aarzilli/nucular"
	"github.com/aarzilli/nucular/rect"
	"github.com/raedahgroup/godcr/app"
	"github.com/raedahgroup/godcr/nuklear/nuklog"
	"github.com/raedahgroup/godcr/nuklear/styles"
	"github.com/raedahgroup/godcr/nuklear/widgets"
)

const navWidth = 200

type Desktop struct {
	walletMiddleware app.WalletMiddleware
	currentPage      string
	pageChanged      bool
	navPages         map[string]navPageHandler
	standalonePages  map[string]standalonePageHandler
	quitApp          func()
}

func LaunchApp(ctx context.Context, walletMiddleware app.WalletMiddleware) error {
	desktop := &Desktop{
		walletMiddleware: walletMiddleware,
		pageChanged:      true,
		currentPage:      "overview", // todo rename to sync
	}

	// initialize master window and set style
	masterWindow := nucular.NewMasterWindow(nucular.WindowNoScrollbar, app.Name, desktop.render)
	masterWindow.SetStyle(styles.MasterWindowStyle())

	// initialize fonts for later use
	err := styles.InitFonts()
	if err != nil {
		return err
	}

	// register nav page handlers
	navPages := getNavPages()
	desktop.navPages = make(map[string]navPageHandler, len(navPages))
	for _, page := range navPages {
		desktop.navPages[page.name] = page.handler
	}

	// register standalone page handlers
	standalonePages := getStandalonePages()
	desktop.standalonePages = make(map[string]standalonePageHandler, len(standalonePages))
	for _, page := range standalonePages {
		desktop.standalonePages[page.name] = page.handler
	}

	// open wallet and start blockchain syncing in background
	walletExists, err := walletMiddleware.OpenWalletIfExist(ctx)
	if err != nil {
		return err
	}

	if !walletExists {
		desktop.currentPage = "createwallet"
	}

	desktop.quitApp = masterWindow.Close

	// draw master window
	masterWindow.Main()
	return nil
}

func (desktop *Desktop) render(window *nucular.Window) {
	if handler, isStandalonePage := desktop.standalonePages[desktop.currentPage]; isStandalonePage {
		desktop.renderStandalonePage(window, handler)
	} else if handler, isNavPage := desktop.navPages[desktop.currentPage]; isNavPage {
		desktop.renderNavPage(window, handler)
	} else {
		errorMessage := fmt.Sprintf("Page not properly set up: %s", desktop.currentPage)
		nuklog.LogError(errors.New(errorMessage))

		w := &widgets.Window{window}
		w.DisplayErrorMessage(errorMessage)
	}
}

func (desktop *Desktop) renderStandalonePage(window *nucular.Window, handler standalonePageHandler) {
	window.Row(0).Dynamic(1)

	if desktop.pageChanged {
		handler.BeforeRender()
		desktop.pageChanged = false
	}

	handler.Render(window, desktop.walletMiddleware, desktop.changePage)
}

func (desktop *Desktop) renderNavPage(window *nucular.Window, handler navPageHandler) {
	// this creates the space on the window that will hold 2 widgets
	// the navigation section on the window and the main page content
	entireWindow := window.Row(0).SpaceBegin(2)

	desktop.renderNavSection(window, entireWindow.H)
	renderPageContentSection(window, entireWindow.W, entireWindow.H)

	// ensure that the handler's BeforeRender function is called only once per page call
	// as it initializes page variables
	if desktop.pageChanged {
		handler.BeforeRender()
		desktop.pageChanged = false
	}

	handler.Render(window, desktop.walletMiddleware)
}

func (desktop *Desktop) renderNavSection(window *nucular.Window, maxHeight int) {
	navSection := rect.Rect{
		X: 0,
		Y: 0,
		W: navWidth,
		H: maxHeight,
	}
	window.LayoutSpacePushScaled(navSection)

	// set the window to use the background, font color and other styles for drawing the nav items/buttons
	styles.SetNavStyle(window.Master())

	// then create a group window and draw the nav buttons
	widgets.NoScrollGroupWindow("nav-group-window", window, func(navGroupWindow *widgets.Window) {
		navGroupWindow.AddSpacing(0, 10)
		navGroupWindow.AddColoredLabel(fmt.Sprintf("%s %s", app.DisplayName, desktop.walletMiddleware.NetType()),
			styles.DecredLightBlueColor, widgets.CenterAlign)
		navGroupWindow.AddSpacing(0, 10)

		for _, page := range getNavPages() {
			navGroupWindow.AddBigButton(page.label, func() {
				desktop.changePage(window, page.name)
			})
		}

		// add exit button
		navGroupWindow.AddBigButton("Exit", func() {
			go desktop.quitApp()
		})
	})
}

func renderPageContentSection(window *nucular.Window, maxWidth, maxHeight int) {
	pageSection := rect.Rect{
		X: navWidth,
		Y: 0,
		W: maxWidth - navWidth,
		H: maxHeight,
	}
	window.LayoutSpacePushScaled(pageSection)

	styles.SetPageStyle(window.Master())
}

func (desktop *Desktop) changePage(window *nucular.Window, newPage string) {
	desktop.currentPage = newPage
	desktop.pageChanged = true
	window.Master().Changed()
}
