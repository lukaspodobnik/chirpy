package main

import (
	"fmt"
	"net/http"
)

func main() {
	fmt.Println("Starting server...")
	mux := http.NewServeMux()
	server := http.Server{
		Addr:    ":8080",
		Handler: mux,
	}
	err := server.ListenAndServe()
	fmt.Println(err)
}
