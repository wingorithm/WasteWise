package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"

	vision "cloud.google.com/go/vision/apiv1"
	"github.com/gorilla/mux"
)

type App struct {
	r *mux.Router
}

func NewApp() *App {
	return &App{
		r: mux.NewRouter(),
	}
}

func (a *App) configure_routes() {
	a.r.HandleFunc("/", idle_page).Methods("GET")
	a.r.HandleFunc("/welcome", welcome_page).Methods("GET")
	a.r.HandleFunc("/scan", scan_page).Methods("GET")
}

func idle_page(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "IDLE PAGE 1")
}

func welcome_page(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
	w.Write([]byte(`{"message": "welcome page"}`))
}

func scan_page(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	fmt.Fprintf(w, "SCAN PAGE 3")
	var labelDesc []string

	labelDesc = google_vision("../Backend/test_image.jpg", labelDesc)

	w.WriteHeader(http.StatusOK)

	for _, label := range labelDesc {
		w.Write([]byte(label))
	}
}

func google_vision(filename string, labelDesc []string) []string {

	ctx := context.Background()

	// Creates a client.
	client, err := vision.NewImageAnnotatorClient(ctx)
	if err != nil {
		log.Fatalf("Failed to create client: %v", err)
	}
	defer client.Close()

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

	for _, label := range labels {
		labelDesc = append(labelDesc, label.Description)
	}
	return labelDesc
}
