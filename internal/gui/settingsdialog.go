package gui

import (
	"math"

	"github.com/gotk3/gotk3/glib"

	"github.com/gotk3/gotk3/gdk"
	"github.com/gotk3/gotk3/gtk"

	"github.com/UwUNote/uwunote/internal/config"
	"github.com/UwUNote/uwunote/internal/globconstants"
)

type settingsDialogContainer struct {
	window             *gtk.Window
	initializeFunction func(appConfigToUse *config.AppConfig)
}

var settingsDialog *settingsDialogContainer

//ShowSettingsDialog shows a settingsdialog and saves them on close
func ShowSettingsDialog() {
	if settingsDialog != nil && settingsDialog.window != nil && settingsDialog.window.IsVisible() {
		settingsDialog.window.Present()
		return
	}

	appConfig := config.GetAppConfigCopy()

	//All GTKErrors are ignored for now

	//DefaultNoteX
	defaultNoteXLabel, _ := gtk.LabelNew("Default horizontal note position")

	defaultNoteXLabel.SetHAlign(gtk.ALIGN_START)
	defaultNoteXLabel.SetTooltipText("Default horizontal position for non relative positioned note windows.")

	defaultNoteXSpinner, _ := gtk.SpinButtonNewWithRange(math.MinInt32, math.MaxInt32, 5)
	expandAndAlignRight(&defaultNoteXSpinner.Widget)

	//DefaultNoteY
	defaultNoteYLabel, _ := gtk.LabelNew("Default vertical note position")

	defaultNoteYLabel.SetHAlign(gtk.ALIGN_START)
	defaultNoteYLabel.SetTooltipText("Default vertical position for non relative positioned note windows.")

	defaultNoteYSpinner, _ := gtk.SpinButtonNewWithRange(math.MinInt32, math.MaxInt32, 5)
	expandAndAlignRight(&defaultNoteYSpinner.Widget)

	//DefaultNoteWidth
	defaultNoteWidthLabel, _ := gtk.LabelNew("Default note width")

	defaultNoteWidthLabel.SetHAlign(gtk.ALIGN_START)
	defaultNoteWidthLabel.SetTooltipText("Default width for note windows.")

	defaultNoteWidthSpinner, _ := gtk.SpinButtonNewWithRange(math.MinInt32, math.MaxInt32, 5)
	expandAndAlignRight(&defaultNoteWidthSpinner.Widget)

	//DefaultNoteHeight
	defaultNoteHeightLabel, _ := gtk.LabelNew("Default note height")

	defaultNoteHeightLabel.SetHAlign(gtk.ALIGN_START)
	defaultNoteHeightLabel.SetTooltipText("Default height for note windows.")

	defaultNoteHeightSpinner, _ := gtk.SpinButtonNewWithRange(math.MinInt32, math.MaxInt32, 5)
	expandAndAlignRight(&defaultNoteHeightSpinner.Widget)

	//appearanceSettings Tab
	appearanceSettingsTab, _ := gtk.GridNew()
	setupTab(appearanceSettingsTab)

	appearanceSettingsTab.Add(defaultNoteXLabel)
	appearanceSettingsTab.AttachNextTo(defaultNoteXSpinner, defaultNoteXLabel, gtk.POS_RIGHT, 1, 1)

	appearanceSettingsTab.AttachNextTo(defaultNoteYLabel, defaultNoteXLabel, gtk.POS_BOTTOM, 1, 1)
	appearanceSettingsTab.AttachNextTo(defaultNoteYSpinner, defaultNoteYLabel, gtk.POS_RIGHT, 1, 1)

	appearanceSettingsTab.AttachNextTo(defaultNoteWidthLabel, defaultNoteYLabel, gtk.POS_BOTTOM, 1, 1)
	appearanceSettingsTab.AttachNextTo(defaultNoteWidthSpinner, defaultNoteWidthLabel, gtk.POS_RIGHT, 1, 1)

	appearanceSettingsTab.AttachNextTo(defaultNoteHeightLabel, defaultNoteWidthLabel, gtk.POS_BOTTOM, 1, 1)
	appearanceSettingsTab.AttachNextTo(defaultNoteHeightSpinner, defaultNoteHeightLabel, gtk.POS_RIGHT, 1, 1)

	//NoteDirectory
	noteDirectoryLabel, _ := gtk.LabelNew("Note directory")
	noteDirectoryPicker, _ := gtk.FileChooserButtonNew("Test", gtk.FILE_CHOOSER_ACTION_SELECT_FOLDER)

	noteDirectoryLabel.SetHAlign(gtk.ALIGN_START)
	noteDirectoryToolTip := "Determines the folder in which all notes will be saved."
	noteDirectoryLabel.SetTooltipText(noteDirectoryToolTip)
	noteDirectoryPicker.SetTooltipText(noteDirectoryToolTip)
	expandAndAlignRight(&noteDirectoryPicker.Widget)

	//AskBeforeNoteDeletion
	askBeforeNoteDeletionLabel, _ := gtk.LabelNew("Confirm note deletion")
	askBeforeNoteDeletionSwitch, _ := gtk.SwitchNew()

	askBeforeNoteDeletionLabel.SetHAlign(gtk.ALIGN_START)
	askBeforeNoteDeletionToolTip := "Shows a dialog before deleting a note, to make a sure you don't accidentally delete a note."
	askBeforeNoteDeletionLabel.SetTooltipText(askBeforeNoteDeletionToolTip)
	askBeforeNoteDeletionSwitch.SetTooltipText(askBeforeNoteDeletionToolTip)
	expandAndAlignRight(&askBeforeNoteDeletionSwitch.Widget)

	//DeleteNotesToTrashbin
	deleteNotesToTrashbinLabel, _ := gtk.LabelNew("Use system trashbin")
	deleteNotesToTrashbinSwitch, _ := gtk.SwitchNew()

	deleteNotesToTrashbinLabel.SetHAlign(gtk.ALIGN_START)
	deleteNotesToTrashbinToolTip := "Decides wether the systems trashbin will be used, this makes notes recoverable."
	deleteNotesToTrashbinLabel.SetTooltipText(deleteNotesToTrashbinToolTip)
	deleteNotesToTrashbinSwitch.SetTooltipText(deleteNotesToTrashbinToolTip)
	expandAndAlignRight(&deleteNotesToTrashbinSwitch.Widget)

	//ShowTrayIcon
	showTrayIconLabel, _ := gtk.LabelNew("Show tray icon")
	showTrayIconSwitch, _ := gtk.SwitchNew()

	showTrayIconLabel.SetHAlign(gtk.ALIGN_START)
	showTrayIconToolTip := "Shows a tray icon in the systems tray area."
	showTrayIconLabel.SetTooltipText(showTrayIconToolTip)
	showTrayIconSwitch.SetTooltipText(showTrayIconToolTip)
	expandAndAlignRight(&showTrayIconSwitch.Widget)

	//GeneralSettings Tab
	generalSettingsTab, _ := gtk.GridNew()
	setupTab(generalSettingsTab)

	generalSettingsTab.Add(askBeforeNoteDeletionLabel)
	generalSettingsTab.AttachNextTo(askBeforeNoteDeletionSwitch, askBeforeNoteDeletionLabel, gtk.POS_RIGHT, 1, 1)

	generalSettingsTab.AttachNextTo(deleteNotesToTrashbinLabel, askBeforeNoteDeletionLabel, gtk.POS_BOTTOM, 1, 1)
	generalSettingsTab.AttachNextTo(deleteNotesToTrashbinSwitch, deleteNotesToTrashbinLabel, gtk.POS_RIGHT, 1, 1)

	generalSettingsTab.AttachNextTo(noteDirectoryLabel, deleteNotesToTrashbinLabel, gtk.POS_BOTTOM, 1, 1)
	generalSettingsTab.AttachNextTo(noteDirectoryPicker, noteDirectoryLabel, gtk.POS_RIGHT, 1, 1)

	//AutoSaveDelay
	autoSaveDelayLabel, _ := gtk.LabelNew("Autosave delay")

	autoSaveDelayLabel.SetHAlign(gtk.ALIGN_START)
	autoSaveDelayLabel.SetTooltipText("Time to wait before automatically saving the edited note after typing.")

	autoSaveDelaySpinner, _ := gtk.SpinButtonNewWithRange(0, math.MaxFloat64, 50)
	expandAndAlignRight(&autoSaveDelaySpinner.Widget)

	//AutoSave
	autoSaveLabel, _ := gtk.LabelNew("Autosave")
	autoSaveSwitch, _ := gtk.SwitchNew()

	autoSaveLabel.SetHAlign(gtk.ALIGN_START)
	autoSaveLabel.SetTooltipText("Automatically saves the content of the edited note after typing.")

	autoSaveSwitch.ConnectAfter("notify::active", func(button *gtk.Switch) {
		autoSaveDelaySpinner.SetSensitive(button.GetActive())
	})
	expandAndAlignRight(&autoSaveSwitch.Widget)

	//AutoIndent
	autoIndentLabel, _ := gtk.LabelNew("Autoindent")
	autoIndentSwitch, _ := gtk.SwitchNew()

	autoIndentLabel.SetHAlign(gtk.ALIGN_START)
	autoIndentToolTip := "Automatically indents a new line by the same amount of tabs, that the previous line was indented with."
	autoIndentLabel.SetTooltipText(autoIndentToolTip)
	autoIndentSwitch.SetTooltipText(autoIndentToolTip)
	expandAndAlignRight(&autoIndentSwitch.Widget)

	//WrapMode
	wrapModeLabel, _ := gtk.LabelNew("Textwrapping")

	const (
		columnID = iota
		columnText
	)
	wrapModeItems, _ := gtk.ListStoreNew(glib.TYPE_INT, glib.TYPE_STRING)
	addItem := func(wrapMode gtk.WrapMode, text string) {

		appendIter := wrapModeItems.Append()
		wrapModeItems.SetValue(appendIter, columnID, wrapMode)
		wrapModeItems.SetValue(appendIter, columnText, text)
	}

	addItem(gtk.WRAP_NONE, "None")
	addItem(gtk.WRAP_CHAR, "Character")
	addItem(gtk.WRAP_WORD, "Word")
	addItem(gtk.WRAP_WORD_CHAR, "Wordcharacters")

	wrapModeComboBox, _ := gtk.ComboBoxNewWithModel(wrapModeItems)

	wrapModeRenderer, _ := gtk.CellRendererTextNew()
	wrapModeComboBox.PackStart(wrapModeRenderer, true)
	wrapModeComboBox.AddAttribute(wrapModeRenderer, "text", 1)
	expandAndAlignRight(&wrapModeComboBox.Widget)

	wrapModeLabel.SetHAlign(gtk.ALIGN_START)
	const wrapModeToolTipText = "Determines wether the editor breaks the text if it out of bounds and how the text will be broken."
	wrapModeLabel.SetTooltipText(wrapModeToolTipText)
	wrapModeComboBox.SetTooltipText(wrapModeToolTipText)

	//EditorSettings Tab
	editorSettingsTab, _ := gtk.GridNew()
	setupTab(editorSettingsTab)

	editorSettingsTab.Add(autoSaveLabel)
	editorSettingsTab.AttachNextTo(autoSaveSwitch, autoSaveLabel, gtk.POS_RIGHT, 1, 1)

	editorSettingsTab.AttachNextTo(autoSaveDelayLabel, autoSaveLabel, gtk.POS_BOTTOM, 1, 1)
	editorSettingsTab.AttachNextTo(autoSaveDelaySpinner, autoSaveDelayLabel, gtk.POS_RIGHT, 1, 1)

	editorSettingsTab.AttachNextTo(autoIndentLabel, autoSaveDelayLabel, gtk.POS_BOTTOM, 1, 1)
	editorSettingsTab.AttachNextTo(autoIndentSwitch, autoIndentLabel, gtk.POS_RIGHT, 1, 1)

	editorSettingsTab.AttachNextTo(wrapModeLabel, autoIndentLabel, gtk.POS_BOTTOM, 1, 1)
	editorSettingsTab.AttachNextTo(wrapModeComboBox, wrapModeLabel, gtk.POS_RIGHT, 1, 1)

	//TabContainer
	settingsTabContainer, _ := gtk.NotebookNew()
	settingsTabContainer.SetVExpand(true)
	settingsTabContainer.SetHExpand(true)

	generalSettingsTabLabel, _ := gtk.LabelNew("General")
	settingsTabContainer.AppendPage(generalSettingsTab, generalSettingsTabLabel)

	editorSettingsTabLabel, _ := gtk.LabelNew("Editor")
	settingsTabContainer.AppendPage(editorSettingsTab, editorSettingsTabLabel)

	appearanceSettingsTabLabel, _ := gtk.LabelNew("Appearance")
	settingsTabContainer.AppendPage(appearanceSettingsTab, appearanceSettingsTabLabel)

	initializeFunction := func(appConfigToUse *config.AppConfig) {
		deleteNotesToTrashbinSwitch.SetActive(appConfigToUse.DeleteNotesToTrashbin)
		askBeforeNoteDeletionSwitch.SetActive(appConfigToUse.AskBeforeNoteDeletion)
		noteDirectoryPicker.SetFilename(appConfigToUse.NoteDirectory)
		showTrayIconSwitch.SetActive(appConfigToUse.ShowTrayIcon)

		autoIndentSwitch.SetActive(appConfigToUse.AutoIndent)
		autoSaveSwitch.SetActive(appConfigToUse.AutoSaveAfterTyping)
		autoSaveDelaySpinner.SetValue(float64(appConfigToUse.AutoSaveAfterTypingDelay))
		autoSaveDelaySpinner.SetSensitive(appConfigToUse.AutoSaveAfterTyping)
		wrapModeComboBox.SetActive(int(appConfigToUse.WrapMode))

		defaultNoteXSpinner.SetValue(float64(appConfigToUse.DefaultNoteX))
		defaultNoteYSpinner.SetValue(float64(appConfigToUse.DefaultNoteY))
		defaultNoteWidthSpinner.SetValue(float64(appConfigToUse.DefaultNoteWidth))
		defaultNoteHeightSpinner.SetValue(float64(appConfigToUse.DefaultNoteHeight))
	}

	initializeFunction(&appConfig)

	//ResetToDefaults Button
	resetToDefaultsButton, _ := gtk.ButtonNewWithLabel("Reset to defaults")
	resetToDefaultsButton.Connect("clicked", func() {
		defaultAppConfig := config.GetAppConfigDefaults()
		initializeFunction(&defaultAppConfig)
	})

	//Windowsetup
	settingsWindow, _ := gtk.WindowNew(gtk.WINDOW_TOPLEVEL)
	settingsWindow.SetTitle(globconstants.ApplicationName + " - Settings")

	pixbufLoader, _ := gdk.PixbufLoaderNew()
	iconAsPixbuf, _ := pixbufLoader.WriteAndReturnPixbuf(AppIcon)
	settingsWindow.SetIcon(iconAsPixbuf)

	settingsWindow.SetSkipTaskbarHint(true)
	settingsWindow.SetSkipPagerHint(true)

	settingsWindow.SetPosition(gtk.WIN_POS_CENTER)

	//Save on close
	settingsWindow.Connect("destroy", func() {
		//GeneralSettings
		appConfig.AskBeforeNoteDeletion = deleteNotesToTrashbinSwitch.GetActive()
		appConfig.AskBeforeNoteDeletion = askBeforeNoteDeletionSwitch.GetActive()
		appConfig.NoteDirectory = noteDirectoryPicker.GetFilename()
		appConfig.ShowTrayIcon = showTrayIconSwitch.GetActive()

		//EditorSettings
		appConfig.AutoSaveAfterTyping = autoSaveSwitch.GetActive()

		currentAutoSaveDelaySpinnerValue := autoSaveDelaySpinner.GetValueAsInt()
		if currentAutoSaveDelaySpinnerValue > math.MaxInt32 {
			appConfig.AutoSaveAfterTypingDelay = math.MaxInt32
		} else if currentAutoSaveDelaySpinnerValue < 0 {
			appConfig.AutoSaveAfterTypingDelay = 0
		} else {
			appConfig.AutoSaveAfterTypingDelay = currentAutoSaveDelaySpinnerValue
		}

		appConfig.AutoIndent = autoIndentSwitch.GetActive()
		appConfig.WrapMode = gtk.WrapMode(wrapModeComboBox.GetActive())

		//AppearanceSettings
		appConfig.DefaultNoteX = defaultNoteXSpinner.GetValueAsInt()
		appConfig.DefaultNoteY = defaultNoteYSpinner.GetValueAsInt()
		appConfig.DefaultNoteWidth = defaultNoteWidthSpinner.GetValueAsInt()
		appConfig.DefaultNoteHeight = defaultNoteHeightSpinner.GetValueAsInt()

		config.PersistAppConfig(&appConfig)
		config.LoadAppConfig()
	})

	settingsPanel, _ := gtk.BoxNew(gtk.ORIENTATION_VERTICAL, 0)
	settingsPanel.Add(settingsTabContainer)
	settingsPanel.Add(resetToDefaultsButton)

	settingsWindow.Add(settingsPanel)
	settingsWindow.ShowAll()

	settingsDialog = &settingsDialogContainer{
		window:             settingsWindow,
		initializeFunction: initializeFunction,
	}
}

func expandAndAlignRight(widget *gtk.Widget) {
	widget.SetHAlign(gtk.ALIGN_END)
	widget.SetHExpand(true)
}

func setupTab(tab *gtk.Grid) {
	tab.SetColumnSpacing(30)
	tab.SetRowSpacing(5)

	const margin = 10

	tab.SetMarginBottom(margin)
	tab.SetMarginEnd(margin)
	tab.SetMarginStart(margin)
	tab.SetMarginTop(margin)
}
