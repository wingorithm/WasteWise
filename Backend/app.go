package main

import (
	"context"
	"encoding/base64"
	"image"
	"image/jpeg"
	"io"
	"log"
	"math"
	"net/http"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"time"

	vision "cloud.google.com/go/vision/apiv1"
	"github.com/gorilla/mux"
	"github.com/gorilla/websocket"
)

func NewApp() *App {
	return &App{
		r: mux.NewRouter(),
	}
}

type App struct {
	r *mux.Router
}

type Event struct {
	event int `json:"event"`
}

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

var recycleWhiteList = map[string]int{
	"recycleable":            5,
	"Plastic":                5,
	"Bottled water":          2,
	"Plastic bottle":         3,
	"Water bottle":           1,
	"Drinkware":              1,
	"Bottle":                 4,
	"Packing materials":      3,
	"Paper bag":              3,
	"Box":                    5,
	"Carton":                 5,
	"Cardboard":              5,
	"Package":                3,
	"Shipping box":           3,
	"delivery":               2,
	"Packaging":              2,
	"Tin,":                   2,
	"Beverage can":           5,
	"Aluminum can,":          5,
	"Tin can":                5,
	"Soft drink":             2,
	"Carbonated soft drinks": 1,
	"Metal":                  2,
	"Nickel":                 1,
	"Laboratory":             1,
	"Mineral water":          2,
	"Packaging and labeling": -2,
	"Plastic bag":            -2,
	"Logo":                   -2,
	"Brand":                  -2,
	"Snack":                  -1,
}
var organikWhiteList = map[string]int{
	"Vegetable":              3,
	"Food":                   3,
	"Fruit":                  3,
	"Plant":                  1,
	"Natural foods":          1,
	"Junk food":              1,
	"Gluten":                 1,
	"Staple food":            1,
	"Ingredient":             2,
	"Peel":                   1,
	"Banana":                 1,
	"Dish":                   1,
	"Soil":                   1,
	"Herb":                   1,
	"Packaging and labeling": -2,
	"Plastic bag":            -2,
	"Logo":                   -2,
	"Brand":                  -2,
	"Snack":                  -1,
}
var blackList = map[string]int{
	"Nose":       1,
	"Hair":       1,
	"Face":       1,
	"Skin":       1,
	"Head":       1,
	"Hand":       1,
	"Eye":        1,
	"Lip":        1,
	"Mouth":      1,
	"Photograph": 1,
	"Happy":      1,
	"Finger":     1,
	"Gesture":    1,
	"Organ":      1,
	"Joint":      1,
	"Shoulder":   1,
	"finger":     1,
	"Thumb":      1,
	"Jaw":        1,
}
var indexSlices = map[int]int{
	1: 1, //Recycle
	2: 1, //Organic
	3: 1, //Lainya
}

// State for IoT
var state = 999
var lastState = 999
var trashState = 0
var trashlastState = 0

// State For Classify
var classState = 0
var lastclassState = 1
var counter = 0

var maxObjectRecycle = 0
var macObjectOrganik = 0
var macObjectlainya = 0

func (a *App) configure_routes() {
	a.r.HandleFunc("/ws", websocketHandle)
	a.r.HandleFunc("/iot", iotHandler).Methods("POST")
	a.r.HandleFunc("/iot2", iot2Handler).Methods("POST")
	a.r.HandleFunc("/imageHandler", classify_page).Methods("POST")
}

func websocketHandle(w http.ResponseWriter, r *http.Request) {
	upgrader.CheckOrigin = func(r *http.Request) bool { return true }

	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Println(err)
	}

	log.Println("Client connected")
	go stateHandler(conn)
	go stateTrashHandler(conn)
	go stateClassHandler(conn)
}

func stateHandler(conn *websocket.Conn) {
	log.Println("handler1 runnning")
	for {
		switch state {
		case 1:
			state = 999
			if lastState != state {
				if err := conn.WriteMessage(websocket.TextMessage, []byte("intro")); err != nil {
					log.Println(err)
					return
				}
			}
			if classState == 0 {
				time.Sleep(8 * time.Second)
				// classify_page(conn)
				classState = 1 //-> siap terima file jolee
			}
		case 0:
			state = 999
			if lastState != state {
				if err := conn.WriteMessage(websocket.TextMessage, []byte("idle")); err != nil {
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
				if err := conn.WriteMessage(websocket.TextMessage, []byte("reward")); err != nil {
					log.Println(err)
					return
				}
			}
		case 3:
			trashState = 999
			if trashlastState != trashState {
				if err := conn.WriteMessage(websocket.TextMessage, []byte("reward")); err != nil {
					log.Println(err)
					return
				}
			}
		case 4:
			trashState = 999
			if trashlastState != trashState {
				if err := conn.WriteMessage(websocket.TextMessage, []byte("reward")); err != nil {
					log.Println(err)
					return
				}
			}
		}
	}
}

func iotHandler(w http.ResponseWriter, r *http.Request) {
	var e Event
	timeout := 0 * time.Second

	temp, _ := io.ReadAll(r.Body)
	e.event, _ = strconv.Atoi(string(temp))

	w.WriteHeader(http.StatusOK)
	select {
	case <-time.After(timeout):
		if e.event == 1 {
			w.Write([]byte("1"))
			state = 1
			timeout = 14 * time.Second
		} else {
			w.Write([]byte("0"))
			state = 0
		}
	}
	lastState = e.event

	log.Println(trashState)
	log.Println(trashlastState)
	log.Println(lastState)
	log.Println(state)
}

func iot2Handler(w http.ResponseWriter, r *http.Request) {
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
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)

		time.Sleep(10 * time.Second)
		state = 0
	}
	log.Println(trashState)
	log.Println(trashlastState)
	log.Println(lastState)
	log.Println(state)
}

