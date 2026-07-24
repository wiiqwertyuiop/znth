package toolbar

import (
	"znth/audio"
	"znth/model"
	"znth/state"

	"fyne.io/fyne/v2/theme"
	"fyne.io/fyne/v2/widget"
)

var playAction *widget.ToolbarAction

func playButton(state *state.State) *widget.ToolbarAction {
	playAction = widget.NewToolbarAction(theme.MediaPlayIcon(), func() { audio.TogglePlay(state) })

	state.OnPlaybackChange(func(ps model.PlaybackState) {
		if ps == model.PlaybackPlaying {
			playAction.SetIcon(theme.MediaPauseIcon())
		} else {
			playAction.SetIcon(theme.MediaPlayIcon())
		}
	})
	return playAction
}
