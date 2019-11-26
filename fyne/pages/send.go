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
	"fyne.io/fyne/theme"
	"fyne.io/fyne/widget"

	"github.com/decred/dcrd/dcrutil"
	"github.com/raedahgroup/dcrlibwallet"
	"github.com/raedahgroup/godcr/fyne/assets"
	"github.com/raedahgroup/godcr/fyne/layouts"
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
}

var sendPage sendPageDynamicData

func sendPageContent(multiWallet *dcrlibwallet.MultiWallet, window fyne.Window) (pageContent fyne.CanvasObject) {
	icons, err := assets.GetIcons(assets.Reveal, assets.Conceal, assets.InfoIcon, assets.MoreIcon, assets.ReceiveAccountIcon, assets.CollapseIcon, assets.CollapseDropdown, assets.ExpandDropdown, assets.DownArrow, assets.Alert)
	if err != nil {
		return widget.NewLabelWithStyle(err.Error(), fyne.TextAlignCenter, fyne.TextStyle{Bold: true})
	}

	var sendPageContents *widget.Box

	sendPage.errorLabel = widgets.NewButton(color.RGBA{237, 109, 71, 255}, "", nil)
	sendPage.errorLabel.Container.Hide()

	showErrorLabel := func(value string) {
		sendPage.errorLabel.SetText(value)
		sendPage.errorLabel.SetMinSize(sendPage.errorLabel.MinSize().Add(fyne.NewSize(20, 8)))
		sendPage.errorLabel.Container.Show()
		widget.Refresh(sendPageContents)

		time.AfterFunc(time.Second*5, func() {
			sendPage.errorLabel.Container.Hide()
			widget.Refresh(sendPageContents)
		})
	}

	successLabelContainer := widgets.NewButton(color.RGBA{65, 190, 83, 255}, "Transaction sent", nil)
	successLabelContainer.Container.Hide()

	// define base widget consisting of label, more icon and info button
	sendLabel := widget.NewLabelWithStyle("Send DCR", fyne.TextAlignTrailing, fyne.TextStyle{Bold: true, Italic: true})
	clickabelInfoIcon := widgets.NewImageButton(icons[assets.InfoIcon], nil, func() {

	})
	clickabelMoreIcon := widgets.NewImageButton(icons[assets.MoreIcon], nil, func() {

	})

	baseWidgets := widget.NewHBox(sendLabel, layout.NewSpacer(), clickabelInfoIcon, clickabelMoreIcon)

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

	transactionAuthor := selectedWallet.NewUnsignedTx(0, dcrlibwallet.DefaultRequiredConfirmations)
	transactionAuthor.AddSendDestination(temporaryAddress, 0, true)

	var sendToAccount *widgets.ClickableBox
	var amountEntry *widget.Entry
	var amountErrorLabel *canvas.Text

	var transactionFeeLabel *widget.Label
	var totalCostLabel *widget.Label
	var balanceAfterSendLabel *widget.Label
	var transactionSize *widget.Label

	// this function is called when the sending wallet account is changed.
	onSendingAccountChange := func() {
		selectedWallet = multiWallet.WalletWithID(sendPage.sendingSelectedWalletID)
		accountNumber, err := selectedWallet.AccountNumber(sendPage.sendingSelectedAccountLabel.Text)
		if err != nil {
			showErrorLabel("Could not get account, " + err.Error())
			return
		}

		transactionAuthor = selectedWallet.NewUnsignedTx(int32(accountNumber), dcrlibwallet.DefaultRequiredConfirmations)
		fmt.Println("Changed account", selectedWallet.Name, sendPage.sendingSelectedAccountLabel.Text)
		transactionAuthor.AddSendDestination(temporaryAddress, 0, true)

		sendPage.spendableLabel.Text = "Spendable: " + sendPage.sendingSelectedAccountBalanceLabel.Text

		balance, err := selectedWallet.GetAccountBalance(int32(accountNumber), dcrlibwallet.DefaultRequiredConfirmations)
		if err != nil {
			showErrorLabel("could not retrieve account balance")
			log.Println("could not retrieve account balance", err.Error())
			return
		}

		amountInAccount = dcrlibwallet.AmountCoin(balance.Total)
		// reset amount entry
		amountEntry.OnChanged(amountEntry.Text)

		widget.Refresh(sendPageContents)
	}

	// we still need a suitable name for this
	sendPage.sendingSelectedAccountLabel = widget.NewLabel(selectedWalletAccounts.Acc[0].Name)
	sendPage.sendingSelectedAccountBalanceLabel = widget.NewLabel(dcrutil.Amount(selectedWalletAccounts.Acc[0].TotalBalance).String())

	sendingAccountClickableBox := createAccountDropdown(onSendingAccountChange, "Sending account", icons[assets.ReceiveAccountIcon], icons[assets.CollapseIcon], multiWallet, openedWalletIDs, &sendPage.sendingSelectedWalletID, sendPage.sendingAccountBoxes, sendPage.sendingSelectedAccountLabel, sendPage.sendingSelectedAccountBalanceLabel, sendingSelectedWalletLabel)
	sendingAccountGroup := widget.NewGroup("From", sendingAccountClickableBox)

	sendPage.selfSendingSelectedAccountLabel = widget.NewLabel(selectedWalletAccounts.Acc[0].Name)
	sendPage.selfSendingSelectedAccountBalanceLabel = widget.NewLabel(dcrutil.Amount(selectedWalletAccounts.Acc[0].TotalBalance).String())

	selfSendingToAccountClickableBox := createAccountDropdown(nil, "Receiving account", icons[assets.ReceiveAccountIcon], icons[assets.CollapseIcon], multiWallet, openedWalletIDs, &sendPage.selfSendingSelectedWalletID, sendPage.selfSendingAccountBoxes, sendPage.selfSendingSelectedAccountLabel, sendPage.selfSendingSelectedAccountBalanceLabel, selfSendingSelectedWalletLabel)
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
			canvas.Refresh(destinationAddressErrorLabel)
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
		widget.Refresh(sendPageContents)
	}

	destinationAddressEntryGroup := widget.NewGroup("To", fyne.NewContainerWithLayout(
		layout.NewFixedGridLayout(fyne.NewSize(widget.NewLabel(temporaryAddress).MinSize().Width, destinationAddressEntry.MinSize().Height)), destinationAddressEntry),
		destinationAddressErrorLabel)

	sendToAccountLabel := canvas.NewText("Send to account", color.RGBA{R: 41, G: 112, B: 255, A: 255})
	sendToAccountLabel.TextSize = 12

	destinationBox := widget.NewHBox(destinationAddressEntryGroup, selfSendingToAccountGroup)

	// This hides self sending account dropdown or destination address entry.
	sendToAccount = widgets.NewClickableBox(widget.NewVBox(sendToAccountLabel), func() {
		if sendToAccountLabel.Text == "Send to account" {
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

		widget.Refresh(sendPageContents)
		amountEntry.OnChanged(amountEntry.Text)
	})
	destinationBox.Append(widget.NewVBox(sendToAccount)) // placed it in a VBox so as to center object

	amountEntry = widget.NewEntry()
	amountEntry.SetPlaceHolder("0 DCR")

	transactionFeeLabel = widget.NewLabel("- DCR")
	transactionSize = widget.NewLabelWithStyle("0 bytes", fyne.TextAlignLeading, fyne.TextStyle{})
	transactionInfoform := fyne.NewContainerWithLayout(layout.NewVBoxLayout())
	transactionInfoform.AddObject(fyne.NewContainerWithLayout(layout.NewHBoxLayout(), widget.NewLabel("Processing time"), widgets.NewHSpacer(46), layout.NewSpacer(), widget.NewLabelWithStyle("Approx. 10 mins (2 blocks)", fyne.TextAlignLeading, fyne.TextStyle{})))
	transactionInfoform.AddObject(fyne.NewContainerWithLayout(layout.NewHBoxLayout(), widget.NewLabel("Fee rate"), layout.NewSpacer(), widget.NewLabelWithStyle("0.0001 DCR/byte", fyne.TextAlignLeading, fyne.TextStyle{})))
	transactionInfoform.AddObject(fyne.NewContainerWithLayout(layout.NewHBoxLayout(), widget.NewLabel("Transaction size"), layout.NewSpacer(), transactionSize))

	paintedtransactionInfoform := fyne.NewContainerWithLayout(layout.NewBorderLayout(nil, nil, nil, nil), canvas.NewRectangle(color.RGBA{158, 158, 158, 0xff}), transactionInfoform)
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

		widget.Refresh(sendPageContents)
	})

	amountErrorLabel = canvas.NewText("", color.RGBA{237, 109, 71, 255})
	amountErrorLabel.TextSize = 14
	amountErrorLabel.Hide()

	transactionFeeBox := widget.NewHBox(widget.NewLabel("Transaction fee"), widgets.NewHSpacer(149), transactionFeeLabel, transactionSizeDropdown)

	amountEntryGroup = widget.NewGroup("Amount", fyne.NewContainerWithLayout(
		layout.NewFixedGridLayout(fyne.NewSize(widget.NewLabel("12345678.12345678").MinSize().Width, amountEntry.MinSize().Height)), amountEntry),
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
			transactionFeeLabel.SetText(fmt.Sprintf("- DCR"))
			totalCostLabel.SetText("- DCR")
			balanceAfterSendLabel.SetText("- DCR")
			transactionSize.SetText("0 bytes")

			nextButton.Disable()
			canvas.Refresh(sendPageContents)
			return
		}

		transactionAuthor.UpdateSendDestination(0, temporaryAddress, dcrlibwallet.AmountAtom(amountInFloat), false)

		feeAndSize, err := transactionAuthor.EstimateFeeAndSize()
		if err != nil {
			if err.Error() == "insufficient_balance" {
				amountErrorLabel.Text = "Insufficient balance"
				amountErrorLabel.Show()
				canvas.Refresh(amountErrorLabel)

			} else {
				showErrorLabel("Could not retrieve transaction fee and size")

				log.Println(fmt.Sprintf("could not retrieve transaction fee and size %s", err.Error()))
			}

			transactionFeeLabel.SetText(fmt.Sprintf("- DCR"))
			totalCostLabel.SetText("- DCR")
			balanceAfterSendLabel.SetText("- DCR")
			transactionSize.SetText("0 bytes")

			nextButton.Disable()
			widget.Refresh(sendPageContents)
			return
		}

		if !amountErrorLabel.Hidden {
			amountErrorLabel.Hide()
			canvas.Refresh(amountErrorLabel)
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
		widget.Refresh(sendPageContents)
	}

	costAndBalanceAfterSendBox := widget.NewVBox()
	totalCostLabel = widget.NewLabelWithStyle("- DCR", fyne.TextAlignLeading, fyne.TextStyle{})
	balanceAfterSendLabel = widget.NewLabelWithStyle("- DCR", fyne.TextAlignLeading, fyne.TextStyle{})
	costAndBalanceAfterSendBox.Append(widget.NewHBox(widget.NewLabel("Total cost"), layout.NewSpacer(), totalCostLabel))
	costAndBalanceAfterSendBox.Append(widget.NewHBox(widget.NewLabel("Balance after send"), layout.NewSpacer(), balanceAfterSendLabel))
	costAndBalanceAfterSendContainer := fyne.NewContainerWithLayout(layout.NewFixedGridLayout(fyne.NewSize(328, 48)), costAndBalanceAfterSendBox)

	nextButton = widgets.NewButton(color.RGBA{41, 112, 255, 255}, "Next", func() {
		if multiWallet.ConnectedPeers() <= 0 {
			showErrorLabel("Not Connected To Decred Network")

			return
		}

		// if sending to self
		var sendingAddress, selfSendingSelectedWalletName string
		var amountInFloat float64
		var err error

		amountInFloat, err = strconv.ParseFloat(amountEntry.Text, 64)
		if err != nil {
			showErrorLabel("Could not parse float")
			return
		}

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
				log.Println("could not get self sending account reason: " + err.Error())

				return
			}

		} else {
			sendingAddress = destinationAddressEntry.Text
		}

		transactionAuthor.UpdateSendDestination(0, sendingAddress, dcrlibwallet.AmountAtom(amountInFloat), false)

		confirmationWindow(amountEntry, destinationAddressEntry, icons[assets.DownArrow], icons[assets.Alert], icons[assets.Reveal], icons[assets.Conceal],
			window, selectedWallet.Name, selfSendingSelectedWalletName, totalCostLabel.Text, transactionFeeLabel.Text, balanceAfterSendLabel.Text,
			destinationAddressEntryGroup.Hidden, transactionAuthor, successLabelContainer)
	})

	nextButton.SetMinSize(nextButton.MinSize().Add(fyne.NewSize(0, 20)))
	nextButton.Disable()

	sendPageContents = widget.NewVBox(
		baseWidgets,
		widget.NewHBox(layout.NewSpacer(), successLabelContainer.Container, sendPage.errorLabel.Container, layout.NewSpacer()),
		sendingAccountGroup,
		widgets.NewVSpacer(8),
		destinationBox,
		widgets.NewVSpacer(8),
		widget.NewHBox(amountEntryGroup, widget.NewVBox(sendPage.spendableLabel)),
		widgets.NewVSpacer(12),
		costAndBalanceAfterSendContainer,
		widgets.NewVSpacer(15),
		nextButton.Container)

	pageContent = widget.NewHBox(widgets.NewHSpacer(10), sendPageContents)
	return
}

