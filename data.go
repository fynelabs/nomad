package main

import (
	"net/url"
	"strconv"
	"time"

	"fyne.io/fyne/v2"
)

const preferenceKeyPrefix = "city."

type city struct {
	name, country string
	unsplash      photo
	localTime     time.Time
}

func newCity(name, country string, unsplash photo, loc *time.Location) *city {
	t := time.Now().In(loc)
	return &city{
		name: name, country: country, localTime: t,
		unsplash: unsplash,
	}
}

type cityStore struct {
	prefs fyne.Preferences
	list  []*city
}

func newCityPhoto(p fyne.Preferences, prefix string) photo {
	cache := p.String(prefix + "cache")
	description := p.String(prefix + "photoDescription")
	photographerName := p.String(prefix + "photographerName")
	photographerPortfolio := p.String(prefix + "photographerPortfolio")
	photographerPortfolioUrl, err := url.Parse(photographerPortfolio)
	if err != nil {
		fyne.LogError("Failed to parse photographer portfolio uri: "+photographerPortfolio, err)
	}
	downloaded := p.String(prefix + "photoDownloaded")
	downloadedUrl, err := url.Parse(downloaded)
	if err != nil {
		fyne.LogError("Failed to parse photo uri: "+downloaded, err)
	}
	link := p.String(prefix + "photoLink")
	linkUrl, err := url.Parse(link)
	if err != nil {
		fyne.LogError("Failed to parse photo uri: "+link, err)
	}

	return photo{
		cache:                 cache,
		description:           description,
		photographerName:      photographerName,
		photographerPortfolio: photographerPortfolioUrl,
		photoDownloaded:       downloadedUrl,
		photoLink:             linkUrl,
	}
}

func newCityStore(p fyne.Preferences) *cityStore {
	s := &cityStore{prefs: p}

	zone, _ := time.LoadLocation("Europe/London")
	count := p.Int(preferenceKeyPrefix + "count")
	if count == 0 {
		s.list = []*city{
			newCity("Edinburgh", "United Kingdom", photo{}, zone),
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
			unsplash := newCityPhoto(p, prefix)
			s.list[i] = newCity(name, country, unsplash, zone)
		}
	}

	return s
}

func (s *cityStore) cities() []*city {
	return s.list
}

func (s *cityStore) savePhoto(prefix string, unsplash *photo) {
	s.prefs.SetString(prefix+"cache", unsplash.cache)
	s.prefs.SetString(prefix+"photoDescription", unsplash.description)
	s.prefs.SetString(prefix+"photographerName", unsplash.photographerName)
	if unsplash.photographerPortfolio != nil {
		s.prefs.SetString(prefix+"photographerPortfolio", unsplash.photographerPortfolio.String())
	}
	if unsplash.photoDownloaded != nil {
		s.prefs.SetString(prefix+"photoDownloaded", unsplash.photoDownloaded.String())
	}
	if unsplash.photoLink != nil {
		s.prefs.SetString(prefix+"photoLink", unsplash.photoLink.String())
	}

}

func (s *cityStore) save() {
	s.prefs.SetInt(preferenceKeyPrefix+"count", len(s.list))

	for i, c := range s.list {
		prefix := preferenceKeyPrefix + strconv.Itoa(i) + "."
		s.prefs.SetString(prefix+"name", c.name)
		s.prefs.SetString(prefix+"country", c.country)
		s.prefs.SetString(prefix+"zone", c.localTime.Location().String())
		s.savePhoto(prefix, &c.unsplash)
	}
}
