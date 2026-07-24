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
		drawMixers(mixers, p.Stems)
		mixers.Refresh()
	})
	return container.NewScroll(mixers)
}

func drawMixers(mixers *fyne.Container, stems []model.Stem) {

	for _, stem := range stems {
		var channel *fyne.Container
		if stem.Id == "Master" {
			channel = drawMaster(stem)
		} else {
			channel = drawInstrument(stem)
		}
		mixers.Add(channel)
	}
}
