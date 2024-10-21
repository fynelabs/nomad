package main

import (
	"time"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/theme"
	"fyne.io/fyne/v2/widget"
)

// timeMenu extends widget.Menu
type timeMenu struct {
	*widget.Menu

	width, height float32
}

// NewSizedMenu creates a sized menu
func newTimeMenu(m *fyne.Menu, s fyne.Size) *timeMenu {
	wid := &timeMenu{Menu: widget.NewMenu(m)}
	wid.width = s.Width
	wid.height = s.Height
	wid.ExtendBaseWidget(wid)
	return wid
}

// MinSize defines the minimum size of this menu
func (s *timeMenu) MinSize() fyne.Size {
	return fyne.NewSize(s.width, s.height)
}

// CreateRenderer applies the custom layout
func (s *timeMenu) CreateRenderer() fyne.WidgetRenderer {
	r := s.Menu.CreateRenderer()
	return &sizedMenuRenderer{WidgetRenderer: r, sized: s}
}

type sizedMenuRenderer struct {
	fyne.WidgetRenderer

	sized    *timeMenu
	scrolled bool
}

// Layout sets size and position
func (r *sizedMenuRenderer) Layout(_ fyne.Size) {
	pos := fyne.NewPos(-theme.Padding(), -theme.Padding())
	size := fyne.NewSize(r.sized.width, r.sized.height)
	for _, o := range r.Objects() {
		o.Move(pos)
		o.Resize(size)
	}
}

func (r *sizedMenuRenderer) Objects() []fyne.CanvasObject {
	objs := r.WidgetRenderer.Objects()

	scroll := objs[1].(*container.Scroll)
	if !r.scrolled {
		r.scrolled = true

		textHeight := widget.NewLabel("").MinSize().Height + theme.Padding()
		yOff := textHeight * -1.25 // 2.25 lines would be time offset, but we have the "Now" too

		hr := time.Now().Hour()
		yOff += textHeight*(float32(hr)*4) - theme.Padding()

		mn := time.Now().Minute()
		if mn >= 15 {
			yOff += textHeight
		}
		if mn >= 30 {
			yOff += textHeight
		}
		if mn >= 45 {
			yOff += textHeight
		}

		scroll.Offset = fyne.NewPos(0, yOff)
	}

	return objs
}
