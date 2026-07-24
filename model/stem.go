package model

type Channels struct {
	MasterVolume float32
	Stems        []Stem
}

type Stem struct {
	Id           string
	Data         []float32
	Info         WavInfo
	VolumeAdjust float32
}

type SavedStemData struct {
	VolumeAdjust float32
}
