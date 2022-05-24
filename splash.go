package main

import (
	"fmt"
	"image/color"
	"time"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/canvas"
	"fyne.io/fyne/v2/container"
	xWidget "fyne.io/x/fyne/widget"
)

func (n *nomad) makeSplash() fyne.CanvasObject {

	text := canvas.NewText("NOMAD", color.White)
	text.TextSize = 50
	text.TextStyle = fyne.TextStyle{Italic: true, Bold: true}

	gif, err := xWidget.NewAnimatedGifFromResource(resourceGlobeSpinnerGif)
	gif.SetMinSize(fyne.NewSize(50, 50))
	gif.Start()

	if err != nil {
		fmt.Println(err)
	}

	vBox := container.NewVBox(
		container.NewCenter(gif),
		container.NewCenter(text),
	)

	return container.NewMax(
		canvas.NewRectangle(&color.NRGBA{0x18, 0x0C, 0x27, 0xFF}),
		container.NewCenter(vBox),
	)
}

func (n *nomad) fadeSplash(obj fyne.CanvasObject) {
	time.Sleep(time.Second * 2)
	obj.Hide()
	n.main.Content().Refresh()
}
