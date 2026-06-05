package main

import (
	"net/http"
)

func main() {
	mux := http.NewServeMux()

	// File server serving current directory
	fileServer := http.FileServer(http.Dir("."))

	// Handle root path
	mux.Handle("/", fileServer)

	server := &http.Server{
		Addr:    ":8080",
		Handler: mux,
	}

	server.ListenAndServe()
}
