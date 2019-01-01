package routes

import (
	"html/template"
	"log"

	"github.com/go-chi/chi"
	"github.com/raedahgroup/godcr/app"
)

// Routes holds data required to process web server routes and display appropriate content on a page
type Routes struct {
	walletMiddleware app.WalletMiddleware
	templates        map[string]*template.Template
	blockchain       *Blockchain
}

// Setup prepares page templates and creates route handlers, returns syncBlockchain function
func Setup(walletMiddleware app.WalletMiddleware, router chi.Router) func() {
	routes := &Routes{
		walletMiddleware: walletMiddleware,
		templates:        map[string]*template.Template{},
		blockchain:       &Blockchain{},
	}

	routes.loadTemplates()
	routes.loadRoutes(router)

	return routes.syncBlockchain
}

func (routes *Routes) loadTemplates() {
	layout := "web/views/layout.html"
	utils := "web/views/utils.html"

	for _, tmpl := range templates() {
		parsedTemplate, err := template.New(tmpl.name).Funcs(templateFuncMap()).ParseFiles(tmpl.path, layout, utils)
		if err != nil {
			log.Fatalf("error loading templates: %s", err.Error())
		}

		routes.templates[tmpl.name] = parsedTemplate
	}
}

func (routes *Routes) loadRoutes(router chi.Router) {
	router.Get("/createwallet", routes.createWalletPage)
	router.Post("/createwallet", routes.createWallet)

	// use router group for routes that require wallet to be loaded before being accessed
	router.Group(routes.registerRoutesRequiringWallet)
}

func (routes *Routes) registerRoutesRequiringWallet(router chi.Router) {
	// this middleware checks if wallet is loaded before executing handlers for following routes
	// if wallet is not loaded, it tries to load it, if that fails, it shows an error page instead
	router.Use(routes.walletLoaderMiddleware())

	router.Get("/", routes.balancePage)
	router.Get("/send", routes.sendPage)
	router.Post("/send", routes.submitSendTxForm)
	router.Get("/receive", routes.receivePage)
	router.Get("/generate-address/{accountNumber}", routes.generateReceiveAddress)
	router.Get("/unspent-outputs/{accountNumber}", routes.getUnspentOutputs)
	router.Get("/history", routes.historyPage)
}
