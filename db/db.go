package db

import (
	"log"
	"os"

	"github.com/jmoiron/sqlx"
	"github.com/joho/godotenv"
	_ "github.com/lib/pq"
)

var DB *sqlx.DB

func InitDB() {
	err := godotenv.Load()
	if err != nil {
		log.Panic(err)
	}

	dbURL := os.Getenv("DATABASE_URL")

	DB = sqlx.MustConnect("postgres", dbURL)
	if err := DB.Ping(); err != nil {
		log.Panic(err)
	}
}
