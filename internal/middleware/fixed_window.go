package middleware

import (
	"encoding/json"
	"net/http"
	"sync"
	"time"
)

type FixedWindow struct {
	endTime         time.Time
	windowLength    time.Duration
	maxRequests     int
	currentRequests int
	mu              sync.Mutex
}

func (fw *FixedWindow) Allow() bool {
	fw.mu.Lock()
	defer fw.mu.Unlock()

	now := time.Now()
	if now.After(fw.endTime) {
		fw.currentRequests = 0
		fw.endTime = now.Add(fw.windowLength)
	}

	if fw.currentRequests < fw.maxRequests {
		fw.currentRequests++
		return true
	}

	return false
}

var fixWin = FixedWindow{
	endTime:         time.Now().Add(100 * time.Second),
	windowLength:    100 * time.Second,
	maxRequests:     100,
	currentRequests: 0,
}

func FixedWindowMW(h http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if !fixWin.Allow() {
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusTooManyRequests)
			json.NewEncoder(w).Encode(map[string]string{"msg": "Too Many Requests"})
			return
		}
		h(w, r)
	}
}