func confirmationWindow(amountEntry, destinationAddressEntry *widget.Entry, downArrow, alert, reveal, conceal fyne.Resource, window fyne.Window,
	selectedWalletName, sendingToSelfSelectedWalletName string, totalCostText, transactionFeeText,
	balanceAfterSendText string, sendingToSelf bool, transactionAuthor *dcrlibwallet.TxAuthor, showSuccess *widgets.Button) {

	var confirmationPagePopup *widget.PopUp

	confirmLabel := canvas.NewText("Confirm to send", color.Black)
	confirmLabel.TextStyle.Bold = true
	confirmLabel.TextSize = 20

	errorLabel := canvas.NewText("Failed to send. Please try again.", color.White)
	errorLabel.Alignment = fyne.TextAlignCenter
	errorBar := canvas.NewRectangle(color.RGBA{237, 109, 71, 255})
	errorBar.SetMinSize(errorLabel.MinSize().Add(fyne.NewSize(20, 16)))

	errorLabelContainer := fyne.NewContainerWithLayout(layout.NewBorderLayout(nil, nil, nil, nil), errorBar, errorLabel)
	errorLabelContainer.Hide()

	accountSelectionPopupHeader := widget.NewHBox(
		widgets.NewImageButton(theme.CancelIcon(), nil, func() { confirmationPagePopup.Hide() }),
		widgets.NewHSpacer(9),
		confirmLabel,
		widgets.NewHSpacer(170),
	)
	sendingSelectedWalletLabel := widget.NewLabelWithStyle(fmt.Sprintf("%s (%s)", sendPage.sendingSelectedAccountLabel.Text, selectedWalletName), fyne.TextAlignTrailing, fyne.TextStyle{Bold: true})

	trailingDotForAmount := strings.Split(amountEntry.Text, ".")
	// if amount is a float
	amountLabelBox := fyne.NewContainerWithLayout(layouts.NewHBox(0))
	if len(trailingDotForAmount) > 1 && len(trailingDotForAmount[1]) > 2 {
		trailingAmountLabel := canvas.NewText(trailingDotForAmount[1][2:]+" DCR", color.Black)
		trailingAmountLabel.TextStyle = fyne.TextStyle{Bold: true, Monospace: true}
		trailingAmountLabel.TextSize = 15

		leadingAmountLabel := canvas.NewText(trailingDotForAmount[0]+"."+trailingDotForAmount[1][:2], color.Black)
		leadingAmountLabel.TextStyle = fyne.TextStyle{Bold: true, Monospace: true}
		leadingAmountLabel.TextSize = 20

		amountLabelBox.AddObject(leadingAmountLabel)
		amountLabelBox.AddObject(trailingAmountLabel)

	} else {
		amountLabel := canvas.NewText(amountEntry.Text, color.Black)
		amountLabel.TextStyle = fyne.TextStyle{Bold: true, Monospace: true}
		amountLabel.TextSize = 20

		DCRLabel := canvas.NewText("DCR", color.Black)
		DCRLabel.TextStyle = fyne.TextStyle{Bold: true, Monospace: true}
		DCRLabel.TextSize = 15

		amountLabelBox.Layout = layouts.NewHBox(5)
		amountLabelBox.AddObject(amountLabel)
		amountLabelBox.AddObject(DCRLabel)
	}

	toDestination := "To destination address"
	destinationAddress := destinationAddressEntry.Text

	if sendingToSelf {
		toDestination = "To self"
		destinationAddress = sendPage.selfSendingSelectedAccountLabel.Text + " (" + sendingToSelfSelectedWalletName + ")"
	}

	sendButton := widgets.NewButton(color.RGBA{41, 112, 255, 255}, "Send "+amountEntry.Text+" DCR", func() {
		errorLabel := canvas.NewText("Wrong spending password. Please try again.", color.RGBA{237, 109, 71, 255})
		errorLabel.Alignment = fyne.TextAlignCenter
		errorLabel.TextSize = 12
		errorLabel.Hide()

		var confirmButton *widgets.Button

		walletPassword := widget.NewPasswordEntry()
		walletPassword.SetPlaceHolder("Spending Password")
		walletPassword.OnChanged = func(value string) {
			if value == "" {
				confirmButton.Disable()
			} else if confirmButton.Disabled() {
				confirmButton.Enable()
			}
		}

		var sendingPasswordPopup *widget.PopUp
		var popupContent *widget.Box

		cancelLabel := canvas.NewText("Cancel", color.RGBA{41, 112, 255, 255})
		cancelLabel.TextStyle.Bold = true

		cancelButton := widgets.NewClickableBox(widget.NewHBox(cancelLabel), func() {
			sendingPasswordPopup.Hide()
			confirmationPagePopup.Show()
		})

		confirmButton = widgets.NewButton(color.RGBA{41, 112, 255, 255}, "Confirm", func() {
			confirmButton.Disable()
			cancelButton.Disable()

			hash, err := transactionAuthor.Broadcast([]byte(walletPassword.Text))
			if err != nil {
				// do not exit password popup on invalid passphrase
				if err.Error() == "invalid_passphrase" {
					errorLabel.Show()
					// this is an hack as selective refresh to errorLabel doesn't work
					widget.Refresh(popupContent)
					confirmButton.Enable()
					cancelButton.Disable()
				} else {
					log.Println(err)
					errorLabelContainer.Show()
					sendingPasswordPopup.Hide()
					confirmationPagePopup.Show()
				}
				return
			}

			destinationAddressEntry.SetText("")
			amountEntry.SetText("")

			fmt.Println(hash)
			showSuccess.Container.Show()
			canvas.Refresh(showSuccess.Container)

			sendingPasswordPopup.Hide()

			time.AfterFunc(time.Second*5, func() {
				showSuccess.Container.Hide()
				canvas.Refresh(showSuccess.Container)
			})
		})
		confirmButton.SetMinSize(fyne.NewSize(91, 40))
		confirmButton.Disable()

		var passwordConceal *widgets.ImageButton
		passwordConceal = widgets.NewImageButton(reveal, nil, func() {
			if walletPassword.Password {
				passwordConceal.SetIcon(conceal)
				walletPassword.Password = false
			} else {
				passwordConceal.SetIcon(reveal)
				walletPassword.Password = true
			}
			// reveal texts
			walletPassword.SetText(walletPassword.Text)
		})

		popupContent = widget.NewHBox(
			widgets.NewHSpacer(24),
			widget.NewVBox(
				widgets.NewVSpacer(24),
				widget.NewLabelWithStyle("Confirm to send", fyne.TextAlignLeading, fyne.TextStyle{Bold: true}),
				widgets.NewVSpacer(40),
				fyne.NewContainerWithLayout(layouts.NewPasswordLayout(fyne.NewSize(312, walletPassword.MinSize().Height)), walletPassword, passwordConceal),
				errorLabel,
				widgets.NewVSpacer(20),
				widget.NewHBox(layout.NewSpacer(), cancelButton, widgets.NewHSpacer(24), confirmButton.Container),
				widgets.NewVSpacer(24),
			),
			widgets.NewHSpacer(24),
		)

		sendingPasswordPopup = widget.NewModalPopUp(popupContent, window.Canvas())
		sendingPasswordPopup.Show()
	})

	sendButton.SetMinSize(fyne.NewSize(312, 56))
	sendButton.SetTextSize(18)

	confirmationPageContent := widget.NewVBox(
		widgets.NewVSpacer(18),
		accountSelectionPopupHeader,
		widgets.NewVSpacer(18),
		canvas.NewLine(color.Black),
		widgets.NewVSpacer(8),
		widget.NewHBox(layout.NewSpacer(), errorLabelContainer, layout.NewSpacer()),
		widgets.NewVSpacer(16),
		widget.NewHBox(layout.NewSpacer(), widget.NewLabel("Sending from"), sendingSelectedWalletLabel, layout.NewSpacer()),
		widget.NewHBox(layout.NewSpacer(), amountLabelBox, layout.NewSpacer()),
		widgets.NewVSpacer(10),
		widget.NewHBox(layout.NewSpacer(), widget.NewIcon(downArrow), layout.NewSpacer()),
		widgets.NewVSpacer(10),
		widget.NewLabelWithStyle(toDestination, fyne.TextAlignCenter, fyne.TextStyle{Monospace: true}),
		widget.NewLabelWithStyle(destinationAddress, fyne.TextAlignCenter, fyne.TextStyle{Bold: true}),
		widgets.NewVSpacer(8),
		canvas.NewLine(color.RGBA{230, 234, 237, 255}),
		widget.NewHBox(widget.NewLabel("Transaction fee"), layout.NewSpacer(), widget.NewLabelWithStyle(transactionFeeText, fyne.TextAlignLeading, fyne.TextStyle{Bold: true})),
		canvas.NewLine(color.RGBA{230, 234, 237, 255}),
		widget.NewHBox(widget.NewLabel("Total cost"), layout.NewSpacer(), widget.NewLabelWithStyle(totalCostText, fyne.TextAlignLeading, fyne.TextStyle{Bold: true})),
		widget.NewHBox(widget.NewLabel("Balance after send"), layout.NewSpacer(), widget.NewLabelWithStyle(balanceAfterSendText, fyne.TextAlignLeading, fyne.TextStyle{Bold: true})),
		canvas.NewLine(color.RGBA{230, 234, 237, 255}),

		widget.NewHBox(layout.NewSpacer(),
			widget.NewIcon(alert), widget.NewLabelWithStyle("Your DCR will be sent and CANNOT be undone.", fyne.TextAlignLeading, fyne.TextStyle{Bold: true}), layout.NewSpacer()),
		sendButton.Container,
		widgets.NewVSpacer(18),
	)

	confirmationPagePopup = widget.NewModalPopUp(
		widget.NewHBox(widgets.NewHSpacer(16), confirmationPageContent, widgets.NewHSpacer(16)),
		window.Canvas())

	confirmationPagePopup.Show()
}

