package main

import (
	"fmt"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/canvas"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/layout"
	"fyne.io/fyne/v2/theme"
	"fyne.io/fyne/v2/widget"
)

type location struct {
	widget.BaseWidget
	location *city

	date *widget.Select
	time *widget.SelectEntry
	dots *widget.Button
	menu *fyne.Menu
}

func newLocation(loc *city) *location {
	l := &location{location: loc}
	l.ExtendBaseWidget(l)

	l.date = widget.NewSelect([]string{}, func(string) {})
	l.date.PlaceHolder = loc.localTime.Format("Mon 02 Jan")
	l.time = widget.NewSelectEntry(listTimes())
	l.time.PlaceHolder = "22:00" // longest
	l.time.Wrapping = fyne.TextWrapOff
	l.time.SetText(loc.localTime.Format("15:04"))
	l.menu = fyne.NewMenu("",
		fyne.NewMenuItem("Delete Place", func() { fmt.Println("Delete place") }),
		fyne.NewMenuItem("Photo info", func() { fmt.Println("Photo info") }))

	l.dots = widget.NewButton("...", func() {
		position := fyne.CurrentApp().Driver().AbsolutePositionForObject(l.dots)
		position.Y += l.dots.Size().Height

		widget.ShowPopUpMenuAtPosition(l.menu, fyne.CurrentApp().Driver().CanvasForObject(l.dots), position)
	})

	return l
}

func (l *location) CreateRenderer() fyne.WidgetRenderer {
	bg := canvas.NewImageFromResource(theme.FileImageIcon())
	bg.Translucency = 0.5
	city := widget.NewRichTextFromMarkdown("# " + l.location.name)
	location := widget.NewRichTextFromMarkdown("## " + l.location.country + " · " + l.location.localTime.Format("MST"))
	input := container.NewBorder(nil, nil, l.date, l.time)

	c := container.NewMax(bg,
		container.NewBorder(nil,
			container.NewVBox(container.NewHBox(city, layout.NewSpacer(), l.dots), location, input), nil, nil))
	return widget.NewSimpleRenderer(c)
}

func listTimes() (times []string) {
	for hour := 0; hour < 24; hour++ {
		times = append(times,
			fmt.Sprintf("%02d:00", hour), fmt.Sprintf("%02d:15", hour),
			fmt.Sprintf("%02d:30", hour), fmt.Sprintf("%02d:45", hour))
	}
	return times
}
