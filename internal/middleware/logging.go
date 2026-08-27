package middleware

import (
	"net/http"
	"time"

	"leaderboard/pkg/logger"
)

func LoggingMiddleware(log *logger.Logger) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			start := time.Now()
			next.ServeHTTP(w, r)
			log.Infof("%s %s %s", r.Method, r.URL.Path, time.Since(start))
		})
	}
}
