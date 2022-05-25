//go:generate fyne bundle -o assets.go plus-circle.svg
//go:generate fyne bundle -append -o assets.go globeSpinnerSplash.gif
//go:generate fyne bundle -append -o assets.go WorkSans-BlackItalic.ttf

package main

import (
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/theme"
	"fyne.io/fyne/v2/widget"
)

func (n *nomad) makeAddCell() fyne.CanvasObject {
	add := widget.NewIcon(theme.NewThemedResource(resourcePlusCircleSvg))
	search := widget.NewEntry()
	search.PlaceHolder = "Add a place"
	content := container.NewBorder(container.NewBorder(nil, nil, add, nil, search),
		nil, nil, nil)
	return container.NewPadded(content)
}

func (n *nomad) makeHome() fyne.CanvasObject {
	cells := []fyne.CanvasObject{}
	for _, c := range n.store.cities() {
		cells = append(cells, newLocation(c, n))
	}
	cells = append(cells, n.makeAddCell())

	return container.New(&nomadLayout{}, cells...)
}
