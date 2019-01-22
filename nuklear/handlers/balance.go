package handlers

import (
	"github.com/aarzilli/nucular"
	"github.com/raedahgroup/godcr/app/walletcore"
	"github.com/raedahgroup/godcr/nuklear/helpers"
)

type BalanceHandler struct {
	err         error
	isRendering bool
	accounts    []*walletcore.Account
}

func (handler *BalanceHandler) BeforeRender() {
	handler.err = nil
	handler.accounts = nil
	handler.isRendering = false
}

func (handler *BalanceHandler) Render(w *nucular.Window, wallet walletcore.Wallet) {
	if !handler.isRendering {
		handler.isRendering = true
		handler.accounts, handler.err = wallet.AccountsOverview(walletcore.DefaultRequiredConfirmations)
	}

	// draw page
	if page := helpers.NewWindow("Balance Page", w, 0); page != nil {
		page.DrawHeader("Balance")

		// content window
		if content := page.ContentWindow("Balance"); content != nil {
			if handler.err != nil {
				content.SetErrorMessage(handler.err.Error())
			} else {
				content.Row(20).Ratio(0.12, 0.12, 0.15, 0.15, 0.26, 0.20)
				content.Label("Account", "LC")
				content.Label("Total", "LC")
				content.Label("Spendable", "LC")
				content.Label("Locked", "LC")
				content.Label("Voting Authority", "LC")
				content.Label("Unconfirmed", "LC")

				// rows
				for _, v := range handler.accounts {
					content.Label(v.Name, "LC")
					content.Label(helpers.AmountToString(v.Balance.Total.ToCoin()), "LC")
					content.Label(helpers.AmountToString(v.Balance.Spendable.ToCoin()), "LC")
					content.Label(helpers.AmountToString(v.Balance.LockedByTickets.ToCoin()), "LC")
					content.Label(helpers.AmountToString(v.Balance.VotingAuthority.ToCoin()), "LC")
					content.Label(helpers.AmountToString(v.Balance.Unconfirmed.ToCoin()), "LC")
				}
			}
			content.End()
		}
		page.End()
	}
}
