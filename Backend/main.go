package main

import (
	"fmt"
	"net/http"

	"github.com/rs/cors"
)

func main() {
	myApp := NewApp()
	myApp.configure_routes()
	handler := cors.AllowAll().Handler(myApp.r)

	server := &http.Server{
		Addr:    ":8080",
		Handler: handler,
	}

	err := server.ListenAndServe()
	if err != nil {
		fmt.Println("Error starting the server:", err)
	}
}

// func enableCors(w *http.ResponseWriter, r *http.Request) {
// 	(*w).Header().Set("Access-Control-Allow-Origin", "*")
// 	(*w).Header().Set("Access-Control-Allow-Methods", "GET, POST, DELETE, OPTIONS")
// 	(*w).Header().Set("Access-Control-Allow-Headers", "Content-Type")

// 	// Handle preflight (OPTIONS) requests
// 	if r.Method == "OPTIONS" {
// 		(*w).WriteHeader(http.StatusOK)
// 		return
// 	}
// }

// wastewise-410108

// go vision package
// https://pkg.go.dev/cloud.google.com/go/vision/apiv1?utm_source=godoc
// cloud.google.com/go/vision

// Web APP Documentation go
// https://go.dev/doc/articles/wiki/#tmp_1

// Project Structure Go
// https://blog.logrocket.com/flat-structure-vs-layered-architecture-structuring-your-go-app/

// go to firestore
// https://gist.github.com/tiebingzhang/b7c6284d3f5e6eab901010377f924f3f
// https://cloud.google.com/firestore/docs/create-database-server-client-library
