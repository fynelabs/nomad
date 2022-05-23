package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/widget"
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

func (n *nomad) makeHome() fyne.CanvasObject {

	input := widget.NewEntry()
	input.SetPlaceHolder("Add a Place")

	input.OnChanged = func(s string) {
		url := "http://geodb-free-service.wirefreethought.com//v1/geo/cities?limit=5&offset=0&namePrefix="
		url += s

		data := &GeoDBdata{}
		getJson(url, data)

		for i := 0; i < len(data.Data); i++ {
			fmt.Println(data.Data[i].Name)
		}
	}

	c := container.NewVBox(widget.NewLabel("Home"),
		input)

	return c
}
