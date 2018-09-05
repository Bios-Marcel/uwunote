// +build windows

package gui

import (
	"unsafe"

	"github.com/gotk3/gotk3/gdk"
	"github.com/gotk3/gotk3/gtk"
	"golang.org/x/sys/windows"
)

var (
	mod                  = windows.NewLazyDLL("user32.dll")
	setWindowPosFunction = mod.NewProc("SetWindowPos")
	findWindowFunction   = mod.NewProc("FindWindowA")
)

//ToBackground moves the window to the background, preserving its size and position
func ToBackground(gtkWindow *gtk.Window) error {
	title, _ := gtkWindow.GetTitle()
	titleAsByteArray := []byte(title)
	hwnd, _, findError := findWindowFunction.Call(0, uintptr(unsafe.Pointer(&titleAsByteArray[0])))
	if findError != nil && findError.Error() != "The operation completed successfully." {
		return findError
	}

	x, y := gtkWindow.GetPosition()
	width, height := gtkWindow.GetSize()
	_, _, setPosError := setWindowPosFunction.Call(hwnd, uintptr(1), uintptr(x), uintptr(y), uintptr(width), uintptr(height), 0)

	if setPosError != nil && setPosError.Error() != "The operation completed successfully." {
		return setPosError
	}

	return nil
}

func makeWindowSturdy(window *gtk.Window) {
	//The Utility type causes the window to be unminimizable, not appear on the
	//taskbar,not appear in Alt + Tab and stick on all desktops.
	window.SetTypeHint(gdk.WINDOW_TYPE_HINT_UTILITY)

	ToBackground(window)
}
