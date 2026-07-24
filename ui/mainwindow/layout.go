package mainwindow

import (
	"znth/ui"
	"znth/ui/mixers"
	"znth/ui/setlist"
	"znth/ui/statusbar"
	"znth/ui/toolbar"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/canvas"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/theme"
)

func createLayout(ui ui.UI) *fyne.Container {

	// Borders
	toolbarBorder := canvas.NewLine(theme.Color(theme.ColorNameSeparator))
	statusBorder := canvas.NewLine(theme.Color(theme.ColorNameSeparator))

	toolbarBorder.StrokeWidth = 1
	statusBorder.StrokeWidth = 1

	// Toolbar
	toolbar := container.NewBorder(
		nil,
		toolbarBorder,
		nil,
		nil,
		toolbar.Create(ui),
	)

	// Split view for setlist and mixers
	mainContent := container.NewHSplit(
		setlist.Create(ui.State),
		mixers.Create(ui.State),
	)

	// Start with left panel at 10% width
	mainContent.SetOffset(0.10)

	// Bottom status bar
	statusBar := container.NewBorder(
		statusBorder,
		nil,
		nil,
		nil,
		statusbar.Create(ui.State),
	)

	// Full window layout
	return container.NewBorder(
		toolbar,
		statusBar,
		nil,
		nil,
		mainContent,
	)
}
