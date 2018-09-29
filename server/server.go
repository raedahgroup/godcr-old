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

func New(address string, walletClient *walletrpcclient.Client) {
	s := &Server{
		walletClient: walletClient,
	}

	router := chi.NewRouter()
	s.registerHandlers(router)
	log.Fatal(http.ListenAndServe(address, router))
}

func (s *Server) registerHandlers(router *chi.Mux) {

}
