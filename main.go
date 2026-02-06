package main

import (
	"net"
	"net/http"
)

func main() {
	host := ""
	port := "8080"
	address := net.JoinHostPort(host, port)

	mux := http.NewServeMux()
	mux.HandleFunc("GET /unlimited", handler)
	mux.HandleFunc("GET /token-bucket", handler)
	mux.HandleFunc("GET /leaky-bucket", handler)
	mux.HandleFunc("GET /window", handler)
	mux.HandleFunc("GET /sliding-window", handler)

	if err := http.ListenAndServe(address, mux); err != nil {
		panic("could not serve")
	}
}

func handler(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
}
