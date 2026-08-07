package mixers

import (
	"znth/model"
	"znth/state"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
)

func Create(state *state.State) *container.Scroll {

	mixers := container.NewHBox()
	mixersScroll := container.NewScroll(mixers)

	var cleanUpFunctions []func()

	state.OnProjectChange(func(p model.Project) {

		// cleanup old sliders
		for _, cleanup := range cleanUpFunctions {
			cleanup()
		}
		cleanUpFunctions = nil

		mixers.RemoveAll()

		cleanUpFunctions = drawMixers(mixers, p.Channels)

		mixers.Refresh()

		// Reset scroll
		mixersScroll.Offset = fyne.NewPos(0, 0)
		mixersScroll.Refresh()
	})

	return mixersScroll
}

func drawMixers(mixers *fyne.Container, channels model.Channels) []func() {

	var cleanUpFunctions []func()

	mixers.Add(drawMaster(channels.MasterVolume))

	for i := range channels.Stems {
		channel, cleanup := drawInstrument(channels.Stems[i])

		mixers.Add(channel)
		cleanUpFunctions = append(cleanUpFunctions, cleanup)
	}

	return cleanUpFunctions
}
