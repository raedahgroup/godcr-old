package pages

import (
	"fmt"
	"github.com/raedahgroup/godcr/fyne/layouts"
	"image/color"
	"log"
	"regexp"
	"sort"
	"strconv"
	"strings"

	"fyne.io/fyne"
	"fyne.io/fyne/canvas"
	"fyne.io/fyne/layout"
	"fyne.io/fyne/theme"
	"fyne.io/fyne/widget"

	"github.com/decred/dcrd/dcrutil"
	"github.com/raedahgroup/dcrlibwallet"

	"github.com/raedahgroup/godcr/fyne/assets"
	"github.com/raedahgroup/godcr/fyne/widgets"
)

type sendPageDynamicData struct {
	// houses all clickable box
	sendingAccountBoxes     []*widget.Box
	selfSendingAccountBoxes []*widget.Box

	errorLabel     *canvas.Text
	spendableLabel *canvas.Text

	selfSendingSelectedAccountLabel        *widget.Label
	selfSendingSelectedAccountBalanceLabel *widget.Label

	sendingSelectedAccountLabel        *widget.Label
	sendingSelectedAccountBalanceLabel *widget.Label
}

var sendPage sendPageDynamicData

func sendPageContent(multiWallet *dcrlibwallet.MultiWallet, window fyne.Window) (pageContent fyne.CanvasObject) {
	icons, err := assets.GetIcons(assets.Reveal, assets.Conceal, assets.InfoIcon, assets.MoreIcon, assets.ReceiveAccountIcon, assets.CollapseIcon, assets.CollapseDropdown, assets.ExpandDropdown, assets.DownArrow, assets.Alert)
	if err != nil {
		return widget.NewLabelWithStyle(err.Error(), fyne.TextAlignCenter, fyne.TextStyle{Bold: true})
	}

	sendPage.errorLabel = canvas.NewText("", color.RGBA{237, 109, 71, 255})

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

	var sendingSelectedWalletID = openedWalletIDs[0]
	var sendingToSelfSelectedWalletID = openedWalletIDs[0]

	var selectedWallet = multiWallet.WalletWithID(sendingSelectedWalletID)

	selectedWalletAccounts, err := selectedWallet.GetAccountsRaw(dcrlibwallet.DefaultRequiredConfirmations)
	if err != nil {
		return widget.NewLabel(fmt.Sprintf("Error: %s", err.Error()))
	}

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
	var form *fyne.Container

	// this function is called when the sending wallet account is changed.
	onSendingAccountChange := func() {
		selectedWallet = multiWallet.WalletWithID(sendingSelectedWalletID)
		accountNumber, err := selectedWallet.AccountNumber(sendPage.sendingSelectedAccountLabel.Text)
		if err != nil {
			sendPage.errorLabel.Text = "Could not get account, " + err.Error()
			sendPage.errorLabel.Show()
			canvas.Refresh(sendPage.errorLabel)
		}

		transactionAuthor = selectedWallet.NewUnsignedTx(int32(accountNumber), dcrlibwallet.DefaultRequiredConfirmations)
		fmt.Println("Changed account", selectedWallet.Name, sendPage.sendingSelectedAccountLabel.Text)
		transactionAuthor.AddSendDestination(temporaryAddress, 0, true)

		sendPage.spendableLabel.Text = "Spendable: " + sendPage.sendingSelectedAccountBalanceLabel.Text

		balance, err := selectedWallet.GetAccountBalance(int32(accountNumber), dcrlibwallet.DefaultRequiredConfirmations)
		if err != nil {
			log.Println("could not retrieve account balance")
			return
		}
		amountInAccount = dcrlibwallet.AmountCoin(balance.Total)
		// reset amount entry
		amountEntry.OnChanged(amountEntry.Text)

		canvas.Refresh(form)
		canvas.Refresh(transactionFeeLabel)
		canvas.Refresh(totalCostLabel)
		canvas.Refresh(balanceAfterSendLabel)
		canvas.Refresh(transactionSize)
		canvas.Refresh(sendPage.errorLabel)
		canvas.Refresh(sendPage.spendableLabel)
	}

	// we still need a suitable name for this
	sendPage.sendingSelectedAccountLabel = widget.NewLabel(selectedWalletAccounts.Acc[0].Name)
	sendPage.sendingSelectedAccountBalanceLabel = widget.NewLabel(dcrutil.Amount(selectedWalletAccounts.Acc[0].TotalBalance).String())

	sendingAccountDropdownContent := widget.NewVBox()
	sendingAccountClickableBox := createAccountDropdown(onSendingAccountChange, "Sending account", icons[assets.ReceiveAccountIcon], icons[assets.CollapseIcon], multiWallet, openedWalletIDs, &sendingSelectedWalletID, sendPage.sendingAccountBoxes, sendingAccountDropdownContent, sendPage.sendingSelectedAccountLabel, sendPage.sendingSelectedAccountBalanceLabel)
	sendingAccountGroup := widget.NewGroup("From", sendingAccountClickableBox)

	sendPage.selfSendingSelectedAccountLabel = widget.NewLabel(selectedWalletAccounts.Acc[0].Name)
	sendPage.selfSendingSelectedAccountBalanceLabel = widget.NewLabel(dcrutil.Amount(selectedWalletAccounts.Acc[0].TotalBalance).String())

	selfSendingAccountDropdownContent := widget.NewVBox()
	selfSendingToAccountClickableBox := createAccountDropdown(nil, "Receiving account", icons[assets.ReceiveAccountIcon], icons[assets.CollapseIcon], multiWallet, openedWalletIDs, &sendingToSelfSelectedWalletID, sendPage.selfSendingAccountBoxes, selfSendingAccountDropdownContent, sendPage.selfSendingSelectedAccountLabel, sendPage.selfSendingSelectedAccountBalanceLabel)
	selfSendingToAccountGroup := widget.NewGroup("To", selfSendingToAccountClickableBox) //sendingToAccountClickableBox)
	selfSendingToAccountGroup.Hide()

	destinationAddressEntry := widget.NewEntry()
	destinationAddressEntry.SetPlaceHolder("Destination Address")

	// shows errors related too destination address
	destinationAddressErrorLabel := canvas.NewText("", color.RGBA{237, 109, 71, 255})
	destinationAddressErrorLabel.TextSize = 12
	destinationAddressErrorLabel.Hide()

	var nextButton *widget.Button

	destinationAddressEntry.OnChanged = func(address string) {
		nextButton.Disable()

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
			widget.Refresh(nextButton)
		}
		canvas.Refresh(destinationAddressErrorLabel)
	}

	destinationAddressEntryGroup := widget.NewGroup("To", fyne.NewContainerWithLayout(
		layout.NewFixedGridLayout(fyne.NewSize(widget.NewLabel(temporaryAddress).MinSize().Width, destinationAddressEntry.MinSize().Height)), destinationAddressEntry),
		destinationAddressErrorLabel)

	sendToAccountLabel := canvas.NewText("Send to account", color.RGBA{R: 41, G: 112, B: 255, A: 255})
	sendToAccountLabel.TextSize = 12

	sendToAccount = widgets.NewClickableBox(widget.NewVBox(sendToAccountLabel), func() {
		if sendToAccountLabel.Text == "Send to account" {
			sendToAccountLabel.Text = "Send to address"
			canvas.Refresh(sendToAccountLabel)
			widget.Refresh(sendToAccount)
			selfSendingToAccountGroup.Show()
			destinationAddressEntryGroup.Hide()
		} else {
			sendToAccountLabel.Text = "Send to account"
			canvas.Refresh(sendToAccountLabel)
			widget.Refresh(sendToAccount)
			destinationAddressEntryGroup.Show()
			selfSendingToAccountGroup.Hide()
		}

		amountEntry.OnChanged(amountEntry.Text)
	})

	amountEntry = widget.NewEntry()
	amountEntry.SetPlaceHolder("0 DCR")

	transactionFeeLabel = widget.NewLabel("- DCR")
	transactionSize = widget.NewLabelWithStyle("0 bytes", fyne.TextAlignLeading, fyne.TextStyle{})
	form = fyne.NewContainerWithLayout(layout.NewVBoxLayout())
	form.AddObject(fyne.NewContainerWithLayout(layout.NewHBoxLayout(), widget.NewLabel("Processing time"), widgets.NewHSpacer(46), layout.NewSpacer(), widget.NewLabelWithStyle("Approx. 10 mins (2 blocks)", fyne.TextAlignLeading, fyne.TextStyle{})))
	form.AddObject(fyne.NewContainerWithLayout(layout.NewHBoxLayout(), widget.NewLabel("Fee rate"), layout.NewSpacer(), widget.NewLabelWithStyle("0.0001 DCR/byte", fyne.TextAlignLeading, fyne.TextStyle{})))
	form.AddObject(fyne.NewContainerWithLayout(layout.NewHBoxLayout(), widget.NewLabel("Transaction size"), layout.NewSpacer(), transactionSize))

	paintedForm := fyne.NewContainerWithLayout(layout.NewBorderLayout(nil, nil, nil, nil), canvas.NewRectangle(color.RGBA{158, 158, 158, 0xff}), form)
	paintedForm.Hide()

	var transactionSizeDropdown *widgets.ClickableBox
	transactionSizeDropdown = widgets.NewClickableBox(widget.NewHBox(widget.NewIcon(icons[assets.ExpandDropdown])), func() {
		if paintedForm.Hidden {
			transactionSizeDropdown.Box.Children[0] = widget.NewIcon(icons[assets.CollapseDropdown])
			paintedForm.Show()
			canvas.Refresh(transactionSizeDropdown)
			canvas.Refresh(paintedForm)
		} else {
			transactionSizeDropdown.Box.Children[0] = widget.NewIcon(icons[assets.ExpandDropdown])
			paintedForm.Hide()
			canvas.Refresh(transactionSizeDropdown)
			canvas.Refresh(paintedForm)
		}
	})

	amountErrorLabel = canvas.NewText("", color.RGBA{237, 109, 71, 255})
	amountErrorLabel.TextSize = 14
	amountErrorLabel.Hide()

	amountEntryGroup := widget.NewGroup("Amount", fyne.NewContainerWithLayout(
		layout.NewFixedGridLayout(fyne.NewSize(widget.NewLabel("12345678.12345678").MinSize().Width, amountEntry.MinSize().Height)), amountEntry),
		amountErrorLabel, widgets.NewVSpacer(4),
		widget.NewHBox(widget.NewLabel("Transaction fee"), widgets.NewHSpacer(149), transactionFeeLabel, transactionSizeDropdown),
		paintedForm)

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

		nextButton.Disable()

		if numbers := strings.Split(value, "."); len(numbers) == 2 {
			if len(numbers[1]) > 8 {
				sendPage.errorLabel.Text = "Amount has more than 8 decimal places"
				sendPage.errorLabel.Show()
				canvas.Refresh(sendPage.errorLabel)

				return
			}
		}

		amountInFloat, err := strconv.ParseFloat(value, 64)
		if err != nil && value != "" {
			sendPage.errorLabel.Text = "Could not parse float"
			sendPage.errorLabel.Show()
			canvas.Refresh(sendPage.errorLabel)
			return
		}

		if amountInFloat == 0.0 {
			transactionFeeLabel.SetText(fmt.Sprintf("- DCR"))
			totalCostLabel.SetText("- DCR")
			balanceAfterSendLabel.SetText("- DCR")
			transactionSize.SetText("0 bytes")

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
				sendPage.errorLabel.Text = "could not retrieve transaction fee and size"
				sendPage.errorLabel.Show()
				canvas.Refresh(sendPage.errorLabel)

				log.Println(fmt.Sprintf("could not retrieve transaction fee and size %s", err.Error()))
			}

			transactionFeeLabel.SetText(fmt.Sprintf("- DCR"))
			totalCostLabel.SetText("- DCR")
			balanceAfterSendLabel.SetText("- DCR")
			transactionSize.SetText("0 bytes")

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

		widget.Refresh(amountEntryGroup)
		canvas.Refresh(paintedForm)

		if destinationAddressEntry.Text != "" && destinationAddressErrorLabel.Hidden || destinationAddressEntryGroup.Hidden {
			nextButton.Enable()
			widget.Refresh(nextButton)
		}

		sendPage.errorLabel.Hide()
		canvas.Refresh(sendPage.errorLabel)
	}

	costAndBalanceAfterSendBox := widget.NewVBox()
	totalCostLabel = widget.NewLabelWithStyle("- DCR", fyne.TextAlignLeading, fyne.TextStyle{})
	balanceAfterSendLabel = widget.NewLabelWithStyle("- DCR", fyne.TextAlignLeading, fyne.TextStyle{})
	costAndBalanceAfterSendBox.Append(widget.NewHBox(widget.NewLabel("Total cost"), layout.NewSpacer(), totalCostLabel))
	costAndBalanceAfterSendBox.Append(widget.NewHBox(widget.NewLabel("Balance after send"), layout.NewSpacer(), balanceAfterSendLabel))
	costAndBalanceAfterSendContainer := fyne.NewContainerWithLayout(layout.NewFixedGridLayout(fyne.NewSize(328, 48)), costAndBalanceAfterSendBox)

	nextButton = widget.NewButton("Next", func() {
		var confirmationPagePopup *widget.PopUp

		initiateAfterPopupClose := func() {
			confirmationPagePopup.Hide()
		}

		confirmLabel := canvas.NewText("Confirm to send", color.Black)
		confirmLabel.TextStyle.Bold = true
		confirmLabel.TextSize = 20

		accountSelectionPopupHeader := widget.NewHBox(
			widgets.NewImageButton(theme.CancelIcon(), nil, initiateAfterPopupClose),
			widgets.NewHSpacer(9),
			confirmLabel,
			widgets.NewHSpacer(170),
		)
		sendingSelectedWalletLabel := widget.NewLabelWithStyle(fmt.Sprintf("%s (%s)", sendPage.sendingSelectedAccountLabel.Text, multiWallet.WalletWithID(sendingToSelfSelectedWalletID).Name), fyne.TextAlignTrailing, fyne.TextStyle{Bold: true})

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

		if destinationAddressEntryGroup.Hidden {
			toDestination = "To self"
			destinationAddress = sendPage.selfSendingSelectedAccountLabel.Text + " (" + multiWallet.WalletWithID(sendingToSelfSelectedWalletID).Name + ")"
		}

		blueBar := canvas.NewRectangle(color.RGBA{41, 112, 255, 255})
		blueBar.SetMinSize(fyne.NewSize(312, 56))
		sendLabel := canvas.NewText("Send "+amountEntry.Text+" DCR", color.White)
		sendLabel.Alignment = fyne.TextAlignCenter
		sendLabel.TextSize = 18

		confirmationPageContent := widget.NewVBox(
			widgets.NewVSpacer(18),
			accountSelectionPopupHeader,
			widgets.NewVSpacer(18),
			canvas.NewLine(color.Black),
			widgets.NewVSpacer(24),
			widget.NewHBox(layout.NewSpacer(), widget.NewLabel("Sending from"), sendingSelectedWalletLabel, layout.NewSpacer()),
			widget.NewHBox(layout.NewSpacer(), amountLabelBox, layout.NewSpacer()),
			widgets.NewVSpacer(10),
			widget.NewHBox(layout.NewSpacer(), widget.NewIcon(icons[assets.DownArrow]), layout.NewSpacer()),
			widgets.NewVSpacer(10),
			widget.NewLabelWithStyle(toDestination, fyne.TextAlignCenter, fyne.TextStyle{Monospace: true}),
			widget.NewLabelWithStyle(destinationAddress, fyne.TextAlignCenter, fyne.TextStyle{Bold: true}),
			widgets.NewVSpacer(8),
			canvas.NewLine(color.RGBA{230, 234, 237, 255}),
			widget.NewHBox(widget.NewLabel("Transaction fee"), layout.NewSpacer(), widget.NewLabelWithStyle(transactionFeeLabel.Text, fyne.TextAlignLeading, fyne.TextStyle{Bold: true})),
			canvas.NewLine(color.RGBA{230, 234, 237, 255}),
			widget.NewHBox(widget.NewLabel("Total cost"), layout.NewSpacer(), widget.NewLabelWithStyle(totalCostLabel.Text, fyne.TextAlignLeading, fyne.TextStyle{Bold: true})),
			widget.NewHBox(widget.NewLabel("Balance after send"), layout.NewSpacer(), widget.NewLabelWithStyle(balanceAfterSendLabel.Text, fyne.TextAlignLeading, fyne.TextStyle{Bold: true})),
			canvas.NewLine(color.RGBA{230, 234, 237, 255}),
			widget.NewHBox(layout.NewSpacer(), widget.NewIcon(icons[assets.Alert]), widget.NewLabelWithStyle("Your DCR will be send and CANNOT be undone.", fyne.TextAlignLeading, fyne.TextStyle{Bold: true}), layout.NewSpacer()),

			widgets.NewClickableBox(
				widget.NewVBox(
					fyne.NewContainerWithLayout(layout.NewBorderLayout(nil, nil, nil, nil), blueBar, sendLabel)),

				func() {
					confirmButtonBar := canvas.NewRectangle(color.RGBA{196, 203, 210, 255})
					confirmButtonBar.SetMinSize(fyne.NewSize(91, 40))

					confirmButtonLabel := canvas.NewText("Confirm", color.White)
					confirmButtonLabel.TextStyle.Bold = true
					confirmButtonLabel.Alignment = fyne.TextAlignCenter

					walletPassword := widget.NewPasswordEntry()
					walletPassword.SetPlaceHolder("Spending Password")
					walletPassword.OnChanged = func(value string) {
						if value == "" {
							confirmButtonBar.FillColor = color.RGBA{196, 203, 210, 255}
							canvas.Refresh(confirmButtonBar)
						} else if confirmButtonBar.FillColor != blueBar.FillColor {
							confirmButtonBar.FillColor = blueBar.FillColor
							canvas.Refresh(confirmButtonBar)
						}
					}

					cancelLabel := canvas.NewText("Cancel", color.RGBA{41, 112, 255, 255})
					cancelLabel.TextStyle.Bold = true

					cancelButton := widgets.NewClickableBox(widget.NewHBox(cancelLabel), func() {
						walletPassword.SetText("")
						confirmationPagePopup.Show()
					})

					var sendingPasswordPopup *widget.PopUp

					var confirmButton *widgets.ClickableBox
					confirmButton = widgets.NewClickableBox(
						widget.NewVBox(fyne.NewContainerWithLayout(layout.NewBorderLayout(nil, nil, nil, nil), confirmButtonBar, confirmButtonLabel)), func() {
							confirmButton.OnTapped = nil
							cancelButton.OnTapped = nil
							fmt.Println("Works")
							sendingPasswordPopup.Hide()
						})

					var passwordConceal *widgets.ImageButton
					passwordConceal = widgets.NewImageButton(icons[assets.Reveal], nil, func() {
						if walletPassword.Password {
							passwordConceal.SetIcon(icons[assets.Conceal])
							walletPassword.Password = false
						} else {
							passwordConceal.SetIcon(icons[assets.Reveal])
							walletPassword.Password = true
						}
						// reveal texts
						walletPassword.SetText(walletPassword.Text)
					})

					popupContent := widget.NewHBox(
						widgets.NewHSpacer(24),
						widget.NewVBox(
							widgets.NewVSpacer(24),
							widget.NewLabelWithStyle("Confirm to send", fyne.TextAlignLeading, fyne.TextStyle{Bold: true}),
							widgets.NewVSpacer(40),
							fyne.NewContainerWithLayout(layouts.NewPasswordLayout(fyne.NewSize(312, walletPassword.MinSize().Height)), walletPassword, passwordConceal),
							widgets.NewVSpacer(20),
							widget.NewHBox(layout.NewSpacer(), cancelButton, widgets.NewHSpacer(24), confirmButton),
							widgets.NewVSpacer(24),
						),
						widgets.NewHSpacer(24),
					)

					sendingPasswordPopup = widget.NewModalPopUp(popupContent, window.Canvas())
					sendingPasswordPopup.Show()
				}),
			widgets.NewVSpacer(18),
		)

		confirmationPagePopup = widget.NewModalPopUp(
			widget.NewHBox(widgets.NewHSpacer(16), confirmationPageContent, widgets.NewHSpacer(16)),
			window.Canvas())

		confirmationPagePopup.Show()
	})
	nextButton.Disable()

	sendPageContents := widget.NewVBox(baseWidgets, sendingAccountGroup,
		widgets.NewVSpacer(8),
		widget.NewHBox(destinationAddressEntryGroup, selfSendingToAccountGroup, widget.NewVBox(sendToAccount)), // this is an hack to center mouse inputs
		widgets.NewVSpacer(8),
		widget.NewHBox(amountEntryGroup, widget.NewVBox(sendPage.spendableLabel)),
		widgets.NewVSpacer(12),
		costAndBalanceAfterSendContainer,
		widgets.NewVSpacer(16),
		sendPage.errorLabel,
		widgets.NewVSpacer(10),
		nextButton)

	pageContent = widget.NewHBox(widgets.NewHSpacer(10), sendPageContents)
	return
}