func createAccountDropdown(initFunction func(), accountLabel string, receiveAccountIcon, collapseIcon fyne.Resource,
	multiWallet *dcrlibwallet.MultiWallet, walletIDs []int, sendingSelectedWalletID *int,
	accountBoxes []*widget.Box, selectedAccountLabel *widget.Label,
	selectedAccountBalanceLabel *widget.Label, selectedWalletLabel *canvas.Text) (accountClickableBox *widgets.ClickableBox) {

	dropdownContent := widget.NewVBox()

	selectAccountBox := widget.NewHBox(
		widgets.NewHSpacer(15),
		widget.NewVBox(widgets.NewVSpacer(10), widget.NewIcon(receiveAccountIcon)),
		widgets.NewHSpacer(20),
		fyne.NewContainerWithLayout(layouts.NewVBox(12), selectedAccountLabel, selectedWalletLabel),
		widgets.NewHSpacer(30),
		widget.NewVBox(widgets.NewVSpacer(4), selectedAccountBalanceLabel),
		widgets.NewHSpacer(8),
		widget.NewVBox(widgets.NewVSpacer(6), widget.NewIcon(collapseIcon)),
	)

	// TODO make wallets and account in a scrollabel container
	dropdownContentWithScroller := widget.NewScrollContainer(dropdownContent)
	accountSelectionPopup := widget.NewPopUp(dropdownContentWithScroller, fyne.CurrentApp().Driver().AllWindows()[0].Canvas())

	accountSelectionPopupHeader := widget.NewHBox(
		widgets.NewHSpacer(16),
		widgets.NewImageButton(theme.CancelIcon(), nil, func() { accountSelectionPopup.Hide() }),
		widgets.NewHSpacer(16),
		widget.NewLabelWithStyle(accountLabel, fyne.TextAlignLeading, fyne.TextStyle{Bold: true}),
		layout.NewSpacer(),
	)

	dropdownContent.Append(widget.NewVBox(
		widgets.NewVSpacer(5),
		accountSelectionPopupHeader,
		widgets.NewVSpacer(5),
		canvas.NewLine(color.Black),
	))

	// we cant access the children of group widget, proposed hack is to
	// create a vertical box array where all accounts would be placed,
	// then when we want to hide checkmarks we call all children of accountbox and hide checkmark icon except selected
	for walletIndex, walletID := range walletIDs {
		getAllWalletAccountsInBox(initFunction, dropdownContent, selectedAccountLabel, selectedAccountBalanceLabel, selectedWalletLabel,
			multiWallet.WalletWithID(walletID), walletIndex, walletID, sendingSelectedWalletID, accountBoxes, receiveAccountIcon, accountSelectionPopup)
	}
	accountSelectionPopup.Hide()

	accountClickableBox = widgets.NewClickableBox(selectAccountBox, func() {
		accountSelectionPopup.Move(fyne.CurrentApp().Driver().AbsolutePositionForObject(
			accountClickableBox).Add(fyne.NewPos(0, accountClickableBox.Size().Height)))

		accountSelectionPopup.Resize(fyne.NewSize(dropdownContent.MinSize().Width+10, fyne.Max(dropdownContent.MinSize().Height, 60)))
		widget.Refresh(selectAccountBox)
		accountSelectionPopup.Show()
	})

	return
}

