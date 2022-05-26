package main

import (
	"fmt"
	"io"
	"log"
	"net/http"
	"net/url"
	"time"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/canvas"
	"fyne.io/fyne/v2/storage"
	"github.com/hbagdi/go-unsplash/unsplash"
)

type photo struct {
	cache                 string
	description           string
	photographerName      string
	photographerPortfolio *url.URL
	photoDownloaded       *url.URL
	photoLink             *url.URL
}

type unsplashSession struct {
	storage   fyne.Storage
	store     *cityStore
	client_id string
}

func newUnsplashSession(storage fyne.Storage, store *cityStore) *unsplashSession {
	client_id := secret()
	if client_id == "" {
		return nil
	}
	return &unsplashSession{storage: storage, store: store, client_id: client_id}
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
	return user.PortfolioURL.URL
}

func getUrl(photo unsplash.Photo) *url.URL {
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
	u := unsplash.NewWithClientID(&client, us.client_id)

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
		cache:                 "unsplash-photo-cache-" + city + "-" + country + ".jpg",
		description:           getString((*photos.Results)[0].Description),
		photographerName:      getPhotographerName((*photos.Results)[0].Photographer),
		photographerPortfolio: getPhotographerPortfolio((*photos.Results)[0].Photographer),
		photoDownloaded:       getUrl((*photos.Results)[0]),
		photoLink:             (*photos.Results)[0].Links.Self.URL,
	}, nil
}

func (us *unsplashSession) download(p photo) *canvas.Image {
	if p.photoDownloaded == nil {
		log.Println("No photo download target")
		return nil
	}

	httpResponse, err := http.Get(p.photoDownloaded.String())
	if err != nil {
		fyne.LogError("Unexpected error", err)
		return nil
	}
	defer httpResponse.Body.Close()

	httpResponse.Header.Set("Accept-Version", "v1")
	httpResponse.Header.Set("Authorization", fmt.Sprintf("CLIENT-IS %v", us.client_id))

	childURI, err := storage.Child(us.storage.RootURI(), p.cache)
	if err != nil {
		fyne.LogError("Unexpected error", err)
		return nil
	}
	write, err := storage.Writer(childURI)
	if err != nil {
		fyne.LogError("Unexpected error", err)
		return nil
	}

	reader := io.TeeReader(httpResponse.Body, write)

	return canvas.NewImageFromReader(reader, p.photoDownloaded.String())
}

func (us *unsplashSession) cached(cache string) *canvas.Image {
	childURI, err := storage.Child(us.storage.RootURI(), cache)
	if err != nil {
		fyne.LogError("Unable to get a URI for "+cache, err)
		return nil
	}

	reader, err := storage.Reader(childURI)
	if err != nil {
		fyne.LogError("Unable to open a reader for "+cache, err)
		return nil
	}

	return canvas.NewImageFromReader(reader, cache)
}

func (us *unsplashSession) get(location *city) *canvas.Image {
	var r *canvas.Image

	if location.unsplash.cache != "" {
		r = us.cached(location.unsplash.cache)
	}

	if r == nil {
		metadata, err := us.fetchMetadata(location.name, location.country)
		if err != nil {
			fyne.LogError("Unable to find a picture for "+location.name+"["+location.country+"]", err)
			return nil
		}

		r = us.download(metadata)
		if err != nil {
			fyne.LogError("Unable to create an image for "+location.name+"["+location.country+"]: "+location.unsplash.cache, err)
			return nil
		}
		location.unsplash = metadata
		us.store.save()
	}
	return r
}
