package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"

	vision "cloud.google.com/go/vision/apiv1"
	// https://pkg.go.dev/github.com/szank/gphoto#GetNewGPhotoCamera
	// https://pkg.go.dev/github.com/saljam/webcam#section-readme
)

func handler(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "Hi there, I love %s!", r.URL.Path[1:])
}

func main() {
	mux := http.NewServeMux()
	mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprint(w, "Welcome page")
	})
	mux.HandleFunc("/guide1", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprint(w, "guide1 page")
	})
	mux.HandleFunc("/google_api_vision", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprint(w, "google_api_vision")
	})

	server := http.Server{
		Addr:    "localhost:8080",
		Handler: mux,
	}
	err := server.ListenAndServe()
	if err != nil {
		panic(err)
	}
}

func google_vision() {
	ctx := context.Background()

	// Creates a client.
	client, err := vision.NewImageAnnotatorClient(ctx)
	if err != nil {
		log.Fatalf("Failed to create client: %v", err)
	}
	defer client.Close()

	// Sets the name of the image file to annotate.
	filename := "../Backend/test_image.jpg"

	file, err := os.Open(filename)
	if err != nil {
		log.Fatalf("Failed to read file: %v", err)
	}
	defer file.Close()
	image, err := vision.NewImageFromReader(file)
	if err != nil {
		log.Fatalf("Failed to create image: %v", err)
	}

	labels, err := client.DetectLabels(ctx, image, nil, 10)
	if err != nil {
		log.Fatalf("Failed to detect labels: %v", err)
	}

	fmt.Println("Labels:")
	for _, label := range labels {
		fmt.Println(label.Description)
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
