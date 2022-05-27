package main

import "fyne.io/fyne/v2"

const (
	minWidth  float32 = 300
	minHeight float32 = 200
	cellSpace float32 = 16
)

type nomadLayout struct{}

func (l *nomadLayout) Layout(objs []fyne.CanvasObject, s fyne.Size) {
	outer := cellSpace
	if fyne.CurrentDevice().IsMobile() {
		outer = 0
	}
	cols := int((s.Width - outer) / (minWidth + outer))
	colWidth := (s.Width-outer)/float32(cols) - outer
	cellSize := fyne.NewSize(colWidth, minHeight)

	offset := 0
	pos := fyne.Position{X: outer, Y: outer}
	for _, o := range objs {
		o.Resize(cellSize)
		o.Move(pos)

		offset++
		pos.X += colWidth + cellSpace
		if offset >= cols {
			offset = 0
			pos.X = outer
			pos.Y += minHeight + cellSpace
		}
	}
}

func (l *nomadLayout) MinSize(_ []fyne.CanvasObject) fyne.Size {
	if fyne.CurrentDevice().IsMobile() {
		return fyne.NewSize(minWidth, minHeight*2+cellSpace)
	}

	return fyne.NewSize(minWidth+cellSpace*2, minHeight*2+cellSpace*3)
}
