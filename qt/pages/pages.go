package pages

import (
	"context"
	"github.com/raedahgroup/godcr/app/walletcore"
	"github.com/therecipe/qt/widgets"
)

type Page interface {
	Setup() *widgets.QWidget
}

// pages that do not use the regular `Setup` method should extend this struct and define a custom setup method
type pageStub struct {}

func (_ pageStub) Setup() *widgets.QWidget {
	return nil
}

type WalletPage interface {
	Page
	SetupWithWallet(ctx context.Context, wallet walletcore.Wallet) *widgets.QWidget
}

func AllPages() map[string]Page {
	return map[string]Page{
		"Status": &statusPage{},
		"Balance": &balancePage{},
	}
}