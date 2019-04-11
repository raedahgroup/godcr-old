package pagehandlers

import (
	"context"

	"github.com/aarzilli/nucular"
	"github.com/raedahgroup/godcr/app/walletcore"
	"github.com/raedahgroup/godcr/nuklear/styles"
	"github.com/raedahgroup/godcr/nuklear/widgets"
)

type HistoryHandler struct {
	err                    error
	isRendering            bool
	transactions           []*walletcore.Transaction
	hasFetchedTransactions bool

	// pagination params
	startBlockHeight int32
	endBlockHeight   int32
}

func (handler *HistoryHandler) BeforeRender() {
	handler.err = nil
	handler.transactions = nil
	handler.isRendering = false
	handler.hasFetchedTransactions = false

	handler.startBlockHeight = -1
	handler.endBlockHeight = 0
}

func (handler *HistoryHandler) Render(window *nucular.Window, wallet walletcore.Wallet) {
	if !handler.isRendering {
		handler.isRendering = true
		go handler.fetchHistory(wallet, window)
	}

	widgets.PageContentWindowDefaultPadding("History", window, func(contentWindow *widgets.Window) {
		if handler.hasFetchedTransactions {
			if handler.err != nil {
				contentWindow.DisplayErrorMessage(handler.err.Error())
			} else {
				contentWindow.UseFontAndResetToPrevious(styles.NavFont, func() {
					contentWindow.Row(styles.LabelHeight).Static(100, 60, 70, 70, 40, 200)
					contentWindow.Label("Date", "LC")
					contentWindow.Label("Direction", "LC")
					contentWindow.Label("Amount", "LC")
					contentWindow.Label("Fee", "LC")
					contentWindow.Label("Type", "LC")
					contentWindow.Label("Hash", "LC")
				})

				for _, tx := range handler.transactions {
					contentWindow.Label(tx.FormattedTime, "LC")
					contentWindow.Label(tx.Direction.String(), "LC")
					contentWindow.Label(tx.Amount, "LC")
					contentWindow.Label(tx.Fee, "LC")
					contentWindow.Label(tx.Type, "LC")
					contentWindow.Label(tx.Hash, "LC")
				}
			}
		} else {
			widgets.ShowLoadingWidget(contentWindow.Window)
		}
	})
}

func (handler *HistoryHandler) fetchHistory(wallet walletcore.Wallet, window *nucular.Window) {
	handler.transactions, handler.endBlockHeight, handler.err = wallet.TransactionHistory(context.Background(), handler.startBlockHeight, walletcore.TransactionHistoryCountPerPage)
	handler.hasFetchedTransactions = true
	window.Master().Changed()
}
