package http

import (
	"fmt"
	"net/http"
)

func NewHTTPServer() *http.Server {
	mux := http.NewServeMux()
	mux.HandleFunc("/echo", func(w http.ResponseWriter, r *http.Request) {
		fmt.Println("HTTP server request received")
		_, _ = w.Write([]byte("echo from HTTP!\n"))
	})

	return &http.Server{
		Handler: mux,
	}
}
