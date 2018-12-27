package routes

import (
	"encoding/base64"
	"fmt"
	"net/http"
	"strconv"

	"github.com/go-chi/chi"
	qrcode "github.com/skip2/go-qrcode"
)

func (routes *Routes) GetBalance(res http.ResponseWriter, req *http.Request) {
	data := map[string]interface{}{}

	result, err := routes.walletMiddleware.AccountsOverview()
	if err != nil {
		data["error"] = err
	} else {
		data["accounts"] = result
	}

	routes.render("balance.html", data, res)
}

func (routes *Routes) GetSend(res http.ResponseWriter, req *http.Request) {
	data := map[string]interface{}{}

	accounts, err := routes.walletMiddleware.AccountsOverview()
	if err != nil {
		data["error"] = err
	} else {
		data["accounts"] = accounts
	}

	routes.render("send.html", data, res)
}

func (routes *Routes) PostSend(res http.ResponseWriter, req *http.Request) {
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

	var txHash string
	if len(utxos) > 0 {
		txHash, err = routes.walletMiddleware.SendFromUTXOs(utxos, amount, sourceAccount, destAddress, passphrase)
	} else {
		txHash, err = routes.walletMiddleware.SendFromAccount(amount, sourceAccount, destAddress, passphrase)
	}

	if err != nil {
		data["error"] = err.Error()
		return
	}

	data["txHash"] = txHash
}

func (routes *Routes) GetReceive(res http.ResponseWriter, req *http.Request) {
	data := map[string]interface{}{}

	accounts, err := routes.walletMiddleware.AccountsOverview()
	if err != nil {
		data["error"] = err
	} else {
		data["accounts"] = accounts
	}

	routes.render("receive.html", data, res)
}

// GetReceiveGenerate calls walletrpcclient to  generate an address where DCR can be sent to
// this function is called via ajax
func (routes *Routes) GetReceiveGenerate(res http.ResponseWriter, req *http.Request) {
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

func (routes *Routes) GetUnspentOutputs(res http.ResponseWriter, req *http.Request) {
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

func (routes *Routes) GetHistory(res http.ResponseWriter, req *http.Request) {
	data := map[string]interface{}{}

	txns, err := routes.walletMiddleware.TransactionHistory()
	if err != nil {
		data["error"] = err
	} else {
		data["result"] = txns
	}

	routes.render("history.html", data, res)
}
