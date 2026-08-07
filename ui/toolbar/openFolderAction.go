package toolbar

import (
	"znth/components"
	"znth/ui"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/dialog"
	"fyne.io/fyne/v2/theme"
	"fyne.io/fyne/v2/widget"
)

func openFolderButton(ui ui.UI) *widget.Button {
	button := widget.NewButtonWithIcon("", theme.FolderOpenIcon(), func() {
		dialog.NewFolderOpen(func(reader fyne.ListableURI, err error) {

			if err != nil {
				ui.State.StatusBarTextChange("Folder selection error: " + err.Error())
				return
			}

			if reader == nil {
				ui.State.StatusBarTextChange("Folder selection cancelled")
				return
			}

			components.AddSong(reader, ui.State)
		}, ui.Window).Show()
	})

	// Optional: make it look more like a toolbar button
	button.Importance = widget.LowImportance

	return button
}
