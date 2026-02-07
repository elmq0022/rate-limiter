package middleware

import "net/http"

func SlidingWindowCounterMW(h http.HandlerFunc) http.HandlerFunc {
	return h
}
