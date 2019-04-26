package fyne

import (
	"context"

	"fyne.io/fyne"
	"fyne.io/fyne/app"
	"fyne.io/fyne/layout"
	"fyne.io/fyne/widget"

	godcrApp "github.com/raedahgroup/godcr/app"
	"github.com/raedahgroup/godcr/fyne/styles"
)

const (
	menuSectionWidth                 = 200
	menuSectionPageSectionSeparation = 20
)

// menuSection represents the menu
type nav struct {
	*fyne.Container
	buttons []*widget.Button
}

// pageSection represents the current page
type page struct {
}

type Desktop struct {
	fyne.App

	ctx              context.Context
	walletMiddleware godcrApp.WalletMiddleware

	window fyne.Window
	nav    *nav
}

func LaunchApp(ctx context.Context, walletMiddleware godcrApp.WalletMiddleware) error {
	fyneApp := app.New()
	fyneApp.Settings().SetTheme(styles.NewTheme())

	desktop := &Desktop{
		App:              fyneApp,
		ctx:              ctx,
		walletMiddleware: walletMiddleware,
		window:           fyneApp.NewWindow(godcrApp.DisplayName),
		nav:              &nav{},
	}

	// TODO: start syncing

	desktop.window.SetContent(fyne.NewContainerWithLayout(
		layout.NewHBoxLayout(),
		desktop.renderNavWindow(500),
		desktop.renderPageContentWindow(),
	))
	desktop.window.Show()

	desktop.Run()
	return nil
}

func (desktop *Desktop) renderNavWindow(height int) *fyne.Container {
	menuGroup := widget.NewVBox()
	for _, page := range getPages() {
		button := widget.NewButton(page.Title, desktop.changePage(page))
		menuGroup.Append(button)
		desktop.nav.buttons = append(desktop.nav.buttons, button)
	}

	menuGroup.Append(widget.NewButton("Exit", desktop.Quit))

	size := fyne.NewSize(menuSectionWidth, height)
	container := fyne.NewContainerWithLayout(
		layout.NewFixedGridLayout(size),
		menuGroup,
	)
	desktop.nav.Container = container

	return container
}

func (desktop *Desktop) renderPageContentWindow() *widget.Box {
	title := widget.NewLabel(godcrApp.DisplayName)
	return widget.NewVBox(title)
}

func (desktop *Desktop) render() {

}

func (desktop *Desktop) changePage(currentPage *Page) func() {
	return func() {
		// highlight current menu item
		for _, button := range desktop.nav.buttons {
			if button.Text == currentPage.Title {
				button.Style = widget.PrimaryButton
			} else {
				button.Style = widget.DefaultButton
			}
		}
		desktop.window.Canvas().Refresh(desktop.nav)
	}
}
