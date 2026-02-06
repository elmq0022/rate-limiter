package main

import (
	"encoding/json"
	"net"
	"net/http"
)

func main() {
	host := ""
	port := "8080"
	address := net.JoinHostPort(host, port)

	mux := http.NewServeMux()
	mux.HandleFunc("GET /unlimited", handler)
	mux.HandleFunc("GET /token-bucket", tokenBucket(handler))
	mux.HandleFunc("GET /leaky-bucket", leakyBucket(handler))
	mux.HandleFunc("GET /window", fixedWindow(handler))
	mux.HandleFunc("GET /sliding-window", slidingWindow(handler))

	if err := http.ListenAndServe(address, mux); err != nil {
		panic("could not serve")
	}
}

func handler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]string{"msg": "200 OK"})
}

func tokenBucket(h http.HandlerFunc) http.HandlerFunc {
	return h
}

func leakyBucket(h http.HandlerFunc) http.HandlerFunc {
	return h
}

func fixedWindow(h http.HandlerFunc) http.HandlerFunc {
	return h
}

func slidingWindow(h http.HandlerFunc) http.HandlerFunc {
	return h
}
