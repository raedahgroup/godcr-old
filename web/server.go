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
	"github.com/raedahgroup/dcrcli/app"
)

type Server struct {
	wallet    app.WalletMiddleware
	templates map[string]*template.Template
}

func StartHttpServer(wallet app.WalletMiddleware, address string) {
	server := &Server{
		wallet:    wallet,
		templates: map[string]*template.Template{},
	}
	router := chi.NewRouter()

	// setup static file serving
	workDir, _ := os.Getwd()
	filesDir := filepath.Join(workDir, "web/public")
	makeStaticFileServer(router, "/static", http.Dir(filesDir))

	// setup templated pages
	server.loadTemplates()
	// create route group for page handlers
	router.Group(server.registerHandlers)

	fmt.Printf("starting http server on %s\n", address)
	err := http.ListenAndServe(address, router)
	if err != nil {
		fmt.Println("Error starting web server")
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}

	// todo begin blockchain sync, after which subscribe to receive block updates while the server is running
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

func (s *Server) registerHandlers(r chi.Router) {
	// this middleware checks if wallet is loaded before executing handlers for following routes
	// if wallet is not loaded, it tries to load it, if that fails, it shows an error page instead
	r.Use(s.makeWalletLoaderMiddleware())

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

// walletLoaderFn checks if wallet is not open and attempts to open it
// if an error occurs while attempting to open wallet, an error page is displayed and the actual route handler is not called
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
