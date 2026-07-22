package main

import (
	"fmt"
	"image/color"
	"os"
	"path/filepath"
	"strings"
	"znth/audio"
	"znth/model"
	"znth/ui"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/canvas"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/dialog"
	"fyne.io/fyne/v2/theme"
	"fyne.io/fyne/v2/widget"
	"github.com/gordonklaus/portaudio"
)

var projectFolder string
var stream *portaudio.Stream = nil
var info = widget.NewLabel("Nothing loaded...")
var playing = false
var playAction *widget.ToolbarAction

var stems []model.Stem

var mixers = container.NewHBox()

var w fyne.Window

func main() {

	// App Init
	a := app.NewWithID("com.example.znth")

	// Setup audio
	portaudio.Initialize()
	defer portaudio.Terminate()
	defer func() {
		if stream != nil {
			stream.Stop()
			stream.Close()
			stream = nil
		}
	}()

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
		createToolbar(),
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

func createToolbar() *widget.Toolbar {
	playAction = widget.NewToolbarAction(theme.MediaPlayIcon(), func() {
		if stream == nil {
			stream = audio.StartStream(stems)
		}
		if !playing {
			stream.Start()
			playAction.SetIcon(theme.MediaPauseIcon())
			playing = true
		} else {
			stream.Stop()
			playAction.SetIcon(theme.MediaPlayIcon())
			playing = false
		}
	})

	toolbar := widget.NewToolbar(
		widget.NewToolbarAction(theme.FolderOpenIcon(), func() {
			dialog.NewFolderOpen(func(reader fyne.ListableURI, err error) {

				if err != nil {
					fmt.Println("Folder selection error:", err)
					return
				}

				if reader == nil {
					fmt.Println("Folder selection cancelled")
					return
				}

				if stream != nil {
					stream.Stop()
					stream.Close()
					stream = nil
					mixers.RemoveAll()
					stems = nil
					audio.SetMusicPosition(0)
					playAction.SetIcon(theme.MediaPlayIcon())
				}
				stems = loadProjectFolder(reader.Path())
				info.SetText("Loaded " + reader.Path())
			}, w).Show()
		}),
		widget.NewToolbarAction(theme.FileApplicationIcon(), func() {}),
		widget.NewToolbarAction(theme.DocumentSaveIcon(), func() {}),
		widget.NewToolbarSeparator(),
		playAction,
		widget.NewToolbarAction(theme.MediaStopIcon(), func() {
			if playing {
				// MAKE SHARED
				stream.Stop()
				playAction.SetIcon(theme.MediaPlayIcon())
				playing = false
			}
			audio.SetMusicPosition(0)
		}),
	)
	return toolbar
}

func loadProjectFolder(folder string) []model.Stem {
	files, err := os.ReadDir(folder)
	if err != nil {
		panic(err)
	}

	// Master Volume
	volume := ui.NewVerticalSlider(0, 100, 20)
	volume.OnChanged = func(v float64) {
		audio.SetMasterVolume(float32(v / 50))
	}

	border := canvas.NewRectangle(color.NRGBA{
		R: 60,
		G: 60,
		B: 60,
		A: 255,
	})
	border.StrokeColor = color.Black
	border.StrokeWidth = 1

	bg := canvas.NewRectangle(color.NRGBA{
		R: 144,
		G: 0,
		B: 0,
		A: 144,
	})

	label := container.NewPadded(widget.NewLabel("Master"))

	nameLabel := container.NewMax(
		bg,
		label,
	)

	channel := container.NewStack(
		border,
		container.NewBorder(
			nil,
			nameLabel,
			nil,
			nil,
			volume,
		),
	)

	mixers.Add(channel)

	var stems []model.Stem
	for _, file := range files {
		if file.IsDir() {
			continue
		}

		if strings.ToLower(filepath.Ext(file.Name())) == ".wav" {
			fullPath := filepath.Join(folder, file.Name())
			println("Loading... ", fullPath)
			data, info, err := audio.LoadWavFloat32(fullPath)
			if err != nil {
				panic(err)
			}

			index := len(stems)

			volume := ui.NewVerticalSlider(0, 100, 50)
			stem := model.Stem{
				Data:         data,
				Info:         info,
				VolumeAdjust: 1,
			}

			volume.OnChanged = func(v float64) {
				stems[index].VolumeAdjust = float32(v / 50)
			}

			border := canvas.NewRectangle(color.NRGBA{
				R: 60,
				G: 60,
				B: 60,
				A: 255,
			})
			border.StrokeColor = color.Black
			border.StrokeWidth = 1

			bg := canvas.NewRectangle(color.NRGBA{
				R: 144,
				G: 238,
				B: 144,
				A: 144,
			})

			label := container.NewPadded(widget.NewLabel(file.Name()))

			nameLabel := container.NewMax(
				bg,
				label,
			)

			channel := container.NewStack(
				border,
				container.NewBorder(
					nil,
					nameLabel,
					nil,
					nil,
					volume,
				),
			)

			mixers.Add(channel)

			stems = append(stems, stem)
		}
	}
	mixers.Refresh()
	return stems
}
