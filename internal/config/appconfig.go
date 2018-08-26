package config

import (
	"encoding/json"
	"io"
	"io/ioutil"
	"os"
	"path/filepath"

	"github.com/skratchdot/open-golang/open"

	"github.com/UwUNote/uwunote/internal/util"
)

//AppConfig contains all possible configuration values
type AppConfig struct {
	AskBeforeNoteDeletion bool

	AutoSaveAfterTyping      bool
	AutoSaveAfterTypingDelay int

	AutoIndent bool

	NoteDirectory string

	ShowTrayIcon bool

	DefaultNoteWidth  int
	DefaultNoteHeight int
	DefaultNoteX      int
	DefaultNoteY      int
}

var (
	appConfigPath    = filepath.Join(configPath, "app.json")
	appConfiguration = AppConfig{
		AskBeforeNoteDeletion: true,

		AutoSaveAfterTyping:      true,
		AutoSaveAfterTypingDelay: 3000,

		AutoIndent: true,

		NoteDirectory: filepath.Join(util.HomeDir, "notes"),

		ShowTrayIcon: true,

		DefaultNoteWidth:  300,
		DefaultNoteHeight: 350,
		DefaultNoteX:      0,
		DefaultNoteY:      0,
	}
)

//LoadAppConfig loads the configuration or creates a default one if none is present
func LoadAppConfig() error {
	file, openError := os.Open(appConfigPath)
	if openError != nil && os.IsNotExist(openError) {
		appConfigurationJSON, _ := json.MarshalIndent(&appConfiguration, "", "\t")
		writeError := ioutil.WriteFile(appConfigPath, appConfigurationJSON, os.ModePerm)

		if writeError != nil {
			return writeError
		}
	} else if openError == nil || os.IsExist(openError) {
		defer file.Close()
		decoder := json.NewDecoder(file)
		appConfigLoadError := decoder.Decode(&appConfiguration)

		if appConfigLoadError != nil && appConfigLoadError != io.EOF {
			return appConfigLoadError
		}
	}

	return nil
}

//GetAppConfig returns a copy off the apps configuration
func GetAppConfig() AppConfig {
	return appConfiguration
}

//OpenAppConfig opens the configuration file with its given application
func OpenAppConfig() {
	open.Run(appConfigPath)
}
