package middleware

import (
	"log"
	"net/http"
)

func LoggerMiddleware(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		log.Printf("METHOD: %s, PATH: %s, ADDR: %s, USERAGENT: %s",
			r.Method, r.URL.Path, r.RemoteAddr, r.UserAgent())

		next(w, r)
	}
}
