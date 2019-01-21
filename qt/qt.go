package qt

import (
	"context"
	"os"

	"github.com/raedahgroup/godcr/app"
	"github.com/raedahgroup/godcr/qt/pages"

	"github.com/therecipe/qt/widgets"
)

const (
	minWindowWidth  = 600
	minWindowHeight = 400
)

func LaunchApp(ctx context.Context, walletMiddleware app.WalletMiddleware) error {
	// needs to be called once before you can start using the QWidgets
	qtApp := widgets.NewQApplication(len(os.Args), os.Args)

	// create a window and set the title
	window := widgets.NewQMainWindow(nil, 0)
	window.SetMinimumSize2(minWindowWidth, minWindowHeight)
	window.SetWindowTitle(app.Name)

	// todo check if wallet exists and if not, show a create wallet page instead

	err := walletMiddleware.OpenWallet()
	if err != nil {
		return err
	}

	// todo run blockchain sync in a goroutine and ensure no page is accessible until sync is completed

	// create tab widget to hold pages for different godcr functions
	tabWidget := widgets.NewQTabWidget(window)
	window.SetCentralWidget(tabWidget)

	for pageName, page := range pages.AllPages() {
		pageWidget := page.Setup()
		if walletPage, ok := page.(pages.WalletPage); ok {
			pageWidget = walletPage.SetupWithWallet(ctx, walletMiddleware)
		}
		pageWidget.SetAccessibleName(pageName)
		tabWidget.AddTab(pageWidget, pageName)
	}

	// make the window visible
	window.Show()

	// start the main Qt event loop and block until app.Exit() is called or the window is closed by the user
	qtApp.Exec()
	return nil
}
