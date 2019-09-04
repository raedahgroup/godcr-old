package pagehandlers

import (
	"context"
	"errors"
	"fmt"
	"strconv"

	"github.com/aarzilli/nucular"
	"github.com/raedahgroup/dcrlibwallet"
	"github.com/raedahgroup/godcr/app/config"
	"github.com/raedahgroup/godcr/app/walletcore"
	"github.com/raedahgroup/godcr/nuklear/styles"
	"github.com/raedahgroup/godcr/nuklear/widgets"
)

const numTicketsInputWidth = 50

type StakingHandler struct {
	wallet walletcore.Wallet

	stakeInfoFetchError error
	stakeInfo           *walletcore.StakeInfo

	spendUnconfirmed      bool
	accountSelector       *widgets.AccountSelector
	numTicketsInput       *nucular.TextEditor
	numTicketsInputErrStr string

	isPurchasingTickets    bool
	purchasedTicketsHashes []string
	purchaseTicketsError   error
}

func (handler *StakingHandler) BeforeRender(wallet walletcore.Wallet, settings *config.Settings, refreshWindowDisplay func()) bool {
	handler.wallet = wallet

	handler.stakeInfoFetchError = nil
	handler.stakeInfo = nil

	// fetch stake info data in background as it could take long for wallets with much txs
	go func() {
		handler.stakeInfo, handler.stakeInfoFetchError = wallet.StakeInfo(context.Background())
		refreshWindowDisplay()
	}()

	handler.spendUnconfirmed = false // todo should use the value in settings
	handler.accountSelector = widgets.AccountSelectorWidget("From:", handler.spendUnconfirmed, true, wallet)
	handler.numTicketsInput = &nucular.TextEditor{}
	handler.numTicketsInput.Flags = nucular.EditClipboard | nucular.EditSimple

	handler.isPurchasingTickets = false
	handler.purchasedTicketsHashes = nil
	handler.purchaseTicketsError = nil

	handler.resetPurchaseTicketsForm()

	return true
}

func (handler *StakingHandler) Render(window *nucular.Window) {
	widgets.PageContentWindowDefaultPadding("Staking", window, func(contentWindow *widgets.Window) {
		handler.displayStakeInfo(contentWindow)
		contentWindow.AddHorizontalSpace(20)
		handler.displayPurchaseTicketForm(contentWindow)
	})
}

func (handler *StakingHandler) displayStakeInfo(contentWindow *widgets.Window) {
	contentWindow.AddLabelWithFont("Stake Info", widgets.LeftCenterAlign, styles.BoldPageContentFont)

	if handler.stakeInfoFetchError != nil {
		contentWindow.DisplayErrorMessage("Error fetching stake info", handler.stakeInfoFetchError)
	} else {
		stakingTable := widgets.NewTable()

		// add table header using nav font
		stakingTable.AddRowWithFont(styles.NavFont,
			widgets.NewLabelTableCell("Expired", widgets.LeftCenterAlign),
			widgets.NewLabelTableCell("Immature", widgets.LeftCenterAlign),
			widgets.NewLabelTableCell("Live", widgets.LeftCenterAlign),
			widgets.NewLabelTableCell("Revoked", widgets.LeftCenterAlign),
			widgets.NewLabelTableCell("Unmined", widgets.LeftCenterAlign),
			widgets.NewLabelTableCell("Unspent", widgets.LeftCenterAlign),
			widgets.NewLabelTableCell("AllmempoolTix", widgets.LeftCenterAlign),
			widgets.NewLabelTableCell("PoolSize", widgets.LeftCenterAlign),
			widgets.NewLabelTableCell("Missed", widgets.LeftCenterAlign),
			widgets.NewLabelTableCell("Voted", widgets.LeftCenterAlign),
			widgets.NewLabelTableCell("Total Subsidy", widgets.LeftCenterAlign),
		)

		// stake info was loaded in background in BeforeRender
		// the data may not have been loaded at this time
		if handler.stakeInfo != nil {
			stakingTable.AddRow(
				widgets.NewLabelTableCell(strconv.Itoa(int(handler.stakeInfo.Expired)), widgets.LeftCenterAlign),
				widgets.NewLabelTableCell(strconv.Itoa(int(handler.stakeInfo.Immature)), widgets.LeftCenterAlign),
				widgets.NewLabelTableCell(strconv.Itoa(int(handler.stakeInfo.Live)), widgets.LeftCenterAlign),
				widgets.NewLabelTableCell(strconv.Itoa(int(handler.stakeInfo.Revoked)), widgets.LeftCenterAlign),
				widgets.NewLabelTableCell(strconv.Itoa(int(handler.stakeInfo.OwnMempoolTix)), widgets.LeftCenterAlign),
				widgets.NewLabelTableCell(strconv.Itoa(int(handler.stakeInfo.Unspent)), widgets.LeftCenterAlign),
				widgets.NewLabelTableCell(strconv.Itoa(int(handler.stakeInfo.AllMempoolTix)), widgets.LeftCenterAlign),
				widgets.NewLabelTableCell(strconv.Itoa(int(handler.stakeInfo.PoolSize)), widgets.LeftCenterAlign),
				widgets.NewLabelTableCell(strconv.Itoa(int(handler.stakeInfo.Missed)), widgets.LeftCenterAlign),
				widgets.NewLabelTableCell(strconv.Itoa(int(handler.stakeInfo.Voted)), widgets.LeftCenterAlign),
				widgets.NewLabelTableCell(handler.stakeInfo.TotalSubsidy, widgets.LeftCenterAlign),
			)
		}

		stakingTable.Render(contentWindow)
	}
}

