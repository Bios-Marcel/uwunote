package main

import (
	"log"

	"os"
	"path/filepath"

	"github.com/gotk3/gotk3/gtk"
)

var (
	//Will lateron be customizable
	notePath = filepath.FromSlash(os.Getenv("HOME") + string(os.PathSeparator) + "notes")
)

func main() {
	// Initialize GTK without parsing any command line arguments.
	gtk.Init(nil)

	generateNoteWindows()

	// Begin executing the GTK main loop.  This blocks until
	// gtk.MainQuit() is run.
	gtk.Main()
}

func generateNoteWindows() {
	var files []string

	err := filepath.Walk(notePath, func(path string, info os.FileInfo, err error) error {
		//We wil ignore all notes that lie on second level
		if info != nil && info.IsDir() {
			return nil
		}

		files = append(files, path)
		return nil
	})

	if err != nil {
		log.Fatal("Error creating notes.")
		panic(err)
	}

	for index, file := range files {
		pos := index * 310
		createWindowForNote(file, pos, 0)
	}
}

// Panics if the given value isn't nil
func panicOnError(possibleError error) {
	if possibleError != nil {
		log.Fatal(possibleError)
		panic(possibleError)
	}
}
