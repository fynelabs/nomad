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
	localTime     time.Time

	date *widget.Select
	time *widget.SelectEntry
}

func newLocation(name, country string, z *time.Location) *location {
	t := time.Now().In(z)
	return &location{name: name, country: country, localTime: t}
}

func (l *location) makeUI() fyne.CanvasObject {
	bg := canvas.NewImageFromResource(theme.FileImageIcon())
	bg.Translucency = 0.5
	city := widget.NewRichTextFromMarkdown("# " + l.name)
	location := widget.NewRichTextFromMarkdown("## " + l.country + " · " + l.localTime.Format("MST"))
	l.date = widget.NewSelect([]string{}, func(string) {})
	l.date.PlaceHolder = l.localTime.Format("Mon 02 Jan")
	l.time = widget.NewSelectEntry(listTimes())
	l.time.PlaceHolder = "22:00" // longest
	l.time.Wrapping = fyne.TextWrapOff
	l.time.SetText(l.localTime.Format("15:04"))
	input := container.NewBorder(nil, nil, l.date, l.time)

	return container.NewMax(bg,
		container.NewBorder(nil,
			container.NewVBox(city, location, input), nil, nil))
}
