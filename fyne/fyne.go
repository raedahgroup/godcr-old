package fyne

import (
	"context"

	"fyne.io/fyne"
	"fyne.io/fyne/app"
	"fyne.io/fyne/layout"
	"fyne.io/fyne/widget"

	godcrApp "github.com/raedahgroup/godcr/app"
	"github.com/raedahgroup/godcr/fyne/pages"
)

type fyneApp struct {
	fyne.App

	ctx context.Context
	walletMiddleware godcrApp.WalletMiddleware

	mainWindow fyne.Window
	menuButtons []*widget.Button
	menuSection *widget.Group
	pageContentSection *widget.Box
}

func LaunchApp(ctx context.Context, walletMiddleware godcrApp.WalletMiddleware) error {
	// open wallet and start sync operation in background before loading app window
	walletExists, err := walletMiddleware.OpenWalletIfExist(ctx)
	if err != nil {
		return err
	}

	this := &fyneApp{
		App: app.New(),
		ctx: ctx,
		walletMiddleware:walletMiddleware,
	}

	this.mainWindow = this.NewWindow(godcrApp.DisplayName)
	mainWindowContent := fyne.NewContainerWithLayout(layout.NewHBoxLayout())

	menuOptionsHolder := widget.NewVBox()
	navPages := pages.NavPages()
	this.menuButtons = make([]*widget.Button, len(navPages))
	for i, page := range pages.NavPages() {
		this.menuButtons[i] = widget.NewButton(page.Title, this.displayPageFunc(page))
		menuOptionsHolder.Append(this.menuButtons[i])
	}

	// add exit menu option
	menuOptionsHolder.Append(widget.NewButton("Exit", this.Quit))

	// add menu to main window
	this.menuSection = widget.NewGroup("Menu", menuOptionsHolder)
	menuSectionLayout := layout.NewFixedGridLayout(fyne.NewSize(200, this.menuSection.MinSize().Height))
	menuSectionContainer := fyne.NewContainerWithLayout(menuSectionLayout, this.menuSection)
	mainWindowContent.AddObject(menuSectionContainer)

	// add page content to main window
	this.pageContentSection = widget.NewVBox()
	pageContentLayout := layout.NewGridLayout(1)
	pageContentContainer := fyne.NewContainerWithLayout(pageContentLayout, this.pageContentSection)
	mainWindowContent.AddObject(pageContentContainer)

	this.mainWindow.SetContent(mainWindowContent)
	this.mainWindow.CenterOnScreen()
	this.mainWindow.SetFullScreen(true)

	_ = walletExists
	// main window will be displayed after sync completes
	// if there's no wallet, the create wallet window will trigger the sync operation after a wallet is created
	//if !walletExists {
	//	this.showCreateWalletWindow()
	//} else {
	//	this.showSyncWindow()
	//}

	this.mainWindow.Show()

	// fyneApp.Run() blocks until the app is exited, before returning nil error to the caller of this LaunchApp function
	this.Run()

	return nil
}

// displayPageFunc returns the function that will be triggered to display a page
func (app *fyneApp) displayPageFunc(page *pages.Page) func() {
	return func() {
		app.highlightCurrentPageMenuButton(page.Title)

		pageTitle := widget.NewLabel(page.Title)
		pageContent := page.Content(app.walletMiddleware)

		app.pageContentSection.Children = []fyne.CanvasObject{
			pageTitle,
			pageContent,
		}

		app.mainWindow.Canvas().Refresh(app.pageContentSection)
		app.mainWindow.Canvas().Refresh(app.pageContentSection)

		// resize page content section
		//pageContentSize := pageContent.MinSize()
		//pageTitleSize := pageTitle.MinSize()
		//
		//pageContentSize.Height += pageTitleSize.Height
		//if pageTitleSize.Width > pageContentSize.Width {
		//	pageContentSize.Width = pageTitleSize.Width
		//}
		//
		//app.pageContentSection.Resize(pageContentSize)
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

}
