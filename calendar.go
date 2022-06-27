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

type calendar struct {
	widget.BaseWidget

	monthPrevious *widget.Button
	monthNext     *widget.Button
	monthLabel    *widget.RichText
	canvas        fyne.Canvas

	t time.Time

	day   int
	month int
	year  int

	dates *fyne.Container

	callback func(time.Time)
}

func (c *calendar) daysOfMonth() []fyne.CanvasObject {
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

		dayNum := d.Day()
		s := fmt.Sprint(dayNum)
		var b fyne.CanvasObject = widget.NewButton(s, func() {

			selectedDate := c.dateForButton(dayNum)

			c.setCachedDateInfo(selectedDate)

			c.callback(selectedDate)

			c.hideOverlay()
		})

		buttons = append(buttons, b)
	}

	return buttons
}

func (c *calendar) dateForButton(dayNum int) time.Time {
	oldName, off := globalAppTime.Zone()
	return time.Date(c.year, time.Month(c.month), dayNum, globalAppTime.Hour(), globalAppTime.Minute(), 0, 0, time.FixedZone(oldName, off)).In(c.t.Location())
}

func (c *calendar) hideOverlay() {
	overlayList := c.canvas.Overlays().List()
	overlayList[0].Hide()
}

func (c *calendar) setCachedDateInfo(dateToSet time.Time) {
	c.day = dateToSet.Day()
	c.month = int(dateToSet.Month())
	c.year = dateToSet.Year()
}

func (c *calendar) monthYear() string {
	return time.Month(c.month).String() + " " + strconv.Itoa(c.year)
}

func (c *calendar) fullDate() string {
	d, _ := time.Parse("2006-1-2", strconv.Itoa(c.year)+"-"+strconv.Itoa(c.month)+"-"+strconv.Itoa(c.day))
	return d.Format("Mon 02 Jan 2006")
}

func (c *calendar) calendarObjects() []fyne.CanvasObject {
	textSize := float32(8)
	columnHeadings := []fyne.CanvasObject{}
	for i := 0; i < 7; i++ {
		j := i + 1
		if j == 7 {
			j = 0
		}

		var canvasObject fyne.CanvasObject = canvas.NewText(strings.ToUpper(time.Weekday(j).String()[:3]), color.NRGBA{0xFF, 0xFF, 0xFF, 0xBF})
		canvasObject.(*canvas.Text).TextSize = textSize
		canvasObject.(*canvas.Text).Alignment = fyne.TextAlignCenter
		columnHeadings = append(columnHeadings, canvasObject)
	}
	columnHeadings = append(columnHeadings, c.daysOfMonth()...)

	return columnHeadings
}

func (c *calendar) showAtPos(canvas fyne.Canvas, pos fyne.Position) {
	c.canvas = canvas
	widget.ShowPopUpAtPosition(c, canvas, pos)
}

func (c *calendar) CreateRenderer() fyne.WidgetRenderer {

	c.monthPrevious = widget.NewButtonWithIcon("", theme.NavigateBackIcon(), func() {
		c.month--
		if c.month < 1 {
			c.month = 12
			c.year--
		}
		c.monthLabel.ParseMarkdown(c.monthYear())

		c.dates.Objects = c.calendarObjects()
	})
	c.monthNext = widget.NewButtonWithIcon("", theme.NavigateNextIcon(), func() {
		c.month++
		if c.month > 12 {
			c.month = 1
			c.year++
		}
		c.monthLabel.ParseMarkdown(c.monthYear())

		c.dates.Objects = c.calendarObjects()
	})

	c.monthLabel = widget.NewRichTextFromMarkdown(c.monthYear())

	nav := container.New(layout.NewBorderLayout(nil, nil, c.monthPrevious, c.monthNext),
		c.monthPrevious, c.monthNext, container.NewCenter(c.monthLabel))

	c.dates = container.New(newCalendarLayout(32), c.calendarObjects()...)

	dateContainer := container.NewVBox(nav, c.dates)

	return widget.NewSimpleRenderer(dateContainer)
}

func newCalendar(t time.Time, callback func(time.Time)) *calendar {
	c := &calendar{day: t.Day(), month: int(t.Month()), year: t.Year(), t: t, callback: callback}
	c.ExtendBaseWidget(c)

	return c
}
