package pagehandlers

import (
	"fmt"
	"image/color"
	"strconv"

	"github.com/aarzilli/nucular"
	"github.com/decred/dcrd/chaincfg/chainhash"
	"github.com/decred/dcrd/dcrutil"
	"github.com/raedahgroup/dcrlibwallet"
	"github.com/raedahgroup/godcr/nuklear/styles"
	"github.com/raedahgroup/godcr/nuklear/widgets"
)

const (
	dividerHeight = 10
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
			hash, _ := chainhash.NewHashFromStr(handler.selectedTxHash)
			handler.selectedTxDetails, handler.fetchTxDetailsError = handler.wallet.GetTransactionRaw(hash[:])
			handler.isFetchingTxDetails = false
			window.Master().Changed()
		}()
	}

	widgets.PageContentWindowDefaultPadding("Transaction Details", window, func(contentWindow *widgets.Window) {
		if handler.fetchTxDetailsError != nil {
			contentWindow.DisplayErrorMessage("Error fetching transaction details", handler.fetchTxDetailsError)
		} else if handler.selectedTxDetails != nil {
			handler.displayTransactionDetails(contentWindow)
		} else if handler.isFetchingTxDetails {
			contentWindow.DisplayIsLoadingMessage()
		}
	})
}

func (handler *HistoryHandler) displayTransactionDetails(contentWindow *widgets.Window) {
	if handler.selectedTxHash == "" {
		return
	}

	var confirmations int32 = 0
	if handler.selectedTxDetails.BlockHeight != -1 {
		confirmations = handler.wallet.GetBestBlock() - handler.selectedTxDetails.BlockHeight + 1
	}

	var spendUnconfirmed = handler.wallet.ReadBoolConfigValueForKey(dcrlibwallet.SpendUnconfirmedConfigKey)

	var status string
	var statusColor color.RGBA
	if spendUnconfirmed || confirmations > dcrlibwallet.DefaultRequiredConfirmations {
		status = "Confirmed"
		statusColor = styles.DecredGreenColor
	} else {
		status = "Pending"
		statusColor = styles.DecredOrangeColor
	}

	// we create our tables here so that we are able to calculate our window height using table data
	txDetailsTable := widgets.NewTable()
	txDetailsTable.AddRow(
		widgets.NewLabelTableCell("Confirmations", "LC"),
		widgets.NewLabelTableCell(strconv.Itoa(int(confirmations)), "LC"),
	)
	txDetailsTable.AddRow(
		widgets.NewLabelTableCell("Hash", "LC"),
		widgets.NewLabelTableCell(handler.selectedTxDetails.Hash, "LC"),
	)
	txDetailsTable.AddRow(
		widgets.NewLabelTableCell("Block Height", "LC"),
		widgets.NewLabelTableCell(strconv.Itoa(int(handler.selectedTxDetails.BlockHeight)), "LC"),
	)
	txDetailsTable.AddRow(
		widgets.NewLabelTableCell("Direction", "LC"),
		widgets.NewLabelTableCell(dcrlibwallet.TransactionDirectionName(handler.selectedTxDetails.Direction), "LC"),
	)
	txDetailsTable.AddRow(
		widgets.NewLabelTableCell("Type", "LC"),
		widgets.NewLabelTableCell(handler.selectedTxDetails.Type, "LC"),
	)
	txDetailsTable.AddRow(
		widgets.NewLabelTableCell("Amount", "LC"),
		widgets.NewLabelTableCell(dcrutil.Amount(handler.selectedTxDetails.Amount).String(), "LC"),
	)
	txDetailsTable.AddRow(
		widgets.NewLabelTableCell("Size", "LC"),
		widgets.NewLabelTableCell(strconv.Itoa(handler.selectedTxDetails.Size)+" Bytes", "LC"),
	)
	txDetailsTable.AddRow(
		widgets.NewLabelTableCell("Fee", "LC"),
		widgets.NewLabelTableCell(dcrutil.Amount(handler.selectedTxDetails.Fee).String(), "LC"),
	)
	txDetailsTable.AddRow(
		widgets.NewLabelTableCell("Fee Rate", "LC"),
		widgets.NewLabelTableCell(dcrutil.Amount(handler.selectedTxDetails.FeeRate).String(), "LC"),
	)
	txDetailsTable.AddRow(
		widgets.NewLabelTableCell("Status", "LC"),
		widgets.NewColoredLabelTableCell(status, "LC", statusColor),
	)
	txDetailsTable.AddRow(
		widgets.NewLabelTableCell("Date", "LC"),
		widgets.NewLabelTableCell(fmt.Sprintf("%s UTC", dcrlibwallet.FormatUTCTime(handler.selectedTxDetails.Timestamp)), "LC"),
	)

	txInputsTable := widgets.NewTable()
	txInputsTable.AddRowWithFont(styles.NavFont,
		widgets.NewLabelTableCell("Previous Outpoint", "LC"),
		widgets.NewLabelTableCell("Account", "LC"),
		widgets.NewLabelTableCell("Amount", "LC"),
	)

	for _, input := range handler.selectedTxDetails.Inputs {
		txInputsTable.AddRow(
			widgets.NewLabelTableCell(input.PreviousTransactionHash, "LC"),
			widgets.NewLabelTableCell(input.AccountName, "LC"),
			widgets.NewLabelTableCell(dcrutil.Amount(input.Amount).String(), "LC"),
		)
	}

	txOutputsTable := widgets.NewTable()
	txOutputsTable.AddRowWithFont(styles.NavFont,
		widgets.NewLabelTableCell("Address", "LC"),
		widgets.NewLabelTableCell("Account", "LC"),
		widgets.NewLabelTableCell("Value", "LC"),
		widgets.NewLabelTableCell("Type", "LC"),
	)

	for _, output := range handler.selectedTxDetails.Outputs {
		txOutputsTable.AddRow(
			widgets.NewLabelTableCell(output.Address, "LC"),
			widgets.NewLabelTableCell(output.AccountName, "LC"),
			widgets.NewLabelTableCell(dcrutil.Amount(output.Amount).String(), "LC"),
			widgets.NewLabelTableCell(output.ScriptType, "LC"),
		)
	}

	// calculate additionally used horizontal space
	hSpace := (dividerHeight * 2) + (widgets.TableRowHeight * 3) // 2 horizontal spaces + 3 lines of text (2 table headers, 1 breadcrumb)

	contentWindow.Window.Row(handler.calculateTxDetailsPageHeight(txDetailsTable.Height(), txInputsTable.Height(), txOutputsTable.Height(), hSpace)).Static(730)
	widgets.NoScrollGroupWindow("tx-details-group-1", contentWindow.Window, func(window *widgets.Window) {
		breadcrumb := []*widgets.Breadcrumb{
			{
				Text: "History",
				Action: func(text string, window *widgets.Window) {
					handler.clearTxDetails()
					window.Master().Changed()
				},
			},
			{
				Text:   "Transaction Details",
				Action: nil,
			},
		}
		window.AddBreadcrumb(breadcrumb)

		txDetailsTable.Render(window)

		window.AddHorizontalSpace(dividerHeight)
		window.AddLabelWithFont("Inputs", "LC", styles.BoldPageContentFont)
		txInputsTable.Render(window)

		window.AddHorizontalSpace(dividerHeight)
		window.AddLabelWithFont("Outputs", "LC", styles.BoldPageContentFont)

		txOutputsTable.Render(window)
	})
}

func (handler *HistoryHandler) calculateTxDetailsPageHeight(tableHeights ...int) int {
	var totalTableHeight int

	for i := range tableHeights {
		totalTableHeight += tableHeights[i]
	}

	return totalTableHeight
}
