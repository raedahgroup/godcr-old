package pagehandlers

import (
	"context"
	"fmt"
	"strconv"

	"github.com/aarzilli/nucular"
	"github.com/decred/dcrd/dcrutil"
	"github.com/raedahgroup/godcr/app/walletcore"
	"github.com/raedahgroup/godcr/nuklear/styles"
	"github.com/raedahgroup/godcr/nuklear/widgets"
)

type HistoryHandler struct {
	err                    error
	ctx                    context.Context
	transactions           []*walletcore.Transaction
	isFetchingTransactions bool
	nextBlockHeight        int32

	transactionHash              string
	transactionDetails           *walletcore.TransactionDetails
	isFetchingTransactionDetails bool

	wallet walletcore.Wallet
}

func (handler *HistoryHandler) BeforeRender(wallet walletcore.Wallet, refreshWindowDisplay func()) bool {
	// todo: caller should ideally pass a context parameter, propagated from main.go
	handler.ctx = context.Background()

	handler.isFetchingTransactions = true
	handler.err = nil
	handler.transactions = nil

	handler.transactionHash = ""
	handler.transactionDetails = nil
	handler.isFetchingTransactionDetails = false

	handler.wallet = wallet

	go handler.fetchTransactions(wallet, refreshWindowDisplay)

	return true
}

func (handler *HistoryHandler) Render(window *nucular.Window) {
	if handler.transactionHash == "" {
		handler.renderHistoryPage(window)
		return
	}
	handler.renderTransactionDetailsPage(window)
}

