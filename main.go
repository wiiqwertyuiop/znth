package main

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"znth/audio"
	"znth/model"
	"znth/ui"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
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

func main() {

	a := app.NewWithID("com.example.znth")

	portaudio.Initialize()

	defer portaudio.Terminate()
	defer func() {
		if stream != nil {
			stream.Stop()
			stream.Close()
			stream = nil
		}
	}()

	w := a.NewWindow("Backing Track")
	w.Resize(fyne.NewSize(800, 800))

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
				}
				stems = loadProjectFolder(reader.Path())
				info.SetText("Loaded " + reader.Path())
			}, w).Show()
		}),
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

	center := container.NewVBox(
		info,
		mixers,
	)

	scroll := container.NewScroll(center)
	content := container.NewBorder(toolbar, nil, nil, nil, scroll)
	w.SetContent(content)
	w.ShowAndRun()
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

	channel := container.NewVBox(
		volume,
		widget.NewLabel("Master"),
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

			channel := container.NewVBox(
				volume,
				widget.NewLabel(file.Name()),
			)

			mixers.Add(channel)

			stems = append(stems, stem)
		}
	}
	mixers.Refresh()
	return stems
}
