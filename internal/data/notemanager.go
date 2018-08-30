package data

import (
	"io/ioutil"
	"os"

	"github.com/Bios-Marcel/wastebasket"
	uuid "github.com/satori/go.uuid"

	"github.com/UwUNote/uwunote/internal/config"
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
	if config.GetAppConfig().DeleteNotesToTrashbin {
		return wastebasket.Trash(file)
	}

	return os.Remove(file)
}

//CreateNote generates a new notefile and returns an error on failure.
func CreateNote() (*string, error) {
	fileName := uuid.Must(uuid.NewV4())
	newNotePath := config.GetAppConfig().NoteDirectory + string(os.PathSeparator) + fileName.String()
	newFile, createError := os.Create(newNotePath)
	if createError != nil {
		return nil, createError
	}
	defer newFile.Close()

	return &newNotePath, nil
}

//GetAmountOfNotes returns the amount of existing notes or `0` and an error.
//Only top level files are counted.
func GetAmountOfNotes() (int, error) {
	files, err := ioutil.ReadDir(config.GetAppConfig().NoteDirectory)

	if err != nil {
		return 0, err
	}

	amountOfFilesOnTopLevel := len(files)
	for _, fileInfo := range files {
		if fileInfo.IsDir() {
			amountOfFilesOnTopLevel--
		}
	}

	return amountOfFilesOnTopLevel, nil
}
