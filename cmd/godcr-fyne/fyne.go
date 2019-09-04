package main

import (
	"github.com/raedahgroup/godcr/fyne"
	"github.com/decred/dcrd/dcrutil"
)

// todo these would ideally be defined in app/app.go for use by all the interfaces
// but
const appDisplayName = "GoDCR"
var defaultAppDataDir = dcrutil.AppDataDir("godcr", false)

func main() {
	//// Initialize log rotation.  After log rotation has been initialized, the
	//// logger variables may be used.
	//initLogRotator(config.LogFile)
	//defer func() {
	//	if logRotator != nil {
	//		logRotator.Close()
	//	}
	//}()
	//
	//// Parse, validate, and set debug log level(s).
	//if err := parseAndSetDebugLevels(appConfig.DebugLevel); err != nil {
	//	err := fmt.Errorf("loadConfig: %s", err.Error())
	//	fmt.Fprintln(os.Stderr, err)
	//	os.Exit(1)
	//	return
	//}

	//log.Info("Launching desktop app with fyne")
	gui := fyne.InitializeUserInterface(appDisplayName, defaultAppDataDir, "testnet3")
	gui.Launch()
}
