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

func newCity(name, country, photoCache, photoDescription, photographerName string, photographerPortfolio, photoDownloaded, photoLink *url.URL, loc *time.Location) *city {
	t := time.Now().In(loc)
	return &city{
		name: name, country: country, localTime: t,
		unsplash: photo{
			photoCache:            photoCache,
			description:           photoDescription,
			photographerName:      photographerName,
			photographerPortfolio: photographerPortfolio,
			photoDownloaded:       photoDownloaded,
			photoLink:             photoLink,
		},
	}
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
			newCity("Edinburgh", "United Kingdom", "", "", "", nil, nil, nil, zone),
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
			photoCache := p.String(prefix + "photoCache")
			photoDescription := p.String(prefix + "photoDescription")
			photographerName := p.String(prefix + "photographerName")
			photographerPortfolio := p.String(prefix + "photographerPortfolio")
			photographerPortfolioUrl, err := url.Parse(photographerPortfolio)
			if err != nil {
				fyne.LogError("Failed to parse photographer portfolio uri: "+photographerPortfolio, err)
			}
			photoDownloaded := p.String(prefix + "photoDownloaded")
			photoDownloadedUrl, err := url.Parse(photoDownloaded)
			if err != nil {
				fyne.LogError("Failed to parse photo uri: "+photoDownloaded, err)
			}
			photoLink := p.String(prefix + "photoLink")
			photoLinkUrl, err := url.Parse(photoLink)
			if err != nil {
				fyne.LogError("Failed to parse photo uri: "+photoLink, err)
			}
			s.list[i] = newCity(name, country, photoCache, photoDescription, photographerName, photographerPortfolioUrl, photoDownloadedUrl, photoLinkUrl, zone)
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
		s.prefs.SetString(prefix+"photoCache", c.unsplash.photoCache)
		s.prefs.SetString(prefix+"photoDescription", c.unsplash.description)
		s.prefs.SetString(prefix+"photographerName", c.unsplash.photographerName)
		if c.unsplash.photographerPortfolio != nil {
			s.prefs.SetString(prefix+"photographerPortfolio", c.unsplash.photographerPortfolio.String())
		}
		if c.unsplash.photoDownloaded != nil {
			s.prefs.SetString(prefix+"photoDownloaded", c.unsplash.photoDownloaded.String())
		}
		if c.unsplash.photoLink != nil {
			s.prefs.SetString(prefix+"photoLink", c.unsplash.photoLink.String())
		}
	}
}
