package routes

import (
	"log"
	"text/template"

	"github.com/go-chi/chi"
	"github.com/raedahgroup/dcrcli/app"
)

// Routes holds data required to process web server routes and display appropriate content on a page
type Routes struct {
	walletMiddleware    app.WalletMiddleware
	templates map[string]*template.Template
}

// Setup prepares page templates and creates route handlers
func Setup(walletMiddleware app.WalletMiddleware, router chi.Router) {
	routes := &Routes{
		walletMiddleware:walletMiddleware,
		templates: map[string]*template.Template{},
	}

	routes.loadTemplates()
	router.Group(routes.handlers) // use router group for page handlers
}

func (routes *Routes) loadTemplates() {
	layout := "web/views/layout.html"
	utils := "web/views/utils.html"

	for _, tmpl := range templates() {
		parsedTemplate, err := template.New(tmpl.name).ParseFiles(tmpl.path, layout, utils)
		if err != nil {
			log.Fatalf("error loading templates: %s", err.Error())
		}

		routes.templates[tmpl.name] = parsedTemplate
	}
}

func (routes *Routes) handlers(router chi.Router) {
	// this middleware checks if wallet is loaded before executing handlers for following routes
	// if wallet is not loaded, it tries to load it, if that fails, it shows an error page instead
	router.Use(routes.makeWalletLoaderMiddleware())

	router.Get("/", routes.GetBalance)
	router.Get("/send", routes.GetSend)
	router.Post("/send", routes.PostSend)
	router.Get("/receive", routes.GetReceive)
	router.Get("/receive/generate/{accountNumber}", routes.GetReceiveGenerate)
	router.Get("/outputs/unspent/{accountNumber}", routes.GetUnspentOutputs)
	router.Get("/history", routes.GetHistory)
}
