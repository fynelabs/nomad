package main

import (
	"image/color"
	"time"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/canvas"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/theme"
)

//nomadTheme.go?//
type myTheme struct{}

func (m myTheme) Color(name fyne.ThemeColorName, variant fyne.ThemeVariant) color.Color {
	return theme.DefaultTheme().Color(name, variant)
}

func (m myTheme) Font(style fyne.TextStyle) fyne.Resource {

	return resourceOpenSansExtraBoldItalicTtf
}

func (m myTheme) Icon(name fyne.ThemeIconName) fyne.Resource {
	return theme.DefaultTheme().Icon(name)
}

func (m myTheme) Size(name fyne.ThemeSizeName) float32 {
	return theme.DefaultTheme().Size(name)
}

//end of nomadTheme.go

//unsplash - token from personal account
// h, err := http.Get("https://api.unsplash.com/photos/random/?client_id=QNaGi11Ne9AcixHFCnUj4XsE6TFcpJaDT0ppweJNRdU")
// if err != nil {
// }
// fmt.Println(h)

func (n *nomad) makeSplash(app fyne.App) fyne.CanvasObject {

	text := canvas.NewText("NOMAD", color.White)
	text.TextSize = 70

	var _ fyne.Theme = (*myTheme)(nil)
	app.Settings().SetTheme(&myTheme{})

	globe := canvas.NewImageFromFile("static/images/globe.png")
	globe.FillMode = canvas.ImageFillOriginal

	vBox := container.NewVBox(
		container.NewCenter(globe),
		container.NewCenter(text),
	)

	return container.NewMax(

		// canvas.NewImageFromFile("static/images/splashPlaceholder.png"),
		container.NewCenter(vBox),
	)
}

func (n *nomad) fadeSplash(obj fyne.CanvasObject) {
	time.Sleep(time.Second * 10)
	obj.Hide()
	n.main.Content().Refresh()
}
