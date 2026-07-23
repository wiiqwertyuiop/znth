package model

type Stem struct {
	Id           string
	Data         []float32
	Info         WavInfo
	VolumeAdjust float32
}

type SavedStemData struct {
	VolumeAdjust float32
}
