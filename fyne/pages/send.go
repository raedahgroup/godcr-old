package pages

import (
	"fmt"
	"image/color"
	"log"
	"regexp"
	"strconv"
	"strings"

	"github.com/raedahgroup/dcrlibwallet"

	"fyne.io/fyne"
	"fyne.io/fyne/canvas"
	"fyne.io/fyne/layout"
	"fyne.io/fyne/theme"
	"fyne.io/fyne/widget"
	"github.com/decred/dcrd/dcrutil"
	"github.com/raedahgroup/godcr/fyne/assets"
	"github.com/raedahgroup/godcr/fyne/widgets"
)

type sendPageDynamicData struct {
	// houses all clickable box
	receivingAccountDropdownContent *widget.Box
	sendingAccountDropdownContent   *widget.Box

	errorLabel                           *canvas.Text
	spendableLabel                       *canvas.Text
	sendingSelectedAccountLabel          *widget.Label
	sendingSelectedAccountBalanceLabel   *widget.Label
	receivingSelectedAccountLabel        *widget.Label
	receivingSelectedAccountBalanceLabel *widget.Label
}

var sendPage sendPageDynamicData

func sendPageContent(dcrlw *dcrlibwallet.LibWallet) fyne.CanvasObject {
	icons, err := assets.GetIcons(assets.InfoIcon, assets.MoreIcon, assets.ReceiveAccountIcon, assets.CollapseIcon, assets.CollapseDropdown, assets.ExpandDropdown)
	if err != nil {
		return widget.NewLabelWithStyle(err.Error(), fyne.TextAlignCenter, fyne.TextStyle{Bold: true})
	}

	sendPage.errorLabel = canvas.NewText("", color.RGBA{255, 0, 0, 0})

	// define base widget consisting of label, more icon and info button
	sendLabel := widget.NewLabelWithStyle("Send DCR", fyne.TextAlignTrailing, fyne.TextStyle{Bold: true, Italic: true})
	clickabelInfoIcon := widgets.NewImageButton(icons[assets.InfoIcon], nil, func() {

	})
	clickabelMoreIcon := widgets.NewImageButton(icons[assets.MoreIcon], nil, func() {

	})

	baseWidgets := widget.NewHBox(sendLabel, layout.NewSpacer(), clickabelInfoIcon, clickabelMoreIcon)

	accounts, err := dcrlw.GetAccountsRaw(dcrlibwallet.DefaultRequiredConfirmations)
	if err != nil {
		return widget.NewLabelWithStyle("could not retrieve account, "+err.Error(), fyne.TextAlignCenter, fyne.TextStyle{Bold: true})
	}

	sendPage.spendableLabel = canvas.NewText("Spendable: "+dcrutil.Amount(accounts.Acc[0].TotalBalance).String(), color.White)
	sendPage.spendableLabel.TextSize = 12

	accountNumber, err := dcrlw.AccountNumber("default")
	if err != nil {
		log.Println("could not retrieve account details", err.Error())
		return widget.NewLabel("could not retrieve account details")
	}

	temporaryAddress, err := dcrlw.CurrentAddress(int32(accountNumber))
	if err != nil {
		log.Println("could not retrieve account details", err.Error())
		return widget.NewLabel("could not retrieve account details")
	}
	amountInAccount := dcrlibwallet.AmountCoin(accounts.Acc[0].TotalBalance)

	transactionAuthor := dcrlw.NewUnsignedTx(int32(accountNumber), dcrlibwallet.DefaultRequiredConfirmations)
	transactionAuthor.AddSendDestination(temporaryAddress, 0, true)

	onAccountChange := func() {
		acct, err := dcrlw.AccountNumber(sendPage.receivingSelectedAccountLabel.Text)
		if err != nil {
			sendPage.errorLabel.Text = "Could not get account, " + err.Error()
		}

		transactionAuthor.SetSourceAccount(int32(acct))
		temporaryAddress, err = dcrlw.CurrentAddress(int32(acct))
		if err != nil {
			log.Println("could not retrieve account details", err.Error())
			return
		}
		fmt.Println("Changed account")
		transactionAuthor.UpdateSendDestination(0, temporaryAddress, 0, true)

		sendPage.spendableLabel.Text = "Spendable: " + sendPage.receivingSelectedAccountBalanceLabel.Text
		canvas.Refresh(sendPage.spendableLabel)

		balance, err := dcrlw.GetAccountBalance(int32(acct), dcrlibwallet.DefaultRequiredConfirmations)
		if err != nil {
			log.Println("could not retrieve account balance")
			return
		}
		amountInAccount = dcrlibwallet.AmountCoin(balance.Total)
	}

	// we still need a suitable name for this
	sendPage.receivingSelectedAccountLabel = widget.NewLabel(accounts.Acc[0].Name)
	sendPage.receivingSelectedAccountBalanceLabel = widget.NewLabel(dcrutil.Amount(accounts.Acc[0].TotalBalance).String())
	sendPage.receivingAccountDropdownContent = widget.NewVBox()
	receivingAccountClickableBox := createAccountDropdown(onAccountChange, icons[assets.ReceiveAccountIcon], icons[assets.CollapseIcon], accounts, sendPage.receivingAccountDropdownContent, sendPage.receivingSelectedAccountLabel, sendPage.receivingSelectedAccountBalanceLabel)
	receivingAccountGroup := widget.NewGroup("From", receivingAccountClickableBox)

	sendPage.sendingSelectedAccountLabel = widget.NewLabel(accounts.Acc[0].Name)
	sendPage.sendingSelectedAccountBalanceLabel = widget.NewLabel(dcrutil.Amount(accounts.Acc[0].TotalBalance).String())
	sendPage.sendingAccountDropdownContent = widget.NewVBox()
	sendingToAccountClickableBox := createAccountDropdown(nil, icons[assets.ReceiveAccountIcon], icons[assets.CollapseIcon], accounts, sendPage.sendingAccountDropdownContent, sendPage.sendingSelectedAccountLabel, sendPage.sendingSelectedAccountBalanceLabel)
	sendingToAccountGroup := widget.NewGroup("To", sendingToAccountClickableBox) //sendingToAccountClickableBox)
	sendingToAccountGroup.Hide()

	destinationAddressEntry := widget.NewEntry()
	destinationAddressEntry.SetPlaceHolder("Destination Address")
	// shows errors related too destination address
	destinationAddressErrorLabel := canvas.NewText("", color.RGBA{237, 109, 71, 255})
	destinationAddressErrorLabel.TextSize = 12
	destinationAddressErrorLabel.Hide()

	destinationAddressEntry.OnChanged = func(address string) {
		_, err := dcrutil.DecodeAddress(address)
		if err != nil {
			destinationAddressErrorLabel.Text = "Invalid address"
			destinationAddressErrorLabel.Show()
		} else {
			destinationAddressErrorLabel.Hide()
		}
		canvas.Refresh(destinationAddressErrorLabel)
	}

	destinationAddressEntryGroup := widget.NewGroup("To", fyne.NewContainerWithLayout(
		layout.NewFixedGridLayout(fyne.NewSize(widget.NewLabel("TsfDLrRkk9ciUuwfp2b8PawwnukYD7yAjGd").MinSize().Width, destinationAddressEntry.MinSize().Height)), destinationAddressEntry),
		destinationAddressErrorLabel)

	sendToAccountLabel := canvas.NewText("Send to account", color.RGBA{R: 41, G: 112, B: 255, A: 255})
	sendToAccountLabel.TextSize = 14
	sendToAccount := widgets.NewClickableBox(widget.NewVBox(sendToAccountLabel), func() {
		if sendToAccountLabel.Text == "Send to account" {
			sendToAccountLabel.Text = "Send to address"
			canvas.Refresh(sendToAccountLabel)
			sendingToAccountGroup.Show()
			destinationAddressEntryGroup.Hide()
		} else {
			sendToAccountLabel.Text = "Send to account"
			canvas.Refresh(sendToAccountLabel)
			destinationAddressEntryGroup.Show()
			sendingToAccountGroup.Hide()
		}
	})

	amountEntry := widget.NewEntry()
	amountEntry.SetPlaceHolder("0 DCR")

	transactionFeeLabel := widget.NewLabel("Transaction fee                - DCR")
	transactionSize := widget.NewLabelWithStyle("0 bytes", fyne.TextAlignLeading, fyne.TextStyle{})
	form := fyne.NewContainerWithLayout(layout.NewVBoxLayout())
	form.AddObject(fyne.NewContainerWithLayout(layout.NewHBoxLayout(), widget.NewLabel("Processing time"), widgets.NewHSpacer(46), layout.NewSpacer(), widget.NewLabelWithStyle("Approx. 10 mins (2 blocks)", fyne.TextAlignLeading, fyne.TextStyle{})))
	form.AddObject(fyne.NewContainerWithLayout(layout.NewHBoxLayout(), widget.NewLabel("Fee rate"), layout.NewSpacer(), widget.NewLabelWithStyle("0.0001 DCR/byte", fyne.TextAlignLeading, fyne.TextStyle{})))
	form.AddObject(fyne.NewContainerWithLayout(layout.NewHBoxLayout(), widget.NewLabel("Transaction size"), layout.NewSpacer(), transactionSize))

	paintedForm := fyne.NewContainerWithLayout(layout.NewBorderLayout(nil, nil, nil, nil), canvas.NewRectangle(color.RGBA{158, 158, 158, 0xff}), form)
	paintedForm.Hide()

	var transactionSizeDropdown *widgets.ClickableBox
	transactionSizeDropdown = widgets.NewClickableBox(widget.NewVBox(widget.NewIcon(icons[assets.ExpandDropdown])), func() {
		if paintedForm.Hidden {
			transactionSizeDropdown.Box.Children[0] = widget.NewIcon(icons[assets.CollapseDropdown])
			canvas.Refresh(transactionSizeDropdown.Box.Children[0])
			paintedForm.Show()
		} else {
			transactionSizeDropdown.Box.Children[0] = widget.NewIcon(icons[assets.ExpandDropdown])
			canvas.Refresh(transactionSizeDropdown.Box.Children[0])
			paintedForm.Hide()
		}
	})

	costAndBalanceAfterSendBox := widget.NewVBox()
	totalCostLabel := widget.NewLabelWithStyle("- DCR", fyne.TextAlignLeading, fyne.TextStyle{})
	balanceAfterSendLabel := widget.NewLabelWithStyle("- DCR", fyne.TextAlignLeading, fyne.TextStyle{})
	costAndBalanceAfterSendBox.Append(widget.NewHBox(widget.NewLabel("Total cost"), layout.NewSpacer(), totalCostLabel))
	costAndBalanceAfterSendBox.Append(widget.NewHBox(widget.NewLabel("Balance after send"), layout.NewSpacer(), balanceAfterSendLabel))
	costAndBalanceAfterSendContainer := fyne.NewContainerWithLayout(layout.NewFixedGridLayout(fyne.NewSize(328, 48)), costAndBalanceAfterSendBox)

	amountErrorLabel := canvas.NewText("", color.RGBA{237, 109, 71, 255})
	amountErrorLabel.TextSize = 12
	amountErrorLabel.Hide()

	amountEntryGroup := widget.NewGroup("Amount", fyne.NewContainerWithLayout(layout.NewFixedGridLayout(
		fyne.NewSize(widget.NewLabel("12345678.12345678").MinSize().Width, amountEntry.MinSize().Height)), amountEntry), amountErrorLabel, widgets.NewVSpacer(4), widget.NewHBox(transactionFeeLabel, transactionSizeDropdown), paintedForm)

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
				value = value[:amountEntry.CursorColumn-1] + value[amountEntry.CursorColumn:]
				//todo: using setText, cursor column count doesnt increase or reduce. Create an issue on this
				amountEntry.CursorColumn--
				amountEntry.SetText(value)
			}
			return
		}

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
			transactionFeeLabel.SetText(fmt.Sprintf("Transaction fee                - DCR"))
			totalCostLabel.SetText("- DCR")
			balanceAfterSendLabel.SetText("- DCR")
			transactionSize.SetText("0 bytes")
		}

		if destinationAddressEntry.Hidden {
			fmt.Println("sending account for sending dropdown is", sendPage.sendingSelectedAccountLabel.Text)
			transactionAuthor.UpdateSendDestination(0, sendPage.sendingSelectedAccountLabel.Text, dcrlibwallet.AmountAtom(amountInFloat), false)
		} else {
			if destinationAddressEntry.Text == "" {
				accountNumber, err = dcrlw.AccountNumber(sendPage.receivingSelectedAccountLabel.Text)
				if err != nil {
					log.Println("Could not get account", sendPage.receivingSelectedAccountLabel.Text)
					return
				}
				temporaryAddress, err := dcrlw.CurrentAddress(int32(accountNumber))
				if err != nil {
					log.Println("Could not get account", sendPage.receivingSelectedAccountLabel.Text)
					return
				}
				transactionAuthor.UpdateSendDestination(0, temporaryAddress, dcrlibwallet.AmountAtom(amountInFloat), false)
			} else {
				if destinationAddressErrorLabel.Hidden {
					transactionAuthor.UpdateSendDestination(0, destinationAddressEntry.Text, dcrlibwallet.AmountAtom(amountInFloat), false)
				} else {
					// return is address in entry is incorrect
					return
				}
			}
		}

		// error was not handled as users might decide not to input address
		feeAndSize, err := transactionAuthor.EstimateFeeAndSize()
		if err != nil {
			if err.Error() == "insufficient_balance" {
				amountErrorLabel.Text = "Insufficient balance"
				amountErrorLabel.Show()
				canvas.Refresh(amountErrorLabel)
			} else {
				log.Println(fmt.Sprintf("could not retrieve transaction fee and size %s", err.Error()))
			}
			return
		}

		if !amountErrorLabel.Hidden {
			amountErrorLabel.Hide()
			canvas.Refresh(amountErrorLabel)
		}

		transactionFeeLabel.SetText(fmt.Sprintf("Transaction fee                %f DCR", feeAndSize.Fee.DcrValue))
		totalCost := feeAndSize.Fee.DcrValue + amountInFloat
		totalCostLabel.SetText(fmt.Sprintf("%f DCR", totalCost))
		balanceAfterSendLabel.SetText(fmt.Sprintf("%f DCR", amountInAccount-totalCost))
		transactionSize.SetText(fmt.Sprintf("%d bytes", feeAndSize.EstimatedSignedSize))

		// transactionAuthor.AddSendDestination("TsfDLrRkk9ciUuwfp2b8PawwnukYD7yAjGd", dcrlibwallet.AmountAtom(10), false)
		// amnt, err := transactionAuthor.EstimateMaxSendAmount()
		// fmt.Println(amnt.DcrValue, err)
		// fee, _ := transactionAuthor.EstimateFeeAndSize()
		// fmt.Println(fee.Fee.DcrValue, fee.EstimatedSignedSize)
	}

	//amountEntryGroup:=widget.NewGroup("Amount", )

	submit := widget.NewButton("Next", func() {
		fmt.Println(sendPage.receivingSelectedAccountLabel.Text, sendPage.receivingSelectedAccountBalanceLabel.Text)
	})
	submit.Disable()

	sendPageContents := widget.NewVBox(baseWidgets, receivingAccountGroup,
		widgets.NewVSpacer(8),
		widget.NewHBox(sendingToAccountGroup, destinationAddressEntryGroup, sendToAccount),
		widgets.NewVSpacer(8),
		widget.NewHBox(amountEntryGroup, widget.NewVBox(sendPage.spendableLabel)),
		widgets.NewVSpacer(12), costAndBalanceAfterSendContainer, widgets.NewVSpacer(16), submit)

	return widget.NewHBox(widgets.NewHSpacer(10), sendPageContents)
}

