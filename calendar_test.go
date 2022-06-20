package main

import (
	"strconv"
	"testing"
	"time"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/test"
	"fyne.io/fyne/v2/widget"
	"github.com/stretchr/testify/assert"
)

func TestNewCalendar(t *testing.T) {
	now := time.Now()
	c := newCalendar(now, func(time.Time) {})
	assert.Equal(t, now.Day(), c.day)
	assert.Equal(t, int(now.Month()), c.month)
	assert.Equal(t, now.Year(), c.year)

	_ = test.WidgetRenderer(c) // and render
	assert.Equal(t, now.Format("January 2006"), c.monthLabel.Segments[0].(*widget.TextSegment).Text)
}

func TestNewCalendar_ButtonDate(t *testing.T) {
	date := time.Now()
	c := newCalendar(date, func(time.Time) {})
	_ = test.WidgetRenderer(c) // and render

	endNextMonth := date.AddDate(0, 1, 0).AddDate(0, 0, -(date.Day() - 1))
	last := endNextMonth.AddDate(0, 0, -1)

	var firstDate *widget.Button
	for _, b := range c.dates.Objects {
		if nonBlank, ok := b.(*widget.Button); ok {
			firstDate = nonBlank
			break
		}
	}

	assert.Equal(t, "1", firstDate.Text)
	lastDate := c.dates.Objects[len(c.dates.Objects)-1].(*widget.Button)
	assert.Equal(t, strconv.Itoa(last.Day()), lastDate.Text)
}

func TestNewCalendar_Next(t *testing.T) {
	date := time.Now()
	c := newCalendar(date, func(time.Time) {})
	_ = test.WidgetRenderer(c) // and render

	assert.Equal(t, date.Format("January 2006"), c.monthLabel.Segments[0].(*widget.TextSegment).Text)

	test.Tap(c.monthNext)
	date = date.AddDate(0, 1, 0)
	assert.Equal(t, date.Format("January 2006"), c.monthLabel.Segments[0].(*widget.TextSegment).Text)
}

func TestNewCalendar_Previous(t *testing.T) {
	date := time.Now()
	c := newCalendar(date, func(time.Time) {})
	_ = test.WidgetRenderer(c) // and render

	assert.Equal(t, date.Format("January 2006"), c.monthLabel.Segments[0].(*widget.TextSegment).Text)

	test.Tap(c.monthPrevious)
	date = date.AddDate(0, -1, 0)
	assert.Equal(t, date.Format("January 2006"), c.monthLabel.Segments[0].(*widget.TextSegment).Text)
}

func TestNewCalendar_ShowHide(t *testing.T) {
	win := test.NewCanvas()
	c := newCalendar(time.Now(), func(time.Time) {})
	pos := fyne.NewPos(10, 10)

	c.showAtPos(win, pos)
	assert.Equal(t, 1, len(win.Overlays().List()))

	c.hideOverlay()
	assert.Equal(t, 0, len(win.Overlays().List()))
}
