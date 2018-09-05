// +build windows darwin

package internal

import (
	"github.com/getlantern/systray"
)

func startWithTrayIcon() {
	systray.Run(
		func() {
			systemTrayRun()

			startAndInitGtk()
		},
		func() {})
}
