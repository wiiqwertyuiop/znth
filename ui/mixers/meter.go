package mixers

import (
	"image/color"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/canvas"
	"fyne.io/fyne/v2/container"
)

func createMeterTicks() fyne.CanvasObject {
	ticks := []float32{
		dbPosition(6),
		dbPosition(0),
		dbPosition(-12),
		dbPosition(-24),
		dbPosition(-48),
		dbPosition(-60),
	}

	objects := make([]fyne.CanvasObject, 0, len(ticks))

	for range ticks {
		line := canvas.NewLine(color.RGBA{R: 255, G: 255, B: 255, A: 100})
		line.StrokeWidth = 3
		objects = append(objects, line)
	}

	return container.New(&tickLayout{
		positions: ticks,
	}, objects...)
}

type tickLayout struct {
	positions []float32
}

func (l *tickLayout) Layout(objects []fyne.CanvasObject, size fyne.Size) {
	for i, obj := range objects {
		y := size.Height * l.positions[i]

		line := obj.(*canvas.Line)
		line.Position1 = fyne.NewPos(0, y)
		line.Position2 = fyne.NewPos(size.Width, y)
		line.Refresh()
	}
}

func (l *tickLayout) MinSize(objects []fyne.CanvasObject) fyne.Size {
	return fyne.NewSize(10, 0)
}

func dbPosition(db float32) float32 {
	return 1 - ((db + 60) / 66)
}
