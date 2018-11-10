// +build !no_systray,!windows,!mac

package internal

import (
	"github.com/getlantern/systray"
)

func startWithTrayIcon() {
	systray.Run(
		func() {
			systemTrayRun()

			//Only on linux gtk will be used, therefore we needn't init gtk on linux, as systray does that.
			start()
		},
		func() {})
}
