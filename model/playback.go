package model

import "sync/atomic"

type PlaybackState string

const (
	PlaybackStopped PlaybackState = "Stopped"
	PlaybackPlaying PlaybackState = "Playing"
	PlaybackPaused  PlaybackState = "Paused"
)

type Playback struct {
	State    PlaybackState
	Position atomic.Int64
}
