package main

import (
	"image/color"
	"time"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/canvas"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/widget"
)

func (n *nomad) makeSplash() fyne.CanvasObject {
	return container.NewMax(
		canvas.NewRectangle(&color.NRGBA{0xFF, 0x85, 0x00, 0xFF}),
		container.NewCenter(widget.NewLabel("NOMAD")))
}

func (n *nomad) fadeSplash(obj fyne.CanvasObject) {
	time.Sleep(time.Second)
	obj.Hide()
	n.main.Content().Refresh()
}