func stateClassHandler(conn *websocket.Conn) {
	for {
		if counter == 5 {
			log.Println(macObjectlainya)
			log.Println(maxObjectRecycle)
			log.Println(macObjectOrganik)
			res := math.Max(float64(macObjectlainya), math.Max(float64(maxObjectRecycle), float64(macObjectOrganik)))
			time.Sleep(2 * time.Second)

			if res == float64(maxObjectRecycle) {
				if err := conn.WriteMessage(websocket.TextMessage, []byte("info:1")); err != nil {
					log.Println(err)
					return
				}
				maxObjectRecycle = 0
				classState = 0
			} else if res == float64(macObjectOrganik) {
				if err := conn.WriteMessage(websocket.TextMessage, []byte("info:2")); err != nil {
					log.Println(err)
					return
				}
				macObjectOrganik = 0
				classState = 0
			} else {
				if err := conn.WriteMessage(websocket.TextMessage, []byte("info:3")); err != nil {
					log.Println(err)
					return
				}
				macObjectlainya = 0
				classState = 0
			}
			counter = 0
		}
	}
}

func classify_page(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}
	// t := TrashObject{}
	var labelDesc []string //untuk

	body, err := io.ReadAll(r.Body)
	if err != nil {
		http.Error(w, "Error reading request body", http.StatusInternalServerError)
		return
	}

	filename := DecodeBase64Image(string(body))

	if classState == 1 {
		// t = google_vision(filename, labelDesc)
		google_vision(filename, labelDesc)
		counter += 1
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
	}
}

func DecodeBase64Image(body string) string {
	// Decode base64-encoded image
	decodedImage, err := base64.StdEncoding.DecodeString(body)
	if err != nil {
		log.Panic("Error decoding base64:")
	}

	// Create an image.Image from the decoded data
	img, _, err := image.Decode(strings.NewReader(string(decodedImage)))
	if err != nil {
		log.Panic("Error decoding image:")
	}
	filePath := saveFrame(img)
	return filePath
}

func saveFrame(img image.Image) string { //file baru setelah ganti orang
	saveFilePath := "../Frontend/src/assets/sample"

	outputFilePath := filepath.Join(saveFilePath, "decoded_image.jpg")
	outputFile, err := os.Create(outputFilePath)
	if err != nil {
		log.Panic("Error creating output file:")
	}
	defer outputFile.Close()

	err = jpeg.Encode(outputFile, img, nil)
	if err != nil {
		log.Panic("Error encoding image:", err)
	}

	return outputFilePath
}

func google_vision(filename string, labelDesc []string) TrashObject {
	// Timeout harus ada
	// INGET REGION CONFIDENCE
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

	labels, err := client.DetectLabels(ctx, image, nil, 20)
	if err != nil {
		log.Fatalf("Failed to detect labels: %v", err)
	}

	for _, label := range labels {

		if blackList[string(label.Description)] != 1 {
			// log.Println(">>", label.Description)
			// log.Println(label.Score)
			labelDesc = append(labelDesc, label.Description)
		}
	}

	return calculateScore(labelDesc)
}

func calculateScore(labelDesc []string) TrashObject { // {2:Organic , 1: recycleable, 2: lainnya}
	to := TrashObject{}
	var maxRecycle, maxOrganic int = 0, 0

	for _, v := range labelDesc {
		maxRecycle += recycleWhiteList[v]
		maxOrganic += organikWhiteList[v]
	}

	res := math.Max(float64(maxRecycle), float64(maxOrganic))

	log.Println(labelDesc)
	log.Println("==================================")
	log.Println(maxRecycle)
	log.Println(maxOrganic)
	if res == float64(maxRecycle) && res > float64(8) {
		to.event = 1
		to.DetectedAt = time.Now().UTC()
		to.Name = "Recycle"
		maxObjectRecycle += 1
	} else if res == float64(maxOrganic) && res > float64(8) {
		to.event = 2
		to.DetectedAt = time.Now().UTC()
		to.Name = "Organic"
		macObjectOrganik += 1
	} else {
		to.event = 3
		to.DetectedAt = time.Now().UTC()
		to.Name = "Lainya"
		macObjectlainya += 1
	}
	return to
}

func addTrash(w http.ResponseWriter, r *http.Request) {
	// untuk klasifikasi sampah masuk di tong mana
	//-> nanti buat next ke thankyou
	//-> update ke
}
