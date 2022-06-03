package main

import (
	"fmt"
	"image/color"
	"reflect"
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

type calendar struct {
	widget.BaseWidget

	dateButton    *widget.Button
	monthPrevious *widget.Button
	monthNext     *widget.Button
	monthLabel    *widget.RichText
	canvas        fyne.Canvas

	l *location

	day   int
	month int
	year  int

	dates *fyne.Container
}

func daysOfMonth(c *calendar) []fyne.CanvasObject {
	start, _ := time.Parse("2006-1-2", strconv.Itoa(c.year)+"-"+strconv.Itoa(c.month)+"-"+strconv.Itoa(1))

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

			overlayList := c.canvas.Overlays().List()
			overlayList[0].Hide()
			c.day, _ = strconv.Atoi(s)
			fmt.Println("Date selected  = ", c.day, c.month, c.year)

			//selectedTime
			selectedTime := c.l.time.Text
			fmt.Println("selected time:", selectedTime)
			split := strings.Split(selectedTime, ":")
			hour, _ := strconv.Atoi(split[0])
			minute, _ := strconv.Atoi(split[1])

			//do for the other locations
			for i := 0; i < len(c.l.homeContainer.Objects); i++ {
				if reflect.TypeOf(c.l.homeContainer.Objects[i]) != reflect.TypeOf(c.l) {
					continue
				}
				loc := c.l.homeContainer.Objects[i].(*location)
				loc.dateButton.SetText(fullDate(c.l.calendar)) //this

				//create a date
				date := time.Date(c.year, time.Month(c.month), c.day, hour, minute, 0, 0, loc.location.localTime.Location())

				loc.locationTZLabel.Text = strings.ToUpper(loc.location.country + " · " + date.Format("MST"))
				loc.locationTZLabel.TextStyle.Monospace = true
				loc.locationTZLabel.TextSize = 10
				loc.locationTZLabel.Move(fyne.NewPos(theme.Padding()*2, 40)) //first time clicked this label moves ever so slightly
				loc.locationTZLabel.Refresh()

				//need to find time in time zone

				//and change time
				time := fmt.Sprintf("%02d:%02d", date.Hour(), date.Minute())
				loc.time.SetText(time)

			}
		})

		buttons = append(buttons, b)
	}

	return buttons
}

func monthYear(c *calendar) string {
	return time.Month(c.month).String() + " " + strconv.Itoa(c.year)
}

func fullDate(c *calendar) string {
	d, _ := time.Parse("2006-1-2", strconv.Itoa(c.year)+"-"+strconv.Itoa(c.month)+"-"+strconv.Itoa(c.day))
	return d.Weekday().String()[:3] + " " + strconv.Itoa(d.Day()) + " " + d.Month().String() + " " + strconv.Itoa(d.Year())
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

func calendarObjects(c *calendar) []fyne.CanvasObject {
	cH := columnHeadings(8)
	cH = append(cH, daysOfMonth(c)...)

	return cH
}

func newCalendarPopUpAtPos(c *calendar, l *location, canvas fyne.Canvas, pos fyne.Position) {
	c.canvas = canvas
	c.l = l //I don't want to pass this!

	widget.ShowPopUpAtPosition(c, canvas, pos)
}

func (c *calendar) CreateRenderer() fyne.WidgetRenderer {

	c.monthPrevious = widget.NewButtonWithIcon("", theme.NavigateBackIcon(), func() {
		c.month--
		if c.month < 1 {
			c.month = 12
			c.year--
		}
		c.monthLabel.ParseMarkdown(monthYear(c))

		c.dates.Objects = calendarObjects(c)
	})
	c.monthNext = widget.NewButtonWithIcon("", theme.NavigateNextIcon(), func() {
		c.month++
		if c.month > 12 {
			c.month = 1
			c.year++
		}
		c.monthLabel.ParseMarkdown(monthYear(c))

		c.dates.Objects = calendarObjects(c)
	})

	c.monthLabel = widget.NewRichTextFromMarkdown(monthYear(c))

	nav := container.New(layout.NewBorderLayout(nil, nil, c.monthPrevious, c.monthNext),
		c.monthPrevious, c.monthNext, container.NewCenter(c.monthLabel))

	c.dates = container.New(NewCalendarLayout(32), calendarObjects(c)...)

	dateContainer := container.NewVBox(nav, c.dates)
	return widget.NewSimpleRenderer(dateContainer)
}

func newCalendar() *calendar {

	c := &calendar{day: time.Now().Day(), month: int(time.Now().Month()), year: time.Now().Year()}
	c.ExtendBaseWidget(c)

	return c
}
