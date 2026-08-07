package state

import (
	"time"
	"znth/model"
)

type State struct {
	Playback model.Playback
	Project  model.Project

	projectListeners  []func(*model.Project)
	playbackListeners []func(model.PlaybackState)
	statusBarListener func(string)

	mainWindowRenderLoopListeners []func()
}

func New() *State {
	newState := &State{}
	ticker := time.NewTicker(50 * time.Millisecond)
	go func() {
		for range ticker.C {
			for _, listener := range newState.mainWindowRenderLoopListeners {
				listener()
			}
			newState.Project.RenderProjectElements()
		}
	}()
	return newState
}

func (s *State) OnMainWindowRenderLoop(f func()) {
	s.mainWindowRenderLoopListeners = append(s.mainWindowRenderLoopListeners, f)
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
	// Clean up
	s.Project.ProjectCleanup()

	// Assign new project
	s.Project = project
	for _, listener := range s.projectListeners {
		listener(&s.Project)
	}
}

func (s *State) OnProjectChange(listener func(*model.Project)) {
	s.projectListeners = append(s.projectListeners, listener)
}
