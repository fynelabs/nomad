package main

import (
	"fmt"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/canvas"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/storage"
	"fyne.io/fyne/v2/theme"
	"fyne.io/fyne/v2/widget"
)

type location struct {
	widget.BaseWidget
	location *city
	n        *nomad

	date *widget.Select
	time *widget.SelectEntry
}

func newLocation(loc *city, n *nomad) *location {
	l := &location{location: loc, n: n}
	l.ExtendBaseWidget(l)

	l.date = widget.NewSelect([]string{}, func(string) {})
	l.date.PlaceHolder = loc.localTime.Format("Mon 02 Jan")
	l.time = widget.NewSelectEntry(listTimes())
	l.time.PlaceHolder = "22:00" // longest
	l.time.Wrapping = fyne.TextWrapOff
	l.time.SetText(loc.localTime.Format("15:04"))

	return l
}

func (l *location) CreateRenderer() fyne.WidgetRenderer {
	bg := canvas.NewImageFromResource(theme.FileImageIcon())
	bg.Translucency = 0.5
	city := widget.NewRichTextFromMarkdown("# " + l.location.name)
	location := widget.NewRichTextFromMarkdown("## " + l.location.country + " · " + l.location.localTime.Format("MST"))
	input := container.NewBorder(nil, nil, l.date, l.time)

	c := container.NewMax(bg,
		container.NewBorder(nil,
			container.NewVBox(city, location, input), nil, nil))

	go func() {
		var unsplashBg *canvas.Image

		if l.location.unsplash.photoCache != "" {
			childURI, err := storage.Child(l.n.storage.RootURI(), l.location.unsplash.photoCache)
			if err != nil {
				fyne.LogError("Unable to get a URI for "+l.location.unsplash.photoCache, err)
			} else {
				reader, err := storage.Reader(childURI)
				if err != nil {
					fyne.LogError("Unable to open a reader for "+l.location.unsplash.photoCache, err)
				} else {
					unsplashBg = canvas.NewImageFromReader(reader, l.location.unsplash.photoDownloaded.String())
				}
			}
		}

		if unsplashBg == nil {
			if l.n.session == nil {
				return
			}

			backgroundPicture, err := l.n.session.getUnsplash(l.location.name, l.location.country)
			if err != nil {
				fyne.LogError("Unable to find a picture for "+l.location.name+"["+l.location.country+"]", err)
				return
			}

			unsplashBg = l.n.session.getUnsplashImage(backgroundPicture)
			if unsplashBg == nil {
				fyne.LogError("Unable to create an image for "+l.location.name+"["+l.location.country+"]: "+backgroundPicture.photoCache, err)
				return
			}
			l.location.unsplash = backgroundPicture
			l.n.store.save()
		}

		c.Objects[0] = unsplashBg
		c.Refresh()
	}()

	return widget.NewSimpleRenderer(c)
}

func listTimes() (times []string) {
	for hour := 0; hour < 24; hour++ {
		times = append(times,
			fmt.Sprintf("%02d:00", hour), fmt.Sprintf("%02d:15", hour),
			fmt.Sprintf("%02d:30", hour), fmt.Sprintf("%02d:45", hour))
	}
	return times
}
