package main

import (
	"znth/audio"
	"znth/ui/mainwindow"

	"fyne.io/fyne/v2/app"
)

func main() {
	fyneApp := app.NewWithID("com.znth.znth")

	// Initialize audio
	audio.Initialize()
	defer audio.Shutdown()

	mainwindow.Create(fyneApp)
	fyneApp.Run()
}
