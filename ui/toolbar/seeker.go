package toolbar

import (
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/widget"
)

func createSeeker() fyne.CanvasObject {
	slider := widget.NewSlider(0, 100)

	return container.NewBorder(
		nil,
		nil,
		container.NewCenter(widget.NewLabel("0:00")),
		container.NewCenter(widget.NewLabel("3:42")),
		slider,
	)
}
