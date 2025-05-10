package main

import (
	"flag"
	"fmt"
	"net/http"
	"os"

	"github.com/gorilla/websocket"
)

type Background struct {
	filename string
	name     string
}

var (
	fDataDir  string
	fHttpPort int
	bgs       []Background
	// The current background
	current Background
	// A channel for background updates
	cc chan string
)

func init() {
	flag.StringVar(&fDataDir, "datadir", "datadir", "directory with image files")
	flag.IntVar(&fHttpPort, "port", 8080, "HTTP port to listen on")
}

// Backgrounds returns a slice of Background structs.
func Backgrounds(datadir string) []Background {
	files, err := os.ReadDir(datadir)
	if err != nil {
		fmt.Printf("Error reading datadir %v: %v", datadir, err)
		return nil
	}
	bgs = make([]Background, 0, len(files))
	for _, file := range files {
		if file.IsDir() {
			continue
		}
		name := file.Name()
		bgs = append(bgs, Background{
			filename: datadir + "/" + name,
			name:     name[:len(name)-4],
		})
	}
	return bgs
}

// HandleRoot handles the root URL.
func HandleRoot(w http.ResponseWriter, r *http.Request) {
	content := fmt.Sprintf("<html><body><style> body { background-image: url('/b?name=%v'); background-repeat: no-repeat; background-size: contain; background-position: center; background-color: black; }</style></body></html>", current.name)
	fmt.Fprint(w, content)
}

// HandleAdmin handles the admin URL.
func HandleAdmin(w http.ResponseWriter, r *http.Request) {
	action := r.URL.Query().Get("action")
	if action == "set" {
		name := r.URL.Query().Get("name")
		if name == "" {
			w.WriteHeader(http.StatusBadRequest)
			return
		}
		for _, bg := range bgs {
			if bg.name == name {
				UpdateCurrentBackground(bg)
				break
			}
		}
	}
	for _, bg := range bgs {
		line := fmt.Sprintf("<a href=\"/b?name=%s\">%s</a> - <a href=\"/a?action=set&name=%s\">[set]</a><br/>", bg.name, bg.name, bg.name)
		fmt.Fprintln(w, line)
	}
}

// HandleBackground handles the background URL.
func HandleBackground(w http.ResponseWriter, r *http.Request) {
	name := r.URL.Query().Get("name")
	if name == "" {
		w.WriteHeader(http.StatusNotFound)
		return
	}
	for _, bg := range bgs {
		if bg.name == name {
			http.ServeFile(w, r, bg.filename)
			return
		}
	}
}

var upgrader = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
}

func HandleWebSocketCurrent(w http.ResponseWriter, r *http.Request) {
	fmt.Printf("WebSocket connection from %v\n", r.RemoteAddr)
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		fmt.Println("Error upgrading connection:", err)
		return
	}
	defer conn.Close()

	err = conn.WriteJSON(current.name)
	if err != nil {
		fmt.Println("Error writing JSON:", err)
		return
	}

	for {
		fmt.Println("Waiting for background update...")
		curr := <-cc
		err = conn.WriteJSON(curr)
		if err != nil {
			fmt.Println("Error writing JSON:", err)
			return
		}
	}
}

func UpdateCurrentBackground(newbg Background) {
	select {
	case cc <- newbg.name:
		fmt.Printf("Background updated to %v\n", newbg.name)
		current = newbg
	default:
		fmt.Println("Background update channel is full.")
		return
	}
}

func main() {
	flag.Parse()
	fmt.Printf("Reading files from %v\n", fDataDir)
	bgs := Backgrounds(fDataDir)
	if bgs == nil {
		fmt.Println("No backgrounds found")
		return
	}
	for _, bg := range bgs {
		fmt.Printf("Background: %v\n", bg.name)
	}
	fmt.Printf("Total backgrounds: %v\n", len(bgs))

	current = bgs[0]

	cc = make(chan string, 1)

	// The 'admin' page is for setting the background
	http.HandleFunc("/a", HandleAdmin)
	// The 'b' endpoint is for serving the background images
	http.HandleFunc("/b", HandleBackground)
	// The 'current' endpoint is for WebSocket connections
	// It sends the current background name to the client as it is updated.
	http.HandleFunc("/current", HandleWebSocketCurrent)
	// The root endpoint serves the HTML page with the current background
	http.HandleFunc("/", HandleRoot)

	fmt.Printf("Listening on port %d\n", fHttpPort)

	http.ListenAndServe(fmt.Sprintf(":%d", fHttpPort), nil)
}
