package pages

import (
	"context"
	"fmt"
	"strings"

	"fyne.io/fyne"
	"fyne.io/fyne/dialog"
	"fyne.io/fyne/layout"
	"fyne.io/fyne/widget"

	"github.com/raedahgroup/dcrlibwallet"
	"github.com/raedahgroup/dcrlibwallet/utils"
	godcrApp "github.com/raedahgroup/godcr/app"
	dlw "github.com/raedahgroup/godcr/app/walletmediums/dcrlibwallet"
	"github.com/raedahgroup/godcr/fyne/widgets"
)

// ticketPageData contains widgets that needs to be updated realtime
type ticketPageData struct {
	stakeInfoLabel          *widget.Label
	ticketsTable            widgets.TableStruct
	ticketsListMessageLabel *widget.Label
}

var ticket ticketPageData

// initTicketPage does not load any data on init,
// allowing the content refreshing goroutine on menu.go to load this page content when it is navigated to.
// Loading data at this point (aka on fyne launch), causes the tickets table to overlay the overview page.
func initTicketPage() fyne.CanvasObject {
	ticket.stakeInfoLabel = widget.NewLabel("")
	ticket.ticketsListMessageLabel = widget.NewLabel("")

	var txTable widgets.TableStruct
	heading := widget.NewHBox(
		widget.NewLabelWithStyle("Date (UTC)", fyne.TextAlignCenter, fyne.TextStyle{Bold: true}),
		widget.NewLabelWithStyle("Hash", fyne.TextAlignCenter, fyne.TextStyle{Bold: true}),
		widget.NewLabelWithStyle("Status", fyne.TextAlignCenter, fyne.TextStyle{Bold: true}),
		widget.NewLabelWithStyle("Block", fyne.TextAlignCenter, fyne.TextStyle{Bold: true}))
	txTable.NewTable(heading)
	ticket.ticketsTable = txTable

	output := widget.NewVBox(
		widgets.NewVSpacer(20),
		widget.NewLabelWithStyle("Staking Summary", fyne.TextAlignLeading, fyne.TextStyle{Bold: true, Italic:true}),
		ticket.stakeInfoLabel,
		widgets.NewVSpacer(20),
		widget.NewLabelWithStyle("Your Tickets", fyne.TextAlignLeading, fyne.TextStyle{Bold: true, Italic:true}),
		fyne.NewContainerWithLayout(
			// this is a hack, should be able to resize table when dynamic content size changes
			layout.NewFixedGridLayout(fyne.NewSize(700, 200)),
			ticket.ticketsTable.Container),
		widgets.NewVSpacer(20),
		ticket.ticketsListMessageLabel,
	)

	return widget.NewHBox(widgets.NewHSpacer(20), output)
}

func ticketPageUpdates(wallet godcrApp.WalletMiddleware, window fyne.Window) {
	stakeInfo, err := wallet.StakeInfo(context.Background())
	if err != nil {
		ticket.stakeInfoLabel.SetText(fmt.Sprintf("Error loading stake info: %s", err.Error()))
	} else {
		stakeInfoText := fmt.Sprintf(
			"Own Mempool Tickets: %d   Immature Tickets: %d   Unspent Tickets: %d   Live Tickets: %d \n"+
				"Missed Tickets: %d   Expired Tickets: %d   Revoked Tickets: %d \n"+
				"Voted Tickets: %d   Total Rewards: %s",
			stakeInfo.OwnMempoolTix, stakeInfo.Immature, stakeInfo.Unspent, stakeInfo.Live,
			stakeInfo.Missed, stakeInfo.Expired, stakeInfo.Revoked,
			stakeInfo.Voted, stakeInfo.TotalSubsidy)
		ticket.stakeInfoLabel.SetText(stakeInfoText)
	}

	tickets, err := wallet.GetTickets()
	if err != nil {
		ticket.ticketsTable.Container.Hide()
		ticket.ticketsListMessageLabel.SetText(fmt.Sprintf("Error retrieving tickets: %s", err.Error()))
		return
	}

	var txTable widgets.TableStruct
	var hBox []*widget.Box
	for _, t := range tickets {
		trimmedHash := t.Ticket.Hash.String()[:25] + "..."
		hBox = append(hBox, widget.NewHBox(
			widget.NewLabelWithStyle(utils.FormatUTCTime(t.Ticket.Timestamp), fyne.TextAlignCenter, fyne.TextStyle{}),
			widget.NewLabelWithStyle(trimmedHash, fyne.TextAlignLeading, fyne.TextStyle{}),
			widget.NewLabelWithStyle(t.Status, fyne.TextAlignCenter, fyne.TextStyle{}),
			widget.NewLabelWithStyle(fmt.Sprintf("%d", t.BlockHeight), fyne.TextAlignLeading, fyne.TextStyle{}),
		))
	}

	heading := widget.NewHBox(
		widget.NewLabelWithStyle("Date (UTC)", fyne.TextAlignLeading, fyne.TextStyle{Bold: true}),
		widget.NewLabelWithStyle("Hash", fyne.TextAlignCenter, fyne.TextStyle{Bold: true}),
		widget.NewLabelWithStyle("Status", fyne.TextAlignCenter, fyne.TextStyle{Bold: true}),
		widget.NewLabelWithStyle("Block", fyne.TextAlignCenter, fyne.TextStyle{Bold: true}))
	txTable.NewTable(heading, hBox...)

	ticket.ticketsTable.Result.Children = txTable.Result.Children
	widget.Refresh(ticket.ticketsTable.Result)

	ticket.ticketsTable.Container.Show()

	processVSPTickets(tickets, wallet, window)
}

