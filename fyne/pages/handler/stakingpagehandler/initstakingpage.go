package stakingpagehandler

import (
	"fyne.io/fyne"

	"fyne.io/fyne/widget"

	"github.com/raedahgroup/dcrlibwallet"
	"github.com/raedahgroup/godcr/fyne/pages/handler/multipagecomponents"
	"github.com/raedahgroup/godcr/fyne/pages/handler/values"
	"github.com/raedahgroup/godcr/fyne/widgets"
)

type StakingPageObjects struct {
	Accounts    multipagecomponents.AccountSelectorStruct
	MultiWallet *dcrlibwallet.MultiWallet

	messageLabel        *widgets.BorderedText
	errorLabel          *widgets.BorderedText
	StakingPageContents *widget.Box

	Window fyne.Window
	icons  map[string]*fyne.StaticResource
}

func (stakingPage *StakingPageObjects) InitStakingPage() error {
	err := stakingPage.initBaseObjects()
	if err != nil {
		return err
	}

	summaryLabel := widget.NewLabelWithStyle(values.SummaryText, fyne.TextAlignLeading, fyne.TextStyle{Bold: true, Italic: true})
	stakingPage.StakingPageContents.Append(summaryLabel)

	stakingPage.summaryWalletList()
	stakingPage.getStakingSummary()

	stakingPage.StakingPageContents.Append(widgets.NewVSpacer(values.SpacerSize20))

	PurchaseTicketLabel := widget.NewLabelWithStyle(values.PurchaseTicketsText, fyne.TextAlignLeading, fyne.TextStyle{Bold: true, Italic: true})
	stakingPage.StakingPageContents.Append(PurchaseTicketLabel)

	stakingPage.purchaseTickets()

	stakingPage.StakingPageContents.Append(widgets.NewVSpacer(values.BottomPadding))

	return nil
}
