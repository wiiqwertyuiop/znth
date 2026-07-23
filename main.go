package main

import (
	"znth/audio"
	"znth/components"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/canvas"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/dialog"
	"fyne.io/fyne/v2/driver/desktop"
	"fyne.io/fyne/v2/theme"
	"fyne.io/fyne/v2/widget"
	"github.com/gordonklaus/portaudio"
)

var w fyne.Window

// Can these be moved?
var playAction *widget.ToolbarAction
var setlist *widget.List
var mixers = container.NewHBox()

var info = widget.NewLabel("Nothing loaded...") // global?

var currentSongPath string = ""

func main() {

	// App Init
	a := app.NewWithID("com.example.znth")

	// Setup audio
	portaudio.Initialize()
	defer portaudio.Terminate()
	defer audio.KillStream()

	// Create window
	w = a.NewWindow("Backing Track")
	addShortCuts()
	w.SetContent(createLayout())
	w.CenterOnScreen()
	w.RequestFocus()
	w.Resize(fyne.NewSize(500, 500))
	w.ShowAndRun()
}

func createLayout() *fyne.Container {

	// Borders
	toolbarBorder := canvas.NewLine(theme.SeparatorColor()) // TODO
	statusBorder := canvas.NewLine(theme.SeparatorColor())

	toolbarBorder.StrokeWidth = 1
	statusBorder.StrokeWidth = 1

	// Toolbar
	toolbar := container.NewBorder(
		nil,
		toolbarBorder,
		nil,
		nil,
		createToolbar(),
	)

	// Setlist
	setlist = components.CreateSetlist()
	setlist.OnSelected = func(id widget.ListItemID) {
		path := components.GetSong(id).Location
		loadSong(path)
	}

	// Split view
	split := container.NewHSplit(
		setlist,
		container.NewScroll(mixers),
	)

	// Start with left panel at 20% width
	split.SetOffset(0.20)

	// Status bar
	statusBar := container.NewBorder(
		statusBorder,
		nil,
		nil,
		nil,
		info,
	)

	// Full window layout
	return container.NewBorder(
		toolbar,
		statusBar,
		nil,
		nil,
		split,
	)
}

func createToolbar() *widget.Toolbar {
	playAction = widget.NewToolbarAction(theme.MediaPlayIcon(), func() {
		components.PlayAction()
		if audio.IsPlaying() {
			playAction.SetIcon(theme.MediaPauseIcon())
		} else {
			playAction.SetIcon(theme.MediaPlayIcon())
		}
	})

	return widget.NewToolbar(
		widget.NewToolbarAction(theme.FolderOpenIcon(), func() {
			dialog.NewFolderOpen(func(reader fyne.ListableURI, err error) {

				if err != nil {
					info.SetText("Folder selection error: " + err.Error())
					return
				}

				if reader == nil {
					info.SetText("Folder selection cancelled")
					return
				}
				loadSong(reader.Path())
				components.AddToSetlist(reader)
				setlist.Refresh()
			}, w).Show()
		}),
		widget.NewToolbarAction(theme.FileApplicationIcon(), func() {
			dialog.NewFileOpen(func(reader fyne.URIReadCloser, err error) {
				components.OpenSetlist(reader)
				setlist.Refresh()
			}, w).Show()
		}),
		widget.NewToolbarAction(theme.DocumentSaveIcon(), func() {
			dialog.NewFileSave(func(writer fyne.URIWriteCloser, err error) {
				if err != nil {
					info.SetText("Save error: " + err.Error())
					return
				}
				e := components.SaveSetlist(writer) // TODO
				if e != nil {
					info.SetText(e.Error())
				}
			}, w).Show()
		}),
		widget.NewToolbarSeparator(),
		playAction,
		widget.NewToolbarAction(theme.MediaStopIcon(), func() {
			components.StopAction()
			playAction.SetIcon(theme.MediaPlayIcon())
		}),
	)
}

func loadSong(path string) {
	audio.KillStream()
	mixers.RemoveAll()
	playAction.SetIcon(theme.MediaPlayIcon())

	info.SetText("Loading files... " + path)

	components.LoadProjectFolder(path, mixers) // Todo
	currentSongPath = path

	info.SetText("Loaded succesfully! " + path)
}

func addShortCuts() {
	shortcut := &desktop.CustomShortcut{
		KeyName:  fyne.KeyS,
		Modifier: fyne.KeyModifierControl,
	}

	w.Canvas().AddShortcut(shortcut, func(shortcut fyne.Shortcut) {
		if currentSongPath != "" {
			components.SaveStemData(currentSongPath)
			info.SetText("Saved stem levels!")
		}
	})
}
