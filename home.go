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

func autoCompleteEntry(listWidget widget.List, savedLocations *[]string) *CompletionEntry {
	entry := NewCompletionEntry([]string{})
	entry.SetPlaceHolder("Add a Place")

	cardTexts := []string{}

	entry.OnChanged = func(s string) {

		// completion start for text length >= 2 Some cities have two letter names
		if len(s) < 2 {
			entry.HideCompletion()
			return
		}

		//Get the list of possible completion
		results := []City{}
		cardTexts = []string{}

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

		for _, r := range results {
			timezone := timezonemapper.LatLngToTimezoneString(r.Latitude, r.Longitude)

			loc, _ := time.LoadLocation(timezone)

			t, _ := time.ParseInLocation("2006-01-02 15:04:05", "2010-01-01 00:00:00", loc)

			//second return is offset in seconds, useful?
			tzName, _ := t.Zone()

			s := r.City + "--" + r.Country + "--" + tzName
			cardTexts = append(cardTexts, s)
		}

		// then show them
		entry.SetOptions(cardTexts)
		entry.ShowCompletion()

		entry.OnSubmitted = func(s string) {
			fmt.Println("submitted with " + s)

			//add new card to list
			*savedLocations = append(*savedLocations, s)

			listWidget.Refresh()

			fmt.Println(savedLocations)
		}
	}

	return entry
}

func savedLocationsWidget(savedLocations *[]string) *widget.List {
	list := widget.NewList(
		func() int {
			return len(*savedLocations)
		},
		func() fyne.CanvasObject {
			return widget.NewLabel("template")
		},
		func(i widget.ListItemID, o fyne.CanvasObject) {
			o.(*widget.Label).SetText((*savedLocations)[i])
		})

	return list
}

func (n *nomad) makeAddCell(listWidget widget.List, savedLocations *[]string) fyne.CanvasObject {

	add := widget.NewIcon(theme.NewThemedResource(resourcePlusCircleSvg))
	search := autoCompleteEntry(listWidget, savedLocations)
	search.PlaceHolder = "ADD A PLACE"

	content := container.NewBorder(container.NewBorder(nil, nil, add, nil, search),
		nil, nil, nil)

	return container.NewPadded(content)

}

func (n *nomad) makeHome() fyne.CanvasObject {
	zone, _ := time.LoadLocation("Europe/London")
	cell := newLocation("Edinburgh", "United Kingdom", zone)

	var savedLocations = []string{"waaa", "baaa", "traaa"}
	list := savedLocationsWidget(&savedLocations)

	return container.New(&nomadLayout{},
		cell, n.makeAddCell(*list, &savedLocations), list)
}
