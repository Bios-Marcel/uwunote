// +build !no_systray,go1.6,!go1.7

package internal

import (
	"os"

	"github.com/UwUNote/uwunote/internal/errors"
	"github.com/UwUNote/uwunote/internal/gui"
	"github.com/getlantern/systray"
	"github.com/gotk3/gotk3/glib"
	"github.com/gotk3/gotk3/gtk"
	"github.com/skratchdot/open-golang/open"
)

func buildSystray() {
	systray.SetIcon(gui.AppIcon)
	newNoteItem := systray.AddMenuItem("New note", "Creates a new note")
	systray.AddSeparator()
	settingsItem := systray.AddMenuItem("Settings", "Opens the settings")
	shortcutsItem := systray.AddMenuItem("Shortcuts", "Opens the shortcuts dialog")
	reportIssueItem := systray.AddMenuItem("Report issue", "Report an application related issue on GitHub")

	systray.AddSeparator()
	quitItem := systray.AddMenuItem("Quit", "Closes the application")

	go func() {
		for {
			select {
			case <-newNoteItem.ClickedCh:
				glib.IdleAdd(gui.CreateNoteGUIWithDefaults)

			case <-settingsItem.ClickedCh:
				glib.IdleAdd(gui.ShowSettingsDialog)

			case <-shortcutsItem.ClickedCh:
				glib.IdleAdd(gui.ShowShortcutsDialog)

			case <-reportIssueItem.ClickedCh:
				open.Run(errors.CreateIssueUrl("None"))

			case <-quitItem.ClickedCh:
				glib.IdleAdd(func() {
					gtk.MainQuit()
					os.Exit(0)
				})
			}
		}
	}()
}
