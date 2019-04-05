package fyne

import (
	"context"

	"fyne.io/fyne"
	"fyne.io/fyne/app"
	"fyne.io/fyne/layout"
	"fyne.io/fyne/widget"

	godcrApp "github.com/raedahgroup/godcr/app"
	"github.com/raedahgroup/godcr/fyne/log"
	"github.com/raedahgroup/godcr/fyne/pages"
	"github.com/raedahgroup/godcr/fyne/widgets"
)

const (
	menuSectionWidth                 = 200
	menuSectionPageSectionSeparation = 20
)

type fyneApp struct {
	fyne.App

	ctx              context.Context
	walletMiddleware godcrApp.WalletMiddleware

	mainWindow        fyne.Window
	mainWindowContent fyne.CanvasObject

	menuSectionOnLeft *fyne.Container
	menuButtons       []*widget.Button

	pageSectionOnRight *widget.Box
	pageTitle          *widget.Label
	pageContent        fyne.CanvasObject
}

func LaunchApp(ctx context.Context, walletMiddleware godcrApp.WalletMiddleware) error {
	// open wallet before loading app window in case there's an error while trying to load wallet
	walletExists, err := walletMiddleware.OpenWalletIfExist(ctx)
	if err != nil {
		return err
	}

	this := &fyneApp{
		App:              app.New(),
		ctx:              ctx,
		walletMiddleware: walletMiddleware,
	}

	this.prepareNavSectionOnLeft()
	this.preparePageSectionOnRight()

	// create main window content holder and add menu and page sections, separated with space
	this.mainWindowContent = fyne.NewContainerWithLayout(
		layout.NewHBoxLayout(),
		this.menuSectionOnLeft,
		widgets.NewHSpacer(menuSectionPageSectionSeparation/2),
		this.pageSectionOnRight,
	)

	this.mainWindow = this.NewWindow(godcrApp.DisplayName)

	// main window content will be displayed after sync completes
	// if there's no wallet, the create wallet window will trigger the sync operation after a wallet is created
	if !walletExists {
		this.showCreateWalletWindow()
	} else {
		this.showSyncWindow()
	}

	// fyneApp.Run() blocks until the app is exited, before returning nil error to the caller of this LaunchApp function
	this.Run()

	return nil
}

func (app *fyneApp) prepareNavSectionOnLeft() {
	menuGroup := widget.NewGroup("Menu")

	for _, page := range pages.NavPages() {
		menuButton := widget.NewButton(page.Title, app.displayPageFunc(page))
		menuGroup.Append(menuButton)
		app.menuButtons = append(app.menuButtons, menuButton)
	}

	// add exit menu option
	menuGroup.Append(widget.NewButton("Exit", app.Quit))

	// layout menu using FixedGridLayout to ensure that the provided `menuSectionWidth` is used in rendering the menu group
	menuSectionSize := fyne.NewSize(menuSectionWidth, menuGroup.MinSize().Height)
	app.menuSectionOnLeft = fyne.NewContainerWithLayout(layout.NewFixedGridLayout(menuSectionSize), menuGroup)
}

func (app *fyneApp) preparePageSectionOnRight() {
	// page section contents
	app.pageTitle = widget.NewLabelWithStyle("", fyne.TextAlignLeading, fyne.TextStyle{Italic: true, Bold: true})

	// put page title and scrollable content area in v-box
	app.pageSectionOnRight = widget.NewVBox(app.pageTitle)
}

// displayPageFunc returns the function that will be triggered to display a page
func (app *fyneApp) displayPageFunc(page *pages.Page) func() {
	return func() {
		app.pageTitle.SetText(page.Title)
		app.highlightCurrentPageMenuButton(page.Title)

		if simplePageLoader, ok := page.PageLoader.(pages.SimplePageLoader); ok {
			simplePageLoader.Load(app.updatePageFunc)
		} else if walletPageLoader, ok := page.PageLoader.(pages.WalletPageLoader); ok {
			walletPageLoader.Load(app.ctx, app.walletMiddleware, app.updatePageFunc)
		} else {
			log.PrintError("Page not properly set up: ", page.Title)
		}
	}
}

func (app *fyneApp) highlightCurrentPageMenuButton(currentPage string) {
	for _, menuButton := range app.menuButtons {
		if menuButton.Text == currentPage {
			menuButton.Style = widget.PrimaryButton
		} else {
			menuButton.Style = widget.DefaultButton
		}
	}

	// refresh menu section so the changes made in this function reflects
	app.mainWindow.Canvas().Refresh(app.menuSectionOnLeft)
}

func (app *fyneApp) updatePageFunc(pageContent fyne.CanvasObject) {
	app.pageContent = pageContent
	app.resizeScrollableContainer()
}

// resizeScrollableContainer ensures that
// - the content of each page is wrapped in scrollable container
// - the scrollable container takes the maximum space available
//
// The idea is, if the content size is bigger than the maximum space available,
// scroll bars become visible and more of the content can be seen by scrolling.
//
// We really only need to resize the scroll container when the window is resized
// but can't seem to find a window resize event
func (app *fyneApp) resizeScrollableContainer() {
	// calculate the maximum available width and height to use for scroll container
	windowSize := app.mainWindow.Content().Size()
	scrollAreaWidth := windowSize.Width - menuSectionWidth - menuSectionPageSectionSeparation
	scrollAreaHeight := windowSize.Height - app.pageTitle.Size().Height
	scrollAreaSize := fyne.NewSize(scrollAreaWidth, scrollAreaHeight)

	// use calculated max size to layout the scrollable container
	scrollContainerLayout := layout.NewFixedGridLayout(scrollAreaSize)
	scrollableContainer := fyne.NewContainerWithLayout(scrollContainerLayout, widget.NewScrollContainer(app.pageContent))

	// must clear items and re-add otherwise the added content will not display
	app.pageSectionOnRight.Children = []fyne.CanvasObject{}
	app.pageSectionOnRight.Append(app.pageTitle)
	app.pageSectionOnRight.Append(scrollableContainer)
}
