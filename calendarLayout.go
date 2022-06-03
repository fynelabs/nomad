package main

import (
	"math"

	"fyne.io/fyne/v2"
)

// Declare conformity with Layout interface
var (
	_        fyne.Layout = (*calendarLayout)(nil)
	padding  float32     = 0
	cellSize float64     = 32
)

type calendarLayout struct {
	Cols            int
	vertical, adapt bool
}

func newCalendarLayout(s float64) fyne.Layout {
	cellSize = s
	return &calendarLayout{Cols: 7}
}

func (g *calendarLayout) horizontal() bool {
	if g.adapt {
		return fyne.IsHorizontal(fyne.CurrentDevice().Orientation())
	}

	return !g.vertical
}

func (g *calendarLayout) countRows(objects []fyne.CanvasObject) int {
	count := 0
	for _, child := range objects {
		if child.Visible() {
			count++
		}
	}

	return int(math.Ceil(float64(count) / float64(g.Cols)))
}

// Get the leading (top or left) edge of a grid cell.
// size is the ideal cell size and the offset is which col or row its on.
func getLeading(size float64, offset int) float32 {
	ret := (size + float64(padding)) * float64(offset)

	return float32(math.Round(ret))
}

// Get the trailing (bottom or right) edge of a grid cell.
// size is the ideal cell size and the offset is which col or row its on.
func getTrailing(size float64, offset int) float32 {
	return getLeading(size, offset+1) - padding
}

// Layout is called to pack all child objects into a specified size.
// For a GridLayout this will pack objects into a table format with the number
// of columns specified in our constructor.
func (g *calendarLayout) Layout(objects []fyne.CanvasObject, size fyne.Size) {
	row, col := 0, 0
	i := 0
	for _, child := range objects {
		if !child.Visible() {
			continue
		}

		x1 := getLeading(cellSize, col)
		y1 := getLeading(cellSize, row)
		x2 := getTrailing(cellSize, col)
		y2 := getTrailing(cellSize, row)

		child.Move(fyne.NewPos(x1, y1))
		child.Resize(fyne.NewSize(x2-x1, y2-y1))

		if g.horizontal() {
			if (i+1)%g.Cols == 0 {
				row++
				col = 0
			} else {
				col++
			}
		} else {
			if (i+1)%g.Cols == 0 {
				col++
				row = 0
			} else {
				row++
			}
		}
		i++
	}
}

func (g *calendarLayout) MinSize(objects []fyne.CanvasObject) fyne.Size {
	rows := g.countRows(objects)
	return fyne.NewSize(float32(cellSize+float64(padding))*7, float32(cellSize+float64(padding))*float32(rows))
}
