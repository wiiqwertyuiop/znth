package toolbar

import (
	"znth/audio"
	"znth/components"
	"znth/ui"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/dialog"
	"fyne.io/fyne/v2/theme"
	"fyne.io/fyne/v2/widget"
)

// Can these be moved?
var playAction *widget.ToolbarAction
var setlist *widget.List
var mixers = container.NewHBox()

var info = widget.NewLabel("Nothing loaded...") // global?

var currentSongPath string = ""

func CreateToolbar(ui ui.UI) *widget.Toolbar {
	playAction = widget.NewToolbarAction(theme.MediaPlayIcon(), play)

	return widget.NewToolbar(
		widget.NewToolbarAction(theme.FolderOpenIcon(), func() {
			dialog.NewFolderOpen(func(reader fyne.ListableURI, err error) {

				if err != nil {
					info.SetText("Folder selection error: " + err.Error())
					return
				}

				if reader == nil {
					info.SetText("Folder selection cancelled")
					return
				}
				loadSong(reader.Path())
				components.AddToSetlist(reader)
				setlist.Refresh()
			}, ui.Window).Show()
		}),
		widget.NewToolbarAction(theme.FileApplicationIcon(), func() {
			dialog.NewFileOpen(func(reader fyne.URIReadCloser, err error) {
				components.OpenSetlist(reader)
				setlist.Refresh()
			}, ui.Window).Show()
		}),
		widget.NewToolbarAction(theme.DocumentSaveIcon(), func() {
			dialog.NewFileSave(func(writer fyne.URIWriteCloser, err error) {
				if err != nil {
					info.SetText("Save error: " + err.Error())
					return
				}
				e := components.SaveSetlist(writer) // TODO
				if e != nil {
					info.SetText(e.Error())
				}
			}, ui.Window).Show()
		}),
		widget.NewToolbarSeparator(),
		playAction,
		widget.NewToolbarAction(theme.MediaStopIcon(), func() {
			if audio.IsStreamActive() {
				components.StopAction()
				playAction.SetIcon(theme.MediaPlayIcon())
			}
		}),
	)
}

func loadSong(path string) {
	audio.KillStream()
	mixers.RemoveAll()
	playAction.SetIcon(theme.MediaPlayIcon())

	info.SetText("Loading files... " + path)

	components.LoadProjectFolder(path, mixers) // Todo
	currentSongPath = path

	info.SetText("Loaded succesfully! " + path)
}

func play() {
	if audio.IsStreamActive() {
		components.PlayAction()
		if audio.IsPlaying() {
			playAction.SetIcon(theme.MediaPauseIcon())
		} else {
			playAction.SetIcon(theme.MediaPlayIcon())
		}
	}
}
