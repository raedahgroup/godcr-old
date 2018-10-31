package server

import (
	"fmt"
	"net/http"
	"os"

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

	fmt.Printf("starting http server on %s\n", address)
	err := http.ListenAndServe(address, router)
	if err != nil {
		fmt.Println("Error starting web server")
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}

func (s *Server) registerHandlers(r *chi.Mux) {
	r.Get("/accounts", s.accounts)
	r.Get("/balance/{accountNumber}", s.balance)
	r.Get("/balance/{accountNumber}/{minConf}", s.balance)
}