func (handler *HistoryHandler) renderHistoryPage(window *nucular.Window) {
	widgets.PageContentWindowDefaultPadding("History", window, func(contentWindow *widgets.Window) {
		if handler.err != nil {
			contentWindow.DisplayErrorMessage("Error fetching txs", handler.err)
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
	historyTable := widgets.NewTable()

	// render table header with nav font
	historyTable.AddRowWithFont(styles.NavFont,
		widgets.NewLabelTableCell("#", "LC"),
		widgets.NewLabelTableCell("Date", "LC"),
		widgets.NewLabelTableCell("Direction", "LC"),
		widgets.NewLabelTableCell("Amount", "LC"),
		widgets.NewLabelTableCell("Fee", "LC"),
		widgets.NewLabelTableCell("Type", "LC"),
		widgets.NewLabelTableCell("Hash", "LC"),
	)

	for i, tx := range handler.transactions {
		historyTable.AddRow(
			widgets.NewLabelTableCell(fmt.Sprintf("%d", i+1), "LC"),
			widgets.NewLabelTableCell(tx.FormattedTime, "LC"),
			widgets.NewLabelTableCell(tx.Direction.String(), "LC"),
			widgets.NewLabelTableCell(tx.Amount, "RC"),
			widgets.NewLabelTableCell(tx.Fee, "RC"),
			widgets.NewLabelTableCell(tx.Type, "LC"),
			widgets.NewLinkTableCell(tx.Hash, "Click to see transaction details", handler.gotoTransactionDetails),
		)
	}

	historyTable.Render(contentWindow)
}

func (handler *HistoryHandler) gotoTransactionDetails(txHash string, window *widgets.Window) {
	handler.transactionHash = txHash
	window.Master().Changed()
}

func (handler *HistoryHandler) renderTransactionDetailsPage(window *nucular.Window) {
	if handler.transactionDetails == nil {
		handler.isFetchingTransactionDetails = true
		go func() {
			handler.transactionDetails, handler.err = handler.wallet.GetTransaction(handler.transactionHash)
			handler.isFetchingTransactionDetails = false
			window.Master().Changed()
		}()
	}

	widgets.PageContentWindowDefaultPadding("Transaction Detail", window, func(contentWindow *widgets.Window) {
		if handler.err != nil {
			contentWindow.DisplayErrorMessage("Error fetching transaction details", handler.err)
		} else if handler.transactionDetails != nil {
			handler.displayTransactionDetails(contentWindow)
		}

		// show loading indicator if tx details is being fetched
		if handler.isFetchingTransactionDetails {
			contentWindow.DisplayIsLoadingMessage()
		}
	})
}

func (handler *HistoryHandler) displayTransactionDetails(contentWindow *widgets.Window) {
	contentWindow.Window.Row(0).Static(670, 700)
	widgets.NoScrollGroupWindow("Transaction details column one", contentWindow.Window, func(window *widgets.Window) {
		table := widgets.NewTable()
		table.AddRow(
			widgets.NewLabelTableCell("Confirmations", "LC"),
			widgets.NewLabelTableCell(strconv.Itoa(int(handler.transactionDetails.Confirmations)), "LC"),
		)
		table.AddRow(
			widgets.NewLabelTableCell("Hash", "LC"),
			widgets.NewLabelTableCell(handler.transactionDetails.Hash, "LC"),
		)
		table.AddRow(
			widgets.NewLabelTableCell("Block Height", "LC"),
			widgets.NewLabelTableCell(strconv.Itoa(int(handler.transactionDetails.BlockHeight)), "LC"),
		)
		table.AddRow(
			widgets.NewLabelTableCell("Direction", "LC"),
			widgets.NewLabelTableCell(handler.transactionDetails.Direction.String(), "LC"),
		)
		table.AddRow(
			widgets.NewLabelTableCell("Type", "LC"),
			widgets.NewLabelTableCell(handler.transactionDetails.Type, "LC"),
		)
		table.Render(window)

		window.AddHorizontalSpace(50)

		window.AddLabelWithFont("Inputs", "LC", styles.BoldPageContentFont)

		table = widgets.NewTable()
		table.AddRowWithFont(styles.NavFont,
			widgets.NewLabelTableCell("Previous Outpoint", "LC"),
			widgets.NewLabelTableCell("Amount", "LC"),
		)

		for _, input := range handler.transactionDetails.Inputs {
			table.AddRow(
				widgets.NewLabelTableCell(input.PreviousTransactionHash, "LC"),
				widgets.NewLabelTableCell(dcrutil.Amount(input.AmountIn).String(), "LC"),
			)
		}
		table.Render(window)
	})

	widgets.NoScrollGroupWindow("Transaction details column two", contentWindow.Window, func(window *widgets.Window) {
		table := widgets.NewTable()
		table.AddRow(
			widgets.NewLabelTableCell("Amount", "LC"),
			widgets.NewLabelTableCell(handler.transactionDetails.Amount, "LC"),
		)
		table.AddRow(
			widgets.NewLabelTableCell("Size", "LC"),
			widgets.NewLabelTableCell(strconv.Itoa(handler.transactionDetails.Size)+" Bytes", "LC"),
		)
		table.AddRow(
			widgets.NewLabelTableCell("Fee", "LC"),
			widgets.NewLabelTableCell(handler.transactionDetails.Fee, "LC"),
		)
		table.AddRow(
			widgets.NewLabelTableCell("Fee Rate", "LC"),
			widgets.NewLabelTableCell(handler.transactionDetails.FeeRate.String(), "LC"),
		)
		table.AddRow(
			widgets.NewLabelTableCell("Time", "LC"),
			widgets.NewLabelTableCell(handler.transactionDetails.FormattedTime, "LC"),
		)
		table.Render(window)

		window.AddHorizontalSpace(50)

		window.AddLabelWithFont("Outputs", "LC", styles.BoldPageContentFont)

		table = widgets.NewTable()
		table.AddRowWithFont(styles.NavFont,
			widgets.NewLabelTableCell("Address", "LC"),
			widgets.NewLabelTableCell("Account", "LC"),
			widgets.NewLabelTableCell("Value", "LC"),
			widgets.NewLabelTableCell("Type", "LC"),
		)

		for _, output := range handler.transactionDetails.Outputs {
			for _, address := range output.Addresses {
				account := address.AccountName
				if !address.IsMine {
					account = "external address"
				}

				table.AddRow(
					widgets.NewLabelTableCell(address.Address, "LC"),
					widgets.NewLabelTableCell(account, "LC"),
					widgets.NewLabelTableCell(dcrutil.Amount(output.Value).String(), "LC"),
					widgets.NewLabelTableCell(output.ScriptType, "LC"),
				)
			}
		}
		table.Render(window)
	})
}
