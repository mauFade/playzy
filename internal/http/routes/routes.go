package routes

import (
	"database/sql"
	"net/http"

	"github.com/mauFade/playzy/internal/http/handler"
	"github.com/mauFade/playzy/internal/http/middleware"
	"github.com/mauFade/playzy/internal/repository"
	"github.com/mauFade/playzy/internal/websocket"
)

func ApplyMiddlewares(handler http.HandlerFunc, middlewares ...func(http.HandlerFunc) http.HandlerFunc) http.HandlerFunc {
	for _, middleware := range middlewares {
		handler = middleware(handler)
	}
	return handler
}

func CommonMiddlewares(handler http.HandlerFunc) http.HandlerFunc {
	return ApplyMiddlewares(handler, middleware.LoggerMiddleware, middleware.EnsureAuthenticatedMiddleware)
}

func Router(db *sql.DB) *http.ServeMux {
	createUserHandler := handler.NewCreateUserHandler(db)
	authHandler := handler.NewAuthenticateUserHandler(db)

	createSessionHandler := handler.NewCreateSessionHandler(db)
	listSessionsHandler := handler.NewListAvailableSessionsHandler(db)

	messageRepo := repository.NewMessageRepository(db)
	wsManager := websocket.NewManager(db, messageRepo)
	go wsManager.Start()

	router := http.NewServeMux()

	router.HandleFunc("POST /users", middleware.LoggerMiddleware(createUserHandler.Handle))
	router.HandleFunc("POST /auth", middleware.LoggerMiddleware(authHandler.Handle))

	router.HandleFunc("POST /sessions", CommonMiddlewares(createSessionHandler.Handle))
	router.HandleFunc("GET /sessions", CommonMiddlewares(listSessionsHandler.Handle))

	router.HandleFunc("GET /ws", middleware.LoggerMiddleware(wsManager.ServeWs))
	router.HandleFunc("GET /conversations", CommonMiddlewares(wsManager.GetConversationHandler))

	return router
}
