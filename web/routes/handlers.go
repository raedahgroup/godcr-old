package routes

import (
	"encoding/base64"
	"fmt"
	"net/http"
	"strconv"

	"github.com/decred/dcrd/dcrutil"
	"github.com/go-chi/chi"
	"github.com/raedahgroup/dcrlibwallet"
	"github.com/raedahgroup/godcr/app/walletcore"
	"github.com/skip2/go-qrcode"
)

func (routes *Routes) createWalletPage(res http.ResponseWriter, req *http.Request) {
	seed, err := routes.walletMiddleware.GenerateNewWalletSeed()
	if err != nil {
		routes.renderError(fmt.Sprintf("Error generating seed for new wallet: %s", err.Error()), res)
		return
	}

	data := struct{ Seed string }{seed}
	routes.render("createwallet.html", data, res)
}

func (routes *Routes) createWallet(res http.ResponseWriter, req *http.Request) {
	req.ParseForm()
	seed := req.FormValue("seed")
	passhprase := req.FormValue("password")

	err := routes.walletMiddleware.CreateWallet(passhprase, seed)
	if err != nil {
		routes.renderError(fmt.Sprintf("Error creating wallet: %s", err.Error()), res)
		return
	}

	// wallet created successfully, wallet is now open, perform first sync
	routes.syncBlockchain()

	http.Redirect(res, req, "/", 303)
}

func (routes *Routes) balancePage(res http.ResponseWriter, req *http.Request) {
	accounts, err := routes.walletMiddleware.AccountsOverview(walletcore.DefaultRequiredConfirmations)
	if err != nil {
		routes.renderError(fmt.Sprintf("Error fetching account balance: %s", err.Error()), res)
		return
	}

	req.ParseForm()
	showDetails := req.FormValue("detailed") != ""

	data := map[string]interface{}{
		"accounts": accounts,
		"detailed": showDetails,
	}
	routes.render("balance.html", data, res)
}

func (routes *Routes) sendPage(res http.ResponseWriter, req *http.Request) {
	accounts, err := routes.walletMiddleware.AccountsOverview(walletcore.DefaultRequiredConfirmations)
	if err != nil {
		routes.renderError(fmt.Sprintf("Error fetching accounts: %s", err.Error()), res)
		return
	}

	data := map[string]interface{}{
		"accounts": accounts,
	}
	routes.render("send.html", data, res)
}

func (routes *Routes) submitSendTxForm(res http.ResponseWriter, req *http.Request) {
	data := map[string]interface{}{}
	defer renderJSON(data, res)

	req.ParseForm()
	utxos := req.Form["utxo"]
	selectedAccount := req.FormValue("source-account")
	passphrase := req.FormValue("wallet-passphrase")
	spendUnconfirmed := req.FormValue("spend-unconfirmed")
	useCustom := req.FormValue("use-custom")

	destinationAddresses := req.Form["destination-address"]
	destinationAmounts := req.Form["destination-amount"]

	sendDestinations, err := walletcore.BuildTxDestinations(destinationAddresses, destinationAmounts)
	if err != nil {
		data["error"] = err.Error()
		return
	}

	account, err := strconv.ParseUint(selectedAccount, 10, 32)
	if err != nil {
		data["error"] = err.Error()
		return
	}
	sourceAccount := uint32(account)

	var requiredConfirmations int32 = walletcore.DefaultRequiredConfirmations
	if spendUnconfirmed != "" {
		requiredConfirmations = 0
	}

	var txHash string
	if useCustom != "" {
		changeOutputAddreses := req.Form["change-output-address"]
		changeOutputAmounts := req.Form["change-output-amount"]

		changeDestinations, err := walletcore.BuildTxDestinations(changeOutputAddreses, changeOutputAmounts)
		if err != nil {
			data["error"] = err.Error()
			return
		}

		if len(changeDestinations) < 1 {
			// add at-least one change output
			totalSelectedInputAmountDcr := req.FormValue("totalSelectedInputAmountDcr")

			totalInputAmountDcr, err := strconv.ParseFloat(totalSelectedInputAmountDcr, 64)
			if err != nil {
				data["error"] = err.Error()
				return
			}

			totalInputAmount, err := dcrutil.NewAmount(totalInputAmountDcr)
			if err != nil {
				data["error"] = err.Error()
				return
			}

			changeDestinations, err = walletcore.GetChangeDestinationsWithRandomAmounts(routes.walletMiddleware, 1, int64(totalInputAmount),
				sourceAccount, len(utxos), sendDestinations)
			if err != nil {
				data["error"] = err.Error()
				return
			}
		}

		txHash, err = routes.walletMiddleware.SendFromUTXOs(sourceAccount, requiredConfirmations, utxos, sendDestinations, changeDestinations, passphrase)
	} else {
		txHash, err = routes.walletMiddleware.SendFromAccount(sourceAccount, requiredConfirmations, sendDestinations, passphrase)
	}

	if err != nil {
		data["error"] = err.Error()
		return
	}

	data["txHash"] = txHash
}

