package middleware

import (
	"log"
	"net/http"
	"time"
)

func StructuredLogging(level string) Middleware {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			started := time.Now()
			next.ServeHTTP(w, r)
			log.Printf("level=%s method=%s path=%s request_id=%s duration=%s", level, r.Method, r.URL.Path, RequestIDFromContext(r.Context()), time.Since(started))
		})
	}
}
