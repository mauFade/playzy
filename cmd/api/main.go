package main

import (
	"database/sql"
	"log"
	"net/http"

	_ "github.com/mattn/go-sqlite3"
	"github.com/mauFade/playzy/internal/http/routes"
)

func main() {
	db, err := sql.Open("sqlite3", "./playzy.db")

	if err != nil {
		log.Fatal("Error opening db:", err)
	}
	defer db.Close()

	router := routes.Router(db)

	log.Println("Server is running on port 8080...")
	log.Fatal(http.ListenAndServe(":8080", router))

}
