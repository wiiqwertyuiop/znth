package toolbar

import (
	"znth/audio"
	"znth/state"

	"fyne.io/fyne/v2/theme"
	"fyne.io/fyne/v2/widget"
)

func stopButton(state *state.State) *widget.Button {
	button := widget.NewButtonWithIcon("", theme.MediaStopIcon(), func() {
		audio.Stop(state)
	})

	// Optional: make it look more like a toolbar button
	button.Importance = widget.LowImportance

	return button
}
