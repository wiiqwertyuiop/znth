package components

import (
	"znth/audio"
	"znth/model"
	"znth/state"
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
