package mixers

import (
	"image/color"
	"math"
	"strings"
	"znth/audio"
	"znth/model"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/canvas"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/widget"
)

func drawInstrument(stem *model.Stem, project *model.Project) (*fyne.Container, func()) {

	defaultVolume := stem.VolumeAdjust

	volume := newVerticalSlider(0, 100, audio.GainToSlider(defaultVolume)*100)
	volume.OnChanged = func(v float64) {
		stem.VolumeAdjust = audio.SliderToGain(v / 100.0)
	}

	project.AddToProjectRenderLoop(func() {
		peak := math.Float32frombits(stem.Peak.Load())
		fyne.Do(func() {
			volume.UpdatePeak(peak)
		})
	})

	border := canvas.NewRectangle(color.NRGBA{
		R: 60,
		G: 60,
		B: 60,
		A: 255,
	})
	border.StrokeColor = color.Black
	border.StrokeWidth = 1

	// Determine track color
	bg := canvas.NewRectangle(determineTrackColor(stem.Id))

	labelText := widget.NewLabel(strings.TrimSuffix(stem.Id, ".wav"))
	labelText.TextStyle.Bold = true
	label := container.NewPadded(labelText)

	nameLabel := container.NewStack(
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
		), func() {
			// Clean up function
			volume.OnChanged = nil
			volume = nil
		}

}

func determineTrackColor(name string) color.Color {
	name = strings.ToLower(name)
	if strings.Contains(name, "bass") {
		return color.NRGBA{
			R: 128,
			G: 0,
			B: 128,
			A: 200,
		}
	}

	if strings.Contains(name, "guitar") {
		return color.NRGBA{
			R: 0,
			G: 0,
			B: 128,
			A: 200,
		}
	}

	if strings.Contains(name, "drums") {
		return color.NRGBA{
			R: 255,
			G: 255,
			B: 0,
			A: 100,
		}
	}

	if strings.Contains(name, "vocals") {
		return color.NRGBA{
			R: 255,
			G: 165,
			B: 0,
			A: 200,
		}
	}

	if strings.Contains(name, "synth") {
		return color.NRGBA{
			R: 144,
			G: 238,
			B: 144,
			A: 144,
		}
	}

	return color.NRGBA{
		R: 255,
		G: 255,
		B: 255,
		A: 100,
	}
}
