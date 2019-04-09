package handlers

import (
	"github.com/aarzilli/nucular"
	"github.com/raedahgroup/godcr/app/walletcore"
	"github.com/raedahgroup/godcr/nuklear/helpers"
)

type OverviewHandler struct {
	err         error
	isRendering bool
	accounts    []*walletcore.Account
	detailed    bool
}

func (handler *OverviewHandler) BeforeRender() {
	handler.err = nil
	handler.accounts = nil
	handler.isRendering = false
	handler.detailed = false
}

func (handler *OverviewHandler) Render(window *nucular.Window, wallet walletcore.Wallet) {
	if !handler.isRendering {
		handler.isRendering = true
		handler.accounts, handler.err = wallet.AccountsOverview(walletcore.DefaultRequiredConfirmations)
	}

	if pageWindow := helpers.NewWindow("Overview Page", window, nucular.WindowNoScrollbar); pageWindow != nil {
		pageWindow.DrawHeader("Overview")

		if contentWindow := pageWindow.ContentWindow("Balance"); contentWindow != nil {
			contentWindow.Row(25).Dynamic(1)
			contentWindow.Label("Current Total Balance", "LC")

			if handler.err != nil {
				contentWindow.SetErrorMessage(handler.err.Error())
			} else {
				contentWindow.Row(25).Dynamic(1)
				contentWindow.Label(walletcore.WalletBalance(handler.accounts), "LC")
			}

			contentWindow.End()
		}
		pageWindow.End()
	}
}
