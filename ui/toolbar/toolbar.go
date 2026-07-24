package toolbar

import (
	"znth/components"
	"znth/ui"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/dialog"
	"fyne.io/fyne/v2/theme"
	"fyne.io/fyne/v2/widget"
)

func Create(ui ui.UI) *widget.Toolbar {

	return widget.NewToolbar(
		openFolderButton(ui),
		widget.NewToolbarAction(theme.FileApplicationIcon(), func() {
			dialog.NewFileOpen(func(reader fyne.URIReadCloser, err error) {
				components.OpenSetlist(reader, ui.State)
			}, ui.Window).Show()
		}),
		widget.NewToolbarAction(theme.DocumentSaveIcon(), func() {
			dialog.NewFileSave(func(writer fyne.URIWriteCloser, err error) {
				if err != nil {
					ui.State.StatusBarTextChange("Save error: " + err.Error())
					return
				}
				e := components.SaveSetlist(writer, ui.State) // TODO
				if e != nil {
					ui.State.StatusBarTextChange(e.Error())
				}
			}, ui.Window).Show()
		}),
		widget.NewToolbarSeparator(),
		playButton(ui.State),
		stopButton(ui.State),
	)
}
