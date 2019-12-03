package pages

import (
	"fmt"
	"image/color"
	"log"
	"regexp"
	"sort"
	"strconv"
	"strings"
	"time"

	"fyne.io/fyne"
	"fyne.io/fyne/canvas"
	"fyne.io/fyne/layout"
	"fyne.io/fyne/widget"

	"github.com/decred/dcrd/dcrutil"
	"github.com/raedahgroup/dcrlibwallet"

	"github.com/raedahgroup/godcr/fyne/assets"
	"github.com/raedahgroup/godcr/fyne/layouts"
	"github.com/raedahgroup/godcr/fyne/pages/sendpagehandler"
	"github.com/raedahgroup/godcr/fyne/widgets"
)

type sendPageDynamicData struct {
	// houses all clickable box
	sendingAccountBoxes     []*widget.Box
	selfSendingAccountBoxes []*widget.Box

	errorLabel     *widgets.Button
	spendableLabel *canvas.Text

	selfSendingSelectedAccountLabel        *widget.Label
	selfSendingSelectedAccountBalanceLabel *widget.Label
	selfSendingSelectedWalletLabel         *canvas.Text
	selfSendingSelectedWalletID            int

	sendingSelectedAccountLabel        *widget.Label
	sendingSelectedAccountBalanceLabel *widget.Label
	sendingSelectedWalletID            int

	Contents *widget.Box
}

var sendPage sendPageDynamicData

