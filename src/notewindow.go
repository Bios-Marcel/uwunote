package main

import (
	"io/ioutil"
	"os"
	"path/filepath"
	"time"

	"github.com/gotk3/gotk3/gdk"
	"github.com/gotk3/gotk3/gtk"
	uuid "github.com/satori/go.uuid"

	"github.com/bios-marcel/uwunote/src/config"
	"github.com/bios-marcel/uwunote/src/util"
)

func createWindowForNote(file string, x, y, width, height int) {
	deleteNoteChannel := make(chan bool)

	//Error variable to be reused
	var gtkError error

	// Create a new toplevel window and connect it to the
	// "destroy" signal to exit the GTK main loop when it is destroyed.
	win, gtkError := gtk.WindowNew(gtk.WINDOW_TOPLEVEL)
	util.LogAndExitOnError(gtkError)

	//TODO Is a title necessary at all?
	win.SetTitle(file)

	newButton, gtkError := gtk.ButtonNew()
	util.LogAndExitOnError(gtkError)

	newButton.SetLabel("New")
	newButton.Connect("clicked", func() { CreateNote(x+20, y+20, 300, 350) })
	newButton.SetHExpand(false)

	deleteButton, gtkError := gtk.ButtonNew()
	util.LogAndExitOnError(gtkError)

	deleteButton.SetLabel("Delete")
	deleteButton.Connect("clicked", func() { deleteNote(file, win, deleteNoteChannel) })
	deleteButton.SetHExpand(false)
	deleteButton.SetHAlign(gtk.ALIGN_END)

	topBar, gtkError := gtk.HeaderBarNew()
	util.LogAndExitOnError(gtkError)

	topBar.PackStart(newButton)
	topBar.PackEnd(deleteButton)

	var hAdjustment, vAdjustment *gtk.Adjustment
	textViewScrollPane, gtkError := gtk.ScrolledWindowNew(hAdjustment, vAdjustment)
	util.LogAndExitOnError(gtkError)

	textView, gtkError := gtk.TextViewNew()
	util.LogAndExitOnError(gtkError)

	textView.SetVExpand(true)
	textView.SetHExpand(true)

	//Wrapping the textView in a scrollpane, otherwise the window will expand instead
	textViewScrollPane.Add(textView)

	buffer, gtkError := textView.GetBuffer()
	util.LogAndExitOnError(gtkError)

	fileContent, _ := ioutil.ReadFile(file)
	buffer.SetText(string(fileContent))

	buffer.ConnectAfter("insert-text", func(textBuffer *gtk.TextBuffer, textIter *gtk.TextIter, chars string) {
		if chars != "\r\n" && chars != "\n" {
			return
		}

		//Count tabs on previous line
		textIter.BackwardLine()
		textIter.BackwardChars(textIter.GetLineOffset())
		amountOfTabs := 0
		for {
			if textIter.GetChar() == '\t' {
				amountOfTabs++
				if !textIter.EndsLine() {
					textIter.ForwardChar()
				} else {
					break
				}
			} else {
				break
			}
		}

		//Insert same amounts of tabs from previous line onto next line
		if amountOfTabs > 0 {
			textIter.ForwardLine()

			for i := 0; i < amountOfTabs; i++ {
				textBuffer.Insert(textIter, "\t")
			}
		}
	})

	//Creating the timer beforehand, so its never nil
	saveTimer := time.NewTimer(0)
	saveTimer.Stop()

	go func() {
	SaveLoop:
		for {
			select {
			case <-saveTimer.C:
				saveNote(file, textView)

			case <-deleteNoteChannel:
				break SaveLoop
			}
		}
	}()

	buffer.ConnectAfter("insert-text", func(textBuffer *gtk.TextBuffer, textIter *gtk.TextIter, chars string) {
		saveTimer.Reset(time.Second * 3)
	})

	nodeLayout, gtkError := gtk.BoxNew(gtk.ORIENTATION_VERTICAL, 0)
	util.LogAndExitOnError(gtkError)

	nodeLayout.Add(textViewScrollPane)
	nodeLayout.SetVExpand(true)

	win.SetTitlebar(topBar)
	win.Add(nodeLayout)

	win.Connect("key_release_event", func(window *gtk.Window, event *gdk.Event) {
		//Subtract default modifiers according to:
		//https://developer.gnome.org/gtk3/stable/checklist-modifiers.html
		//modifiers := gtk.AcceleratorGetDefaultModMask()

		keyEvent := gdk.EventKeyNewFromEvent(event)
		if (keyEvent.State() & gdk.GDK_CONTROL_MASK) == gdk.GDK_CONTROL_MASK {
			if keyEvent.KeyVal() == gdk.KEY_s {
				saveNote(file, textView)
			} else if keyEvent.KeyVal() == gdk.KEY_d {
				deleteNote(file, win, deleteNoteChannel)
			} else if keyEvent.KeyVal() == gdk.KEY_n {
				CreateNote(x+20, y+20, 300, 350)
			}
		}
	})

	//Rebuilding behaviour from TYPE_HINT_DESKTOP
	win.SetSkipTaskbarHint(true)
	win.SetSkipPagerHint(true)
	win.SetKeepBelow(true)
	win.Stick()

	//Making the window unminimizable
	win.Connect("window-state-event", func(window *gtk.Window, event *gdk.Event) {
		windowEvent := gdk.EventWindowStateNewFromEvent(event)
		newWindowState := windowEvent.NewWindowState()

		if (newWindowState & gdk.WINDOW_STATE_ICONIFIED) == gdk.WINDOW_STATE_ICONIFIED {
			window.Present()
		}
	})

	win.Move(x, y)
	win.SetDefaultSize(width, height)

	win.Connect("configure-event", func(window *gtk.Window, event *gdk.Event) {
		windowX, windowY := window.GetPosition()
		windowWidth, windowHeight := window.GetSize()

		noteName := filepath.Base(file)
		config.SetWindowDataForFile(noteName, windowX, windowY, windowWidth, windowHeight)

		config.PersistWindowConfiguration()
	})

	// Recursively show all widgets contained in this window.
	win.ShowAll()
}

