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
	fyne.Do(func() {
		state.StatusBarTextChange("Loading files... " + path)
		state.IsLoading(true)
	})

	go func() {
		channels := loadProjectFolder(path)

		project := model.Project{
			Channels:        channels,
			CurrentSongPath: path,
			SongNames:       state.Project.SongNames,
		}

		/* var m runtime.MemStats
		runtime.ReadMemStats(&m)

		fmt.Printf("Alloc: %d MB\n", m.Alloc/1024/1024)
		fmt.Printf("HeapAlloc: %d MB\n", m.HeapAlloc/1024/1024)
		fmt.Printf("HeapInuse: %d MB\n", m.HeapInuse/1024/1024)
		fmt.Printf("HeapSys: %d MB\n", m.HeapSys/1024/1024) */
		fyne.DoAndWait(func() {
			state.SetProject(project)
			state.IsLoading(false)
			state.StatusBarTextChange("Loaded successfully! " + path)
		})

		audio.StartStream(channels.Stems, state)
		audio.Pause(state)
	}()

}

func AddSong(reader fyne.ListableURI, state *state.State) {
	path := reader.Path()

	audio.KillStream(state)
	fyne.Do(func() {
		state.StatusBarTextChange("Loading files... " + path)
		state.IsLoading(true)
	})

	go func() {
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

		fyne.DoAndWait(func() {
			state.SetProject(project)
			state.IsLoading(false)
			state.StatusBarTextChange("Loaded successfully! " + path)
		})

		audio.StartStream(channels.Stems, state)
		audio.Pause(state)
	}()
}
