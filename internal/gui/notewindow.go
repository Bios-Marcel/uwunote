package gui

import (
	"fmt"
	"os"
	"path/filepath"
	"time"

	"github.com/gotk3/gotk3/gdk"
	"github.com/gotk3/gotk3/glib"
	"github.com/gotk3/gotk3/gtk"

	"github.com/UwUNote/uwunote/internal/config"
	"github.com/UwUNote/uwunote/internal/data"
	"github.com/UwUNote/uwunote/internal/globconstants"
	"github.com/UwUNote/uwunote/internal/util"
)

//CreateWindowForNote creates a window at the given position and with the
//given dimensions, that contains the content of the passed file.
func CreateWindowForNote(file string, x, y, width, height int) {
	const defaultXOffsetNewNote = 20
	const defaultYOffsetNewNote = 20

	killSaveRoutineChannel := make(chan bool)

	//Error variable to be reused
	var gtkError error

	appConfig := config.GetAppConfig()

	// Create a new toplevel window and connect it to the
	// "destroy" signal to exit the GTK main loop when it is destroyed.
	win, gtkError := gtk.WindowNew(gtk.WINDOW_TOPLEVEL)
	util.LogAndExitOnError(gtkError)

	//setting the title in order to allow the user to distinguish the windows in his os gui or so
	win.SetTitle(globconstants.ApplicationName + " - " + file)

	newButton, gtkError := gtk.ButtonNew()
	util.LogAndExitOnError(gtkError)

	newButton.SetLabel("New")
	newButton.Connect("clicked", func() {
		currentX, currentY := win.GetPosition()
		CreateNoteGUI(currentX+defaultXOffsetNewNote, currentY+defaultYOffsetNewNote, appConfig.DefaultNoteWidth, appConfig.DefaultNoteHeight, win)
	})
	newButton.SetHExpand(false)

	deleteButton, gtkError := gtk.ButtonNew()
	util.LogAndExitOnError(gtkError)

	deleteButton.SetLabel("Delete")
	deleteButton.Connect("clicked", func() { deleteNoteGUI(appConfig, file, win, killSaveRoutineChannel) })
	deleteButton.SetHExpand(false)
	deleteButton.SetHAlign(gtk.ALIGN_END)

	titleBar, gtkError := gtk.HeaderBarNew()
	util.LogAndExitOnError(gtkError)

	titleBar.PackStart(newButton)
	titleBar.PackEnd(deleteButton)

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

	fileContent, loadError := data.LoadNote(file)
	util.LogAndExitOnError(loadError)
	buffer.SetText(string(fileContent))

	if appConfig.AutoIndent {
		registerAutoIndentListener(buffer)
	}

	if appConfig.AutoSaveAfterTyping {
		//Creating the timer beforehand, so its never nil
		saveTimer := time.NewTimer(0)
		saveTimer.Stop()

		go func() {
		SaveLoop:
			for {
				select {
				case <-saveTimer.C:
					//has to be run in gtk thread in order to show error dialogs
					glib.IdleAdd(func() { saveNoteGUI(win, file, buffer) })

				case <-killSaveRoutineChannel:
					break SaveLoop
				}
			}
		}()

		var delay int
		if appConfig.AutoSaveAfterTypingDelay < 0 {
			delay = 0
		} else {
			delay = appConfig.AutoSaveAfterTypingDelay
		}
		saveTimerDuration := time.Millisecond * time.Duration(delay)
		buffer.ConnectAfter("changed", func(textBuffer *gtk.TextBuffer) {
			saveTimer.Reset(saveTimerDuration)
		})
	}

	nodeLayout, gtkError := gtk.BoxNew(gtk.ORIENTATION_VERTICAL, 0)
	util.LogAndExitOnError(gtkError)

	nodeLayout.Add(textViewScrollPane)
	nodeLayout.SetVExpand(true)

	win.SetTitlebar(titleBar)
	win.Add(nodeLayout)

	win.Connect("key_release_event", func(window *gtk.Window, event *gdk.Event) {
		keyEvent := gdk.EventKeyNewFromEvent(event)
		if (keyEvent.State() & gdk.GDK_CONTROL_MASK) == gdk.GDK_CONTROL_MASK {
			if keyEvent.KeyVal() == gdk.KEY_s {
				saveNoteGUI(win, file, buffer)
			} else if keyEvent.KeyVal() == gdk.KEY_d {
				deleteNoteGUI(appConfig, file, win, killSaveRoutineChannel)
			} else if keyEvent.KeyVal() == gdk.KEY_n {
				currentX, currentY := win.GetPosition()
				CreateNoteGUI(currentX+defaultXOffsetNewNote, currentY+defaultYOffsetNewNote, appConfig.DefaultNoteWidth, appConfig.DefaultNoteHeight, win)
			} else if keyEvent.KeyVal() == gdk.KEY_o {
				ShowSettingsDialog()
			}
		}
	})

	makeWindowSturdy(win)

	noteName := filepath.Base(file)
	registerWindowStatePersister(noteName, win)

	win.Move(x, y)
	win.SetDefaultSize(width, height)

	// Recursively show all widgets contained in this window.
	win.ShowAll()
}

