package middleware

import (
	"net/http"
	"strings"

	"leaderboard/internal/config"
	"leaderboard/pkg/httpx"
)

func AuthMiddleware(cfg *config.Config) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			if strings.HasPrefix(r.URL.Path, "/api/") {
				token := strings.TrimPrefix(r.Header.Get("Authorization"), "Bearer ")
				if token == "" {
					token = r.URL.Query().Get("token")
				}
				if token != cfg.APIToken {
					httpx.Unauthorized(w, "无效的访问凭证")
					return
				}
			}
			next.ServeHTTP(w, r)
		})
	}
}
