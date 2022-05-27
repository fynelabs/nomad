package main

import (
	"image/color"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/canvas"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/theme"
	"fyne.io/fyne/v2/widget"
)

var _ fyne.Widget = (*iconWithPopUpMenu)(nil)
var _ fyne.Tappable = (*iconWithPopUpMenu)(nil)

// handles the tap event.
type iconWithPopUpMenu struct {
	widget.BaseWidget
	container *fyne.Container
	bg        *canvas.Rectangle
	menu      *fyne.Menu
	animation *fyne.Animation
}

func newIconWithPopUpMenu(res fyne.Resource, menu *fyne.Menu) *iconWithPopUpMenu {
	r := &iconWithPopUpMenu{
		bg:   canvas.NewRectangle(color.Transparent),
		menu: menu,
	}
	r.ExtendBaseWidget(r)
	r.animation = fyne.NewAnimation(canvas.DurationStandard, func(done float32) {
		if done >= 1 {
			r.bg.FillColor = color.Transparent
			canvas.Refresh(r.bg)
			return
		}

		mid := r.Size().Width / 2
		size := mid * done
		r.bg.Resize(fyne.NewSize(size*2, r.Size().Height))
		r.bg.Move(fyne.NewPos(mid-size, 0))

		red, green, blue, alpha := theme.PressedColor().RGBA()
		fade := uint8(alpha) - uint8(float32(uint8(alpha))*done)
		r.bg.FillColor = &color.NRGBA{R: uint8(red), G: uint8(green), B: uint8(blue), A: fade}
		canvas.Refresh(r.bg)
	})
	r.animation.Curve = fyne.AnimationEaseOut
	r.container = container.NewMax(r.bg, widget.NewIcon(res))

	return r
}

func (t *iconWithPopUpMenu) CreateRenderer() fyne.WidgetRenderer {
	return widget.NewSimpleRenderer(t.container)
}

func (t *iconWithPopUpMenu) Tapped(e *fyne.PointEvent) {
	t.animation.Stop()
	t.animation.Start()

	position := fyne.CurrentApp().Driver().AbsolutePositionForObject(t)
	position.Y += t.Size().Height

	widget.ShowPopUpMenuAtPosition(t.menu, fyne.CurrentApp().Driver().CanvasForObject(t), position)
}
