package httpapi

import (
	"net/http"
	"time"

	"gitlab.com/arkine/l4/3/internal/logger"
)

func LoggingMiddleware(l *logger.Logger, next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()
		next.ServeHTTP(w, r)
		l.Log(r.Method + " " + r.URL.Path + " " + time.Since(start).String())
	})
}
