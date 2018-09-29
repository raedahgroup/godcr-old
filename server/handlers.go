package server

import (
	"encoding/json"
	"net/http"
)

func (s *Server) renderJSON(res interface{}, w http.ResponseWriter, req *http.Request) {
	data, err := json.Marshal(res)
	if err != nil {
		w.Write([]byte("A server error occurred"))
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.Write(data)
}

func (s *Server) accounts(w http.ResponseWriter, req *http.Request) {
	res, err := s.walletClient.RunCommand("accounts", nil)
	if err != nil {
		r := map[string]string{
			"error": err.Error(),
		}
		s.renderJSON(r, w, req)
		return
	}

	s.renderJSON(res, w, req)
}

func (s *Server) balance(w http.ResponseWriter, req *http.Request) {

}
