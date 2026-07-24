package mainwindow

import (
	"znth/state"
	"znth/ui"

	"fyne.io/fyne/v2"
)

func Create(a fyne.App, state *state.State) {
	// Window setup
	w := a.NewWindow("Backing Track")
	w.CenterOnScreen()
	w.RequestFocus()
	w.Resize(fyne.NewSize(500, 500))

	// Bind state to window
	ui := ui.UI{Window: w, State: state}
	w.SetContent(createLayout(ui))
	addShortCuts(ui)

	w.Show()
}
