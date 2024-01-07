package main

import (
	"fmt"
	"net/http"
)

func main() {
	myApp := NewApp()
	myApp.configure_routes()

	server := &http.Server{
		Addr:    ":8080",
		Handler: myApp.r,
	}

	err := server.ListenAndServe()
	if err != nil {
		fmt.Println("Error starting the server:", err)
	}
}

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
