package components

import (
	"fmt"
	"image/color"
	"os"
	"path/filepath"
	"strings"
	"znth/audio"
	"znth/model"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/canvas"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/dialog"
	"fyne.io/fyne/v2/theme"
	"fyne.io/fyne/v2/widget"
)

var playing = false
var stems []model.Stem
var playAction *widget.ToolbarAction

func CreateToolbar(w fyne.Window, mixers *fyne.Container) *widget.Toolbar {
	playAction = widget.NewToolbarAction(theme.MediaPlayIcon(), func() {
		if !audio.IsStreamActive() {
			audio.StartStream(stems)
		}
		if !playing {
			audio.Play()
			playAction.SetIcon(theme.MediaPauseIcon())
			playing = true
		} else {
			audio.Pause()
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

				if audio.IsStreamActive() {
					audio.KillStream()
					mixers.RemoveAll()
					stems = nil
					audio.SetMusicPosition(0)
					playAction.SetIcon(theme.MediaPlayIcon())
				}
				stems = loadProjectFolder(reader.Path(), mixers)
				//info.SetText("Loaded " + reader.Path())
			}, w).Show()
		}),
		widget.NewToolbarAction(theme.FileApplicationIcon(), func() {}),
		widget.NewToolbarAction(theme.DocumentSaveIcon(), func() {}),
		widget.NewToolbarSeparator(),
		playAction,
		widget.NewToolbarAction(theme.MediaStopIcon(), func() {
			if playing {
				// MAKE SHARED
				audio.Pause()
				playAction.SetIcon(theme.MediaPlayIcon())
				playing = false
			}
			audio.SetMusicPosition(0)
		}),
	)
	return toolbar
}

func loadProjectFolder(folder string, mixers *fyne.Container) []model.Stem {
	files, err := os.ReadDir(folder)
	if err != nil {
		panic(err)
	}

	// Master Volume
	volume := NewVerticalSlider(0, 100, 20)
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

			volume := NewVerticalSlider(0, 100, 50)
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
