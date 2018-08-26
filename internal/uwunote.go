package internal

import (
	"io/ioutil"
	"log"
	"os"
	"path/filepath"

	"github.com/getlantern/systray"
	"github.com/gotk3/gotk3/glib"
	"github.com/gotk3/gotk3/gtk"

	"github.com/UwUNote/uwunote/internal/config"
	"github.com/UwUNote/uwunote/internal/util"
)

//Start initializes gtk and creates a window for every note.
func Start() {
	config.CreateNeccessaryFiles()
	util.LogAndExitOnError(config.LoadAppConfig())

	os.MkdirAll(notePath, os.ModePerm)

	if config.GetAppConfig().ShowTrayIcon {
		startWithTrayIcon()
	} else {
		startAndInitGtk()
	}
}

func systemTrayRun() {
	systray.SetIcon(Icon)
	newNoteItem := systray.AddMenuItem("New note", "Creates a new note")
	settingsItem := systray.AddMenuItem("Settings", "Opens the settings")
	systray.AddSeparator()
	quitItem := systray.AddMenuItem("Quit", "Closes the application")

	go func() {
		for {
			select {
			case <-newNoteItem.ClickedCh:
				glib.IdleAdd(CreateNoteGUIWithDefaults)

			case <-settingsItem.ClickedCh:
				config.OpenAppConfig()

			case <-quitItem.ClickedCh:
				glib.IdleAdd(func() {
					gtk.MainQuit()
					os.Exit(0)
				})
			}
		}
	}()
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
	files, err := ioutil.ReadDir(notePath)

	if err != nil {
		log.Fatal("Error reading notes.")
		panic(err)
	}

	config.LoadWindowConfiguration()
	appConfig := config.GetAppConfig()

	if len(files) == 0 {
		log.Println("Generating a initial note.")
		CreateNoteGUIWithDefaults()
	} else {
		log.Println("Creating windows for existing notes.")
		for _, fileInfo := range files {
			if fileInfo.IsDir() {
				continue
			}

			fileName := fileInfo.Name()
			configForWindow, exists := config.GetWindowDataForFile(fileName)

			pathToNote := filepath.Join(notePath, fileName)
			if exists {
				createWindowForNote(pathToNote, configForWindow.X, configForWindow.Y, configForWindow.Width, configForWindow.Height)
			} else {
				createWindowForNote(pathToNote, appConfig.DefaultNoteX, appConfig.DefaultNoteY, appConfig.DefaultNoteWidth, appConfig.DefaultNoteHeight)
			}
		}
	}
}
