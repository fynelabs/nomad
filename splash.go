package main

import (
	"fmt"
	"image/color"
	"time"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/canvas"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/storage"
	xWidget "fyne.io/x/fyne/widget"
)

func (n *nomad) makeSplash(app fyne.App) fyne.CanvasObject {

	text := canvas.NewText("NOMAD", color.White)
	text.TextSize = 50

	var _ fyne.Theme = (*myTheme)(nil)
	app.Settings().SetTheme(&myTheme{})

	gif, err := xWidget.NewAnimatedGif(storage.NewFileURI("./static/images/globeSpinner.gif"))
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
		container.NewCenter(vBox),
	)
}

func (n *nomad) fadeSplash(obj fyne.CanvasObject, splashFadeTime int) {
	time.Sleep(time.Second * time.Duration(splashFadeTime))
	obj.Hide()
	n.main.Content().Refresh()
}
