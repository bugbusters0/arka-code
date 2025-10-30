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
	dsn := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?parseTime=true",
		cfg.DBUser, cfg.DBPassword, cfg.DBHost, cfg.DBPort, cfg.DBName)

	var err error
	DB, err = sql.Open("mysql", dsn)
	if err != nil {
		log.Fatal("Error al conectar con la base de datos:", err)
	}

	if err = DB.Ping(); err != nil {
		log.Fatal("Error al hacer ping a la base de datos:", err)
	}

	log.Println("✓ Conexión exitosa a MySQL")
}

func Close() {
	if DB != nil {
		DB.Close()
	}
}
