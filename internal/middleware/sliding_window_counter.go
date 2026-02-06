package middleware

import "net/http"

func SlidingWindowCounter(h http.HandlerFunc) http.HandlerFunc {
	return h
}
