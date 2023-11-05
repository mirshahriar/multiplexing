package http

import (
	"fmt"
	"net/http"
)

func NewHTTPServer() *http.Server {
	mux := http.NewServeMux()
	mux.HandleFunc("/echo", func(w http.ResponseWriter, r *http.Request) {
		fmt.Println("HTTP server received request")
		_, _ = w.Write([]byte("Hello from HTTP!"))
	})

	return &http.Server{
		Handler: mux,
	}
}
