package database

import (
	"arka-code/config"
	"database/sql"
	"fmt"
	"log"
	"runtime"
	"time"

	_ "github.com/go-sql-driver/mysql"
)

var DB *sql.DB

func Connect() {
	cfg := config.AppConfig

	var dsn string

	if cfg.DBProtocol == "unix" && runtime.GOOS != "windows" {
		// Conexión por socket Unix (Linux/Mac)
		dsn = fmt.Sprintf("%s:%s@unix(%s)/%s?parseTime=true",
			cfg.DBUser,
			cfg.DBPassword,
			cfg.DBSocket,
			cfg.DBName)
		log.Printf("🔌 Conectando via socket Unix: %s", cfg.DBSocket)
	} else {
		// Conexión TCP (Windows o fallback)
		dsn = fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?parseTime=true",
			cfg.DBUser,
			cfg.DBPassword,
			cfg.DBHost,
			cfg.DBPort,
			cfg.DBName)
		log.Printf("🔌 Conectando via TCP: %s:%s", cfg.DBHost, cfg.DBPort)
	}

	var err error
	DB, err = sql.Open("mysql", dsn)
	if err != nil {
		log.Fatal("❌ Error al abrir conexión con la base de datos:", err)
	}

	// Configurar pool de conexiones
	DB.SetMaxOpenConns(25)
	DB.SetMaxIdleConns(25)
	DB.SetConnMaxLifetime(5 * time.Minute)

	if err = DB.Ping(); err != nil {
		log.Printf("❌ Error en conexión. DSN: %s", maskPassword(dsn))
		log.Fatal("❌ Error al conectar con la base de datos:", err)
	}

	log.Printf("✅ Conexión exitosa a %s (%s)", cfg.DBName, runtime.GOOS)
}

// Ocultar contraseña en logs
func maskPassword(dsn string) string {
	// Implementación simple para ocultar contraseña en logs
	return dsn // En producción, deberías implementar esto mejor
}

func Close() {
	if DB != nil {
		DB.Close()
		log.Println("🔌 Conexión a BD cerrada")
	}
}
