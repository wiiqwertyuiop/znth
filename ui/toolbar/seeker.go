package toolbar

import (
	"fmt"
	"time"
	"znth/model"
	"znth/state"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/widget"
)

func createSeeker(state *state.State) fyne.CanvasObject {
	slider := widget.NewSlider(0, 100)

	ticker := time.NewTicker(50 * time.Millisecond)

	playBackTimeText := widget.NewLabel("3:42")

	slider.OnChangeEnded = func(value float64) {
		targetSample := int64(value * 48000 * 2)

		state.Playback.Position.Store(targetSample)
	}

	state.OnPlaybackChange(func(ps model.PlaybackState) {
		if ps != model.PlaybackPlaying {
			return
		}
		slider.Max = float64(state.Playback.Length)
		slider.Value = 0
		slider.Refresh()
	})

	go func() {
		for range ticker.C {
			position := state.Playback.Position.Load()

			// Set label text
			totalSeconds := int64(float64(position) / float64(48000*2))

			minutes := totalSeconds / 60
			seconds := totalSeconds % 60

			fyne.Do(func() {
				slider.SetValue(float64(totalSeconds))
				playBackTimeText.SetText(fmt.Sprintf("%d:%02d", minutes, seconds))
			})
		}
	}()

	return container.NewBorder(
		nil,
		nil,
		container.NewCenter(widget.NewLabel("0:00")),
		container.NewCenter(playBackTimeText),
		slider,
	)
}
