package pagehandlers

//
//import (
//	"context"
//	"errors"
//	"fmt"
//	"strconv"
//
//	"github.com/aarzilli/nucular"
//	"github.com/raedahgroup/dcrlibwallet"
//	"github.com/raedahgroup/godcr/app/walletcore"
//	"github.com/raedahgroup/godcr/nuklear/styles"
//	"github.com/raedahgroup/godcr/nuklear/widgets"
//)
//
//type StakingHandler struct {
//	isRendering bool
//
//	stakeInfoFetchError error
//	stakeInfo           *walletcore.StakeInfo
//
//	fetchAccountsError   error
//	purchaseTicketsError error
//
//	accountNumbers       []uint32
//	accountOverviews     []string
//	selectedAccountIndex int
//
//	numTicketsInput       nucular.TextEditor
//	numTicketsInputErrStr string
//	spendUnconfirmed      bool
//
//	isPurchasingTickets    bool
//	purchasedTicketsHashes []string
//}
//
//func (handler *StakingHandler) BeforeRender() {
//	handler.stakeInfoFetchError = nil
//	handler.stakeInfo = nil
//
//	handler.fetchAccountsError = nil
//	handler.purchaseTicketsError = nil
//	handler.numTicketsInputErrStr = ""
//
//	handler.accountNumbers = nil
//	handler.accountOverviews = nil
//
//	handler.resetPurchaseTicketsForm()
//	handler.purchasedTicketsHashes = nil
//
//	handler.isRendering = false
//}
//
//func (handler *StakingHandler) Render(window *nucular.Window, wallet walletcore.Wallet) {
//	if !handler.isRendering {
//		handler.isRendering = true
//		handler.stakeInfo, handler.stakeInfoFetchError = wallet.StakeInfo(context.Background())
//		handler.fetchAccounts(wallet)
//	}
//
//	widgets.PageContentWindowDefaultPadding("Staking", window, func(contentWindow *widgets.Window) {
//		handler.displayStakeInfo(contentWindow)
//		handler.displayPurchaseTicketForm(contentWindow, wallet)
//	})
//}
//
//// fetch accounts to display source account dropdown in purchase ticket section
//func (handler *StakingHandler) fetchAccounts(wallet walletcore.Wallet) {
//	accounts, sendErr := wallet.AccountsOverview(walletcore.DefaultRequiredConfirmations)
//	if sendErr != nil {
//		handler.fetchAccountsError = sendErr
//		return
//	}
//
//	numAccounts := len(accounts)
//	handler.accountNumbers = make([]uint32, numAccounts)
//	handler.accountOverviews = make([]string, numAccounts)
//
//	for index, account := range accounts {
//		handler.accountOverviews[index] = account.String()
//		handler.accountNumbers[index] = account.Number
//	}
//	handler.selectedAccountIndex = 0
//}
//
//func (handler *StakingHandler) displayStakeInfo(contentWindow *widgets.Window) {
//	// display section title with nav font
//	contentWindow.UseFontAndResetToPrevious(styles.NavFont, func() {
//		contentWindow.AddLabel("Stake Info", widgets.LeftCenterAlign)
//	})
//
//	if handler.stakeInfoFetchError != nil {
//		contentWindow.DisplayErrorMessage(handler.stakeInfoFetchError.Error())
//	} else {
//		contentWindow.UseFontAndResetToPrevious(styles.NavFont, func() {
//			contentWindow.Row(styles.LabelHeight).Static(43, 48, 35, 46, 46, 43, 80, 46, 43, 43, 67)
//			contentWindow.Label("Expired", widgets.LeftCenterAlign)
//			contentWindow.Label("Immature", widgets.LeftCenterAlign)
//			contentWindow.Label("Live", widgets.LeftCenterAlign)
//			contentWindow.Label("Revoked", widgets.LeftCenterAlign)
//			contentWindow.Label("Unmined", widgets.LeftCenterAlign)
//			contentWindow.Label("Unspent", widgets.LeftCenterAlign)
//			contentWindow.Label("AllmempoolTix", widgets.LeftCenterAlign)
//			contentWindow.Label("PoolSize", widgets.LeftCenterAlign)
//			contentWindow.Label("Missed", widgets.LeftCenterAlign)
//			contentWindow.Label("Voted", widgets.LeftCenterAlign)
//			contentWindow.Label("Total Subsidy", widgets.LeftCenterAlign)
//		})
//
//		contentWindow.Label(strconv.Itoa(int(handler.stakeInfo.Expired)), widgets.LeftCenterAlign)
//		contentWindow.Label(strconv.Itoa(int(handler.stakeInfo.Immature)), widgets.LeftCenterAlign)
//		contentWindow.Label(strconv.Itoa(int(handler.stakeInfo.Live)), widgets.LeftCenterAlign)
//		contentWindow.Label(strconv.Itoa(int(handler.stakeInfo.Revoked)), widgets.LeftCenterAlign)
//		contentWindow.Label(strconv.Itoa(int(handler.stakeInfo.OwnMempoolTix)), widgets.LeftCenterAlign)
//		contentWindow.Label(strconv.Itoa(int(handler.stakeInfo.Unspent)), widgets.LeftCenterAlign)
//		contentWindow.Label(strconv.Itoa(int(handler.stakeInfo.AllMempoolTix)), widgets.LeftCenterAlign)
//		contentWindow.Label(strconv.Itoa(int(handler.stakeInfo.PoolSize)), widgets.LeftCenterAlign)
//		contentWindow.Label(strconv.Itoa(int(handler.stakeInfo.Missed)), widgets.LeftCenterAlign)
//		contentWindow.Label(strconv.Itoa(int(handler.stakeInfo.Voted)), widgets.LeftCenterAlign)
//		contentWindow.Label(handler.stakeInfo.TotalSubsidy, widgets.LeftCenterAlign)
//	}
//}
//
//func (handler *StakingHandler) displayPurchaseTicketForm(contentWindow *widgets.Window, wallet walletcore.Wallet) {
//	// display form title with nav font
//	contentWindow.UseFontAndResetToPrevious(styles.NavFont, func() {
//		contentWindow.AddLabel("Purchase Ticket", widgets.LeftCenterAlign)
//	})
//
//	// if error fetching accounts, no point displaying the form
//	if handler.fetchAccountsError != nil {
//		contentWindow.DisplayErrorMessage(handler.fetchAccountsError.Error())
//		return
//	}
//
//	// display purchase form proper
//	contentWindow.Row(styles.LabelHeight).Static(styles.TextEditorWidth)
//	contentWindow.Label("Source Account", widgets.LeftCenterAlign)
//
//	contentWindow.Row(styles.TextEditorHeight).Static(styles.AccountSelectorWidth)
//	handler.selectedAccountIndex = contentWindow.ComboSimple(handler.accountOverviews, handler.selectedAccountIndex, 25)
//
//	contentWindow.Row(styles.LabelHeight).Static(styles.TextEditorWidth)
//	contentWindow.Label("Number of tickets", widgets.LeftCenterAlign)
//
//	contentWindow.Row(styles.TextEditorHeight).Static(styles.TextEditorWidth)
//	handler.numTicketsInput.Edit(contentWindow.Window)
//
//	if handler.numTicketsInputErrStr != "" {
//		contentWindow.Row(styles.LabelHeight).Static(styles.TextEditorWidth)
//		contentWindow.LabelColored(handler.numTicketsInputErrStr, widgets.LeftCenterAlign, styles.DecredOrangeColor)
//	}
//
//	contentWindow.Row(styles.CheckboxHeight).Static(styles.TextEditorWidth)
//	contentWindow.CheckboxText("Spend Unconfirmed", &handler.spendUnconfirmed)
//
//	// show tickets hashes after successful tickets purchase, or show error message if purchase failed
//	numTickets := len(handler.purchasedTicketsHashes)
//	if numTickets > 0 {
//		contentWindow.Row(styles.LabelHeight).Dynamic(1)
//		contentWindow.LabelColored(fmt.Sprintf("You have purchased %d ticket(s)", numTickets), widgets.LeftCenterAlign, styles.DecredGreenColor)
//
//		for _, ticketHash := range handler.purchasedTicketsHashes {
//			contentWindow.Row(styles.LabelHeight).Dynamic(1)
//			contentWindow.LabelColored(ticketHash, widgets.LeftCenterAlign, styles.DecredGreenColor)
//		}
//	} else if handler.purchaseTicketsError != nil {
//		contentWindow.Row(styles.LabelHeight).Dynamic(1)
//		contentWindow.LabelColored(handler.purchaseTicketsError.Error(), widgets.LeftCenterAlign, styles.DecredOrangeColor)
//	}
//
//	submitButtonText := "Purchase"
//	if handler.isPurchasingTickets {
//		submitButtonText = "Purchasing..."
//	} else if numTickets > 0 {
//		submitButtonText = "Purchase again"
//	}
//
//	contentWindow.Row(styles.ButtonHeight).Static(styles.ButtonWidth)
//	if contentWindow.ButtonText(submitButtonText) {
//		handler.validateAndSubmit(contentWindow.Window, wallet)
//	}
//}
//
//func (handler *StakingHandler) validateAndSubmit(window *nucular.Window, wallet walletcore.Wallet) {
//	if handler.isPurchasingTickets {
//		return
//	}
//
//	if string(handler.numTicketsInput.Buffer) == "" {
//		handler.numTicketsInputErrStr = "Please specify the number of tickets to purchase"
//	} else {
//		passphraseChan := make(chan string)
//		widgets.NewPassphraseWidget().Get(window, passphraseChan)
//
//		go func() {
//			passphrase := <-passphraseChan
//			if passphrase != "" {
//				handler.submit(passphrase, window, wallet)
//			}
//		}()
//		return
//	}
//
//	window.Master().Changed()
//}
//
//func (handler *StakingHandler) submit(passphrase string, window *nucular.Window, wallet walletcore.Wallet) {
//	handler.isPurchasingTickets = true
//	handler.purchaseTicketsError = nil
//	handler.numTicketsInputErrStr = ""
//	window.Master().Changed()
//
//	defer func() {
//		handler.isPurchasingTickets = false
//		window.Master().Changed()
//	}()
//
//	numTickets, sendErr := strconv.ParseUint(string(handler.numTicketsInput.Buffer), 10, 32)
//	if sendErr != nil {
//		handler.purchaseTicketsError = sendErr
//		return
//	}
//
//	sourceAccount := handler.accountNumbers[handler.selectedAccountIndex]
//
//	requiredConfirmations := walletcore.DefaultRequiredConfirmations
//	if handler.spendUnconfirmed {
//		requiredConfirmations = 0
//	}
//
//	request := dcrlibwallet.PurchaseTicketsRequest{
//		RequiredConfirmations: uint32(requiredConfirmations),
//		Passphrase:            []byte(passphrase),
//		NumTickets:            uint32(numTickets),
//		Account:               uint32(sourceAccount),
//	}
//
//	ticketHashes, sendErr := wallet.PurchaseTicket(context.Background(), request)
//	if sendErr != nil {
//		handler.purchaseTicketsError = sendErr
//		return
//	}
//
//	if len(ticketHashes) == 0 {
//		handler.purchaseTicketsError = errors.New("no ticket was purchased")
//		return
//	}
//
//	handler.purchasedTicketsHashes = ticketHashes
//	handler.resetPurchaseTicketsForm()
//}
//
//func (handler *StakingHandler) resetPurchaseTicketsForm() {
//	handler.selectedAccountIndex = 0
//
//	handler.numTicketsInput.Flags = nucular.EditClipboard | nucular.EditSimple
//	handler.numTicketsInput.Buffer = []rune{'1'}
//
//	handler.spendUnconfirmed = false
//
//	handler.isPurchasingTickets = false
//}
