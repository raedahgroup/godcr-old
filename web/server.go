package web

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"strings"
	"text/template"

	"github.com/go-chi/chi"
	"github.com/raedahgroup/dcrcli/core"
)

type Server struct {
	wallet    core.Wallet
	templates map[string]*template.Template
}

func StartHttpServer(address string, wallet core.Wallet) {
	server := &Server{
		wallet:    wallet,
		templates: map[string]*template.Template{},
	}
	router := chi.NewRouter()

	// ensure wallet is loaded before executing following handlers
	router.Use(server.makeWalletLoaderMiddleware())

	// setup static file serving
	workDir, _ := os.Getwd()
	filesDir := filepath.Join(workDir, "web/public")
	makeStaticFileServer(router, "/static", http.Dir(filesDir))

	// setup templated pages
	server.loadTemplates()
	server.registerHandlers(router)

	fmt.Printf("starting http server on %s\n", address)
	err := http.ListenAndServe(address, router)
	if err != nil {
		fmt.Println("Error starting web server")
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}

func makeStaticFileServer(r chi.Router, path string, root http.FileSystem) {
	if strings.ContainsAny(path, "{}*") {
		panic("FileServer does not permit URL parameters.")
	}

	fs := http.StripPrefix(path, http.FileServer(root))

	if path != "/" && path[len(path)-1] != '/' {
		r.Get(path, http.RedirectHandler(path+"/", 301).ServeHTTP)
		path += "/"
	}
	path += "*"

	r.Get(path, http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		fs.ServeHTTP(w, r)
	}))
}

func (s *Server) loadTemplates() {
	layout := "web/views/layout.html"
	utils := "web/views/utils.html"
	funcMap := templateFuncMap()

	tpls := map[string]string{
		"error.html":   "web/views/error.html",
		"balance.html": "web/views/balance.html",
		"send.html":    "web/views/send.html",
		"receive.html": "web/views/receive.html",
		"history.html": "web/views/history.html",
	}

	for i, v := range tpls {
		tpl, err := template.New(i).Funcs(funcMap).ParseFiles(v, layout, utils)
		if err != nil {
			log.Fatalf("error loading templates: %s", err.Error())
		}

		s.templates[i] = tpl
	}
}

func (s *Server) registerHandlers(r *chi.Mux) {
	r.Get("/", s.GetBalance)
	r.Get("/send", s.GetSend)
	r.Post("/send", s.PostSend)
	r.Get("/receive", s.GetReceive)
	r.Get("/receive/generate/{accountNumber}", s.GetReceiveGenerate)
	r.Get("/outputs/unspent/{accountNumber}", s.GetUnspentOutputs)
	r.Get("/history", s.GetHistory)
}

func (s *Server) makeWalletLoaderMiddleware() func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return s.walletLoaderFn(next)
	}
}

func (s *Server) walletLoaderFn(next http.Handler) http.Handler {
	return http.HandlerFunc(func(res http.ResponseWriter, req *http.Request) {
		if !s.wallet.IsWalletOpen() {
			err := s.loadWallet()
			if err != nil {
				s.renderError(err.Error(), res)
				return
			}
		}

		next.ServeHTTP(res, req)
	})
}

func (s *Server) loadWallet() error {
	walletExists, err := s.wallet.WalletExists()
	if err != nil {
		return fmt.Errorf("Error checking for wallet: %s", err.Error())
	}

	if !walletExists {
		return fmt.Errorf("Wallet not created. Please create a wallet to continue. Use `dcrcli create` on terminal")
	}

	err = s.wallet.OpenWallet()
	if err != nil {
		return fmt.Errorf("Failed to open wallet: %s", err.Error())
	}

	return nil
}
