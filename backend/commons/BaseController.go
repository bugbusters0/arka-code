package commons

import (
	"net/http"
)

type BaseController struct {
	// Campos comunes si los necesitas
}

func (bc *BaseController) RenderTemplate(w http.ResponseWriter, tmpl string, data TemplateData) {
	RenderTemplate(w, tmpl, data)
}

func (bc *BaseController) Redirect(w http.ResponseWriter, r *http.Request, url string) {
	http.Redirect(w, r, url, http.StatusFound)
}

func (bc *BaseController) JsonResponse(w http.ResponseWriter, data interface{}, statusCode int) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)
	// json.NewEncoder(w).Encode(data)
}
