package middleware

import (
	"encoding/json"
	"net/http"
	"sync"
	"time"
)

type TokenBucket struct {
	mu             sync.Mutex
	lastAccessTime time.Time
	capacity       float64
	tokens         float64
	refillRate     float64
}

func (tb *TokenBucket) refill() {
	now := time.Now()
	tokens := tb.tokens + float64(now.Sub(tb.lastAccessTime).Milliseconds())*(tb.refillRate/1000)
	tb.tokens = min(tokens, tb.capacity)
	tb.lastAccessTime = now
}

func (tb *TokenBucket) Allow() bool {
	tb.mu.Lock()
	defer tb.mu.Unlock()
	tb.refill()

	if tb.tokens > 0 {
		tb.tokens--
		return true
	}

	return false
}

var tb *TokenBucket = &TokenBucket{
	capacity:       100.0,
	tokens:         100.0,
	refillRate:     100.0 / 60.0,
	lastAccessTime: time.Now(),
}

func TokenBucketMW(h http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if !tb.Allow() {
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusTooManyRequests)
			json.NewEncoder(w).Encode(map[string]string{"msg": "Too Many Requests"})
		} else {
			h(w, r)
		}
	}
}
