package toolbar

import (
	"znth/audio"
	"znth/model"
	"znth/state"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/theme"
	"fyne.io/fyne/v2/widget"
)

var playButtonWidget *widget.Button

func playButton(state *state.State) *widget.Button {
	playButtonWidget = widget.NewButtonWithIcon("", theme.MediaPlayIcon(), func() {
		audio.TogglePlay(state)
	})

	// Optional: make it look flatter like a toolbar button
	playButtonWidget.Importance = widget.LowImportance

	state.OnPlaybackChange(func(ps model.PlaybackState) {
		fyne.Do(func() {
			if ps == model.PlaybackPlaying {
				playButtonWidget.SetIcon(theme.MediaPauseIcon())
			} else {
				playButtonWidget.SetIcon(theme.MediaPlayIcon())
			}
		})
	})

	return playButtonWidget
}
