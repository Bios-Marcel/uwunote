package internal

import (
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"

	"github.com/skratchdot/open-golang/open"

	"github.com/UwUNote/uwunote/internal/errors"
	"github.com/UwUNote/uwunote/internal/updates"

	"github.com/getlantern/systray"
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
	systray.SetIcon(gui.AppIcon)
	newNoteItem := systray.AddMenuItem("New note", "Creates a new note")
	systray.AddSeparator()
	settingsItem := systray.AddMenuItem("Settings", "Opens the settings")
	shortcutsItem := systray.AddMenuItem("Shortcuts", "Opens the shortcuts dialog")
	checkForUpdatesItem := systray.AddMenuItem("Check for updates", "Checks for a newer version of this application")
	reportIssueItem := systray.AddMenuItem("Report issue", "Report an application related issue on GitHub")

	systray.AddSeparator()
	quitItem := systray.AddMenuItem("Quit", "Closes the application")

	go func() {
		for {
			select {
			case <-newNoteItem.ClickedCh:
				glib.IdleAdd(gui.CreateNoteGUIWithDefaults)

			case <-settingsItem.ClickedCh:
				glib.IdleAdd(gui.ShowSettingsDialog)

			case <-shortcutsItem.ClickedCh:
				glib.IdleAdd(gui.ShowShortcutsDialog)

			case <-checkForUpdatesItem.ClickedCh:
				if updates.IsUpdateAvailable() {
					glib.IdleAdd(updates.AskIfTheLatestReleaseShouldBeOpenedInBrowser)
				} else {
					glib.IdleAdd(updates.ShowUpToDateDialog)
				}

			case <-reportIssueItem.ClickedCh:
				open.Run(errors.CreateIssueUrl("None"))

			case <-quitItem.ClickedCh:
				fmt.Println("I am alive")
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
