package web

import (
	"encoding/json"
	"log"
	"net/http"
)

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

func (s *Server) renderError(errorMessage string, res http.ResponseWriter) {
	errorTemplate := s.templates["error.html"]
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
