package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"

	//"fyne.io/fyne/v2/widget"
	"fyne.io/x/fyne/widget"
)

type GeoDBdata struct {
	Data []GeoDBbody
}
type GeoDBbody struct {
	Id          int
	WikiDataId  string
	Kind        string `json:"type"`
	Name        string
	Country     string
	CountryCode string
	Region      string
	RegionCode  int
	Latitude    float32
	Longitude   float32
}

var myClient = &http.Client{Timeout: 10 * time.Second}

func getJson(url string, target interface{}) error {
	r, err := myClient.Get(url)
	if err != nil {
		return err
	}
	defer r.Body.Close()

	return json.NewDecoder(r.Body).Decode(target)
}

func geoDBAPI(s string) {
	url := "http://geodb-free-service.wirefreethought.com//v1/geo/cities?limit=5&offset=0&namePrefix="
	url += s

	data := &GeoDBdata{}
	getJson(url, data)

	for i := 0; i < len(data.Data); i++ {
		fmt.Println(data.Data[i].Name)
	}
}

func autoCompleteEntry() *widget.CompletionEntry {
	entry := widget.NewCompletionEntry([]string{})
	entry.SetPlaceHolder("Add a Place")

	entry.OnChanged = func(s string) {
		// completion start for text length >= 2 Some cities have two letter names
		if len(s) < 2 {
			entry.HideCompletion()
			return
		}

		//Get the list of possible completion
		c := []City{}

		if len(s) < 2 {
			return
		}

		for _, value := range cities {

			if len(value.City) < len(s) {
				continue
			}

			if s == value.City[:len(s)] {
				c = append(c, value)
			}
		}

		// no results
		if len(c) == 0 {
			entry.HideCompletion()
			return
		}

		results := []string{}
		for _, value := range c {
			results = append(results, value.City)
		}

		// then show them
		entry.SetOptions(results)
		entry.ShowCompletion()
	}

	return entry
}

func (n *nomad) makeHome() fyne.CanvasObject {

	c := container.NewVBox(autoCompleteEntry())

	return c
}
