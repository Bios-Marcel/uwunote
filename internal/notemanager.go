package internal

import (
	"io/ioutil"
	"os"

	uuid "github.com/satori/go.uuid"

	"github.com/Bios-Marcel/uwuNote/internal/config"
)

var (
	//Will be customizable at some point
	notePath = config.GetAppConfig().NoteDirectory
)

//LoadNote loads the content of a note and returns an error on failure.
func LoadNote(file string) ([]byte, error) {
	return ioutil.ReadFile(file)
}

//SaveNote saves the given data into the note
func SaveNote(file string, textToSave []byte) error {
	return ioutil.WriteFile(file, textToSave, os.ModeType)
}

//DeleteNote deletes a notes data
func DeleteNote(file string) error {
	return os.Remove(file)
}

//CreateNote generates a new notefile and returns an error on failure.
func CreateNote() (*string, error) {
	fileName := uuid.Must(uuid.NewV4())
	newNotePath := notePath + string(os.PathSeparator) + fileName.String()
	newFile, createError := os.Create(newNotePath)
	if createError != nil {
		return nil, createError
	}
	defer newFile.Close()

	return &newNotePath, nil
}
