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
// Combined intrument and master UI componenets
// Stem meters
// tooltips
// better meter ticks - figure out how to match starting POS with DB. Make sliders matchup better
//
// Optional:
// Change speed?
//

func main() {
	fyneApp := app.NewWithID("com.znth.znth")

	// Initialize audio
	audio.Initialize()
	defer audio.Shutdown()

	mainwindow.Create(fyneApp)
	fyneApp.Run()
}
