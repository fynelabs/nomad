//go:generate fyne bundle -o assets.go plus-circle.svg

package main

import (
	"strings"
	"time"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/theme"
	"fyne.io/fyne/v2/widget"
	xWidget "fyne.io/x/fyne/widget"
	"github.com/zsefvlol/timezonemapper"
)

func autoCompleteEntry() *xWidget.CompletionEntry {
	entry := xWidget.NewCompletionEntry([]string{})
	entry.SetPlaceHolder("Add a Place")

	entry.OnChanged = func(s string) {

		// completion start for text length >= 2 Some cities have two letter names
		if len(s) < 2 {
			entry.HideCompletion()
			return
		}

		//Get the list of possible completion
		results := []City{}

		for _, value := range cities {

			if len(value.City) < len(s) {
				continue
			}

			if strings.EqualFold(s, value.City[:len(s)]) {
				results = append(results, value)
			}
		}

		if len(results) == 0 {
			entry.HideCompletion()
			return
		}

		cardTexts := []string{}

		for _, r := range results {
			timezone := timezonemapper.LatLngToTimezoneString(r.Latitude, r.Longitude)

			loc, _ := time.LoadLocation(timezone)

			t, _ := time.ParseInLocation("2006-01-02 15:04:05", "2010-01-01 00:00:00", loc)

			//second return is offset in seconds, useful?
			tzName, _ := t.Zone()

			s := r.City + " " + r.Country + " " + tzName
			cardTexts = append(cardTexts, s)
		}

		// then show them
		entry.SetOptions(cardTexts)
		entry.ShowCompletion()

	}

	return entry
}

func (n *nomad) makeAddCell() fyne.CanvasObject {
	add := widget.NewIcon(theme.NewThemedResource(resourcePlusCircleSvg))
	search := autoCompleteEntry()
	search.PlaceHolder = "Add a place"
	content := container.NewBorder(container.NewBorder(nil, nil, add, nil, search),
		nil, nil, nil)
	return container.NewPadded(content)
}

func (n *nomad) makeHome() fyne.CanvasObject {
	zone, _ := time.LoadLocation("Europe/London")
	cell := newLocation("Edinburgh", "United Kingdom", zone)

	return container.New(&nomadLayout{},
		cell, n.makeAddCell())
}