func (handler *StakingHandler) displayPurchaseTicketForm(contentWindow *widgets.Window) {
	contentWindow.AddLabelWithFont("Purchase Ticket", widgets.LeftCenterAlign, styles.BoldPageContentFont)

	handler.accountSelector.Render(contentWindow)
	contentWindow.AddCheckbox("Spend Unconfirmed", &handler.spendUnconfirmed, func() {
		// reload account balance and refresh display
		handler.accountSelector = widgets.AccountSelectorWidget("From:", handler.spendUnconfirmed,
			true, handler.wallet)
		handler.accountSelector.Render(contentWindow)
		contentWindow.Master().Changed()
	})

	contentWindow.AddHorizontalSpace(10)
	contentWindow.Row(widgets.EditorHeight).Static(contentWindow.LabelWidth("Number of Tickets"), numTicketsInputWidth)
	contentWindow.AddLabelsToCurrentRow(widgets.NewLabelTableCell("Number of Tickets", widgets.LeftCenterAlign))
	contentWindow.AddEditorToCurrentRow(handler.numTicketsInput)
	if handler.numTicketsInputErrStr != "" {
		contentWindow.DisplayMessage(handler.numTicketsInputErrStr, styles.DecredOrangeColor)
	}

	submitButtonText := "Purchase"
	if handler.isPurchasingTickets {
		submitButtonText = "Purchasing..."
	}
	contentWindow.AddHorizontalSpace(10)
	contentWindow.AddButton(submitButtonText, func() {
		handler.validateAndSubmit(contentWindow.Window)
	})

	// show tickets hashes after successful tickets purchase, or show error message if purchase failed
	contentWindow.AddHorizontalSpace(10)
	numTickets := len(handler.purchasedTicketsHashes)
	if numTickets > 0 {
		successMessage := fmt.Sprintf("You have purchased %d ticket(s)", numTickets)
		contentWindow.AddColoredLabel(successMessage, styles.DecredGreenColor, widgets.LeftCenterAlign)
		for _, ticketHash := range handler.purchasedTicketsHashes {
			contentWindow.AddColoredLabel(ticketHash, styles.DecredGreenColor, widgets.LeftCenterAlign)
		}
	} else if handler.purchaseTicketsError != nil {
		contentWindow.DisplayErrorMessage("Error purchasing ticket", handler.purchaseTicketsError)
	}
}

func (handler *StakingHandler) validateAndSubmit(window *nucular.Window) {
	if handler.isPurchasingTickets {
		return
	}

	if string(handler.numTicketsInput.Buffer) == "" {
		handler.numTicketsInputErrStr = "Please specify the number of tickets to purchase"
		window.Master().Changed()
	} else {
		passphraseChan := make(chan string)
		widgets.NewPassphraseWidget().Get(window, passphraseChan)

		go func() {
			passphrase := <-passphraseChan
			if passphrase != "" {
				handler.submit(passphrase, window)
			}
		}()
		return
	}
}

func (handler *StakingHandler) submit(passphrase string, window *nucular.Window) {
	handler.isPurchasingTickets = true
	handler.purchaseTicketsError = nil
	handler.numTicketsInputErrStr = ""
	window.Master().Changed()

	defer func() {
		handler.isPurchasingTickets = false
		window.Master().Changed()
	}()

	numTickets, sendErr := strconv.ParseUint(string(handler.numTicketsInput.Buffer), 10, 32)
	if sendErr != nil {
		handler.purchaseTicketsError = sendErr
		return
	}

	sourceAccount := handler.accountSelector.GetSelectedAccountNumber()

	requiredConfirmations := walletcore.DefaultRequiredConfirmations
	if handler.spendUnconfirmed {
		requiredConfirmations = 0
	}

	request := dcrlibwallet.PurchaseTicketsRequest{
		RequiredConfirmations: uint32(requiredConfirmations),
		Passphrase:            []byte(passphrase),
		NumTickets:            uint32(numTickets),
		Account:               uint32(sourceAccount),
	}

	ticketHashes, sendErr := handler.wallet.PurchaseTicket(context.Background(), request)
	if sendErr != nil {
		handler.purchaseTicketsError = sendErr
		return
	}

	if len(ticketHashes) == 0 {
		handler.purchaseTicketsError = errors.New("no ticket was purchased")
		return
	}

	handler.purchasedTicketsHashes = ticketHashes
	handler.resetPurchaseTicketsForm()
	window.Master().Changed()
}

func (handler *StakingHandler) resetPurchaseTicketsForm() {
	handler.accountSelector.Reset()

	handler.numTicketsInput.Buffer = []rune{'1'}
	handler.numTicketsInputErrStr = ""

	handler.isPurchasingTickets = false
}
