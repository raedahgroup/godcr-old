package routes

import (
	"encoding/json"
	"github.com/raedahgroup/godcr/app/walletcore"
	"github.com/raedahgroup/godcr/web/weblog"
	"log"
	"net/http"
)

func (routes *Routes) renderPage(tplName string, data map[string]interface{}, res http.ResponseWriter) {
	connectionInfo, err := walletcore.WalletConnectionInfo(routes.walletMiddleware, routes.walletMiddleware.NetType())
	if err != nil {
		weblog.LogError(err)
	}
	data["connectionInfo"] = connectionInfo
	routes.render(tplName, data, res)
}

func (routes *Routes) render(tplName string, data interface{}, res http.ResponseWriter) {
	if tpl, ok := routes.templates[tplName]; ok {
		err := tpl.Execute(res, data)
		if err != nil {
			log.Fatalf("error executing template: %s", err.Error())
		}
		return
	}

	log.Fatalf("template %s is not registered", tplName)
}

func (routes *Routes) renderError(errorMessage string, res http.ResponseWriter) {
	data := map[string]interface{}{
		"error": errorMessage,
	}
	routes.renderPage("error.html", data, res)
}

func (routes *Routes) renderNoWalletError(res http.ResponseWriter) {
	data := map[string]interface{}{
		"noWallet": true,
	}
	routes.renderPage("error.html", data, res)
}

func renderJSON(data interface{}, res http.ResponseWriter) {
	d, err := json.Marshal(data)
	if err != nil {
		log.Fatalf("error marshalling data: %s", err.Error())
	}

	res.Header().Set("Content-Type", "application/json")
	res.Write(d)
}
