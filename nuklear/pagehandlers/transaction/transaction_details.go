package transaction

import (
	"strconv"

	"github.com/aarzilli/nucular"
	"github.com/decred/dcrd/dcrutil"
	"github.com/raedahgroup/godcr/nuklear/styles"
	"github.com/raedahgroup/godcr/nuklear/widgets"
)

func (handler *HistoryHandler) clearTxDetails() {
	handler.selectedTxHash = ""
	handler.selectedTxDetails = nil
	handler.isFetchingTxDetails = false
	handler.fetchTxDetailsError = nil
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
