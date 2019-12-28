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

	textObject := func(text string, copyAble bool, bold bool) *widgets.ClickableBox {
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
		if bold {
			textToCopy.TextStyle = fyne.TextStyle{Bold: true}
		}
		textToCopy.TextSize = 14
		textToCopy.Alignment = fyne.TextAlignTrailing

		return widgets.NewClickableBox(widget.NewHBox(textToCopy),
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

	tableConfirmations := widget.NewHBox(
		textObject("Confirmations:", false, true),
		textObject(strconv.Itoa(int(confirmations)), false, false),
	)

	tableHash := widget.NewHBox(
		textObject("Transaction ID:", false, true),
		textObject(txDetails.Hash, true, false),
	)

	tableBlockHeight := widget.NewHBox(
		textObject("Block Height:", false, true),
		textObject(strconv.Itoa(int(txDetails.BlockHeight)), false, false),
	)
	tableDirection := widget.NewHBox(
		textObject("Direction:", false, true),
		textObject(dcrlibwallet.TransactionDirectionName(txDetails.Direction), false, false),
	)
	tableType := widget.NewHBox(
		textObject("Type:", false, true),
		textObject(txDetails.Type, false, false),
	)
	tableAmount := widget.NewHBox(
		textObject("Amount:", false, true),
		textObject(dcrutil.Amount(txDetails.Amount).String(), false, false),
	)
	tableSize := widget.NewHBox(
		textObject("Size:", false, true),
		textObject(strconv.Itoa(txDetails.Size)+" Bytes", false, false),
	)
	tableFee := widget.NewHBox(
		textObject("Fee:", false, true),
		textObject(dcrutil.Amount(txDetails.Fee).String(), false, false),
	)
	tableFeeRate := widget.NewHBox(
		textObject("Fee Rate:", false, true),
		textObject(dcrutil.Amount(txDetails.FeeRate).String(), false, false),
	)
	tableStatus := widget.NewHBox(
		textObject("Status:", false, true),
		textObject(status, false, false),
	)
	tableDate := widget.NewHBox(
		textObject("Date:", false, true),
		textObject(fmt.Sprintf("%s UTC", dcrlibwallet.FormatUTCTime(txDetails.Timestamp)), false, false),
	)

	var txInput widgets.Table
	inputTableColumnLabels := widget.NewHBox(
		textObject("Previous Outpoint", false, true),
		textObject("Account", false, true),
		textObject("Amount", false, true))

	var inputBox []*widget.Box
	for i := range txDetails.Inputs {
		inputBox = append(inputBox, widget.NewHBox(
			textObject(txDetails.Inputs[i].PreviousOutpoint, true, false),
			textObject(txDetails.Inputs[i].AccountName, false, false),
			textObject(dcrutil.Amount(txDetails.Inputs[i].Amount).String(), false, false),
		))
	}
	txInput.NewTable(inputTableColumnLabels, inputBox...)

	var txOutput widgets.Table
	outputTableColumnLabels := widget.NewHBox(
		textObject("Address", false, true),
		textObject("Account", false, true),
		textObject("Value", false, true),
		textObject("Type", false, true))

	var outputBox []*widget.Box
	for i := range txDetails.Outputs {
		outputBox = append(outputBox, widget.NewHBox(
			textObject(txDetails.Outputs[i].AccountName, false, false),
			textObject(txDetails.Outputs[i].Address, true, false),
			textObject(dcrutil.Amount(txDetails.Outputs[i].Amount).String(), false, false),
			textObject(txDetails.Outputs[i].ScriptType, false, false),
		))
	}
	txOutput.NewTable(outputTableColumnLabels, outputBox...)

	tableData := widget.NewVBox(
		tableConfirmations,
		widgets.NewVSpacer(values.SpacerSize4),
		tableHash,
		widgets.NewVSpacer(values.SpacerSize4),
		tableBlockHeight,
		widgets.NewVSpacer(values.SpacerSize4),
		tableDirection,
		widgets.NewVSpacer(values.SpacerSize4),
		tableType,
		widgets.NewVSpacer(values.SpacerSize4),
		tableAmount,
		widgets.NewVSpacer(values.SpacerSize4),
		tableSize,
		widgets.NewVSpacer(values.SpacerSize4),
		tableFee,
		widgets.NewVSpacer(values.SpacerSize4),
		tableFeeRate,
		widgets.NewVSpacer(values.SpacerSize4),
		tableStatus,
		widgets.NewVSpacer(values.SpacerSize4),
		tableDate,
		widgets.NewVSpacer(values.SpacerSize4),
	)

	link, err := url.Parse(fmt.Sprintf("https://%s.dcrdata.org/tx/%s", values.NetType, txDetails.Hash))
	if err != nil {
		helpers.ErrorHandler(fmt.Sprintf("Error: ", err.Error()), errorMessageLabel)
		txDetailsErrorMethod()
		return
	}

	redirectWidget := widget.NewHBox(
		widget.NewHyperlinkWithStyle("View on dcrdata", link, fyne.TextAlignLeading, fyne.TextStyle{Bold: true}),
		widgets.NewHSpacer(values.SpacerSize20),
		widget.NewIcon(historyPage.icons[assets.RedirectIcon]),
	)

	txDetailsData := widget.NewVBox(
		widgets.NewHSpacer(values.SpacerSize10),
		tableData,
		canvas.NewLine(values.TxdetailsLineColor),
		redirectWidget,
		widgets.NewHSpacer(values.SpacerSize10),
		canvas.NewLine(values.TxdetailsLineColor),
		widget.NewLabelWithStyle("Inputs", fyne.TextAlignLeading, fyne.TextStyle{Bold: true}),
		txInput.Result,
		widgets.NewHSpacer(values.SpacerSize10),
		canvas.NewLine(values.TxdetailsLineColor),
		widget.NewLabelWithStyle("Outputs", fyne.TextAlignLeading, fyne.TextStyle{Bold: true}),
		txOutput.Result,
		widgets.NewHSpacer(values.SpacerSize20),
		widgets.NewVSpacer(values.SpacerSize10),
	)

	txDetailsScrollContainer := widget.NewScrollContainer(txDetailsData)
	txDetailsOutput := widget.NewVBox(
		widgets.NewHSpacer(values.SpacerSize10),
		widget.NewHBox(
			txDetailslabel,
			widgets.NewHSpacer(txDetailsData.MinSize().Width-180),
			minimizeIcon,
		),
		widget.NewHBox(widgets.NewHSpacer(txDetailsScrollContainer.MinSize().Width*12), messageLabel),
		fyne.NewContainerWithLayout(layout.NewFixedGridLayout(fyne.NewSize(txDetailsData.MinSize().Width, txDetailsData.MinSize().Height-200)), txDetailsScrollContainer),
		widgets.NewVSpacer(values.SpacerSize10),
	)

	txDetailsPopUp = widget.NewModalPopUp(fyne.NewContainer(txDetailsOutput), historyPage.Window.Canvas())
}
