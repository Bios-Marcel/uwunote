package internal

import (
	"io/ioutil"
	"log"
	"os"
	"path/filepath"

	"github.com/gotk3/gotk3/gtk"

	"github.com/Bios-Marcel/uwuNote/internal/config"
)

var (
	//Will be customizable at some point
	notePath = filepath.FromSlash(os.Getenv("HOME") + string(os.PathSeparator) + "notes")
)

//Start initializes gtk and creates a window for every note.
func Start() {
	// Initialize GTK without parsing any command line arguments.
	gtk.Init(nil)

	createNeccessaryDirectories()

	config.LoadAppConfig()

	generateNoteWindows()

	// Begin executing the GTK main loop. This blocks until gtk.MainQuit() is run.
	gtk.Main()
}

func createNeccessaryDirectories() {
	os.MkdirAll(notePath, os.ModePerm)
	config.CreateNeccessaryFiles()
}

//Creates a window for every node inside of the notePath
func generateNoteWindows() {
	files, err := ioutil.ReadDir(notePath)

	if err != nil {
		log.Fatal("Error reading notes.")
		panic(err)
	}

	config.LoadWindowConfiguration()

	if len(files) == 0 {
		log.Println("Generating a initial note.")
		CreateNote(0, 0, 300, 350)
	} else {
		log.Println("Creating windows for existing notes.")
		for _, fileInfo := range files {
			if fileInfo.IsDir() {
				continue
			}

			fileName := fileInfo.Name()
			configForWindow, exists := config.GetWindowDataForFile(fileName)

			pathToNote := notePath + string(os.PathSeparator) + fileName
			if exists {
				createWindowForNote(pathToNote, configForWindow.X, configForWindow.Y, configForWindow.Width, configForWindow.Height)
			} else {
				createWindowForNote(pathToNote, 0, 0, 300, 350)
			}
		}
	}
}
