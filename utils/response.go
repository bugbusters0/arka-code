package utils

import (
	"html/template"
	"log"
	"net/http"
)

func RenderTemplate(w http.ResponseWriter, layout string, page string, data interface{}) {
	layoutPath := "views/layouts/" + layout + ".html"
	pagePath := "views/" + page + ".html"

	// Crear funciones personalizadas para templates
	funcMap := template.FuncMap{
		"substr": func(s string, start, length int) string {
			if start < 0 || start >= len(s) {
				return ""
			}
			end := start + length
			if end > len(s) {
				end = len(s)
			}
			return s[start:end]
		},
		"sub": func(a, b float64) float64 {
			return a - b
		},
	}

	tmpl, err := template.New(layout+".html").Funcs(funcMap).ParseFiles(layoutPath, pagePath)
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
