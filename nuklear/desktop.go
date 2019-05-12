package nuklear

import (
	"context"
	"errors"
	"fmt"
	"image"

	"github.com/aarzilli/nucular"
	"github.com/aarzilli/nucular/rect"
	"github.com/raedahgroup/godcr/app"
	"github.com/raedahgroup/godcr/app/config"
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
	currentPage      string
	nextPage         string
	pageChanged      bool
	syncer           *Syncer
	settings         *config.Settings
}

func LaunchApp(ctx context.Context, walletMiddleware app.WalletMiddleware, settings *config.Settings) error {
	desktop := &Desktop{
		walletMiddleware: walletMiddleware,
		pageChanged:      true,
		currentPage:      "overview",
		syncer:           NewSyncer(),
		settings:         settings,
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

	// start syncing in background
	go desktop.syncer.startSyncing(walletMiddleware, masterWindow)

	// draw master window
	masterWindow.Main()
	return nil
}

func (desktop *Desktop) render(window *nucular.Window) {
	if _, isNavPage := desktop.navPages[desktop.currentPage]; isNavPage {
		desktop.renderNavPage(window)
	} else {
		errorMessage := fmt.Sprintf("Page not properly set up: %s", desktop.currentPage)
		nuklog.LogError(errors.New(errorMessage))

		w := &widgets.Window{window}
		w.DisplayMessage(errorMessage, styles.DecredOrangeColor)
	}
}

func (desktop *Desktop) renderNavPage(window *nucular.Window) {
	// this creates the space on the window that will hold 2 widgets
	// the navigation section on the window and the main page content
	entireWindow := window.Row(0).SpaceBegin(2)

	desktop.renderNavWindow(window, entireWindow.H)
	desktop.renderPageContentWindow(window, entireWindow.W, entireWindow.H)
}

func (desktop *Desktop) renderNavWindow(window *nucular.Window, maxHeight int) {
	navSectionRect := rect.Rect{
		X: 0,
		Y: 0,
		W: navWidth,
		H: maxHeight,
	}
	window.LayoutSpacePushScaled(navSectionRect)

	// set style
	styles.SetNavStyle(window.Master())

	// create window and draw nav menu
	widgets.NoScrollGroupWindow("nav-group-window", window, func(navGroupWindow *widgets.Window) {
		navGroupWindow.AddHorizontalSpace(10)
		navGroupWindow.AddColoredLabel(fmt.Sprintf("%s %s", app.DisplayName, desktop.walletMiddleware.NetType()),
			styles.DecredLightBlueColor, widgets.CenterAlign)
		navGroupWindow.AddHorizontalSpace(10)

		for _, page := range getNavPages() {
			if desktop.currentPage == page.name {
				navGroupWindow.AddCurrentNavButton(page.label, func() {
					desktop.changePage(window, page.name)
				})
			} else {
				navGroupWindow.AddBigButton(page.label, func() {
					desktop.changePage(window, page.name)
				})
			}
		}

		// add exit button
		navGroupWindow.AddBigButton("Exit", func() {
			go navGroupWindow.Master().Close()
		})
	})
}

func (desktop *Desktop) renderPageContentWindow(window *nucular.Window, maxWidth, maxHeight int) {
	pageSectionRect := rect.Rect{
		X: navWidth,
		Y: 0,
		W: maxWidth - navWidth,
		H: maxHeight,
	}

	// set style
	styles.SetPageStyle(window.Master())
	window.LayoutSpacePushScaled(pageSectionRect)

	if !desktop.syncer.isDoneSyncing() {
		desktop.syncer.Render(window)
	} else {
		handler := desktop.navPages[desktop.currentPage]
		// ensure that the handler's BeforeRender function is called only once per page call
		// as it initializes page variables
		if desktop.pageChanged {
			handler.BeforeRender(desktop.walletMiddleware, desktop.settings, window.Master().Changed)
			desktop.pageChanged = false
		}
		handler.Render(window)
	}
}

func (desktop *Desktop) changePage(window *nucular.Window, newPage string) {
	desktop.nextPage = newPage
	desktop.currentPage = newPage
	desktop.pageChanged = true
	window.Master().Changed()
}
