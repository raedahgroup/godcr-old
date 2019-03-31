package handlers

import (
	"context"
	"fmt"
	"github.com/aarzilli/nucular"
	"github.com/raedahgroup/godcr/app/walletcore"
	"github.com/raedahgroup/godcr/nuklear/handlers/widgets"
	"github.com/raedahgroup/godcr/nuklear/helpers"
)

type HistoryHandler struct {
	err                    error
	ctx                    context.Context
	isRendering            bool
	transactions           []*walletcore.Transaction
	isFetchingTransactions bool
	nextBlockHeight        int32
}

func (handler *HistoryHandler) BeforeRender() {
	handler.isRendering = false
	handler.isFetchingTransactions = false
	handler.err = nil
	handler.transactions = nil
}

func (handler *HistoryHandler) Render(window *nucular.Window, wallet walletcore.Wallet) {
	// todo: caller should ideally pass a context parameter, propagated from main.go
	handler.ctx = context.Background()

	if !handler.isRendering {
		handler.isRendering = true
		go handler.fetchTransactions(wallet, window.Master())
	}

	// draw page
	if pageWindow := helpers.NewWindow("History Page", window, 0); pageWindow != nil {
		pageWindow.DrawHeader("History")

		// content window
		if contentWindow := pageWindow.ContentWindow("History"); contentWindow != nil {
			// show error or history table before loading indicator
			if handler.err != nil {
				contentWindow.SetErrorMessage(handler.err.Error())
			} else if len(handler.transactions) > 0 {
				handler.displayTransactions(contentWindow, wallet)
			}

			// show loading indicator if tx is being fetched
			if handler.isFetchingTransactions {
				widgets.ShowLoadingWidget(contentWindow.Window)
			}

			contentWindow.End()
		}
		pageWindow.End()
	}
}

func (handler *HistoryHandler) fetchTransactions(wallet walletcore.Wallet, masterWindow nucular.MasterWindow) {
	handler.isFetchingTransactions = true

	if len(handler.transactions) == 0 {
		// first page
		handler.nextBlockHeight = -1
	}

	transactions, endBlockHeight, err := wallet.TransactionHistory(handler.ctx, handler.nextBlockHeight,
		walletcore.TransactionHistoryCountPerPage)

	// next start block should be the block immediately preceding the current end block
	handler.err = err
	handler.transactions = append(handler.transactions, transactions...)
	handler.nextBlockHeight = endBlockHeight - 1

	masterWindow.Changed()

	// load more if possible
	if handler.nextBlockHeight >= 0 {
		handler.fetchTransactions(wallet, masterWindow)
	} else {
		handler.isFetchingTransactions = false
	}
}

func (handler *HistoryHandler) displayTransactions(contentWindow *helpers.Window, wallet walletcore.Wallet) {
	// render table header with nav font
	helpers.SetFont(contentWindow.Window, helpers.NavFont)
	contentWindow.Row(20).Static(25, 80, 60, 70, 70, 80, 300)
	contentWindow.Label("#", "LC")
	contentWindow.Label("Date", "LC")
	contentWindow.Label("Direction", "LC")
	contentWindow.Label("Amount", "LC")
	contentWindow.Label("Fee", "LC")
	contentWindow.Label("Type", "LC")
	contentWindow.Label("Hash", "LC")

	// render rows with content font
	helpers.SetFont(contentWindow.Window, helpers.PageContentFont)

	for i, tx := range handler.transactions {
		contentWindow.Row(15).Static(25, 80, 60, 70, 70, 80, 300)
		contentWindow.Label(fmt.Sprintf("%d", i+1), "LC")
		contentWindow.Label(tx.FormattedTime, "LC")
		contentWindow.Label(tx.Direction.String(), "LC")
		contentWindow.Label(tx.Amount, "RC")
		contentWindow.Label(tx.Fee, "RC")
		contentWindow.Label(tx.Type, "LC")
		contentWindow.Label(tx.Hash, "LC")
	}
}
