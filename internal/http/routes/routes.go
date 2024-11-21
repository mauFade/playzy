package routes

import (
	"database/sql"
	"net/http"

	"github.com/mauFade/playzy/internal/http/handler"
)

func Router(db *sql.DB) *http.ServeMux {
	createUserHandler := handler.NewCreateUserHandler(db)

	router := http.NewServeMux()

	router.HandleFunc("POST /users", createUserHandler.Handle)

	return router
}
