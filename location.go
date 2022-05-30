package main

import (
	"fmt"
	"image/color"
	"strings"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/canvas"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/theme"
	"fyne.io/fyne/v2/widget"
)

type location struct {
	widget.BaseWidget
	location *city

	date *widget.Select
	time *widget.SelectEntry
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

	return l
}

func (l *location) CreateRenderer() fyne.WidgetRenderer {
	bg := canvas.NewImageFromResource(theme.FileImageIcon())
	bg.Translucency = 0.5
	city := widget.NewRichTextFromMarkdown("# " + strings.ToUpper(l.location.name))

	var locationTextColor = color.NRGBA{0xFF, 0xFF, 0xFF, 0xBF}
	location := canvas.NewText(" "+strings.ToUpper(l.location.country)+" · "+l.location.localTime.Format("MST"), locationTextColor)
	location.TextStyle.Monospace = true
	location.TextSize = 10
	location.Move(fyne.NewPos(theme.Padding(), city.MinSize().Height-location.TextSize*.5))
	input := container.NewBorder(nil, nil, l.date, l.time)

	c := container.NewMax(bg,
		container.NewBorder(nil,
			container.NewVBox(container.NewWithoutLayout(city, location), input), nil, nil))
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
