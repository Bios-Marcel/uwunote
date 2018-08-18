package util

import "log"

//Logs the given errors and calls os.exit
func LogAndExitOnError(possibleError error) {
	if possibleError != nil {
		log.Fatal(possibleError)
	}
}
