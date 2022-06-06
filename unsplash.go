package main

import (
	"fmt"
	"image/color"
	"io"
	"net/http"
	"net/url"
	"time"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/canvas"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/storage"
	"fyne.io/fyne/v2/theme"
	"fyne.io/fyne/v2/widget"
	"github.com/hbagdi/go-unsplash/unsplash"
)

type photo struct {
	cache            string
	description      string
	photographerName string
	portfolio        *url.URL
	original         *url.URL
	full             *url.URL
	photoWebsite     *url.URL
}

type unsplashSession struct {
	storage  fyne.Storage
	store    *cityStore
	clientID string
}

func newUnsplashSession(storage fyne.Storage, store *cityStore) *unsplashSession {
	clientID := secret()
	if clientID == "" {
		return nil
	}
	return &unsplashSession{storage: storage, store: store, clientID: clientID}
}

func getString(str *string) string {
	if str == nil {
		return ""
	}
	return *str
}

func getPhotographerName(user *unsplash.User) string {
	if user == nil {
		return "unknown photographer"
	}
	return getString(user.Name)
}

func getPhotographerPortfolio(user *unsplash.User) *url.URL {
	if user == nil {
		return nil
	}
	if user.PortfolioURL == nil {
		return nil
	}
	return user.PortfolioURL.URL
}

func getURL(photo unsplash.Photo) *url.URL {
	if photo.Urls.Small != nil {
		return photo.Urls.Small.URL
	}
	if photo.Urls.Regular != nil {
		return photo.Urls.Regular.URL
	}
	if photo.Urls.Full != nil {
		return photo.Urls.Full.URL
	}
	return nil
}

func (us *unsplashSession) fetchMetadata(city string, country string) (photo, error) {
	client := http.Client{Timeout: time.Duration(60) * time.Second}
	//use the http.Client to instantiate unsplash
	u := unsplash.NewWithClientID(&client, us.clientID)

	opt := unsplash.SearchOpt{
		Page:    1,
		PerPage: 1,
		Query:   city,
	}
	photos, _, err := u.Search.Photos(&opt)
	if err != nil {
		return photo{}, err
	}
	if *photos.Total == 0 {
		opt.Query = country
		photos, _, err = u.Search.Photos(&opt)
		if err != nil {
			return photo{}, err
		}
	}

	return photo{
		cache:            "unsplash-photo-cache-" + city + "-" + country + ".jpg",
		description:      getString((*photos.Results)[0].Description),
		photographerName: getPhotographerName((*photos.Results)[0].Photographer),
		portfolio:        getPhotographerPortfolio((*photos.Results)[0].Photographer),
		original:         getURL((*photos.Results)[0]),
		full:             (*photos.Results)[0].Urls.Full.URL,
		photoWebsite:     (*photos.Results)[0].Links.HTML.URL,
	}, nil
}

func (us *unsplashSession) download(p photo) (*canvas.Image, error) {
	if p.original == nil {
		return nil, fmt.Errorf("no photo download target")
	}

	httpResponse, err := http.Get(p.original.String())
	if err != nil {
		return nil, err
	}
	defer httpResponse.Body.Close()

	httpResponse.Header.Set("Accept-Version", "v1")
	httpResponse.Header.Set("Authorization", fmt.Sprintf("CLIENT-IS %v", us.clientID))

	childURI, err := storage.Child(us.storage.RootURI(), p.cache)
	if err != nil {
		return nil, err
	}
	write, err := storage.Writer(childURI)
	if err != nil {
		return nil, err
	}

	reader := io.TeeReader(httpResponse.Body, write)
	return canvasImage(reader, p.original.String()), nil
}

func (us *unsplashSession) cached(cache string) (*canvas.Image, error) {
	childURI, err := storage.Child(us.storage.RootURI(), cache)
	if err != nil {
		return nil, err
	}

	reader, err := storage.Reader(childURI)
	if err != nil {
		return nil, err
	}

	return canvasImage(reader, cache), nil
}

func (us *unsplashSession) get(location *city) (*canvas.Image, error) {
	if location.unsplash.cache != "" {
		r, err := us.cached(location.unsplash.cache)
		if r != nil {
			return r, nil
		}
		if err != nil {
			fyne.LogError("existing cache corrupted", err)
		}
	}

	metadata, err := us.fetchMetadata(location.name, location.country)
	if err != nil {
		return nil, err
	}

	r, err := us.download(metadata)
	if err != nil {
		return nil, err
	}
	location.unsplash = metadata
	us.store.save()

	return r, nil
}

func canvasImage(r io.Reader, name string) *canvas.Image {
	img := canvas.NewImageFromReader(r, name)
	img.ScaleMode = canvas.ImageScaleFastest
	img.Translucency = 0.15
	return img
}

func (city city) newInfoScreen(c fyne.Canvas) fyne.CanvasObject {
	photographer := canvas.NewText("Photographer", locationTextColor)
	photographer.TextStyle.Monospace = true
	photographer.TextSize = 10
	photographerName := canvas.NewText(city.unsplash.photographerName, color.White)

	location := canvas.NewText("Location", locationTextColor)
	location.TextStyle.Monospace = true
	location.TextSize = 10
	cityCountry := canvas.NewText(city.name+", "+city.country, color.White)

	linkImage := widget.NewHyperlink("View on unsplash", city.unsplash.photoWebsite)

	overlay := container.NewMax()

	bg := canvas.NewImageFromResource(theme.FileImageIcon())
	bg.SetMinSize(c.Size())
	overlay.Add(bg)

	go func() {
		if city.unsplash.full == nil {
			return
		}
		httpResponse, err := http.Get(city.unsplash.full.String())
		if err != nil {
			fyne.LogError("Unable to download full image", err)
			return
		}
		defer httpResponse.Body.Close()

		bg = canvasImage(httpResponse.Body, city.unsplash.full.String())
		fmt.Println("x", bg.Size())

		bg.FillMode = canvas.ImageFillStretch

		overlay.Objects[0] = bg
		//defaults to 0.15 translucency
		overlay.Objects[0].(*canvas.Image).Translucency = 0
	}()

	exitButton := widget.NewButtonWithIcon("", theme.CancelIcon(), func() {
		c.Overlays().Remove(overlay)
	})
	exitButton.Importance = widget.LowImportance

	pulldown := container.NewBorder(
		container.NewMax(canvas.NewRectangle(&color.NRGBA{0x00, 0x00, 0x00, 0x80}),
			container.NewBorder(nil, nil, nil, container.NewBorder(exitButton, nil, nil, nil)),
			container.NewVBox(photographer, photographerName, location, cityCountry, linkImage)), nil, nil, nil)

	overlay.Add(pulldown)
	return overlay
}
