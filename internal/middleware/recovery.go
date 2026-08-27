package middleware

import (
	"net/http"
	"runtime/debug"

	"leaderboard/pkg/httpx"
	"leaderboard/pkg/logger"
)

func RecoveryMiddleware(log *logger.Logger) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			defer func() {
				if rec := recover(); rec != nil {
					log.Errorf("panic: %v\n%s", rec, debug.Stack())
					httpx.InternalError(w, "服务器内部错误")
				}
			}()
			next.ServeHTTP(w, r)
		})
	}
}