func getAllWalletAccountsInBox(initFunction func(), dropdownContent *widget.Box, selectedAccountLabel,
	selectedAccountBalanceLabel *widget.Label, selectedWalletLabel *canvas.Text, wallet *dcrlibwallet.Wallet, walletIndex, walletID int,
	sendingSelectedWalletID *int, accountsBoxes []*widget.Box, receiveIcon fyne.Resource, popup *widget.PopUp) {

	accounts, err := wallet.GetAccountsRaw(dcrlibwallet.DefaultRequiredConfirmations)
	if err != nil {
		return
	}

	var groupedWalletsAccounts = widget.NewGroup(wallet.Name)
	// we cant access children of a group so a box is used
	accountsBox := widget.NewVBox()

	for index, account := range accounts.Acc {
		if account.Name == "imported" {
			continue
		}

		spendableLabel := canvas.NewText("Spendable", color.Black)
		spendableLabel.TextSize = 10

		accountName := account.Name
		accountNameLabel := widget.NewLabel(accountName)
		accountNameLabel.Alignment = fyne.TextAlignLeading
		accountNameBox := widget.NewVBox(
			accountNameLabel,
			widget.NewHBox(widgets.NewHSpacer(1), spendableLabel),
		)

		spendableAmountLabel := canvas.NewText(dcrutil.Amount(account.Balance.Spendable).String(), color.Black)
		spendableAmountLabel.TextSize = 10
		spendableAmountLabel.Alignment = fyne.TextAlignTrailing

		amount := dcrutil.Amount(account.Balance.Total).String()
		accountBalance := amount
		accountBalanceLabel := widget.NewLabel(accountBalance)
		accountBalanceLabel.Alignment = fyne.TextAlignTrailing

		accountBalanceBox := widget.NewVBox(
			accountBalanceLabel,
			spendableAmountLabel,
		)

		checkmarkIcon := widget.NewIcon(theme.ConfirmIcon())
		var spacing fyne.CanvasObject
		if index != 0 || walletID != *sendingSelectedWalletID {
			checkmarkIcon.Hide()
			spacing = widgets.NewHSpacer(35)
		} else {
			spacing = widgets.NewHSpacer(15)
		}

		accountsView := widget.NewHBox(
			widgets.NewHSpacer(15),
			widget.NewIcon(receiveIcon),
			widgets.NewHSpacer(20),
			accountNameBox,
			layout.NewSpacer(),
			widgets.NewHSpacer(30),
			accountBalanceBox,
			widgets.NewHSpacer(30),
			checkmarkIcon,
			spacing,
		)

		accountsBox.Append(widgets.NewClickableBox(accountsView, func() {
			*sendingSelectedWalletID = walletID
			for _, boxes := range accountsBoxes {
				for _, objectsChild := range boxes.Children {
					if box, ok := objectsChild.(*widgets.ClickableBox); !ok {
						continue
					} else {
						if len(box.Children) != 10 {
							continue
						}

						if icon, ok := box.Children[8].(*widget.Icon); !ok {
							continue
						} else {
							icon.Hide()
						}
						if spacing, ok := box.Children[9].(*fyne.Container); !ok {
							continue
						} else {
							spacing.Layout = layout.NewFixedGridLayout(fyne.NewSize(35, 0))
							canvas.Refresh(spacing)
						}
					}

					canvas.Refresh(objectsChild)
				}
			}

			checkmarkIcon.Show()
			if spacing, ok := accountsView.Children[9].(*fyne.Container); !ok {
			} else {
				spacing.Layout = layout.NewFixedGridLayout(fyne.NewSize(15, 0))
				canvas.Refresh(spacing)
			}

			if accountbalanceBox, ok := accountsView.Children[6].(*widget.Box); ok {
				if len(accountbalanceBox.Children) == 2 {
					if balanceLabel, ok := accountbalanceBox.Children[0].(*widget.Label); ok {
						selectedAccountBalanceLabel.SetText(balanceLabel.Text)
					}
				}
			}

			selectedAccountLabel.SetText(accountName)
			selectedWalletLabel.Text = wallet.Name
			canvas.Refresh(selectedWalletLabel)
			fmt.Println("does it work", selectedWalletLabel.Text)

			if initFunction != nil {
				initFunction()
			}
			popup.Hide()
		}))
	}

	accountsBoxes[walletIndex] = accountsBox
	groupedWalletsAccounts.Append(accountsBoxes[walletIndex])
	dropdownContent.Append(groupedWalletsAccounts)
}

