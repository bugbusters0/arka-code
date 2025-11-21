package middleware

import (
	"arka-code/utils"
	"net/http"
)

// RequireAdmin protege rutas que requieren privilegios de administrador
// Flujo:
// - Verifica si el usuario actual tiene rol de administrador usando utils.IsAdmin(r)
// - Si NO es administrador:
//   - Retorna error con mensaje "Acceso denegado. Solo administradores."
//   - Detiene la ejecución
//
// - Si ES administrador: Ejecuta el handler siguiente
// Uso típico: mux.HandleFunc("/admin/configuracion", middleware.RequireAdmin(handler))

func RequireAdmin(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if !utils.IsAdmin(r) {
			http.Error(w, "Acceso denegado. Solo administradores.", http.StatusForbidden)
			return
		}
		next(w, r)
	}
}
