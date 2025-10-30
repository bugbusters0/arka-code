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
	// ========================================
	// SELECCIÓN Y CREACIÓN DE USUARIOS
	// ========================================
	// mux.HandleFunc("/usuario/select", middleware.RequireGuest(
	// 	controllers.UsuarioControllerInstance.ShowSelect,
	// ))

	// mux.HandleFunc("/usuario/login", middleware.RequireGuest(
	// 	controllers.UsuarioControllerInstance.Login,
	// ))

	// mux.HandleFunc("/usuario/create", middleware.RequireGuest(
	// 	controllers.UsuarioControllerInstance.ShowCreate,
	// ))

	// mux.HandleFunc("/usuario/create/submit", middleware.RequireGuest(
	// 	controllers.UsuarioControllerInstance.Create,
	// ))

	// // ========================================
	// // RUTAS PRIVADAS (Requieren autenticación)
	// // ========================================
	// mux.HandleFunc("/logout", middleware.RequireAuth(
	// 	controllers.UsuarioControllerInstance.Logout,
	// ))

	// mux.HandleFunc("/dashboard", middleware.RequireAuth(
	// 	controllers.DashboardControllerInstance.Index,
	// ))

	// // ========================================
	// // CONCEPTOS
	// // ========================================
	// mux.HandleFunc("/conceptos", middleware.RequireAuth(
	// 	controllers.ConceptoControllerInstance.Index,
	// ))

	// mux.HandleFunc("/conceptos/create", middleware.RequireAuth(
	// 	controllers.ConceptoControllerInstance.ShowCreate,
	// ))

	// mux.HandleFunc("/conceptos/create/submit", middleware.RequireAuth(
	// 	controllers.ConceptoControllerInstance.Create,
	// ))

	// mux.HandleFunc("/conceptos/edit", middleware.RequireAuth(
	// 	controllers.ConceptoControllerInstance.ShowEdit,
	// ))

	// mux.HandleFunc("/conceptos/update", middleware.RequireAuth(
	// 	controllers.ConceptoControllerInstance.Update,
	// ))

	// mux.HandleFunc("/conceptos/delete", middleware.RequireAuth(
	// 	controllers.ConceptoControllerInstance.Delete,
	// ))

	// // ========================================
	// // MOVIMIENTOS
	// // ========================================
	// mux.HandleFunc("/movimientos", middleware.RequireAuth(
	// 	controllers.MovimientoControllerInstance.Index,
	// ))

	// mux.HandleFunc("/movimientos/create", middleware.RequireAuth(
	// 	controllers.MovimientoControllerInstance.ShowCreate,
	// ))

	// mux.HandleFunc("/movimientos/create/submit", middleware.RequireAuth(
	// 	controllers.MovimientoControllerInstance.Create,
	// ))

	// mux.HandleFunc("/movimientos/edit", middleware.RequireAuth(
	// 	controllers.MovimientoControllerInstance.ShowEdit,
	// ))

	// mux.HandleFunc("/movimientos/update", middleware.RequireAuth(
	// 	controllers.MovimientoControllerInstance.Update,
	// ))

	// mux.HandleFunc("/movimientos/delete", middleware.RequireAuth(
	// 	controllers.MovimientoControllerInstance.Delete,
	// ))

	// // ========================================
	// // RUTAS SOLO PARA ADMIN
	// // ========================================
	// mux.HandleFunc("/admin/usuarios", middleware.RequireAuth(
	// 	middleware.RequireAdmin(
	// 		controllers.DashboardControllerInstance.AdminUsuarios,
	// 	),
	// ))

	// mux.HandleFunc("/admin/usuarios/rol", middleware.RequireAuth(
	// 	middleware.RequireAdmin(
	// 		controllers.UsuarioControllerInstance.UpdateRol,
	// 	),
	// ))

	// mux.HandleFunc("/admin/usuarios/delete", middleware.RequireAuth(
	// 	middleware.RequireAdmin(
	// 		controllers.UsuarioControllerInstance.SoftDelete,
	// 	),
	// ))

	return mux
}
