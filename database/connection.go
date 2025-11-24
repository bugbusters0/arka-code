package database

import (
	"arka-code/config"
	"database/sql"
	"fmt"
	"log"

	_ "github.com/go-sql-driver/mysql"
)

var DB *sql.DB

func Connect() {
	cfg := config.AppConfig

	// Conexión TCP para Windows / XAMPP
	dsn := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?parseTime=true&charset=utf8mb4&loc=Local",
		cfg.DBUser,
		cfg.DBPassword,
		cfg.DBHost,
		cfg.DBPort,
		cfg.DBName,
	)

	var err error
	DB, err = sql.Open("mysql", dsn)
	if err != nil {
		log.Fatal("Error al conectar con la base de datos:", err)
	}

	if err = DB.Ping(); err != nil {
		log.Fatal("Error al hacer ping a la base de datos:", err)
	}

	log.Println("✓ Conexión exitosa a MySQL por TCP")
}

func Close() {
	if DB != nil {
		DB.Close()
	}
}