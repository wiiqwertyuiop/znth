package main

import (
	_ "embed"

	"fmt"
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

//go:embed Icon.png
var iconData []byte

func main() {
	defer func() {
		if r := recover(); r != nil {
			fmt.Printf("PANIC: %v\n", r)
			fmt.Println("Press Enter to exit...")
			fmt.Scanln()
		}
	}()

	fyneApp := app.NewWithID("com.znth.znth")

	fyneApp.SetIcon(
		fyne.NewStaticResource("Icon.png", iconData),
	)

	// Initialize audio
	audio.Initialize()
	defer audio.Shutdown()

	mainwindow.Create(fyneApp)
	fyneApp.Run()
}