func saveNote(file string, textView *gtk.TextView) {
	//TODO Following errors on which i now panic should probably inform the user about not being able to save.

	currentNoteBuffer, bufferError := textView.GetBuffer()
	if bufferError != nil {
		panic(bufferError)
	}

	iterStart, iterEnd := currentNoteBuffer.GetBounds()
	//TODO Check if I need the "Hidden chars"
	textToSave, textError := currentNoteBuffer.GetText(iterStart, iterEnd, false)
	if textError != nil {
		panic(textError)
	}

	writeError := ioutil.WriteFile(file, []byte(textToSave), os.ModeType)
	if writeError != nil {
		panic(writeError)
	}
}

func deleteNote(file string, win *gtk.Window, deleteNoteChannel chan bool) {
	deleteDialog := gtk.MessageDialogNew(win, gtk.DIALOG_DESTROY_WITH_PARENT, gtk.MESSAGE_QUESTION, gtk.BUTTONS_YES_NO, "Are you sure, that you want to delete this note.")
	choice := deleteDialog.Run()
	deleteDialog.Close()
	if choice == gtk.RESPONSE_YES {
		deleteNoteChannel <- true
		os.Remove(file)
		win.Close()
	}

	//TODO create new note if the last one was deleted?
}

//CreateNote generates a new notefile and opens the corresponding window.
func CreateNote(x, y, width, height int) {
	fileName := uuid.Must(uuid.NewV4())
	newNotePath := notePath + string(os.PathSeparator) + fileName.String()
	os.Create(newNotePath)
	createWindowForNote(newNotePath, x, y, width, height)
}
