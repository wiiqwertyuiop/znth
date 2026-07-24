package toolbar

import (
	"znth/audio"
	"znth/state"

	"fyne.io/fyne/v2/theme"
	"fyne.io/fyne/v2/widget"
)

func stopButton(state *state.State) *widget.ToolbarAction {
	return widget.NewToolbarAction(theme.MediaStopIcon(), func() { audio.Stop(state) })
}
