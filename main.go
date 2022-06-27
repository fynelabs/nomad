package main

import (
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/container"
)

type nomad struct {
	main    fyne.Window
	store   *cityStore
	session *unsplashSession
}

func main() {
	a := app.NewWithID("com.fynelabs.nomad")
	w := a.NewWindow("Nomad")

	selfManage(a, w)

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

	go fadeSplash(splash)
	w.ShowAndRun()
}
