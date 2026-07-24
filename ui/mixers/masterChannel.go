package mixers

import (
	"image/color"
	"znth/audio"
	"znth/model"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/canvas"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/widget"
)

func drawMaster(stem model.Stem) *fyne.Container {

	masterVolume := stem.VolumeAdjust

	// Master Volume
	volume := newVerticalSlider(0, 100, audio.GainToSlider(masterVolume)*100)
	volume.OnChanged = func(v float64) {
		audio.SetMasterVolume(audio.SliderToGain(v / 100.0))
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

	labelText := widget.NewLabel("Master")
	labelText.TextStyle.Bold = true
	label := container.NewPadded(labelText)

	nameLabel := container.NewMax(
		bg,
		label,
	)

	return container.NewStack(
		border,
		container.NewBorder(
			nil,
			nameLabel,
			createMeterTicks(),
			nil,
			volume,
		),
	)
}
