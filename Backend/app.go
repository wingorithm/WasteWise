package main

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"strconv"
	"time"

	vision "cloud.google.com/go/vision/apiv1"
	"github.com/gorilla/mux"
	"github.com/gorilla/websocket"
	"gocv.io/x/gocv"
)

func NewApp() *App {
	return &App{
		r: mux.NewRouter(),
	}
}

type App struct {
	r *mux.Router
}

// scanner <- response 1, Trashobject
// award <- response 1,  && dashboard <- response, trashobject
// state:idle (change frontend)
type Event struct {
	event int `json:"event"`
}

// {0: recycleable, 1: organik, 2: lainnya}
type TrashObject struct {
	Event
	Name       string    `json:"name"`
	DetectedAt time.Time `json:"time"`
}

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
	CheckOrigin:     func(r *http.Request) bool { return true },
}

var whiteList = map[string]string{"name": "Sammy", "animal": "shark", "color": "blue", "location": "ocean"}

// nodemcu -> 1 (ada orang) / 0 (gk ada orang)
var state = 999
var lastState = 999

// nodemcu terima sampah -> 2:organik / 3:recycle / 4:another
var trashState = 0
var trashlastState = 0

func (a *App) configure_routes() {
	a.r.HandleFunc("/ws", websocketHandle)
	a.r.HandleFunc("/iot", iotHandler).Methods("POST")   //human
	a.r.HandleFunc("/iot2", iot2Handler).Methods("POST") //barang masuk

	a.r.HandleFunc("/classify", classify_page) // nanti apus
	// a.r.HandleFunc("/classify", classify_page).Methods("POST")
}

func websocketHandle(w http.ResponseWriter, r *http.Request) {
	// Upgrade the HTTP connection to a WebSocket connection.
	upgrader.CheckOrigin = func(r *http.Request) bool { return true }

	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Println(err)
	}

	log.Println("Client connected")

	go stateHandler(conn)
	go stateTrashHandler(conn)
}

func stateHandler(conn *websocket.Conn) {
	log.Println("handler1 runnning")
	for {
		switch state {
		case 1:
			state = 999
			if lastState != state {
				if err := conn.WriteMessage(websocket.TextMessage, []byte("state:start")); err != nil {
					log.Println(err)
					return
				}
			}
		case 0:
			state = 999
			if lastState != state {
				if err := conn.WriteMessage(websocket.TextMessage, []byte("state:idle")); err != nil {
					log.Println(err)
					return
				}
			}
		}
	}
}

func stateTrashHandler(conn *websocket.Conn) {
	log.Println("handler2 runnning")
	for {
		switch trashState {
		case 2:
			trashState = 999
			if trashlastState != trashState {
				if err := conn.WriteMessage(websocket.TextMessage, []byte("state:awarding")); err != nil {
					log.Println(err)
					return
				}
			}
		case 3:
			trashState = 999
			if trashlastState != trashState {
				if err := conn.WriteMessage(websocket.TextMessage, []byte("state:awarding")); err != nil {
					log.Println(err)
					return
				}
			}
		case 4:
			trashState = 999
			if trashlastState != trashState {
				if err := conn.WriteMessage(websocket.TextMessage, []byte("state:awarding")); err != nil {
					log.Println(err)
					return
				}
			}
		}
	}
}

func iotHandler(w http.ResponseWriter, r *http.Request) {
	enableCors(&w)
	var e Event

	temp, _ := io.ReadAll(r.Body)
	e.event, _ = strconv.Atoi(string(temp))

	if e.event == 1 {
		w.Write([]byte("1"))
		state = 1
	} else {
		w.Write([]byte("0"))
		state = 0
	}
	lastState = e.event

	log.Println(trashState)
	log.Println(trashlastState)
	log.Println(lastState)
	log.Println(state)
}

func iot2Handler(w http.ResponseWriter, r *http.Request) {
	enableCors(&w)
	var e Event

	temp, _ := io.ReadAll(r.Body)
	e.event, _ = strconv.Atoi(string(temp))

	if lastState == 1 {
		switch e.event {
		case 2:
			trashState = 2
		case 3:
			trashState = 3
		case 4:
			trashState = 4
		}
		trashlastState = e.event
	}
	log.Println(trashState)
	log.Println(trashlastState)
	log.Println(lastState)
	log.Println(state)
}

func classify_page(w http.ResponseWriter, r *http.Request) {
	enableCors(&w)

	t := TrashObject{}
	var labelDesc []string //untuk

	webcam, err := gocv.OpenVideoCapture(0)
	if err != nil {
		fmt.Println("Error opening webcam:", err)
		return
	}
	defer webcam.Close()

	img := gocv.NewMat()
	defer img.Close()

	time.Sleep(5 * time.Second)
	if ok := webcam.Read(&img); !ok {
		fmt.Println("Error reading frame from webcam")
	}
	filename := saveFrame(img)

	labelDesc = google_vision(filename, labelDesc)

	t.DetectedAt = time.Now().Local()
	// t.Type = 1

	tJson, _ := json.Marshal(t)
	labelJson, _ := json.Marshal(labelDesc)

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(tJson)
	w.Write(labelJson)
	// for _, label := range labelDesc {
	// 	w.Write([]byte(label))
	// }
}

func saveFrame(frame gocv.Mat) string {
	filename := filepath.Join("../Frontend/src/assets", "frame_image.jpg")
	// filename := filepath.Join("../Frontend/src/assets", fmt.Sprintf("frame_%s.jpg", time.Now().Format("20060102150405")))

	// Write the frame to the file
	if ok := gocv.IMWrite(filename, frame); !ok {
		fmt.Println("Error writing frame to file:", filename)
	} else {
		fmt.Println("Frame saved:", filename)
	}
	return filename
}

func google_vision(filename string, labelDesc []string) []string {
	// Timeout harus ada

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
		fmt.Println(label.Description)
		fmt.Println(label.Score)
	}
	return labelDesc
}

func addTrash(w http.ResponseWriter, r *http.Request) {
	// untuk klasifikasi sampah masuk di tong mana
	//-> nanti buat next ke thankyou
	//-> update ke
}
