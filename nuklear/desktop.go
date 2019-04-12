package nuklear

import (
	"context"
	"errors"
	"fmt"
	"image"

	"github.com/aarzilli/nucular"
	"github.com/aarzilli/nucular/rect"
	"github.com/raedahgroup/godcr/app"
	"github.com/raedahgroup/godcr/nuklear/nuklog"
	"github.com/raedahgroup/godcr/nuklear/styles"
	"github.com/raedahgroup/godcr/nuklear/widgets"
)

const (
	navWidth            = 200
	defaultWindowWidth  = 800
	defaultWindowHeight = 600
)

type Desktop struct {
	walletMiddleware app.WalletMiddleware
	navPages         map[string]navPageHandler
	standalonePages  map[string]standalonePageHandler
	currentPage      string
	nextPage         string
	pageChanged      bool
}

func LaunchApp(ctx context.Context, walletMiddleware app.WalletMiddleware) error {
	desktop := &Desktop{
		walletMiddleware: walletMiddleware,
		pageChanged:      true,
		currentPage:      "sync",
	}

	// initialize master window and set style
	windowSize := image.Point{defaultWindowWidth, defaultWindowHeight}
	masterWindow := nucular.NewMasterWindowSize(nucular.WindowNoScrollbar, app.Name, windowSize, desktop.render)
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
		w.DisplayMessage(errorMessage, styles.DecredOrangeColor)
	}
}

func (desktop *Desktop) renderStandalonePage(window *nucular.Window, handler standalonePageHandler) {
	if desktop.currentPage != desktop.nextPage && desktop.nextPage != "" {
		// page navigation changes may take some seconds to effect
		// causing this method to receive the wrong handler
		desktop.currentPage = desktop.nextPage
		desktop.pageChanged = true
		window.Master().Changed()
		return
	}

	window.Row(0).Dynamic(1)
	styles.SetPageStyle(window.Master())

	if desktop.pageChanged {
		handler.BeforeRender()
		desktop.pageChanged = false
	}

	handler.Render(window, desktop.walletMiddleware, desktop.changePage)
}

func (desktop *Desktop) renderNavPage(window *nucular.Window, handler navPageHandler) {
	if desktop.currentPage != desktop.nextPage && desktop.nextPage != "" {
		// page navigation changes may take some seconds to effect
		// causing this method to receive the wrong handler
		desktop.currentPage = desktop.nextPage
		desktop.pageChanged = true
		window.Master().Changed()
		return
	}

	// this creates the space on the window that will hold 2 widgets
	// the navigation section on the window and the main page content
	entireWindow := window.Row(0).SpaceBegin(2)

	desktop.renderNavSection(window, entireWindow.H)
	renderPageContentSection(window, entireWindow.W, entireWindow.H)

	// Only call handler.Render if the page is not being switched to for the first time (i.e !desktop.pageChanged).
	// If it is (i.e. desktop.pageChanged == true), then only call handler.Render after handler.BeforeRender returns.
	if !desktop.pageChanged || handler.BeforeRender(desktop.walletMiddleware, window.Master().Changed) {
		desktop.pageChanged = false
		handler.Render(window)
	}
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
		navGroupWindow.AddHorizontalSpace(10)
		navGroupWindow.AddColoredLabel(fmt.Sprintf("%s %s", app.DisplayName, desktop.walletMiddleware.NetType()),
			styles.DecredLightBlueColor, widgets.CenterAlign)
		navGroupWindow.AddHorizontalSpace(10)

		for _, page := range getNavPages() {
			navGroupWindow.AddBigButton(page.label, func() {
				desktop.changePage(window, page.name)
			})
		}

		// add exit button
		navGroupWindow.AddBigButton("Exit", func() {
			go navGroupWindow.Master().Close()
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
	desktop.nextPage = newPage
	window.Master().Changed()
}
