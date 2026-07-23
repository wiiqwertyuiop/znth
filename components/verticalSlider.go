package components

import (
	"image/color"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/canvas"
	"fyne.io/fyne/v2/widget"
)

type VerticalSlider struct {
	widget.BaseWidget

	Min   float64
	Max   float64
	Value float64

	OnChanged func(float64)
}

type verticalSliderRenderer struct {
	slider  *VerticalSlider
	track   *canvas.Rectangle
	thumb   *canvas.Rectangle
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

	r.thumb.Resize(fyne.NewSize(20, 10))
	r.thumb.Move(fyne.NewPos(
		(size.Width-20)/2,
		y,
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

	canvas.Refresh(r.track)
	canvas.Refresh(r.thumb)
}
func (r *verticalSliderRenderer) Destroy() {}

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

func NewVerticalSlider(min, max float64, init float64) *VerticalSlider {
	s := &VerticalSlider{
		Min:   min,
		Max:   max,
		Value: init,
	}
	s.ExtendBaseWidget(s)
	return s
}

func (s *VerticalSlider) CreateRenderer() fyne.WidgetRenderer {
	track := canvas.NewRectangle(color.Gray{Y: 100})

	thumb := canvas.NewRectangle(color.White)
	thumb.StrokeColor = color.Black
	thumb.StrokeWidth = 1

	return &verticalSliderRenderer{
		slider: s,
		track:  track,
		thumb:  thumb,
		objects: []fyne.CanvasObject{
			track,
			thumb,
		},
	}
}
