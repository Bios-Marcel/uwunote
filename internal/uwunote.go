package internal

import (
	"io/ioutil"
	"log"
	"os"
	"runtime"

	"github.com/gotk3/gotk3/glib"

	"github.com/getlantern/systray"
	"github.com/gotk3/gotk3/gtk"

	"github.com/UwUNote/uwunote/internal/config"
)

//Start initializes gtk and creates a window for every note.
func Start() {
	config.CreateNeccessaryFiles()
	config.LoadAppConfig()

	os.MkdirAll(notePath, os.ModePerm)

	if config.GetAppConfig().ShowTrayIcon {
		systray.Run(
			func() {
				systemTrayRun()

				//Only on linux (and some others) gtk will be used, therefore we needed init gtk on linux.
				if runtime.GOOS == "windows" || runtime.GOOS == "darwin" {
					startAndInitGtk()
				} else {
					start()
				}
			},
			func() {})
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

func start() {
	glib.IdleAdd(generateNoteWindows)
}

func startAndInitGtk() {
	// Initialize GTK without parsing any command line arguments.
	gtk.Init(nil)

	start()

	// Begin executing the GTK main loop. This blocks until gtk.MainQuit() is run.
	gtk.Main()
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
		CreateNoteGUI(0, 0, 300, 350, nil)
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
