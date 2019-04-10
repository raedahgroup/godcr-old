package pagehandlers

import (
	"context"
	"errors"
	"fmt"
	"strconv"

	"github.com/aarzilli/nucular"
	"github.com/raedahgroup/dcrlibwallet"
	"github.com/raedahgroup/godcr/app/walletcore"
	"github.com/raedahgroup/godcr/nuklear/widgets"
	"github.com/raedahgroup/godcr/nuklear/styles"
)

type StakingHandler struct {
	isRendering bool

	stakeInfoFetchError error
	stakeInfo           *walletcore.StakeInfo

	fetchAccountsError   error
	purchaseTicketsError error

	accountNumbers       []uint32
	accountOverviews     []string
	selectedAccountIndex int

	numTicketsInput       nucular.TextEditor
	numTicketsInputErrStr string
	spendUnconfirmed      bool

	isPurchasingTickets    bool
	purchasedTicketsHashes []string
}

func (handler *StakingHandler) BeforeRender() {
	handler.stakeInfoFetchError = nil
	handler.stakeInfo = nil

	handler.fetchAccountsError = nil
	handler.purchaseTicketsError = nil
	handler.numTicketsInputErrStr = ""

	handler.accountNumbers = nil
	handler.accountOverviews = nil

	handler.resetPurchaseTicketsForm()
	handler.purchasedTicketsHashes = nil

	handler.isRendering = false
}

func (handler *StakingHandler) Render(window *nucular.Window, wallet walletcore.Wallet) {
	if !handler.isRendering {
		handler.isRendering = true
		handler.stakeInfo, handler.stakeInfoFetchError = wallet.StakeInfo(context.Background())
		handler.fetchAccounts(wallet)
	}

	widgets.PageContentWindow("Staking", window, func(contentWindow *widgets.Window) {
		handler.displayStakeInfo(contentWindow)
		handler.displayPurchaseTicketForm(contentWindow.Window, wallet)
	})
}

// fetch accounts to display source account dropdown in purchase ticket section
func (handler *StakingHandler) fetchAccounts(wallet walletcore.Wallet) {
	accounts, err := wallet.AccountsOverview(walletcore.DefaultRequiredConfirmations)
	if err != nil {
		handler.fetchAccountsError = err
		return
	}

	numAccounts := len(accounts)
	handler.accountNumbers = make([]uint32, numAccounts)
	handler.accountOverviews = make([]string, numAccounts)

	for index, account := range accounts {
		handler.accountOverviews[index] = account.String()
		handler.accountNumbers[index] = account.Number
	}
	handler.selectedAccountIndex = 0
}

func (handler *StakingHandler) displayStakeInfo(contentWindow *widgets.Window) {
	// display section title with nav font
	styles.SetFont(contentWindow.Window, styles.NavFont)
	contentWindow.Row(styles.LabelHeight).Dynamic(1)
	contentWindow.Label("Stake Info", "LC")

	// reset page font
	styles.SetFont(contentWindow.Window, styles.PageContentFont)

	if handler.stakeInfoFetchError != nil {
		contentWindow.SetErrorMessage(handler.stakeInfoFetchError.Error())
	} else {
		styles.SetFont(contentWindow.Window, styles.NavFont)
		contentWindow.Row(styles.LabelHeight).Static(43, 48, 35, 46, 46, 43, 80, 46, 43, 43, 67)
		contentWindow.Label("Expired", "LC")
		contentWindow.Label("Immature", "LC")
		contentWindow.Label("Live", "LC")
		contentWindow.Label("Revoked", "LC")
		contentWindow.Label("Unmined", "LC")
		contentWindow.Label("Unspent", "LC")
		contentWindow.Label("AllmempoolTix", "LC")
		contentWindow.Label("PoolSize", "LC")
		contentWindow.Label("Missed", "LC")
		contentWindow.Label("Voted", "LC")
		contentWindow.Label("Total Subsidy", "LC")

		styles.SetFont(contentWindow.Window, styles.PageContentFont)
		contentWindow.Label(strconv.Itoa(int(handler.stakeInfo.Expired)), "LC")
		contentWindow.Label(strconv.Itoa(int(handler.stakeInfo.Immature)), "LC")
		contentWindow.Label(strconv.Itoa(int(handler.stakeInfo.Live)), "LC")
		contentWindow.Label(strconv.Itoa(int(handler.stakeInfo.Revoked)), "LC")
		contentWindow.Label(strconv.Itoa(int(handler.stakeInfo.OwnMempoolTix)), "LC")
		contentWindow.Label(strconv.Itoa(int(handler.stakeInfo.Unspent)), "LC")
		contentWindow.Label(strconv.Itoa(int(handler.stakeInfo.AllMempoolTix)), "LC")
		contentWindow.Label(strconv.Itoa(int(handler.stakeInfo.PoolSize)), "LC")
		contentWindow.Label(strconv.Itoa(int(handler.stakeInfo.Missed)), "LC")
		contentWindow.Label(strconv.Itoa(int(handler.stakeInfo.Voted)), "LC")
		contentWindow.Label(handler.stakeInfo.TotalSubsidy, "LC")
	}
}

