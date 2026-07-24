package mixers

import (
	"image/color"
	"math"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/canvas"
	"fyne.io/fyne/v2/container"
)

func createMeterTicks() fyne.CanvasObject {
	ticks := container.NewWithoutLayout()

	levels := []float32{10, 57, 105, 152, 200, 238, 277, 316, 355}

	for _, y := range levels {
		line := canvas.NewLine(color.RGBA{R: 255, G: 255, B: 255, A: 100})
		line.StrokeWidth = 3
		line.Position1 = fyne.NewPos(0, y)
		if math.Mod(float64(y), 2) != 0 {
			line.Position2 = fyne.NewPos(7, y)
		} else {
			line.Position2 = fyne.NewPos(12, y)
		}

		ticks.Add(line)
	}

	ticks.Resize(fyne.NewSize(10, 100))

	return ticks
}
