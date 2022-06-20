package main

import (
	"strings"
	"testing"
	"time"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/test"

	"github.com/stretchr/testify/assert"
)

func TestNewLocation(t *testing.T) {
	n := &nomad{}
	c := newCity("Test", "Country", photo{}, time.UTC)
	l := newLocation(c, n, fyne.NewContainer())
	_ = test.WidgetRenderer(l)

	utc := time.Now().In(time.UTC)
	assert.Equal(t, utc.Format("15:04"), l.time.Text)
	assert.True(t, strings.Contains(l.locationTZLabel.Text, " UTC"))
}

func TestLocation_PickTime(t *testing.T) {
	n := &nomad{}
	c := newCity("Test", "Country", photo{}, time.UTC)
	l := newLocation(c, n, fyne.NewContainer())
	_ = test.WidgetRenderer(l)

	zone, _ := time.LoadLocation("EST")
	inEST := newCity("City", "America", photo{}, zone)
	l2 := newLocation(inEST, n, l.homeContainer)
	_ = test.WidgetRenderer(l2)
	l.homeContainer.Objects = append(l.homeContainer.Objects, l2)

	l.time.SetText("12:00")
	assert.Equal(t, "07:00", l2.time.Text)
}
