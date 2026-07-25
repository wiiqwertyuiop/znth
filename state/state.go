package state

import "znth/model"

type State struct {
	Playback model.Playback
	Project  model.Project

	projectListeners  []func(model.Project)
	playbackListeners []func(model.PlaybackState)
	statusBarListener func(string)
}

func New() *State {
	return &State{}
}

func (s *State) StatusBarTextChange(str string) {
	if s != nil {
		s.statusBarListener(str)
	}
}

func (s *State) OnStatusBarTextChange(listener func(str string)) {
	s.statusBarListener = listener
}

func (s *State) PlaybackChange(playbackState model.PlaybackState) {
	s.Playback.State = playbackState
	for _, listener := range s.playbackListeners {
		listener(playbackState)
	}
}

func (s *State) OnPlaybackChange(listener func(model.PlaybackState)) {
	s.playbackListeners = append(s.playbackListeners, listener)
}

func (s *State) SetProject(project model.Project) {
	s.Project = project
	for _, listener := range s.projectListeners {
		listener(s.Project)
	}
}

func (s *State) OnProjectChange(listener func(model.Project)) {
	s.projectListeners = append(s.projectListeners, listener)
}
