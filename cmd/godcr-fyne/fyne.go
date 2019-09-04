package main

import (
	"github.com/decred/dcrd/dcrutil"
	"github.com/raedahgroup/godcr/fyne"
)

// todo these would ideally be defined in app/app.go for use by all the interfaces
// but are kept here temporarily since attempting to use any property from the `app` package
// will cause fyne to use dcrlibwallet/wip branch which is unwanted.
const appDisplayName = "GoDCR"
var defaultAppDataDir = dcrutil.AppDataDir("godcr", false)

func main() {
	gui := fyne.InitializeUserInterface(appDisplayName, defaultAppDataDir, "testnet3")
	gui.Launch()
}
