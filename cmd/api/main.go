package main

import (
	"database/sql"
	"log"
	"net/http"

	"github.com/mauFade/playzy/internal/http/handler"
)

func main() {
	db, err := sql.Open("mysql", "books:books@tcp(host:3306)/books")
	if err != nil {
		log.Fatal("Error opening db:", err)
	}
	defer db.Close()

	createUserHandler := handler.NewCreateUserHandler(db)

	http.HandleFunc("/create-user", createUserHandler.Handle)

	log.Println("Server is running on port 8080...")
	log.Fatal(http.ListenAndServe(":8080", nil))
}
