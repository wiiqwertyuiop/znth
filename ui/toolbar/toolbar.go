package toolbar

import (
	"image/color"
	"znth/components"
	"znth/ui"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/canvas"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/dialog"
	"fyne.io/fyne/v2/theme"
	"fyne.io/fyne/v2/widget"
)

func Create(ui ui.UI) *fyne.Container {

	saveSetlistButton := widget.NewButtonWithIcon("", theme.DocumentSaveIcon(), func() {
		dialog.NewFileSave(func(writer fyne.URIWriteCloser, err error) {
			if err != nil {
				ui.State.StatusBarTextChange("Save error: " + err.Error())
				return
			}

			if err := components.SaveSetlist(writer, ui.State); err != nil {
				ui.State.StatusBarTextChange(err.Error())
			}
		}, ui.Window).Show()
	})

	openSetlistButton := widget.NewButtonWithIcon("", theme.FileApplicationIcon(), func() {
		dialog.NewFileOpen(func(reader fyne.URIReadCloser, err error) {
			components.OpenSetlist(reader, ui.State)
		}, ui.Window).Show()
	})

	openSetlistButton.Importance = widget.LowImportance
	saveSetlistButton.Importance = widget.LowImportance

	leftControls := container.NewHBox(
		openFolderButton(ui),
		openSetlistButton,
		saveSetlistButton,
		verticalSeparator(),
		playButton(ui.State),
		stopButton(ui.State),
	)

	return container.NewBorder(
		nil,                    // top
		nil,                    // bottom
		leftControls,           // left
		nil,                    // right
		createSeeker(ui.State), // center
	)
}

func verticalSeparator() fyne.CanvasObject {
	rect := canvas.NewRectangle(color.White)
	rect.SetMinSize(fyne.NewSize(1, 24))

	return rect
}
