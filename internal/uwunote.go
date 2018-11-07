package internal

import (
	"flag"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"

	"github.com/UwUNote/uwunote/internal/errors"

	"github.com/gotk3/gotk3/glib"
	"github.com/gotk3/gotk3/gtk"

	"github.com/UwUNote/uwunote/internal/config"
	"github.com/UwUNote/uwunote/internal/gui"
	"github.com/UwUNote/uwunote/internal/util"
)

//Start initializes gtk and creates a window for every note.
func Start() {
	configpathPointer := flag.String("configdir", filepath.Join(util.HomeDir, ".uwunote"), "configdir <path to folder>")
	flag.Parse()

	config.ConfigPath = *configpathPointer

	config.CreateNeccessaryFiles()
	errors.ShowErrorDialogOnError(config.LoadAppConfig())

	os.MkdirAll(config.GetAppConfig().NoteDirectory, 0755)

	if config.GetAppConfig().ShowTrayIcon {
		startWithTrayIcon()
	} else {
		startAndInitGtk()
	}
}

func systemTrayRun() {
	buildSystray()
}

//startAndInitGtk initializes gtk, invokes `start` and triggers the gtk mainloop.
func startAndInitGtk() {
	// Initialize GTK without parsing any command line arguments.
	gtk.Init(nil)

	start()

	// Begin executing the GTK main loop. This blocks until gtk.MainQuit() is run.
	gtk.Main()
}

func start() {
	glib.IdleAdd(generateNoteWindows)
}

//Creates a window for every note inside of the notePath
func generateNoteWindows() {
	files, errorReadingNotes := ioutil.ReadDir(config.GetAppConfig().NoteDirectory)

	if errorReadingNotes != nil {
		log.Fatalf("Error reading notes (%s).", errorReadingNotes.Error())
	}

	config.LoadWindowConfiguration()
	appConfig := config.GetAppConfig()

	if len(files) == 0 {
		//Creating a single note, as there are no existing notes yet.
		gui.CreateNoteGUIWithDefaults()
		return
	}

	for _, fileInfo := range files {
		if fileInfo.IsDir() {
			continue
		}

		fileName := fileInfo.Name()
		configForNoteWindow, noteWindowConfigExists := config.GetWindowDataForFile(fileName)

		pathToNote := filepath.Join(config.GetAppConfig().NoteDirectory, fileName)
		if noteWindowConfigExists {
			gui.CreateWindowForNote(pathToNote, configForNoteWindow.X, configForNoteWindow.Y, configForNoteWindow.Width, configForNoteWindow.Height)
		} else {
			gui.CreateWindowForNote(pathToNote, appConfig.DefaultNoteX, appConfig.DefaultNoteY, appConfig.DefaultNoteWidth, appConfig.DefaultNoteHeight)
		}
	}
}
