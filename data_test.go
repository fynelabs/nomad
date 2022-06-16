package main

import (
	"testing"
	"time"

	"fyne.io/fyne/v2/test"

	"github.com/stretchr/testify/assert"
)

func TestNewCityStore(t *testing.T) {
	a := test.NewApp()
	s := newCityStore(a.Preferences())

	assert.Equal(t, 1, len(s.list))
}

func TestCityStore_Add(t *testing.T) {
	a := test.NewApp()
	s := newCityStore(a.Preferences())

	s.add(newTestCity())
	assert.Equal(t, 2, len(s.list))

	s = newCityStore(a.Preferences())
	assert.Equal(t, 2, len(s.list))
}

func TestCityStore_Remove(t *testing.T) {
	a := test.NewApp()
	s := newCityStore(a.Preferences())

	s.remove(0)
	assert.Equal(t, 0, len(s.list))

	// reset to always have 1
	s = newCityStore(a.Preferences())
	assert.Equal(t, 1, len(s.list))
}

func TestNewCity(t *testing.T) {
	c := newTestCity()
	assert.Equal(t, time.Now().Location(), c.localTime.Location())
}

func newTestCity() *city {
	return newCity("Local", "Here", photo{}, time.Local)
}
