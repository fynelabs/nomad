package main

import (
	"fmt"
	"image/color"
	"strings"
	"time"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/canvas"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/driver/desktop"
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
	app      *nomad

	button *widget.Button
	dots   *fyne.Container

	dateButton      *widget.Button
	timeButton      *widget.Button
	locationTZLabel *canvas.Text

	calendar      *calendar
	homeContainer *fyne.Container
}

func newLocation(loc *city, n *nomad, homeC *fyne.Container) *location {

	l := &location{location: loc, session: n.session, app: n, homeContainer: homeC}
	l.ExtendBaseWidget(l)

	menu := fyne.NewMenu("",
		fyne.NewMenuItem("Delete Place", func() { l.remove(homeC) }),
		fyne.NewMenuItem("Photo info", func() {
			c := fyne.CurrentApp().Driver().CanvasForObject(l.button)
			info := loc.newInfoScreen(c)
			info.Resize(c.Size())
			c.Overlays().Add(info)
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
		position := fyne.CurrentApp().Driver().AbsolutePositionForObject(l.dateButton)
		position.Y += l.button.Size().Height
		l.calendar.showAtPos(n.main.Canvas(), position)
	})
	l.dateButton.Alignment = widget.ButtonAlignLeading
	l.dateButton.IconPlacement = widget.ButtonIconTrailingText
	l.dateButton.Importance = widget.LowImportance

	l.timeButton = widget.NewButtonWithIcon(loc.localTime.Format("15:04"), theme.MenuDropDownIcon(), func() {
		position := fyne.CurrentApp().Driver().AbsolutePositionForObject(l.timeButton)
		position.Y += l.timeButton.Size().Height
		sizedMenuWidth := float32(104)
		position.X += l.timeButton.Size().Width - sizedMenuWidth - theme.Padding()*2

		times := listTimes()
		menuItems := []*fyne.MenuItem{}
		for _, t := range times {
			v := t
			menuItems = append(menuItems, fyne.NewMenuItem("       "+t, func() {
				l.onTimeSelect(v)
				n.main.Canvas().Overlays().Top().Hide()
			}))
		}
		t := fyne.NewMenu("Times", menuItems...)
		m := newSizedMenu(t, fyne.NewSize(sizedMenuWidth+theme.Padding()*2, minHeight))
		widget.ShowPopUpAtPosition(m, n.main.Canvas(), position)

	})
	l.timeButton.Alignment = widget.ButtonAlignLeading
	l.timeButton.IconPlacement = widget.ButtonIconTrailingText
	l.timeButton.Importance = widget.LowImportance

	return l
}

func (l *location) onTimeSelect(t string) {
	var hour, minute int
	if t == "Now" {
		globalAppTime = time.Now()
		hour = time.Now().Hour()
		minute = time.Now().Minute()
		currentTimeSelected = true
	} else {
		fmt.Sscanf(t, "%d:%d", &hour, &minute)
		currentTimeSelected = false
	}
	localOld := globalAppTime.In(l.location.localTime.Location())
	selectedDate := time.Date(localOld.Year(), localOld.Month(), localOld.Day(), hour, minute, 0, 0, l.location.localTime.Location())

	setDate(selectedDate, l.homeContainer.Objects)
}

func (l *location) CreateRenderer() fyne.WidgetRenderer {
	bg := canvas.NewImageFromResource(theme.FileImageIcon())
	op := canvas.NewRectangle(color.NRGBA{0x00, 0x00, 0x00, 0x59})
	bg.Translucency = 0.5
	city := widget.NewRichTextFromMarkdown("# " + strings.ToUpper(l.location.name))
	city.Move(fyne.NewPos(3, 0))
	l.locationTZLabel = canvas.NewText(strings.ToUpper(l.location.country)+" · "+l.location.localTime.Format("MST"), locationTextColor)
	l.locationTZLabel.TextStyle.Monospace = true
	l.locationTZLabel.TextSize = 10
	l.locationTZLabel.Move(fyne.NewPos(12, 40))

	input := container.NewBorder(nil, nil, l.dateButton, l.timeButton)

	c := container.NewMax(bg, op,
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
	l.timeButton.Text = locDate.Format("15:04")
	l.timeButton.Refresh()
	l.locationTZLabel.Text = strings.ToUpper(l.location.country + " · " + locDate.Format("MST"))
	l.locationTZLabel.Refresh()
	l.dateButton.SetText(locDate.Format("Mon 02 Jan 2006"))
}

func (l *location) remove(homeContainer *fyne.Container) {
	for i := 0; i < len(l.app.store.list); i++ {
		if l.location == l.app.store.list[i] {

			l.app.store.remove(i)
			l.removeLocationFromContainer(homeContainer)

			l.session.removeImageFromCache(l)
			l.updateMenu()

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

func (l *location) updateMenu() {
	if deskApp, ok := fyne.CurrentApp().(desktop.App); ok {
		setupSystrayMenu(deskApp, l.app.main, l.app.store)
	}
}
