package mainwindow

import (
	"znth/audio"
	"znth/components"
	"znth/ui"
	"znth/ui/toolbar"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/canvas"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/driver/desktop"
	"fyne.io/fyne/v2/theme"
	"fyne.io/fyne/v2/widget"
)

// Can these be moved?
var playAction *widget.ToolbarAction
var setlist *widget.List
var mixers = container.NewHBox()

var info = widget.NewLabel("Nothing loaded...") // global?

var currentSongPath string = ""

func createLayout(ui ui.UI) *fyne.Container {

	// Borders
	toolbarBorder := canvas.NewLine(theme.Color(theme.ColorNameSeparator))
	statusBorder := canvas.NewLine(theme.Color(theme.ColorNameSeparator))

	toolbarBorder.StrokeWidth = 1
	statusBorder.StrokeWidth = 1

	// Toolbar
	toolbar := container.NewBorder(
		nil,
		toolbarBorder,
		nil,
		nil,
		toolbar.CreateToolbar(ui),
	)

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
		createSetlist(),
	)
}

func createSetlist() *container.Split {
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
	return split
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

func addShortCuts(ui ui.UI) {
	ctrlS_Shortcut := &desktop.CustomShortcut{
		KeyName:  fyne.KeyS,
		Modifier: fyne.KeyModifierControl,
	}

	ui.Window.Canvas().AddShortcut(ctrlS_Shortcut, func(shortcut fyne.Shortcut) {
		if currentSongPath != "" {
			components.SaveStemData(currentSongPath)
			info.SetText("Saved stem levels!")
		}
	})

	ui.Window.Canvas().SetOnTypedKey(func(ev *fyne.KeyEvent) {
		switch ev.Name {
		case fyne.KeySpace:
			play()
		}
	})
}

func play() {
	if audio.IsStreamActive() {
		components.PlayAction()
		if audio.IsPlaying() {
			playAction.SetIcon(theme.MediaPauseIcon())
		} else {
			playAction.SetIcon(theme.MediaPlayIcon())
		}
	}
}
