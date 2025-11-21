package main

import (
	"arka-code/config"
	"arka-code/database"
	"arka-code/routes"
	"arka-code/scheduler"
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"
)

func main() {
	// Cargar configuración
	config.LoadConfig()

	// Conectar base de datos
	database.Connect()
	defer database.Close()

	// Crear contexto que se cancela con señales del sistema
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	// Iniciar scheduler de tareas automáticas
	schedulerInstance := scheduler.NewScheduler()
	schedulerInstance.Start(ctx) // Pasar el contexto al scheduler

	// Configurar rutas
	router := routes.SetupRoutes()

	// Crear servidor HTTP con timeouts configurados
	server := &http.Server{
		Addr:         ":" + config.AppConfig.Port,
		Handler:      router,
		ReadTimeout:  15 * time.Second,
		WriteTimeout: 15 * time.Second,
		IdleTimeout:  60 * time.Second,
	}

	// Canal para manejar señales de cierre
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, os.Interrupt, syscall.SIGTERM)

	// Iniciar servidor en una goroutine
	go func() {
		log.Printf("🚀 Servidor Arka iniciado en http://localhost:%s", config.AppConfig.Port)
		if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("❌ Error iniciando servidor: %v", err)
		}
	}()

	// Esperar señal de cierre
	<-sigChan
	log.Println("🛑 Recibida señal de cierre, cerrando servidor...")

	// Crear contexto con timeout para el cierre graceful
	shutdownCtx, shutdownCancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer shutdownCancel()

	// Intentar cierre graceful del servidor
	if err := server.Shutdown(shutdownCtx); err != nil {
		log.Printf("⚠️  Error en cierre graceful: %v", err)
	}

	// Cancelar el contexto principal (detendrá el scheduler)
	cancel()

	// Dar tiempo para que las goroutines se cierren
	time.Sleep(2 * time.Second)

	log.Println("✅ Servidor cerrado exitosamente")
}
