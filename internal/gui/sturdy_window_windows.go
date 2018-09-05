// +build windows

package gui

import (
	"github.com/gotk3/gotk3/gdk"
	"github.com/gotk3/gotk3/gtk"
)

func makeWindowSturdy(window *gtk.Window) {
	//The Utility type causes the window to be unminimizable, not appear on the
	//taskbar,not appear in Alt + Tab and stick on all desktops.
	window.SetTypeHint(gdk.WINDOW_TYPE_HINT_UTILITY)

	//TODO Sending the window to background currently doesn't work
}
