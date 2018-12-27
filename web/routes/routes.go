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
	blockchain *Blockchain
}

// Setup prepares page templates and creates route handlers, returns wallet loader function
func Setup(walletMiddleware app.WalletMiddleware, router chi.Router) func() error {
	routes := &Routes{
		walletMiddleware:walletMiddleware,
		templates: map[string]*template.Template{},
		blockchain: &Blockchain{},
	}

	routes.loadTemplates()
	router.Group(routes.handlers) // use router group for page handlers

	return routes.loadWalletAndSyncBlockchain
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
	router.Use(routes.walletLoaderMiddleware())

	router.Get("/", routes.BalancePage)
	router.Get("/send", routes.SendPage)
	router.Post("/send", routes.SubmitSendTxForm)
	router.Get("/receive", routes.ReceivePage)
	router.Get("/generate-address/{accountNumber}", routes.GenerateReceiveAddress)
	router.Get("/unspent-outputs/{accountNumber}", routes.GetUnspentOutputs)
	router.Get("/history", routes.HistoryPage)
}
