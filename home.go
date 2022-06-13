//go:generate fyne bundle -o assets.go plus-circle.svg
//go:generate fyne bundle -append -o assets.go globeSpinnerSplash.gif
//go:generate fyne bundle -append -o assets.go WorkSans-BlackItalic.ttf
//go:generate fyne bundle -append -o assets.go WorkSans-Black.ttf
//go:generate fyne bundle -append -o assets.go WorkSans-Bold.ttf
//go:generate fyne bundle -append -o assets.go WorkSans-Regular.ttf

package main

import (
	"errors"
	"strings"
	"time"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/theme"
	"fyne.io/fyne/v2/widget"
	"github.com/tidwall/cities"
	"github.com/zsefvlol/timezonemapper"
)

var (
	globalAppTime            = time.Now()
	currentTimeSelected bool = true
)

func (n *nomad) autoCompleteEntry(homeContainer *fyne.Container) *CompletionEntry {

	entry := NewCompletionEntry([]string{})
	entry.SetPlaceHolder("Add a Place")

	entry.OnChanged = func(s string) {

		// completion start for text length >= 2 Some cities have two letter names
		if len(s) < 2 {
			entry.HideCompletion()
			return
		}

		results := []cities.City{}

		for _, value := range cities.Cities {

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

			s := r.City + "--" + r.Country + "--" + timezone
			cardTexts = append(cardTexts, s)
		}

		// then show them
		entry.SetOptions(cardTexts)
		entry.ShowCompletion()

		entry.Entry.SetText(strings.ToUpper(entry.Entry.Text))
	}

	entry.OnSubmitted = func(s string) {

		//reset entry
		entry.SetText("")
		entry.PlaceHolder = "ADD A PLACE"
		split := strings.Split(s, "--")
		for i := 0; i < len(homeContainer.Objects); i++ {
			l, ok := homeContainer.Objects[i].(*location)
			if !ok {
				continue
			}
			if l.location.name == split[0] && l.location.country == split[1] {
				return
			}
		}

		if len(split) != 3 {
			err1 := errors.New(s)
			fyne.LogError("Location entry string incorrect format", err1)
		} else {
			zone, _ := time.LoadLocation(split[2])
			c := newCity(split[0], split[1], photo{}, zone)

			n.store.list = append(n.store.list, c)
			n.store.save()

			l := newLocation(c, n, homeContainer)
			homeContainer.Objects = append(homeContainer.Objects[:len(homeContainer.Objects)-1], l, homeContainer.Objects[len(homeContainer.Objects)-1])
		}
	}

	return entry
}

func setDate(dateToSet time.Time, containerObjects []fyne.CanvasObject) {
	globalAppTime = dateToSet

	for i := 0; i < len(containerObjects); i++ {
		loc, ok := containerObjects[i].(*location)
		if !ok {
			continue
		}
		locDate := dateToSet.In(loc.location.localTime.Location())
		loc.updateCountry(locDate)
	}
}

func (n *nomad) makeAddCell(homeContainer *fyne.Container) fyne.CanvasObject {

	add := widget.NewIcon(theme.NewPrimaryThemedResource(resourcePlusCircleSvg))
	search := n.autoCompleteEntry(homeContainer)

	content := container.NewBorder(container.NewBorder(nil, nil, add, nil, search),
		nil, nil, nil)

	return container.NewPadded(content)
}

func (n *nomad) makeHome() fyne.CanvasObject {
	layout := &nomadLayout{}
	homeContainer := container.New(layout)

	cells := []fyne.CanvasObject{}
	for _, c := range n.store.cities() {
		cells = append(cells, newLocation(c, n, homeContainer))
	}

	homeContainer.Objects = append(cells, n.makeAddCell(homeContainer))
	scroll := container.NewVScroll(homeContainer)
	scroll.SetMinSize(layout.minOuterSize())

	clockTick(homeContainer.Objects)

	return scroll
}

func clockTick(containerObjects []fyne.CanvasObject) {
	ticker := time.NewTicker(time.Second)
	go func() {
		for t := range ticker.C {
			if !currentTimeSelected {
				continue
			}
			globalAppTime = t
			for i := 0; i < len(containerObjects); i++ {
				l, ok := containerObjects[i].(*location)
				if !ok {
					continue
				}
				local := t.In(l.location.localTime.Location()).Format("15:04")
				l.time.Text = local
				l.time.Refresh()
			}
		}
	}()
}
