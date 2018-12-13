package web

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"strings"

	"text/template"

	"github.com/go-chi/chi"
	ws "github.com/raedahgroup/dcrcli/walletsource"
)

type Server struct {
	walletSource ws.WalletSource
	templates    map[string]*template.Template
}

func StartHttpServer(address string, walletSource ws.WalletSource) {
	server := &Server{
		walletSource: walletSource,
		templates:    map[string]*template.Template{},
	}

	// load templates
	server.loadTemplates()

	router := chi.NewRouter()
	workDir, _ := os.Getwd()
	filesDir := filepath.Join(workDir, "public")
	FileServer(router, "/static", http.Dir(filesDir))
	server.registerHandlers(router)

	fmt.Printf("starting http server on %s\n", address)
	err := http.ListenAndServe(address, router)
	if err != nil {
		fmt.Println("Error starting web server")
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}

func (s *Server) loadTemplates() {
	layout := "web/views/layout.html"
	utils := "web/views/utils.html"
	funcMap := templateFuncMap()

	tpls := map[string]string{
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

func (s *Server) render(tplName string, data map[string]interface{}, res http.ResponseWriter) {
	if tpl, ok := s.templates[tplName]; ok {
		err := tpl.Execute(res, data)
		if err != nil {
			log.Fatalf("error executing template: %s", err.Error())
		}
		return
	}

	log.Fatalf("template %s is not registered", tplName)
}

func FileServer(r chi.Router, path string, root http.FileSystem) {
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

func (s *Server) registerHandlers(r *chi.Mux) {
	r.Get("/", s.GetBalance)
	r.Get("/send", s.GetSend)
	r.Post("/send", s.PostSend)
	r.Get("/receive", s.GetReceive)
	r.Get("/receive/generate/{accountNumber}", s.GetReceiveGenerate)
	r.Get("/outputs/unspent/{accountNumber}", s.GetUnspentOutputs)
	r.Get("/history", s.GetHistory)
}

func renderJSON(data interface{}, res http.ResponseWriter) {
	d, err := json.Marshal(data)
	if err != nil {
		log.Fatalf("error marshalling data: %s", err.Error())
	}

	res.Header().Set("Content-Type", "application/json")
	res.Write(d)
}
