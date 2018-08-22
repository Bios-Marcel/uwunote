package util

import (
	"github.com/mitchellh/go-homedir"
)

var (
	//HomeDir is the path to the users home directory, usually C:\users\user or /home/user
	HomeDir = getHomeDir()
)

func getHomeDir() string {
	//TODO Consider a better solution

	dir, err := homedir.Dir()
	if err != nil {
		return ""
	}

	return dir
}
