package routes

import (
	"encoding/json"
	"log"
	"net/http"
)

func (routes *Routes) render(tplName string, data map[string]interface{}, res http.ResponseWriter) {
	// append blockchain status to data
	data["blockchainStatus"] = routes.blockchain.report()

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
	errorTemplate := routes.templates["error.html"]
	err := errorTemplate.Execute(res, errorMessage)
	if err != nil {
		log.Fatalf("error executing template: %s", err.Error())
	}
}

func renderJSON(data interface{}, res http.ResponseWriter) {
	d, err := json.Marshal(data)
	if err != nil {
		log.Fatalf("error marshalling data: %s", err.Error())
	}

	res.Header().Set("Content-Type", "application/json")
	res.Write(d)
}
