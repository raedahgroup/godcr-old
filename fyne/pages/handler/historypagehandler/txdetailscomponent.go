package historypagehandler

import (
	"fmt"
	"net/url"
	"strconv"
	"strings"
	"time"

	"fyne.io/fyne"
	"fyne.io/fyne/canvas"
	"fyne.io/fyne/layout"
	"fyne.io/fyne/theme"
	"fyne.io/fyne/widget"

	"github.com/decred/dcrd/chaincfg/chainhash"
	"github.com/decred/dcrd/dcrutil"
	"github.com/raedahgroup/dcrlibwallet"
	"github.com/raedahgroup/godcr/fyne/assets"
	"github.com/raedahgroup/godcr/fyne/helpers"
	"github.com/raedahgroup/godcr/fyne/pages/handler/values"
	"github.com/raedahgroup/godcr/fyne/widgets"
)

func (historyPage *HistoryPageData) fetchTxDetails(hash string) {
	messageLabel := widget.NewLabelWithStyle("Fetching data..", fyne.TextAlignCenter, fyne.TextStyle{})
	time.AfterFunc(time.Millisecond*300, func() {
		if historyPage.TabMenu.CurrentTabIndex() == 1 {
			messageLabel.Hide()
		}
	})

	txDetailslabel := widget.NewLabelWithStyle("Transaction Details", fyne.TextAlignLeading, fyne.TextStyle{Bold: true, Italic: true})

	var txDetailsPopUp *widget.PopUp
	minimizeIcon := widgets.NewImageButton(theme.CancelIcon(), nil, func() { txDetailsPopUp.Hide() })
	errorMessageLabel := widget.NewLabelWithStyle("", fyne.TextAlignCenter, fyne.TextStyle{Bold: true})

	txDetailsErrorMethod := func() {
		txErrorDetailsOutput := widget.NewVBox(
			widget.NewHBox(
				txDetailslabel,
				widgets.NewHSpacer(txDetailslabel.MinSize().Width*3),
				minimizeIcon,
			),
			widget.NewHBox(errorMessageLabel),
		)
		txDetailsPopUp = widget.NewModalPopUp(widget.NewVBox(fyne.NewContainer(txErrorDetailsOutput)), historyPage.Window.Canvas())
		txDetailsPopUp.Show()
	}

	chainHash, err := chainhash.NewHashFromStr(hash)
	if err != nil {
		helpers.ErrorHandler(fmt.Sprintf("Error fetching generating chainhash from for \n %s \n %s ", hash, err.Error()), errorMessageLabel)
		txDetailsErrorMethod()
		return
	}

	txDetails, err := historyPage.MultiWallet.WalletWithID(historyPage.selectedWalletID).GetTransactionRaw(chainHash[:])
	if err != nil {
		helpers.ErrorHandler(fmt.Sprintf("Error fetching transaction details for \n %s \n %s ", hash, err.Error()), errorMessageLabel)
		txDetailsErrorMethod()
		return
	}

	var confirmations int32 = 0
	if txDetails.BlockHeight != -1 {
		confirmations = historyPage.MultiWallet.WalletWithID(historyPage.selectedWalletID).GetBestBlock() - txDetails.BlockHeight + 1
	}

	var status string
	var spendUnconfirmed = historyPage.MultiWallet.ReadBoolConfigValueForKey(dcrlibwallet.SpendUnconfirmedConfigKey, true)
	if spendUnconfirmed || confirmations > dcrlibwallet.DefaultRequiredConfirmations {
		status = "Confirmed"
	} else {
		status = "Pending"
	}

	textObject := func(text string, copyAble bool, align fyne.TextAlign) *widgets.ClickableBox {
		var textToCopy *canvas.Text
		if copyAble {
			if strings.Contains(text, ":") {
				txt := strings.Split(text, ":")
				txt1, txt2 := txt[0], txt[1]
				trimmedText := txt1[:25] + "..." + txt1[len(txt1)-25:] + ":" + txt2
				textToCopy = canvas.NewText(trimmedText, values.Blue)
			} else {
				textToCopy = canvas.NewText(text, values.Blue)
			}
		} else {
			textToCopy = canvas.NewText(text, values.DefaultTextColor)
		}
		textToCopy.TextSize = 14
		textToCopy.Alignment = align

		return widgets.NewClickableBox(widget.NewVBox(widgets.NewVSpacer(1), textToCopy),
			func() {
				messageLabel.SetText("Data Copied")
				clipboard := historyPage.Window.Clipboard()
				clipboard.SetContent(text)
				messageLabel.Show()

				time.AfterFunc(time.Second*2, func() {
					if historyPage.TabMenu.CurrentTabIndex() == 1 {
						messageLabel.Hide()
					}
				})
			},
		)
	}

	txDetailsForm := widget.NewForm()
	txDetailsForm.Append("Fee: ", widget.NewLabelWithStyle(dcrutil.Amount(txDetails.Fee).String(), fyne.TextAlignLeading, fyne.TextStyle{}))
	txDetailsForm.Append("Date: ", widget.NewLabelWithStyle(fmt.Sprintf("%s UTC", dcrlibwallet.FormatUTCTime(txDetails.Timestamp)), fyne.TextAlignLeading, fyne.TextStyle{}))
	txDetailsForm.Append("Type: ", widget.NewLabelWithStyle(txDetails.Type, fyne.TextAlignLeading, fyne.TextStyle{}))
	txDetailsForm.Append("Size: ", widget.NewLabelWithStyle(strconv.Itoa(txDetails.Size)+" Bytes", fyne.TextAlignLeading, fyne.TextStyle{}))
	txDetailsForm.Append("Amount: ", widget.NewLabelWithStyle(dcrutil.Amount(txDetails.Amount).String(), fyne.TextAlignLeading, fyne.TextStyle{}))
	txDetailsForm.Append("Status: ", widget.NewLabelWithStyle(status, fyne.TextAlignLeading, fyne.TextStyle{}))
	txDetailsForm.Append("Fee Rate: ", widget.NewLabelWithStyle(dcrutil.Amount(txDetails.FeeRate).String(), fyne.TextAlignLeading, fyne.TextStyle{}))
	txDetailsForm.Append("Direction: ", widget.NewLabelWithStyle(dcrlibwallet.TransactionDirectionName(txDetails.Direction), fyne.TextAlignLeading, fyne.TextStyle{}))
	txDetailsForm.Append("Block Height: ", widget.NewLabelWithStyle(strconv.Itoa(int(txDetails.BlockHeight)), fyne.TextAlignLeading, fyne.TextStyle{}))
	txDetailsForm.Append("Confirmations: ", widget.NewLabelWithStyle(strconv.Itoa(int(confirmations)), fyne.TextAlignLeading, fyne.TextStyle{}))
	txDetailsForm.Append("Transaction ID: ", textObject(txDetails.Hash, true, fyne.TextAlignLeading))

	var txInput widgets.Table
	inputTableColumnLabels := widget.NewHBox(
		widget.NewLabelWithStyle("Previous Outpoint", fyne.TextAlignLeading, fyne.TextStyle{Bold: true}),
		widget.NewLabelWithStyle("Account", fyne.TextAlignCenter, fyne.TextStyle{Bold: true}),
		widget.NewLabelWithStyle("Amount", fyne.TextAlignCenter, fyne.TextStyle{Bold: true}))

	var inputBox []*widget.Box
	for i := range txDetails.Inputs {
		inputBox = append(inputBox, widget.NewHBox(
			textObject(txDetails.Inputs[i].PreviousOutpoint, true, fyne.TextAlignLeading),
			widget.NewLabelWithStyle(txDetails.Inputs[i].AccountName, fyne.TextAlignCenter, fyne.TextStyle{}),
			widget.NewLabelWithStyle(dcrutil.Amount(txDetails.Inputs[i].Amount).String(), fyne.TextAlignCenter, fyne.TextStyle{}),
		))
	}
	txInput.NewTable(inputTableColumnLabels, inputBox...)

	var txOutput widgets.Table
	outputTableColumnLabels := widget.NewHBox(
		widget.NewLabelWithStyle("Account", fyne.TextAlignCenter, fyne.TextStyle{Bold: true}),
		widget.NewLabelWithStyle("Address", fyne.TextAlignCenter, fyne.TextStyle{Bold: true}),
		widget.NewLabelWithStyle("Value", fyne.TextAlignCenter, fyne.TextStyle{Bold: true}),
		widget.NewLabelWithStyle("Type", fyne.TextAlignCenter, fyne.TextStyle{Bold: true}))

	var outputBox []*widget.Box
	for i := range txDetails.Outputs {
		outputBox = append(outputBox, widget.NewHBox(
			widget.NewLabelWithStyle(txDetails.Outputs[i].AccountName, fyne.TextAlignCenter, fyne.TextStyle{}),
			textObject(txDetails.Outputs[i].Address, true, fyne.TextAlignCenter),
			widget.NewLabelWithStyle(dcrutil.Amount(txDetails.Outputs[i].Amount).String(), fyne.TextAlignTrailing, fyne.TextStyle{}),
			widget.NewLabelWithStyle(txDetails.Outputs[i].ScriptType, fyne.TextAlignCenter, fyne.TextStyle{}),
		))
	}
	txOutput.NewTable(outputTableColumnLabels, outputBox...)

	link, err := url.Parse(fmt.Sprintf("https://%s.dcrdata.org/tx/%s", values.NetType, txDetails.Hash))
	if err != nil {
		helpers.ErrorHandler(fmt.Sprintf("Error: ", err.Error()), errorMessageLabel)
		txDetailsErrorMethod()
		return
	}

	redirectWidget := widget.NewHBox(
		widget.NewHyperlinkWithStyle("View on dcrdata", link, fyne.TextAlignLeading, fyne.TextStyle{Bold: true}),
		widgets.NewHSpacer(5),
		widget.NewIcon(historyPage.icons[assets.RedirectIcon]),
	)

	txDetailsData := widget.NewVBox(
		widgets.NewHSpacer(10),
		txDetailsForm,
		canvas.NewLine(values.TxdetailsLineColor),
		widget.NewHBox(layout.NewSpacer(), redirectWidget, layout.NewSpacer()),
		widgets.NewHSpacer(10),
		canvas.NewLine(values.TxdetailsLineColor),
		widget.NewLabelWithStyle("Inputs", fyne.TextAlignLeading, fyne.TextStyle{Bold: true}),
		txInput.Result,
		widgets.NewHSpacer(10),
		canvas.NewLine(values.TxdetailsLineColor),
		widget.NewLabelWithStyle("Outputs", fyne.TextAlignLeading, fyne.TextStyle{Bold: true}),
		txOutput.Result,
		widgets.NewHSpacer(20),
		widgets.NewVSpacer(10),
	)

	txDetailsScrollContainer := widget.NewScrollContainer(txDetailsData)
	txDetailsOutput := widget.NewVBox(
		widgets.NewHSpacer(10),
		widget.NewHBox(
			txDetailslabel,
			widgets.NewHSpacer(txDetailsData.MinSize().Width-180),
			minimizeIcon,
		),
		widget.NewHBox(widgets.NewHSpacer(txDetailsScrollContainer.MinSize().Width*13), messageLabel),
		fyne.NewContainerWithLayout(layout.NewFixedGridLayout(fyne.NewSize(txDetailsData.MinSize().Width, txDetailsData.MinSize().Height-200)), txDetailsScrollContainer),
		widgets.NewVSpacer(10),
	)

	txDetailsPopUp = widget.NewModalPopUp(fyne.NewContainer(txDetailsOutput), historyPage.Window.Canvas())
}
