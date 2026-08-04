package mainwindow

import (
	"image/color"
	"znth/state"
	"znth/ui"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/theme"
)

func Create(a fyne.App) {
	// Window setup
	w := a.NewWindow("Backing Track")
	a.Settings().SetTheme(DarkTheme{})

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

type DarkTheme struct{}

func (DarkTheme) Color(name fyne.ThemeColorName, variant fyne.ThemeVariant) color.Color {
	return theme.DefaultTheme().Color(name, theme.VariantDark)
}

func (DarkTheme) Font(style fyne.TextStyle) fyne.Resource {
	return theme.DefaultTheme().Font(style)
}

func (DarkTheme) Icon(name fyne.ThemeIconName) fyne.Resource {
	return theme.DefaultTheme().Icon(name)
}

func (DarkTheme) Size(name fyne.ThemeSizeName) float32 {
	return theme.DefaultTheme().Size(name)
}
