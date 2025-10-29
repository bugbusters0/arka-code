package commons

import (
	"html/template"
	"net/http"
	"strconv"

	"github.com/gorilla/sessions"
)

var Store = sessions.NewCookieStore([]byte("arka-code-secret-key-2024"))

// Sesiones
func IsAuthenticated(r *http.Request) bool {
	session, _ := Store.Get(r, "session")
	auth, ok := session.Values["authenticated"].(bool)
	return ok && auth
}

func Login(w http.ResponseWriter, r *http.Request, user interface{}) {
	session, _ := Store.Get(r, "session")
	u := user.(map[string]interface{})

	session.Values["authenticated"] = true
	session.Values["user_id"] = u["id"]
	session.Values["family_id"] = u["family_id"]
	session.Values["user_type"] = u["type"]
	session.Save(r, w)
}

func Logout(w http.ResponseWriter, r *http.Request) {
	session, _ := Store.Get(r, "session")
	session.Values["authenticated"] = false
	session.Values["user_id"] = nil
	session.Values["family_id"] = nil
	session.Values["user_type"] = nil
	session.Values["miembro_id"] = nil
	session.Save(r, w)
}

func GetUserID(r *http.Request) int {
	session, _ := Store.Get(r, "session")
	userID, ok := session.Values["user_id"].(int)
	if !ok {
		return 0
	}
	return userID
}

func GetFamilyID(r *http.Request) int {
	session, _ := Store.Get(r, "session")
	familyID, ok := session.Values["family_id"].(int)
	if !ok {
		return 0
	}
	return familyID
}

func GetUserType(r *http.Request) string {
	session, _ := Store.Get(r, "session")
	userType, ok := session.Values["user_type"].(string)
	if !ok {
		return ""
	}
	return userType
}

func GetMiembroID(r *http.Request) int {
	session, _ := Store.Get(r, "session")
	miembroID, ok := session.Values["miembro_id"].(int)
	if !ok {
		return GetUserID(r) // Fallback al user_id
	}
	return miembroID
}

func SetMiembroID(w http.ResponseWriter, r *http.Request, miembroID int) {
	session, _ := Store.Get(r, "session")
	session.Values["miembro_id"] = miembroID
	session.Save(r, w)
}

// Mensajes flash
func SetError(r *http.Request, message string) {
	session, _ := Store.Get(r, "session")
	session.AddFlash(message, "errors")
	session.Save(r, w)
}

func SetValidationErrors(r *http.Request, errors map[string]string) {
	session, _ := Store.Get(r, "session")
	session.Values["validation_errors"] = errors
	session.Save(r, w)
}

func GetValidationErrors(r *http.Request) map[string]string {
	session, _ := Store.Get(r, "session")
	errors, ok := session.Values["validation_errors"].(map[string]string)
	if !ok {
		return nil
	}
	// Limpiar después de leer
	delete(session.Values, "validation_errors")
	session.Save(r, w)
	return errors
}

// Helpers
func StringToInt(s string) int {
	val, _ := strconv.Atoi(s)
	return val
}

func GetFlashMessages(w http.ResponseWriter, r *http.Request, key string) []interface{} {
	session, _ := Store.Get(r, "session")
	flashes := session.Flashes(key)
	if len(flashes) > 0 {
		session.Save(r, w)
	}
	return flashes
}

// TemplateData estructura para pasar datos a los templates
type TemplateData struct {
	Title            string
	IsAuthenticated  bool
	ShowHeader       bool
	ShowFooter       bool
	FlashError       string
	FlashSuccess     string
	ValidationErrors map[string]string
	FieldErrors      map[string]string
	FormData         map[string]string
	Data             interface{} // Datos específicos de la página
}

// RenderTemplate renderiza templates con layout
func RenderTemplate(w http.ResponseWriter, tmpl string, data TemplateData) {
	// Configurar valores por defecto
	if data.Title == "" {
		data.Title = "Arka Code"
	}
	if data.ShowHeader {
		data.ShowHeader = true
	}
	if data.ShowFooter {
		data.ShowFooter = true
	}

	// Parsear template base y el template específico
	t, err := template.ParseFiles(
		"templates/layouts/base.html",
		"templates/"+tmpl,
	)
	if err != nil {
		http.Error(w, "Error al cargar template: "+err.Error(), http.StatusInternalServerError)
		return
	}

	// Ejecutar template
	err = t.ExecuteTemplate(w, "base", data)
	if err != nil {
		http.Error(w, "Error al renderizar template: "+err.Error(), http.StatusInternalServerError)
	}
}

// Helper para crear TemplateData básico
func NewTemplateData(r *http.Request) TemplateData {
	return TemplateData{
		IsAuthenticated:  IsAuthenticated(r),
		ShowHeader:       true,
		ShowFooter:       true,
		ValidationErrors: GetValidationErrors(r),
	}
}
