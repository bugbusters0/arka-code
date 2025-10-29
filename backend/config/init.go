package config

import (
	"database/sql"
	"log"
	"os"

	_ "github.com/lib/pq"
)

var (
	DB      *sql.DB
	ROOT    string
	URLROOT string
)

func Init() {
	// Configurar paths
	ROOT = getRootPath()
	URLROOT = "http://localhost:8080" // Cambiar según entorno

	// Inicializar base de datos
	initDatabase()
}

func getRootPath() string {
	dir, err := os.Getwd()
	if err != nil {
		log.Fatal(err)
	}
	return dir
}

func initDatabase() {
	var err error
	connStr := "user=tu_usuario dbname=arka_code password=tu_password host=localhost sslmode=disable"
	DB, err = sql.Open("postgres", connStr)
	if err != nil {
		log.Fatal("Error conectando a la base de datos:", err)
	}

	err = DB.Ping()
	if err != nil {
		log.Fatal("Error haciendo ping:", err)
	}

	log.Println("Conectado a PostgreSQL")
}
