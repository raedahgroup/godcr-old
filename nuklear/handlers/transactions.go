package handlers

import (
	"context"
	"github.com/aarzilli/nucular"
	"github.com/raedahgroup/godcr/app"
	"github.com/raedahgroup/godcr/app/walletcore"
	"github.com/raedahgroup/godcr/nuklear/handlers/widgets"
	"github.com/raedahgroup/godcr/nuklear/helpers"
)

type HistoryHandler struct {
	err                    error
	ctx                    context.Context
	isRendering            bool
	transactions           []*walletcore.Transaction
	nextBlockHeight        int32
	hasFetchedTransactions bool
}

func (handler *HistoryHandler) BeforeRender() {
	handler.err = nil
	handler.ctx = context.Background()
	handler.nextBlockHeight = -1
	handler.transactions = nil
	handler.isRendering = false
	handler.hasFetchedTransactions = false
}

func (handler *HistoryHandler) Render(window *nucular.Window, walletMiddleware app.WalletMiddleware) {
	if !handler.isRendering {
		handler.isRendering = true
		go handler.fetchTransactions(walletMiddleware, window)
	}

	// draw page
	if pageWindow := helpers.NewWindow("History Page", window, 0); pageWindow != nil {
		pageWindow.DrawHeader("History")

		// content window
		if contentWindow := pageWindow.ContentWindow("History"); contentWindow != nil {
			if handler.hasFetchedTransactions {
				if handler.err != nil {
					contentWindow.SetErrorMessage(handler.err.Error())
				} else {
					helpers.SetFont(window, helpers.NavFont)
					contentWindow.Row(20).Ratio(0.18, 0.08, 0.15, 0.15, 0.15, 0.7)
					contentWindow.Label("Date", "LC")
					contentWindow.Label("Amount", "LC")
					contentWindow.Label("Fee", "LC")
					contentWindow.Label("Direction", "LC")
					contentWindow.Label("Type", "LC")
					contentWindow.Label("Hash", "LC")

					for _, tx := range handler.transactions {
						helpers.SetFont(window, helpers.PageContentFont)
						contentWindow.Row(20).Ratio(0.18, 0.08, 0.15, 0.15, 0.15, 0.7)
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

func (handler *HistoryHandler) fetchTransactions(wallet walletcore.Wallet, window *nucular.Window) {
	handler.transactions, handler.nextBlockHeight, handler.err = wallet.TransactionHistory(handler.ctx, handler.nextBlockHeight, walletcore.TransactionHistoryCountPerPage)
	handler.nextBlockHeight--
	handler.hasFetchedTransactions = true
	window.Master().Changed()
}
