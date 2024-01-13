package main

import (
	"context"
	"fmt"
	"io"
	"log"
	"math"
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

// Scoring 1-3
var recycleWhiteList = map[string]int{ // plastic bottle || can || cardboard
	"recycleable":            3,
	"Bottled water":          1,
	"Plastic bottle":         3,
	"Water bottle":           1,
	"Drinkware":              1,
	"Bottle":                 2,
	"Packing materials":      1,
	"Paper bag":              1,
	"Box":                    1,
	"Carton":                 2,
	"Cardboard":              3,
	"Tin,":                   1,
	"Beverage can":           3,
	"Aluminum can,":          3,
	"Tin can":                3,
	"Soft drink":             1,
	"Carbonated soft drinks": 1,
	"Metal":                  2,
	"Nickel":                 1,
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
}

// nodemcu -> 1 (ada orang) / 0 (gk ada orang)
var state = 999
var lastState = 999

// nodemcu terima sampah -> 2:organik / 3:recycle / 4:another
var trashState = 0
var trashlastState = 0

// State For Classify
var classState = 0

func (a *App) configure_routes() {
	a.r.HandleFunc("/ws", websocketHandle)
	a.r.HandleFunc("/iot", iotHandler).Methods("POST")   //human
	a.r.HandleFunc("/iot2", iot2Handler).Methods("POST") //barang masuk

	// a.r.HandleFunc("/classify", classify_page) // nanti apus
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
			log.Println("state")
			if lastState != state {
				if err := conn.WriteMessage(websocket.TextMessage, []byte("intro")); err != nil {
					log.Println(err)
					return
				}
			}
			if classState == 0 {
				log.Println("mauuuu mulai camera")
				time.Sleep(8 * time.Second)
				classify_page(conn)
				classState = 1
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
	enableCors(&w)
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
		classState = 0
		w.WriteHeader(http.StatusOK)
		time.Sleep(10 * time.Second)
		state = 0
	}
	log.Println(trashState)
	log.Println(trashlastState)
	log.Println(lastState)
	log.Println(state)
}

func classify_page(conn *websocket.Conn) {
	// tSlice := []TrashObject{}
	t := TrashObject{}
	var labelDesc []string //untuk

	// devices, err := gocv.Enume()
	// if err != nil {
	// 	fmt.Println("Error enumerating video capture devices:", err)
	// 	return
	// }

	// // Print information about available cameras
	// fmt.Println("Available video capture devices:")
	// for i, device := range devices {
	// 	fmt.Printf("%d: %s\n", i, device)
	// }

	log.Println("mulai camera")
	webcam, err := gocv.OpenVideoCapture(1)
	// 2 -> obs
	if err != nil {
		fmt.Println("Error opening webcam:", err)
		return
	}
	log.Println("mulai camera berhasil")

	img := gocv.NewMat()
	defer img.Close()
	maxObjectRecycle := 0
	macObjectOrganik := 0
	macObjectlainya := 0

	if classState == 0 {
		log.Println("masuk camera1")
		for { //take 5 image
			log.Println("masuk camera")
			time.Sleep(1 * time.Second)
			if ok := webcam.Read(&img); !ok {
				fmt.Println("Error reading frame from webcam")
			}
			filename := saveFrame(img)
			t = google_vision(filename, labelDesc)
			// t = google_vision("./sample image/snacl.jpg", labelDesc)

			if t.event == 1 {
				maxObjectRecycle += 1
			} else if t.event == 2 {
				macObjectOrganik += 1
			} else {
				macObjectlainya += 1
			}

			if (maxObjectRecycle + macObjectOrganik + macObjectlainya) == 5 {
				break
			}
		}
		webcam.Close()
	}

	res := math.Max(float64(macObjectlainya), math.Max(float64(maxObjectRecycle), float64(macObjectOrganik)))
	log.Println(t.event)
	log.Println(t.Name)
	time.Sleep(2 * time.Second)

	if res == float64(maxObjectRecycle) {
		if err := conn.WriteMessage(websocket.TextMessage, []byte("info:1")); err != nil {
			log.Println(err)
			return
		}
		classState = 0
	} else if res == float64(macObjectOrganik) {
		if err := conn.WriteMessage(websocket.TextMessage, []byte("info:2")); err != nil {
			log.Println(err)
			return
		}
		classState = 0
	} else {
		if err := conn.WriteMessage(websocket.TextMessage, []byte("info:3")); err != nil {
			log.Println(err)
			return
		}
		classState = 0
	}
}

func saveFrame(frame gocv.Mat) string { //file baru setelah ganti orang
	log.Println("masuk camera")
	filename := filepath.Join("../Frontend/src/assets/sample", "frame_image.jpg")
	// filename := filepath.Join("../Frontend/src/assets", fmt.Sprintf("frame_%s.jpg", time.Now().Format("20060102150405")))

	// Write the frame to the file
	if ok := gocv.IMWrite(filename, frame); !ok {
		fmt.Println("Error writing frame to file:", filename)
	} else {
		fmt.Println("Frame saved:", filename)
	}
	return filename
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
			labelDesc = append(labelDesc, label.Description)
		}
		fmt.Println(label.Description)
		fmt.Println(label.Score)
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

	len := len(labelDesc)
	res := math.Max(float64(maxRecycle), float64(maxOrganic))

	if res == float64(maxRecycle) && res > float64(len) {
		to.event = 1
		to.DetectedAt = time.Now().UTC()
		to.Name = "Recycle"
	} else if res == float64(maxOrganic) && res > float64(len) {
		to.event = 2
		to.DetectedAt = time.Now().UTC()
		to.Name = "Organic"
	} else {
		to.event = 3
		to.DetectedAt = time.Now().UTC()
		to.Name = "Lainya"
	}
	return to
}

func addTrash(w http.ResponseWriter, r *http.Request) {
	// untuk klasifikasi sampah masuk di tong mana
	//-> nanti buat next ke thankyou
	//-> update ke
}
