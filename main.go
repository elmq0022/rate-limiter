package main

import (
	"encoding/json"
	"net"
	"net/http"

	"github.com/elmq0022/rate-limiter/internal/middleware"
)

func main() {
	host := ""
	port := "8080"
	address := net.JoinHostPort(host, port)

	mux := http.NewServeMux()
	mux.HandleFunc("GET /unlimited", handler)
	mux.HandleFunc("GET /token-bucket", middleware.TokenBucket(handler))
	mux.HandleFunc("GET /window", middleware.FixedWindow(handler))
	mux.HandleFunc("GET /sliding-window", middleware.SlidingWindowCounter(handler))

	if err := http.ListenAndServe(address, mux); err != nil {
		panic("could not serve")
	}
}

func handler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]string{"msg": "200 OK"})
}
