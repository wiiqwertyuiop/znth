package main

import (
	"znth/audio"
	"znth/state"
	"znth/ui/mainwindow"

	"fyne.io/fyne/v2/app"
)

func main() {
	fyneApp := app.NewWithID("com.znth.znth")

	// Initialize audio
	audio.Initialize()
	defer audio.Shutdown()

	// Create new app state
	state := state.Initialize()

	mainwindow.Create(fyneApp, state)
	fyneApp.Run()
}
