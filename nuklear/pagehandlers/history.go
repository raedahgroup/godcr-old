package pagehandlers

import (
	"context"
	"fmt"

	"github.com/aarzilli/nucular"
	"github.com/raedahgroup/godcr/app/walletcore"
	"github.com/raedahgroup/godcr/nuklear/widgets"
	"github.com/raedahgroup/godcr/nuklear/styles"
)

type HistoryHandler struct {
	err                    error
	ctx                    context.Context
	transactions           []*walletcore.Transaction
	isFetchingTransactions bool
	nextBlockHeight        int32
}

func (handler *HistoryHandler) BeforeRender(wallet walletcore.Wallet, refreshWindowDisplay func()) {
	// todo: caller should ideally pass a context parameter, propagated from main.go
	handler.ctx = context.Background()

	handler.isFetchingTransactions = false
	handler.err = nil
	handler.transactions = nil

	go handler.fetchTransactions(wallet, refreshWindowDisplay)
}

func (handler *HistoryHandler) Render(window *nucular.Window) {
	widgets.PageContentWindowDefaultPadding("History", window, func(contentWindow *widgets.Window) {
		if handler.err != nil {
			contentWindow.DisplayErrorMessage(handler.err.Error())
		} else if len(handler.transactions) > 0 {
			handler.displayTransactions(contentWindow)
		}

		// show loading indicator if tx is being fetched
		if handler.isFetchingTransactions {
			contentWindow.DisplayIsLoadingMessage()
		}
	})
}

func (handler *HistoryHandler) fetchTransactions(wallet walletcore.Wallet, refreshWindowDisplay func()) {
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

	refreshWindowDisplay()

	// load more if possible
	if handler.nextBlockHeight >= 0 {
		handler.fetchTransactions(wallet, refreshWindowDisplay)
	} else {
		handler.isFetchingTransactions = false
	}
}

func (handler *HistoryHandler) displayTransactions(contentWindow *widgets.Window) {
	// render table header with nav font
	contentWindow.UseFontAndResetToPrevious(styles.NavFont, func() {
		contentWindow.Row(20).Static(25, 80, 60, 70, 70, 80, 300)
		contentWindow.Label("#", "LC")
		contentWindow.Label("Date", "LC")
		contentWindow.Label("Direction", "LC")
		contentWindow.Label("Amount", "LC")
		contentWindow.Label("Fee", "LC")
		contentWindow.Label("Type", "LC")
		contentWindow.Label("Hash", "LC")
	})

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