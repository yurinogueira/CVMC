package middleware

import (
	"context"
	"crypto/rand"
	"encoding/hex"
	"net/http"
)

type requestIDKey struct{}
type correlationIDKey struct{}

func RequestID(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		requestID := valueOrNew(r.Header.Get("X-Request-ID"))
		correlationID := valueOrNew(r.Header.Get("X-Correlation-ID"))

		w.Header().Set("X-Request-ID", requestID)
		w.Header().Set("X-Correlation-ID", correlationID)
		ctx := context.WithValue(r.Context(), requestIDKey{}, requestID)
		ctx = context.WithValue(ctx, correlationIDKey{}, correlationID)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

func RequestIDFromContext(ctx context.Context) string {
	value, _ := ctx.Value(requestIDKey{}).(string)
	return value
}

func CorrelationIDFromContext(ctx context.Context) string {
	value, _ := ctx.Value(correlationIDKey{}).(string)
	return value
}

func valueOrNew(value string) string {
	if value != "" {
		return value
	}
	return newID()
}

func newID() string {
	buf := make([]byte, 16)
	if _, err := rand.Read(buf); err != nil {
		return "unknown"
	}
	return hex.EncodeToString(buf)
}
