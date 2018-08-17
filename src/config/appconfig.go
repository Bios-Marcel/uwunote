package config

import (
	"encoding/json"
	"io"
	"io/ioutil"
	"os"

	"github.com/Bios-Marcel/uwuNote/src/util"
)

//AppConfig contains all possible configuration values
type AppConfig struct {
	AskBeforeNoteDeletion bool

	AutoSaveAfterTyping      bool
	AutoSaveAfterTypingDelay int

	AutoIndent bool
}

var (
	appConfigPath    = configPath + string(os.PathSeparator) + "app.json"
	appConfiguration = AppConfig{
		AskBeforeNoteDeletion: true,

		AutoSaveAfterTyping:      true,
		AutoSaveAfterTypingDelay: 3000,

		AutoIndent: true,
	}
)

//LoadAppConfig loads the configuration or creates a default one if none is present
func LoadAppConfig() {
	file, openError := os.Open(appConfigPath)
	if openError != nil && os.IsNotExist(openError) {
		appConfigurationJSON, _ := json.MarshalIndent(&appConfiguration, "", "\t")
		writeError := ioutil.WriteFile(appConfigPath, appConfigurationJSON, os.ModePerm)
		//TODO Better way?
		util.LogAndExitOnError(writeError)
	} else if openError == nil || os.IsExist(openError) {
		defer file.Close()
		decoder := json.NewDecoder(file)
		appConfigLoadError := decoder.Decode(&appConfiguration)

		if appConfigLoadError != io.EOF {
			util.LogAndExitOnError(appConfigLoadError)
		}
	}
}

//GetAppConfig returns a copy off the apps configuration
func GetAppConfig() AppConfig {
	return appConfiguration
}
