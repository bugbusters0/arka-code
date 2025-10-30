package main

import (
	"arka-code/config"
	"arka-code/database"
	"arka-code/routes"
	"log"
	"net/http"
)

func main() {
	// Cargar configuración
	config.LoadConfig()

	// Conectar base de datos
	database.Connect()
	defer database.Close()

	// Configurar rutas
	router := routes.SetupRoutes()

	log.Printf("🚀 Servidor Arka iniciado en http://localhost:%s", config.AppConfig.Port)
	log.Fatal(http.ListenAndServe(":"+config.AppConfig.Port, router))
}
