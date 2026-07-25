package components

import (
	"slices"
	"znth/audio"
	"znth/model"
	"znth/state"

	"fyne.io/fyne/v2"
)

func LoadSong(path string, state *state.State) {
	audio.KillStream(state)
	state.StatusBarTextChange("Loading files... " + path)

	channels := loadProjectFolder(path)

	project := model.Project{
		Channels:        channels,
		CurrentSongPath: path,
		SongNames:       state.Project.SongNames,
	}

	state.SetProject(project)

	audio.StartStream(channels.Stems, state)
	audio.Pause(state)

	state.StatusBarTextChange("Loaded succesfully! " + path)
}

func AddSong(reader fyne.ListableURI, state *state.State) {
	path := reader.Path()

	audio.KillStream(state)
	state.StatusBarTextChange("Loading files... " + path)

	channels := loadProjectFolder(path)

	song := model.SongDetails{
		Name:     reader.Name(),
		Location: path,
	}

	// TODO: dont modify state like this
	state.Project.SongNames = slices.Insert(state.Project.SongNames, 0, song)

	project := model.Project{
		Channels:        channels,
		CurrentSongPath: path,
		SongNames:       state.Project.SongNames,
	}

	state.SetProject(project)

	audio.StartStream(channels.Stems, state)
	audio.Pause(state)

	state.StatusBarTextChange("Loaded succesfully! " + path)
}
