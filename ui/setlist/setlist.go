package setlist

import (
	"image/color"
	"znth/components"
	"znth/model"
	"znth/state"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/canvas"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/widget"
)

func Create(state *state.State) fyne.CanvasObject {
	setlist := widget.NewList(
		func() int {
			if len(state.Project.SongNames) == 0 {
				return 1
			}

			return len(state.Project.SongNames)
		},
		func() fyne.CanvasObject {
			label := widget.NewLabel("")
			label.Wrapping = fyne.TextWrap(fyne.TextTruncateClip)
			label.TextStyle.Bold = true
			return label
		},
		func(i widget.ListItemID, o fyne.CanvasObject) {
			label := o.(*widget.Label)

			if len(state.Project.SongNames) == 0 {
				label.SetText("[Empty setlist...]")
				return
			}

			label.SetText(state.Project.SongNames[i].Name)
		},
	)

	setlist.OnSelected = func(id widget.ListItemID) {
		if len(state.Project.SongNames) == 0 {
			return
		}

		path := state.Project.SongNames[id].Location
		components.LoadSong(path, state)
	}

	state.OnProjectChange(func(p *model.Project) {
		setlist.Refresh()
	})

	title := canvas.NewText("Setlist", color.Black)
	title.TextSize = 16
	title.TextStyle.Bold = true
	title.Alignment = fyne.TextAlignCenter

	bg := canvas.NewRectangle(color.NRGBA{
		R: 250,
		G: 250,
		B: 250,
		A: 255,
	})

	separator := canvas.NewRectangle(color.Black)
	separator.SetMinSize(fyne.NewSize(0, 1))

	headerContent := container.NewCenter(title)

	header := container.NewBorder(
		nil,
		separator,
		nil,
		nil,
		container.NewStack(
			bg,
			headerContent,
		),
	)

	return container.NewBorder(
		header, // top
		nil,    // bottom
		nil,    // left
		nil,    // right
		setlist,
	)
}
