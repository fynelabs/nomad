package main

import "fyne.io/fyne/v2"

const (
	minWidth  float32 = 300
	minHeight float32 = 200
)

type nomadLayout struct{}

func (l *nomadLayout) Layout(objs []fyne.CanvasObject, s fyne.Size) {
	cols := int(s.Width / minWidth)
	colWidth := s.Width / float32(cols)
	cellSize := fyne.NewSize(colWidth, minHeight)

	offset := 0
	pos := fyne.Position{}
	for _, o := range objs {
		o.Resize(cellSize)
		o.Move(pos)

		offset++
		pos.X += colWidth
		if offset >= cols {
			offset = 0
			pos.X = 0
			pos.Y += minHeight
		}
	}
}

func (l *nomadLayout) MinSize(_ []fyne.CanvasObject) fyne.Size {
	return fyne.NewSize(minWidth, minHeight*2)
}
