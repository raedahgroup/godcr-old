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

func (handler *OverviewHandler) Render(w *nucular.Window, wallet walletcore.Wallet) {
	if !handler.isRendering {
		handler.isRendering = true
		handler.accounts, handler.err = wallet.AccountsOverview(walletcore.DefaultRequiredConfirmations)
	}

	// draw page
	if page := helpers.NewWindow("Overview Page", w, 0); page != nil {
		page.DrawHeader("Overview")

		if contentWindow := page.ContentWindow("Balance"); contentWindow != nil {
			contentWindow.DrawHeader("Current Total Balance")

			if handler.err != nil {
				contentWindow.SetErrorMessage(handler.err.Error())
			} else {
				contentWindow.Row(25).Dynamic(1)
				contentWindow.Label(walletcore.WalletBalance(handler.accounts), "LC")
			}

			contentWindow.End()
		}
		page.End()
	}
}
