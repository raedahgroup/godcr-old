package fyne

import (
	"context"

	fyneApp "fyne.io/fyne/app"
	"fyne.io/fyne/widget"
	"github.com/raedahgroup/godcr/app"
	"github.com/raedahgroup/godcr/fyne/pages"
	"github.com/raedahgroup/godcr/fyne/styles"
	"github.com/raedahgroup/godcr/fyne/widgets"
)

const (
	navWidth                        = 200
	pageWidth                       = 600
	defaultWindowHeight             = 600
	navSectionPageSectionSeparation = 20
)

type navSection struct {
	*widgets.Container
	buttons []*widget.Button
}

type pageSection struct {
	*widgets.Box
	container *widgets.Container
}

type content struct {
	window      *widgets.Window
	navSection  *navSection
	pageSection *pageSection
}

type Desktop struct {
	*content

	ctx              context.Context
	walletMiddleware app.WalletMiddleware
	currentPage      *pages.Page
	syncer           *Syncer
}

func LaunchApp(ctx context.Context, walletMiddleware app.WalletMiddleware) error {
	fyneApp := fyneApp.New()

	desktop := &Desktop{
		ctx:              ctx,
		walletMiddleware: walletMiddleware,
		content: &content{
			window: widgets.NewWindow(app.DisplayName, fyneApp),
		},
		currentPage: pages.GetPages()[0],
		syncer:      NewSyncer(),
	}

	desktop.window.Settings().SetTheme(styles.NewTheme())
	desktop.render()

	return nil
}

func (desktop *Desktop) render() {
	desktop.window.Settings().SetTheme(styles.NewTheme())

	mainContainer := widgets.NewHBoxContainer()
	mainContainer.AddChildContainer(desktop.getNavSection())
	mainContainer.AddChildContainer(desktop.getContentSection())

	changePageFunc := desktop.changePage(desktop.currentPage)
	// start syncing
	desktop.syncer.startSyncing(desktop.walletMiddleware, changePageFunc)
	// render sync view
	desktop.syncer.Render(desktop.pageSection.Box)
	desktop.window.Render(mainContainer)
}

func (desktop *Desktop) getNavSection() *widgets.Container {
	navBox := widgets.NewVBox()
	pages := pages.GetPages()

	desktop.navSection = &navSection{
		buttons:   make([]*widget.Button, len(pages)+1),
		Container: widgets.NewFixedGridLayout(navWidth, defaultWindowHeight),
	}

	for index, page := range pages {
		if index == 0 {
			desktop.currentPage = page
		}

		button := navBox.AddButton(page.Title, desktop.changePage(page))
		desktop.navSection.buttons[index] = button
	}

	// add exit button
	exitButton := navBox.AddButton("Exit", desktop.window.Close)
	desktop.navSection.buttons[len(desktop.navSection.buttons)-1] = exitButton
	desktop.navSection.Container.AddBox(navBox)

	return desktop.navSection.Container
}

func (desktop *Desktop) getContentSection() *widgets.Container {
	container := widgets.NewFixedGridLayout(pageWidth, defaultWindowHeight)
	contentBox := widgets.NewVBox()
	contentBox.SetParent(desktop.content.window)
	container.AddBox(contentBox)

	desktop.pageSection = &pageSection{
		Box:       contentBox,
		container: container,
	}

	return container
}

func (desktop *Desktop) changePage(page *pages.Page) func() {
	return func() {
		// if not done syncing, return
		if !desktop.syncer.isDoneSyncing() {
			return
		}

		// highlight current menu item
		for _, button := range desktop.navSection.buttons {
			if button.Text == page.Title {
				button.Style = widget.PrimaryButton
			} else {
				button.Style = widget.DefaultButton
			}
		}

		desktop.window.RefreshContainer(desktop.navSection.Container)
		// empty page container
		desktop.pageSection.Empty()
		desktop.pageSection.SetTitle(page.Title)

		// render curent page
		page.Handler.Render(desktop.ctx, desktop.walletMiddleware, desktop.pageSection.Box)
	}
}
