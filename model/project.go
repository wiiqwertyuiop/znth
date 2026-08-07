package model

type Project struct {
	Channels        Channels
	CurrentSongPath string
	SongNames       []SongDetails

	projectRenderLoopListeners []func()
}

func (p *Project) AddToProjectRenderLoop(f func()) {
	p.projectRenderLoopListeners = append(p.projectRenderLoopListeners, f)
}

func (p *Project) RenderProjectElements() {
	for _, listener := range p.projectRenderLoopListeners {
		listener()
	}
}

func (p *Project) ProjectCleanup() {
	for _, stem := range p.Channels.Stems {
		if stem != nil {
			stem.Data = nil
		}
	}

	p.Channels.Stems = nil
	p.projectRenderLoopListeners = nil
}
