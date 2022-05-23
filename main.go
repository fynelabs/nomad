package main

import (
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/container"
)

type nomad struct {
	main fyne.Window
}

func main() {
	a := app.NewWithID("com.fynelabs.nomad")
	w := a.NewWindow("Nomad")
	ui := &nomad{main: w}

	splash := ui.makeSplash()
	w.SetContent(container.NewMax(ui.makeHome(), splash))
	w.SetPadded(false)
	w.Resize(fyne.NewSize(300, 500))

	go ui.fadeSplash(splash)
	w.ShowAndRun()
}
