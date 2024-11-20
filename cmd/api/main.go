package main

import (
	"database/sql"
	"log"
	"net/http"

	"github.com/mauFade/playzy/internal/http/handler"

	_ "github.com/mattn/go-sqlite3"
)

func main() {
	db, err := sql.Open("sqlite3", "./playzy.db")

	if err != nil {
		log.Fatal("Error opening db:", err)
	}
	defer db.Close()

	createUserHandler := handler.NewCreateUserHandler(db)

	router := http.NewServeMux()

	router.HandleFunc("POST /users", createUserHandler.Handle)

	log.Println("Server is running on port 8080...")
	log.Fatal(http.ListenAndServe(":8080", router))

}
