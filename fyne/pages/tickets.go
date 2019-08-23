package pages

import (
	"context"
	"fmt"

	"fyne.io/fyne"
	"fyne.io/fyne/layout"
	"fyne.io/fyne/widget"

	"github.com/raedahgroup/dcrlibwallet/utils"
	godcrApp "github.com/raedahgroup/godcr/app"
	"github.com/raedahgroup/godcr/fyne/widgets"
)

// ticketPageData contains widgets that needs to be updated realtime
type ticketPageData struct {
	stakeInfoLabel        *widget.Label
	ticketsTable          widgets.TableStruct
	loadTicketsErrorLabel *widget.Label
}

var ticket ticketPageData

func ticketPageUpdates(wallet godcrApp.WalletMiddleware) {
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
		ticket.loadTicketsErrorLabel.SetText(fmt.Sprintf("Error retrieving tickets: %s", err.Error()))
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

	ticket.loadTicketsErrorLabel.Hide()
	ticket.ticketsTable.Container.Show()
}

// initTicketPage does not load any data on init,
// allowing the content refreshing goroutine on menu.go to load this page content when it is navigated to.
// Loading data at this point (aka on fyne launch), causes the tickets table to overlay the overview page.
func initTicketPage() fyne.CanvasObject {
	ticket.stakeInfoLabel = widget.NewLabel("")
	ticket.loadTicketsErrorLabel = widget.NewLabel("")

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
		widget.NewLabelWithStyle("Staking Summary", fyne.TextAlignLeading, fyne.TextStyle{Bold: true}),
		ticket.stakeInfoLabel,
		widgets.NewVSpacer(10),
		ticket.loadTicketsErrorLabel,
		fyne.NewContainerWithLayout(
			// this is a hack, should be able to resize table when dynamic content size changes
			layout.NewFixedGridLayout(fyne.NewSize(700, 500)),
			ticket.ticketsTable.Container),
	)

	return widget.NewHBox(widgets.NewHSpacer(20), output)
}
