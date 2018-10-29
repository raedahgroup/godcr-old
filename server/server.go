package server

import (
	"log"
	"net/http"

	"github.com/go-chi/chi"
	"github.com/raedahgroup/dcrcli/walletrpcclient"
)

type Server struct {
	walletClient *walletrpcclient.Client
}

func StartHttpServer(address string, walletClient *walletrpcclient.Client) {
	server := &Server{
		walletClient: walletClient,
	}

	router := chi.NewRouter()
	server.registerHandlers(router)
	log.Fatal(http.ListenAndServe(address, router))
}

func (s *Server) registerHandlers(r *chi.Mux) {
	r.Get("/accounts", s.accounts)
	r.Get("/overview", s.overview)
	r.Get("/version", s.version)
	r.Get("/balance/{accountNumber}", s.balance)
	r.Get("/balance/{accountNumber}/{minConf}", s.balance)
}
