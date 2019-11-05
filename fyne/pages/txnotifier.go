package pages

import (
	"fmt"

	"fyne.io/fyne/widget"
)

type listener struct {
	tabMenu *widget.TabContainer
	// Place widgets that should be dynamically updated on fyne.
	// Only update currently viewed tab and only update widget children.
}

func (test *listener) OnTransaction(transaction string) {
	if test.tabMenu.CurrentTabIndex() == 0 {
		// place overview page dynamic data here
	} else if test.tabMenu.CurrentTabIndex() == 2 {
		// place send page dynamic data here
	} else if test.tabMenu.CurrentTabIndex() == 3 {
		// place receive page dynamic data here
	} else if test.tabMenu.CurrentTabIndex() == 2 {
		// place account page dynamic data here
	}
}

func (test *listener) OnTransactionConfirmed(hash string, height int32) {
	fmt.Println("working OnTransactionConfirmed", hash, height)
	if test.tabMenu.CurrentTabIndex() == 0 {
		// place overview page dynamic data here
	} else if test.tabMenu.CurrentTabIndex() == 2 {
		// place send page dynamic data here
	} else if test.tabMenu.CurrentTabIndex() == 3 {
		// place receive page dynamic data here
	} else if test.tabMenu.CurrentTabIndex() == 2 {
		// place account page dynamic data here
	}
}

func (test *listener) OnBlockAttached(height int32, timestamp int64) {
	fmt.Println("working OnBlockAttached", height, timestamp)
	if test.tabMenu.CurrentTabIndex() == 0 {
		// place overview page dynamic data here
	} else if test.tabMenu.CurrentTabIndex() == 2 {
		// place send page dynamic data here
	} else if test.tabMenu.CurrentTabIndex() == 3 {
		// place receive page dynamic data here
	} else if test.tabMenu.CurrentTabIndex() == 2 {
		// place account page dynamic data here
	}
}

func (app *AppInterface) walletNotificationListener(dcrlistener *listener) {
	// pass dynamic variables.
	dcrlistener.tabMenu = app.tabMenu
	app.Dcrlw.TransactionNotification(dcrlistener)
}
