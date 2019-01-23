package handlers

import (
	"github.com/aarzilli/nucular"
	"github.com/raedahgroup/godcr/app/walletcore"
	"github.com/raedahgroup/godcr/nuklear/handlers/widgets"
	"github.com/raedahgroup/godcr/nuklear/helpers"
)

type TransactionsHandler struct {
	err                    error
	isRendering            bool
	transactions           []*walletcore.Transaction
	hasFetchedTransactions bool
}

func (handler *TransactionsHandler) BeforeRender() {
	handler.err = nil
	handler.transactions = nil
	handler.isRendering = false
	handler.hasFetchedTransactions = false
}

func (handler *TransactionsHandler) Render(window *nucular.Window, wallet walletcore.Wallet) {
	if !handler.isRendering {
		handler.isRendering = true
		go handler.fetchTransactions(wallet, window)
	}

	// draw page
	if pageWindow := helpers.NewWindow("Transactions Page", window, 0); pageWindow != nil {
		pageWindow.DrawHeader("Transactions")

		// content window
		if contentWindow := pageWindow.ContentWindow("Transactions"); contentWindow != nil {
			if handler.hasFetchedTransactions {
				if handler.err != nil {
					contentWindow.SetErrorMessage(handler.err.Error())
				} else {
					helpers.SetFont(window, helpers.NavFont)
					contentWindow.Row(20).Ratio(0.18, 0.08, 0.1, 0.15, 0.15, 0.7)
					contentWindow.Label("Date", "LC")
					contentWindow.Label("Amount", "LC")
					contentWindow.Label("Fee", "LC")
					contentWindow.Label("Direction", "LC")
					contentWindow.Label("Type", "LC")
					contentWindow.Label("Hash", "LC")

					for _, tx := range handler.transactions {
						helpers.SetFont(window, helpers.PageContentFont)
						contentWindow.Row(20).Ratio(0.18, 0.08, 0.1, 0.15, 0.15, 0.7)
						contentWindow.Label(tx.FormattedTime, "LC")
						contentWindow.Label(helpers.AmountToString(tx.Amount.ToCoin()), "LC")
						contentWindow.Label(helpers.AmountToString(tx.Fee.ToCoin()), "LC")
						contentWindow.Label(tx.Direction.String(), "LC")
						contentWindow.Label(tx.Type, "LC")
						contentWindow.Label(tx.Hash, "LC")
					}
				}
			} else {
				widgets.Run("loading", contentWindow.Window, func() {})
			}
			contentWindow.End()
		}
		pageWindow.End()
	}
}

func (handler *TransactionsHandler) fetchTransactions(wallet walletcore.Wallet, window *nucular.Window) {
	handler.transactions, handler.err = wallet.TransactionHistory()
	handler.hasFetchedTransactions = true
	window.Master().Changed()
}