func processVSPTickets(tickets []*dcrlibwallet.TicketInfo, wallet godcrApp.WalletMiddleware, window fyne.Window) {
	allVSPTicketHashes := make([]dcrlibwallet.VSPTicketPurchaseInfoRequest, 0)
	errors := make([]string, 0)

	for _, ticket := range tickets {
		hash := ticket.Ticket.Hash.String()
		txDetails, err := wallet.GetTransaction(hash)
		if err != nil {
			errors = append(errors, fmt.Sprintf("Invalid ticket (%s): %s", hash, err.Error()))
		}

		// If the first stakesubmission output does not belong to this wallet,
		// this is likely a vsp ticket, and need to import redeem script.
		if txDetails.Outputs[0].AccountNumber < 0 {
			var walletOwnedSstxCommitOutAddr string
			for _, txOut := range txDetails.Outputs {
				if txOut.ScriptType == "stakesubmission" {
					walletOwnedSstxCommitOutAddr = txOut.Address
				}
			}

			allVSPTicketHashes = append(allVSPTicketHashes, dcrlibwallet.VSPTicketPurchaseInfoRequest{
				TicketHash: hash,
				TicketOwnerCommitmentAddr: walletOwnedSstxCommitOutAddr,
			})
		}
	}

	if len(errors) > 0 {
		errorMessage := "Errors:\n"+strings.Join(errors, "\n")
		ticket.ticketsListMessageLabel.SetText(errorMessage)
	}

	if len(allVSPTicketHashes) == 0 {
		return
	}

	updateReport := func(newReport string) {
		report := ticket.ticketsListMessageLabel.Text
		if report != "" {
			report = fmt.Sprintf("%s\n\n%s", report, newReport)
		} else {
			report = newReport
		}
		ticket.ticketsListMessageLabel.SetText(report)
	}

	report := "Tickets requiring redeem script:"
	for _, t := range allVSPTicketHashes {
		report += fmt.Sprintf("\n%s", t.TicketHash)
	}
	updateReport(report)

	// first confirm that vsp host info is configured
	libwallet, ok := wallet.(*dlw.DcrWalletLib)
	if !ok {
		updateReport("Redeem script recovery not yet implemented for dcrwallet rpc.")
		return
	}

	var vspHost string
	if err := libwallet.ReadFromSettings(dcrlibwallet.VSPHostSettingsKey, &vspHost); err != nil {
		updateReport(fmt.Sprintf("Error reading vsp host configuration: %s", err.Error()))
		return
	}

	if vspHost == "" {
		updateReport("Set VSP Host information from Settings to recover missing redeem scripts.")
		return
	}

	updateReport(fmt.Sprintf("\nAttempting to download redeem scripts from %s...", vspHost))

	passphrase := requestPassphrase(window)
	if passphrase == "" {
		return
	}

	importErrors := libwallet.ImportRedeemScriptsForTickets(allVSPTicketHashes, vspHost, passphrase)
	var finalReport string
	if len(importErrors) > 0 {
		finalReport = fmt.Sprintf("Redeem scripts recovery completed with %d errors:", len(importErrors))
		for _, err := range importErrors {
			finalReport += "\n"+err.Error()
		}
	} else {
		finalReport = "Redeem scripts recovery completed."
	}
	updateReport(finalReport)
}

func requestPassphrase(window fyne.Window) string {
	passphraseChan := make(chan string)

	passphraseEntry := widget.NewPasswordEntry()
	modalForm := &widget.Form{
		OnSubmit: func() {
			if passphraseEntry.Text != "" {
				passphraseChan <- passphraseEntry.Text
			}
		},
		OnCancel: func() {
			passphraseChan <- ""
		},
	}
	modalForm.Append("Enter passphrase to proceed", passphraseEntry)

	dialog.ShowCustom("Fetch Missing Scripts", "Ignore", modalForm, window)

	return <- passphraseChan
}
