package main

import (
	"znth/audio"
	"znth/components"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/canvas"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/theme"
	"fyne.io/fyne/v2/widget"
	"github.com/gordonklaus/portaudio"
)

var projectFolder string
var info = widget.NewLabel("Nothing loaded...")

var mixers = container.NewHBox()

var w fyne.Window

func main() {

	// App Init
	a := app.NewWithID("com.example.znth")

	// Setup audio
	portaudio.Initialize()
	defer portaudio.Terminate()
	defer audio.KillStream()

	// Create window
	w = a.NewWindow("Backing Track")
	w.SetContent(createLayout())
	w.CenterOnScreen()
	w.RequestFocus()
	w.Resize(fyne.NewSize(500, 500))
	w.ShowAndRun()
}

func createLayout() *fyne.Container {

	// Borders
	toolbarBorder := canvas.NewLine(theme.SeparatorColor())
	statusBorder := canvas.NewLine(theme.SeparatorColor())

	toolbarBorder.StrokeWidth = 1
	statusBorder.StrokeWidth = 1

	// Toolbar
	toolbar := container.NewBorder(
		nil,
		toolbarBorder,
		nil,
		nil,
		components.CreateToolbar(w, mixers),
	)

	// Setlist data
	songNames := []string{
		"Song One",
		"Song Two",
		"Very Long Song Name That Could Overflow",
		"Another Song",
		"Song Five",
		"Song Six",
		"Song Seven",
		"Song Eight",
	}

	// Setlist
	setlist := widget.NewList(
		func() int {
			return len(songNames)
		},
		func() fyne.CanvasObject {
			label := widget.NewLabel("")
			label.Wrapping = fyne.TextTruncate
			return label
		},
		func(i widget.ListItemID, o fyne.CanvasObject) {
			o.(*widget.Label).SetText(songNames[i])
		},
	)

	// Status bar
	statusBar := container.NewBorder(
		statusBorder,
		nil,
		nil,
		nil,
		info,
	)

	// Split view
	split := container.NewHSplit(
		setlist,
		container.NewScroll(mixers),
	)

	// Start with left panel at 25% width
	split.SetOffset(0.25)

	// Full window layout
	return container.NewBorder(
		toolbar,
		statusBar,
		nil,
		nil,
		split,
	)
}
