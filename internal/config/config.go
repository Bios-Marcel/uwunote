package config

import (
	"os"
	"path/filepath"

	"github.com/UwUNote /uwunote/internal/util"
)

var (
	configPath = filepath.Join(util.HomeDir, ".uwunote")
)

//CreateNeccessaryFiles creates the config folder
func CreateNeccessaryFiles() {
	os.MkdirAll(configPath, os.ModePerm)
}
