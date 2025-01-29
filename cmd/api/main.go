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
	"github.com/rs/cors"
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

	db.Exec("CREATE TABLE IF NOT EXISTS users (id UUID PRIMARY KEY, name VARCHAR NOT NULL, email VARCHAR NOT NULL, phone VARCHAR NOT NULL, password VARCHAR NOT NULL, gamertag VARCHAR NOT NULL, is_deleted BOOLEAN NOT NULL, deleted_at TIMESTAMP NULL, updated_at TIMESTAMP NOT NULL, created_at TIMESTAMP NOT NULL)")
	db.Exec("CREATE TABLE IF NOT EXISTS sessions (id UUID PRIMARY KEY, game VARCHAR NOT NULL, user_id UUID NOT NULL, objective VARCHAR NOT NULL, rank VARCHAR NULL, is_ranked BOOLEAN NOT NULL, updated_at TIMESTAMP NOT NULL, created_at TIMESTAMP NOT NULL, CONSTRAINT fk_user FOREIGN KEY (user_id) REFERENCES users (id) ON DELETE CASCADE)")

	if err != nil {
		log.Fatal("Error opening db:", err)
	}
	defer db.Close()

	router := routes.Router(db)

	c := cors.New(cors.Options{
		AllowedOrigins:   []string{"http://localhost:3000"},
		AllowCredentials: true,
		AllowedHeaders:   []string{"Content-Type", "Authorization"},
		AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		Debug:            true,
	})

	handler := c.Handler(router)

	log.Println("Server is running on port 8080...")
	log.Fatal(http.ListenAndServe(":8080", handler))

}
