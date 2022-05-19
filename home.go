package main

import (
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/widget"
)

func (n *nomad) makeHome() fyne.CanvasObject {
	return widget.NewLabel("Home")
}
