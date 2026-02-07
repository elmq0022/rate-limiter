package middleware

import "net/http"

func SlidingWindowLogMW(h http.HandlerFunc) http.HandlerFunc {
	return h
}
