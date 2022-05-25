package main

import (
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
	s.list = []*city{
		newCity("Edinburgh", "United Kingdom", zone),
	}

	return s
}

func (s *cityStore) cities() []*city {
	return s.list
}
