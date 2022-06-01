package main

import (
	"fmt"
	"image/color"
	"os"
	"strings"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/canvas"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/layout"
	"fyne.io/fyne/v2/theme"
	"fyne.io/fyne/v2/widget"
)

var (
	locationTextColor = color.NRGBA{0xFF, 0xFF, 0xFF, 0xBF}
)

type location struct {
	widget.BaseWidget
	location *city
	session  *unsplashSession

	date   *widget.Select
	time   *widget.SelectEntry
	button *widget.Button
	dots   *fyne.Container
}

func newLocation(loc *city, session *unsplashSession, n *nomad, homeContainer *fyne.Container) *location {
	l := &location{location: loc, session: session}
	l.ExtendBaseWidget(l)

	l.date = widget.NewSelect([]string{}, func(string) {})
	l.date.PlaceHolder = loc.localTime.Format("Mon 02 Jan")
	l.time = widget.NewSelectEntry(listTimes())
	l.time.PlaceHolder = "22:00" // longest
	l.time.Wrapping = fyne.TextWrapOff
	l.time.SetText(loc.localTime.Format("15:04"))

	menu := fyne.NewMenu("",
		fyne.NewMenuItem("Delete Place", func() {
			for i := 0; i < len(n.store.list); i++ {
				if l.location == n.store.list[i] {

					n.store.list = append(n.store.list[:i], n.store.list[i+1:]...)
					n.store.save()

					for j := 0; j < len(homeContainer.Objects)-1; j++ {
						if l.location.name == homeContainer.Objects[j].(*location).location.name {
							homeContainer.Objects = append(homeContainer.Objects[:j], homeContainer.Objects[j+1:]...)
							break
						}
					}

					imageLocation := session.storage.RootURI().String() + "/" + loc.unsplash.cache
					//session.storage.RootURI() gives file location prefixed with file://
					e := os.Remove(strings.Split(imageLocation, "//")[1])
					if e != nil {
						fyne.LogError("Image could not be deleted from cache", e)
					}
				}
			}
		}),
		fyne.NewMenuItem("Photo info", func() { fmt.Println("Photo info") }))

	l.button = widget.NewButtonWithIcon("", theme.MoreHorizontalIcon(), func() {
		position := fyne.CurrentApp().Driver().AbsolutePositionForObject(l.button)
		position.Y += l.button.Size().Height

		widget.ShowPopUpMenuAtPosition(menu, fyne.CurrentApp().Driver().CanvasForObject(l.button), position)
	})
	l.button.Importance = widget.LowImportance

	l.dots = container.NewVBox(layout.NewSpacer(), l.button)

	return l
}

func (l *location) CreateRenderer() fyne.WidgetRenderer {
	bg := canvas.NewImageFromResource(theme.FileImageIcon())
	bg.Translucency = 0.5
	city := widget.NewRichTextFromMarkdown("# " + strings.ToUpper(l.location.name))

	location := canvas.NewText(" "+strings.ToUpper(l.location.country)+" · "+l.location.localTime.Format("MST"), locationTextColor)
	location.TextStyle.Monospace = true
	location.TextSize = 10
	location.Move(fyne.NewPos(theme.Padding(), city.MinSize().Height-location.TextSize*.5))
	input := container.NewBorder(nil, nil, l.date, l.time)

	c := container.NewMax(bg,
		container.NewBorder(nil,
			container.NewVBox(container.NewHBox(container.NewWithoutLayout(city, location), layout.NewSpacer(), l.dots), input), nil, nil))

	go func() {
		if l.session == nil {
			return
		}

		unsplashBg, err := l.session.get(l.location)
		if err != nil {
			fyne.LogError("unable to build Unsplash image", err)
			return
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
