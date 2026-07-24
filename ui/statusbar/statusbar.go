package statusbar

import (
	"znth/state"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/widget"
)

func Create(state *state.State) *widget.Label {
	info := widget.NewLabel("Open a setlist or stem folder to get started.")
	info.TextStyle.Bold = true

	state.OnStatusBarTextChange(func(str string) {
		fyne.Do(func() {
			info.SetText(str)
		})
	})

	return info
}
