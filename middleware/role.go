package middleware

import (
	"arka-code/utils"
	"net/http"
)

// RequireAdmin - Solo permite acceso a administradores (rol = 1)
func RequireAdmin(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if !utils.IsAdmin(r) {
			http.Error(w, "Acceso denegado. Solo administradores.", http.StatusForbidden)
			return
		}
		next(w, r)
	}
}
