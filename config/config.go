package config

import (
	"log"
	"os"
	"runtime"
)

type Config struct {
	Port       string
	DBHost     string
	DBPort     string
	DBUser     string
	DBPassword string
	DBName     string
	SessionKey string
	DBProtocol string // "tcp" o "unix"
	DBSocket   string // Solo para Unix
}

var AppConfig *Config

func LoadConfig() {
	AppConfig = &Config{
		Port:       getEnv("PORT", "8080"),
		DBHost:     getEnv("DB_HOST", "localhost"),
		DBPort:     getEnv("DB_PORT", "3306"),
		DBUser:     getEnv("DB_USER", "root"),
		DBPassword: getEnv("DB_PASSWORD", ""),
		DBName:     getEnv("DB_NAME", "arka_go"),
		SessionKey: getEnv("SESSION_KEY", "mi-clave-secreta"),
		DBProtocol: getEnv("DB_PROTOCOL", "tcp"), // Por defecto TCP
		DBSocket:   getDefaultSocketPath(),
	}

	log.Printf("✅ Configuración cargada - SO: %s, Protocolo: %s", runtime.GOOS, AppConfig.DBProtocol)
}

func getDefaultSocketPath() string {
	switch runtime.GOOS {
	case "linux":
		return "/run/mysqld/mysqld.sock"
	case "darwin": // Mac
		return "/tmp/mysql.sock"
	default: // Windows y otros
		return ""
	}
}

func getEnv(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}
