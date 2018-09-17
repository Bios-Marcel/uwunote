package config

import (
	"encoding/json"
	"io"
	"io/ioutil"
	"os"
	"path/filepath"

	"github.com/UwUNote/uwunote/internal/util"
)

//AppConfig contains all possible configuration values
type AppConfig struct {
	DeleteNotesToTrashbin bool
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
	appConfiguration = GetAppConfigDefaults()
)

func getAppConfigPath() string {
	return filepath.Join(ConfigPath, "app.json")
}

//GetAppConfigDefaults returns all defaults for the application configuration
func GetAppConfigDefaults() AppConfig {
	return AppConfig{
		DeleteNotesToTrashbin: true,
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
}

//LoadAppConfig loads the configuration or creates a default one if none is present
func LoadAppConfig() error {
	file, openError := os.Open(getAppConfigPath())
	if openError != nil && os.IsNotExist(openError) {
		writeError := PersistAppConfig(&appConfiguration)

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

//GetAppConfigCopy returns a copy off the apps configuration
func GetAppConfigCopy() AppConfig {
	return appConfiguration
}

//GetAppConfig returns the loaded AppConfig
func GetAppConfig() *AppConfig {
	return &appConfiguration
}

//PersistAppConfig writes the given AppConfig as JSON to the filesystem.
func PersistAppConfig(config *AppConfig) error {
	appConfigurationJSON, _ := json.MarshalIndent(&config, "", "\t")
	return ioutil.WriteFile(getAppConfigPath(), appConfigurationJSON, os.ModePerm)
}
