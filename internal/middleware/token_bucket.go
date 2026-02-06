package middleware

import "net/http"

func TokenBucket(h http.HandlerFunc) http.HandlerFunc {
	return h
}
