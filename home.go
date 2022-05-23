//go:generate fyne bundle -o assets.go plus-circle.svg

package main

import (
	"fmt"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/canvas"
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
	return container.NewGridWithColumns(1,
		n.makeLocationCell(), n.makeAddCell())
}

func (n *nomad) makeLocationCell() fyne.CanvasObject {
	bg := canvas.NewImageFromResource(theme.FileImageIcon())
	bg.Translucency = 0.5
	city := widget.NewRichTextFromMarkdown("# Edinburgh")
	location := widget.NewRichTextFromMarkdown("## United Kingdom · BST")
	date := widget.NewSelect([]string{}, func(string) {})
	date.PlaceHolder = "Sun 01 May"
	time := widget.NewSelect(listTimes(), func(string) {})
	time.PlaceHolder = "22:00"
	time.SetSelected("13:00")
	input := container.NewBorder(nil, nil, date, time)

	return container.NewMax(bg,
		container.NewBorder(nil,
			container.NewVBox(city, location, input), nil, nil))
}

func listTimes() (times []string) {
	for hour := 0; hour < 24; hour++ {
		times = append(times,
			fmt.Sprintf("%2d:00", hour), fmt.Sprintf("%2d:15", hour),
			fmt.Sprintf("%2d:30", hour), fmt.Sprintf("%2d:45", hour))
	}
	return times
}
