package main

import (
	"znth/audio"
	"znth/ui/mainwindow"

	"fyne.io/fyne/v2/app"
)

// TODO:
// don't access event directly
//	- add song event? load song event?
//
// Stem meters
// Seeking
// Change speed?
//
// Optional:
// tooltips
// better meter ticks
//

func main() {
	fyneApp := app.NewWithID("com.znth.znth")

	// Initialize audio
	audio.Initialize()
	defer audio.Shutdown()

	mainwindow.Create(fyneApp)
	fyneApp.Run()
}