func registerAutoIndentListener(buffer *gtk.TextBuffer) {
	buffer.ConnectAfter("insert-text", func(textBuffer *gtk.TextBuffer, textIter *gtk.TextIter, chars string) {
		if chars != "\r\n" && chars != "\n" {
			return
		}

		//Count tabs on previous line
		textIter.BackwardLine()
		textIter.BackwardChars(textIter.GetLineOffset())
		amountOfTabs := 0
		for {
			if textIter.GetChar() != '\t' {
				break
			}

			amountOfTabs++
			if textIter.EndsLine() {
				break
			}

			textIter.ForwardChar()
		}

		//Insert same amounts of tabs from previous line onto next line
		if amountOfTabs > 0 {
			textIter.ForwardLine()

			for i := 0; i < amountOfTabs; i++ {
				textBuffer.Insert(textIter, "\t")
			}
		}
	})

}

func registerWindowStatePersister(identifier string, window *gtk.Window) {
	window.Connect("configure-event", func(window *gtk.Window, event *gdk.Event) {
		windowX, windowY := window.GetPosition()
		windowWidth, windowHeight := window.GetSize()

		config.SetWindowDataForFile(identifier, windowX, windowY, windowWidth, windowHeight)

		config.PersistWindowConfiguration()
	})
}

func saveNoteGUI(window *gtk.Window, file string, textBuffer *gtk.TextBuffer) {
	displaySaveError := func(err error) {
		message := fmt.Sprintf("Error saving note '%s' (%s).", file, err.Error())
		dialog := gtk.MessageDialogNew(window, gtk.DIALOG_DESTROY_WITH_PARENT, gtk.MESSAGE_ERROR, gtk.BUTTONS_OK, message)
		dialog.Run()
		dialog.Destroy()
	}

	iterStart, iterEnd := textBuffer.GetBounds()
	textToSave, textError := textBuffer.GetText(iterStart, iterEnd, true)

	if textError != nil {
		displaySaveError(textError)
	}

	writeError := data.SaveNote(file, []byte(textToSave))
	if writeError != nil {
		displaySaveError(writeError)
	}
}

func deleteNoteGUI(appConfig *config.AppConfig, file string, win *gtk.Window, killSaveRoutineChannel chan bool) {
	if appConfig.AskBeforeNoteDeletion {
		deleteDialog := gtk.MessageDialogNew(win, gtk.DIALOG_DESTROY_WITH_PARENT, gtk.MESSAGE_QUESTION, gtk.BUTTONS_YES_NO, "Are you sure, that you want to delete this note.")
		choice := deleteDialog.Run()
		deleteDialog.Destroy()
		if choice != gtk.RESPONSE_YES {
			return
		}
	}

	deleteError := data.DeleteNote(file)
	if deleteError == nil {
		killSaveRoutineChannel <- true
		win.Close()

		glib.IdleAdd(func() {
			if !appConfig.ShowTrayIcon {
				noNotesDialog := gtk.MessageDialogNew(win, gtk.DIALOG_DESTROY_WITH_PARENT, gtk.MESSAGE_QUESTION, gtk.BUTTONS_NONE, "All notes have been deleted.\nDo you want to create a new note or close the application?")

				const responseNewNote = 0
				const responseCloseApplicaion = 1

				noNotesDialog.AddButton("Create new note", responseNewNote)
				noNotesDialog.AddButton("Close application", responseCloseApplicaion)

				choice := noNotesDialog.Run()
				noNotesDialog.Destroy()
				if choice == responseNewNote {
					//Gonna ignore this error for now, as it probably means the users tinkered with the files manually
					amountOfNotes, _ := data.GetAmountOfNotes()
					if amountOfNotes == 0 {
						CreateNoteGUIWithDefaults()
					}
				} else {
					gtk.MainQuit()
					os.Exit(0)
				}
			}
		})
	} else {
		message := fmt.Sprintf("Error deleting note '%s' (%s).", file, deleteError.Error())
		dialog := gtk.MessageDialogNew(win, gtk.DIALOG_DESTROY_WITH_PARENT, gtk.MESSAGE_ERROR, gtk.BUTTONS_OK, message)
		dialog.Run()
		dialog.Destroy()
	}
}

//CreateNoteGUIWithDefaults generates a new notefile and opens the corresponding window.
func CreateNoteGUIWithDefaults() {
	appConfig := config.GetAppConfig()
	CreateNoteGUI(appConfig.DefaultNoteX, appConfig.DefaultNoteY, appConfig.DefaultNoteWidth, appConfig.DefaultNoteHeight, nil)
}

//CreateNoteGUI generates a new notefile and opens the corresponding window.
func CreateNoteGUI(x, y, width, height int, nullableParent *gtk.Window) {
	newNotePath, createError := data.CreateNote()

	if createError == nil {
		CreateWindowForNote(*newNotePath, x, y, width, height)
	} else {
		message := fmt.Sprintf("Error creating new note (%s).", createError.Error())
		dialog := gtk.MessageDialogNew(nullableParent, gtk.DIALOG_DESTROY_WITH_PARENT, gtk.MESSAGE_ERROR, gtk.BUTTONS_OK, message)
		dialog.Run()
		dialog.Destroy()
	}
}
