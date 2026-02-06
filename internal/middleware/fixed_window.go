package middleware

import "net/http"

func FixedWindow(h http.HandlerFunc) http.HandlerFunc {
	return h
}
