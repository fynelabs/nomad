package main

import "fyne.io/fyne/v2"

type sizedMenu struct {
}

func (d *sizedMenu) MinSize(objects []fyne.CanvasObject) fyne.Size {
	return fyne.NewSize(100, 20)
}

func (s *sizedMenu) Layout(objects []fyne.CanvasObject, containerSize fyne.Size) {
	for _, o := range objects {
		o.Move(fyne.NewPos(-200, -100))  // can change position
		o.Resize(fyne.NewSize(100, 200)) // not working
	}
}
