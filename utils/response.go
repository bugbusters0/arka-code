package utils

import (
	"html/template"
	"log"
	"net/http"
	"time"
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

		"subtract": func(a, b float64) float64 {
			return a - b
		},
		"sub": func(a, b float64) float64 {
			return a - b
		},
		"dateMes": func(fechaStr string) string {
			fecha, err := time.Parse("2006-01-02", fechaStr)
			if err != nil {
				return ""
			}
			meses := []string{
				"", "Enero", "Febrero", "Marzo", "Abril", "Mayo", "Junio",
				"Julio", "Agosto", "Septiembre", "Octubre", "Noviembre", "Diciembre",
			}
			return meses[fecha.Month()]
		},
		"dateYear": func(fechaStr string) string {
			fecha, err := time.Parse("2006-01-02", fechaStr)
			if err != nil {
				return ""
			}
			return fecha.Format("2006")
		},
		"add": func(a, b float64) float64 {
			return a + b
		},
		"mul": func(a, b float64) float64 {
			return a * b
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
