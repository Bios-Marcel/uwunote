package internal

import (
	"io/ioutil"
	"log"
	"os"

	"github.com/gotk3/gotk3/gtk"

	"github.com/Bios-Marcel/uwuNote/internal/config"
)

//Start initializes gtk and creates a window for every note.
func Start() {
	// Initialize GTK without parsing any command line arguments.
	gtk.Init(nil)

	startInternal()

	// Begin executing the GTK main loop. This blocks until gtk.MainQuit() is run.
	gtk.Main()
}

func startInternal() {
	createNeccessaryDirectories()

	config.LoadAppConfig()

	generateNoteWindows()
}

func createNeccessaryDirectories() {
	os.MkdirAll(notePath, os.ModePerm)
	config.CreateNeccessaryFiles()
}

//Creates a window for every note inside of the notePath
func generateNoteWindows() {
	files, err := ioutil.ReadDir(notePath)

	if err != nil {
		log.Fatal("Error reading notes.")
		panic(err)
	}

	config.LoadWindowConfiguration()

	if len(files) == 0 {
		log.Println("Generating a initial note.")
		CreateNoteGUI(0, 0, 300, 350)
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
