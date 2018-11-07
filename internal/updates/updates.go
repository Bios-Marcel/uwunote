// +build go1.7

package updates

import (
	"context"
	"fmt"

	"github.com/Masterminds/semver"
	"github.com/google/go-github/github"
	"github.com/gotk3/gotk3/gtk"
	"github.com/skratchdot/open-golang/open"
)

//AppVersion is the version of the client, by default it is always empty
var AppVersion string

//IsUpdateAvailable checks wether there is a release with a higher version number and returns true if so.
func IsUpdateAvailable() bool {
	currentVersion := VersionAsSemver()
	if currentVersion == nil {
		return false
	}

	client := github.NewClient(nil)
	release, response, err := client.Repositories.GetLatestRelease(context.Background(), "UwUNote", "uwunote")
	if response.StatusCode == 404 {
		//There is no latest release (no release at all).
		return false
	}

	if err != nil {
		//TODO Show error dialog
		return false
	}

	tagAsSemver := semver.MustParse(release.GetTagName())

	return tagAsSemver.GreaterThan(currentVersion)
}

//VersionAsSemver creates a semver.Version of the current AppVersion string.
//If the AppVersion string is invalid or empty it returns nil
func VersionAsSemver() *semver.Version {
	if len(AppVersion) == 0 {
		//Omitting check, since this seems to be a development build.
		return nil
	}

	currentVersion, parseError := semver.NewVersion(AppVersion)
	if parseError != nil {
		//Omitting check, since this seems to be a development build.
		return nil
	}

	return currentVersion
}

//ShowUpToDateDialog informs the user that his installation is on the latest version.
func ShowUpToDateDialog() {
	message := ""
	if VersionAsSemver() == nil {
		message = "You are currently running a development build."
	} else {
		message = fmt.Sprintf("Your installation is up to date. You are currently running version %s.", AppVersion)
	}

	upToDateDialog := gtk.MessageDialogNew(nil, gtk.DIALOG_DESTROY_WITH_PARENT, gtk.MESSAGE_INFO, gtk.BUTTONS_OK, message)
	upToDateDialog.Run()
	upToDateDialog.Destroy()
}

//AskIfTheLatestReleaseShouldBeOpenedInBrowser opens the latest release of UwUNote/uwunote in the users browser if the user wants to.
func AskIfTheLatestReleaseShouldBeOpenedInBrowser() {
	openUpdateDialog := gtk.MessageDialogNew(nil, gtk.DIALOG_DESTROY_WITH_PARENT, gtk.MESSAGE_QUESTION, gtk.BUTTONS_YES_NO, "There is an update available, do you want to show it in your browser?")
	choice := openUpdateDialog.Run()
	openUpdateDialog.Destroy()
	if choice == gtk.RESPONSE_YES {
		open.Run("https://github.com/UwUNote/uwunote/releases/latest")
	}
}
