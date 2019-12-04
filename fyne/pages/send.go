package pages

import (
	"fmt"
	"image/color"
	"log"
	"sort"
	"strconv"
	"time"

	"fyne.io/fyne"
	"fyne.io/fyne/canvas"
	"fyne.io/fyne/layout"
	"fyne.io/fyne/widget"

	"github.com/decred/dcrd/dcrutil"
	"github.com/raedahgroup/dcrlibwallet"

	"github.com/raedahgroup/godcr/fyne/assets"
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

	successLabelContainer := widgets.NewButton(color.RGBA{65, 190, 83, 255}, "Transaction sent", nil)
	successLabelContainer.SetMinSize(successLabelContainer.MinSize().Add(fyne.NewSize(20, 16)))
	successLabelContainer.Container.Hide()

	openedWalletIDs := multiWallet.OpenedWalletIDsRaw()
	if len(openedWalletIDs) == 0 {
		return widget.NewHBox(widgets.NewHSpacer(10), widget.NewLabelWithStyle("Could not retrieve wallets", fyne.TextAlignCenter, fyne.TextStyle{Bold: true}))
	}
	sort.Ints(openedWalletIDs)

	var selectedWallet = multiWallet.WalletWithID(openedWalletIDs[0])

	selectedWalletAccounts, err := selectedWallet.GetAccountsRaw(dcrlibwallet.DefaultRequiredConfirmations)
	if err != nil {
		log.Println(fmt.Sprintf("Error while getting accounts for wallet %s", err.Error()))
		return widget.NewLabel("Error while getting accounts for wallet")
	}

	sendingSelectedWalletLabel := canvas.NewText(selectedWallet.Name, color.RGBA{137, 151, 165, 255})
	selfSendingSelectedWalletLabel := canvas.NewText(selectedWallet.Name, color.RGBA{137, 151, 165, 255})

	temporaryAddress, err := selectedWallet.CurrentAddress(selectedWalletAccounts.Acc[0].Number)
	if err != nil {
		log.Println("could not retrieve account details", err.Error())
		return widget.NewLabel("could not retrieve account details")
	}

	amountInAccount := dcrlibwallet.AmountCoin(selectedWalletAccounts.Acc[0].TotalBalance)

	transactionAuthor := selectedWallet.NewUnsignedTx(0, dcrlibwallet.DefaultRequiredConfirmations)
	transactionAuthor.AddSendDestination(temporaryAddress, 0, true)

	totalCostLabel := widget.NewLabel("- DCR")
	balanceAfterSendLabel := widget.NewLabel("- DCR")
	transactionFeeLabel := widget.NewLabel("- DCR")
	transactionSizeLabel := widget.NewLabel("0 bytes")

	nextButton := widgets.NewButton(color.RGBA{41, 112, 255, 255}, "Next", nil)
	initDynamicContent(openedWalletIDs, selectedWalletAccounts)

	destinationAddressEntry, destinationAddressEntryErrorLabel := widget.NewEntry(), canvas.NewText("", color.RGBA{237, 109, 71, 255})

	amountEntryContainer, amountEntry, isAmountErrorLabelHidden := sendpagehandler.AmountEntryComponents(sendPage.errorLabel, showErrorLabel, &temporaryAddress, &amountInAccount, transactionAuthor,
		transactionFeeLabel, totalCostLabel, balanceAfterSendLabel, transactionSizeLabel, &destinationAddressEntry.Text, &destinationAddressEntry.Hidden,
		&destinationAddressEntryErrorLabel.Hidden, sendPage.Contents, nextButton, sendPage.spendableLabel, multiWallet)

	// this function is called when the sending wallet account is changed.
	onSendingAccountChange := func() {
		selectedWallet = multiWallet.WalletWithID(sendPage.sendingSelectedWalletID)

		accountNumber, err := selectedWallet.AccountNumber(sendPage.sendingSelectedAccountLabel.Text)
		if err != nil {
			showErrorLabel("Could not get accounts")
			log.Println("could not get accounts on account change, reason:", err.Error())
			return
		}

		itransactionAuthor := selectedWallet.NewUnsignedTx(int32(accountNumber), dcrlibwallet.DefaultRequiredConfirmations)
		itransactionAuthor.AddSendDestination(temporaryAddress, 0, true)
		*transactionAuthor = *itransactionAuthor

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
		transactionFeeLabel.Refresh()
		transactionSizeLabel.Refresh()
		sendPage.Contents.Refresh()
		amountEntry.OnChanged(amountEntry.Text)
	}

	fromAccountSelector := sendpagehandler.FromAccountSelector(onSendingAccountChange, "Sending account",
		icons[assets.ReceiveAccountIcon], icons[assets.CollapseIcon], multiWallet, openedWalletIDs, &sendPage.sendingSelectedWalletID,
		sendPage.sendingAccountBoxes, sendPage.sendingSelectedAccountLabel, sendPage.sendingSelectedAccountBalanceLabel,
		sendingSelectedWalletLabel, sendPage.Contents)

	toAccountSelector := sendpagehandler.SendingDestinationComponents(destinationAddressEntry, destinationAddressEntryErrorLabel, sendPage.Contents.Refresh,
		"Receiving account", icons[assets.ReceiveAccountIcon], icons[assets.CollapseIcon], multiWallet, openedWalletIDs, &sendPage.selfSendingSelectedWalletID,
		sendPage.selfSendingAccountBoxes, sendPage.selfSendingSelectedAccountLabel, sendPage.selfSendingSelectedAccountBalanceLabel,
		selfSendingSelectedWalletLabel, transactionFeeLabel, totalCostLabel, balanceAfterSendLabel, transactionSizeLabel, amountEntry, &isAmountErrorLabelHidden, sendPage.Contents, nextButton)

	transactionInfoContainer := sendpagehandler.TransactionDetails(icons[assets.CollapseDropdown], icons[assets.ExpandDropdown],
		transactionFeeLabel, transactionSizeLabel, totalCostLabel, balanceAfterSendLabel, sendPage.Contents)

	nextButton.Container.OnTapped = func() {
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
		if destinationAddressEntry.Hidden {
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
			sendPage.selfSendingSelectedAccountLabel.Text, destinationAddressEntry.Hidden, transactionAuthor, successLabelContainer, sendPage.Contents)
	}

	nextButton.SetMinSize(nextButton.MinSize().Add(fyne.NewSize(0, 20)))
	nextButton.Disable()

	baseWidgets := sendpagehandler.BaseWidgets(icons[assets.InfoIcon], icons[assets.MoreIcon], amountEntry, destinationAddressEntry, window)

	sendPage.Contents.Append(widgets.NewVSpacer(10))
	sendPage.Contents.Append(baseWidgets)
	sendPage.Contents.Append(widgets.NewVSpacer(10))
	sendPage.Contents.Append(widget.NewHBox(layout.NewSpacer(), successLabelContainer.Container, sendPage.errorLabel.Container, layout.NewSpacer()))
	sendPage.Contents.Append(fromAccountSelector)
	sendPage.Contents.Append(widgets.NewVSpacer(10))
	sendPage.Contents.Append(toAccountSelector)
	sendPage.Contents.Append(widgets.NewVSpacer(10))
	sendPage.Contents.Append(amountEntryContainer)
	sendPage.Contents.Append(widgets.NewVSpacer(12))
	sendPage.Contents.Append(transactionInfoContainer)
	sendPage.Contents.Append(widgets.NewVSpacer(15))
	sendPage.Contents.Append(nextButton.Container)

	return widget.NewHBox(widgets.NewHSpacer(20), sendPage.Contents)
}

