package config

import (
	"log"
	"os"
)

type Config struct {
	Port       string
	DBHost     string
	DBPort     string
	DBUser     string
	DBPassword string
	DBName     string
	SessionKey string
	DBSocket   string // Agregar esto
}

var AppConfig *Config

func LoadConfig() {
	AppConfig = &Config{
		Port:       getEnv("PORT", "8080"),
		DBHost:     getEnv("DB_HOST", "localhost"),
		DBPort:     getEnv("DB_PORT", "3306"),
		DBUser:     getEnv("DB_USER", "root"),
		DBPassword: getEnv("DB_PASSWORD", ""),
		DBName: getEnv("DB_NAME", "arka"),
		SessionKey: getEnv("SESSION_KEY", "mi-clave-secreta"),
		DBSocket: getEnv("DB_SOCKET", ""), // Socket por defecto
	}

	log.Printf("✅ Configuración cargada - DB: %s@%s", AppConfig.DBUser, AppConfig.DBSocket)
}

func getEnv(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}