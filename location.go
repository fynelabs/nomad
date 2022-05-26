package main

import (
	"fmt"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/canvas"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/theme"
	"fyne.io/fyne/v2/widget"
)

type location struct {
	widget.BaseWidget
	location *city
	session  *unsplashSession

	date *widget.Select
	time *widget.SelectEntry
}

func newLocation(loc *city, session *unsplashSession) *location {
	l := &location{location: loc, session: session}
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
	city := widget.NewRichTextFromMarkdown("# " + l.location.name)
	location := widget.NewRichTextFromMarkdown("## " + l.location.country + " · " + l.location.localTime.Format("MST"))
	input := container.NewBorder(nil, nil, l.date, l.time)

	c := container.NewMax(bg,
		container.NewBorder(nil,
			container.NewVBox(city, location, input), nil, nil))

	go func() {
		if l.session == nil {
			return
		}

		unsplashBg := l.session.get(l.location)
		if unsplashBg == nil {
			return
		}

		c.Objects[0] = unsplashBg
		c.Refresh()
	}()

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
