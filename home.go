//go:generate fyne bundle -o assets.go plus-circle.svg

package main

import (
	"fmt"
	"strings"
	"time"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/theme"
	"fyne.io/fyne/v2/widget"
	"github.com/zsefvlol/timezonemapper"
)

func autoCompleteEntry(n *nomad, savedLocations *[]string) *CompletionEntry {
	entry := NewCompletionEntry([]string{})
	entry.SetPlaceHolder("Add a Place")

	cardTexts := []string{}

	entry.OnChanged = func(s string) {

		// completion start for text length >= 2 Some cities have two letter names
		if len(s) < 2 {
			entry.HideCompletion()
			return
		}

		results := []City{}

		for _, value := range Cities {

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

		cardTexts = []string{}
		for _, r := range results {
			timezone := timezonemapper.LatLngToTimezoneString(r.Latitude, r.Longitude)

			s := r.City + "--" + r.Country + "--" + timezone
			cardTexts = append(cardTexts, s)
		}

		// then show them
		entry.SetOptions(cardTexts)
		entry.ShowCompletion()

		entry.OnSubmitted = func(s string) {
			*savedLocations = append(*savedLocations, s)
			entry.SetText("")
			entry.PlaceHolder = "ADD A PLACE"

			fmt.Println(savedLocations)
		}
	}

	return entry
}

func (n *nomad) makeAddCell(savedLocations *[]string) fyne.CanvasObject {

	add := widget.NewIcon(theme.NewThemedResource(resourcePlusCircleSvg))
	search := autoCompleteEntry(n, savedLocations)

	content := container.NewBorder(container.NewBorder(nil, nil, add, nil, search),
		nil, nil, nil)

	return container.NewPadded(content)
}

func savedCells(savedLocations *[]string) []*location {

	cells := []*location{}

	for i, _ := range *savedLocations {

		split := strings.Split((*savedLocations)[i], "--")
		zone, _ := time.LoadLocation(split[2])
		cell := newLocation(split[0], split[1], zone)
		cells = append(cells, cell)
	}
	return cells
}

func (n *nomad) makeHome() fyne.CanvasObject {
	zone, _ := time.LoadLocation("Europe/London")
	cell := newLocation("Edinburgh", "United Kingdom", zone)

	var savedLocations = []string{"City Test--Country--Europe/London"}

	cells := savedCells(&savedLocations)

	return container.New(
		&nomadLayout{},
		cell,
		cells[0],
		n.makeAddCell(&savedLocations))
}
