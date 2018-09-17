package config

import (
	"os"
)

//ConfigPath path to save non-data files in by default
var ConfigPath string

//CreateNeccessaryFiles creates the config folder
func CreateNeccessaryFiles() {
	os.MkdirAll(ConfigPath, os.ModePerm)
}
