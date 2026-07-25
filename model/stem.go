package model

type Channels struct {
	MasterVolume float32
	Stems        []*Stem
}

type Stem struct {
	Id           string
	Data         []int16
	Info         WavInfo
	VolumeAdjust float32
}

type SavedStemData struct {
	VolumeAdjust float32
}