func initDynamicContent(openedWalletIDs []int, selectedWalletAccounts *dcrlibwallet.Accounts) {
	sendPage.errorLabel = widgets.NewButton(color.RGBA{237, 109, 71, 255}, "", nil)
	sendPage.errorLabel.Container.Hide()

	sendPage.Contents = widget.NewVBox()

	sendPage.sendingSelectedWalletID = openedWalletIDs[0]
	sendPage.selfSendingSelectedWalletID = openedWalletIDs[0]

	sendPage.selfSendingAccountBoxes = make([]*widget.Box, len(openedWalletIDs))
	sendPage.sendingAccountBoxes = make([]*widget.Box, len(openedWalletIDs))

	sendPage.spendableLabel = canvas.NewText("Spendable: "+dcrutil.Amount(selectedWalletAccounts.Acc[0].Balance.Spendable).String(), color.Black)
	sendPage.spendableLabel.TextSize = 12

	sendPage.sendingSelectedAccountLabel = widget.NewLabel(selectedWalletAccounts.Acc[0].Name)
	sendPage.sendingSelectedAccountBalanceLabel = widget.NewLabel(dcrutil.Amount(selectedWalletAccounts.Acc[0].TotalBalance).String())

	sendPage.selfSendingSelectedAccountLabel = widget.NewLabel(selectedWalletAccounts.Acc[0].Name)
	sendPage.selfSendingSelectedAccountBalanceLabel = widget.NewLabel(dcrutil.Amount(selectedWalletAccounts.Acc[0].TotalBalance).String())
}

func showErrorLabel(value string) {
	sendPage.errorLabel.SetText(value)
	sendPage.errorLabel.SetMinSize(sendPage.errorLabel.MinSize().Add(fyne.NewSize(20, 8)))
	sendPage.errorLabel.Container.Show()
	sendPage.Contents.Refresh()

	time.AfterFunc(time.Second*5, func() {
		sendPage.errorLabel.Container.Hide()
		sendPage.Contents.Refresh()
	})
}
