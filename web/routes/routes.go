package routes

import (
	"context"
	"fmt"
	"html/template"
	"log"
	"os"

	"github.com/go-chi/chi"
	"github.com/raedahgroup/godcr/app"
	"github.com/raedahgroup/godcr/app/config"
)

// Routes holds data required to process web server routes and display appropriate content on a page
type Routes struct {
	walletMiddleware app.WalletMiddleware
	walletExists     bool
	templates        map[string]*template.Template
	blockchain       *Blockchain
	ctx              context.Context
	cnfg             *config.Config
}

// OpenWalletAndSetupRoutes attempts to open the wallet, prepares page templates and creates route handlers
// returns syncBlockchain function
func OpenWalletAndSetupRoutes(ctx context.Context, walletMiddleware app.WalletMiddleware, router chi.Router, appConfig *config.Config) (func(), error) {
	walletExists, err := walletMiddleware.OpenWalletIfExist(ctx)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Failed to open %s wallet: %s\n", walletMiddleware.NetType(), err.Error())
		return nil, err
	}

	routes := &Routes{
		walletMiddleware: walletMiddleware,
		templates:        map[string]*template.Template{},
		blockchain:       &Blockchain{},
		ctx:              ctx,
<<<<<<< 4d5ac49f674f8cba0232d8a03321cd218a22eb5d
		walletExists:     walletExists,
=======
		walletExists: 	  walletExists,
		cnfg:             appConfig,
>>>>>>> spending unconfirmed funds settings
	}

	routes.loadTemplates()
	routes.loadRoutes(router)

	return routes.syncBlockchain, nil
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

	router.Get("/", routes.overviewPage)
	router.Get("/send", routes.sendPage)
	router.Post("/send", routes.submitSendTxForm)
	router.Get("/receive", routes.receivePage)
	router.Get("/generate-address/{accountNumber}", routes.generateReceiveAddress)
	router.Get("/unspent-outputs/{accountNumber}", routes.getUnspentOutputs)
	router.Get("/random-change-outputs", routes.getRandomChangeOutputs)
	router.Get("/history", routes.historyPage)
	router.Get("/next-history-page", routes.getNextHistoryPage)
	router.Get("/transaction-details/{hash}", routes.transactionDetailsPage)
	router.Get("/staking", routes.stakingPage)
	router.Post("/purchase-tickets", routes.submitPurchaseTicketsForm)
	router.Get("/settings", routes.settingsPage)
	router.Post("/settings/change-password", routes.changeSpendingPassword)
	router.Post("/settings/update-spend-unconfirmed", routes.updateSpendUnconfirmedFundSetting)
}
