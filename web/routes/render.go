package routes

import (
	"encoding/json"
	"log"
	"net/http"
)

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
	routes.render("error.html", data, res)
}

func (routes *Routes) renderNoWalletError(res http.ResponseWriter) {
	data := map[string]interface{}{
		"noWallet": true,
	}
	routes.render("error.html", data, res)
}

func renderJSON(data interface{}, res http.ResponseWriter) {
	d, err := json.Marshal(data)
	if err != nil {
		log.Fatalf("error marshalling data: %s", err.Error())
	}

	res.Header().Set("Content-Type", "application/json")
	res.Write(d)
}
