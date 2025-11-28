package utils

import (
	"arka-code/config"
	"arka-code/entities"
	"log"
	"net/http"

	"github.com/gorilla/sessions"
)

var store *sessions.CookieStore

func InitSession() {
	store = sessions.NewCookieStore([]byte(config.AppConfig.SessionKey))
	store.Options = &sessions.Options{
		Path:     "/",                  // Cookie válida para todo el sitio
		MaxAge:   86400 * 7,            // Expira en 7 días (86400 segundos = 1 día)
		HttpOnly: true,                 // No accesible desde JavaScript (seguridad)
		Secure:   false,                // false = funciona en HTTP, true = solo HTTPS
		SameSite: http.SameSiteLaxMode, // Protección contra CSRF
	}
}

func GetSession(r *http.Request) (*sessions.Session, error) {
	return store.Get(r, "arka-session") // "arka-session" = nombre de la cookie
}

func SetFamiliaSession(w http.ResponseWriter, r *http.Request, familia *entities.Familia) error {
	session, err := GetSession(r)
	if err != nil {
		return err
	}
	// Guarda estos datos en la sesión:
	session.Values["authenticated_familia"] = true
	session.Values["correo_familia"] = familia.Correo

	return session.Save(r, w) // Guarda en cookies
}
func SetUserSession(w http.ResponseWriter, r *http.Request, usuario *entities.Usuario) error {
	session, err := GetSession(r)
	if err != nil {
		return err
	}
	// Guarda estos datos en la sesión:
	session.Values["authenticated_usuario"] = true
	session.Values["nombre_usuario"] = usuario.NombreUsuario
	session.Values["rol"] = usuario.Rol
	session.Values["nombre_personal_usuario"] = usuario.NombrePersonal

	return session.Save(r, w) // Guarda en cookies
}

func ClearSessionFamilia(w http.ResponseWriter, r *http.Request) error {
	session, err := GetSession(r)
	if err != nil {
		return err
	}

	session.Values["authenticated_familia"] = false
	delete(session.Values, "correo_familia")
	session.Options.MaxAge = -1

	return session.Save(r, w)
}
func ClearSessionUsuario(w http.ResponseWriter, r *http.Request) error {
	session, err := GetSession(r)
	if err != nil {
		return err
	}

	session.Values["authenticated_usuario"] = false
	delete(session.Values, "nombre_usuario")
	delete(session.Values, "rol")
	delete(session.Values, "nombre_personal_usuario")
	session.Options.MaxAge = -1

	return session.Save(r, w)
}

func IsAuthenticated(r *http.Request) bool {
	session, err := GetSession(r)
	if err != nil {
		log.Printf("❌ Error obteniendo sesión: %v", err)
		return false
	}

	auth, ok := session.Values["authenticated_familia"].(bool)
	log.Printf("🔍 Sesión - authenticated_familia: %v, ok: %v, valores: %+v", auth, ok, session.Values)
	return ok && auth // true si está autenticado
}

func GetSessionData(r *http.Request) (*entities.SessionData, bool) {
	session, err := GetSession(r)
	if err != nil {
		return nil, false
	}

	nombreUsuario, ok1 := session.Values["nombre_usuario"].(string)
	correoFamilia, ok2 := session.Values["correo_familia"].(string)
	rolInt, ok3 := session.Values["rol"].(int8)
	nombrePersonal, ok4 := session.Values["nombre_personal_usuario"].(string)

	if !ok1 || !ok2 || !ok3 || !ok4 {
		return nil, false
	}

	return &entities.SessionData{
		NombreUsuario:  nombreUsuario,
		CorreoFamilia:  correoFamilia,
		Rol:            rolInt,
		NombrePersonal: nombrePersonal,
	}, true
}

func GetNombreUsuario(r *http.Request) (string, bool) {
	session, err := GetSession(r)
	if err != nil {
		return "", false
	}

	nombreUsuario, ok := session.Values["nombre_usuario"].(string)
	return nombreUsuario, ok
}

func GetCorreoFamilia(r *http.Request) (string, bool) {
	session, err := GetSession(r)
	if err != nil {
		return "", false
	}

	correoFamilia, ok := session.Values["correo_familia"].(string)
	return correoFamilia, ok
}

func GetRol(r *http.Request) (int8, bool) {
	session, err := GetSession(r)
	if err != nil {
		return 0, false
	}

	rol, ok := session.Values["rol"].(int8) // 1 admin - 0 miembro
	return rol, ok
}

func IsAdmin(r *http.Request) bool {
	rol, ok := GetRol(r)
	return ok && rol == 1
}

func IsFamiliaAuthenticated(r *http.Request) bool {
	session, err := GetSession(r)
	if err != nil {
		log.Printf("❌ Error obteniendo sesión: %v", err)
		return false
	}

	auth, ok := session.Values["authenticated_familia"].(bool)
	log.Printf("🔍 Sesión Familia - authenticated_familia: %v, ok: %v", auth, ok)
	return ok && auth
}

// Función para verificar si un usuario está autenticado
func IsUsuarioAuthenticated(r *http.Request) bool {
	session, err := GetSession(r)
	if err != nil {
		return false
	}

	auth, ok := session.Values["authenticated_usuario"].(bool)
	return ok && auth
}

// Función para obtener correo de familia de la sesión
func GetCorreoFamiliaSession(r *http.Request) (string, bool) {
	session, err := GetSession(r)
	if err != nil {
		return "", false
	}

	correoFamilia, ok := session.Values["correo_familia"].(string)
	return correoFamilia, ok
}
