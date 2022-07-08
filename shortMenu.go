package main

import (
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/theme"
	"fyne.io/fyne/v2/widget"
)

const sizedMenuWidth = 104

// SizedMenu extends widget.Menu
type SizedMenu struct {
	*widget.Menu
}

// NewShortMenu creates a sized menu
func NewShortMenu(m *fyne.Menu) *SizedMenu {
	wid := &SizedMenu{widget.NewMenu(m)}
	wid.ExtendBaseWidget(wid)
	return wid
}

// MinSize defines the minimum size of this menu
func (s *SizedMenu) MinSize() fyne.Size {
	return fyne.NewSize(sizedMenuWidth, minHeight+cellSpace-theme.Padding()*2)
}

// CreateRenderer applies the custom layout
func (s *SizedMenu) CreateRenderer() fyne.WidgetRenderer {
	r := s.Menu.CreateRenderer()
	return &sizedMenuRenderer{r}
}

type sizedMenuRenderer struct {
	fyne.WidgetRenderer
}

// Layout sets size and position
func (r *sizedMenuRenderer) Layout(_ fyne.Size) {
	pos := fyne.NewPos(-theme.Padding(), -theme.Padding())
	size := fyne.NewSize(sizedMenuWidth+theme.Padding()*2, minHeight+cellSpace)
	for _, o := range r.Objects() {
		o.Move(pos)
		o.Resize(size)
	}
}