func (routes *Routes) receivePage(res http.ResponseWriter, req *http.Request) {
	accounts, err := routes.walletMiddleware.AccountsOverview(walletcore.DefaultRequiredConfirmations)
	if err != nil {
		routes.renderError(fmt.Sprintf("Error fetching accounts: %s", err.Error()), res)
		return
	}

	data := map[string]interface{}{
		"accounts": accounts,
	}
	routes.render("receive.html", data, res)
}

func (routes *Routes) generateReceiveAddress(res http.ResponseWriter, req *http.Request) {
	data := map[string]interface{}{}
	defer renderJSON(data, res)

	accountNumberStr := chi.URLParam(req, "accountNumber")
	accountNumber, err := strconv.ParseUint(accountNumberStr, 10, 32)
	if err != nil {
		data["success"] = false
		data["message"] = err.Error()
		return
	}

	address, err := routes.walletMiddleware.ReceiveAddress(uint32(accountNumber))
	if err != nil {
		data["success"] = false
		data["message"] = err.Error()
		return
	}

	png, err := qrcode.Encode(address, qrcode.Medium, 256)
	if err != nil {
		data["success"] = false
		data["message"] = err.Error()
		return
	}

	// encode to base64
	encodedStr := base64.StdEncoding.EncodeToString(png)
	imgStr := "data:image/png;base64," + encodedStr

	data["success"] = true
	data["address"] = address
	data["imageStr"] = fmt.Sprintf(`<img src="%s" />`, imgStr)
}

func (routes *Routes) getUnspentOutputs(res http.ResponseWriter, req *http.Request) {
	data := map[string]interface{}{}
	defer renderJSON(data, res)

	accountNumberStr := chi.URLParam(req, "accountNumber")
	accountNumber, err := strconv.ParseUint(accountNumberStr, 10, 32)
	if err != nil {
		data["success"] = false
		data["message"] = err.Error()
		return
	}

	requiredConfirmations := walletcore.DefaultRequiredConfirmations

	getUnconfirmed := req.URL.Query().Get("getUnconfirmed")
	if getUnconfirmed == "true" {
		requiredConfirmations = 0
	}

	utxos, err := routes.walletMiddleware.UnspentOutputs(uint32(accountNumber), 0, int32(requiredConfirmations))
	if err != nil {
		data["success"] = false
		data["message"] = err.Error()
		return
	}

	data["success"] = true
	data["message"] = utxos
}

func (routes *Routes) getRandomChangeOutputs(res http.ResponseWriter, req *http.Request) {
	data := map[string]interface{}{}
	defer renderJSON(data, res)

	req.ParseForm()
	utxos := req.Form["utxo"]
	totalSelectedInputAmountDcr := req.FormValue("totalSelectedInputAmountDcr")
	selectedAccount := req.FormValue("source-account")
	nChangeOutputsStr := req.FormValue("nChangeOutput")

	account, err := strconv.ParseUint(selectedAccount, 10, 32)
	if err != nil {
		data["error"] = err.Error()
		return
	}
	sourceAccount := uint32(account)

	nChangeOutputs, err := strconv.ParseInt(nChangeOutputsStr, 10, 32)
	if err != nil {
		data["error"] = err.Error()
		return
	}

	destinationAddresses := req.Form["destination-address"]
	destinationAmounts := req.Form["destination-amount"]

	destinations, err := walletcore.BuildTxDestinations(destinationAddresses, destinationAmounts)
	if err != nil {
		data["error"] = err.Error()
		return
	}

	changeOutputAddreses := req.Form["change-output-address"]
	changeOutputAmounts := req.Form["change-output-amount"]

	existingChangeDestinations, err := walletcore.BuildTxDestinations(changeOutputAddreses, changeOutputAmounts)
	if err != nil {
		data["error"] = err.Error()
		return
	}

	destinations = append(destinations, existingChangeDestinations...)

	totalInputAmountDcr, err := strconv.ParseFloat(totalSelectedInputAmountDcr, 64)
	if err != nil {
		data["error"] = err.Error()
		return
	}

	totalInputAmount, err := dcrutil.NewAmount(totalInputAmountDcr)
	if err != nil {
		data["error"] = err.Error()
		return
	}

	changeOutputDestinations, err := walletcore.GetChangeDestinationsWithRandomAmounts(routes.walletMiddleware, int(nChangeOutputs), int64(totalInputAmount), sourceAccount, len(utxos), destinations)
	if err != nil {
		data["error"] = err.Error()
		return
	}
	data["message"] = changeOutputDestinations
}

