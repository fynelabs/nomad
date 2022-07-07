package main

import (
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/theme"
	"fyne.io/fyne/v2/widget"
)

const sizedMenuWidth = 104

type sizedMenu struct {
	*widget.Menu
}

func NewShortMenu(m *fyne.Menu) *sizedMenu {
	wid := &sizedMenu{widget.NewMenu(m)}
	wid.ExtendBaseWidget(wid)
	return wid
}

func (s *sizedMenu) MinSize() fyne.Size {
	return fyne.NewSize(sizedMenuWidth, minHeight+cellSpace-theme.Padding()*2)
}

func (s *sizedMenu) CreateRenderer() fyne.WidgetRenderer {
	r := s.Menu.CreateRenderer()
	return &sizedMenuRenderer{r}
}

type sizedMenuRenderer struct {
	fyne.WidgetRenderer
}

func (r *sizedMenuRenderer) Layout(_ fyne.Size) {
	pos := fyne.NewPos(-theme.Padding(), -theme.Padding())
	size := fyne.NewSize(sizedMenuWidth+theme.Padding()*2, minHeight+cellSpace)
	for _, o := range r.Objects() {
		o.Move(pos)
		o.Resize(size)
	}
}
