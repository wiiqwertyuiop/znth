package ui

import (
	"znth/state"

	"fyne.io/fyne/v2"
)

type UI struct {
	Window fyne.Window
	State  *state.State
}
