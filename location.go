package main

import (
	"fmt"
	"image/color"
	"strings"
	"time"

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

	time   *widget.SelectEntry
	button *widget.Button
	dots   *fyne.Container

	dateButton      *widget.Button
	locationTZLabel *canvas.Text

	calendar      *calendar
	homeContainer *fyne.Container
}

func newLocation(loc *city, n *nomad, homeC *fyne.Container) *location {

	l := &location{location: loc, session: n.session, homeContainer: homeC}
	l.ExtendBaseWidget(l)

	l.time = widget.NewSelectEntry(listTimes())
	l.time.PlaceHolder = "22:00" // longest
	l.time.Wrapping = fyne.TextWrapOff
	l.time.SetText(loc.localTime.Format("15:04"))
	l.time.OnChanged = func(s string) {
		var hour, minute int
		if s == "Now" {
			globalAppTime = time.Now()
			hour = time.Now().Hour()
			minute = time.Now().Minute()
			currentTimeSelected = true
		} else {
			fmt.Sscanf(s, "%d:%d", &hour, &minute)
			currentTimeSelected = false
		}
		localOld := globalAppTime.In(l.location.localTime.Location())
		selectedDate := time.Date(localOld.Year(), localOld.Month(), localOld.Day(), hour, minute, 0, 0, l.location.localTime.Location())

		setDate(selectedDate, l.homeContainer.Objects)
	}

	menu := fyne.NewMenu("",
		fyne.NewMenuItem("Delete Place", func() { l.remove(homeC, n) }),
		fyne.NewMenuItem("Photo info", func() {
			c := fyne.CurrentApp().Driver().CanvasForObject(l.button)
			c.Overlays().Add(loc.newInfoScreen(c))
		}))

	l.button = widget.NewButtonWithIcon("", theme.MoreHorizontalIcon(), func() {
		position := fyne.CurrentApp().Driver().AbsolutePositionForObject(l.button)
		position.Y += l.button.Size().Height

		widget.ShowPopUpMenuAtPosition(menu, fyne.CurrentApp().Driver().CanvasForObject(l.button), position)
	})
	l.button.Importance = widget.LowImportance

	l.dots = container.NewVBox(layout.NewSpacer(), l.button)

	l.calendar = newCalendar(loc.localTime, func(t time.Time) {
		currentTimeSelected = false
		setDate(t, l.homeContainer.Objects)
	})

	l.dateButton = widget.NewButtonWithIcon(l.calendar.fullDate(), theme.MenuDropDownIcon(), func() {
		l.calendar.newCalendarPopUpAtPos(n.main.Canvas(), fyne.NewPos(0, l.Size().Height))
	})
	l.dateButton.Alignment = widget.ButtonAlignLeading
	l.dateButton.IconPlacement = widget.ButtonIconTrailingText
	l.dateButton.Importance = widget.LowImportance

	return l
}

func (l *location) CreateRenderer() fyne.WidgetRenderer {
	bg := canvas.NewImageFromResource(theme.FileImageIcon())
	bg.Translucency = 0.5
	city := widget.NewRichTextFromMarkdown("# " + l.location.name)
	l.locationTZLabel = canvas.NewText(strings.ToUpper(l.location.country)+" · "+l.location.localTime.Format("MST"), locationTextColor)
	l.locationTZLabel.TextStyle.Monospace = true
	l.locationTZLabel.TextSize = 10
	l.locationTZLabel.Move(fyne.NewPos(theme.Padding()*2, 40))
	input := container.NewBorder(nil, nil, l.dateButton, l.time)

	c := container.NewMax(bg,
		container.NewBorder(nil,
			container.NewVBox(container.NewHBox(container.NewWithoutLayout(city, l.locationTZLabel), layout.NewSpacer(), l.dots), input), nil, nil))

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
	times = append(times, "Now")
	for hour := 0; hour < 24; hour++ {
		times = append(times,
			fmt.Sprintf("%02d:00", hour), fmt.Sprintf("%02d:15", hour),
			fmt.Sprintf("%02d:30", hour), fmt.Sprintf("%02d:45", hour))
	}
	return times
}

func (l *location) updateLocation(locDate time.Time) {
	l.time.Text = locDate.Format("15:04")
	l.time.Refresh()
	l.locationTZLabel.Text = strings.ToUpper(l.location.country + " · " + locDate.Format("MST"))
	l.locationTZLabel.Refresh()
	l.dateButton.SetText(locDate.Format("Mon 02 Jan 2006"))
}

func (l *location) remove(homeContainer *fyne.Container, n *nomad) {
	for i := 0; i < len(n.store.list); i++ {
		if l.location == n.store.list[i] {

			n.store.remove(i)
			l.removeLocationFromContainer(homeContainer)

			l.session.removeImageFromCache(l)

			break
		}
	}
}

func (l *location) removeLocationFromContainer(homeContainer *fyne.Container) {
	for j := 0; j < len(homeContainer.Objects)-1; j++ {
		if l.location.name == homeContainer.Objects[j].(*location).location.name {
			homeContainer.Remove(homeContainer.Objects[j])
			break
		}
	}
}
