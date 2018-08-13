package main

import (
	"encoding/json"
	"io"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"

	"github.com/gotk3/gotk3/gtk"
)

var (
	//Will lateron be customizable
	notePath         = filepath.FromSlash(os.Getenv("HOME") + string(os.PathSeparator) + "notes")
	configPath       = notePath + string(os.PathSeparator) + "config"
	windowConfigPath = configPath + string(os.PathSeparator) + "windows.json"
	//Configuration contains the positions and sizes for all notes
	Configuration = WindowConfig{}
)

//WindowConfig contains a map of WindowData entries in a map
type WindowConfig struct {
	Data map[string]WindowData
}

//WindowData contains the sosition (x and y) and size of a Window
type WindowData struct {
	X int
	Y int

	Width  int
	Height int
}

func main() {
	// Initialize GTK without parsing any command line arguments.
	gtk.Init(nil)

	generateNoteWindows()

	// Begin executing the GTK main loop.  This blocks until
	// gtk.MainQuit() is run.
	gtk.Main()
}

func generateNoteWindows() {
	files, err := ioutil.ReadDir(notePath)

	if err != nil {
		log.Fatal("Error reading notes.")
		panic(err)
	}

	log.Println("Get window configuration")
	file, openError := os.Open(windowConfigPath)
	if openError == nil || os.IsExist(openError) {
		defer file.Close()
		decoder := json.NewDecoder(file)
		windowConfigLoadError := decoder.Decode(&Configuration)
		if windowConfigLoadError != io.EOF {
			logAndExit(windowConfigLoadError)
		}
	}

	if Configuration.Data == nil {
		//Creating an empty map to prevent nil pointer references
		Configuration.Data = make(map[string]WindowData)
	}

	log.Println("Create note windows")
	for _, fileInfo := range files {
		if fileInfo.IsDir() {
			continue
		}

		file := fileInfo.Name()
		configForWindow, exists := Configuration.Data[file]

		if exists {
			createWindowForNote(notePath+string(os.PathSeparator)+file, configForWindow.X, configForWindow.Y, configForWindow.Width, configForWindow.Height)
		} else {
			createWindowForNote(notePath+string(os.PathSeparator)+file, 0, 0, 300, 350)
		}
	}
}

func persistWindowConfiguration() {
	windowConfigurationJSON, _ := json.Marshal(&Configuration)
	writeError := ioutil.WriteFile(windowConfigPath, windowConfigurationJSON, os.ModePerm)
	//TODO Better way?
	logAndExit(writeError)
}

func logAndExit(possibleError error) {
	if possibleError != nil {
		log.Fatal(possibleError)
	}
}
