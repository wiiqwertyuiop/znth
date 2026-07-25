package mixers

import (
	"znth/model"
	"znth/state"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
)

func Create(state *state.State) *container.Scroll {

	mixers := container.NewHBox()

	state.OnProjectChange(func(p model.Project) {
		mixers.RemoveAll()
		drawMixers(mixers, p.Channels)
		mixers.Refresh()
	})

	return container.NewScroll(mixers)
}

func drawMixers(mixers *fyne.Container, channels model.Channels) {
	mixers.Add(drawMaster(channels.MasterVolume))
	for i := range channels.Stems {
		mixers.Add(drawInstrument(&channels.Stems[i]))
	}
}
