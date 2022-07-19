package main

import (
	"math"

	"fyne.io/fyne/v2"
)

const (
	minWidth  float32 = 300
	minHeight float32 = 200
	cellSpace float32 = 16
)

type nomadLayout struct {
	cols int
}

func (l *nomadLayout) Layout(objs []fyne.CanvasObject, s fyne.Size) {
	outer := cellSpace
	if fyne.CurrentDevice().IsMobile() {
		outer = 0
	}
	l.cols = int((s.Width - outer*2 + cellSpace) / (minWidth + cellSpace))
	colWidth := (s.Width-outer*2+cellSpace)/float32(l.cols) - cellSpace
	cellSize := fyne.NewSize(colWidth, minHeight)

	offset := 0
	pos := fyne.Position{X: outer, Y: outer}
	for _, o := range objs {
		o.Resize(cellSize)
		o.Move(pos)

		offset++
		pos.X += colWidth + cellSpace
		if offset >= l.cols {
			offset = 0
			pos.X = outer
			pos.Y += minHeight + cellSpace
		}
	}
}

func (l *nomadLayout) MinSize(cells []fyne.CanvasObject) fyne.Size {
	// we calculate how much is required to scroll to fit in all of the cells
	cols := l.cols
	if cols < 1 { // possibly not layed out yet
		cols = 1
	}
	rows := int(math.Ceil(float64(len(cells)) / float64(l.cols)))
	height := float32((minHeight+cellSpace*4)*float32(rows) - cellSpace)
	if fyne.CurrentDevice().IsMobile() {
		return fyne.NewSize(minWidth, height)
	}

	return fyne.NewSize(minWidth, height+cellSpace*2)
}

func (l *nomadLayout) minOuterSize() fyne.Size {
	if fyne.CurrentDevice().IsMobile() {
		return fyne.NewSize(minWidth, minHeight)
	}

	return fyne.NewSize(minWidth+cellSpace*2, minHeight+cellSpace*2)
}