func sendPageContent(multiWallet *dcrlibwallet.MultiWallet, window fyne.Window) (pageContent fyne.CanvasObject) {
	icons, err := assets.GetIcons(assets.Reveal, assets.Conceal, assets.InfoIcon, assets.MoreIcon,
		assets.ReceiveAccountIcon, assets.CollapseIcon, assets.CollapseDropdown, assets.ExpandDropdown, assets.DownArrow, assets.Alert)
	if err != nil {
		return widget.NewLabelWithStyle(err.Error(), fyne.TextAlignCenter, fyne.TextStyle{Bold: true})
	}

	sendPage.errorLabel = widgets.NewButton(color.RGBA{237, 109, 71, 255}, "", nil)
	sendPage.errorLabel.Container.Hide()

	sendPage.Contents = widget.NewVBox()

	showErrorLabel := func(value string) {
		sendPage.errorLabel.SetText(value)
		sendPage.errorLabel.SetMinSize(sendPage.errorLabel.MinSize().Add(fyne.NewSize(20, 8)))
		sendPage.errorLabel.Container.Show()
		sendPage.Contents.Refresh()

		time.AfterFunc(time.Second*5, func() {
			sendPage.errorLabel.Container.Hide()
			sendPage.Contents.Refresh()
		})
	}

	refresher := func(objects ...fyne.Widget) {
		for _, object := range objects {
			object.Refresh()
		}
	}

	successLabelContainer := widgets.NewButton(color.RGBA{65, 190, 83, 255}, "Transaction sent", nil)
	successLabelContainer.SetMinSize(successLabelContainer.MinSize().Add(fyne.NewSize(20, 16)))
	successLabelContainer.Container.Hide()

	openedWalletIDs := multiWallet.OpenedWalletIDsRaw()
	if len(openedWalletIDs) == 0 {
		return widget.NewHBox(widgets.NewHSpacer(10), widget.NewLabelWithStyle("Could not retrieve wallets", fyne.TextAlignCenter, fyne.TextStyle{Bold: true}))
	}
	sort.Ints(openedWalletIDs)

	sendPage.sendingSelectedWalletID = openedWalletIDs[0]
	sendPage.selfSendingSelectedWalletID = openedWalletIDs[0]

	sendPage.selfSendingAccountBoxes = make([]*widget.Box, len(openedWalletIDs))
	sendPage.sendingAccountBoxes = make([]*widget.Box, len(openedWalletIDs))

	var selectedWallet = multiWallet.WalletWithID(sendPage.sendingSelectedWalletID)

	selectedWalletAccounts, err := selectedWallet.GetAccountsRaw(dcrlibwallet.DefaultRequiredConfirmations)
	if err != nil {
		log.Println(fmt.Sprintf("Error while getting accounts for wallet %s", err.Error()))
		return widget.NewLabel("Error while getting accounts for wallet")
	}

	sendingSelectedWalletLabel := canvas.NewText(selectedWallet.Name, color.RGBA{137, 151, 165, 255})
	selfSendingSelectedWalletLabel := canvas.NewText(selectedWallet.Name, color.RGBA{137, 151, 165, 255})

	sendPage.spendableLabel = canvas.NewText("Spendable: "+dcrutil.Amount(selectedWalletAccounts.Acc[0].Balance.Spendable).String(), color.Black)
	sendPage.spendableLabel.TextSize = 12

	temporaryAddress, err := selectedWallet.CurrentAddress(selectedWalletAccounts.Acc[0].Number)
	if err != nil {
		log.Println("could not retrieve account details", err.Error())
		return widget.NewLabel("could not retrieve account details")
	}

	amountInAccount := dcrlibwallet.AmountCoin(selectedWalletAccounts.Acc[0].TotalBalance)

	var transactionAuthor *dcrlibwallet.TxAuthor
	initiateTxAuthor := func(accountNumber int32) {
		transactionAuthor = selectedWallet.NewUnsignedTx(accountNumber, dcrlibwallet.DefaultRequiredConfirmations)
		transactionAuthor.AddSendDestination(temporaryAddress, 0, true)
	}
	initiateTxAuthor(0)

	var amountEntry *widget.Entry
	var amountErrorLabel *canvas.Text

	costAndBalanceAfterSendBox := widget.NewVBox()

	totalCostLabel := widget.NewLabel("- DCR")
	costAndBalanceAfterSendBox.Append(widget.NewHBox(widget.NewLabel("Total cost"), layout.NewSpacer(), totalCostLabel))

	balanceAfterSendLabel := widget.NewLabel("- DCR")
	costAndBalanceAfterSendBox.Append(widget.NewHBox(widget.NewLabel("Balance after send"), layout.NewSpacer(), balanceAfterSendLabel))

	costAndBalanceAfterSendContainer := fyne.NewContainerWithLayout(layout.NewFixedGridLayout(fyne.NewSize(328, 48)), costAndBalanceAfterSendBox)

	transactionFeeLabel := widget.NewLabel("- DCR")
	transactionSize := widget.NewLabel("0 bytes")

	var transactionFeeBox *widget.Box

	// this function is called when the sending wallet account is changed.
	onSendingAccountChange := func() {
		selectedWallet = multiWallet.WalletWithID(sendPage.sendingSelectedWalletID)

		accountNumber, err := selectedWallet.AccountNumber(sendPage.sendingSelectedAccountLabel.Text)
		if err != nil {
			showErrorLabel("Could not get accounts")
			log.Println("could not get accounts on account change, reason:", err.Error())
			return
		}
		initiateTxAuthor(int32(accountNumber))

		balance, err := selectedWallet.GetAccountBalance(int32(accountNumber), dcrlibwallet.DefaultRequiredConfirmations)
		if err != nil {
			showErrorLabel("could not retrieve account balance")
			log.Println("could not retrieve account balance on account change, reason:", err.Error())
			return
		}

		sendPage.spendableLabel.Text = "Spendable: " + dcrutil.Amount(balance.Spendable).String()
		sendPage.spendableLabel.Refresh()

		amountInAccount = dcrlibwallet.AmountCoin(balance.Total)

		// reset amount entry
		costAndBalanceAfterSendContainer.Refresh()
		transactionFeeLabel.Refresh()
		transactionSize.Refresh()
		transactionFeeBox.Refresh()
		sendPage.Contents.Refresh()
		amountEntry.OnChanged(amountEntry.Text)
	}

	sendPage.sendingSelectedAccountLabel = widget.NewLabel(selectedWalletAccounts.Acc[0].Name)
	sendPage.sendingSelectedAccountBalanceLabel = widget.NewLabel(dcrutil.Amount(selectedWalletAccounts.Acc[0].TotalBalance).String())

	fromAccountSelector := sendpagehandler.FromAccountSelector(onSendingAccountChange, "Sending account",
		icons[assets.ReceiveAccountIcon], icons[assets.CollapseIcon], multiWallet, openedWalletIDs, &sendPage.sendingSelectedWalletID,
		sendPage.sendingAccountBoxes, sendPage.sendingSelectedAccountLabel, sendPage.sendingSelectedAccountBalanceLabel,
		sendingSelectedWalletLabel, sendPage.Contents)

	sendPage.selfSendingSelectedAccountLabel = widget.NewLabel(selectedWalletAccounts.Acc[0].Name)
	sendPage.selfSendingSelectedAccountBalanceLabel = widget.NewLabel(dcrutil.Amount(selectedWalletAccounts.Acc[0].TotalBalance).String())

	selfSendingToAccountClickableBox := sendpagehandler.CreateAccountSelector(sendPage.Contents.Refresh, "Receiving account",
		icons[assets.ReceiveAccountIcon], icons[assets.CollapseIcon], multiWallet, openedWalletIDs, &sendPage.selfSendingSelectedWalletID,
		sendPage.selfSendingAccountBoxes, sendPage.selfSendingSelectedAccountLabel, sendPage.selfSendingSelectedAccountBalanceLabel,
		selfSendingSelectedWalletLabel, sendPage.Contents)

	selfSendingToAccountGroup := widget.NewGroup("To", selfSendingToAccountClickableBox)
	selfSendingToAccountGroup.Hide()

	destinationAddressEntry := widget.NewEntry()
	destinationAddressEntry.SetPlaceHolder("Destination Address")

	// shows errors related too destination address
	destinationAddressErrorLabel := canvas.NewText("", color.RGBA{237, 109, 71, 255})
	destinationAddressErrorLabel.TextSize = 12
	destinationAddressErrorLabel.Hide()

	var nextButton *widgets.Button

	destinationAddressEntry.OnChanged = func(address string) {
		if destinationAddressEntry.Text == "" {
			destinationAddressErrorLabel.Hide()
			sendPage.Contents.Refresh()
			return
		}

		_, err := dcrutil.DecodeAddress(address)
		if err != nil {
			destinationAddressErrorLabel.Text = "Invalid address"
			destinationAddressErrorLabel.Show()
		} else {
			destinationAddressErrorLabel.Hide()
		}

		if amountEntry.Text != "" && amountErrorLabel.Hidden && destinationAddressErrorLabel.Hidden {
			nextButton.Enable()
		} else {
			nextButton.Disable()
		}

		sendPage.Contents.Refresh()
	}

	destinationAddressEntryGroup := widget.NewGroup("To", fyne.NewContainerWithLayout(
		layout.NewFixedGridLayout(fyne.NewSize(widget.NewLabel(temporaryAddress).MinSize().Width, destinationAddressEntry.MinSize().Height)),
		destinationAddressEntry),
		destinationAddressErrorLabel)

	sendToAccountLabel := canvas.NewText("Send to account", color.RGBA{R: 41, G: 112, B: 255, A: 255})
	sendToAccountLabel.TextSize = 12

	destinationBox := widget.NewHBox(destinationAddressEntryGroup, selfSendingToAccountGroup, layout.NewSpacer())

	// This hides self sending account dropdown or destination address entry.
	sendToAccount := widgets.NewClickableBox(widget.NewVBox(sendToAccountLabel), func() {
		if selfSendingToAccountGroup.Hidden {
			sendToAccountLabel.Text = "Send to address"
			selfSendingToAccountGroup.Show()
			destinationAddressEntryGroup.Hide()

			if amountEntry.Text != "" {
				nextButton.Enable()
			} else {
				nextButton.Disable()
			}
		} else {
			sendToAccountLabel.Text = "Send to account"
			destinationAddressEntryGroup.Show()
			selfSendingToAccountGroup.Hide()

			if amountEntry.Text != "" && destinationAddressEntry.Text != "" && destinationAddressErrorLabel.Hidden {
				nextButton.Enable()
			} else {
				nextButton.Disable()
			}
		}

		sendPage.Contents.Refresh()
		amountEntry.OnChanged(amountEntry.Text)
	})

	destinationBox.Append(widget.NewVBox(sendToAccount)) // placed it in a VBox so as to center object

	amountEntry = widget.NewEntry()
	amountEntry.SetPlaceHolder("0 DCR")

	transactionInfoform := fyne.NewContainerWithLayout(layout.NewVBoxLayout())

	transactionInfoform.AddObject(fyne.NewContainerWithLayout(
		layout.NewHBoxLayout(), widget.NewLabel("Processing time"), widgets.NewHSpacer(46),
		layout.NewSpacer(), widget.NewLabelWithStyle("Approx. 10 mins (2 blocks)", fyne.TextAlignLeading, fyne.TextStyle{})))

	transactionInfoform.AddObject(fyne.NewContainerWithLayout(
		layout.NewHBoxLayout(), widget.NewLabel("Fee rate"), layout.NewSpacer(),
		widget.NewLabelWithStyle("0.0001 DCR/byte", fyne.TextAlignLeading, fyne.TextStyle{})))

	transactionInfoform.AddObject(fyne.NewContainerWithLayout(layout.NewHBoxLayout(),
		widget.NewLabel("Transaction size"), layout.NewSpacer(), transactionSize))

	paintedtransactionInfoform := fyne.NewContainerWithLayout(layout.NewBorderLayout(nil, nil, nil, nil),
		canvas.NewRectangle(color.RGBA{158, 158, 158, 0xff}), transactionInfoform)

	paintedtransactionInfoform.Hide()

	var transactionSizeDropdown *widgets.ClickableBox
	var amountEntryGroup *widget.Group

	transactionSizeDropdown = widgets.NewClickableBox(widget.NewHBox(widget.NewIcon(icons[assets.ExpandDropdown])), func() {
		if paintedtransactionInfoform.Hidden {
			transactionSizeDropdown.Box.Children[0] = widget.NewIcon(icons[assets.CollapseDropdown])
			paintedtransactionInfoform.Show()
		} else {
			transactionSizeDropdown.Box.Children[0] = widget.NewIcon(icons[assets.ExpandDropdown])
			paintedtransactionInfoform.Hide()
		}

		transactionFeeBox.Refresh()
		transactionSizeDropdown.Refresh()
		paintedtransactionInfoform.Refresh()
		sendPage.Contents.Refresh()
	})

	amountErrorLabel = canvas.NewText("", color.RGBA{237, 109, 71, 255})
	amountErrorLabel.TextSize = 12
	amountErrorLabel.Hide()

	maxButton := widgets.NewButton(color.RGBA{61, 88, 115, 255}, "MAX", func() {
		transactionAuthor.UpdateSendDestination(0, temporaryAddress, 0, true)

		maxAmount, err := transactionAuthor.EstimateMaxSendAmount()
		if err != nil {
			if err.Error() == dcrlibwallet.ErrInsufficientBalance {
				amountErrorLabel.Text = "Not enough funds"
				if !multiWallet.IsSynced() {
					amountErrorLabel.Text = "Not enough funds (or not connected)."
				}

				amountErrorLabel.Show()
				sendPage.Contents.Refresh()
				return
			}
		}

		amountErrorLabel.Hide()
		sendPage.Contents.Refresh()
		amountEntry.SetText(fmt.Sprintf("%f", maxAmount.DcrValue-0.000012))
	})

	maxButton.SetTextSize(9)
	maxButton.SetMinSize(maxButton.MinSize().Add(fyne.NewSize(8, 8)))

	transactionFeeBox = widget.NewHBox(widget.NewLabel("Transaction fee"), widgets.NewHSpacer(149), transactionFeeLabel, transactionSizeDropdown)

	amountEntryGroup = widget.NewGroup("Amount", fyne.NewContainerWithLayout(
		layouts.NewPasswordLayout(
			fyne.NewSize(widget.NewLabel("12345678.12345678").MinSize().Width+maxButton.MinSize().Width, amountEntry.MinSize().Height)),
		amountEntry, maxButton.Container),
		amountErrorLabel, widgets.NewVSpacer(4),
		transactionFeeBox,
		paintedtransactionInfoform)

	// amount entry accepts only floats
	amountEntryExpression, err := regexp.Compile("^\\d*\\.?\\d*$")
	if err != nil {
		log.Println(err)
	}

	amountEntry.OnChanged = func(value string) {
		if len(value) > 0 && !amountEntryExpression.MatchString(value) {
			if len(value) == 1 {
				amountEntry.SetText("")
			} else {
				//fix issue with crash on paste here
				value = value[:amountEntry.CursorColumn-1] + value[amountEntry.CursorColumn:]
				//todo: using setText, cursor column count doesnt increase or reduce. Create an issue on this
				amountEntry.CursorColumn--
				amountEntry.SetText(value)
			}
			return
		}

		if numbers := strings.Split(value, "."); len(numbers) == 2 {
			if len(numbers[1]) > 8 {
				showErrorLabel("Amount has more than 8 decimal places.")
				return
			}
		}

		amountInFloat, err := strconv.ParseFloat(value, 64)
		if err != nil && value != "" {
			showErrorLabel("Could not parse float")
			return
		}

		if amountInFloat == 0.0 {
			transactionFeeLabel.SetText("- DCR")
			totalCostLabel.SetText("- DCR")
			balanceAfterSendLabel.SetText("- DCR")
			transactionSize.SetText("0 bytes")
			amountErrorLabel.Hide()

			nextButton.Disable()
			refresher(costAndBalanceAfterSendBox, transactionFeeLabel, totalCostLabel, balanceAfterSendLabel, transactionSize)
			paintedtransactionInfoform.Refresh()
			sendPage.Contents.Refresh()
			return
		}

		transactionAuthor.UpdateSendDestination(0, temporaryAddress, dcrlibwallet.AmountAtom(amountInFloat), false)
		feeAndSize, err := transactionAuthor.EstimateFeeAndSize()
		if err != nil {
			if err.Error() == dcrlibwallet.ErrInsufficientBalance {
				amountErrorLabel.Text = "Insufficient balance"
				amountErrorLabel.Show()
			} else {
				showErrorLabel("Could not retrieve transaction fee and size")
				log.Println(fmt.Sprintf("could not retrieve transaction fee and size %s", err.Error()))
			}

			transactionFeeLabel.SetText(fmt.Sprintf("- DCR"))
			totalCostLabel.SetText("- DCR")
			balanceAfterSendLabel.SetText("- DCR")
			transactionSize.SetText("0 bytes")

			nextButton.Disable()
			paintedtransactionInfoform.Refresh()
			refresher(transactionFeeLabel, totalCostLabel, balanceAfterSendLabel, transactionSize, costAndBalanceAfterSendBox)
			sendPage.Contents.Refresh()
			return
		}

		if !amountErrorLabel.Hidden {
			amountErrorLabel.Hide()
		}

		transactionFeeLabel.SetText(fmt.Sprintf("%f DCR", feeAndSize.Fee.DcrValue))
		totalCostLabel.SetText(fmt.Sprintf("%f DCR", feeAndSize.Fee.DcrValue+amountInFloat))
		balanceAfterSendLabel.SetText(fmt.Sprintf("%f DCR", amountInAccount-(feeAndSize.Fee.DcrValue+amountInFloat)))
		transactionSize.SetText(fmt.Sprintf("%d bytes", feeAndSize.EstimatedSignedSize))

		if destinationAddressEntry.Text != "" && destinationAddressErrorLabel.Hidden || destinationAddressEntryGroup.Hidden {
			nextButton.Enable()
		} else {
			nextButton.Disable()
		}

		sendPage.errorLabel.Container.Hide()
		paintedtransactionInfoform.Refresh()
		refresher(transactionFeeLabel, totalCostLabel, balanceAfterSendLabel, transactionSize, costAndBalanceAfterSendBox)
		sendPage.Contents.Refresh()
	}

	nextButton = widgets.NewButton(color.RGBA{41, 112, 255, 255}, "Next", func() {
		if multiWallet.ConnectedPeers() <= 0 {
			showErrorLabel("Not Connected To Decred Network")

			return
		}

		var sendingAddress, selfSendingSelectedWalletName string
		var amountInFloat float64
		var err error

		amountInFloat, err = strconv.ParseFloat(amountEntry.Text, 64)
		if err != nil {
			showErrorLabel("Could not parse float")
			return
		}

		// if sending to self
		if destinationAddressEntryGroup.Hidden {
			sendingToSelfSelectedWallet := multiWallet.WalletWithID(sendPage.selfSendingSelectedWalletID)
			if sendingToSelfSelectedWallet == nil {
				showErrorLabel("Selected wallet is invalid")
				return
			}

			var accountNo uint32
			selfSendingSelectedWalletName = sendingToSelfSelectedWallet.Name
			accountNo, err = sendingToSelfSelectedWallet.AccountNumber(sendPage.selfSendingSelectedAccountLabel.Text)
			if err != nil {
				showErrorLabel("Could not get self sending account")
				log.Println("could not get self sending account, reason: ", err.Error())

				return
			}

			sendingAddress, err = sendingToSelfSelectedWallet.CurrentAddress(int32(accountNo))
			if err != nil {
				showErrorLabel("could not get self sending account")
				log.Println("could not get self sending account reason:", err.Error())

				return
			}

		} else {
			sendingAddress = destinationAddressEntry.Text
		}

		transactionAuthor.UpdateSendDestination(0, sendingAddress, dcrlibwallet.AmountAtom(amountInFloat), false)

		sendpagehandler.ConfirmationWindow(amountEntry, destinationAddressEntry, icons[assets.DownArrow], icons[assets.Alert], icons[assets.Reveal], icons[assets.Conceal],
			window, selectedWallet.Name, selfSendingSelectedWalletName, totalCostLabel.Text, transactionFeeLabel.Text, balanceAfterSendLabel.Text, sendPage.sendingSelectedAccountLabel.Text,
			sendPage.selfSendingSelectedAccountLabel.Text, destinationAddressEntryGroup.Hidden, transactionAuthor, successLabelContainer, sendPage.Contents)
	})

	nextButton.SetMinSize(nextButton.MinSize().Add(fyne.NewSize(0, 20)))
	nextButton.Disable()

	// define base widget consisting of label, more icon and info button
	sendLabel := widget.NewLabelWithStyle("Send DCR", fyne.TextAlignTrailing, fyne.TextStyle{Bold: true, Italic: true})

	dialogLabel := widget.NewLabelWithStyle("Input the destination \nwallet address and the amount in \nDCR to send funds.", fyne.TextAlignLeading, fyne.TextStyle{})

	var clickabelInfoIcon *widgets.ImageButton
	clickabelInfoIcon = widgets.NewImageButton(icons[assets.InfoIcon], nil, func() {
		var popup *widget.PopUp
		confirmationText := canvas.NewText("Got it", color.RGBA{41, 112, 255, 255})
		confirmationText.TextStyle.Bold = true

		dialog := widget.NewVBox(
			widgets.NewVSpacer(12),
			widget.NewLabelWithStyle("Send DCR", fyne.TextAlignLeading, fyne.TextStyle{Bold: true}),
			widgets.NewVSpacer(30),
			dialogLabel,
			widget.NewHBox(layout.NewSpacer(), widgets.NewClickableBox(widget.NewVBox(confirmationText), func() { popup.Hide() })),
			widgets.NewVSpacer(10))

		popup = widget.NewModalPopUp(widget.NewHBox(widgets.NewHSpacer(24), dialog, widgets.NewHSpacer(20)), window.Canvas())
	})

	var clickabelMoreIcon *widgets.ImageButton
	clickabelMoreIcon = widgets.NewImageButton(icons[assets.MoreIcon], nil, func() {
		var popup *widget.PopUp
		popup = widget.NewPopUp(widgets.NewButton(color.White, "Clear all fields", func() {
			amountEntry.SetText("")
			destinationAddressEntry.SetText("")
			popup.Hide()

		}).Container, window.Canvas())
		popup.Move(fyne.CurrentApp().Driver().AbsolutePositionForObject(
			clickabelMoreIcon).Add(fyne.NewPos(clickabelMoreIcon.MinSize().Width, clickabelMoreIcon.MinSize().Height)))
	})

	baseWidgets := widget.NewHBox(sendLabel, layout.NewSpacer(), clickabelInfoIcon, clickabelMoreIcon)

	sendPage.Contents.Append(widgets.NewVSpacer(10))
	sendPage.Contents.Append(baseWidgets)
	sendPage.Contents.Append(widgets.NewVSpacer(10))
	sendPage.Contents.Append(widget.NewHBox(layout.NewSpacer(), successLabelContainer.Container, sendPage.errorLabel.Container, layout.NewSpacer()))
	sendPage.Contents.Append(fromAccountSelector)
	sendPage.Contents.Append(widgets.NewVSpacer(8))
	sendPage.Contents.Append(destinationBox)
	sendPage.Contents.Append(widgets.NewVSpacer(8))
	sendPage.Contents.Append(widget.NewHBox(amountEntryGroup, widget.NewVBox(sendPage.spendableLabel)))
	sendPage.Contents.Append(widgets.NewVSpacer(12))
	sendPage.Contents.Append(costAndBalanceAfterSendContainer)
	sendPage.Contents.Append(widgets.NewVSpacer(15))
	sendPage.Contents.Append(nextButton.Container)

	return widget.NewHBox(widgets.NewHSpacer(15), sendPage.Contents)
}
