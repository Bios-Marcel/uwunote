// +build !windows

package gui

import (
	"github.com/gotk3/gotk3/gdk"
	"github.com/gotk3/gotk3/gtk"
)

func makeWindowSturdy(window *gtk.Window) {
	//Rebuilding behaviour from TYPE_HINT_DESKTOP
	window.SetSkipTaskbarHint(true)
	window.SetSkipPagerHint(true)
	window.SetKeepBelow(true)
	window.Stick()

	//HACK Making the window effectively unminimizable, but doesn't reliably work
	window.Connect("window-state-event", func(w *gtk.Window, event *gdk.Event) {
		windowEvent := gdk.EventWindowStateNewFromEvent(event)
		newWindowState := windowEvent.NewWindowState()

		if (newWindowState & gdk.WINDOW_STATE_ICONIFIED) == gdk.WINDOW_STATE_ICONIFIED {
			w.Present()
		}
	})
}

func showWindowSturdy(window *gtk.Window) {
	makeWindowSturdy(window)
	window.ShowAll()
}
