package routes

import (
	"arka-code/controllers"
	"arka-code/middleware"
	"arka-code/utils"
	"net/http"
)

func SetupRoutes() *http.ServeMux {
	utils.InitSession()
	mux := http.NewServeMux()

	// Archivos estáticos
	fs := http.FileServer(http.Dir("./static"))
	mux.Handle("/static/", http.StripPrefix("/static/", fs))

	// ========================================
	// RUTAS PÚBLICAS
	// ========================================
	mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		if utils.IsAuthenticated(r) {
			http.Redirect(w, r, "/conceptos", http.StatusSeeOther)
		} else {
			http.Redirect(w, r, "/login", http.StatusSeeOther)
		}
	})

	// ========================================
	// AUTENTICACIÓN DE FAMILIA
	// ========================================
	mux.HandleFunc("/registro", middleware.RequireGuest(
		controllers.AuthControllerInstance.ShowRegister,
	))

	mux.HandleFunc("/registro/submit", middleware.RequireGuest(
		controllers.AuthControllerInstance.Register,
	))

	mux.HandleFunc("/login", middleware.RequireGuest(
		controllers.AuthControllerInstance.ShowLogin,
	))

	mux.HandleFunc("/login/submit", middleware.RequireGuest(
		controllers.AuthControllerInstance.Login,
	))

	mux.HandleFunc("/seleccionar-perfil", middleware.RequireFamiliaAuth(
		controllers.PerfilControllerInstance.ShowSeleccionarPerfil,
	))

	mux.HandleFunc("/seleccionar-perfil/submit", middleware.RequireFamiliaAuth(
		controllers.PerfilControllerInstance.SeleccionarPerfil,
	))

	mux.HandleFunc("/logout-familia", middleware.RequireFamiliaAuth(
		controllers.AuthControllerInstance.LogoutFamilia,
	))

	mux.HandleFunc("/logout-usuario", middleware.RequireFamiliaAuth(
		controllers.AuthControllerInstance.LogoutUsuario,
	))

	mux.HandleFunc("/conceptos", middleware.RequireUsuarioAuth(
		controllers.ConceptoControllerInstance.Index,
	))

	mux.HandleFunc("/conceptos/crear", middleware.RequireUsuarioAuth(
		controllers.ConceptoControllerInstance.Crear,
	))

	mux.HandleFunc("/conceptos/editar", middleware.RequireUsuarioAuth(
		controllers.ConceptoControllerInstance.Editar,
	))

	mux.HandleFunc("/conceptos/deshabilitar", middleware.RequireUsuarioAuth(
		controllers.ConceptoControllerInstance.Deshabilitar,
	))

	mux.HandleFunc("/conceptos/habilitar", middleware.RequireUsuarioAuth(
		controllers.ConceptoControllerInstance.Habilitar,
	))

	mux.HandleFunc("/movimientos", middleware.RequireUsuarioAuth(
		controllers.MovimientoControllerInstance.Index,
	))

	mux.HandleFunc("/movimientos/crear", middleware.RequireUsuarioAuth(
		controllers.MovimientoControllerInstance.Crear,
	))

	mux.HandleFunc("/movimientos/editar", middleware.RequireUsuarioAuth(
		controllers.MovimientoControllerInstance.Editar,
	))

	mux.HandleFunc("/movimientos/eliminar", middleware.RequireUsuarioAuth(
		controllers.MovimientoControllerInstance.Eliminar,
	))

	mux.HandleFunc("/balance", middleware.RequireUsuarioAuth(
		controllers.BalanceControllerInstance.Index,
	))

	mux.HandleFunc("/perfiles", middleware.RequireUsuarioAuth(
		controllers.PerfilControllerInstance.Index,
	))

	mux.HandleFunc("/perfiles/crear", middleware.RequireUsuarioAuth(
		middleware.RequireAdmin(
			controllers.PerfilControllerInstance.Crear,
		),
	))

	mux.HandleFunc("/perfiles/editar", middleware.RequireUsuarioAuth(
		controllers.PerfilControllerInstance.Editar,
	))

	mux.HandleFunc("/perfiles/deshabilitar", middleware.RequireUsuarioAuth(
		middleware.RequireAdmin(
			controllers.PerfilControllerInstance.Deshabilitar,
		),
	))

	mux.HandleFunc("/perfiles/limites/crear", middleware.RequireUsuarioAuth(
		controllers.PerfilControllerInstance.CrearLimite,
	))

	mux.HandleFunc("/perfiles/limites/editar", middleware.RequireUsuarioAuth(
		controllers.PerfilControllerInstance.EditarLimite,
	))

	mux.HandleFunc("/perfiles/limites/eliminar", middleware.RequireUsuarioAuth(
		controllers.PerfilControllerInstance.EliminarLimite,
	))

	return mux
}
