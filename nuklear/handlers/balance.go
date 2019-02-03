package handlers

import (
	"fmt"

	"github.com/aarzilli/nucular"
	"github.com/raedahgroup/godcr/app/walletcore"
	"github.com/raedahgroup/godcr/nuklear/helpers"
)

type BalanceHandler struct {
	err         error
	isRendering bool
	accounts    []*walletcore.Account
	detailed    bool
}

func (handler *BalanceHandler) BeforeRender() {
	handler.err = nil
	handler.accounts = nil
	handler.isRendering = false
	handler.detailed = false
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
				detailsCheckboxText := "Show details"
				if handler.detailed {
					detailsCheckboxText = "Hide details"
				}

				content.Row(20).Dynamic(1)
				if content.CheckboxText(detailsCheckboxText, &handler.detailed) {
					content.Master().Changed()
				}

				if !handler.detailed && len(handler.accounts) == 1 {
					handler.showSimpleView(content.Window)
				} else {
					handler.showDetailedView(content.Window)
				}
			}
			content.End()
		}
		page.End()
	}
}

func (handler *BalanceHandler) showSimpleView(window *nucular.Window) {
	helpers.SetFont(window, helpers.PageHeaderFont)
	window.Row(25).Dynamic(1)
	window.Label(handler.simpleBalance(handler.accounts[0].Balance, false), "LC")
}

func (handler *BalanceHandler) showDetailedView(window *nucular.Window) {
	helpers.SetFont(window, helpers.NavFont)
	window.Row(20).Ratio(0.16, 0.14, 0.15, 0.15, 0.2, 0.20)
	window.Label("Account Name", "LC")
	window.Label("Balance", "LC")

	if handler.detailed {
		window.Label("Spendable", "LC")
		window.Label("Locked", "LC")
		window.Label("Voting Authority", "LC")
		window.Label("Unconfirmed", "LC")
	}

	// rows
	helpers.SetFont(window, helpers.PageContentFont)
	for _, account := range handler.accounts {
		window.Row(20).Ratio(0.16, 0.14, 0.15, 0.15, 0.2, 0.20)
		window.Label(account.Name, "LC")
		window.Label(handler.simpleBalance(account.Balance, handler.detailed), "LC")

		if handler.detailed {
			window.Label(account.Balance.Spendable.String(), "LC")
			window.Label(account.Balance.LockedByTickets.String(), "LC")
			window.Label(account.Balance.VotingAuthority.String(), "LC")
			window.Label(account.Balance.Unconfirmed.String(), "LC")
		}
	}
}

func (handler *BalanceHandler) simpleBalance(balance *walletcore.Balance, detailed bool) string {
	if detailed || balance.Total == balance.Spendable {
		return balance.Total.String()
	} else {
		return fmt.Sprintf("Total %s (Spendable %s)", balance.Total.String(), balance.Spendable.String())
	}
}
