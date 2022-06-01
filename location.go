package main

import (
	"fmt"
	"image/color"
	"strconv"
	"strings"
	"time"

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

	selectedDay   int
	selectedMonth int
	selectedYear  int

	dateButton    *widget.Button
	monthLabel    *widget.RichText
	monthPrevious *widget.Button
	monthNext     *widget.Button

	dateContainer *fyne.Container
	calendar      *fyne.Container
}

func daysOfMonth(l *location) []fyne.CanvasObject {
	start, _ := time.Parse("2006-1-2", strconv.Itoa(l.selectedYear)+"-"+strconv.Itoa(l.selectedMonth)+"-"+strconv.Itoa(1))

	buttons := []fyne.CanvasObject{}

	//account for Go time pkg starting on sunday at index 0
	dayIndex := int(start.Weekday())
	if dayIndex == 0 {
		dayIndex += 7
	}

	//add spacers if week doesn't start on Monday
	for i := 0; i < dayIndex-1; i++ {
		buttons = append(buttons, layout.NewSpacer())
	}

	for d := start; d.Month() == start.Month(); d = d.AddDate(0, 0, 1) {

		s := fmt.Sprint(d.Day())
		var b fyne.CanvasObject = widget.NewButton(s, func() {
			//functionality for task #12 "Change time using calendar and time picker affecting all city"
			//to go here
			fmt.Println("Date selected  = "+s, d.Month(), d.Year())
		})

		buttons = append(buttons, b)
	}

	return buttons
}

func monthYear(l *location) string {
	return time.Month(l.selectedMonth).String() + " " + strconv.Itoa(l.selectedYear)
}

func dayMonthYear(l *location) string {
	d, _ := time.Parse("2006-1-2", strconv.Itoa(l.selectedYear)+"-"+strconv.Itoa(l.selectedMonth)+"-"+strconv.Itoa(l.selectedDay))
	return d.Weekday().String()[:3] + " " + d.Month().String() + " " + strconv.Itoa(d.Year())
}

func columnHeadings(textSize float32) []fyne.CanvasObject {
	l := []fyne.CanvasObject{}
	for i := 0; i < 7; i++ {
		j := i + 1
		if j == 7 {
			j = 0
		}

		var canvasObject fyne.CanvasObject = canvas.NewText(strings.ToUpper(time.Weekday(j).String()[:3]), color.NRGBA{0xFF, 0xFF, 0xFF, 0xBF})
		canvasObject.(*canvas.Text).TextSize = textSize
		canvasObject.(*canvas.Text).Alignment = fyne.TextAlignCenter
		l = append(l, canvasObject)
	}

	return l
}

func calendarObjects(l *location) []fyne.CanvasObject {
	c := columnHeadings(8)
	c = append(c, daysOfMonth(l)...)

	return c
}

func navigation(l *location) {

	l.monthPrevious = widget.NewButtonWithIcon("", theme.NavigateBackIcon(), func() {
		l.selectedMonth--
		if l.selectedMonth < 1 {
			l.selectedMonth = 12
			l.selectedYear--
		}
		l.monthLabel.ParseMarkdown(monthYear(l))

		l.calendar.Objects = calendarObjects(l)
	})
	l.monthNext = widget.NewButtonWithIcon("", theme.NavigateNextIcon(), func() {
		l.selectedMonth++
		if l.selectedMonth > 12 {
			l.selectedMonth = 1
			l.selectedYear++
		}
		l.monthLabel.ParseMarkdown(monthYear(l))

		l.calendar.Objects = calendarObjects(l)
	})

	l.monthLabel = widget.NewRichTextFromMarkdown(monthYear(l))
}

func calendar(l *location) {
	l.calendar = container.New(NewCalendarLayout(32), calendarObjects(l)...)

	b := container.New(layout.NewBorderLayout(nil, nil, l.monthPrevious, l.monthNext),
		l.monthPrevious, l.monthNext, container.NewCenter(l.monthLabel))

	l.dateContainer = container.NewVBox(b, l.calendar)
}

func newLocation(loc *city, session *unsplashSession, n *nomad) *location {
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

	l.selectedDay = time.Now().Day()
	l.selectedMonth = int(time.Now().Month())
	l.selectedYear = time.Now().Year()

	navigation(l)
	calendar(l)

	l.dateButton = widget.NewButton(dayMonthYear(l), func() {
		widget.ShowPopUpAtPosition(l.dateContainer, n.main.Canvas(), fyne.NewPos(0, l.Size().Height)) //wait for merge to use "homeContainer" for second object position
	})

	return l
}

func (l *location) CreateRenderer() fyne.WidgetRenderer {
	bg := canvas.NewImageFromResource(theme.FileImageIcon())
	bg.Translucency = 0.5
	city := widget.NewRichTextFromMarkdown("# " + l.location.name)
	location := widget.NewRichTextFromMarkdown("## " + l.location.country + " · " + l.location.localTime.Format("MST"))

	input := container.NewBorder(nil, nil, l.dateButton, l.time)
	c := container.NewMax(bg,
		container.NewBorder(nil,
			container.NewVBox(city, location, input), nil, nil))

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