func createAccountDropdown(initFunction func(), accountLabel string, receiveAccountIcon, collapseIcon fyne.Resource,
	multiWallet *dcrlibwallet.MultiWallet, walletIDs []int, sendingSelectedWalletID *int,
	accountBoxes []*widget.Box, dropdownContent *widget.Box, selectedAccountLabel *widget.Label,
	selectedAccountBalanceLabel *widget.Label) (accountClickableBox *widgets.ClickableBox) {

	selectAccountBox := widget.NewHBox(
		widgets.NewHSpacer(15),
		widget.NewIcon(receiveAccountIcon),
		widgets.NewHSpacer(20),
		selectedAccountLabel,
		widgets.NewHSpacer(30),
		selectedAccountBalanceLabel,
		widgets.NewHSpacer(8),
		widget.NewIcon(collapseIcon),
	)

	// TODO make wallets and account in a scrollabel container
	accountSelectionPopup := widget.NewPopUp(dropdownContent, fyne.CurrentApp().Driver().AllWindows()[0].Canvas())
	accountBoxes = make([]*widget.Box, len(walletIDs))

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
		getAllWalletAccountsInBox(initFunction, dropdownContent, selectedAccountLabel, selectedAccountBalanceLabel,
			multiWallet.WalletWithID(walletID), walletIndex, walletID, sendingSelectedWalletID, accountBoxes, receiveAccountIcon, accountSelectionPopup)
	}
	accountSelectionPopup.Hide()

	accountClickableBox = widgets.NewClickableBox(selectAccountBox, func() {
		accountSelectionPopup.Move(fyne.CurrentApp().Driver().AbsolutePositionForObject(
			accountClickableBox).Add(fyne.NewPos(0, accountClickableBox.Size().Height)))
		accountSelectionPopup.Show()
	})
	return
}

func getAllWalletAccountsInBox(initFunction func(), dropdownContent *widget.Box, selectedAccountLabel,
	selectedAccountBalanceLabel *widget.Label, wallet *dcrlibwallet.Wallet, walletIndex, walletID int,
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

func updateAccountDropdownContent(accountBox *widget.Box, account *dcrlibwallet.Accounts) {
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
