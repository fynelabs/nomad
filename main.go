package main

import (
	"time"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/driver/desktop"
)

type nomad struct {
	main    fyne.Window
	store   *cityStore
	session *unsplashSession
}

func main() {
	a := app.NewWithID("com.fynelabs.nomad")
	w := a.NewWindow("Nomad")

	store := newCityStore(a.Preferences())
	session := newUnsplashSession(a.Storage(), store)
	ui := &nomad{main: w, store: store, session: session}

	var _ fyne.Theme = (*myTheme)(nil)
	a.Settings().SetTheme(&myTheme{})

	splash := makeSplash()
	w.SetContent(container.NewMax(ui.makeHome(), splash))
	w.SetPadded(false)
	w.Resize(fyne.NewSize(300, 500))
	w.SetIcon(resourceIconPng)

	if deskApp, ok := a.(desktop.App); ok {
		w.SetCloseIntercept(func() {
			w.Hide() // don't close the window if system tray used
		})

		setupSystray(deskApp, w, store)
	}
	go fadeSplash(splash)
	w.ShowAndRun()
}

func setupSystray(a desktop.App, w fyne.Window, store *cityStore) {
	a.SetSystemTrayIcon(resourceIconPng)
	setupSystrayMenu(a, w, store)

	go func() {
		for range time.NewTicker(time.Minute).C {
			setupSystrayMenu(a, w, store)
		}
	}()
}

func setupSystrayMenu(a desktop.App, w fyne.Window, store *cityStore) {
	times := make([]*fyne.MenuItem, len(store.list)+2)
	times[0] = fyne.NewMenuItem("Show Nomad", func() {
		w.Show()
	})
	times[1] = fyne.NewMenuItemSeparator()

	for i, item := range store.list {
		locDate := time.Now().In(item.localTime.Location())
		local := locDate.Format("15:04")
		label := item.name + ": " + local

		localTime := fyne.NewMenuItem(label, nil)
		localTime.Disabled = true
		times[i+2] = localTime
	}

	a.SetSystemTrayMenu(fyne.NewMenu("Times", times...))
}
