//go:generate fyne bundle -o assets.go plus-circle.svg
//go:generate fyne bundle -append -o assets.go globeSpinnerSplash.gif
//go:generate fyne bundle -append -o assets.go WorkSans-BlackItalic.ttf
//go:generate fyne bundle -append -o assets.go WorkSans-Black.ttf
//go:generate fyne bundle -append -o assets.go WorkSans-Bold.ttf
//go:generate fyne bundle -append -o assets.go WorkSans-Regular.ttf
//go:generate fyne bundle -append -o assets.go Icon.png

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
	"fyne.io/fyne/v2/theme"
	"fyne.io/fyne/v2/widget"

	xWidget "fyne.io/x/fyne/widget"
	"github.com/tidwall/cities"
	"github.com/zsefvlol/timezonemapper"
)

var (
	globalAppTime            = time.Now()
	currentTimeSelected bool = true
)

func (n *nomad) autoCompleteEntry(homeContainer *fyne.Container) *xWidget.CompletionEntry {
	results := []cities.City{}
	entry := xWidget.NewCompletionEntry([]string{})
	entry.SetPlaceHolder("ADD A PLACE")

	entry.CustomCreate = func() fyne.CanvasObject {
		city := widget.NewRichTextFromMarkdown("City Lowercase")
		city.Move(fyne.NewPos(0, 0))

		location := canvas.NewText("COUNTRY - BST", color.NRGBA{0xFF, 0xFF, 0xFF, 0xBF})
		location.TextStyle.Monospace = true
		location.TextSize = 10
		location.Move(fyne.NewPos(theme.Padding()*2, -theme.Padding()*2))

		return container.NewVBox(city, container.NewWithoutLayout(location))

	}
	entry.CustomUpdate = func(i widget.ListItemID, o fyne.CanvasObject) {
		//bug catch
		if i > len(results)-1 {
			fmt.Println("Crashes if not caught here")
			return
		}

		r := results[i]
		timezone := timezonemapper.LatLngToTimezoneString(r.Latitude, r.Longitude)

		c := o.(*fyne.Container)
		city := c.Objects[0].(*widget.RichText)
		city.ParseMarkdown(r.City)

		countryAndTZ := c.Objects[1].(*fyne.Container).Objects[0].(*canvas.Text)
		z, _ := time.LoadLocation(timezone)
		t := time.Now().In(z)
		countryAndTZ.Text = (strings.ToUpper(r.Country) + " · " + t.Format("MST"))
	}

	entry.OnChanged = func(s string) {
		if index, err := strconv.Atoi(s); len(s) > 0 && err == nil {
			if index < len(results) {
				entry.SetText("")

				chooseCity(results[index], n, homeContainer)
			}
			return
		}
		entry.Entry.SetText(strings.ToUpper(entry.Entry.Text))

		// completion start for text length >= 2 Some cities have two letter names
		if len(s) < 2 {
			entry.HideCompletion()
			return
		}

		newResults := []cities.City{}
		for _, value := range cities.Cities {
			if len(value.City) < len(s) {
				continue
			}

			if strings.EqualFold(s, value.City[:len(s)]) {
				newResults = append(newResults, value)
			}
		}

		results = newResults
		if len(newResults) == 0 {
			entry.HideCompletion()
			return
		}

		cardTexts := []string{}
		for i := range results {
			cardTexts = append(cardTexts, strconv.Itoa(i))
		}

		// then show them
		entry.SetOptions(cardTexts)
		entry.ShowCompletion()
	}

	entry.OnSubmitted = func(s string) {
		entry.SetText("")
		if len(results) <= 0 {
			return
		}

		r := results[0]
		chooseCity(r, n, homeContainer)
	}

	return entry
}

func chooseCity(chosen cities.City, n *nomad, homeContainer *fyne.Container) {
	for i := 0; i < len(homeContainer.Objects); i++ {
		l, ok := homeContainer.Objects[i].(*location)
		if !ok {
			continue
		}
		if l.location.name == chosen.City && l.location.country == chosen.Country {
			return
		}
	}

	timezone := timezonemapper.LatLngToTimezoneString(chosen.Latitude, chosen.Longitude)
	zone, _ := time.LoadLocation(timezone)
	c := newCity(chosen.City, chosen.Country, photo{}, zone)

	n.store.add(c)
	l := newLocation(c, n, homeContainer)
	homeContainer.Objects = append(homeContainer.Objects[:len(homeContainer.Objects)-1], l, homeContainer.Objects[len(homeContainer.Objects)-1])
	l.updateMenu()
}

func setDate(dateToSet time.Time, containerObjects []fyne.CanvasObject) {
	globalAppTime = dateToSet

	for i := 0; i < len(containerObjects); i++ {
		loc, ok := containerObjects[i].(*location)
		if !ok {
			continue
		}
		locDate := dateToSet.In(loc.location.localTime.Location())
		loc.updateLocation(locDate)
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
	for _, c := range n.store.list {
		cells = append(cells, newLocation(c, n, homeContainer))
	}

	homeContainer.Objects = append(cells, n.makeAddCell(homeContainer))
	scroll := container.NewVScroll(homeContainer)
	scroll.SetMinSize(layout.minOuterSize())

	startClockTick(homeContainer.Objects)

	return scroll
}

func startClockTick(containerObjects []fyne.CanvasObject) {
	ticker := time.NewTicker(time.Second)
	go func() {
		for t := range ticker.C {
			if !currentTimeSelected {
				continue
			}
			globalAppTime = t
			setDate(t, containerObjects)
		}
	}()
}
