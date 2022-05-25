package main

import (
	"strconv"
	"time"

	"fyne.io/fyne/v2"
)

const preferenceKeyPrefix = "city."

type city struct {
	name, country string
	localTime     time.Time
}

func newCity(name, country string, loc *time.Location) *city {
	t := time.Now().In(loc)
	return &city{name: name, country: country, localTime: t}
}

type cityStore struct {
	prefs fyne.Preferences
	list  []*city
}

func newCityStore(p fyne.Preferences) *cityStore {
	s := &cityStore{prefs: p}

	zone, _ := time.LoadLocation("Europe/London")
	count := p.Int(preferenceKeyPrefix + "count")
	if count == 0 {
		s.list = []*city{
			newCity("Edinburgh", "United Kingdom", zone),
		}
		s.save()
	} else {
		s.list = make([]*city, count)
		for i := 0; i < count; i++ {
			prefix := preferenceKeyPrefix + strconv.Itoa(i) + "."
			name := p.StringWithFallback(prefix+"name", "No City")
			country := p.StringWithFallback(prefix+"country", "No Country")
			zoneName := p.StringWithFallback(prefix+"zone", "UTC")
			zone, err := time.LoadLocation(zoneName)
			if err != nil {
				fyne.LogError("Failed to load timezone "+zoneName, err)
			}
			s.list[i] = newCity(name, country, zone)
		}
	}

	return s
}

func (s *cityStore) cities() []*city {
	return s.list
}

func (s *cityStore) save() {
	s.prefs.SetInt(preferenceKeyPrefix+"count", len(s.list))

	for i, c := range s.list {
		prefix := preferenceKeyPrefix + strconv.Itoa(i) + "."
		s.prefs.SetString(prefix+"name", c.name)
		s.prefs.SetString(prefix+"country", c.country)
		s.prefs.SetString(prefix+"zone", c.localTime.Location().String())
	}
}
