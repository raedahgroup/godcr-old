package server

import (
	"log"
	"net/http"

	"github.com/go-chi/chi"
	"github.com/raedahgroup/dcrcli/walletrpcclient"
)

type Server struct {
	walletClient *walletrpcclient.Client
	address      string
}

func New(address string, walletClient *walletrpcclient.Client) *Server {
	return &Server{
		address:      address,
		walletClient: walletClient,
	}
}

func (s *Server) registerHandlers(r *chi.Mux) {
	r.Get("/accounts", s.accounts)
	r.Get("/overview", s.overview)
	r.Get("/version", s.version)
	r.Get("/balance/{accountNumber}", s.balance)
	r.Get("/balance/{accountNumber}/{minConf}", s.balance)
}

func (s *Server) Serve() {
	router := chi.NewRouter()
	s.registerHandlers(router)
	log.Fatal(http.ListenAndServe(s.address, router))
}
