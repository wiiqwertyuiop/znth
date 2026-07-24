package components

import (
	"znth/audio"
	"znth/model"
	"znth/state"
)

func LoadSong(path string, state *state.State) {
	audio.KillStream(state)
	state.StatusBarTextChange("Loading files... " + path)

	stems := loadProjectFolder(path)

	project := model.Project{
		Stems:           stems,
		CurrentSongPath: path,
	}

	state.SetProject(project)

	audio.StartStream(stems, state)
	audio.Pause(state)

	state.StatusBarTextChange("Loaded succesfully! " + path)
}
