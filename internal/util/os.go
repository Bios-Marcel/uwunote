package util

import (
	"github.com/mitchellh/go-homedir"
)

var (
	//HomeDir is the path to the users home directory, usually C:\users\user or /home/user
	HomeDir = getHomeDir()
)

func getHomeDir() string {
	//This is a severe error and therefore a panic is okay here!
	dir, homeDirError := homedir.Dir()
	if homeDirError != nil {
		panic(homeDirError)
	}

	return dir
}
