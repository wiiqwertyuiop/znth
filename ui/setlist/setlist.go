package setlist

import (
	"znth/components"
	"znth/state"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/widget"
)

func Create(state *state.State) *widget.List {
	// Setlist
	setlist := widget.NewList(
		func() int {
			return len(state.Project.SongNames)
		},
		func() fyne.CanvasObject {
			label := widget.NewLabel("")
			label.Wrapping = fyne.TextWrap(fyne.TextTruncateClip)
			label.TextStyle.Bold = true
			return label
		},
		func(i widget.ListItemID, o fyne.CanvasObject) {
			o.(*widget.Label).SetText(state.Project.SongNames[i].Name)
		},
	)
	setlist.OnSelected = func(id widget.ListItemID) {
		path := state.Project.SongNames[id].Location
		components.LoadSong(path, state)
	}

	return setlist
}
