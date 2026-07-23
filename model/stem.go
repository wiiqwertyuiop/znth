package model

type Stem struct {
	Id           int
	Data         []float32
	Info         WavInfo
	VolumeAdjust float32
}

type SavedStemData struct {
	VolumeAdjust float32
}
