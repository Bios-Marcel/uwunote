package gui

import (
	"fmt"
	"os"
	"path/filepath"
	"time"

	"github.com/gotk3/gotk3/gdk"
	"github.com/gotk3/gotk3/glib"
	"github.com/gotk3/gotk3/gtk"
	"github.com/skratchdot/open-golang/open"

	"github.com/UwUNote/uwunote/internal/config"
	"github.com/UwUNote/uwunote/internal/data"
	"github.com/UwUNote/uwunote/internal/errors"
)

var (
	appMenu *gtk.Menu
)

func getAppMenu() *gtk.Menu {
	if appMenu == nil {
		var gtkError error
		appMenu, gtkError = gtk.MenuNew()

		newNoteItem, gtkError := gtk.MenuItemNewWithLabel("New note")
		errors.ShowErrorDialogOnError(gtkError)
		newNoteItem.SetTooltipText("Creates a new note")
		newNoteItem.Connect("activate", func() {
			CreateNoteGUIWithDefaults()
		})

		separatorOne, gtkError := gtk.SeparatorMenuItemNew()
		errors.ShowErrorDialogOnError(gtkError)

		settingsItem, gtkError := gtk.MenuItemNewWithLabel("Settings")
		errors.ShowErrorDialogOnError(gtkError)
		settingsItem.SetTooltipText("Opens the settings")
		settingsItem.Connect("activate", func() {
			ShowSettingsDialog()
		})

		shortcutsItem, gtkError := gtk.MenuItemNewWithLabel("Shortcuts")
		errors.ShowErrorDialogOnError(gtkError)
		shortcutsItem.SetTooltipText("Opens the shortcuts dialog")
		shortcutsItem.Connect("activate", func() {
			ShowShortcutsDialog()
		})

		reportIssueItem, gtkError := gtk.MenuItemNewWithLabel("Report issue")
		errors.ShowErrorDialogOnError(gtkError)
		reportIssueItem.SetTooltipText("Report an application related issue on GitHub")
		reportIssueItem.Connect("activate", func() {
			open.Run(errors.CreateIssueUrl("None"))
		})

		separatorTwo, gtkError := gtk.SeparatorMenuItemNew()
		errors.ShowErrorDialogOnError(gtkError)

		quitItem, gtkError := gtk.MenuItemNewWithLabel("Quit")
		errors.ShowErrorDialogOnError(gtkError)
		quitItem.SetTooltipText("Closes the application")
		quitItem.Connect("activate", func() {
			gtk.MainQuit()
			os.Exit(0)
		})

		appMenu.Append(newNoteItem)
		appMenu.Append(separatorOne)
		appMenu.Append(settingsItem)
		appMenu.Append(shortcutsItem)
		appMenu.Append(reportIssueItem)
		appMenu.Append(separatorTwo)
		appMenu.Append(quitItem)
		appMenu.ShowAll()
	}

	fmt.Println(appMenu)

	return appMenu
}

//CreateWindowForNote creates a window at the given position and with the
//given dimensions, that contains the content of the passed file.
func CreateWindowForNote(file string, x, y, width, height int) {
	const defaultXOffsetNewNote = 20
	const defaultYOffsetNewNote = 20

	killSaveRoutineChannel := make(chan bool)

	//Error variable to be reused
	var gtkError error

	appConfig := config.GetAppConfig()

	win, gtkError := gtk.WindowNew(gtk.WINDOW_TOPLEVEL)
	errors.ShowErrorDialogOnError(gtkError)

	newButton, gtkError := gtk.ButtonNew()
	errors.ShowErrorDialogOnError(gtkError)

	newButton.SetLabel("New")
	newButton.Connect("clicked", func() {
		currentX, currentY := win.GetPosition()
		CreateNoteGUI(currentX+defaultXOffsetNewNote, currentY+defaultYOffsetNewNote, appConfig.DefaultNoteWidth, appConfig.DefaultNoteHeight, win)
	})
	newButton.SetHExpand(false)

	deleteButton, gtkError := gtk.ButtonNew()
	errors.ShowErrorDialogOnError(gtkError)

	deleteButton.SetLabel("Delete")
	deleteButton.Connect("clicked", func() { deleteNoteGUI(appConfig, file, win, killSaveRoutineChannel) })
	deleteButton.SetHExpand(false)
	deleteButton.SetHAlign(gtk.ALIGN_END)

	titleBar, gtkError := gtk.HeaderBarNew()
	errors.ShowErrorDialogOnError(gtkError)

	titleBar.PackStart(newButton)
	menuButton, gtkError := gtk.MenuButtonNew()
	errors.ShowErrorDialogOnError(gtkError)
	menuButton.SetPopup(getAppMenu())
	titleBar.Add(menuButton)
	titleBar.PackEnd(deleteButton)

	var hAdjustment, vAdjustment *gtk.Adjustment
	textViewScrollPane, gtkError := gtk.ScrolledWindowNew(hAdjustment, vAdjustment)
	errors.ShowErrorDialogOnError(gtkError)

	textView, gtkError := gtk.TextViewNew()
	errors.ShowErrorDialogOnError(gtkError)

	textView.SetVExpand(true)
	textView.SetHExpand(true)
	textView.SetLeftMargin(5)
	textView.SetRightMargin(5)
	textView.SetWrapMode(appConfig.WrapMode)

	//Wrapping the textView in a scrollpane, otherwise the window will expand instead
	textViewScrollPane.Add(textView)

	buffer, gtkError := textView.GetBuffer()
	errors.ShowErrorDialogOnError(gtkError)

	fileContent, loadError := data.LoadNote(file)
	errors.ShowErrorDialogOnError(loadError)
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
	errors.ShowErrorDialogOnError(gtkError)

	nodeLayout.Add(textViewScrollPane)
	nodeLayout.SetVExpand(true)

	win.SetTitlebar(titleBar)
	win.Add(nodeLayout)

	win.Connect("key_release_event", func(window *gtk.Window, event *gdk.Event) {
		keyEvent := gdk.EventKeyNewFromEvent(event)
		keyEventState := keyEvent.State()

		if (keyEventState & gdk.GDK_CONTROL_MASK) == gdk.GDK_CONTROL_MASK {
			keyVal := keyEvent.KeyVal()

			if keyVal == gdk.KEY_s {
				saveNoteGUI(win, file, buffer)
			} else if keyVal == gdk.KEY_d {
				deleteNoteGUI(appConfig, file, win, killSaveRoutineChannel)
			} else if keyVal == gdk.KEY_n {
				currentX, currentY := win.GetPosition()
				CreateNoteGUI(currentX+defaultXOffsetNewNote, currentY+defaultYOffsetNewNote, appConfig.DefaultNoteWidth, appConfig.DefaultNoteHeight, win)
			} else if keyVal == gdk.KEY_o {
				ShowSettingsDialog()
			}
		} else if (keyEventState & uint(gtk.AcceleratorGetDefaultModMask())) == 0 {
			if keyEvent.KeyVal() == gdk.KEY_F1 {
				ShowShortcutsDialog()
			}
		}
	})

	noteName := filepath.Base(file)
	registerWindowStatePersister(noteName, win)

	win.Move(x, y)
	win.SetDefaultSize(width, height)

	showWindowSturdy(win)
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
		config.DeleteWindowDataForFile(filepath.Base(file))
		config.PersistWindowConfiguration()

		glib.IdleAdd(func() {
			amountOfNotes, _ := data.GetAmountOfNotes()
			if amountOfNotes == 0 {
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
		errors.ShowErrorDialogWithMessage("Error creating new note.", createError)
	}
}
