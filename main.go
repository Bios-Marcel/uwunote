package main

import (
	"fmt"
	"log"

	"github.com/gotk3/gotk3/gdk"
	"github.com/gotk3/gotk3/gtk"
)

func main() {
	//Error variable to be reused
	var gtkError error

	// Initialize GTK without parsing any command line arguments.
	gtk.Init(nil)

	// Create a new toplevel window and connect it to the
	// "destroy" signal to exit the GTK main loop when it is destroyed.
	win, gtkError := gtk.WindowNew(gtk.WINDOW_TOPLEVEL)
	panicOnError(gtkError)

	win.Connect("destroy", func() {
		gtk.MainQuit()
	})

	win.SetTitle("Sticky background window")

	//Remove taskbar icon
	win.SetSkipTaskbarHint(true)
	//Remove platform windowframe
	win.SetDecorated(false)
	//Make window stick on desktop
	win.SetTypeHint(gdk.WINDOW_TYPE_HINT_DESKTOP)

	win.SetPosition(gtk.WIN_POS_CENTER)
	//win.Maximize()

	saveButton, gtkError := gtk.ButtonNew()
	panicOnError(gtkError)

	saveButton.SetLabel("Save")
	saveButton.Connect("clicked", func() {
		fmt.Println("Clicked")
	})

	var hAdjustment, vAdjustment *gtk.Adjustment
	textViewScrollPane, gtkError := gtk.ScrolledWindowNew(hAdjustment, vAdjustment)
	panicOnError(gtkError)

	textView, gtkError := gtk.TextViewNew()
	panicOnError(gtkError)

	textView.SetVExpand(true)
	textView.SetHExpand(true)

	//Wrapping the textView in a scrollpane, otherwise the window will expand instead
	textViewScrollPane.Add(textView)

	buffer, gtkError := textView.GetBuffer()
	panicOnError(gtkError)

	buffer.SetText("Oh, hello.")

	nodeLayout, gtkError := gtk.BoxNew(gtk.ORIENTATION_VERTICAL, 0)
	panicOnError(gtkError)

	nodeLayout.Add(saveButton)

	nodeLayout.Add(textViewScrollPane)
	nodeLayout.SetVExpand(true)

	win.Add(nodeLayout)

	// Set the default window size.
	win.SetDefaultSize(800, 600)

	// Recursively show all widgets contained in this window.
	win.ShowAll()

	// Begin executing the GTK main loop.  This blocks until
	// gtk.MainQuit() is run.
	gtk.Main()
}

// Panics if the given value isn't nil
func panicOnError(possibleError error) {
	if possibleError != nil {
		log.Fatal(possibleError)
		panic(possibleError)
	}
}
