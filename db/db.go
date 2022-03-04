package db

import (
	"fmt"
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
	var (
		host     = os.Getenv("PG_HOST")
		port     = os.Getenv("PG_PORT")
		user     = os.Getenv("PG_USER")
		password = os.Getenv("PG_PASSWORD")
		dbname   = os.Getenv("PG_DBNAME")
	)

	psqlConn := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=disable",
		host, port, user, password, dbname)
	DB = sqlx.MustConnect("postgres", psqlConn)
	if err := DB.Ping(); err != nil {
		log.Panic(err)
	}
}
