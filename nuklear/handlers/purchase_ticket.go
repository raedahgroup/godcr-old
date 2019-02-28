package handlers

import (
	"context"
	"errors"
	"fmt"
	"image/color"
	"strconv"

	"github.com/aarzilli/nucular"
	"github.com/raedahgroup/dcrlibwallet"
	"github.com/raedahgroup/godcr/app"
	"github.com/raedahgroup/godcr/app/walletcore"
	"github.com/raedahgroup/godcr/nuklear/handlers/widgets"
	"github.com/raedahgroup/godcr/nuklear/helpers"
)

type PurchaseTicketsHandler struct {
	err         error
	isRendering bool

	numTicketsInput       nucular.TextEditor
	numTicketsInputErrStr string
	spendUnconfirmed      bool

	accountNumbers   []uint32
	accountOverviews []string

	selectedAccountIndex int

	isSubmitting bool

	ticketHashes []string
}

func (handler *PurchaseTicketsHandler) BeforeRender() {
	handler.isRendering = false
	handler.ticketHashes = nil
	handler.numTicketsInput.Flags = nucular.EditSimple
	handler.accountNumbers = nil
	handler.accountOverviews = nil
	handler.resetForm()
}

func (handler *PurchaseTicketsHandler) Render(window *nucular.Window, wallet app.WalletMiddleware, changePageFunc func(string)) {
	if !handler.isRendering {
		handler.isRendering = true
		handler.fetchAccounts(wallet)
	}

	// draw page
	if pageWindow := helpers.NewWindow("Purchase Tickets Page", window, 0); pageWindow != nil {
		pageWindow.DrawHeader("Purchase Tickets")

		// content window
		if contentWindow := pageWindow.ContentWindow("Purchase Tickets"); contentWindow != nil {
			numTickets := len(handler.ticketHashes)

			if numTickets > 0 {
				contentWindow.Row(10).Dynamic(1)
				contentWindow.LabelColored(fmt.Sprintf("You have purchased %d ticket(s)", numTickets), "LC", color.RGBA{40, 167, 69, 255})

				for i := range handler.ticketHashes {
					contentWindow.Row(10).Dynamic(1)
					contentWindow.LabelColored(handler.ticketHashes[i], "LC", color.RGBA{40, 167, 69, 255})
				}
			}

			if handler.err != nil {
				contentWindow.Row(30).Dynamic(1)
				contentWindow.LabelColored(handler.err.Error(), "LC", color.RGBA{205, 32, 32, 255})
			}

			contentWindow.Row(10).Dynamic(2)
			contentWindow.Label("Source Account", "LC")

			contentWindow.Row(25).Dynamic(1)
			handler.selectedAccountIndex = contentWindow.ComboSimple(handler.accountOverviews, handler.selectedAccountIndex, 25)

			contentWindow.Row(15).Dynamic(2)
			contentWindow.Label("Number of tickets", "LC")

			contentWindow.Row(25).Dynamic(2)
			handler.numTicketsInput.Edit(contentWindow.Window)

			contentWindow.Row(20).Dynamic(1)
			contentWindow.CheckboxText("Spend Unconfirmed", &handler.spendUnconfirmed)

			contentWindow.Row(25).Dynamic(3)
			submitButtonText := "Submit"
			if handler.isSubmitting {
				submitButtonText = "Submitting..."
			}
			if contentWindow.ButtonText(submitButtonText) {
				handler.validateAndSubmit(window, wallet)
			}

			contentWindow.End()
		}
		pageWindow.End()
	}
}

// fetch accounts for select source account field
func (handler *PurchaseTicketsHandler) fetchAccounts(wallet walletcore.Wallet) {
	accounts, err := wallet.AccountsOverview(walletcore.DefaultRequiredConfirmations)
	if err != nil {
		handler.err = err
		return
	}

	numAccounts := len(accounts)
	handler.accountNumbers = make([]uint32, numAccounts)
	handler.accountOverviews = make([]string, numAccounts)

	for index, account := range accounts {
		handler.accountOverviews[index] = fmt.Sprintf("%s - Total %s (Spendable %s)", account.Name, account.Balance.Total.String(), account.Balance.Spendable.String())
		handler.accountNumbers[index] = account.Number
	}
	handler.selectedAccountIndex = 0
}

func (handler *PurchaseTicketsHandler) validateAndSubmit(window *nucular.Window, wallet walletcore.Wallet) {
	isClean := true

	if string(handler.numTicketsInput.Buffer) == "" {
		handler.numTicketsInputErrStr = "Please specify the number of tickets to purchase"
		isClean = false
	}

	if isClean {
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

func (handler *PurchaseTicketsHandler) submit(passphrase string, window *nucular.Window, wallet walletcore.Wallet) {
	handler.isSubmitting = true
	window.Master().Changed()

	defer func() {
		handler.isSubmitting = false
		window.Master().Changed()
	}()

	numTickets, err := strconv.ParseUint(string(handler.numTicketsInput.Buffer), 10, 32)
	if err != nil {
		handler.err = err
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

	ticketHashes, err := wallet.PurchaseTickets(context.Background(), request)
	if err != nil {
		handler.err = err
		return
	}

	if len(ticketHashes) == 0 {
		handler.err = errors.New("no ticket was purchased")
		return
	}

	handler.ticketHashes = ticketHashes
	handler.resetForm()
}

func (handler *PurchaseTicketsHandler) resetForm() {
	handler.err = nil
	handler.numTicketsInputErrStr = ""
	handler.spendUnconfirmed = false
	handler.selectedAccountIndex = 0
	handler.isSubmitting = false
	handler.numTicketsInput.Buffer = []rune{'1'}
}
