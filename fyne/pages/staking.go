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
	godcrApp "github.com/raedahgroup/godcr/app"
	"github.com/raedahgroup/godcr/app/walletcore"
	"github.com/raedahgroup/godcr/fyne/widgets"
)

type stakingPageData struct {
	stakeInfoLabel *widget.Label
}

var staking stakingPageData

func stakingPageReloadData(wallet godcrApp.WalletMiddleware) {
	widget.Refresh(staking.stakeInfoLabel) // necessary to prevent the following SetText() call from throwing an error

	stakeInfo, err := wallet.StakeInfo(context.Background())
	if err != nil {
		staking.stakeInfoLabel.SetText(fmt.Sprintf("Error loading stake info: %s", err.Error()))
	} else {
		stakeInfoText := fmt.Sprintf("unmined %d   immature %d   live %d   voted %d   missed %d   expired %d \n"+
			"revoked %d   unspent %d   allmempooltix %d   poolsize %d   total subsidy %s",
			stakeInfo.OwnMempoolTix, stakeInfo.Immature, stakeInfo.Live, stakeInfo.Voted, stakeInfo.Missed, stakeInfo.Expired,
			stakeInfo.Revoked, stakeInfo.Unspent, stakeInfo.AllMempoolTix, stakeInfo.PoolSize, stakeInfo.TotalSubsidy)
		staking.stakeInfoLabel.SetText(stakeInfoText)
	}
}

func stakingPage(wallet godcrApp.WalletMiddleware) fyne.CanvasObject {
	pageTitleLabel := widget.NewLabelWithStyle("Staking", fyne.TextAlignLeading, fyne.TextStyle{Italic: true, Bold: true})

	summarySectionTitleLabel := widget.NewLabelWithStyle("Summary", fyne.TextAlignLeading, fyne.TextStyle{Bold: true})
	staking.stakeInfoLabel = widget.NewLabel("")
	stakingPageReloadData(wallet)

	ticketsSectionTitleLabel := widget.NewLabelWithStyle("Purchase Ticket(s)",
		fyne.TextAlignLeading, fyne.TextStyle{Bold: true})

	purchaseForm := widget.NewForm()

	accountDropdown := accountSelectionWidget(wallet)
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
	ticketPurchaseError := func(errorMessage string) {
		ticketsLabel.SetText(fmt.Sprintf("Error purchasing tickets: %s.", errorMessage))
		ticketsLabel.Show()
		widget.Refresh(ticketsLabel)
	}

	submitFormButton := widget.NewButton("Submit", func() {
		sourceAccount, _ := wallet.AccountNumber(accountDropdown.Selected)
		nTickets, err := strconv.Atoi(numberOfTicketsEntry.Text)
		if err != nil {
			ticketPurchaseError("invalid number of tickets value")
			return
		}

		requiredConfirmations := walletcore.DefaultRequiredConfirmations
		if spendUnconfirmedCheck.Checked {
			requiredConfirmations = 0
		}

		purchaseTicketsRequest := dcrlibwallet.PurchaseTicketsRequest{
			Account:               sourceAccount,
			NumTickets:            uint32(nTickets),
			RequiredConfirmations: uint32(requiredConfirmations),
			VSPHost:               vspHostEntry.Text,
			Passphrase:            []byte(passphraseEntry.Text),
		}

		ticketHashes, err := wallet.PurchaseTicket(context.Background(), purchaseTicketsRequest)
		if err != nil {
			ticketPurchaseError(err.Error())
		} else {
			numberOfTicketsEntry.SetText("")
			vspHostEntry.SetText("")
			passphraseEntry.SetText("")
			widget.Refresh(purchaseForm)

			ticketsLabel.SetText(fmt.Sprintf("Success: \n%s", strings.Join(ticketHashes, "\n")))
			widget.Refresh(ticketsLabel)
			stakingPageReloadData(wallet)
		}
	})

	output := widget.NewVBox(
		pageTitleLabel,
		widgets.NewVSpacer(10),
		summarySectionTitleLabel,
		staking.stakeInfoLabel,
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

func accountSelectionWidget(wallet godcrApp.WalletMiddleware) *widget.Select {
	accounts, err := wallet.AccountsOverview(walletcore.DefaultRequiredConfirmations)
	if err != nil {
		fmt.Println(err.Error())
		return nil
	}

	var options []string
	for _, account := range accounts {
		options = append(options, account.Name)
	}

	return widget.NewSelect(options, nil)
}
