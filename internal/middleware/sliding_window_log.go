package middleware

import (
	"encoding/json"
	"net/http"
	"sync"
	"time"
)

type SlidingWindowLog struct {
	maxRequests int
	requestLog  []time.Time
	logDuration time.Duration
	mu          sync.Mutex
}

func (s *SlidingWindowLog) Allow() bool {
	s.mu.Lock()
	defer s.mu.Unlock()

	now := time.Now()
	start := now.Add(-s.logDuration)

	index := 0
	for i, t := range s.requestLog {
		if t.Before(start) {
			index = i + 1
		} else {
			break
		}
	}

	n := copy(s.requestLog, s.requestLog[index:])
	s.requestLog = s.requestLog[:n]

	if n < s.maxRequests {
		s.requestLog = append(s.requestLog, now)
		return true
	}
	return false
}

var swl = SlidingWindowLog{
	maxRequests: 100,
	requestLog:  []time.Time{},
	logDuration: 100 * time.Second,
}

func SlidingWindowLogMW(h http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if !swl.Allow() {
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusTooManyRequests)
			json.NewEncoder(w).Encode(map[string]string{"msg": "Too Many Requests"})
			return
		}
		h(w, r)
	}
}
