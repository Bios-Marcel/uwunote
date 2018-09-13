package updates

import (
	"context"

	"github.com/Masterminds/semver"
	"github.com/google/go-github/github"
	"github.com/gotk3/gotk3/gtk"
	"github.com/skratchdot/open-golang/open"
)

//AppVersion is the version of the client, by default it is always empty
var AppVersion string

//IsUpdateAvailable checks wether there is a release with a higher version number and returns true if so.
func IsUpdateAvailable() bool {
	if len(AppVersion) == 0 {
		//Omitting check, since this seems to be a development build.
		return false
	}

	curentVersion, parseError := semver.NewVersion(AppVersion)
	if parseError != nil {
		//Omitting check, since this seems to be a development build.
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

	return tagAsSemver.GreaterThan(curentVersion)
}

//ShowUpToDateDialog informs the user that his installation is on the latest version.
func ShowUpToDateDialog() {
	upToDateDialog := gtk.MessageDialogNew(nil, gtk.DIALOG_DESTROY_WITH_PARENT, gtk.MESSAGE_INFO, gtk.BUTTONS_OK, "Your installation is up to date.")
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
