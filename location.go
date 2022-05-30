package main

import (
	"fmt"
	"image/color"
	"strings"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/canvas"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/layout"
	"fyne.io/fyne/v2/theme"
	"fyne.io/fyne/v2/widget"
)

var (
	locationTextColor = color.NRGBA{0xFF, 0xFF, 0xFF, 0xBF}
)

type location struct {
	widget.BaseWidget
	location *city
	session  *unsplashSession

	date *widget.Select
	time *widget.SelectEntry
	dots *fyne.Container
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

	menu := fyne.NewMenu("",
		fyne.NewMenuItem("Delete Place", func() { fmt.Println("Delete place") }),
		fyne.NewMenuItem("Photo info", func() { fmt.Println("Photo info") }))

	l.dots = container.NewVBox(layout.NewSpacer(), newIconWithPopUpMenu(theme.MoreHorizontalIcon(), menu))

	return l
}

func (l *location) CreateRenderer() fyne.WidgetRenderer {
	bg := canvas.NewImageFromResource(theme.FileImageIcon())
	bg.Translucency = 0.5
	city := widget.NewRichTextFromMarkdown("# " + strings.ToUpper(l.location.name))

	location := canvas.NewText(" "+strings.ToUpper(l.location.country)+" · "+l.location.localTime.Format("MST"), locationTextColor)
	location.TextStyle.Monospace = true
	location.TextSize = 10
	location.Move(fyne.NewPos(theme.Padding(), city.MinSize().Height-location.TextSize*.5))
	input := container.NewBorder(nil, nil, l.date, l.time)

	c := container.NewMax(bg,
		container.NewBorder(nil,
			container.NewVBox(container.NewWithoutLayout(container.NewHBox(city, layout.NewSpacer(), l.dots), location), input), nil, nil))

	go func() {
		if l.session == nil {
			return
		}

		unsplashBg, err := l.session.get(l.location)
		if err != nil {
			fyne.LogError("unable to build Unsplash image", err)
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
