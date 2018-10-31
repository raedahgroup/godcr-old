package web

import (
	"fmt"
	"net/http"
	"strconv"

	"github.com/go-chi/chi"
	qrcode "github.com/skip2/go-qrcode"
)

func (s *Server) GetBalance(res http.ResponseWriter, req *http.Request) {
	data := map[string]interface{}{}

	result, err := s.walletClient.Balance()
	if err != nil {
		data["error"] = err
	} else {
		data["result"] = result
	}
	s.render("balance.html", data, res)
}

func (s *Server) GetSend(res http.ResponseWriter, req *http.Request) {
	data := map[string]interface{}{}

	balances, err := s.walletClient.Balance()
	if err != nil {
		data["error"] = err
	} else {
		data["balances"] = balances
	}
	s.render("send.html", data, res)
}

func (s *Server) PostSend(res http.ResponseWriter, req *http.Request) {
	data := map[string]interface{}{}
	defer s.render("send.html", data, res)

	amountStr := req.FormValue("amount")
	sourceAccountStr := req.FormValue("sourceAccount")
	destinationAddressStr := req.FormValue("destinationAddress")
	passphraseStr := req.FormValue("walletPassphrase")

	amount, err := strconv.ParseInt(amountStr, 10, 64)
	if err != nil {
		data["error"] = err
		return
	}

	sourceAccount, err := strconv.ParseUint(sourceAccountStr, 10, 32)
	if err != nil {
		data["error"] = err
		return
	}

	result, err := s.walletClient.Send(amount, uint32(sourceAccount), destinationAddressStr, passphraseStr)
	if err != nil {
		data["error"] = err
		return
	}

	data["success"] = result.TransactionHash
}

func (s *Server) GetReceive(res http.ResponseWriter, req *http.Request) {
	data := map[string]interface{}{}

	accounts, err := s.walletClient.Balance()
	if err != nil {
		data["error"] = err
	} else {
		data["accounts"] = accounts
	}
	s.render("receive.html", data, res)
}

func (s *Server) GetReceiveGenerate(res http.ResponseWriter, req *http.Request) {
	data := map[string]interface{}{}
	defer renderJSON(data, res)

	accountNumberStr := chi.URLParam(req, "account")
	accountNumber, err := strconv.ParseUint(accountNumberStr, 10, 32)
	if err != nil {
		data["success"] = false
		data["message"] = err.Error()
		return
	}

	addr, err := s.walletClient.Receive(uint32(accountNumber))
	if err != nil {
		data["success"] = false
		data["message"] = err.Error()
		return
	}

	png, err := qrcode.Encode(addr.Address, qrcode.Medium, 256)
	if err != nil {
		data["success"] = false
		data["message"] = err.Error()
		return
	}

	imgStr := "data:image/png;base64," + string(png)

	data["success"] = true
	data["address"] = addr.Address
	data["imageStr"] = fmt.Sprintf(`<img src="%s" />`, imgStr)
}
