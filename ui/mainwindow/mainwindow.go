package mainwindow

import (
	"znth/state"
	"znth/ui"

	"fyne.io/fyne/v2"
)

func Create(a fyne.App) {
	// Window setup
	w := a.NewWindow("Backing Track")
	w.CenterOnScreen()
	w.RequestFocus()
	w.Resize(fyne.NewSize(800, 500))

	// Create window state
	state := state.New()

	// Bind state to window
	ui := ui.UI{Window: w, State: state}

	// Create and set content
	w.SetContent(createLayout(ui))
	addShortCuts(ui)

	// Show window
	w.Show()
}
