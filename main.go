package main

import (
	"os"
	"znth/audio"
	"znth/ui/mainwindow"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
)

// TODO:
// tooltips
//
// massive code clean up
//
// Combined intrument and master UI componenets
// better meter ticks - figure out how to match starting POS with DB. Make sliders matchup better
//
// don't access event directly
//	- add song event? load song event?
//
// Optional:
// Change speed?
//

func main() {
	fyneApp := app.NewWithID("com.znth.znth")

	data, err := os.ReadFile("Icon.png")
	if err != nil {
		panic(err)
	}

	fyneApp.SetIcon(fyne.NewStaticResource("Icon.png", data))

	// Initialize audio
	audio.Initialize()
	defer audio.Shutdown()

	mainwindow.Create(fyneApp)
	fyneApp.Run()
}
