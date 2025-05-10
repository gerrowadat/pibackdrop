package main

import (
	_ "embed"
	"flag"
	"fmt"
	"net/http"
	"os"
	"sync"
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
	// mutex for current.
	cm sync.Mutex
	//go:embed index.html
	indexHtml []byte
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
	w.Write(indexHtml)
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
				cm.Lock()
				defer cm.Unlock()
				current = bg
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

func HandleCurrent(w http.ResponseWriter, r *http.Request) {
	cm.Lock()
	defer cm.Unlock()
	fmt.Fprint(w, current.name)
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

	if len(bgs) > 0 {
		current = bgs[0]
	}

	// The 'admin' page is for setting the background
	http.HandleFunc("/a", HandleAdmin)
	// The 'b' endpoint is for serving the background images
	http.HandleFunc("/b", HandleBackground)
	// Just sends the name of the current background.
	http.HandleFunc("/c", HandleCurrent)
	// The root endpoint serves the HTML page with the current background
	http.HandleFunc("/", HandleRoot)

	fmt.Printf("Listening on port %d\n", fHttpPort)

	http.ListenAndServe(fmt.Sprintf(":%d", fHttpPort), nil)
}
