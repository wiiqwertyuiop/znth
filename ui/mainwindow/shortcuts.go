package mainwindow

import (
	"znth/audio"
	"znth/components"
	"znth/ui"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/driver/desktop"
)

var fullscreen = false

func addShortCuts(ui ui.UI) {
	ctrlS_Shortcut := &desktop.CustomShortcut{
		KeyName:  fyne.KeyS,
		Modifier: fyne.KeyModifierControl,
	}

	ui.Window.Canvas().AddShortcut(ctrlS_Shortcut, func(shortcut fyne.Shortcut) {
		if ui.State.Project.CurrentSongPath != "" {
			components.SaveStemData(ui.State.Project.CurrentSongPath, ui.State.Project.Channels.Stems)
			ui.State.StatusBarTextChange("Saved stem levels!")
		}
	})

	ui.Window.Canvas().SetOnTypedKey(func(ev *fyne.KeyEvent) {
		switch ev.Name {
		case fyne.KeySpace:
			audio.TogglePlay(ui.State)
		case fyne.KeyF11:
			fullscreen = !fullscreen
			ui.Window.SetFullScreen(fullscreen)
		}

	})
}
