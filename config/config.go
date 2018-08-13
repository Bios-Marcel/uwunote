package config

import (
	"os"
	"path/filepath"
)

var (
	configPath = filepath.FromSlash(os.Getenv("HOME") + string(os.PathSeparator) + ".uwunotes")
)
