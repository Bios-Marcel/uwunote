package config

import (
	"encoding/json"
	"io"
	"io/ioutil"
	"log"
	"os"

	"github.com/bios-marcel/uwunote/util"
)

var (
	windowConfigPath = configPath + string(os.PathSeparator) + "windows.json"
	//Configuration contains the positions and sizes for all notes
	configuration = WindowConfig{}
)

//WindowConfig contains a map of WindowData entries in a map
type WindowConfig struct {
	Data map[string]WindowData
}

//WindowData contains the sosition (x and y) and size of a Window
type WindowData struct {
	X int
	Y int

	Width  int
	Height int
}

//LoadWindowConfiguration loads the window configuration from its path.
func LoadWindowConfiguration() {
	log.Println("Loading window configuration")
	file, openError := os.Open(windowConfigPath)
	if openError == nil || os.IsExist(openError) {
		defer file.Close()
		decoder := json.NewDecoder(file)
		windowConfigLoadError := decoder.Decode(&configuration)

		if windowConfigLoadError != io.EOF {
			util.LogAndExitOnError(windowConfigLoadError)
		}
	}

	if configuration.Data == nil {
		//Creating an empty map to prevent nil pointer references
		configuration.Data = make(map[string]WindowData)
	}
}

//PersistWindowConfiguration saves the current window configuration to its iven path
func PersistWindowConfiguration() {
	windowConfigurationJSON, _ := json.Marshal(&configuration)
	writeError := ioutil.WriteFile(windowConfigPath, windowConfigurationJSON, os.ModePerm)
	//TODO Better way?
	util.LogAndExitOnError(writeError)
}

//GetWindowDataForFile retrieves the window-config entry for the given file
func GetWindowDataForFile(file string) (*WindowData, bool) {
	data, exists := configuration.Data[file]
	return &data, exists
}

//SetWindowDataForFile sets coordinates and size in the window-config for the given file
func SetWindowDataForFile(file string, x, y, width, height int) {
	configForWindow, exists := configuration.Data[file]

	if exists {
		configForWindow.X = x
		configForWindow.Y = y

		configForWindow.Width = width
		configForWindow.Height = height

		configuration.Data[file] = configForWindow
	} else {
		configuration.Data[file] = WindowData{
			X:      x,
			Y:      y,
			Width:  width,
			Height: height,
		}
	}
}
