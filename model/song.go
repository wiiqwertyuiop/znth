package model

type SongDetails struct {
	Name     string
	Location string
}

type Setlist struct {
	Data []SongDetails
}

type SongSave struct {
	Stems map[string]SavedStemData
}
