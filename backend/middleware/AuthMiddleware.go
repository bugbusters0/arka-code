package middleware

import (
	"arka-code/backend/commons"
	"arka-code/backend/config"
	"net/http"
	"strings"
)

func AuthMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Rutas públicas
		if isPublicRoute(r.URL.Path) {
			next.ServeHTTP(w, r)
			return
		}

		// Verificar autenticación
		if !commons.IsAuthenticated(r) {
			http.Redirect(w, r, config.URLROOT+"/login", http.StatusFound)
			return
		}

		// Middleware específico para familyOnly
		if r.URL.Path == "/seleccionar-perfil" {
			if commons.GetUserType(r) != "family" {
				http.Redirect(w, r, config.URLROOT+"/concepto/gasto", http.StatusFound)
				return
			}
		}

		next.ServeHTTP(w, r)
	})
}

func isPublicRoute(path string) bool {
	publicRoutes := []string{
		"/login",
		"/registro",
		"/public/",
		"/logout",
	}

	for _, route := range publicRoutes {
		if path == route || strings.HasPrefix(path, route) {
			return true
		}
	}
	return false
}

// Funciones estáticas como en tu PHP
func Guest() func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			if commons.IsAuthenticated(r) {
				http.Redirect(w, r, config.URLROOT+"/seleccionar-perfil", http.StatusFound)
				return
			}
			next.ServeHTTP(w, r)
		})
	}
}

func FamilyOnly() func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			if commons.GetUserType(r) != "family" {
				http.Redirect(w, r, config.URLROOT+"/concepto/gasto", http.StatusFound)
				return
			}
			next.ServeHTTP(w, r)
		})
	}
}
