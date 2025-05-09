package main

import (
	"flag"
	"fmt"
	"net/http"
	"os"
)

type Background struct {
	filename string
	name     string
}

var (
	fDataDir    string
	fReloadFile string
	bgs         []Background
	current     Background
)

func init() {
	flag.StringVar(&fDataDir, "datadir", "datadir", "directory with image files")
	flag.StringVar(&fReloadFile, "reloadfile", "reloadplease", "file to touch when kioks should reload")
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

func CreateReloadFile() {
	if fReloadFile == "" {
		return
	}
	if _, err := os.Stat(fReloadFile); os.IsNotExist(err) {
		fmt.Printf("Creating reload file %v\n", fReloadFile)
		f, err := os.Create(fReloadFile)
		if err != nil {
			fmt.Printf("Error creating reload file %v: %v", fReloadFile, err)
			return
		}
		f.Close()
	} else {
		fmt.Printf("Reload file %v already exists\n", fReloadFile)
	}
}

// HandleRoot handles the root URL.
func HandleRoot(w http.ResponseWriter, r *http.Request) {
	content := fmt.Sprintf("<html><body><style> body { background-image: url('/b?name=%v'); background-repeat: no-repeat; background-size: contain; background-position: center; background-color: black; }</style></body></html>", current.name)
	fmt.Fprintf(w, content)
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
				current = bg
				fmt.Println("Set current background to", bg.name)
				CreateReloadFile()
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

	http.HandleFunc("/a", HandleAdmin)
	http.HandleFunc("/b", HandleBackground)
	http.HandleFunc("/", HandleRoot)

	http.ListenAndServe(":8080", nil)
}
