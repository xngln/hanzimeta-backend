package db

import (
	"log"
	"os"

	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
)

var DB *sqlx.DB

func InitDB() {
	dbURL := os.Getenv("DATABASE_URL")

	DB = sqlx.MustConnect("postgres", dbURL)
	if err := DB.Ping(); err != nil {
		log.Panic(err)
	}
}
