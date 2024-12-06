package routes

import (
	"database/sql"
	"net/http"

	"github.com/mauFade/playzy/internal/http/handler"
	"github.com/mauFade/playzy/internal/http/middleware"
)

func Router(db *sql.DB) *http.ServeMux {
	createUserHandler := handler.NewCreateUserHandler(db)
	authHandler := handler.NewAuthenticateUserHandler(db)

	createSessionHandler := handler.NewCreateSessionHandler(db)

	router := http.NewServeMux()

	router.HandleFunc("POST /users", middleware.LoggerMiddleware(createUserHandler.Handle))
	router.HandleFunc("POST /auth", middleware.LoggerMiddleware(authHandler.Handle))

	router.HandleFunc("POST /sessions", middleware.LoggerMiddleware(middleware.EnsureAuthenticatedMiddleware(createSessionHandler.Handle)))

	return router
}
