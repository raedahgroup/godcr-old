package pages

// import (
// 	// "fmt"
// 	// "strings"
// 	// "strconv"

// 	"fyne.io/fyne"
// 	// "fyne.io/fyne/widget"
// 	// "fyne.io/fyne/layout"

// 	// "github.com/raedahgroup/godcr/fyne/widgets"
// )

func displayTransactionDetails() {
	// var confirmations int32 = 0
	// if handler.selectedTxDetails.BlockHeight != -1 {
	// 	confirmations = handler.wallet.GetBestBlock() - handler.selectedTxDetails.BlockHeight + 1
	// }

	// var spendUnconfirmed = handler.wallet.ReadBoolConfigValueForKey(dcrlibwallet.SpendUnconfirmedConfigKey)

	// var status string
	// // var statusColor color.RGBA
	// if spendUnconfirmed || confirmations > dcrlibwallet.DefaultRequiredConfirmations {
	// 	status = "Confirmed"
	// 	// statusColor = styles.DecredGreenColor
	// } else {
	// 	status = "Pending"
	// 	// statusColor = styles.DecredOrangeColor
	// }

	// tableConfirmations := widget.NewHBox(
	// 	widget.NewLabelWithStyle("Confirmations", fyne.TextAlignCenter, fyne.TextStyle{Bold: true}),
	// 	// widget.NewLabelWithStyle(strconv.Itoa(int(confirmations)), fyne.TextAlignCenter, fyne.TextStyle{Bold: true}),
	// )
	// tableHash := widget.NewHBox(
	// 	widget.NewLabelWithStyle("Hash", fyne.TextAlignCenter, fyne.TextStyle{Bold: true}),
	// 	// widget.NewLabelWithStyle(handler.selectedTxDetails.Hash, fyne.TextAlignCenter, fyne.TextStyle{Bold: true}),
	// )
	// tableBlockHeight := widget.NewHBox(
	// 	widget.NewLabelWithStyle("Block Height", fyne.TextAlignCenter, fyne.TextStyle{Bold: true}),
	// 	// widget.NewLabelWithStyle(strconv.Itoa(int(handler.selectedTxDetails.BlockHeight)), fyne.TextAlignCenter, fyne.TextStyle{Bold: true}),
	// )
	// tableDirection := widget.NewHBox(
	// 	widget.NewLabelWithStyle("Direction", fyne.TextAlignCenter, fyne.TextStyle{Bold: true}),
	// 	// widget.NewLabelWithStyle(dcrlibwallet.TransactionDirectionName(handler.selectedTxDetails.Direction), fyne.TextAlignCenter, fyne.TextStyle{Bold: true}),
	// )
	// tableType := widget.NewHBox(
	// 	widget.NewLabelWithStyle("Type", fyne.TextAlignCenter, fyne.TextStyle{Bold: true}),
	// 	// widget.NewLabelWithStyle(handler.selectedTxDetails.Type, fyne.TextAlignCenter, fyne.TextStyle{Bold: true}),
	// )
	// tableAmount := widget.NewHBox(
	// 	widget.NewLabelWithStyle("Amount", fyne.TextAlignCenter, fyne.TextStyle{Bold: true}),
	// 	// widget.NewLabelWithStyle(dcrutil.Amount(handler.selectedTxDetails.Amount).String(), fyne.TextAlignCenter, fyne.TextStyle{Bold: true}),
	// )
	// tableSize := widget.NewHBox(
	// 	widget.NewLabelWithStyle("Size", fyne.TextAlignCenter, fyne.TextStyle{Bold: true}),
	// 	// widget.NewLabelWithStyle(strconv.Itoa(handler.selectedTxDetails.Size)+" Bytes", fyne.TextAlignCenter, fyne.TextStyle{Bold: true}),
	// )
	// tableFee := widget.NewHBox(
	// 	widget.NewLabelWithStyle("Fee", fyne.TextAlignCenter, fyne.TextStyle{Bold: true}),
	// 	// widget.NewLabelWithStyle(dcrutil.Amount(handler.selectedTxDetails.Fee).String(), fyne.TextAlignCenter, fyne.TextStyle{Bold: true}),
	// )
	// tableFeeRate := widget.NewHBox(
	// 	widget.NewLabelWithStyle("Fee Rate", fyne.TextAlignCenter, fyne.TextStyle{Bold: true}),
	// 	// widget.NewLabelWithStyle(dcrutil.Amount(handler.selectedTxDetails.FeeRate).String(), fyne.TextAlignCenter, fyne.TextStyle{Bold: true}),
	// )
	// tableStatus := widget.NewHBox(
	// 	widget.NewLabelWithStyle("Status", fyne.TextAlignCenter, fyne.TextStyle{Bold: true}),
	// 	// widget.NewLabelWithStyle(status, fyne.TextAlignCenter, fyne.TextStyle{Bold: true}),
	// )
	// tableDate := widget.NewHBox(
	// 	widget.NewLabelWithStyle("Date", fyne.TextAlignCenter, fyne.TextStyle{Bold: true}),
	// 	// widget.NewLabelWithStyle(fmt.Sprintf("%s UTC", dcrlibwallet.FormatUTCTime(handler.selectedTxDetails.Timestamp)), fyne.TextAlignCenter, fyne.TextStyle{Bold: true}),
	// )

	// tableData := widget.NewVBox(
	// 	tableConfirmations,
	// 	tableHash,
	// 	tableBlockHeight,
	// 	tableDirection,
	// 	tableType,
	// 	tableAmount,
	// 	tableSize,
	// 	tableFee,
	// 	tableFeeRate,
	// 	tableStatus,
	// 	tableDate,
	// )

	// var txInput widgets.TableStruct
	// heading := widget.NewHBox(
	// 	widget.NewLabelWithStyle("Previous Outpoint", fyne.TextAlignCenter, fyne.TextStyle{Bold: true}),
	// 	widget.NewLabelWithStyle("Account", fyne.TextAlignCenter, fyne.TextStyle{Bold: true}),
	// 	widget.NewLabelWithStyle("Amount", fyne.TextAlignCenter, fyne.TextStyle{Bold: true}))

	// var inputBox []*widget.Box
	// for i := range txDetails.Inputs {
	// 	inputBox = append(inputBox, widget.NewHBox(
	// 		widget.NewLabelWithStyle(txDetails.Inputs[i].PreviousOutpoint, fyne.TextAlignLeading, fyne.TextStyle{}),
	// 		widget.NewLabelWithStyle(txDetails.Inputs[i].AccountName, fyne.TextAlignCenter, fyne.TextStyle{}),
	// 		widget.NewLabelWithStyle(dcrutil.Amount(txDetails.Inputs[i].Amount).String(), fyne.TextAlignTrailing, fyne.TextStyle{}),
	// 	))
	// }
	// txInput.NewTable(heading, inputBox...)

	// var txOutput widgets.TableStruct
	// heading = widget.NewHBox(
	// 	widget.NewLabelWithStyle("Address", fyne.TextAlignCenter, fyne.TextStyle{Bold: true}),
	// 	widget.NewLabelWithStyle("Account", fyne.TextAlignCenter, fyne.TextStyle{Bold: true}),
	// 	widget.NewLabelWithStyle("Value", fyne.TextAlignCenter, fyne.TextStyle{Bold: true}),
	// 	widget.NewLabelWithStyle("Type", fyne.TextAlignCenter, fyne.TextStyle{Bold: true}))

	// var outputBox []*widget.Box
	// for i := range txDetails.Outputs {
	// 	outputBox = append(outputBox, widget.NewHBox(
	// 		widget.NewLabelWithStyle(txDetails.Outputs[i].Address, fyne.TextAlignLeading, fyne.TextStyle{}),
	// 		widget.NewLabelWithStyle(txDetails.Outputs[i].AccountName, fyne.TextAlignCenter, fyne.TextStyle{}),
	// 		widget.NewLabelWithStyle(dcrutil.Amount(txDetails.Outputs[i].Amount).String(), fyne.TextAlignTrailing, fyne.TextStyle{}),
	// 		widget.NewLabelWithStyle(txDetails.Outputs[i].ScriptType, fyne.TextAlignCenter, fyne.TextStyle{})))
	// }
	// txOutput.NewTable(heading, outputBox...)

	// output := widget.NewHBox(widgets.NewHSpacer(10), widget.NewVBox(
	// 	form,
	// 	widget.NewLabelWithStyle("Inputs", fyne.TextAlignLeading, fyne.TextStyle{Bold: true}),
	// 	txInput.Result,
	// 	widgets.NewVSpacer(10),
	// 	widget.NewLabelWithStyle("Outputs", fyne.TextAlignLeading, fyne.TextStyle{Bold: true}),
	// 	txOutput.Result,
	// ), widgets.NewHSpacer(10))

	// scrollContainer := widget.NewScrollContainer(output)
	// scrollContainer.Resize(fyne.NewSize(scrollContainer.MinSize().Width, 500))
	// popUp := widget.NewPopUp(widget.NewVBox(label, fyne.NewContainer(scrollContainer)),
	// 	window.Canvas())
	// popUp.Show()
	return
}