package middleware

import (
	"net/http"
	"sync"
	"time"

	"leaderboard/pkg/httpx"
)

type RateLimiter struct {
	mu     sync.Mutex
	visits map[string][]time.Time
	limit  int
	window time.Duration
}

func NewRateLimiter(limit int, window time.Duration) *RateLimiter {
	return &RateLimiter{
		visits: make(map[string][]time.Time),
		limit:  limit,
		window: window,
	}
}

func (rl *RateLimiter) Allow(ip string) bool {
	now := time.Now()
	rl.mu.Lock()
	defer rl.mu.Unlock()
	list := rl.visits[ip]
	var fresh []time.Time
	for _, t := range list {
		if now.Sub(t) < rl.window {
			fresh = append(fresh, t)
		}
	}
	if len(fresh) >= rl.limit {
		rl.visits[ip] = fresh
		return false
	}
	rl.visits[ip] = append(fresh, now)
	return true
}

func RateLimitMiddleware(rl *RateLimiter) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			ip := r.RemoteAddr
			if !rl.Allow(ip) {
				httpx.Error(w, http.StatusTooManyRequests, 429, "请求过于频繁")
				return
			}
			next.ServeHTTP(w, r)
		})
	}
}
