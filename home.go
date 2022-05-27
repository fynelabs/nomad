//go:generate fyne bundle -o assets.go plus-circle.svg
//go:generate fyne bundle -append -o assets.go globeSpinnerSplash.gif
//go:generate fyne bundle -append -o assets.go WorkSans-BlackItalic.ttf
//go:generate fyne bundle -append -o assets.go WorkSans-Black.ttf
//go:generate fyne bundle -append -o assets.go WorkSans-Bold.ttf
//go:generate fyne bundle -append -o assets.go WorkSans-Regular.ttf

package main

import (
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/theme"
	"fyne.io/fyne/v2/widget"
)

func (n *nomad) makeAddCell() fyne.CanvasObject {
	add := widget.NewIcon(theme.NewPrimaryThemedResource(resourcePlusCircleSvg))
	search := widget.NewEntry()
	search.PlaceHolder = "ADD A PLACE"

	content := container.NewBorder(container.NewBorder(nil, nil, add, nil, search),
		nil, nil, nil)
	return container.NewPadded(content)
}

func (n *nomad) makeHome() fyne.CanvasObject {

	cells := []fyne.CanvasObject{}
	for _, c := range n.store.cities() {
		cells = append(cells, newLocation(c))
	}
	cells = append(cells, n.makeAddCell())

	layout := &nomadLayout{}
	scroll := container.NewVScroll(container.New(layout, cells...))
	scroll.SetMinSize(layout.minOuterSize())
	return scroll
}
