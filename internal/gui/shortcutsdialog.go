package gui

import (
	"fmt"

	"github.com/gotk3/gotk3/gtk"
)

var shortcutsDialogUIFile = `
<?xml version="1.0" encoding="UTF-8"?>
<interface>
  <object class="GtkShortcutsWindow" id="shortcuts-uwunote">
    <property name="modal">1</property>
    <child>
      <object class="GtkShortcutsSection">
        <property name="visible">1</property>
        <property name="section-name">shortcuts</property>
        <property name="max-height">12</property>
        <child>
          <object class="GtkShortcutsGroup">
            <property name="visible">1</property>
            <property name="title" translatable="yes">General</property>
            <child>
              <object class="GtkShortcutsShortcut">
                <property name="visible">1</property>
                <property name="accelerator">F1</property>
                <property name="title" translatable="yes">Show shortcuts</property>
              </object>
            </child>
            <child>
              <object class="GtkShortcutsShortcut">
                <property name="visible">1</property>
                <property name="accelerator">&lt;ctrl&gt;O</property>
                <property name="title" translatable="yes">Open Settings</property>
              </object>
            </child>
          </object>
        </child>
        <child>
          <object class="GtkShortcutsGroup">
            <property name="visible">1</property>
            <property name="title" translatable="yes">Notes</property>
            <child>
              <object class="GtkShortcutsShortcut">
                <property name="visible">1</property>
                <property name="accelerator">&lt;ctrl&gt;n</property>
                <property name="title" translatable="yes">New</property>
              </object>
            </child>
            <child>
              <object class="GtkShortcutsShortcut">
                <property name="visible">1</property>
                <property name="accelerator">&lt;ctrl&gt;s</property>
                <property name="title" translatable="yes">Save</property>
              </object>
            </child>
            <child>
              <object class="GtkShortcutsShortcut">
                <property name="visible">1</property>
                <property name="accelerator">&lt;ctrl&gt;d</property>
                <property name="title" translatable="yes">Delete</property>
              </object>
            </child>
          </object>
        </child>
      </object>
    </child>
  </object>
</interface>
`

//ShowShortcutsDialog shows the applications shortcuts dialog.
func ShowShortcutsDialog() {
	builder, _ := gtk.BuilderNew()
	builder.AddFromString(shortcutsDialogUIFile)

	showErrorDialog := func(errorMessage string) {
		message := fmt.Sprintf("Error showing shortcuts dialog (%s).", errorMessage)
		dialog := gtk.MessageDialogNew(nil, gtk.DIALOG_DESTROY_WITH_PARENT, gtk.MESSAGE_ERROR, gtk.BUTTONS_OK, message)
		dialog.Run()
		dialog.Destroy()
	}

	window, builderError := builder.GetObject("shortcuts-uwunote")
	if builderError != nil {
		showErrorDialog(builderError.Error())
		return
	}

	windowCast, ok := window.(*gtk.ShortcutsWindow)
	if ok {
		windowCast.SetResizable(true)
		windowCast.ShowAll()
	} else {
		showErrorDialog("Casting to ShortcutsWindow failed.")
	}
}
