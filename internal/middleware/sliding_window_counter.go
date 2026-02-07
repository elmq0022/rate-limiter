package middleware

import (
	"encoding/json"
	"net/http"
	"sync"
	"time"
)

// https://codingchallenges.fyi/challenges/challenge-rate-limiter
//
// We use a weighted count of the current and previous windows counts to determine the count for the sliding window.
// This helps smooth out the impact of burst of traffic. For example, if the current window is 40% through,
// then we weight the previous window’s count by 60% and add that to the current window count.

type SlidingWindowCounter struct {
	window1start   time.Time
	window1count   float64
	window2start   time.Time
	window2count   float64
	maxRequests    float64
	windowDuration time.Duration
	mu             sync.Mutex
}

func (s *SlidingWindowCounter) AdjustWindows(now time.Time) {
	// now is in window 2 - do nothing

	// now is after window 2 + 2 x duration the RESET
	if now.After(s.window2start.Add(2 * s.windowDuration)) {
		s.window1count = 0
		s.window1start = now.Add(-s.windowDuration)
		s.window2count = 0
		s.window2start = now
	}

	// if now is after window 2 but less than
	// window 2 plus duration then SLIDE
	if now.After(s.window2start.Add(s.windowDuration)) {
		s.window1count = s.window2count
		s.window1start = s.window2start
		s.window2count = 0
		s.window2start = s.window2start.Add(s.windowDuration)
	}
}

func (s *SlidingWindowCounter) Allow() bool {
	s.mu.Lock()
	defer s.mu.Unlock()

	now := time.Now()
	s.AdjustWindows(now)

	elapsed := s.window2start.Add(s.windowDuration).Sub(now)
	w1 := float64(elapsed.Seconds()) / float64(s.windowDuration.Seconds())
	weightedRequests := w1*s.window1count + s.window2count

	if weightedRequests <= s.maxRequests {
		s.window2count++
		return true
	}
	return false
}

func NewSlidingWindowCounter(d time.Duration, maxRequests float64) *SlidingWindowCounter {
	now := time.Now()
	return &SlidingWindowCounter{
		window1start:   now.Add(-d),
		window2start:   now,
		windowDuration: d,
		window1count:   0,
		window2count:   0,
		maxRequests:    maxRequests,
	}
}

var swc = NewSlidingWindowCounter(time.Minute, 100)

func SlidingWindowCounterMW(h http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if !swc.Allow() {
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusTooManyRequests)
			json.NewEncoder(w).Encode(map[string]string{"msg": "Too Many Requests"})
			return
		}
		h(w, r)
	}
}
