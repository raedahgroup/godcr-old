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
	// acct, _ := dcrlw.AccountNumber("default")
	// txauthor := dcrlw.NewUnsignedTx(int32(acct), 0)
	// txauthor.AddSendDestination("TsfDLrRkk9ciUuwfp2b8PawwnukYD7yAjGd", dcrlibwallet.AmountAtom(10), false)
	// amnt, err := txauthor.EstimateMaxSendAmount()
	// fmt.Println(amnt.DcrValue, err)
	// fee, _ := txauthor.EstimateFeeAndSize()
	// fmt.Println(fee.Fee.DcrValue, fee.EstimatedSignedSize)

	// hash, err := txauthor.Broadcast([]byte("admin"))
	// fmt.Println(hash, err)
	// fmt.Println(chainhash.NewHash(hash))

	icons, err := assets.GetIcons(assets.InfoIcon, assets.MoreIcon, assets.ReceiveAccountIcon, assets.CollapseIcon)
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

	accountNumber, _ := dcrlw.AccountNumber("default")
	transactionAuthor := dcrlw.NewUnsignedTx(int32(accountNumber), dcrlibwallet.DefaultRequiredConfirmations)
	onAccountChange := func() {
		acct, err := dcrlw.AccountNumber(sendPage.receivingSelectedAccountLabel.Text)
		if err != nil {
			sendPage.errorLabel.Text = "Could not get account, " + err.Error()
		}
		transactionAuthor.SetSourceAccount(int32(acct))
		fmt.Println("Changed account")

		sendPage.spendableLabel.Text = "Spendable: " + sendPage.receivingSelectedAccountBalanceLabel.Text
		canvas.Refresh(sendPage.spendableLabel)
	}

	sendPage.receivingSelectedAccountLabel = widget.NewLabel(accounts.Acc[0].Name)
	sendPage.receivingSelectedAccountBalanceLabel = widget.NewLabel(dcrutil.Amount(accounts.Acc[0].TotalBalance).String())
	sendPage.receivingAccountDropdownContent = widget.NewVBox()
	receivingAccountClickableBox := createAccountDropdown(onAccountChange, icons[assets.ReceiveAccountIcon], icons[assets.CollapseIcon], accounts, sendPage.receivingAccountDropdownContent, sendPage.receivingSelectedAccountLabel, sendPage.receivingSelectedAccountBalanceLabel)

	receivingAccountGroup := widget.NewGroup("From", receivingAccountClickableBox)

	sendPage.sendingSelectedAccountLabel = widget.NewLabel(accounts.Acc[0].Name)
	sendPage.sendingSelectedAccountBalanceLabel = widget.NewLabel(dcrutil.Amount(accounts.Acc[0].TotalBalance).String())
	sendPage.sendingAccountDropdownContent = widget.NewVBox()
	sendingToAccountClickableBox := createAccountDropdown(nil, icons[assets.ReceiveAccountIcon], icons[assets.CollapseIcon], accounts, sendPage.sendingAccountDropdownContent, sendPage.sendingSelectedAccountLabel, sendPage.sendingSelectedAccountBalanceLabel)

	// destinationAddress := widget.NewEntry()

	// sendingAccountsDropdown := widgets.NewClickableBox(receivingAccountTab, func() {

	// })

	// var accountDropdown *widgets.ClickableBox
	// accountDropdown = widgets.NewClickableBox(accountTab, func() {
	// 	accountSelectionPopup.Move(fyne.CurrentApp().Driver().AbsolutePositionForObject(
	// 		accountDropdown).Add(fyne.NewPos(0, accountDropdown.Size().Height)))
	// 	accountSelectionPopup.Show()
	// })

	destinationAddressEntry := widget.NewEntry()
	destinationAddressEntry.SetPlaceHolder("Destination Address")

	sendingToAccountGroup := widget.NewGroup("To", sendingToAccountClickableBox) //sendingToAccountClickableBox)
	sendingToAccountGroup.Hide()
	destinationAddressEntryGroup := widget.NewGroup("To", fyne.NewContainerWithLayout(
		layout.NewFixedGridLayout(fyne.NewSize(widget.NewLabel("TsfDLrRkk9ciUuwfp2b8PawwnukYD7yAjGd").MinSize().Width, destinationAddressEntry.MinSize().Height)), destinationAddressEntry))

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
	amountEntryGroup := widget.NewGroup("Amount", fyne.NewContainerWithLayout(layout.NewFixedGridLayout(
		fyne.NewSize(widget.NewLabel("12345678.12345678").MinSize().Width, amountEntry.MinSize().Height)), amountEntry))

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
				//amountEntry.CursorColumn--
				amountEntry.SetText(value)
			}
			return
		}

		if numbers := strings.Split(value, "."); len(numbers) == 2 {
			if len(numbers[1]) > 8 {
				sendPage.errorLabel.Text = "Amount has more than 8 decimal places"
				canvas.Refresh(sendPage.errorLabel)
				return
			}
		}

		amountInFloat, err := strconv.ParseFloat(value, 64)
		if err != nil {
			sendPage.errorLabel.Text = "Could not parse float"
			canvas.Refresh(sendPage.errorLabel)
			return
		}

		if destinationAddressEntry.Hidden {
			transactionAuthor.AddSendDestination(sendPage.sendingSelectedAccountLabel.Text, dcrlibwallet.AmountAtom(amountInFloat), false)
		} else {
			transactionAuthor.AddSendDestination(destinationAddressEntry.Text, dcrlibwallet.AmountAtom(amountInFloat), false)
		}

		// transactionAuthor.AddSendDestination("TsfDLrRkk9ciUuwfp2b8PawwnukYD7yAjGd", dcrlibwallet.AmountAtom(10), false)
		// amnt, err := transactionAuthor.EstimateMaxSendAmount()
		// fmt.Println(amnt.DcrValue, err)
		// fee, _ := transactionAuthor.EstimateFeeAndSize()
		// fmt.Println(fee.Fee.DcrValue, fee.EstimatedSignedSize)
	}

	//amountEntryGroup:=widget.NewGroup("Amount", )

	submit := widget.NewButton("Submit", func() {
		fmt.Println(sendPage.receivingSelectedAccountLabel.Text, sendPage.receivingSelectedAccountBalanceLabel.Text)
	})

	return widget.NewHBox(widgets.NewHSpacer(10), widget.NewVBox(baseWidgets, receivingAccountGroup, widget.NewHBox(sendingToAccountGroup, destinationAddressEntryGroup, sendToAccount),
		widget.NewHBox(amountEntryGroup, widget.NewVBox(sendPage.spendableLabel)), widgets.NewVSpacer(10), submit))
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
