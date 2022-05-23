package main

import (
	"time"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/canvas"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/theme"
	"fyne.io/fyne/v2/widget"
)

type location struct {
	name, country string
	zone          *time.Location

	date, time *widget.Select
}

func newLocation(name, country string, z *time.Location) *location {
	return &location{name: name, country: country, zone: z}
}

func (l *location) makeUI() fyne.CanvasObject {
	bg := canvas.NewImageFromResource(theme.FileImageIcon())
	bg.Translucency = 0.5
	city := widget.NewRichTextFromMarkdown("# " + l.name)
	location := widget.NewRichTextFromMarkdown("## " + l.country + " · " + l.zone.String())
	l.date = widget.NewSelect([]string{}, func(string) {})
	l.date.PlaceHolder = "Sun 01 May"
	l.time = widget.NewSelect(listTimes(), func(string) {})
	l.time.PlaceHolder = "22:00" // longest
	l.time.SetSelected("13:00")
	input := container.NewBorder(nil, nil, l.date, l.time)

	return container.NewMax(bg,
		container.NewBorder(nil,
			container.NewVBox(city, location, input), nil, nil))
}
