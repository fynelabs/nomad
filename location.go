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

	time   *widget.SelectEntry
	button *widget.Button
	dots   *fyne.Container

	dateButton *widget.Button

	calendar *calendar
}

func newLocation(loc *city, session *unsplashSession, canvas fyne.Canvas) *location {
	l := &location{location: loc, session: session}
	l.ExtendBaseWidget(l)

	l.time = widget.NewSelectEntry(listTimes())
	l.time.PlaceHolder = "22:00" // longest
	l.time.Wrapping = fyne.TextWrapOff
	l.time.SetText(loc.localTime.Format("15:04"))

	menu := fyne.NewMenu("",
		fyne.NewMenuItem("Delete Place", func() { fmt.Println("Delete place") }),
		fyne.NewMenuItem("Photo info", func() { fmt.Println("Photo info") }))

	l.button = widget.NewButtonWithIcon("", theme.MoreHorizontalIcon(), func() {
		position := fyne.CurrentApp().Driver().AbsolutePositionForObject(l.button)
		position.Y += l.button.Size().Height

		widget.ShowPopUpMenuAtPosition(menu, fyne.CurrentApp().Driver().CanvasForObject(l.button), position)
	})
	l.button.Importance = widget.LowImportance

	l.dots = container.NewVBox(layout.NewSpacer(), l.button)

	l.calendar = newCalendar()

	l.dateButton = widget.NewButtonWithIcon(dayMonthYear(l.calendar), theme.MenuDropDownIcon(), func() {
		newCalendarPopUpAtPos(l.calendar, canvas, fyne.NewPos(0, l.Size().Height))
	})
	l.dateButton.Alignment = widget.ButtonAlignLeading
	l.dateButton.IconPlacement = widget.ButtonIconTrailingText
	l.dateButton.Importance = widget.LowImportance

	return l
}

func (l *location) CreateRenderer() fyne.WidgetRenderer {
	bg := canvas.NewImageFromResource(theme.FileImageIcon())
	bg.Translucency = 0.5
	city := widget.NewRichTextFromMarkdown("# " + l.location.name)
	location := canvas.NewText(" "+strings.ToUpper(l.location.country)+" · "+l.location.localTime.Format("MST"), locationTextColor)
	location.TextStyle.Monospace = true
	location.TextSize = 10
	location.Move(fyne.NewPos(theme.Padding(), city.MinSize().Height-location.TextSize*.5))
	input := container.NewBorder(nil, nil, l.dateButton, l.time)

	c := container.NewMax(bg,
		container.NewBorder(nil,
			container.NewVBox(container.NewHBox(container.NewWithoutLayout(city, location), layout.NewSpacer(), l.dots), input), nil, nil))

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
