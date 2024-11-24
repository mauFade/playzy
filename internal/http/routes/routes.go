package routes

import (
	"database/sql"
	"net/http"

	"github.com/mauFade/playzy/internal/http/handler"
)

func Router(db *sql.DB) *http.ServeMux {
	createUserHandler := handler.NewCreateUserHandler(db)
	authHanlder := handler.NewAuthenticateUserHandler(db)

	router := http.NewServeMux()

	router.HandleFunc("POST /users", createUserHandler.Handle)
	router.HandleFunc("POST /auth", authHanlder.Handle)

	return router
}
