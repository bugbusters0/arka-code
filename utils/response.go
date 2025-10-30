package utils

import (
	"html/template"
	"log"
	"net/http"
)

func RenderTemplate(w http.ResponseWriter, layout string, page string, data interface{}) {
	layoutPath := "views/layouts/" + layout + ".html"
	pagePath := "views/" + page + ".html"

	tmpl, err := template.ParseFiles(layoutPath, pagePath)
	if err != nil {
		log.Printf("Error parseando templates: %v", err)
		http.Error(w, "Error cargando template: "+err.Error(), http.StatusInternalServerError)
		return
	}

	err = tmpl.ExecuteTemplate(w, layout, data)
	if err != nil {
		log.Printf("Error ejecutando template: %v", err)
		http.Error(w, "Error ejecutando template: "+err.Error(), http.StatusInternalServerError)
		return
	}
}