func (handler *StakingHandler) displayPurchaseTicketForm(contentWindow *nucular.Window, wallet walletcore.Wallet) {
	// display form title with nav font
	styles.SetFont(contentWindow, styles.NavFont)
	contentWindow.Row(styles.LabelHeight).Dynamic(1)
	contentWindow.Label("Purchase Ticket", "LC")

	// reset page font
	styles.SetFont(contentWindow, styles.PageContentFont)

	// if error fetching accounts, no point displaying the form
	if handler.fetchAccountsError != nil {
		contentWindow.Row(styles.LabelHeight).Dynamic(1)
		contentWindow.LabelColored(handler.fetchAccountsError.Error(), "LC", styles.DecredOrangeColor)
		return
	}

	// display purchase form proper
	contentWindow.Row(styles.LabelHeight).Static(styles.TextEditorWidth)
	contentWindow.Label("Source Account", "LC")

	contentWindow.Row(styles.TextEditorHeight).Static(styles.AccountSelectorWidth)
	handler.selectedAccountIndex = contentWindow.ComboSimple(handler.accountOverviews, handler.selectedAccountIndex, 25)

	contentWindow.Row(styles.LabelHeight).Static(styles.TextEditorWidth)
	contentWindow.Label("Number of tickets", "LC")

	contentWindow.Row(styles.TextEditorHeight).Static(styles.TextEditorWidth)
	handler.numTicketsInput.Edit(contentWindow)

	if handler.numTicketsInputErrStr != "" {
		contentWindow.Row(styles.LabelHeight).Static(styles.TextEditorWidth)
		contentWindow.LabelColored(handler.numTicketsInputErrStr, "LC", styles.DecredOrangeColor)
	}

	contentWindow.Row(styles.CheckboxHeight).Static(styles.TextEditorWidth)
	contentWindow.CheckboxText("Spend Unconfirmed", &handler.spendUnconfirmed)

	// show tickets hashes after successful tickets purchase, or show error message if purchase failed
	numTickets := len(handler.purchasedTicketsHashes)
	if numTickets > 0 {
		contentWindow.Row(styles.LabelHeight).Dynamic(1)
		contentWindow.LabelColored(fmt.Sprintf("You have purchased %d ticket(s)", numTickets), "LC", styles.DecredGreenColor)

		for _, ticketHash := range handler.purchasedTicketsHashes {
			contentWindow.Row(styles.LabelHeight).Dynamic(1)
			contentWindow.LabelColored(ticketHash, "LC", styles.DecredGreenColor)
		}
	} else if handler.purchaseTicketsError != nil {
		contentWindow.Row(styles.LabelHeight).Dynamic(1)
		contentWindow.LabelColored(handler.purchaseTicketsError.Error(), "LC", styles.DecredOrangeColor)
	}

	submitButtonText := "Purchase"
	if handler.isPurchasingTickets {
		submitButtonText = "Purchasing..."
	} else if numTickets > 0 {
		submitButtonText = "Purchase again"
	}

	contentWindow.Row(styles.ButtonHeight).Static(styles.ButtonWidth)
	if contentWindow.ButtonText(submitButtonText) {
		handler.validateAndSubmit(contentWindow, wallet)
	}
}

func (handler *StakingHandler) validateAndSubmit(window *nucular.Window, wallet walletcore.Wallet) {
	if handler.isPurchasingTickets {
		return
	}

	if string(handler.numTicketsInput.Buffer) == "" {
		handler.numTicketsInputErrStr = "Please specify the number of tickets to purchase"
	} else {
		passphraseChan := make(chan string)
		widgets.NewPassphraseWidget().Get(window, passphraseChan)

		go func() {
			passphrase := <-passphraseChan
			if passphrase != "" {
				handler.submit(passphrase, window, wallet)
			}
		}()
		return
	}

	window.Master().Changed()
}

func (handler *StakingHandler) submit(passphrase string, window *nucular.Window, wallet walletcore.Wallet) {
	handler.isPurchasingTickets = true
	handler.purchaseTicketsError = nil
	handler.numTicketsInputErrStr = ""
	window.Master().Changed()

	defer func() {
		handler.isPurchasingTickets = false
		window.Master().Changed()
	}()

	numTickets, err := strconv.ParseUint(string(handler.numTicketsInput.Buffer), 10, 32)
	if err != nil {
		handler.purchaseTicketsError = err
		return
	}

	sourceAccount := handler.accountNumbers[handler.selectedAccountIndex]

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

	ticketHashes, err := wallet.PurchaseTicket(context.Background(), request)
	if err != nil {
		handler.purchaseTicketsError = err
		return
	}

	if len(ticketHashes) == 0 {
		handler.purchaseTicketsError = errors.New("no ticket was purchased")
		return
	}

	handler.purchasedTicketsHashes = ticketHashes
	handler.resetPurchaseTicketsForm()
}

func (handler *StakingHandler) resetPurchaseTicketsForm() {
	handler.selectedAccountIndex = 0

	handler.numTicketsInput.Flags = nucular.EditClipboard | nucular.EditSimple
	handler.numTicketsInput.Buffer = []rune{'1'}

	handler.spendUnconfirmed = false

	handler.isPurchasingTickets = false
}
