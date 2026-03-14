package middleware

import (
	"net/http"
	"sync"
	"time"
)

type client struct {
	requests int
	reset    time.Time
}

type RateLimiter struct {
	mu      sync.Mutex
	clients map[string]*client
	limit   int
	window  time.Duration
}

func NewRateLimiter(limit int, window time.Duration) *RateLimiter {
	return &RateLimiter{
		clients: make(map[string]*client),
		limit:   limit,
		window:  window,
	}
}

func (rl *RateLimiter) Middleware(next http.Handler) http.Handler {

	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

		apiKeyID, ok := r.Context().Value(APIKeyIDKey).(string)
		if !ok {
			http.Error(w, "missing api key", http.StatusUnauthorized)
			return
		}

		rl.mu.Lock()
		defer rl.mu.Unlock()

		c, exists := rl.clients[apiKeyID]

		if !exists || time.Now().After(c.reset) {
			rl.clients[apiKeyID] = &client{
				requests: 1,
				reset:    time.Now().Add(rl.window),
			}
		} else {

			if c.requests >= rl.limit {
				http.Error(w, "rate limit exceeded", http.StatusTooManyRequests)
				return
			}

			c.requests++
		}

		next.ServeHTTP(w, r)
	})
}
