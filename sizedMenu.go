package main

import (
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/theme"
	"fyne.io/fyne/v2/widget"
)

var width float32
var height float32

// sizedMenu extends widget.Menu
type sizedMenu struct {
	*widget.Menu
}

// NewSizedMenu creates a sized menu
func newSizedMenu(m *fyne.Menu, s fyne.Size) *sizedMenu {
	wid := &sizedMenu{widget.NewMenu(m)}
	width = s.Width
	height = s.Height
	wid.ExtendBaseWidget(wid)
	return wid
}

// MinSize defines the minimum size of this menu
func (s *sizedMenu) MinSize() fyne.Size {
	return fyne.NewSize(width, height)
}

// CreateRenderer applies the custom layout
func (s *sizedMenu) CreateRenderer() fyne.WidgetRenderer {
	r := s.Menu.CreateRenderer()
	return &sizedMenuRenderer{r}
}

type sizedMenuRenderer struct {
	fyne.WidgetRenderer
}

// Layout sets size and position
func (r *sizedMenuRenderer) Layout(_ fyne.Size) {
	pos := fyne.NewPos(-theme.Padding(), -theme.Padding())
	size := fyne.NewSize(width, height)
	for _, o := range r.Objects() {
		o.Move(pos)
		o.Resize(size)
	}
}
