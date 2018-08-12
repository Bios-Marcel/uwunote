package main

import (
	"io/ioutil"
	"log"

	"os"
	"path/filepath"

	"github.com/gotk3/gotk3/gdk"
	"github.com/gotk3/gotk3/gtk"

	"github.com/satori/go.uuid"
)

var (
	//Will lateron be customizable
	notePath = filepath.FromSlash(os.Getenv("HOME") + string(os.PathSeparator) + "notes")
)

func main() {
	// Initialize GTK without parsing any command line arguments.
	gtk.Init(nil)

	generateNoteWindows()

	// Begin executing the GTK main loop.  This blocks until
	// gtk.MainQuit() is run.
	gtk.Main()
}

func generateNoteWindows() {
	var files []string

	err := filepath.Walk(notePath, func(path string, info os.FileInfo, err error) error {
		//We wil ignore all notes that lie on second level
		if info != nil && info.IsDir() {
			return nil
		}

		files = append(files, path)
		return nil
	})

	if err != nil {
		log.Fatal("Error creating notes.")
		panic(err)
	}

	for index, file := range files {
		pos := index * 310
		createWindowForNote(file, pos, 0)
	}
}

func createWindowForNote(file string, x int, y int) {
	//Error variable to be reused
	var gtkError error

	// Create a new toplevel window and connect it to the
	// "destroy" signal to exit the GTK main loop when it is destroyed.
	win, gtkError := gtk.WindowNew(gtk.WINDOW_TOPLEVEL)
	panicOnError(gtkError)

	// The app isn't killable for now.
	/*win.Connect("destroy", func() {
		gtk.MainQuit()
	})*/

	win.SetTitle("Sticky background window")

	//Make window stick on desktop
	//TODO Fix problems: Prevents DND, Allows offscreen placement
	win.SetTypeHint(gdk.WINDOW_TYPE_HINT_DESKTOP)

	newButton, gtkError := gtk.ButtonNew()
	panicOnError(gtkError)

	newButton.SetLabel("New")
	newButton.Connect("clicked", func() {
		fileName := uuid.Must(uuid.NewV4())
		newNotePath := notePath + string(os.PathSeparator) + fileName.String() + ".md"
		os.Create(newNotePath)
		createWindowForNote(newNotePath, x+20, y+20)
	})
	newButton.SetHExpand(false)

	deleteButton, gtkError := gtk.ButtonNew()
	panicOnError(gtkError)

	deleteButton.SetLabel("Delete")
	deleteButton.Connect("clicked", func() {
		os.Remove(file)
		win.Destroy()
	})
	deleteButton.SetHExpand(false)
	deleteButton.SetHAlign(gtk.ALIGN_END)

	topBar, gtkError := gtk.HeaderBarNew()
	panicOnError(gtkError)

	topBar.Add(newButton)

	sep, gtkError := gtk.SeparatorToolItemNew()
	panicOnError(gtkError)
	sep.SetHExpand(true)
	topBar.Add(sep)

	topBar.Add(deleteButton)

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

	fileContent, _ := ioutil.ReadFile(file)
	buffer.SetText(string(fileContent))

	nodeLayout, gtkError := gtk.BoxNew(gtk.ORIENTATION_VERTICAL, 0)
	panicOnError(gtkError)

	nodeLayout.Add(textViewScrollPane)
	nodeLayout.SetVExpand(true)

	win.SetTitlebar(topBar)
	win.Add(nodeLayout)

	win.SetResizable(true)
	win.Move(x, y)
	win.SetDefaultSize(300, 350)

	// Recursively show all widgets contained in this window.
	win.ShowAll()
}

// Panics if the given value isn't nil
func panicOnError(possibleError error) {
	if possibleError != nil {
		log.Fatal(possibleError)
		panic(possibleError)
	}
}