func (routes *Routes) historyPage(res http.ResponseWriter, req *http.Request) {
	txns, err := routes.walletMiddleware.TransactionHistory()
	if err != nil {
		routes.renderError(fmt.Sprintf("Error fetching history: %s", err.Error()), res)
		return
	}

	data := map[string]interface{}{
		"result": txns,
	}
	routes.render("history.html", data, res)
}

func (routes *Routes) transactionDetailsPage(res http.ResponseWriter, req *http.Request) {
	hash := chi.URLParam(req, "hash")
	tx, err := routes.walletMiddleware.GetTransaction(hash)

	if err != nil {
		routes.renderError(fmt.Sprintf("Error fetching transaction: %s", err.Error()), res)
		return
	}

	data := map[string]interface{}{
		"tx": tx,
	}
	routes.render("transaction_details.html", data, res)
}

func (routes *Routes) stakingPage(res http.ResponseWriter, req *http.Request) {
	stakeInfo, err := routes.walletMiddleware.StakeInfo(routes.ctx)
	if err != nil {
		routes.renderError(fmt.Sprintf("Error fetching stake info: %s", err.Error()), res)
		return
	}

	accounts, err := routes.walletMiddleware.AccountsOverview(walletcore.DefaultRequiredConfirmations)
	if err != nil {
		routes.renderError(fmt.Sprintf("Error fetching accounts: %s", err.Error()), res)
		return
	}

	ticketPrice, err := routes.walletMiddleware.TicketPrice(routes.ctx)
	if err != nil {
		routes.renderError(fmt.Sprintf("Error fetching ticket price: %s", err.Error()), res)
		return
	}

	data := map[string]interface{}{
		"stakeinfo": stakeInfo,
		"accounts":    accounts,
		"ticketPrice": dcrutil.Amount(ticketPrice).ToCoin(),
	}
	routes.render("staking.html", data, res)
}

func (routes *Routes) submitPurchaseTicketsForm(res http.ResponseWriter, req *http.Request) {
	data := map[string]interface{}{}
	defer renderJSON(data, res)

	req.ParseForm()
	walletPassphrase := req.FormValue("wallet-passphrase")
	numTicketsStr := req.FormValue("number-of-tickets")
	sourceAccountStr := req.FormValue("source-account")
	spendUnconfirmed := req.FormValue("spend-unconfirmed")

	numTickets, err := strconv.ParseUint(numTicketsStr, 10, 32)
	if err != nil {
		data["success"] = false
		data["message"] = err.Error()
		return
	}

	sourceAccount, err := strconv.ParseUint(sourceAccountStr, 10, 32)
	if err != nil {
		data["success"] = false
		data["message"] = err.Error()
		return
	}

	requiredConfirmations := walletcore.DefaultRequiredConfirmations
	if spendUnconfirmed != "" {
		requiredConfirmations = 0
	}

	request := dcrlibwallet.PurchaseTicketsRequest{
		RequiredConfirmations: uint32(requiredConfirmations),
		Passphrase:            []byte(walletPassphrase),
		NumTickets:            uint32(numTickets),
		Account:               uint32(sourceAccount),
	}

	ticketHashes, err := routes.walletMiddleware.PurchaseTicket(routes.ctx, request)
	if err != nil {
		data["success"] = false
		data["message"] = err.Error()
		return
	}

	if len(ticketHashes) == 0 {
		data["success"] = false
		data["message"] = "no ticket was purchased"
		return
	}

	data["success"] = true
	data["message"] = ticketHashes
}
