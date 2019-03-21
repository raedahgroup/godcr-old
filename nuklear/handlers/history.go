package handlers

import (
	"github.com/aarzilli/nucular"
	"github.com/raedahgroup/godcr/app"
	"github.com/raedahgroup/godcr/app/walletcore"
	"github.com/raedahgroup/godcr/nuklear/handlers/widgets"
	"github.com/raedahgroup/godcr/nuklear/helpers"
)

type HistoryHandler struct {
	err                    error
	isRendering            bool
	transactions           []*walletcore.Transaction
	hasFetchedTransactions bool
}

func (handler *HistoryHandler) BeforeRender() {
	handler.err = nil
	handler.transactions = nil
	handler.isRendering = false
	handler.hasFetchedTransactions = false
}

func (handler *HistoryHandler) Render(window *nucular.Window, walletMiddleware app.WalletMiddleware) {
	if !handler.isRendering {
		handler.isRendering = true
		go handler.fetchTransactions(walletMiddleware, window)
	}

	if pageWindow := helpers.NewWindow("History Page", window, nucular.WindowNoScrollbar); pageWindow != nil {
		pageWindow.DrawHeader("History")

		if contentWindow := pageWindow.ContentWindow("History"); contentWindow != nil {
			if handler.hasFetchedTransactions {
				if handler.err != nil {
					contentWindow.SetErrorMessage(handler.err.Error())
				} else {
					helpers.SetFont(window, helpers.NavFont)
					contentWindow.Row(20).Static(100, 70, 70, 60, 50, 280)
					contentWindow.Label("Date", "LC")
					contentWindow.Label("Amount", "LC")
					contentWindow.Label("Fee", "LC")
					contentWindow.Label("Direction", "LC")
					contentWindow.Label("Type", "LC")
					contentWindow.Label("Hash", "LC")

					for _, tx := range handler.transactions {
						helpers.SetFont(window, helpers.PageContentFont)
						contentWindow.Label(tx.FormattedTime, "LC")
						contentWindow.Label(tx.Amount.String(), "LC")
						contentWindow.Label(tx.Fee.String(), "LC")
						contentWindow.Label(tx.Direction.String(), "LC")
						contentWindow.Label(tx.Type, "LC")
						contentWindow.Label(tx.Hash, "LC")
					}
				}
			} else {
				widgets.ShowLoadingWidget(contentWindow.Window)
			}
			contentWindow.End()
		}
		pageWindow.End()
	}
}

func (handler *HistoryHandler) fetchTransactions(walletMiddleware app.WalletMiddleware, window *nucular.Window) {
	handler.transactions, handler.err = walletMiddleware.TransactionHistory()
	handler.hasFetchedTransactions = true
	window.Master().Changed()
}
