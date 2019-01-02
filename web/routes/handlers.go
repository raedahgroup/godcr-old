package routes

import (
	"encoding/base64"
	"fmt"
	"net/http"
	"strconv"

	"github.com/go-chi/chi"
	"github.com/raedahgroup/dcrlibwallet/txhelper"
	qrcode "github.com/skip2/go-qrcode"
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
	accounts, err := routes.walletMiddleware.AccountsOverview()
	if err != nil {
		routes.renderError(fmt.Sprintf("Error fetching account balance: %s", err.Error()), res)
		return
	}

	data := map[string]interface{}{
		"accounts": accounts,
		"detailed": chi.URLParam(req, "detailed") != "",
	}
	routes.render("balance.html", data, res)
}

func (routes *Routes) sendPage(res http.ResponseWriter, req *http.Request) {
	accounts, err := routes.walletMiddleware.AccountsOverview()
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
	utxos := req.Form["tx"]
	amountStr := req.FormValue("amount")
	selectedAccount := req.FormValue("sourceAccount")
	destAddress := req.FormValue("destinationAddress")
	passphrase := req.FormValue("walletPassphrase")

	amount, err := strconv.ParseFloat(amountStr, 64)
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

	sendDestinations := []txhelper.TransactionDestination{{
		Amount:  amount,
		Address: destAddress,
	}}

	var txHash string
	if len(utxos) > 0 {
		txHash, err = routes.walletMiddleware.SendFromUTXOs(sourceAccount, utxos, sendDestinations, passphrase)
	} else {
		txHash, err = routes.walletMiddleware.SendFromAccount(sourceAccount, sendDestinations, passphrase)
	}

	if err != nil {
		data["error"] = err.Error()
		return
	}

	data["txHash"] = txHash
}

func (routes *Routes) receivePage(res http.ResponseWriter, req *http.Request) {
	accounts, err := routes.walletMiddleware.AccountsOverview()
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

	address, err := routes.walletMiddleware.GenerateReceiveAddress(uint32(accountNumber))
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

	utxos, err := routes.walletMiddleware.UnspentOutputs(uint32(accountNumber), 0)
	if err != nil {
		data["success"] = false
		data["message"] = err.Error()
		return
	}

	data["success"] = true
	data["message"] = utxos
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
