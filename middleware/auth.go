package middleware

import (
	"arka-code/utils"
	"log"
	"net/http"
)

func RequireFamiliaAuth(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if !utils.IsFamiliaAuthenticated(r) {
			log.Printf("🔐 Familia no autenticada, redirigiendo a login")
			http.Redirect(w, r, "/login", http.StatusSeeOther)
			return
		}
		next(w, r)
	}
}

func RequireUsuarioAuth(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if !utils.IsUsuarioAuthenticated(r) {
			log.Printf("🔐 Usuario no autenticado, redirigiendo a seleccionar-perfil")
			http.Redirect(w, r, "/seleccionar-perfil", http.StatusSeeOther)
			return
		}
		next(w, r)
	}
}

func RequireAuth(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		isAuth := utils.IsAuthenticated(r)
		log.Printf("🔐 RequireAuth - Ruta: %s, Autenticado: %v", r.URL.Path, isAuth)

		if !isAuth {
			log.Printf("🔀 Redirigiendo a login - No autenticado")
			http.Redirect(w, r, "/login", http.StatusSeeOther)
			return
		}
		next(w, r)
	}
}

func RequireGuest(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		isAuth := utils.IsAuthenticated(r)
		log.Printf("👤 RequireGuest - Ruta: %s, Autenticado: %v", r.URL.Path, isAuth)

		if isAuth {
			log.Printf("🔀 Redirigiendo a dashboard - Ya autenticado")
			http.Redirect(w, r, "/seleccionar-perfil", http.StatusSeeOther)
			return
		}
		log.Printf("✅ Permitiendo acceso a ruta pública")
		next(w, r)
	}
}
