package web

import (
	"encoding/base64"
	"fmt"
	"net/http"
	"strconv"

	"github.com/go-chi/chi"
	qrcode "github.com/skip2/go-qrcode"
)

func (s *Server) GetBalance(res http.ResponseWriter, req *http.Request) {
	data := map[string]interface{}{}

	result, err := s.walletSource.AccountsOverview()
	if err != nil {
		data["error"] = err
	} else {
		data["accounts"] = result
	}

	s.render("balance.html", data, res)
}

func (s *Server) GetSend(res http.ResponseWriter, req *http.Request) {
	data := map[string]interface{}{}

	accounts, err := s.walletSource.AccountsOverview()
	if err != nil {
		data["error"] = err
	} else {
		data["accounts"] = accounts
	}

	s.render("send.html", data, res)
}

func (s *Server) PostSend(res http.ResponseWriter, req *http.Request) {
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
		txHash, err = s.walletSource.SendFromUTXOs(utxos, amount, sourceAccount, destAddress, passphrase)
	} else {
		txHash, err = s.walletSource.SendFromAccount(amount, sourceAccount, destAddress, passphrase)
	}

	if err != nil {
		data["error"] = err.Error()
		return
	}

	data["txHash"] = txHash
}

func (s *Server) GetReceive(res http.ResponseWriter, req *http.Request) {
	data := map[string]interface{}{}

	accounts, err := s.walletSource.AccountsOverview()
	if err != nil {
		data["error"] = err
	} else {
		data["accounts"] = accounts
	}

	s.render("receive.html", data, res)
}

// GetReceiveGenerate calls walletrpcclient to  generate an address where DCR can be sent to
// this function is called via ajax
func (s *Server) GetReceiveGenerate(res http.ResponseWriter, req *http.Request) {
	data := map[string]interface{}{}
	defer renderJSON(data, res)

	accountNumberStr := chi.URLParam(req, "accountNumber")
	accountNumber, err := strconv.ParseUint(accountNumberStr, 10, 32)
	if err != nil {
		data["success"] = false
		data["message"] = err.Error()
		return
	}

	address, err := s.walletSource.GenerateReceiveAddress(uint32(accountNumber))
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

func (s *Server) GetUnspentOutputs(res http.ResponseWriter, req *http.Request) {
	data := map[string]interface{}{}
	defer renderJSON(data, res)

	accountNumberStr := chi.URLParam(req, "accountNumber")
	accountNumber, err := strconv.ParseUint(accountNumberStr, 10, 32)
	if err != nil {
		data["success"] = false
		data["message"] = err.Error()
		return
	}

	utxos, err := s.walletSource.UnspentOutputs(uint32(accountNumber), 0)
	if err != nil {
		data["success"] = false
		data["message"] = err.Error()
		return
	}

	data["success"] = true
	data["message"] = utxos
}

func (s *Server) GetHistory(res http.ResponseWriter, req *http.Request) {
	data := map[string]interface{}{}

	txns, err := s.walletClient.GetTransactions()
	if err != nil {
		data["error"] = err
	} else {
		data["result"] = txns
	}

	s.render("history.html", data, res)
}
