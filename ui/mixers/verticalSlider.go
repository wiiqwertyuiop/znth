package mixers

import (
	"image/color"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/canvas"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/widget"
)

type VerticalSlider struct {
	widget.BaseWidget

	Min   float64
	Max   float64
	Value float64

	Peak           float32
	PeakHoldFrames int

	OnChanged func(float64)
}

type verticalSliderRenderer struct {
	slider  *VerticalSlider
	track   *canvas.Rectangle
	peak    *canvas.Rectangle
	thumb   *fyne.Container
	objects []fyne.CanvasObject
}

func (r *verticalSliderRenderer) Layout(size fyne.Size) {

	trackWidth := float32(6)

	r.track.Resize(fyne.NewSize(trackWidth, size.Height))
	r.track.Move(fyne.NewPos(
		(size.Width-trackWidth)/2,
		0,
	))

	// Position thumb based on slider value
	pct := (r.slider.Value - r.slider.Min) /
		(r.slider.Max - r.slider.Min)

	y := size.Height - float32(pct)*size.Height

	r.thumb.Resize(fyne.NewSize(25, 50))
	r.thumb.Move(fyne.NewPos(
		(size.Width-25)/2,
		y-20,
	))
}

func (r *verticalSliderRenderer) MinSize() fyne.Size {
	return fyne.NewSize(30, 150)
}

func (r *verticalSliderRenderer) Objects() []fyne.CanvasObject {
	return r.objects
}

func (s *VerticalSlider) DragEnd() {}

func (r *verticalSliderRenderer) Refresh() {
	r.Layout(r.slider.Size())

	trackWidth := float32(6)

	peakHeight := r.slider.Peak * r.slider.Size().Height

	r.peak.Resize(fyne.NewSize(
		trackWidth,
		peakHeight,
	))

	r.peak.Move(fyne.NewPos(
		(r.slider.Size().Width-trackWidth-10)/2,
		r.slider.Size().Height-peakHeight,
	))

	canvas.Refresh(r.track)
	canvas.Refresh(r.peak)
	canvas.Refresh(r.thumb)
}

func (r *verticalSliderRenderer) Destroy() {}

func (s *VerticalSlider) UpdatePeak(input float32) {

	if s == nil {
		return
	}

	if input < 0 {
		s.Peak = 0
		s.Refresh()
		return
	}

	const threshold = float32(0.1)

	if input > s.Peak {
		// Fast attack
		s.Peak = input
		s.PeakHoldFrames = 2
		return
	}

	// Non-linear decay
	decay := float32(0.97)

	// Hold if the signal is still close
	if s.Peak-input < threshold {
		decay = float32(0.99)
	}

	if s.Peak < 0.2 {
		decay = 0.85
	} else if s.Peak < 0.5 {
		decay = 0.90
	}

	s.Peak *= decay

	if s.Peak < 0.001 {
		s.Peak = 0
	}

	s.Refresh()
}

func (s *VerticalSlider) Dragged(e *fyne.DragEvent) {
	size := s.Size()

	// Current mouse Y relative to the widget
	y := e.Position.Y

	// Clamp to the slider
	if y < 0 {
		y = 0
	}
	if y > size.Height {
		y = size.Height
	}

	// Convert to 0..1 (top = 1, bottom = 0)
	pct := 1.0 - float64(y/size.Height)

	// Convert to slider range
	s.Value = s.Min + pct*(s.Max-s.Min)

	if s.OnChanged != nil {
		s.OnChanged(s.Value)
	}

	s.Refresh()
}

func newVerticalSlider(min, max float64, init float64) *VerticalSlider {
	s := &VerticalSlider{
		Min:   min,
		Max:   max,
		Value: init,
	}
	s.ExtendBaseWidget(s)
	return s
}

func (s *VerticalSlider) CreateRenderer() fyne.WidgetRenderer {
	track := canvas.NewRectangle(color.RGBA{R: 255, G: 255, B: 255, A: 80})
	peak := canvas.NewRectangle(color.RGBA{R: 0, G: 255, B: 0, A: 80})

	rect := canvas.NewRectangle(color.RGBA{R: 255, G: 255, B: 255, A: 200})
	rect.CornerRadius = 10
	rect.Resize(fyne.NewSize(25, 50))

	line := canvas.NewLine(color.Black)
	line.StrokeWidth = 3
	line.Position1 = fyne.NewPos(0, 25)
	line.Position2 = fyne.NewPos(25, 25)

	thumb := container.NewWithoutLayout(
		rect,
		line,
	)

	return &verticalSliderRenderer{
		slider: s,
		track:  track,
		peak:   peak,
		thumb:  thumb,
		objects: []fyne.CanvasObject{
			track,
			peak,
			thumb,
		},
	}
}
