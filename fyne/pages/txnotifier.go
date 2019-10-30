package pages

import (
	"fmt"
)

type listener struct {
	// Place widgets that should be dynamically updated on fyne.
	// Only update currently viewed tab and only update widget children.
}

func (test *listener) OnTransaction(transaction string) {
	fmt.Println("working OnTransaction ", transaction)
}
func (test *listener) OnTransactionConfirmed(hash string, height int32) {
	fmt.Println("working OnTransactionConfirmed", hash, height)
}
func (test *listener) OnBlockAttached(height int32, timestamp int64) {
	fmt.Println("working OnBlockAttached", height, timestamp)
}

func (app *AppInterface) walletNotificationListener(dcrlistener *listener) {
	// pass dynamic variables.
	app.Dcrlw.TransactionNotification(dcrlistener)
}
