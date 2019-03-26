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
	err                     error
	ctx                     context.Context
	isRendering             bool
	transactions            []*walletcore.Transaction
	isFetchingTransactions  bool
	loadingNextPage         bool
	loadingPreviousPage     bool
	pagesBlockHeightHistory []int32
	pagesTxCount            []int
	nextBlockHeight         int32
}

func (handler *HistoryHandler) BeforeRender() {
	handler.isRendering = false
	handler.loadingNextPage = false
	handler.loadingPreviousPage = false
	handler.isFetchingTransactions = false
	handler.err = nil
	handler.transactions = nil
	handler.pagesBlockHeightHistory = nil
	handler.pagesTxCount = nil
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
			// show error or history table first
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

	if len(handler.pagesBlockHeightHistory) == 0 {
		// first page
		handler.nextBlockHeight = -1
	}

	// record current blockheight for use in previous/next page navigation
	handler.pagesBlockHeightHistory = append(handler.pagesBlockHeightHistory, handler.nextBlockHeight)

	// record current page count for use in correctly calculating s/n later when displaying the transactions
	// don't record current page count if we're navigating backwards since the current tx list will soon be discarded
	if !handler.loadingPreviousPage {
		handler.pagesTxCount = append(handler.pagesTxCount, len(handler.transactions))
	}

	var endBlockHeight int32
	handler.transactions, endBlockHeight, handler.err = wallet.TransactionHistory(handler.ctx, handler.nextBlockHeight,
		walletcore.TransactionHistoryCountPerPage)

	// next start block should be the block immediately preceding the current end block
	handler.nextBlockHeight = endBlockHeight - 1

	handler.isFetchingTransactions = false
	handler.loadingPreviousPage = false
	handler.loadingNextPage = false

	masterWindow.Changed()
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

	// s/n calculation
	totalTxDisplayed := 0
	for _, count := range handler.pagesTxCount {
		totalTxDisplayed += count
	}

	// render rows with content font
	helpers.SetFont(contentWindow.Window, helpers.PageContentFont)
	for i, tx := range handler.transactions {
		n := totalTxDisplayed + i + 1

		contentWindow.Row(20).Static(25, 80, 60, 70, 70, 80, 300)
		contentWindow.Label(fmt.Sprintf("%d", n), "LC")
		contentWindow.Label(tx.FormattedTime, "LC")
		contentWindow.Label(tx.Direction.String(), "LC")
		contentWindow.Label(tx.Amount.String(), "LC")
		contentWindow.Label(tx.Fee.String(), "LC")
		contentWindow.Label(tx.Type, "LC")
		contentWindow.Label(tx.Hash, "LC")
	}

	// row for previous and next buttons
	contentWindow.Row(25).Static(100, 75)

	// show previous page button only if this isn't first page
	if len(handler.pagesBlockHeightHistory) > 1 && contentWindow.ButtonText("Previous Page") {
		if handler.loadingPreviousPage || handler.loadingNextPage {
			return
		}

		lastPage := len(handler.pagesBlockHeightHistory) - 2
		handler.nextBlockHeight = handler.pagesBlockHeightHistory[lastPage]

		handler.pagesBlockHeightHistory = handler.pagesBlockHeightHistory[:lastPage]
		handler.pagesTxCount = handler.pagesTxCount[:lastPage+1]

		handler.loadingPreviousPage = true
		handler.fetchTransactions(wallet, contentWindow.Master())
	}

	// only show next page button if there's more data to load
	if handler.nextBlockHeight >= 0 && contentWindow.ButtonText("Next Page") {
		if handler.loadingNextPage || handler.loadingPreviousPage {
			return
		}

		handler.loadingNextPage = true
		handler.fetchTransactions(wallet, contentWindow.Master())
	}
}
