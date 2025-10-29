package main

import (
	"arka-code/backend/commons"
	"arka-code/backend/config"
	"arka-code/backend/controllers"
	"arka-code/backend/middleware"
	"log"
	"net/http"
	"strings"
)

func main() {
	// Inicializar configuración
	config.Init()
	defer config.DB.Close()

	// Configurar rutas
	mux := http.NewServeMux()

	// Servir archivos estáticos desde public/
	publicFS := http.FileServer(http.Dir("./public"))
	mux.Handle("/public/", http.StripPrefix("/public/", publicFS))

	// Rutas EXACTAS como en tu PHP
	mux.HandleFunc("/", homeHandler)
	mux.HandleFunc("/login", loginHandler)
	mux.HandleFunc("/registro", registroHandler)
	mux.HandleFunc("/seleccionar-perfil", seleccionarPerfilHandler)
	mux.HandleFunc("/logout", logoutHandler)
	mux.HandleFunc("/concepto/crear", conceptoCrearHandler)
	mux.HandleFunc("/concepto/gasto", conceptoGastoHandler)
	mux.HandleFunc("/concepto/ingresos", conceptoIngresosHandler)
	mux.HandleFunc("/concepto/guardarConcepto", conceptoGuardarHandler)
	mux.HandleFunc("/concepto/editar/", conceptoEditarHandler)
	mux.HandleFunc("/concepto/deshabilitar/", conceptoDeshabilitarHandler)
	mux.HandleFunc("/concepto/eliminar/", conceptoEliminarHandler)

	// Middleware global
	handler := middleware.AuthMiddleware(mux)

	log.Printf("Servidor ejecutándose en %s", config.URLROOT)
	log.Println("Estructura IDÉNTICA a tu PHP:")
	log.Println("✅ backend/config/")
	log.Println("✅ backend/controllers/")
	log.Println("✅ backend/models/")
	log.Println("✅ backend/middleware/")
	log.Println("✅ backend/validator/")
	log.Println("✅ backend/commons/")
	log.Println("✅ public/")

	http.ListenAndServe(":8080", handler)
}

// Handlers que replican EXACTAMENTE tu switch case de PHP
func homeHandler(w http.ResponseWriter, r *http.Request) {
	http.Redirect(w, r, config.URLROOT+"/login", http.StatusFound)
}

func loginHandler(w http.ResponseWriter, r *http.Request) {
	middleware.Guest()(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		controller := controllers.NewAuthController()

		if r.Method == "POST" {
			controller.EnviarCredenciales(w, r)
		} else {
			controller.ShowLogin(w, r)
		}
	})).ServeHTTP(w, r)
}

func registroHandler(w http.ResponseWriter, r *http.Request) {
	middleware.Guest()(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		controller := controllers.NewAuthController()

		if r.Method == "POST" {
			controller.ProcessRegister(w, r)
		} else {
			controller.ShowRegister(w, r)
		}
	})).ServeHTTP(w, r)
}

func seleccionarPerfilHandler(w http.ResponseWriter, r *http.Request) {
	middleware.FamilyOnly()(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Verificar si ya tiene miembro seleccionado
		if commons.IsAuthenticated(r) && commons.GetMiembroID(r) != commons.GetUserID(r) {
			http.Redirect(w, r, config.URLROOT+"/concepto/gasto", http.StatusFound)
			return
		}

		if r.Method == "POST" {
			perfilController := controllers.NewPerfilController()
			perfilController.ConsultarVerificacion(w, r)
		} else if r.URL.Query().Get("profile") != "" {
			authController := controllers.NewAuthController()
			authController.ProcessSeleccionarPerfil(w, r)
		} else {
			perfilController := controllers.NewPerfilController()
			perfilController.SolicitarPerfiles(w, r)
		}
	})).ServeHTTP(w, r)
}

func logoutHandler(w http.ResponseWriter, r *http.Request) {
	controller := controllers.NewAuthController()
	controller.Logout(w, r)
}

func conceptoCrearHandler(w http.ResponseWriter, r *http.Request) {
	controller := controllers.NewConceptoController()

	if r.Method == "POST" {
		controller.GuardarConcepto(w, r)
	} else {
		controller.MostrarConceptos(w, r, "gasto") // Default como en tu PHP
	}
}

func conceptoGastoHandler(w http.ResponseWriter, r *http.Request) {
	controller := controllers.NewConceptoController()
	controller.MostrarConceptos(w, r, "gasto")
}

func conceptoIngresosHandler(w http.ResponseWriter, r *http.Request) {
	controller := controllers.NewConceptoController()
	controller.MostrarConceptos(w, r, "ingreso")
}

func conceptoGuardarHandler(w http.ResponseWriter, r *http.Request) {
	controller := controllers.NewConceptoController()
	controller.GuardarConcepto(w, r)
}

func conceptoEditarHandler(w http.ResponseWriter, r *http.Request) {
	controller := controllers.NewConceptoController()

	if r.Method == "POST" {
		// Extraer ID de la URL para POST también
		path := r.URL.Path
		idStr := strings.TrimPrefix(path, "/concepto/editar/")
		controllerID := commons.StringToInt(idStr)

		// Lógica de actualización
		nombre := r.FormValue("nombre")
		tipo := r.FormValue("tipo")
		controller.Update(controllerID, nombre, tipo)
		http.Redirect(w, r, config.URLROOT+"/concepto/"+tipo, http.StatusFound)
	} else {
		controller.Editar(w, r)
	}
}

func conceptoDeshabilitarHandler(w http.ResponseWriter, r *http.Request) {
	controller := controllers.NewConceptoController()
	controller.Deshabilitar(w, r)
}

func conceptoEliminarHandler(w http.ResponseWriter, r *http.Request) {
	controller := controllers.NewConceptoController()
	controller.Eliminar(w, r)
}
