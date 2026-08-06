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
	slider.Disable()

	ticker := time.NewTicker(50 * time.Millisecond)

	playBackTimeText := widget.NewLabel("3:42")

	onChanged := func(value float64) {
		targetSample := int64(value * 48000 * 2)
		state.Playback.Position.Store(targetSample)
	}

	state.OnProjectChange(func(p model.Project) {
		slider.Max = float64(int64(len(p.Channels.Stems[0].Data)) / (48000 * 2))
		slider.Value = 0
		slider.Enable()
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
				slider.OnChanged = nil
				slider.SetValue(float64(totalSeconds))
				slider.OnChanged = onChanged
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
