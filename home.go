//go:generate fyne bundle -o assets.go plus-circle.svg

package main

import (
	"time"

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
	zone, _ := time.LoadLocation("Europe/London")
	cell := newLocation("Edinburgh", "United Kingdom", zone)

	c := container.New(&nomadLayout{},
		cell, n.makeAddCell())

	c.Hide()

	return c
}

func (n *nomad) showHome(obj fyne.CanvasObject, splashFadeTime int) {
	time.Sleep(time.Second * time.Duration(splashFadeTime))
	obj.Show()
	n.main.Content().Refresh()
}