func createAccountDropdown(initFunction func(), receiveAccountIcon, collapseIcon fyne.Resource, accounts *dcrlibwallet.Accounts, dropdownContent *widget.Box, selectedAccountLabel *widget.Label, selectedAccountBalanceLabel *widget.Label) (accountClickableBox *widgets.ClickableBox) {
	receivingAccountBox := widget.NewHBox(
		widgets.NewHSpacer(15),
		widget.NewIcon(receiveAccountIcon),
		widgets.NewHSpacer(20),
		selectedAccountLabel,
		widgets.NewHSpacer(30),
		selectedAccountBalanceLabel,
		widgets.NewHSpacer(8),
		widget.NewIcon(collapseIcon),
	)

	receivingAccountSelectionPopup := widget.NewPopUp(dropdownContent, fyne.CurrentApp().Driver().AllWindows()[0].Canvas())
	getAccountInBox(initFunction, dropdownContent, selectedAccountLabel, selectedAccountBalanceLabel,
		accounts, receiveAccountIcon, receivingAccountSelectionPopup)
	receivingAccountSelectionPopup.Hide()

	accountClickableBox = widgets.NewClickableBox(receivingAccountBox, func() {
		receivingAccountSelectionPopup.Move(fyne.CurrentApp().Driver().AbsolutePositionForObject(
			accountClickableBox).Add(fyne.NewPos(0, accountClickableBox.Size().Height)))
		receivingAccountSelectionPopup.Show()
	})

	return
}

func getAccountInBox(initFunction func(), dropdownContent *widget.Box, selectedAccountLabel, selectedAccountBalanceLabel *widget.Label, accounts *dcrlibwallet.Accounts, receiveIcon fyne.Resource, popup *widget.PopUp) {
	for index, account := range accounts.Acc {
		if account.Name == "imported" {
			continue
		}

		spendableLabel := canvas.NewText("Spendable", color.White)
		spendableLabel.TextSize = 10

		accountName := account.Name
		accountNameLabel := widget.NewLabel(accountName)
		accountNameLabel.Alignment = fyne.TextAlignLeading
		accountNameBox := widget.NewVBox(
			accountNameLabel,
			widget.NewHBox(widgets.NewHSpacer(1), spendableLabel),
		)

		spendableAmountLabel := canvas.NewText(dcrutil.Amount(account.Balance.Spendable).String(), color.White)
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
		if index != 0 {
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

		dropdownContent.Append(widgets.NewClickableBox(accountsView, func() {
			for _, children := range dropdownContent.Children {
				if box, ok := children.(*widgets.ClickableBox); !ok {
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
				canvas.Refresh(children)
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
			initFunction()
			popup.Hide()
		}))
	}
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
