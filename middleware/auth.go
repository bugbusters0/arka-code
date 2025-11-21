package middleware

import (
	"arka-code/utils"
	"log"
	"net/http"
)

// RequireFamiliaAuth protege rutas que requieren autenticación a nivel de familia
// Flujo:
// - Verifica si la familia está autenticada usando utils.IsFamiliaAuthenticated(r)
// - Si NO está autenticada:
//   - Registra el intento de acceso no autorizado
//   - Redirige a /login
//   - Detiene la ejecución
//
// - Si está autenticada: Ejecuta el handler siguiente
// Uso típico: mux.HandleFunc("/ruta-protegida", middleware.RequireFamiliaAuth(handler))

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

// RequireUsuarioAuth protege rutas que requieren autenticación a nivel de usuario específico
// Flujo:
// - Verifica si el usuario está autenticado usando utils.IsUsuarioAuthenticated(r)
// - Si NO está autenticado:
//   - Registra el intento de acceso no autorizado
//   - Redirige a /seleccionar-perfil
//   - Detiene la ejecución
//
// - Si está autenticado: Ejecuta el handler siguiente
// Uso típico: mux.HandleFunc("/movimientos", middleware.RequireUsuarioAuth(handler))

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

// RequireAuth protege rutas que requieren cualquier tipo de autenticación (familia o usuario)
// Flujo:
// - Verifica si hay cualquier tipo de autenticación usando utils.IsAuthenticated(r)
// - Registra la ruta y estado de autenticación para debugging
// - Si NO está autenticado:
//   - Registra el intento de acceso no autorizado
//   - Redirige a /login
//   - Detiene la ejecución
//
// - Si está autenticado: Ejecuta el handler siguiente
// Uso típico: mux.HandleFunc("/perfil", middleware.RequireAuth(handler))

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

// RequireGuest protege rutas que solo deben ser accesibles para usuarios NO autenticados
// Flujo:
// - Verifica si hay cualquier tipo de autenticación usando utils.IsAuthenticated(r)
// - Registra la ruta y estado de autenticación para debugging
// - Si ESTÁ autenticado:
//   - Registra el intento de acceso ya autenticado
//   - Redirige a /seleccionar-perfil
//   - Detiene la ejecución
//
// - Si NO está autenticado: Ejecuta el handler siguiente
// Uso típico: mux.HandleFunc("/login", middleware.RequireGuest(handler))
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
