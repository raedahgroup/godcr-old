package fyne

import (
	"fyne.io/fyne"
	"fyne.io/fyne/app"
	"fyne.io/fyne/layout"
	"fyne.io/fyne/widget"

	godcrApp "github.com/raedahgroup/godcr/app"
	"github.com/raedahgroup/godcr/fyne/pages"
)

type fyneApp struct {
	fyne.App
	walletMiddleware godcrApp.WalletMiddleware
	menuButtons []*widget.Button
	pageContentSection *widget.Box
}

func LaunchApp(walletMiddleware godcrApp.WalletMiddleware) {
	this := &fyneApp{
		App: app.New(),
		walletMiddleware:walletMiddleware,
	}

	mainWindow := this.NewWindow(godcrApp.DisplayName)
	mainWindowContent := fyne.NewContainerWithLayout(layout.NewHBoxLayout())

	menuOptionsHolder := widget.NewVBox()
	navPages := pages.NavPages()
	this.menuButtons = make([]*widget.Button, len(navPages))
	for i, page := range pages.NavPages() {
		this.menuButtons[i] = widget.NewButton(page.Title, this.displayPageFunc(page))
		menuOptionsHolder.Append(this.menuButtons[i])
	}

	// add menu to main window
	menuSection := widget.NewGroup("Menu", menuOptionsHolder)
	menuSectionLayout := layout.NewFixedGridLayout(fyne.NewSize(200, menuSection.MinSize().Height))
	menuSectionContainer := fyne.NewContainerWithLayout(menuSectionLayout, menuSection)
	mainWindowContent.AddObject(menuSectionContainer)

	// add page content to main window
	this.pageContentSection = widget.NewVBox()
	pageContentLayout := layout.NewGridLayout(1)
	pageContentContainer := fyne.NewContainerWithLayout(pageContentLayout, this.pageContentSection)
	mainWindowContent.AddObject(pageContentContainer)

	mainWindow.SetContent(mainWindowContent)

	// ShowAndRun blocks until the app is exited, then returns to the caller of this LaunchApp function
	mainWindow.ShowAndRun()
}

// displayPageFunc returns the function that will be triggered to display a page
func (app *fyneApp) displayPageFunc(page *pages.Page) func() {
	return func() {
		app.pageContentSection.Children = []fyne.CanvasObject{
			widget.NewLabel(page.Title),
			page.Content(app.walletMiddleware),
		}
	}
}
