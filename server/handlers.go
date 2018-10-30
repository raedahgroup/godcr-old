package server

import (
	"encoding/json"
	"net/http"

	"github.com/go-chi/chi"
)

func (s *Server) renderJSON(w http.ResponseWriter, req *http.Request, res interface{}) {
	data, err := json.Marshal(res)
	if err != nil {
		w.Write([]byte("A server error occurred"))
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.Write(data)
}

func (s *Server) renderError(w http.ResponseWriter, req *http.Request, err error) {
	res := map[string]string{
		"error": err.Error(),
	}
	s.renderJSON(w, req, res)
}

func (s *Server) accounts(w http.ResponseWriter, req *http.Request) {
	res, err := s.walletClient.RunCommand("accounts", nil)
	if err != nil {
		s.renderError(w, req, err)
		return
	}

	s.renderJSON(w, req, res)
}

func (s *Server) balance(w http.ResponseWriter, req *http.Request) {
	accountNumber := chi.URLParam(req, "accountNumber")
	if accountNumber == "" {
		accountNumber = "0"
	}

	minConf := chi.URLParam(req, "minConf")
	if minConf == "" {
		minConf = "0"
	}

	opts := []string{accountNumber, minConf}
	res, err := s.walletClient.RunCommand("balance", opts)
	if err != nil {
		s.renderError(w, req, err)
		return
	}

	s.renderJSON(w, req, res)
}

func (s *Server) overview(w http.ResponseWriter, req *http.Request) {
	res, err := s.walletClient.RunCommand("overview", nil)
	if err != nil {
		s.renderError(w, req, err)
		return
	}
	s.renderJSON(w, req, res)
}

func (s *Server) version(w http.ResponseWriter, req *http.Request) {
	res, err := s.walletClient.RunCommand("walletversion", nil)
	if err != nil {
		s.renderError(w, req, err)
		return
	}
	s.renderJSON(w, req, res)
}
