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
	fetchHistoryError      error
	ctx                    context.Context
	transactions           []*walletcore.Transaction
	isFetchingTransactions bool
	nextBlockHeight        int32

	selectedTxHash      string
	selectedTxDetails   *walletcore.TransactionDetails
	isFetchingTxDetails bool
	fetchTxDetailsError error

	wallet walletcore.Wallet
}

func (handler *HistoryHandler) BeforeRender(wallet walletcore.Wallet, refreshWindowDisplay func()) bool {
	// todo: caller should ideally pass a context parameter, propagated from main.go
	handler.ctx = context.Background()

	handler.isFetchingTransactions = true
	handler.fetchHistoryError = nil
	handler.transactions = nil

	handler.wallet = wallet

	go handler.fetchTransactions(wallet, refreshWindowDisplay)

	handler.clearTxDetails()

	return true
}

func (handler *HistoryHandler) fetchTransactions(wallet walletcore.Wallet, refreshWindowDisplay func()) {
	if len(handler.transactions) == 0 {
		// first page
		handler.nextBlockHeight = -1
	}

	transactions, endBlockHeight, err := wallet.TransactionHistory(handler.ctx, handler.nextBlockHeight,
		walletcore.TransactionHistoryCountPerPage)

	// next start block should be the block immediately preceding the current end block
	handler.fetchHistoryError = err
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

func (handler *HistoryHandler) clearTxDetails() {
	handler.selectedTxHash = ""
	handler.selectedTxDetails = nil
	handler.isFetchingTxDetails = false
	handler.fetchTxDetailsError = nil
}

func (handler *HistoryHandler) Render(window *nucular.Window) {
	if handler.selectedTxHash == "" {
		handler.renderHistoryPage(window)
		return
	}
	handler.renderTransactionDetailsPage(window)
}

func (handler *HistoryHandler) renderHistoryPage(window *nucular.Window) {
	widgets.PageContentWindowDefaultPadding("History", window, func(contentWindow *widgets.Window) {
		if handler.fetchHistoryError != nil {
			contentWindow.DisplayErrorMessage("Error fetching txs", handler.fetchHistoryError)
		} else if len(handler.transactions) > 0 {
			handler.displayTransactions(contentWindow)
		}

		// show loading indicator if tx is being fetched
		if handler.isFetchingTransactions {
			contentWindow.DisplayIsLoadingMessage()
		}
	})
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
	handler.selectedTxHash = txHash
	window.Master().Changed()
}

func (handler *HistoryHandler) renderTransactionDetailsPage(window *nucular.Window) {
	if handler.selectedTxDetails == nil {
		handler.isFetchingTxDetails = true
		go func() {
			handler.selectedTxDetails, handler.fetchTxDetailsError = handler.wallet.GetTransaction(handler.selectedTxHash)
			handler.isFetchingTxDetails = false
			window.Master().Changed()
		}()
	}

	widgets.PageContentWindowDefaultPadding("Transaction Details", window, func(contentWindow *widgets.Window) {
		contentWindow.AddButton("Back", func() {
			handler.goBackToHistory(contentWindow)
		})

		if handler.fetchTxDetailsError != nil {
			contentWindow.DisplayErrorMessage("Error fetching transaction details", handler.fetchTxDetailsError)
		} else if handler.selectedTxDetails != nil {
			handler.displayTransactionDetails(contentWindow)
		} else if handler.isFetchingTxDetails {
			contentWindow.DisplayIsLoadingMessage()
		}
	})
}

func (handler *HistoryHandler) goBackToHistory(contentWindow *widgets.Window) {
	handler.clearTxDetails()
	contentWindow.Master().Changed()
}

func (handler *HistoryHandler) displayTransactionDetails(contentWindow *widgets.Window) {
	// Create row to hold tx details in 2 columns
	// Each column will display data about the tx in a group window.
	// Row height is calculated based on the max group items total height
	contentWindow.Window.Row(handler.calculateTxDetailsPageHeight()).Static(670, 700)
	widgets.NoScrollGroupWindow("tx-details-col-1", contentWindow.Window, func(window *widgets.Window) {
		txDetailsTable1 := widgets.NewTable()
		txDetailsTable1.AddRow(
			widgets.NewLabelTableCell("Confirmations", "LC"),
			widgets.NewLabelTableCell(strconv.Itoa(int(handler.selectedTxDetails.Confirmations)), "LC"),
		)
		txDetailsTable1.AddRow(
			widgets.NewLabelTableCell("Hash", "LC"),
			widgets.NewLabelTableCell(handler.selectedTxDetails.Hash, "LC"),
		)
		txDetailsTable1.AddRow(
			widgets.NewLabelTableCell("Block Height", "LC"),
			widgets.NewLabelTableCell(strconv.Itoa(int(handler.selectedTxDetails.BlockHeight)), "LC"),
		)
		txDetailsTable1.AddRow(
			widgets.NewLabelTableCell("Direction", "LC"),
			widgets.NewLabelTableCell(handler.selectedTxDetails.Direction.String(), "LC"),
		)
		txDetailsTable1.AddRow(
			widgets.NewLabelTableCell("Type", "LC"),
			widgets.NewLabelTableCell(handler.selectedTxDetails.Type, "LC"),
		)
		txDetailsTable1.Render(window)

		window.AddHorizontalSpace(30)

		window.AddLabelWithFont("Inputs", "LC", styles.BoldPageContentFont)

		txInputsTable := widgets.NewTable()
		txInputsTable.AddRowWithFont(styles.NavFont,
			widgets.NewLabelTableCell("Previous Outpoint", "LC"),
			widgets.NewLabelTableCell("Amount", "LC"),
		)

		for _, input := range handler.selectedTxDetails.Inputs {
			txInputsTable.AddRow(
				widgets.NewLabelTableCell(input.PreviousTransactionHash, "LC"),
				widgets.NewLabelTableCell(dcrutil.Amount(input.AmountIn).String(), "LC"),
			)
		}
		txInputsTable.Render(window)
	})

	widgets.NoScrollGroupWindow("tx-details-col-2", contentWindow.Window, func(window *widgets.Window) {
		txDetailsTable2 := widgets.NewTable()
		txDetailsTable2.AddRow(
			widgets.NewLabelTableCell("Amount", "LC"),
			widgets.NewLabelTableCell(handler.selectedTxDetails.Amount, "LC"),
		)
		txDetailsTable2.AddRow(
			widgets.NewLabelTableCell("Size", "LC"),
			widgets.NewLabelTableCell(strconv.Itoa(handler.selectedTxDetails.Size)+" Bytes", "LC"),
		)
		txDetailsTable2.AddRow(
			widgets.NewLabelTableCell("Fee", "LC"),
			widgets.NewLabelTableCell(handler.selectedTxDetails.Fee, "LC"),
		)
		txDetailsTable2.AddRow(
			widgets.NewLabelTableCell("Fee Rate", "LC"),
			widgets.NewLabelTableCell(handler.selectedTxDetails.FeeRate.String(), "LC"),
		)
		txDetailsTable2.AddRow(
			widgets.NewLabelTableCell("Time", "LC"),
			widgets.NewLabelTableCell(handler.selectedTxDetails.FormattedTime, "LC"),
		)
		txDetailsTable2.Render(window)

		window.AddHorizontalSpace(30)

		window.AddLabelWithFont("Outputs", "LC", styles.BoldPageContentFont)

		txOutputsTable := widgets.NewTable()
		txOutputsTable.AddRowWithFont(styles.NavFont,
			widgets.NewLabelTableCell("Address", "LC"),
			widgets.NewLabelTableCell("Account", "LC"),
			widgets.NewLabelTableCell("Value", "LC"),
			widgets.NewLabelTableCell("Type", "LC"),
		)

		for _, output := range handler.selectedTxDetails.Outputs {
			for _, address := range output.Addresses {
				account := address.AccountName
				if !address.IsMine {
					account = "external address"
				}

				txOutputsTable.AddRow(
					widgets.NewLabelTableCell(address.Address, "LC"),
					widgets.NewLabelTableCell(account, "LC"),
					widgets.NewLabelTableCell(dcrutil.Amount(output.Value).String(), "LC"),
					widgets.NewLabelTableCell(output.ScriptType, "LC"),
				)
			}
		}
		txOutputsTable.Render(window)
	})
}

func (handler *HistoryHandler) calculateTxDetailsPageHeight() int {
	firstSectionHeight := 240
	outputsLen := len(handler.selectedTxDetails.Outputs)
	inputsLen := len(handler.selectedTxDetails.Inputs)

	var secondSectionLines int
	if outputsLen > inputsLen {
		secondSectionLines = outputsLen
	} else {
		secondSectionLines = inputsLen
	}

	return firstSectionHeight + (secondSectionLines * 27)
}
