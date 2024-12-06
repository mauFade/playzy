package main

import (
	"database/sql"
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/joho/godotenv"
	_ "github.com/lib/pq"
	"github.com/mauFade/playzy/internal/http/routes"
)

func init() {
	err := godotenv.Load()

	if err != nil {
		fmt.Println("error loading .env: ", err.Error())
		os.Exit(2)
	}
}

func main() {
	dbHost := os.Getenv("DB_HOST")
	dbPort := os.Getenv("DB_PORT")
	dbUser := os.Getenv("DB_USER")
	dbPassword := os.Getenv("DB_PASSWORD")
	dbName := os.Getenv("DB_NAME")

	connStr := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=disable",
		dbHost, dbPort, dbUser, dbPassword, dbName)

	db, err := sql.Open("postgres", connStr)

	if err != nil {
		log.Fatal("Error opening db:", err)
	}
	defer db.Close()

	router := routes.Router(db)

	log.Println("Server is running on port 8080...")
	log.Fatal(http.ListenAndServe(":8080", router))

}
