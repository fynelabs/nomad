package main

import (
	"testing"
	"time"

	"fyne.io/fyne/v2/driver/software"

	"github.com/stretchr/testify/assert"
)

func TestMakeSplash(t *testing.T) {
	s := makeSplash()
	gif := gifInSplash(s)
	assert.True(t, gif.Visible())

	// test it is running
	img1 := software.Render(gif, nil)
	time.Sleep(time.Second / 5)
	img2 := software.Render(gif, nil)
	assert.NotEqual(t, img1, img2)
}

func TestFadeSplash(t *testing.T) {
	s := makeSplash()
	gif := gifInSplash(s)

	fadeSplash(s)
	assert.False(t, s.Visible())

	// test it is stopped
	img1 := software.Render(gif, nil)
	time.Sleep(time.Second / 5)
	img2 := software.Render(gif, nil)
	assert.Equal(t, img1, img2)
}