func updateContentOnNotification(accountBoxes []*widget.Box, sendingSelectedAccountLabel *widget.Label, selfSendingSelectedAccountBalanceLabel *widget.Label, multiWallet *dcrlibwallet.MultiWallet, selectedWalletID int) {
	selectedWalletIDs := multiWallet.OpenedWalletIDsRaw()
	sort.Ints(selectedWalletIDs)
	if len(selectedWalletIDs) != len(accountBoxes) {
		fmt.Println("Not working")
		return
	}

	for walletIndex, accountBox := range accountBoxes {
		wallet := multiWallet.WalletWithID(selectedWalletIDs[walletIndex])
		if wallet == nil {
			return
		}

		account, err := wallet.GetAccountsRaw(dcrlibwallet.DefaultRequiredConfirmations)
		if err != nil {
			log.Println("could not retrieve accounts on transaction notification")
			continue
		}

		if len(accountBox.Children) != len(account.Acc) {
			continue
		}

		for index, boxContent := range accountBox.Children {
			spendableAmountLabel := canvas.NewText(dcrutil.Amount(account.Acc[index].Balance.Spendable).String(), color.White)
			spendableAmountLabel.TextSize = 10
			spendableAmountLabel.Alignment = fyne.TextAlignTrailing

			accountBalance := dcrutil.Amount(account.Acc[index].Balance.Total).String()
			accountBalanceLabel := widget.NewLabel(accountBalance)
			accountBalanceLabel.Alignment = fyne.TextAlignTrailing

			accountBalanceBox := widget.NewVBox(
				accountBalanceLabel,
				spendableAmountLabel,
			)

			accountBalance = dcrutil.Amount(account.Acc[index].Balance.Total).String()

			if content, ok := boxContent.(*widgets.ClickableBox); ok {
				fmt.Println("worksss")
				content.Box.Children[6] = accountBalanceBox
				widget.Refresh(content.Box)
				widget.Refresh(content)
			}
		}
	}

	wallet := multiWallet.WalletWithID(selectedWalletID)
	if wallet == nil {
		log.Println("could not retrieve selected wallet on transaction notification")
		return
	}

	accountNumber, err := wallet.AccountNumber(sendingSelectedAccountLabel.Text)
	if err != nil {
		log.Println("could not retrieve selected account number on transaction notification")
		return
	}

	account, err := wallet.GetAccount(int32(accountNumber), dcrlibwallet.DefaultRequiredConfirmations)
	if err != nil {
		log.Println("could not retrieve selected account on transaction notification")
		return
	}

	selfSendingSelectedAccountBalanceLabel.SetText(dcrutil.Amount(account.TotalBalance).String())
}
