package routes

import (
	"context"
	"html/template"
	"log"

	"github.com/go-chi/chi"
	"github.com/raedahgroup/dcrlibwallet/defaultsynclistener"
	"github.com/raedahgroup/godcr/app"
	"github.com/raedahgroup/godcr/app/config"
)

// Routes holds data required to process web server routes and display appropriate content on a page
type Routes struct {
	walletMiddleware app.WalletMiddleware
	//walletExists       bool
	templates          map[string]*template.Template
	syncProgressReport *defaultsynclistener.ProgressReport
	ctx                context.Context
	settings           *config.Settings
}

// OpenWalletAndSetupRoutes attempts to open the wallet, prepares page templates and creates route handlers
// returns syncBlockChain function
func OpenWalletAndSetupRoutes(ctx context.Context, walletMiddleware app.WalletMiddleware, router chi.Router, settings *config.Settings) (func(), error) {
	// todo: main.go now requires that the user select a wallet or create one before launching interfaces, so need for this check
	//walletExists, err := walletMiddleware.WalletExists()
	//if err != nil {
	//	return nil, err
	//}
	routes := &Routes{
		walletMiddleware:   walletMiddleware,
		templates:          map[string]*template.Template{},
		syncProgressReport: defaultsynclistener.InitProgressReport(),
		ctx:                ctx,
		//walletExists:       walletExists,
		settings: settings,
	}

	routes.loadTemplates()
	routes.loadRoutes(router)

	return routes.syncBlockChain, nil
}

func (routes *Routes) loadTemplates() {
	layout := "../../web/views/layout.html"
	utils := "../../web/views/utils.html"

	for _, tmpl := range templates() {
		parsedTemplate, err := template.New(tmpl.name).Funcs(templateFuncMap()).ParseFiles(tmpl.path, layout, utils)
		if err != nil {
			log.Fatalf("error loading templates: %s", err.Error())
		}

		routes.templates[tmpl.name] = parsedTemplate
	}
}

func (routes *Routes) loadRoutes(router chi.Router) {
	// todo: main.go now requires that the user select a wallet or create one before launching interfaces, so need for this code
	//router.Get("/createwallet", routes.createWalletPage)
	//router.Post("/createwallet", routes.createWallet)
	router.Get("/settings", routes.settingsPage)
	router.Post("/change-password", routes.changeSpendingPassword)
	router.Put("/settings", routes.updateSetting)
	router.Post("/rescan-blockchain", routes.rescanBlockchain)
	router.Delete("/delete-wallet", routes.deleteWallet)

	router.Get("/ws", routes.wsHandler)
	go routes.waitToSendMessagesToClients()

	// use router group for routes that require wallet to be loaded before being accessed
	router.Group(routes.registerRoutesRequiringWallet)
}

func (routes *Routes) registerRoutesRequiringWallet(router chi.Router) {
	// this middleware checks if wallet is loaded before executing handlers for following routes
	// if wallet is not loaded, it tries to load it, if that fails, it shows an error page instead
	router.Use(routes.walletLoaderMiddleware())

	router.Get("/", routes.overviewPage)
	router.Get("/send", routes.sendPage)
	router.Post("/validate-address", routes.validateAddress)
	router.Get("/tx-fee-and-size", routes.getFeeAndSize)
	router.Post("/send", routes.submitSendTxForm)
	router.Get("/max-send-amount", routes.maxSendAmount)
	router.Get("/receive", routes.receivePage)
	router.Get("/generate-address/{accountNumber}", routes.generateReceiveAddress)
	router.Get("/unspent-outputs/{accountNumber}", routes.getUnspentOutputs)
	router.Get("/random-change-outputs", routes.getRandomChangeOutputs)
	router.Get("/history", routes.historyPage)
	router.Get("/next-history-page", routes.getNextHistoryPage)
	router.Get("/transaction-details/{hash}", routes.transactionDetailsPage)
	router.Get("/staking", routes.stakingPage)
	router.Post("/purchase-tickets", routes.submitPurchaseTicketsForm)
	router.Get("/accounts", routes.accountsPage)
	router.Get("/security", routes.securityPage)
}
