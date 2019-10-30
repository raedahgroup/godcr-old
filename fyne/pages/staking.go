package pages

import (
	"context"
	"fmt"
	"strconv"
	"strings"

	"fyne.io/fyne"
	"fyne.io/fyne/layout"
	"fyne.io/fyne/widget"
	"github.com/raedahgroup/dcrlibwallet"
	"github.com/raedahgroup/godcr/fyne/widgets"
)

const requiredConfirmations = 2

func stakingPageUpdates(dcrlw *dcrlibwallet.LibWallet, stakeInfoLabel *widget.Label) {
	stakeInfo, err := dcrlw.StakeInfo() //wallet.StakeInfo(context.Background())
	widget.Refresh(stakeInfoLabel)
	if err != nil {
		stakeInfoLabel.SetText(fmt.Sprintf("Error loading stake info: %s", err.Error()))
	} else {
		stakeInfoText := fmt.Sprintf("unmined %d   immature %d   live %d   voted %d   missed %d   expired %d \n"+
			"revoked %d   unspent %d   allmempooltix %d   poolsize %d   total subsidy %s",
			stakeInfo.OwnMempoolTix, stakeInfo.Immature, stakeInfo.Live, stakeInfo.Voted, stakeInfo.Missed, stakeInfo.Expired,
			stakeInfo.Revoked, stakeInfo.Unspent, stakeInfo.AllMempoolTix, stakeInfo.PoolSize, stakeInfo.TotalSubsidy)
		stakeInfoLabel.SetText(stakeInfoText)
	}
}

func stakingPageContent(dcrlw *dcrlibwallet.LibWallet) fyne.CanvasObject {
	pageTitleLabel := widget.NewLabelWithStyle("Staking", fyne.TextAlignLeading, fyne.TextStyle{Italic: true, Bold: true})

	summarySectionTitleLabel := widget.NewLabelWithStyle("Summary", fyne.TextAlignLeading, fyne.TextStyle{Bold: true})
	stakeInfoLabel := widget.NewLabel("")
	stakingPageUpdates(dcrlw, stakeInfoLabel)

	ticketsSectionTitleLabel := widget.NewLabelWithStyle("Purchase Ticket(s)",
		fyne.TextAlignLeading, fyne.TextStyle{Bold: true})

	purchaseForm := widget.NewForm()

	accountDropdown := accountSelectionWidget(dcrlw)
	purchaseForm.Append("Source account", accountDropdown)

	numberOfTicketsEntry := widget.NewEntry()
	purchaseForm.Append("No. of tickets", numberOfTicketsEntry)

	spendUnconfirmedCheck := widget.NewCheck("", nil)
	purchaseForm.Append("Spend unconfirmed", spendUnconfirmedCheck)

	vspHostEntry := widget.NewEntry()
	purchaseForm.Append("VSP host (leave empty if not using vsp)", vspHostEntry)

	passphraseEntry := widget.NewPasswordEntry()
	purchaseForm.Append("Passphrase", passphraseEntry)

	ticketsLabel := widget.NewLabel("")
	ticketsLabel.Hide()
	ticketsLabel.Hide()
	ticketPurchaseError := func(errorMessage string) {
		ticketsLabel.Show()
		widget.Refresh(ticketsLabel)
		ticketsLabel.SetText(fmt.Sprintf("Error purchasing tickets: %s.", errorMessage))
	}

	submitFormButton := widget.NewButton("Submit", func() {
		sourceAccount, _ := dcrlw.AccountNumber(accountDropdown.Selected) //wallet.AccountNumber(accountDropdown.Selected)
		nTickets, err := strconv.Atoi(numberOfTicketsEntry.Text)
		if err != nil {
			ticketPurchaseError("invalid number of tickets value")
			return
		}

		stakingRequiredConfirmations := requiredConfirmations
		if spendUnconfirmedCheck.Checked {
			stakingRequiredConfirmations = 0
		}

		purchaseTicketsRequest := &dcrlibwallet.PurchaseTicketsRequest{
			Account:               sourceAccount,
			NumTickets:            uint32(nTickets),
			RequiredConfirmations: uint32(stakingRequiredConfirmations),
			PoolAddress:           vspHostEntry.Text,
			Passphrase:            []byte(passphraseEntry.Text),
		}

		ticketHashes, err := dcrlw.PurchaseTickets(context.Background(), purchaseTicketsRequest) //wallet.PurchaseTicket(context.Background(), purchaseTicketsRequest)
		if err != nil {
			ticketPurchaseError(err.Error())
		} else {
			numberOfTicketsEntry.SetText("")
			vspHostEntry.SetText("")
			passphraseEntry.SetText("")
			ticketsLabel.SetText(fmt.Sprintf("Success: \n%s", strings.Join(ticketHashes, "\n")))

			stakingPageUpdates(dcrlw, stakeInfoLabel)
		}
	})

	output := widget.NewVBox(
		pageTitleLabel,
		widgets.NewVSpacer(10),
		summarySectionTitleLabel,
		stakeInfoLabel,
		widgets.NewVSpacer(10),
		ticketsSectionTitleLabel,
		widgets.NewVSpacer(10),
		purchaseForm,
		widgets.NewVSpacer(10),
		widget.NewHBox(layout.NewSpacer(), submitFormButton, layout.NewSpacer()),
		ticketsLabel,
	)

	return widget.NewHBox(widgets.NewHSpacer(10), output)

}

func accountSelectionWidget(dcrlw *dcrlibwallet.LibWallet) *widget.Select {
	accounts, err := dcrlw.GetAccountsRaw(0) //wallet.AccountsOverview(walletcore.DefaultRequiredConfirmations)
	if err != nil {
		fmt.Println(err.Error())
		return nil
	}

	var options []string
	for _, account := range accounts.Acc {
		options = append(options, account.Name)
	}

	return widget.NewSelect(options, nil)
}
